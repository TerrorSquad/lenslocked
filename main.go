package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/terrorsquad/lenslocked/controllers"
	"github.com/terrorsquad/lenslocked/templates"
	"github.com/terrorsquad/lenslocked/views"
	"log"
	"net/http"
	"os"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	var tpl views.Template

	var baseLayouts = []string{"layout-page.gohtml", "layout-page-tailwind.gohtml"}

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "home.gohtml")...))
	router.Get("/", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "contact.gohtml")...))
	router.Get("/contact", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "faq.gohtml")...))
	router.Get("/faq", controllers.FAQ(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "dummy.gohtml")...))
	router.Get("/dummy", controllers.StaticHandler(tpl, struct {
		DummyData string
	}{
		DummyData: "Lorem ipsum dolor sit amet",
	}))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Page not found", http.StatusNotFound) })

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	address := ""
	if os.Getenv("OS") == "macos" {
		address = "localhost"
	}

	fmt.Println("Server is running on port: " + PORT)
	log.Println("Server is running on port: " + PORT)
	log.Fatal(http.ListenAndServe(address+":"+PORT, router))
}
