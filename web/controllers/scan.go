package controllers

import (
	"log"
	"net/http"

	"github.com/dstuessy/film-scanner/internal/auth"
	"github.com/dstuessy/film-scanner/internal/render"
)

func NewScanHandler(w http.ResponseWriter, r *http.Request) {
	_, err := auth.CheckToken(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := render.RenderPage(w, "/new.html", nil); err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}
