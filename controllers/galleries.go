package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/terrorsquad/lenslocked/context"
	"github.com/terrorsquad/lenslocked/errors"
	"github.com/terrorsquad/lenslocked/models"
	"math/rand"
	"net/http"
	"strconv"
)

type Galleries struct {
	Templates struct {
		Show  Template
		New   Template
		Index Template
		Edit  Template
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
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (galleries *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID     int
		Title  string
		Images []string
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

	data.ID = gallery.ID
	data.Title = gallery.Title
	data.Images = make([]string, 20)
	for i := 0; i < 20; i++ {
		w, h := rand.Intn(500)+200, rand.Intn(500)+200
		catImageUrl := fmt.Sprintf("https://placekitten.com/%d/%d", w, h)
		data.Images[i] = catImageUrl
	}
	galleries.Templates.Show.Execute(w, r, data)
}

func (galleries *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
		ID    int
	}
	var galleryId, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		// TODO: Handle this error better.
		http.Error(w, "Invalid gallery ID", http.StatusInternalServerError)
		return
	}
	var gallery *models.Gallery
	gallery, err = galleries.GalleryService.ByID(galleryId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			err = errors.Public(err, fmt.Sprintf("Gallery with the ID %v was not found.", galleryId))
			galleries.Templates.Edit.Execute(w, r, data, err)
			return
		}
		err = errors.Public(err, fmt.Sprintf("Something went wrong"))
		galleries.Templates.Edit.Execute(w, r, data, err)
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		// TODO: Handle this error better.
		http.Error(w, "You do not have permission to edit this gallery.", http.StatusForbidden)
		return
	}
	data.ID = gallery.ID
	data.Title = gallery.Title

	galleries.Templates.Edit.Execute(w, r, data)
}

func (galleries *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
		ID    string
	}
	var galleryId, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		// TODO: Handle this error better.
		http.Error(w, "Invalid gallery ID", http.StatusInternalServerError)
		return
	}
	var gallery *models.Gallery
	gallery, err = galleries.GalleryService.ByID(galleryId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			err = errors.Public(err, fmt.Sprintf("Gallery with the ID %v was not found.", galleryId))
			galleries.Templates.Edit.Execute(w, r, data, err)
			return
		}
		err = errors.Public(err, fmt.Sprintf("Something went wrong"))
		galleries.Templates.Edit.Execute(w, r, data, err)
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		// TODO: Handle this error better.
		http.Error(w, "You do not have permission to edit this gallery.", http.StatusForbidden)
		return
	}
	data.Title = r.FormValue("title")
	gallery.Title = data.Title
	err = galleries.GalleryService.Update(*gallery)
	if err != nil {
		err = errors.Public(err, "Gallery could not be updated.")
		galleries.Templates.Edit.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (galleries *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
		ID    string
	}
	var galleryId, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		// TODO: Handle this error better.
		http.Error(w, "Invalid gallery ID", http.StatusInternalServerError)
		return
	}
	var gallery *models.Gallery
	gallery, err = galleries.GalleryService.ByID(galleryId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			err = errors.Public(err, fmt.Sprintf("Gallery with the ID %v was not found.", galleryId))
			galleries.Templates.Edit.Execute(w, r, data, err)
			return
		}
		err = errors.Public(err, fmt.Sprintf("Something went wrong"))
		galleries.Templates.Edit.Execute(w, r, data, err)
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		// TODO: Handle this error better.
		http.Error(w, "You do not have permission to delete this gallery.", http.StatusForbidden)
		return
	}
	err = galleries.GalleryService.Delete(galleryId)
	if err != nil {
		err = errors.Public(err, "Gallery could not be deleted.")
		galleries.Templates.Edit.Execute(w, r, data, err)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (galleries *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int
		Title string
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())

	userGalleries, err := galleries.GalleryService.ByUserID(user.ID)
	if err != nil {
		err = errors.Public(err, "Galleries could not be retrieved.")
		galleries.Templates.Index.Execute(w, r, data, err)
		return
	}
	data.Galleries = make([]Gallery, len(userGalleries))
	for i, gallery := range userGalleries {
		data.Galleries[i] = Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		}
	}

	galleries.Templates.Index.Execute(w, r, data)
}
