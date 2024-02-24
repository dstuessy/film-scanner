package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/drive"
	"github.com/dstuessy/film-scanner/internal/render"
)

func NewWorkspaceHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.CheckToken(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	fileSrv, err := drive.GetDriveFileService(token, drive.GetContext())
	if err != nil {
		if errors.Is(err, new(auth.TokenExpiredError)) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else {
			log.Fatal(err)
		}
	}

	folder, err := drive.CreateFolder(fileSrv, drive.DriveDirName, "")
	if err != nil {
		log.Fatal(err)
	}

	files, err := drive.ListFiles(fileSrv, folder.Id)
	if err != nil {
		log.Fatal(err)
	}

	if err := render.RenderComponent(w, "/files.html", files.Files); err != nil {
		log.Fatal(err)
	}
}
