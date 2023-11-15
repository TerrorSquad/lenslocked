package views

import (
	"bytes"
	"errors"
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

type public interface {
	Public() string
}

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
			"errors": func() []string {
				return nil
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

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error displaying this page", http.StatusInternalServerError)
		return
	}
	errMessages := errorMessages(errs...)
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return csrf.TemplateField(r), nil
			},
			"currentUser": func() (*models.User, error) {
				return context.User(r.Context()), nil
			},
			"errors": func() []string {
				return errMessages
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

func errorMessages(errs ...error) []string {
	var messages []string

	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			messages = append(messages, pubErr.Public())
		} else {
			messages = append(messages, "Something went wrong.")
		}

	}

	return messages
}

type Template struct {
	htmlTpl *template.Template
}
