package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func executeTemplate(w http.ResponseWriter, filepath string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	tmp, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "home.gohtml"), nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "contact.gohtml"), nil)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "faq.gohtml"), nil)
}

type MetaData struct {
	Address string
	Phone   string
}
type User struct {
	Name     string
	Age      int
	Email    string
	MetaData MetaData
}
type Data struct {
	User    User
	Integer int
	Float   float64
	Bool    bool
	Map     map[string]string
	Slice   []string
}

func playgroundHandler(w http.ResponseWriter, r *http.Request) {

	var userData = User{
		Name:  "John Doe",
		Age:   30,
		Email: "john.doe@example.com",
		MetaData: MetaData{
			Address: "Hollywood Boulevard 42",
			Phone:   "555-1234-5678",
		},
	}

	var data = Data{
		User:    userData,
		Integer: 42,
		Float:   3.14,
		Bool:    true,
		Map: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		Slice: []string{"a", "b", "c"},
	}

	executeTemplate(w, filepath.Join("templates", "playground.gohtml"), data)
}

func main() {
	router := setupRouter()
	http.ListenAndServe("localhost:3000", router)
	fmt.Println("Server is running on port 3000")
}

func setupRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", homeHandler)
	router.Get("/contact", contactHandler)
	router.Get("/faq", faqHandler)
	router.Get("/playground", playgroundHandler)
	router.NotFound(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Page not found", http.StatusNotFound) })
	return router
}
