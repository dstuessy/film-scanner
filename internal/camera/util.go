package camera

import (
	"gocv.io/x/gocv"
	"image"
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
