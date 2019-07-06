package webserver

func (s *Server) routes() {
	s.router.HandleFunc("/api/v1.0/reset4g", s.handleReset4G)
	s.router.HandleFunc("/api/v1.0/status4g", s.handleStatus4G)
	s.router.HandleFunc("/api/v1.0/ping", s.handlePing)
	s.router.HandleFunc("/api/v1.0/service/{service}/restart", s.handleServiceRestart)
	s.router.HandleFunc("/api/v1.0/routes", s.handleRoutes)
	s.router.HandleFunc("/api/v1.0/route/{route}", s.handleRoute)
	s.router.HandleFunc("/api/v1.0/route/{route}/enable", s.handleRouteEnable)
	s.router.HandleFunc("/api/v1.0/route/{route}/disable", s.handleRouteDisable)
}
