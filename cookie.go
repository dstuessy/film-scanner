package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

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
