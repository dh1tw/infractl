package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dh1tw/infractl/microtik"

	"github.com/gorilla/mux"
)

type Server struct {
	// db     *someDatabase
	router  *mux.Router
	mt      *microtik.Microtik
	address string
	port    int
}

// Option is the type used for functional options
type Option func(*Server)

// New returns an instance of a Server configured according to the provided options.
func New(opts ...Option) *Server {

	mtConfig := microtik.Config{
		Address:  "192.168.0.1",
		Port:     8728,
		Username: "admin",
		Password: "admin",
	}

	s := &Server{
		address: "localhost",
		port:    6556,
		router:  mux.NewRouter(),
		mt:      microtik.New(mtConfig),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Init() error {

	return nil
}

// reset the 4G modem attached to the Microtik Router (needs to be supported
// the routerboard)
func (s *Server) reset4g(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	err := s.mt.Reset4G()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (s *Server) checkConnectivity(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	err := s.mt.Reset4G()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

// ListenHTTP starts a HTTP Server on a given network adapter / port and
// sets a HTTP handler.
// Since this function contains an endless loop, it should be executed
// in a go routine. If the listener can not be initialized, it will
// close the errorCh channel.
func (s *Server) ListenHTTP(errorCh chan<- struct{}) {

	defer close(errorCh)

	// load the HTTP routes with their respective endpoints
	s.routes()

	url := fmt.Sprintf("%s:%d", s.address, s.port)

	// Listen for incoming connections.
	log.Printf("listening on %s for HTTP connections\n", url)

	err := http.ListenAndServe(url, s.router)
	if err != nil {
		log.Println(err)
		return
	}
}
