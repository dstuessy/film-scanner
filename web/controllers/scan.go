package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/render"
	"github.com/gorilla/mux"
)

func NewScanHandler(w http.ResponseWriter, r *http.Request) {
	_, err := auth.CheckToken(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	projectId, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println(fmt.Sprintf("Project id not found in URL: %s", r.URL.Path))
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	data := struct {
		ProjectId string
	}{
		ProjectId: projectId,
	}

	if err := render.RenderPage(w, "/new.html", data); err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}
