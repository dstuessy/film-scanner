package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Url string
	}{
		Url: oauthConf.AuthCodeURL("state"),
	}

	err := renderPage(w, "/login.html", data)
	if err != nil {
		log.Fatal(err)
	}
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
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
}
