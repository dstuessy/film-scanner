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

	// imgCmd := fmt.Sprintf(os.Getenv("STILL_IMG_COMMAND"), imgLoc)
	imgCmds := strings.ReplaceAll(os.Getenv("STILL_IMG_COMMAND"), "{image}", imgLoc)
	for _, imgCmd := range strings.Split(imgCmds, ";") {
		slicedCmd := strings.Split(strings.Trim(imgCmd, " "), " ")

		log.Println("Capturing still image with command:", imgCmd)

		cmd := exec.Command(slicedCmd[0], slicedCmd[1:]...)

		output, err := cmd.Output()
		log.Println(output)
		if err != nil {
			log.Println("Failed to capture still image with command:", imgCmd)
			return nil, err
		}
	}

	log.Println("Captured still image")

	finalLoc := BuildFileName(imgLoc)
	tiff, err := os.ReadFile(finalLoc)
	return tiff, err
}

func BuildFileName(name string) string {
	return fmt.Sprintf("%s%s", name, os.Getenv("STILL_IMG_EXT"))
}

func GetMimeType() string {
	return os.Getenv("STILL_IMG_MIME")
}
