package main

import (
	"errors"
	"log"
	"net/http"

	"google.golang.org/api/drive/v3"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := checkToken(w, r)
	if err != nil {
		return
	}

	pd, err := getPhotoDirectory(r)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		log.Fatal(err)
	}

	srv, err := getDriveService(token, getContext())
	if err != nil {
		if errors.Is(err, new(TokenExpiredError)) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			log.Fatal(err)
		}
	}

	files, err := listFiles(srv, "")
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		PhotoDirectory string
		Files          []*drive.File
	}{
		PhotoDirectory: pd,
		Files:          files.Files,
	}

	if err := renderPage(w, "/index.html", data); err != nil {
		log.Fatal(err)
	}
}
