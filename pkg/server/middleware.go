package server

import (
	"net/http"
)

func secureMiddleware(f func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "user")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Redirect unauthenticated user to login page
		if userId := session.Values["user"]; userId == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		f(w, r)
	}
}
