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

	srv, err := getDriveFileService(token, getContext())
	if err != nil {
		if errors.Is(err, new(TokenExpiredError)) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			log.Fatal(err)
		}
	}

	dir, err := findFolder(srv, driveDirName)
	if err != nil {
		log.Fatal(err)
	}

	dirId := ""

	if dir != nil {
		dirId = dir.Id
	}

	files, err := listFiles(srv, dirId)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Directory *drive.File
		Files     []*drive.File
	}{
		Directory: dir,
		Files:     files.Files,
	}

	if err := renderPage(w, "/index.html", data); err != nil {
		log.Fatal(err)
	}
}
