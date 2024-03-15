package camera

import (
	"errors"
	"image"
	"log"
	"os"
	"strconv"
	"time"

	"gocv.io/x/gocv"
)

var webcam *gocv.VideoCapture

var stream = make(chan ImageData)

var FrameInterval = 60 * time.Millisecond

type ImageData struct {
	Rows int
	Cols int
	Data []byte
}

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

func StartStream() error {
	lastFrame := ImageData{}

	OpenCamera()

	go func() {
		for {
			if webcam == nil || !webcam.IsOpened() {
				stream <- lastFrame
				continue
			}

			img, err := captureFrame()
			if err != nil {
				log.Println(err)
				log.Println("Closing stream")
				close(stream)
				break
			}

			stream <- img
			lastFrame = img

			time.Sleep(FrameInterval)
		}
	}()

	return nil
}

func OpenCamera() error {
	c, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return err
	}

	if os.Getenv("CAM_WIDTH") != "" && os.Getenv("CAM_HEIGHT") != "" {
		w, err := strconv.ParseFloat(os.Getenv("CAM_WIDTH"), 64)
		if err != nil {
			return err
		}
		c.Set(gocv.VideoCaptureFrameWidth, w)

		h, err := strconv.ParseFloat(os.Getenv("CAM_HEIGHT"), 64)
		if err != nil {
			return err
		}
		c.Set(gocv.VideoCaptureFrameHeight, h)
	}

	webcam = c

	return nil
}

func CloseCamera() error {
	if err := webcam.Close(); err != nil {
		return err
	}
	webcam = nil
	return nil
}

func GetStream() chan ImageData {
	return stream
}

func captureFrame() (ImageData, error) {
	mat := gocv.NewMat()
	defer mat.Close()

	if ok := webcam.Read(&mat); !ok {
		return DataFromMat(mat), errors.New("cannot read from webcam")
	}
	if mat.Empty() {
		return DataFromMat(mat), errors.New("empty frame")
	}

	return DataFromMat(mat), nil
}
