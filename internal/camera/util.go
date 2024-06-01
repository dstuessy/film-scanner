package camera

import (
	"errors"
	"image"
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

func DataFromMat(bgr gocv.Mat) ImageData {
	rgb := gocv.NewMat()
	defer rgb.Close()
	gocv.CvtColor(bgr, &rgb, gocv.ColorBGRToRGB)

	return ImageData{
		Rows: rgb.Rows(),
		Cols: rgb.Cols(),
		Data: rgb.ToBytes(),
	}
}

func ResizeData(img ImageData, scale float64) (ImageData, error) {
	mat, err := gocv.NewMatFromBytes(img.Rows, img.Cols, gocv.MatTypeCV8UC3, img.Data)
	defer mat.Close()
	if err != nil {
		return ImageData{}, err
	}

	gocv.Resize(mat, &mat, image.Point{}, scale, scale, gocv.InterpolationArea)

	return DataFromMat(mat), nil
}

func EncodeJpeg(img ImageData) ([]byte, error) {
	mat, err := gocv.NewMatFromBytes(img.Rows, img.Cols, gocv.MatTypeCV8UC3, img.Data)
	defer mat.Close()
	if err != nil {
		return nil, err
	}

	jpeg, err := gocv.IMEncode(".jpg", mat)
	if err != nil {
		return nil, err
	}

	return jpeg.GetBytes(), nil
}

func thresholdImage(img gocv.Mat, threshold int, mask gocv.Mat) gocv.Mat {
	thresholded := gocv.NewMat()
	gocv.Threshold(img, &thresholded, float32(threshold), 255, gocv.ThresholdBinary)

	gocv.BitwiseAnd(thresholded, mask, &thresholded)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	gocv.MorphologyExWithParams(thresholded, &thresholded, gocv.MorphClose, kernel, 3, 0)

	return thresholded
}

func createIgnoreMask(img gocv.Mat) gocv.Mat {
	mat := gocv.NewMat()
	defer mat.Close()
	gocv.CvtColor(img, &mat, gocv.ColorBGRToGray)

	gocv.FastNlMeansDenoising(mat, &mat)

	gocv.EqualizeHist(mat, &mat)

	gocv.Threshold(mat, &mat, 250, 255, gocv.ThresholdBinary)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	brightnessMask := gocv.NewMat()
	defer brightnessMask.Close()
	gocv.MorphologyExWithParams(mat, &brightnessMask, gocv.MorphClose, kernel, 3, 0)

	hsv := gocv.NewMat()
	defer hsv.Close()
	gocv.CvtColor(img, &hsv, gocv.ColorBGRToHSV)
	gocv.GaussianBlur(hsv, &hsv, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	saturationMask := gocv.NewMat()
	defer saturationMask.Close()
	gocv.InRangeWithScalar(hsv, gocv.NewScalar(0, 0, 0, 255), gocv.NewScalar(255, 7, 255, 255), &saturationMask)

	ignoreMask := gocv.NewMat()
	gocv.BitwiseOr(brightnessMask, saturationMask, &ignoreMask)
	gocv.BitwiseNot(ignoreMask, &ignoreMask)

	return ignoreMask
}

func findLargestContourRect(img gocv.Mat) gocv.RotatedRect {
	var adequateRect gocv.RotatedRect
	var largestArea float64

	contours := gocv.FindContours(img, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	for i := 0; i < contours.Size(); i++ {
		contour := contours.At(i)
		area := gocv.ContourArea(contour)

		if area > largestArea {
			adequateRect = gocv.MinAreaRect(contour)
		}
	}

	return adequateRect
}

func getSmallestRect(rects []gocv.RotatedRect) gocv.RotatedRect {
	var smallestRect gocv.RotatedRect
	var smallestArea float64 = 0

	for _, r := range rects {
		area := float64(r.Width) * float64(r.Height)
		if smallestArea == 0 || area < smallestArea {
			smallestRect = r
			smallestArea = area
		}
	}

	return smallestRect
}

func CropFilm(img gocv.Mat, minCropRatio, maxCropRatio, trimming float64, debug bool) (gocv.Mat, error) {
	ignoreMask := createIgnoreMask(img)
	defer ignoreMask.Close()

	imgHeight := float64(img.Rows())
	imgWidth := imgHeight * 1.5

	log.Print("image width, height", imgWidth, imgHeight)

	minCropWidth := int(minCropRatio * imgWidth)
	minCropHeight := int(minCropRatio * imgHeight)
	maxCropWidth := int(maxCropRatio * imgWidth)
	maxCropHeight := int(maxCropRatio * imgHeight)

	log.Println("finding edges:")

	cropRects := make([]gocv.RotatedRect, 0)

	for threshold := 0; threshold <= 250; threshold += 5 {
		gray := gocv.NewMat()
		gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
		t := thresholdImage(gray, threshold, ignoreMask)
		defer t.Close()

		log.Println("------------------------")
		log.Println("| threshold:", threshold)

		r := findLargestContourRect(t)

		log.Println("| width, height:", r.Width, r.Height)

		if debug {
			log.Println("drawing rectangle", r.BoundingRect)
			gocv.Rectangle(&img, r.BoundingRect, color.RGBA{0, 255, 0, 1}, 2)
		}

		if r.Width < minCropWidth || r.Height < minCropHeight || r.Width > maxCropWidth || r.Height > maxCropHeight {
			continue
		}

		log.Println("| it fits!")
		log.Println("| minCropWidth, maxCropWidth:", minCropWidth, maxCropWidth)
		log.Println("| r:", r)

		cropRects = append(cropRects, r)
	}

	if debug {
		centerX := int(img.Cols() / 2)
		centerY := int(img.Rows() / 2)

		// drawing the expected minimum crop area
		points := []image.Point{
			{X: centerX - minCropWidth/2, Y: centerY - minCropHeight/2},
			{X: centerX + minCropWidth/2, Y: centerY - minCropHeight/2},
			{X: centerX + minCropWidth/2, Y: centerY + minCropHeight/2},
			{X: centerX - minCropWidth/2, Y: centerY + minCropHeight/2},
		}
		minContour := gocv.NewPointVectorFromPoints(points)
		cs := gocv.NewPointsVector()
		cs.Append(minContour)
		gocv.DrawContours(&img, cs, -1, color.RGBA{255, 0, 0, 1}, 2)

		// drawing the expected maximum crop area
		points = []image.Point{
			{X: centerX - maxCropWidth/2, Y: centerY - maxCropHeight/2},
			{X: centerX + maxCropWidth/2, Y: centerY - maxCropHeight/2},
			{X: centerX + maxCropWidth/2, Y: centerY + maxCropHeight/2},
			{X: centerX - maxCropWidth/2, Y: centerY + maxCropHeight/2},
		}
		maxContour := gocv.NewPointVectorFromPoints(points)
		cs = gocv.NewPointsVector()
		cs.Append(maxContour)
		gocv.DrawContours(&img, cs, -1, color.RGBA{0, 0, 255, 1}, 2)
	}

	if len(cropRects) == 0 {
		temp := gocv.NewMat()
		defer temp.Close()
		return temp, errors.New("No crop found")
	}

	smallestRect := getSmallestRect(cropRects)

	cropped := img.Region(smallestRect.BoundingRect)

	return cropped, nil
}
