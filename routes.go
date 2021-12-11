package main

func (s *Server) routes() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/about", s.handleAbout())
}
