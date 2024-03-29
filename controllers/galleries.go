package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/terrorsquad/lenslocked/context"
	"github.com/terrorsquad/lenslocked/errors"
	"github.com/terrorsquad/lenslocked/models"
	"net/http"
	"net/url"
	"path/filepath"
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
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (controller *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := controller.galleryById(w, r)
	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
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
			GalleryID:       image.GalleryID,
			Filename:        image.FileName,
			FilenameEscaped: url.PathEscape(image.FileName),
		})
	}
	controller.Templates.Show.Execute(w, r, data)
}

func (controller *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := controller.galleryById(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		Title  string
		ID     int
		Images []Image
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
			GalleryID:       image.GalleryID,
			Filename:        image.FileName,
			FilenameEscaped: url.PathEscape(image.FileName),
		})
	}

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
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
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
	http.Redirect(w, r, "/galleries", http.StatusFound)
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
	filename := controller.filename(w, r)
	galleryId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Invalid gallery ID", http.StatusInternalServerError)
		return
	}
	image, err := controller.GalleryService.Image(galleryId, filename)

	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, image.Path)
}

func (controller *Galleries) UploadImage(w http.ResponseWriter, r *http.Request) {
	gallery, err := controller.galleryById(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "Image could not be uploaded", http.StatusInternalServerError)
		return
	}
	fileHeaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		fmt.Printf("Attempting to upload %v for gallery %d\n", fileHeader.Filename, gallery.ID)

		err = controller.GalleryService.CreateImage(gallery.ID, fileHeader.Filename, file)
		if err != nil {
			var fileError models.FileError
			if errors.As(err, &fileError) {
				msg := fmt.Sprintf("%v has an invalid content type or extensions. Only png, gif and jpg files can be uploaded", fileHeader.Filename)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/galleries/%d/edit", gallery.ID), http.StatusFound)
	}
}

func (controller *Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := controller.filename(w, r)
	gallery, err := controller.galleryById(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = controller.GalleryService.DeleteImage(gallery.ID, filename)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/galleries/%d/edit", gallery.ID), http.StatusFound)
}

func (controller *Galleries) filename(w http.ResponseWriter, r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	filename = filepath.Base(filename)
	return filename
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
		http.Error(w, "You do not have access to this gallery.", http.StatusForbidden)
		return errors.Public(nil, "You do not have access to this gallery.")
	}
	return nil
}
