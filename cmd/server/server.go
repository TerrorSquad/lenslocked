package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/terrorsquad/lenslocked/controllers"
	"github.com/terrorsquad/lenslocked/migrations"
	"github.com/terrorsquad/lenslocked/models"
	"github.com/terrorsquad/lenslocked/templates"
	"github.com/terrorsquad/lenslocked/views"
	"log"
	"net/http"
	"os"
	"strconv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
		Port    string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	port := os.Getenv("SERVER_PORT")
	address := os.Getenv("SERVER_ADDRESS")
	cfg.Server.Address = address
	cfg.Server.Port = port
	return cfg, nil
}

func main() {
	// Load environment variables
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	// Set up database connection
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup services

	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	passwordResetService := &models.PasswordResetService{
		DB: db,
	}

	emailService, err := models.NewEmailService(cfg.SMTP)
	galleryService := &models.GalleryService{
		DB: db,
	}

	// Setup middleware

	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfKey := []byte(cfg.CSRF.Key)
	csrfMw := csrf.Protect(
		csrfKey,
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// Setup controllers

	usersController := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}
	galleriesController := controllers.Galleries{
		GalleryService: galleryService,
	}
	var baseLayouts = []string{"layouts/layout-page.gohtml", "layouts/layout-page-tailwind.gohtml"}
	usersController.Templates.SignIn = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/signin.gohtml")...))
	usersController.Templates.New = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/signup.gohtml")...))
	usersController.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/forgot-password.gohtml")...))
	usersController.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/check-your-email.gohtml")...))
	usersController.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/reset-password.gohtml")...))

	galleriesController.Templates.New = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "galleries/new.gohtml")...))
	galleriesController.Templates.Edit = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "galleries/edit.gohtml")...))
	galleriesController.Templates.Show = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "galleries/show.gohtml")...))
	galleriesController.Templates.Index = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "galleries/index.gohtml")...))

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
	router.Get("/reset-pw", usersController.ResetPassword)
	router.Post("/reset-pw", usersController.ProcessResetPassword)
	router.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersController.CurrentUser)
	})

	router.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesController.Show)
		r.Get("/{id}/images/{filename}", galleriesController.Image)
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/new", galleriesController.New)
			r.Get("/", galleriesController.Index)
			r.Post("/", galleriesController.Create)
			r.Get("/{id}/edit", galleriesController.Edit)
			r.Post("/{id}", galleriesController.Update)
			r.Post("/{id}/delete", galleriesController.Delete)
			r.Post("/{id}/images/{filename}/delete", galleriesController.DeleteImage)
			r.Post("/{id}/images/", galleriesController.UploadImage)
		})
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Page not found", http.StatusNotFound) })

	fmt.Println("Server is running on port: " + cfg.Server.Port)
	log.Println("Server is running on port: " + cfg.Server.Port)
	log.Fatal(http.ListenAndServe(cfg.Server.Address+":"+cfg.Server.Port, router))
}
