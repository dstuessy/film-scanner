package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/camera"
	"github.com/dstuessy/film-scanner/web/controllers"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	auth.Setup()

	if err := camera.Open(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer camera.Close()
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("web/assets"))
	r.Handle("/assets/", http.StripPrefix("/assets/", fs))

	r.HandleFunc("/", controllers.HomeHandler)

	r.HandleFunc("/project/{id}", controllers.ProjectHandler)

	r.HandleFunc("/project/{id}/scan", controllers.NewScanHandler)

	r.HandleFunc("/login", controllers.LoginHandler)

	r.HandleFunc("/oauth2callback", controllers.AuthCallbackHandler)

	r.HandleFunc("/resource/workspace/create", controllers.NewWorkspaceHandler)

	r.HandleFunc("/resource/project/create", controllers.NewProjectHandler)

	r.HandleFunc("/capture/stream", controllers.StreamHandler)

	r.HandleFunc("/capture/scan", controllers.CaptureScanHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
