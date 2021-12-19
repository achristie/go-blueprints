package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/achristie/go-blueprints/ch1/trace"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/stretchr/objx"
)

func main() {
	addr := flag.String("addr", ":8081", "address of app")
	secret := flag.String("secret", "test", "client secret for gh")
	flag.Parse()

	goth.UseProviders(
		github.New("7dbca52d1333166da68e", *secret, "http://localhost:3000/auth/github/callback"),
	)
	gothic.GetProviderName = func(req *http.Request) (string, error) {
		return "github", nil
	}

	r := newRoom(UseGravatar)
	r.tracer = trace.New(os.Stdout)

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.Handle("/avatars/",
		http.StripPrefix("/avatars", http.FileServer(http.Dir("./avatars"))))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/auth/", loginHandler)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)
	go r.run()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}
