package controllers

import (
	"github.com/terrorsquad/lenslocked/views"
	"net/http"
)

func StaticHandler(tpl views.Template, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, data)
	}
}
