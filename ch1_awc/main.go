package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"sync"
)

func main() {
	room := newRoom()
	// room.run()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", room)
	go room.run()
	log.Fatal(http.ListenAndServe(":8082", nil))
}

type templateHandler struct {
	once     sync.Once
	filename string
	tmpl     *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.tmpl = template.Must(template.ParseFiles(path.Join("templates", t.filename)))
	})
	t.tmpl.Execute(w, r)
}
