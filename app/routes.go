package webserver

func (s *Server) routes() {
	s.router.HandleFunc("/api/v1.0/reset4g", s.handleReset4G).Methods("GET")
	s.router.HandleFunc("/api/v1.0/status4g", s.handleStatus4G).Methods("GET")
	s.router.HandleFunc("/api/v1.0/ping", s.handlePing).Methods("GET")
}
