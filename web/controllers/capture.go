package controllers

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/camera"
	"github.com/dstuessy/film-scanner/internal/drive"
	"github.com/dstuessy/film-scanner/internal/tiff"
	"gocv.io/x/gocv"
)

const boundaryWord = "MJPEGBOUNDARY"

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	_, err := auth.CheckToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundaryWord))
	w.Header().Set("Cache-Control", "no-cache")

	for {
		time.Sleep(camera.FrameInterval)

		img := <-camera.GetStream()
		smallImg := gocv.NewMat()

		gocv.Resize(img, &smallImg, image.Point{}, 0.5, 0.5, gocv.InterpolationArea)

		jpeg, err := gocv.IMEncode(".jpg", smallImg)
		smallImg.Close()
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		jpegBytes := jpeg.GetBytes()

		header := strings.Join([]string{
			fmt.Sprintf("\r\n--%s", boundaryWord),
			"Content-Type: image/jpeg",
			fmt.Sprintf("Content-Length: %d", len(jpegBytes)),
			"X-Timestamp: 0.000000",
			"\r\n",
		}, "\r\n")

		frame := make([]byte, len(header)+len(jpegBytes))

		copy(frame, header)
		copy(frame[len(header):], jpegBytes)

		if _, err := w.Write(frame); err != nil {
			log.Println(err)
			break
		}
	}

	log.Println("Stream disconnected")
}

func CaptureScanHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.CheckToken(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectId := r.URL.Query()["project"]
	if len(projectId) == 0 {
		log.Println(fmt.Sprintf("Project id not found in URL: %s", r.URL.String()))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	srv, err := drive.GetDriveFileService(token, drive.GetContext())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	img := <-camera.GetStream()

	tiff, err := tiff.EncodeTiff(img)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	name := fmt.Sprintf("image-%d.tiff", time.Now().Unix())
	if _, err := drive.SaveImage(srv, tiff, name, projectId[0]); err != nil {
		log.Println(err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}
