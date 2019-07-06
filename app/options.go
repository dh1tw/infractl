package webserver

import (
	"strings"
	"time"

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

// ErrorCh is a functional option which provides channel to the webserver
// which will be closed if an unrecoverable error happens
func ErrorCh(ch chan struct{}) func(*Server) {
	return func(s *Server) {
		s.errorCh = ch
	}
}

// Microtik is a functional option which sets an instance of a Microtik
// router
func Microtik(m *microtik.Microtik) func(*Server) {
	return func(s *Server) {
		s.microtik = m
	}
}

// Mf823Address is a functional option which sets the ip address of a
// a ZTE MF823 4G USB modem
func Mf823Address(a string) func(*Server) {
	return func(s *Server) {
		s.mf823Address = a
	}
}

// Mf823Parameters is a functional option which sets the parameters which
// will be retrieved from a ZTE MF823 4G USB modem
func Mf823Parameters(p []string) func(*Server) {
	return func(s *Server) {
		s.mf823Parameters = p
	}
}

// PingAddress sets the hosts to be pinged
func PingAddress(addresses []string) func(*Server) {
	return func(s *Server) {
		s.pingHosts = addresses
	}
}

// PingEnabled enables the background job to ping the provided hosts
// in the defined interval. Otherwise the hosts will only be pinged
// on the execution of the corresponing http call
func PingEnabled(enabled bool) func(*Server) {
	return func(s *Server) {
		s.pingEnabled = enabled
	}
}

// PingInterval sets the interval in which the provided hosts will be pinged
func PingInterval(interval time.Duration) func(*Server) {
	return func(s *Server) {
		s.pingInterval = interval
	}
}

// Service authorizes the webserver to control a systemd service. Services can
// either be specified with or without the extension ".service"
func Service(serviceName string) func(*Server) {
	return func(s *Server) {
		sName := strings.ToLower(serviceName)
		sName = strings.Replace(sName, ".service", "", 1)
		s.services[sName] = struct{}{}
	}
}

// Route authorizes the webserver to access a microtik route
func Route(routeName string) func(*Server) {
	return func(s *Server) {
		rName := strings.ToLower(routeName)
		s.mtRoutes = append(s.mtRoutes, rName)
	}
}
