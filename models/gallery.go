package models

import (
	"database/sql"
	"fmt"
)

type Gallery struct {
	ID     uint
	UserID uint
	title  string
}

type GalleryService struct {
	DB *sql.DB
}

func (service *GalleryService) Create(userID uint, title string) (*Gallery, error) {
	gallery := Gallery{UserID: userID, title: title}
	row := service.DB.QueryRow(`INSERT INTO galleries (user_id, title) VALUES ($1, $2) RETURNING id;`, gallery.UserID, gallery.title)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}
