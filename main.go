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
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	var tpl views.Template

	var baseLayouts = []string{"layouts/layout-page.gohtml", "layouts/layout-page-tailwind.gohtml"}

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/home.gohtml")...))
	router.Get("/", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/contact.gohtml")...))
	router.Get("/contact", controllers.StaticHandler(tpl, nil))

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/faq.gohtml")...))
	router.Get("/faq", controllers.FAQ(tpl))

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

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	usersController := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersController.Templates.SignIn = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/signin.gohtml")...))
	usersController.Templates.New = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/signup.gohtml")...))
	router.Get("/signup", usersController.New)
	router.Post("/users", usersController.Create)
	router.Get("/signin", usersController.SignIn)
	router.Post("/signin", usersController.Authenticate)
	router.Get("/users/me", usersController.CurrentUser)
	router.Post("/signout", usersController.SignOut)

	tpl = views.Must(views.ParseFS(templates.FS, append(baseLayouts, "pages/dummy.gohtml")...))
	router.Get("/dummy", controllers.StaticHandler(tpl, struct {
		DummyData string
	}{
		DummyData: "Lorem ipsum dolor sit amet",
	}))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Page not found", http.StatusNotFound) })

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
	fmt.Println("Server is running on port: " + PORT)
	log.Println("Server is running on port: " + PORT)
	log.Fatal(http.ListenAndServe(address+":"+PORT, csrfMw(umw.SetUser(router))))
}
