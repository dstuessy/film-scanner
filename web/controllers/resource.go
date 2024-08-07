package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/cache"
	"github.com/dstuessy/film-scanner/internal/drive"
	"github.com/dstuessy/film-scanner/internal/render"
	"github.com/gorilla/mux"
	gdrive "google.golang.org/api/drive/v3"
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
	return
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

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
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

	projectId, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println(fmt.Sprintf("Project id not found in URL: %s", r.URL.Path))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	dir, err := drive.GetFile(fileSrv, projectId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	page := r.URL.Query().Get("page")

	if page == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	files, err := drive.ListFiles(fileSrv, projectId, page)
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	log.Println("next page", files.NextPageToken)

	data := struct {
		Directory     *gdrive.File
		NextPageToken string
		Files         []*gdrive.File
	}{
		Directory:     dir,
		NextPageToken: files.NextPageToken,
		Files:         files.Files,
	}

	if err := render.RenderComponent(w, "/files.html", data); err != nil {
		log.Println(err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
	}
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
		return
	}

	drive.DeleteFile(fileSrv, fileId)
	w.Header().Set("HX-Refresh", "true")
	return
}

func DeleteCacheFileHandler(w http.ResponseWriter, r *http.Request) {
	_, err := auth.CheckToken(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	projectId, ok := mux.Vars(r)["project"]
	if !ok {
		log.Println(fmt.Sprintf("Project id not found in URL: %s", r.URL.Path))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	fileName, ok := mux.Vars(r)["file"]
	if !ok {
		log.Println(fmt.Sprintf("File name not found in URL: %s", r.URL.Path))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if cache.DeleteImage(projectId, fileName) != nil {
		log.Println(err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	return
}

func UploadCacheHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.CheckToken(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	projectId, ok := mux.Vars(r)["project"]
	if !ok {
		log.Println(fmt.Sprintf("Project id not found in URL: %s", r.URL.Path))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	files, err := cache.ReadProject(projectId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	srv, err := drive.GetDriveFileService(token, drive.GetContext())
	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		jpeg, err := cache.ReadImage(projectId, file)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		if _, err := drive.SaveImage(srv, jpeg, file, projectId); err != nil {
			log.Println(err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		if cache.DeleteImage(projectId, file) != nil {
			log.Println(err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("HX-Refresh", "true")
	return
}
