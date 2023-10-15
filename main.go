package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<p>To get in touch, please send an email to <a href=\"mailto:test@example.com\">John Doe</a>.</p>")
}

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
		break
	case "/contact":
		contactHandler(w, r)
		break
	default:
		http.NotFound(w, r)
	}
}

func main() {
	var router Router
	fmt.Println("Server is running on port 3000")
	http.ListenAndServe("localhost:3000", router)
}
