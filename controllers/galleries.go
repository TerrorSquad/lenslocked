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
		Show  Template
		New   Template
		Index Template
		Edit  Template
	}
	GalleryService *models.GalleryService
}

func (controller *Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	controller.Templates.New.Execute(w, r, data)
}

func (controller *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	user := context.User(r.Context())
	gallery, err := controller.GalleryService.Create(user.ID, data.Title)
	if err != nil {
		err = errors.Public(err, "Gallery could not be created.")
		controller.Templates.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/controller/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (controller *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := controller.galleryById(w, r)
	type Image struct {
		GalleryID int
		FileName  string
	}
	var data struct {
		ID     int
		Title  string
		Images []Image
	}
	if err != nil {
		return
	}
	data.ID = gallery.ID
	data.Title = gallery.Title

	images, err := controller.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID: image.GalleryID,
			FileName:  image.FileName,
		})
	}
	controller.Templates.Show.Execute(w, r, data)
}

func (controller *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := controller.galleryById(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	var data struct {
		Title string
		ID    int
	}
	data.ID = gallery.ID
	data.Title = gallery.Title

	controller.Templates.Edit.Execute(w, r, data)
}

func (controller *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := controller.galleryById(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	var data struct {
		Title string
		ID    string
	}
	data.Title = r.FormValue("title")
	gallery.Title = data.Title
	err = controller.GalleryService.Update(*gallery)
	if err != nil {
		err = errors.Public(err, "Gallery could not be updated.")
		controller.Templates.Edit.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/controller/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (controller *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := controller.galleryById(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	err = controller.GalleryService.Delete(gallery.ID)
	if err != nil {
		err = errors.Public(err, "Gallery could not be deleted.")
		var data struct {
			Title string
			ID    string
		}
		controller.Templates.Edit.Execute(w, r, data, err)
		return
	}
	http.Redirect(w, r, "/controller", http.StatusFound)
}

func (controller *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int
		Title string
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())

	userGalleries, err := controller.GalleryService.ByUserID(user.ID)
	if err != nil {
		err = errors.Public(err, "Galleries could not be retrieved.")
		controller.Templates.Index.Execute(w, r, data, err)
		return
	}
	data.Galleries = make([]Gallery, len(userGalleries))
	for i, gallery := range userGalleries {
		data.Galleries[i] = Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		}
	}

	controller.Templates.Index.Execute(w, r, data)
}

func (controller *Galleries) Image(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	galleryId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusInternalServerError)
		return
	}
	images, err := controller.GalleryService.Images(galleryId)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var requestedImage models.Image
	var imageFound = false
	for _, image := range images {
		if image.FileName == filename {
			requestedImage = image
			imageFound = true
			break
		}
	}
	if !imageFound {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, requestedImage.Path)
}

type galleryOption func(http.ResponseWriter, *http.Request, *models.Gallery) error

func (controller *Galleries) galleryById(w http.ResponseWriter, r *http.Request, options ...galleryOption) (*models.Gallery, error) {
	var galleryId, err = strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		// TODO: Handle this error better.
		http.Error(w, "Invalid gallery ID", http.StatusInternalServerError)
		return nil, err
	}
	var gallery *models.Gallery
	gallery, err = controller.GalleryService.ByID(galleryId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something went wrong", http.StatusNotFound)
		return nil, err
	}
	for _, option := range options {
		err = option(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}
	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())

	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit this gallery.", http.StatusForbidden)
		return errors.Public(nil, "You do not have access to this gallery.")
	}
	return nil
}
