package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/drive"
	"github.com/gorilla/mux"
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

	if _, err := drive.CreateFolder(fileSrv, drive.DriveDirName, ""); err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	w.Header().Set("HX-Redirect", "/")
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

	r.ParseForm()

	folder, err := drive.CreateFolder(fileSrv, r.Form.Get("projectName"), workspaceDir.Id)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	w.Header().Set("HX-Redirect", fmt.Sprintf("/project/%s", folder.Id))
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.CheckToken(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fileId, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println(fmt.Sprintf("File id not found in URL: %s", r.URL.Path))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	fileSrv, err := drive.GetDriveFileService(token, drive.GetContext())
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	drive.DeleteFile(fileSrv, fileId)
	w.Header().Set("HX-Refresh", "true")
}
