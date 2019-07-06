package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dh1tw/infractl/connectivity"

	"github.com/dh1tw/infractl/mf823"
)

func (s *Server) handlePing(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	s.RLock()
	res := s.pingResults
	enabled := s.pingEnabled
	hosts := s.pingHosts
	s.RUnlock()

	// execute the ping if the background pinging has not been enabled
	if !enabled {
		res = connectivity.PingHosts(hosts)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to encode ping data to json"))
	}
}

// handleReset4G the 4G modem attached to the Microtik Router (needs to be supported
// the routerboard)
func (s *Server) handleReset4G(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	s.RLock()
	defer s.RUnlock()

	if s.microtik == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("no microtik instance configured"))
		return
	}

	err := s.microtik.Reset4G()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

//retrieve the requested status from a ZTE MF823 4G modem
func (s *Server) handleStatus4G(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	s.RLock()
	addr := s.mf823Address
	params := s.mf823Parameters
	s.RUnlock()

	if len(addr) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("address of mf823 not set or no query parameters provided")))
		return
	}

	resp, err := mf823.Status(addr, params...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	j, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

	}
	w.Write(j)
}
