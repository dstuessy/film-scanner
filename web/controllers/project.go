package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/drive"
	"github.com/dstuessy/film-scanner/internal/render"
	"github.com/gorilla/mux"
	gdrive "google.golang.org/api/drive/v3"
)

type Breadcrumb struct {
	Name string
	Link string
}

func ProjectHandler(w http.ResponseWriter, r *http.Request) {
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

	projectId, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println(fmt.Sprintf("Project id not found in URL: %s", r.URL.Path))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	dir, err := drive.GetFile(srv, projectId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	dirId := ""
	dirname := ""

	if dir != nil {
		dirId = dir.Id
		dirname = dir.Name
	}

	files, err := drive.ListFiles(srv, dirId, "")
	if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Directory     *gdrive.File
		Breadcrumbs   []Breadcrumb
		NextPageToken string
		Files         []*gdrive.File
	}{
		Directory: dir,
		Breadcrumbs: []Breadcrumb{
			{Name: drive.DriveDirName, Link: "/"},
			{Name: dirname, Link: ""},
		},
		NextPageToken: files.NextPageToken,
		Files:         files.Files,
	}

	if err := render.RenderPage(w, "/project.html", data); err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
}
