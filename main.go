package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/camera"
	"github.com/dstuessy/film-scanner/internal/drive"
	"github.com/dstuessy/film-scanner/internal/render"
	"github.com/joho/godotenv"
)

const driveDirName = "Open Scanner"

const boundaryWord = "MJPEGBOUNDARY"

var frameInterval time.Duration

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	auth.Setup()

	frameInterval = 50 * time.Millisecond

	if err := camera.Open(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer camera.Close()

	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/login", loginHandler)

	http.HandleFunc("/oauth2callback", authCallbackHandler)

	http.HandleFunc("/scan/new", func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.CheckToken(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		if err := render.RenderPage(w, "/new.html", nil); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/resource/workspace/create", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.CheckToken(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		fileSrv, err := drive.GetDriveFileService(token, drive.GetContext())
		if err != nil {
			if errors.Is(err, new(auth.TokenExpiredError)) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			} else {
				log.Fatal(err)
			}
		}

		folder, err := drive.CreateFolder(fileSrv, driveDirName, "")
		if err != nil {
			log.Fatal(err)
		}

		files, err := drive.ListFiles(fileSrv, folder.Id)
		if err != nil {
			log.Fatal(err)
		}

		if err := render.RenderComponent(w, "/files.html", files.Files); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/preview", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/scan/capture", func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.CheckToken(w, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		srv, err := drive.GetDriveFileService(token, drive.GetContext())
		if err != nil {
			log.Fatal(err)
		}

		dir, err := drive.FindFolder(srv, driveDirName)
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
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
