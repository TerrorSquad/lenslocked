package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/terrorsquad/lenslocked/context"
	"github.com/terrorsquad/lenslocked/errors"
	"github.com/terrorsquad/lenslocked/models"
	"net/http"
	"strconv"
)

type Galleries struct {
	Templates struct {
		New  Template
		Show Template
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

func (galleries *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	user := context.User(r.Context())
	gallery, err := galleries.GalleryService.Create(user.ID, data.Title)
	if err != nil {
		err = errors.Public(err, "Gallery could not be created.")
		galleries.Templates.New.Execute(w, r, data, err)
		return
	}
	http.Redirect(w, r, "/galleries/"+strconv.Itoa(gallery.ID), http.StatusFound)
}

func (galleries *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	var data struct {
	}
	var galleryId, err = strconv.Atoi(chi.URLParam(r, "id"))
	var gallery *models.Gallery
	gallery, err = galleries.GalleryService.ByID(galleryId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			err = errors.Public(err, fmt.Sprintf("Gallery with the ID %v was not found.", galleryId))
			galleries.Templates.Show.Execute(w, r, data, err)
		} else {
			err = errors.Public(err, fmt.Sprintf("Something went wrong"))
			galleries.Templates.Show.Execute(w, r, data, err)
		}
	}

	galleries.Templates.Show.Execute(w, r, gallery)
}
