package controllers

import (
	"fmt"
	"github.com/terrorsquad/lenslocked/context"
	"github.com/terrorsquad/lenslocked/models"
	"net/http"
)

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := readCookie(r, CookieSession)
		if err != nil {
			//fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(tokenCookie)
		if err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			// TODO: Add a flash message to tell the user why they were redirected.
			return
		}
		next.ServeHTTP(w, r)
	})
}
