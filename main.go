package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dstuessy/film-scanner/internal/camera"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const layoutDir = "web/layout"
const pageDir = "web/pages"
const componentDir = "web/components"

const accessTokenCookieName = "access_token"

const driveDirName = "Open Scanner"

const boundaryWord = "MJPEGBOUNDARY"

var frameInterval time.Duration

var oauthConf *oauth2.Config

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/drive.file",
		},
		Endpoint: google.Endpoint,
	}

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
		_, err := checkToken(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		if err := renderPage(w, "/new.html", nil); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/resource/workspace/create", func(w http.ResponseWriter, r *http.Request) {
		token, err := checkToken(w, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

		fileSrv, err := getDriveFileService(token, getContext())
		if err != nil {
			if errors.Is(err, new(TokenExpiredError)) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			} else {
				log.Fatal(err)
			}
		}

		folder, err := createFolder(fileSrv, driveDirName, "")
		if err != nil {
			log.Fatal(err)
		}

		files, err := listFiles(fileSrv, folder.Id)
		if err != nil {
			log.Fatal(err)
		}

		if err := renderComponent(w, "/files.html", files.Files); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/preview", func(w http.ResponseWriter, r *http.Request) {
		_, err := checkToken(w, r)
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

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
