package app

import (
	"github.com/dh1tw/infractl/microtik"
)

// Address is a functional option to set the address of the webserver
func Address(addr string) func(*Server) {
	return func(s *Server) {
		s.address = addr
	}
}

// Port is a functional option to set the port on which the webserver
// will listen
func Port(port int) func(*Server) {
	return func(s *Server) {
		s.port = port
	}
}

// Microtik is a functional option which sets an instance of a Microtik
// router
func Microtik(m *microtik.Microtik) func(*Server) {
	return func(s *Server) {
		s.mt = m
	}
}
