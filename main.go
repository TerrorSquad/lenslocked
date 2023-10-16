package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/terrorsquad/lenslocked/controllers"
	"github.com/terrorsquad/lenslocked/templates"
	"github.com/terrorsquad/lenslocked/views"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	var tpl views.Template

	tpl = views.Must(views.ParseFS(templates.FS, "home.gohtml"))
	router.Get("/", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml"))
	router.Get("/contact", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml"))
	router.Get("/faq", controllers.FAQ(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "dummy.gohtml"))
	router.Get("/dummy", controllers.StaticHandler(tpl, struct {
		DummyData string
	}{
		DummyData: "Lorem ipsum dolor sit amet",
	}))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Page not found", http.StatusNotFound) })

	fmt.Println("Server is running on port 3000")
	http.ListenAndServe("localhost:3000", router)
}
