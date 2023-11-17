package controllers

import (
	"github.com/terrorsquad/lenslocked/models"
	"net/http"
)

type Galleries struct {
	Templates struct {
		New Template
	}
	GalleryService *models.GalleryService
}

func (galleries *Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	galleries.Templates.New.Execute(w, r, data)
}
