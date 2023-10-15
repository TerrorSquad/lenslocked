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

func pathHandler(w http.ResponseWriter, r *http.Request) {
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

type Router struct{}

//func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	switch r.URL.Path {
//	case "/":
//		homeHandler(w, r)
//		break
//	case "/contact":
//		contactHandler(w, r)
//		break
//	default:
//		http.NotFound(w, r)
//	}
//}

func main() {
	fmt.Println("Server is running on port 3000")

	//// First way: store the handler in a variable
	//// long form
	//var router http.HandlerFunc = pathHandler
	//// short form
	//router := http.HandlerFunc(pathHandler)
	//// pass the router to the ListenAndServe function
	//http.ListenAndServe("localhost:3000", router)

	// Second way: use http.HandlerFunc as a type conversion directly.
	// Second argument is not a function call, but a conversion, similar to (int) 42 in C
	http.ListenAndServe("localhost:3000", http.HandlerFunc(pathHandler))
}
