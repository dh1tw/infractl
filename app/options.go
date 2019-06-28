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
		s.config.microtik = m
	}
}

// Mf823Address is a functional option which sets the ip address of a
// a ZTE MF823 4G USB modem
func Mf823Address(a string) func(*Server) {
	return func(s *Server) {
		s.config.mf823Address = a
	}
}

// Mf823Parameters is a functional option which sets the parameters which
// will be retrieved from a ZTE MF823 4G USB modem
func Mf823Parameters(p []string) func(*Server) {
	return func(s *Server) {
		s.config.mf823Parameters = p
	}
}
