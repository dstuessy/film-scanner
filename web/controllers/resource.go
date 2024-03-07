package controllers

import (
	"fmt"
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
		return
	}

	fileSrv, err := drive.GetDriveFileService(token, drive.GetContext())
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	folder, err := drive.CreateFolder(fileSrv, drive.DriveDirName, "")
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	files, err := drive.ListFiles(fileSrv, folder.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	if err := render.RenderComponent(w, "/files.html", files.Files); err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func NewProjectHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.CheckToken(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fileSrv, err := drive.GetDriveFileService(token, drive.GetContext())
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	workspaceDir, err := drive.GetWorkspaceDir(fileSrv)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	folder, err := drive.CreateFolder(fileSrv, r.URL.Query().Get("dirname"), workspaceDir.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/project/%s", folder.Id))
	// http.Redirect(w, r, fmt.Sprintf("/workspace/%s", folder.Id), http.StatusSeeOther)
}
