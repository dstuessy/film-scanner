package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/camera"
	"github.com/dstuessy/film-scanner/internal/drive"
	// "github.com/dstuessy/film-scanner/internal/tiff"
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

		if len(img.Data) == 0 {
			log.Println("Stream Empty")
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		smallImg, err := camera.ResizeData(img, 0.5)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		jpeg, err := camera.EncodeJpeg(smallImg)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		header := strings.Join([]string{
			fmt.Sprintf("\r\n--%s", boundaryWord),
			"Content-Type: image/jpeg",
			fmt.Sprintf("Content-Length: %d", len(jpeg)),
			"X-Timestamp: 0.000000",
			"\r\n",
		}, "\r\n")

		frame := make([]byte, len(header)+len(jpeg))

		copy(frame, header)
		copy(frame[len(header):], jpeg)

		if _, err := w.Write(frame); err != nil {
			log.Println(err)
			break
		}
	}

	log.Println("Stream disconnected")
	return
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

	if err := camera.CloseCamera(); err != nil {
		log.Println(err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if camera.IsCameraOpen() {
			return
		}

		if err := camera.OpenCamera(); err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
	}()

	jpeg, err := camera.CaptureStill()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	name := fmt.Sprintf("image-%d.jpeg", time.Now().Unix())
	if _, err := drive.SaveImage(srv, jpeg, name, projectId[0]); err != nil {
		log.Println(err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	return
}
