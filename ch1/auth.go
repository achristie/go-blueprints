package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/objx"
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
		gothic.GetProviderName = func(req *http.Request) (string, error) {
			return "github", nil
		}
		gothic.BeginAuthHandler(w, r)
	case "callback":
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not complete authentication: %s", err), http.StatusBadRequest)
			return
		}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email))
		userId := fmt.Sprintf("%x", m.Sum(nil))
		authCookie := objx.New(map[string]interface{}{
			"userid":     userId,
			"name":       user.Name,
			"avatar_url": user.AvatarURL,
			"email":      user.Email,
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookie,
			Path:  "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		http.Error(w, fmt.Sprintf("Not implemented: %s", r.URL.Path), http.StatusNotFound)
	}
}
