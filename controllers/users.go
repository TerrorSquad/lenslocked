package controllers

import (
	"github.com/terrorsquad/lenslocked/views"
	"net/http"
)

type Users struct {
	Templates struct {
		New views.Template
	}
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.Templates.New.Execute(w, nil)
}
