package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const AccessTokenCookieName = "access_token"

var OauthConf *oauth2.Config

type TokenExpiredError struct{}

func (m *TokenExpiredError) Error() string {
	return "Token Expired"
}

func Setup() {
	OauthConf = &oauth2.Config{
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

func CheckToken(w http.ResponseWriter, r *http.Request) (*oauth2.Token, error) {
	token, err := GetToken(r)
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

func GetToken(r *http.Request) (*oauth2.Token, error) {
	cookie, err := r.Cookie(AccessTokenCookieName)
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
