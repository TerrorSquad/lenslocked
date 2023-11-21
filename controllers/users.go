package controllers

import (
	"fmt"
	"github.com/terrorsquad/lenslocked/context"
	"github.com/terrorsquad/lenslocked/errors"
	"github.com/terrorsquad/lenslocked/models"
	"net/http"
	"net/url"
)

type Users struct {
	Templates struct {
		New            Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

func (controller Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	controller.Templates.New.Execute(w, r, data)
}

func (controller Users) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	password := r.FormValue("password")
	user, err := controller.UserService.Create(data.Email, password)
	if err != nil {
		if errors.Is(err, models.ErrEmailTaken) {
			err = errors.Public(err, "That email address is already associated with an account")
		}
		controller.Templates.New.Execute(w, r, data, err)
		return
	}
	session, err := controller.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: Long term, we should show a warning to the user that their session was not created.
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries/", http.StatusFound)
	fmt.Fprintf(w, "User created: %+v", user)
}

func (controller Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	controller.Templates.SignIn.Execute(w, r, data)
}

func (controller Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := controller.UserService.Authenticate(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := controller.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (controller Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = controller.SessionService.Delete(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)

	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (controller Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Fprintf(w, "User: %s\n", user.Email)
	// TODO: Render a template with the user's information
}

func (controller Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	controller.Templates.ForgotPassword.Execute(w, r, data)
}

func (controller Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := controller.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: Handle other cases in the future
		// e.g. user does not exist
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": []string{pwReset.Token},
	}
	resetURL := "https://localhost:8080/reset-pw?" + vals.Encode()
	err = controller.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	// Don't render the reset token here. We need the user to confirm that they have received the email.
	// Render a page that tells them to check their email.
	controller.Templates.CheckYourEmail.Execute(w, r, data)
}

func (controller Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	controller.Templates.ResetPassword.Execute(w, r, data)
}
func (controller Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := controller.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		// TODO: distinguish between invalid token and other errors
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	err = controller.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Sign the user in now that their password has been reset
	// Any errors should redirect the user to the sign in page.
	session, err := controller.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries", http.StatusFound)
}
