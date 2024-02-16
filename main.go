package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const layoutDir = "web/layout"
const pageDir = "web/pages"
const componentDir = "web/components"

const accessTokenCookieName = "access_token"

var oauthConf *oauth2.Config

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

func renderComponent(w http.ResponseWriter, path string, data interface{}) error {
	cmpTmpl, err := template.ParseFiles(fmt.Sprintf("%s%s", componentDir, path))
	if err != nil {
		return err
	}

	if err := cmpTmpl.Execute(w, data); err != nil {
		return err
	}

	return nil
}

func getToken(r *http.Request) (*oauth2.Token, error) {
	cookie, err := r.Cookie(accessTokenCookieName)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.URLEncoding.DecodeString(cookie.Value)

	token := &oauth2.Token{}
	if err := json.Unmarshal(decoded, token); err != nil {
		return nil, err
	}

	return token, err
}

func checkToken(w http.ResponseWriter, r *http.Request) (*oauth2.Token, error) {
	token, err := getToken(r)
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
	}
	return token, err
}

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		token, err := checkToken(w, r)
		if err != nil {
			return
		}

		ctx := context.Background()

		srv, err := drive.NewService(ctx, option.WithTokenSource(oauthConf.TokenSource(ctx, token)))
		if err != nil {
			log.Fatalf("Unable to retrieve Drive client: %v", err)
		}

		files, err2 := srv.Files.List().
			PageSize(10).
			Q("mimeType='application/vnd.google-apps.folder' and name contains 'Photography'").
			Fields("nextPageToken, files(id, name)").
			Spaces("drive").
			Do()
		if err2 != nil {
			log.Fatalf("Unable to retrieve files: %v", err2)
		}

		fmt.Println(files.Files, files.HTTPStatusCode)

		data := struct {
			Files []*drive.File
		}{
			Files: files.Files,
		}

		if err := renderPage(w, "/index.html", data); err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Url string
		}{
			Url: oauthConf.AuthCodeURL("state"),
		}

		err := renderPage(w, "/login.html", data)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		tok, err := oauthConf.Exchange(
			context.Background(), r.URL.Query().Get("code"))
		if err != nil {
			log.Fatal(err)
		}

		tokenJson, err := json.Marshal(tok)
		if err != nil {
			log.Fatal(err)
		}

		cookie := &http.Cookie{
			Name:    accessTokenCookieName,
			Value:   base64.URLEncoding.EncodeToString(tokenJson),
			Expires: time.Now().Add(time.Hour * 24),
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
