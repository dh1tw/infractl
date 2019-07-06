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

	"github.com/gorilla/mux"
)

type Server struct {
	sync.RWMutex
	router          *mux.Router
	address         string
	port            int
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
		router:          mux.NewRouter(),
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

	s.routes()

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

	// Listen for incoming connections.
	log.Printf("listening on %s for HTTP connections\n", url)

	err := http.ListenAndServe(url, s.apiRedirectRouter(s.router))
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Server) startPing(interval time.Duration) {

	ticker := time.NewTicker(interval).C

	for {
		s.RLock()
		hosts := s.pingHosts
		s.RUnlock()

		res := connectivity.PingHosts(hosts)
		s.Lock()
		s.pingResults = res
		s.Unlock()
		// make sure the hosts are pinged immediately after startup
		<-ticker
	}
}
