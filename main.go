package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8080", "port to listen on")
	flag.Parse()

	server := NewServer()
	server.routes()

	log.Println("Listening on port", *addr)
	if err := http.ListenAndServe(*addr, server.router); err != nil {
		log.Fatal("ListenAndServe:", err)
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
