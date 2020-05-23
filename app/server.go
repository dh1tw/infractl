package webserver

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/dh1tw/infractl/connectivity"
	"github.com/dh1tw/infractl/microtik"
	"github.com/markbates/pkger"

	"github.com/gorilla/mux"
)

type Server struct {
	sync.Mutex
	router          *mux.Router
	address         string
	port            int
	fileServer      http.Handler
	apiVersion      string
	apiMatch        *regexp.Regexp
	errorCh         chan struct{}
	closeOnce       sync.Once
	microtik        *microtik.Microtik
	mf823Address    string
	mf823Parameters []string
	pingEnabled     bool
	pingInterval    time.Duration
	pingHosts       []string
	pingResults     connectivity.PingResults
	services        map[string]struct{}
	mtRoutes        []string
}

// Option is the type used for functional options
type Option func(*Server)

// New returns an instance of a Server configured according to the provided options.
func New(opts ...Option) *Server {

	s := &Server{
		address:         "localhost",
		port:            6556,
		apiVersion:      "1.0",
		apiMatch:        regexp.MustCompile(`api\/v\d\.\d\/`),
		router:          mux.NewRouter().StrictSlash(true),
		pingHosts:       []string{},
		pingEnabled:     false,
		pingInterval:    time.Duration(time.Second * 10),
		pingResults:     make(connectivity.PingResults),
		mf823Parameters: []string{},
		errorCh:         make(chan struct{}),
		services:        make(map[string]struct{}),
		mtRoutes:        []string{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) close() {
	s.closeOnce.Do(func() {
		close(s.errorCh)
	})
}

// Serve starts a HTTP Server on a given network adapter / port and
// sets a HTTP handler.
// Since this function contains an endless loop, it should be executed
// in a go routine. If the listener can not be initialized, it will
// close the errorCh channel.
func (s *Server) Serve() {

	defer s.close()

	if s.pingEnabled {
		log.Printf("start pinging hosts in %v interval\n", s.pingInterval)
		go s.startPing(s.pingInterval)
	}

	url := fmt.Sprintf("%s:%d", s.address, s.port)

	s.fileServer = http.FileServer(pkger.Dir("/web/dist"))

	s.routes()

	// Listen for incoming connections.
	log.Printf("listening on %s for HTTP connections\n", url)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         url,
		Handler:      s.apiRedirectRouter((s.router)),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Server) startPing(interval time.Duration) {

	ticker := time.NewTicker(interval).C

	for {
		s.Lock()
		hosts := s.pingHosts
		s.Unlock()

		res := connectivity.PingHosts(hosts)
		s.Lock()
		s.pingResults = res
		s.Unlock()
		// make sure the hosts are pinged immediately after startup
		<-ticker
	}
}
