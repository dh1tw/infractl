package app

func (s *Server) routes() {
	s.router.HandleFunc("/api/v1.0/reset4g", s.reset4g).Methods("GET")
	s.router.HandleFunc("/api/v1.0/check-connectivity", s.reset4g).Methods("GET")
}
