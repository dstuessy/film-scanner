package camera

import (
	"errors"

	"gocv.io/x/gocv"
)

var webcam *gocv.VideoCapture

func Setup() error {
	c, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return err
	}

	webcam = c

	return nil
}

func Close() error {
	return webcam.Close()
}

func CaptureFrame() ([]byte, error) {
	img := gocv.NewMat()
	defer img.Close()

	if ok := webcam.Read(&img); !ok {
		return nil, errors.New("cannot read from webcam")
	}
	if img.Empty() {
		return nil, errors.New("empty frame")
	}
	buf, err := gocv.IMEncode(".jpg", img)
	if err != nil {
		return nil, err
	}

	return buf.GetBytes(), nil
}
