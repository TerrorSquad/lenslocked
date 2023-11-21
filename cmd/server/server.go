package main

import (
	"database/sql"
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
	err = run(cfg)
	if err != nil {
		panic(err)
	}
}

func run(cfg config) error {
	// Set up database connection
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return err
	}

	server := initServer(cfg, db)

	fmt.Println("Server is running on address:port: " + cfg.Server.Address + ":" + cfg.Server.Port)
	return http.ListenAndServe(cfg.Server.Address+":"+cfg.Server.Port, server.Router)
}

type Services struct {
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
	GalleryService       *models.GalleryService
}
type Controllers struct {
	UsersController     controllers.Users
	GalleriesController controllers.Galleries
}
type Middleware struct {
	UserMiddleware controllers.UserMiddleware
	CSRFMiddleware func(http.Handler) http.Handler
}

type Server struct {
	Services    *Services
	Middleware  *Middleware
	Controllers *Controllers
	BaseLayouts []string
	Router      *chi.Mux
}

func initServices(cfg config, db *sql.DB) *Services {
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
	if err != nil {
		panic(err)
	}
	return &Services{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
		GalleryService:       galleryService,
	}
}

func initServer(cfg config, db *sql.DB) Server {
	server := Server{}
	var baseLayouts = []string{"layouts/layout-page.gohtml", "layouts/layout-page-tailwind.gohtml"}
	server.BaseLayouts = baseLayouts
	server.Services = initServices(cfg, db)
	server.Middleware = initMiddleware(cfg, server)
	server.Controllers = initControllers(server)
	server.Router = initRouter(server)
	return server
}

func initRouter(server Server) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(server.Middleware.CSRFMiddleware)
	router.Use(server.Middleware.UserMiddleware.SetUser)
	var tpl views.Template

	tpl = views.Must(views.ParseFS(templates.FS, append(server.BaseLayouts, "pages/home.gohtml")...))
	router.Get("/", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.ParseFS(templates.FS, append(server.BaseLayouts, "pages/contact.gohtml")...))
	router.Get("/contact", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.ParseFS(templates.FS, append(server.BaseLayouts, "pages/faq.gohtml")...))
	router.Get("/faq", controllers.FAQ(tpl))

	router.Get("/signup", server.Controllers.UsersController.New)
	router.Post("/users", server.Controllers.UsersController.Create)
	router.Get("/signin", server.Controllers.UsersController.SignIn)
	router.Post("/signin", server.Controllers.UsersController.ProcessSignIn)
	router.Post("/signout", server.Controllers.UsersController.ProcessSignOut)
	router.Get("/forgot-pw", server.Controllers.UsersController.ForgotPassword)
	router.Post("/forgot-pw", server.Controllers.UsersController.ProcessForgotPassword)
	router.Get("/reset-pw", server.Controllers.UsersController.ResetPassword)
	router.Post("/reset-pw", server.Controllers.UsersController.ProcessResetPassword)
	router.Route("/users/me", func(r chi.Router) {
		r.Use(server.Middleware.UserMiddleware.RequireUser)
		r.Get("/", server.Controllers.UsersController.CurrentUser)
	})

	router.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", server.Controllers.GalleriesController.Show)
		r.Get("/{id}/images/{filename}", server.Controllers.GalleriesController.Image)
		r.Group(func(r chi.Router) {
			r.Use(server.Middleware.UserMiddleware.RequireUser)
			r.Get("/new", server.Controllers.GalleriesController.New)
			r.Get("/", server.Controllers.GalleriesController.Index)
			r.Post("/", server.Controllers.GalleriesController.Create)
			r.Get("/{id}/edit", server.Controllers.GalleriesController.Edit)
			r.Post("/{id}", server.Controllers.GalleriesController.Update)
			r.Post("/{id}/delete", server.Controllers.GalleriesController.Delete)
			r.Post("/{id}/images/{filename}/delete", server.Controllers.GalleriesController.DeleteImage)
			r.Post("/{id}/images/", server.Controllers.GalleriesController.UploadImage)
		})
	})

	assetsHandler := http.FileServer(http.Dir("assets"))
	router.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Page not found", http.StatusNotFound) })
	return router
}

func initControllers(server Server) *Controllers {
	usersController := controllers.Users{
		UserService:          server.Services.UserService,
		SessionService:       server.Services.SessionService,
		PasswordResetService: server.Services.PasswordResetService,
		EmailService:         server.Services.EmailService,
	}
	galleriesController := controllers.Galleries{
		GalleryService: server.Services.GalleryService,
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

	return &Controllers{
		UsersController:     usersController,
		GalleriesController: galleriesController,
	}
}

func initMiddleware(cfg config, server Server) *Middleware {
	umw := controllers.UserMiddleware{
		SessionService: server.Services.SessionService,
	}

	csrfKey := []byte(cfg.CSRF.Key)
	csrfMw := csrf.Protect(
		csrfKey,
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)
	return &Middleware{
		UserMiddleware: umw,
		CSRFMiddleware: csrfMw,
	}
}
