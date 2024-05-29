package camera

import (
	"errors"
	"gocv.io/x/gocv"
	"log"
	"time"
)

var stream = make(chan ImageData)

func StartStream() error {
	lastFrame := ImageData{}

	OpenCamera()

	go func() {
		for {
			if !IsCameraOpen() {
				stream <- lastFrame
				continue
			}

			img, err := captureFrame()
			if err != nil {
				log.Println(err)
				log.Println("Closing stream")
				continue
			}

			stream <- img
			lastFrame = img

			time.Sleep(FrameInterval)
		}
	}()

	return nil
}

func GetStream() chan ImageData {
	return stream
}

func captureFrame() (ImageData, error) {
	mat := gocv.NewMat()
	defer mat.Close()

	if ok := webcam.Read(&mat); !ok {
		return ImageData{}, errors.New("Cannot read from webcam")
	}
	if mat.Empty() {
		return ImageData{}, errors.New("Empty frame")
	}

	return DataFromMat(mat), nil
}
