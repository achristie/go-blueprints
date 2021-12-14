package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.next.ServeHTTP(w, r)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	provider := segs[2]
	action := segs[3]
	switch action {
	case "login":
		provider, err := goth.GetProvider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s", provider, err), http.StatusBadRequest)
		}
		gothic.BeginAuthHandler(w, r)
	default:
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not complete authentication: %s", err), http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: "blah",
			Path:  "/",
		})
		fmt.Fprintf(w, "%v", user)
		// w.Header().Set("Location", "/chat")
		// w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
