package camera

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"gocv.io/x/gocv"
)

var webcam *gocv.VideoCapture

var stream = make(chan gocv.Mat)

var FrameInterval = 60 * time.Millisecond

func Open() error {
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

			img.Close()
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
