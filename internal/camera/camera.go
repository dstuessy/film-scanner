package camera

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

var webcam *gocv.VideoCapture

var stream = make(chan ImageData)

var FrameInterval = 60 * time.Millisecond

var tmpdir string

type ImageData struct {
	Rows int
	Cols int
	Data []byte
}

func SetupTempDir() error {
	dir, err := os.MkdirTemp("", os.Getenv("STILL_IMG_DIR"))
	if err != nil {
		log.Println("Error creating temp dir:", err)
		return err
	}
	log.Println("Created temp dir:", dir)

	tmpdir = dir
	return nil
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

	gocv.BitwiseNot(mat, &mat)

	wb := contrib.NewSimpleWB()

	wb.BalanceWhite(mat, &mat)

	jpeg, err := gocv.IMEncode(".jpg", mat)
	if err != nil {
		return nil, err
	}

	return jpeg.GetBytes(), nil
}
