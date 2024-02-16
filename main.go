package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const layoutDir = "web/layout"
const pageDir = "web/pages"

const accessTokenCookieName = "access_token"

func renderPage(w http.ResponseWriter, path string, data interface{}) error {
	layoutTmpl, err := template.ParseFiles(fmt.Sprintf("%s/master.html", layoutDir))
	if err != nil {
		return err
	}

	pageTmpl, err := template.Must(layoutTmpl.Clone()).ParseFiles(fmt.Sprintf("%s%s", pageDir, path))
	if err != nil {
		return err
	}

	if err := pageTmpl.Execute(w, data); err != nil {
		return err
	}
	return nil
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	conf := &oauth2.Config{
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

	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(accessTokenCookieName)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				log.Println("Access Token cookie not found")
				log.Println(err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			default:
				log.Println(err)
				http.Error(w, "server error", http.StatusInternalServerError)
			}
		} else if err := renderPage(w, "/index.html", nil); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Url string
		}{
			Url: conf.AuthCodeURL("state"),
		}

		err := renderPage(w, "/login.html", data)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		tok, err := conf.Exchange(
			context.Background(), r.URL.Query().Get("code"))
		if err != nil {
			log.Fatal(err)
		}

		cookie := &http.Cookie{
			Name:    accessTokenCookieName,
			Value:   tok.AccessToken,
			Expires: time.Now().Add(time.Hour * 24),
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
