package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/cache"
	"github.com/dstuessy/film-scanner/internal/camera"
	"github.com/dstuessy/film-scanner/web/controllers"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := camera.SetupTempDir(); err != nil {
		log.Fatal(err)
	}

	if err := cache.SetupCacheDir(); err != nil {
		log.Fatal(err)
	}

	auth.Setup()

	if err := camera.StartStream(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer camera.CloseCamera()
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("web/assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	r.HandleFunc("/", controllers.HomeHandler)

	r.HandleFunc("/project/{id}", controllers.ProjectHandler)

	r.HandleFunc("/project/{id}/scan", controllers.NewScanHandler)

	r.HandleFunc("/login", controllers.LoginHandler)

	r.HandleFunc("/oauth2callback", controllers.AuthCallbackHandler)

	r.HandleFunc("/resource/workspace/create", controllers.NewWorkspaceHandler)

	r.HandleFunc("/resource/project/create", controllers.NewProjectHandler)

	r.HandleFunc("/resource/project/{id}", controllers.GetProjectHandler)

	r.HandleFunc("/resource/file/{id}/delete", controllers.DeleteFileHandler)

	r.HandleFunc("/resource/cache/{project}/upload", controllers.UploadCacheHandler)

	r.HandleFunc("/resource/cache/{project}/file/{file}/delete", controllers.DeleteCacheFileHandler)

	r.HandleFunc("/capture/stream", controllers.StreamHandler)

	r.HandleFunc("/capture/scan", controllers.CaptureScanHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
