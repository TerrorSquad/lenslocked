package models

import (
	"database/sql"
	"fmt"
	"github.com/terrorsquad/lenslocked/errors"
)

type Gallery struct {
	ID     uint
	UserID uint
	Title  string
}

type GalleryService struct {
	DB *sql.DB
}

func (service *GalleryService) Create(userID uint, title string) (*Gallery, error) {
	gallery := Gallery{UserID: userID, Title: title}
	row := service.DB.QueryRow(`INSERT INTO galleries (user_id, title) VALUES ($1, $2) RETURNING id;`, gallery.UserID, gallery.Title)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (service *GalleryService) ByID(id uint) (*Gallery, error) {
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

func (service *GalleryService) ByUserID(userID uint) ([]Gallery, error) {
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
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	return galleries, nil
}
