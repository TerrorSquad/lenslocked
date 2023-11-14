package views

import (
	"bytes"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/terrorsquad/lenslocked/context"
	"github.com/terrorsquad/lenslocked/models"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(path.Base(patterns[0]))
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return ``, fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return ``, fmt.Errorf("currentUser not implemented")
			},
			"errors": func() (template.HTML, error) {
				return ``, fmt.Errorf("errors not implemented")
			},
		},
	)
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing embedded template: %w", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

//func Parse(filepath string) (Template, error) {
//	tpl, err := template.ParseFiles(filepath)
//	if err != nil {
//		return Template{}, fmt.Errorf("parsing template: %w", err)
//	}
//	return Template{
//		htmlTpl: tpl,
//	}, nil
//}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error displaying this page", http.StatusInternalServerError)
		return
	}
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return csrf.TemplateField(r), nil
			},
			"currentUser": func() (*models.User, error) {
				return context.User(r.Context()), nil
			},
			"errors": func() []string {
				return []string{}
			},
		},
	)
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

type Template struct {
	htmlTpl *template.Template
}
