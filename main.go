package main

import (
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
			"https://www.googleapis.com/auth/drive",
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

	http.HandleFunc("/photo-directory/set", func(w http.ResponseWriter, r *http.Request) {
		photoDirectory := r.FormValue("photo_directory")

		fmt.Println(photoDirectory)

		if photoDirectory != "" {
			setPhotoDirectory(w, photoDirectory)
		}

		err := renderComponent(w, "/files.html", nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
