package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const layoutDir = "web/layout"
const pageDir = "web/pages"
const componentDir = "web/components"

const accessTokenCookieName = "access_token"
const photoDirectoryCookieName = "photo_directory"

const driveDirName = "Open Scanner"

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
}

func main() {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/login", loginHandler)

	http.HandleFunc("/oauth2callback", authCallbackHandler)

	http.HandleFunc("/resource/photo-directory/create", func(w http.ResponseWriter, r *http.Request) {
		token, err := checkToken(w, r)
		if err != nil {
			log.Fatal(err)
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

		if renderComponent(w, "/files.html", files) != nil {
			log.Fatal(err)
		}
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
