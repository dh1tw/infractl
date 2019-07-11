package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dh1tw/infractl/microtik"

	"github.com/dh1tw/infractl/services"

	"github.com/dh1tw/infractl/connectivity"
	"github.com/gorilla/mux"

	"github.com/dh1tw/infractl/mf823"
)

func (s *Server) handlePing(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.Lock()
	res := s.pingResults
	enabled := s.pingEnabled
	hosts := s.pingHosts
	s.Unlock()

	// execute the ping if the background pinging has not been enabled
	if !enabled {
		res = connectivity.PingHosts(hosts)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to encode ping data to json"))
	}
}

// handleReset4G the 4G modem attached to the Microtik Router (needs to be supported
// the routerboard)
func (s *Server) handleReset4G(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.Lock()
	defer s.Unlock()

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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.Lock()
	addr := s.mf823Address
	params := s.mf823Parameters
	s.Unlock()

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

func (s *Server) handleServiceRestart(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.Lock()
	defer s.Unlock()

	vars := mux.Vars(req)
	sName := strings.ToLower(vars["service"])
	sName = strings.Replace(sName, ".service", "", 1)

	if _, ok := s.services[sName]; !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("unauthorized service")))
		return
	}

	if err := services.Restart(sName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(err.Error())))
		return
	}
}

func (s *Server) handleRoutes(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.Lock()
	defer s.Unlock()

	if s.microtik == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("no microtik instance configured"))
		return
	}

	results := make(map[string]microtik.RouteResult)

	for _, route := range s.mtRoutes {
		res, err := s.microtik.RouteStatus(route)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		results[route] = res
	}

	j, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

	}
	w.Write(j)
}

func (s *Server) handleRoute(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.Lock()
	defer s.Unlock()

	if s.microtik == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("no microtik instance configured"))
		return
	}

	vars := mux.Vars(req)
	rName := strings.ToLower(vars["route"])

	res, err := s.microtik.RouteStatus(rName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	j, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

	}
	w.Write(j)
}

func (s *Server) handleRouteEnable(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.Lock()
	defer s.Unlock()

	vars := mux.Vars(req)
	rName := strings.ToLower(vars["route"])

	err := s.microtik.SetRoute(rName, "disabled=false")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (s *Server) handleRouteDisable(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.Lock()
	defer s.Unlock()

	vars := mux.Vars(req)
	rName := strings.ToLower(vars["route"])

	err := s.microtik.SetRoute(rName, "disabled=true")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
