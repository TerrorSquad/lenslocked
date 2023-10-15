package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/terrorsquad/lenslocked/controllers"
	"github.com/terrorsquad/lenslocked/views"
	"net/http"
	"path/filepath"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	var tpl views.Template

	tpl = views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))
	router.Get("/", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))
	router.Get("/contact", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))
	router.Get("/faq", controllers.StaticHandler(tpl, nil))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Page not found", http.StatusNotFound) })

	http.ListenAndServe("localhost:3000", router)
	fmt.Println("Server is running on port 3000")
}
