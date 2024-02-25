package camera

import (
	"errors"
	"log"
	"time"

	"gocv.io/x/gocv"
)

var webcam *gocv.VideoCapture

var stream = make(chan []byte)

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

func GetStream() chan []byte {
	return stream
}

func captureFrame() ([]byte, error) {
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
