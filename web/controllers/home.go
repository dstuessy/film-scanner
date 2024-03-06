package controllers

import (
	"log"
	"net/http"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/drive"
	"github.com/dstuessy/film-scanner/internal/render"
	gdrive "google.golang.org/api/drive/v3"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.CheckToken(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	srv, err := drive.GetDriveFileService(token, drive.GetContext())
	if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	dir, err := drive.FindFolder(srv, drive.DriveDirName)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	dirId := ""

	if dir != nil {
		dirId = dir.Id
	}

	files, err := drive.ListFiles(srv, dirId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Directory   *gdrive.File
		Breadcrumbs []string
		Files       []*gdrive.File
	}{
		Directory:   dir,
		Breadcrumbs: []string{drive.DriveDirName},
		Files:       files.Files,
	}

	if err := render.RenderPage(w, "/index.html", data); err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
}
