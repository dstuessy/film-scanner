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
)

const boundaryWord = "MJPEGBOUNDARY"

var frameInterval = 50 * time.Millisecond

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	_, err := auth.CheckToken(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", boundaryWord))
	w.Header().Set("Cache-Control", "no-cache")

	for {
		time.Sleep(frameInterval)

		img, err := camera.CaptureFrame()
		if err != nil {
			log.Println(err)
		}

		header := strings.Join([]string{
			fmt.Sprintf("\r\n--%s", boundaryWord),
			"Content-Type: image/jpeg",
			fmt.Sprintf("Content-Length: %d", len(img)),
			"X-Timestamp: 0.000000",
			"\r\n",
		}, "\r\n")

		frame := make([]byte, len(header)+len(img))

		copy(frame, header)
		copy(frame[len(header):], img)

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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	srv, err := drive.GetDriveFileService(token, drive.GetContext())
	if err != nil {
		log.Fatal(err)
	}

	dir, err := drive.FindFolder(srv, drive.DriveDirName)
	if err != nil {
		log.Fatal(err)
	}

	img, err := camera.CaptureFrame()
	if err != nil {
		log.Println(err)
	}

	name := fmt.Sprintf("image-%d.jpg", time.Now().Unix())
	if _, err := drive.SaveImage(srv, img, name, dir.Id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
