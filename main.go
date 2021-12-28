package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/achristie/go-blueprints/meander"
)

func main() {
	addr := flag.String("addr", ":8080", "port to listen on")
	flag.Parse()

	server := NewServer()
	server.routes()

	http.HanndleFunc("/recommendations", cors(func(w http.ResponseWriter, r *http.Request) {
		q := &meander.Query{
			Journey: strings.Split(r.URL.Query().Get("journey"), "|"),
		}
		var err error
		q.Lat, err = strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		q.Lng, err = strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		q.Radius, err = strconv.Atoi(r.URL.Query().Get("radius"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		q.CostRangeStr = r.URL.Query().Get("cost")
		places := q.Run()
		respond(w, r, places)
	}))

	log.Println("Listening on port", *addr)
	if err := http.ListenAndServe(*addr, server.router); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello!")
	}
}

func (s *Server) handleAbout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "about")
	}
}

type Server struct {
	router *http.ServeMux
}

func NewServer() *Server {
	return &Server{router: http.NewServeMux()}
}
