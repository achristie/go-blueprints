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
)

func main() {
	addr := flag.String("addr", ":8081", "address of app")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.Handle("/login", &templateHandler{filename: "login.html"})

	http.Handle("/auth/login/google")
	http.Handle("/auth/callback/google")
	http.Handle("/auth/login/github")
	http.Handle("/auth/callback/github")

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
	t.templ.Execute(w, r)
}
