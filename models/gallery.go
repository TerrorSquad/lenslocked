package models

import (
	"database/sql"
	"fmt"
	"github.com/terrorsquad/lenslocked/errors"
	"path/filepath"
	"strings"
)

type Image struct {
	GalleryID int
	Path      string
	FileName  string
}

const (
	imagesDir = "images"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB

	// ImagesDir is used to tell the GalleryService where to store and locate the images for a gallery.
	// If not set, it will default to "images/" directory.
	ImagesDir string
}

func (service *GalleryService) Create(userID int, title string) (*Gallery, error) {
	gallery := Gallery{UserID: userID, Title: title}
	row := service.DB.QueryRow(`INSERT INTO galleries (user_id, title) VALUES ($1, $2) RETURNING id;`, gallery.UserID, gallery.Title)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

// ByID will look up a gallery by the provided ID.
//
// Possible errors:
//
// - ErrNotFound - No gallery exists with the provided ID.
// - Other errors that may be returned from the database.
func (service *GalleryService) ByID(id int) (*Gallery, error) {
	// TODO: Add validation to ensure that the gallery ID is not 0.
	gallery := Gallery{ID: id}
	row := service.DB.QueryRow(`SELECT user_id, title FROM galleries WHERE id = $1;`, gallery.ID)
	err := row.Scan(&gallery.UserID, &gallery.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by ID: %w", err)
	}
	return &gallery, nil
}

func (service *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := service.DB.Query(`SELECT id, title FROM galleries WHERE user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user ID: %w", err)
	}
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{UserID: userID}
		err := rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	// Checks for errors encountered during iteration that were not otherwise detected by rows.Scan().
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	return galleries, nil
}

func (service *GalleryService) Update(gallery Gallery) error {
	_, err := service.DB.Exec(`UPDATE galleries SET title = $1 WHERE id = $2;`, gallery.Title, gallery.ID)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (service *GalleryService) Delete(id int) error {
	_, err := service.DB.Exec(`DELETE FROM galleries WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}
	return nil
}

func (service *GalleryService) Images(galleryId int) ([]Image, error) {
	globPattern := filepath.Join(service.galleryDir(galleryId), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}
	var images []Image
	var extensions = service.extensions()
	for _, file := range allFiles {
		if hasExtension(file, extensions) {
			images = append(images, Image{
				Path:      file,
				FileName:  filepath.Base(file),
				GalleryID: galleryId,
			})
		}
	}
	return images, nil
}

func (service *GalleryService) extensions() []string {
	return []string{".jpg", ".png", ".jpeg", ".gif"}
}

func (service *GalleryService) galleryDir(id int) string {
	images := service.ImagesDir
	if images == "" {
		images = imagesDir
	}
	return filepath.Join(images, fmt.Sprintf("gallery-%d", id))
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(filepath.Ext(file))
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}
