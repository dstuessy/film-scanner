package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"
)

func getPhotoDirectory(r *http.Request) (string, error) {
	cookie, err := r.Cookie(photoDirectoryCookieName)
	if err != nil {
		return "", err
	}

	decoded, err := base64.URLEncoding.DecodeString(cookie.Value)

	return string(decoded), err
}

func setPhotoDirectory(w http.ResponseWriter, directory string) {
	encoded := base64.URLEncoding.EncodeToString([]byte(directory))

	cookie := &http.Cookie{
		Name:  photoDirectoryCookieName,
		Value: encoded,
	}

	http.SetCookie(w, cookie)
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
