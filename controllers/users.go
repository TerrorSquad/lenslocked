package controllers

import (
	"fmt"
	"github.com/terrorsquad/lenslocked/models"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Email: "+r.FormValue("email"))
	fmt.Fprintln(w, "Password: "+r.FormValue("password"))
}
