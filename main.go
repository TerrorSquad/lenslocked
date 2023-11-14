package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/terrorsquad/lenslocked/controllers"
	"github.com/terrorsquad/lenslocked/migrations"
	"github.com/terrorsquad/lenslocked/models"
	"github.com/terrorsquad/lenslocked/templates"
	"github.com/terrorsquad/lenslocked/views"
	"log"
	"net/http"
	"os"
	"runtime"
)

func main() {
	// Set up database connection

	cfg := models.DefaultPostgresConfig()
	fmt.Println(cfg.String())
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup services

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// Setup middleware

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	address := ""
	if runtime.GOOS == "darwin" {
		address = "localhost"
	}

	csrfKey := []byte("gA29bm9uY2UgY2FsbCB0aGlzIGlzIGEgY29va2ll")
	csrfMw := csrf.Protect(
		csrfKey,
		// TODO: Fix this before deploying to production
		csrf.Secure(false),
	)

	// Setup controllers

	usersController := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	var baseLayouts = []string{"layouts/layout-page.gohtml", "layouts/layout-page-tailwind.gohtml"}
	usersController.Templates.SignIn = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/signin.gohtml")...))
	usersController.Templates.New = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/signup.gohtml")...))
	usersController.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/forgot_password.gohtml")...))

	// Setup router and routes

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(csrfMw)
	router.Use(umw.SetUser)
	var tpl views.Template

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/home.gohtml")...))
	router.Get("/", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/contact.gohtml")...))
	router.Get("/contact", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/faq.gohtml")...))
	router.Get("/faq", controllers.FAQ(tpl))

	router.Get("/signup", usersController.New)
	router.Post("/users", usersController.Create)
	router.Get("/signin", usersController.SignIn)
	router.Post("/signin", usersController.ProcessSignIn)
	router.Post("/signout", usersController.ProcessSignOut)
	router.Get("/forgot-pw", usersController.ForgotPassword)
	router.Post("/forgot-pw", usersController.ProcessForgotPassword)
	router.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersController.CurrentUser)
	})
	router.NotFound(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Page not found", http.StatusNotFound) })

	fmt.Println("Server is running on port: " + PORT)
	log.Println("Server is running on port: " + PORT)
	log.Fatal(http.ListenAndServe(address+":"+PORT, router))
}
