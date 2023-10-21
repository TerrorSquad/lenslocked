package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Template, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, data)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "What is this?",
			Answer:   "This is a simple image gallery built in Go.",
		},
		{
			Question: "Who made this?",
			Answer:   "This was built by Jon Calhoun for his upcoming book on web development in Go.",
		},
		{
			Question: "How was this made?",
			Answer:   "This was built using the Go standard library and the Tailwind CSS framework.",
		},
		{
			Question: "Can I help?",
			Answer:   "Not yet, but <i>soon</i>. Follow @joncalhoun on Twitter for more updates.",
		},
	}

	return StaticHandler(tpl, questions)
}
