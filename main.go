package main

import (
	"fmt"
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

	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", controllers.HomeHandler)

	http.HandleFunc("/login", controllers.LoginHandler)

	http.HandleFunc("/oauth2callback", controllers.AuthCallbackHandler)

	http.HandleFunc("/scan/new", controllers.NewScanHandler)

	http.HandleFunc("/resource/workspace/create", controllers.NewWorkspaceHandler)

	http.HandleFunc("/capture/stream", controllers.StreamHandler)

	http.HandleFunc("/capture/scan", controllers.CaptureScanHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
