package camera

import (
	"errors"
	"log"
	"time"

	"gocv.io/x/gocv"
)

var webcam *gocv.VideoCapture

var stream = make(chan gocv.Mat)

var FrameInterval = 50 * time.Millisecond

func Open() error {
	c, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return err
	}

	webcam = c

	go func() {
		for {
			img, err := captureFrame()
			if err != nil {
				log.Println(err)
				log.Println("Closing stream")
				close(stream)
				break
			}

			stream <- img

			time.Sleep(FrameInterval)
		}
	}()

	return nil
}

func Close() error {
	return webcam.Close()
}

func GetStream() chan gocv.Mat {
	return stream
}

func captureFrame() (gocv.Mat, error) {
	img := gocv.NewMat()

	if ok := webcam.Read(&img); !ok {
		return img, errors.New("cannot read from webcam")
	}
	if img.Empty() {
		return img, errors.New("empty frame")
	}

	return img, nil
}
