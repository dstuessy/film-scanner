package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/drive"
	"github.com/dstuessy/film-scanner/internal/render"
	gdrive "google.golang.org/api/drive/v3"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.CheckToken(w, r)
	if err != nil {
		return
	}

	srv, err := drive.GetDriveFileService(token, drive.GetContext())
	if err != nil {
		if errors.Is(err, new(auth.TokenExpiredError)) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			log.Fatal(err)
		}
	}

	dir, err := drive.FindFolder(srv, drive.DriveDirName)
	if err != nil {
		log.Fatal(err)
	}

	dirId := ""

	if dir != nil {
		dirId = dir.Id
	}

	files, err := drive.ListFiles(srv, dirId)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Directory *gdrive.File
		Files     []*gdrive.File
	}{
		Directory: dir,
		Files:     files.Files,
	}

	if err := render.RenderPage(w, "/index.html", data); err != nil {
		log.Fatal(err)
	}
}
