package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/render"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Url string
	}{
		Url: auth.OauthConf.AuthCodeURL("state"),
	}

	err := render.RenderPage(w, "/login.html", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
}

func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.OauthConf.Exchange(
		context.Background(), r.URL.Query().Get("code"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenJson, err := json.Marshal(tok)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:    auth.AccessTokenCookieName,
		Value:   base64.URLEncoding.EncodeToString(tokenJson),
		Expires: time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
