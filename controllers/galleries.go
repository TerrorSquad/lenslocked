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
	galleries.Templates.New.Execute(w, r, nil)
}
