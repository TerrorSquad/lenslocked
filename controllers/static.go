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
			Answer:   "This was built by Goran Ninkovic - Senior Full Stack developer.",
		},
		{
			Question: "How was this made?",
			Answer:   "This was built using the Go standard library and the Tailwind CSS framework.",
		},
		{
			Question: "Can I help?",
			Answer:   "Not yet, but <i>soon</i>.",
		},
	}

	return StaticHandler(tpl, questions)
}
