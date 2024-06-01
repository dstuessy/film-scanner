package camera

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
	// "gocv.io/x/gocv/contrib"
)

var webcam *gocv.VideoCapture

var FrameInterval = 60 * time.Millisecond

type ImageData struct {
	Rows int
	Cols int
	Data []byte
}

func IsCameraOpen() bool {
	if webcam == nil {
		return false
	}

	return webcam.IsOpened()
}

func OpenCamera() error {
	if webcam != nil {
		return nil
	}

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
		webcam = nil
		return err
	}
	webcam = nil
	return nil
}

func CaptureStill() ([]byte, error) {
	if webcam != nil {
		return nil, errors.New("Camera is still open for streaming")
	}

	if os.Getenv("STILL_IMG_COMMAND") == "" {
		return nil, errors.New("STILL_IMG_COMMAND not set")
	}

	imgName := fmt.Sprintf(os.Getenv("STILL_IMG_NAME"), time.Now().Unix())
	imgLoc := fmt.Sprintf("%s/%s", tmpdir, imgName)

	imgCmd := fmt.Sprintf(os.Getenv("STILL_IMG_COMMAND"), imgLoc)
	slicedCmd := strings.Split(imgCmd, " ")

	log.Println("Capturing still image with command:", imgCmd)

	cmd := exec.Command(slicedCmd[0], slicedCmd[1:]...)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	log.Println("Captured still image", fmt.Sprintf("%s", output))

	mat := gocv.IMRead(imgLoc, gocv.IMReadColor)
	defer mat.Close()

	frame := mat

	outerFrame, err := CropFilm(mat, 0.6, 0.7, 0.1, false)
	if err == nil {
		log.Println("film:", err)
		frame = outerFrame
	}

	film, err := CropFilm(frame, 0.7, 0.9, 0, true)
	if err == nil {
		log.Println("film:", err)
		frame = film
	}

	wb := contrib.NewSimpleWB()
	wb.BalanceWhite(frame, &frame)

	gocv.BitwiseNot(frame, &frame)

	jpeg, err := gocv.IMEncode(".jpg", frame)
	if err != nil {
		return nil, err
	}

	return jpeg.GetBytes(), nil
}
