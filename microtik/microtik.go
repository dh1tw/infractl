package microtik

import (
	"fmt"

	"gopkg.in/routeros.v2"
)

type Microtik struct {
	*routeros.Client
	config   Config
	routeIDs map[string]string
}

type Config struct {
	Address  string
	Port     int
	Username string
	Password string
}

type RouteResult map[string]bool

type Option func(m *Microtik)

func New(c Config, opts ...Option) *Microtik {

	m := &Microtik{
		config:   c,
		routeIDs: make(map[string]string),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Microtik) connect() error {
	url := fmt.Sprintf("%s:%d", m.config.Address, m.config.Port)
	c, err := routeros.Dial(url, m.config.Username, m.config.Password)
	if err != nil {
		return err
	}
	m.Client = c
	return nil
}

func (m *Microtik) Reset4G() error {

	if err := m.connect(); err != nil {
		return err
	}

	reply, err := m.Run("/system/routerboard/usb/power-reset", "?duration=5s")
	if err != nil {
		return err
	}

	if len(reply.Re) > 0 {
		return fmt.Errorf("Error: %v", reply)
	}

	m.Client.Close()
	m.Client = nil
	return nil
}

// RouteStatus allows to query this status of a particular ip/route on a
// microtik router. The routeID needs to be previously registered at construction
// of the Microtik object. The function returns three parameters:
// 1. route disabled (bool)
// 2. route active (bool)
// 3. error
func (m *Microtik) RouteStatus(name string) (RouteResult, error) {

	rComment, ok := m.routeIDs[name]
	if !ok {
		return nil, fmt.Errorf("unknown route %s", name)
	}

	if err := m.connect(); err != nil {
		return nil, err
	}

	reply, err := m.Run("/ip/route/print")
	if err != nil {
		return nil, err
	}

	route, err := getRoute(reply, rComment)

	if route == nil {
		return nil, fmt.Errorf("unable to find route %s", name)
	}

	disabled := false
	d, ok := route["disabled"]
	if !ok {
		return nil, fmt.Errorf("unable to determine parameter disabled of route %s", name)
	}
	if d == "true" {
		disabled = true
	}

	active := false
	a, ok := route["active"]
	if !ok {
		return nil, fmt.Errorf("unable to determine parameter active of route %s", name)
	}
	if a == "true" {
		active = true
	}

	res := RouteResult{
		"disabled": disabled,
		"active":   active,
	}

	return res, nil
}

func getRoute(reply *routeros.Reply, comment string) (map[string]string, error) {

	if reply == nil {
		return nil, fmt.Errorf("router response nil")
	}

	if len(reply.Re) == 0 {
		return nil, fmt.Errorf("router response empty")
	}

	var route map[string]string

	for _, r := range reply.Re {
		c, ok := r.Map["comment"]
		if !ok {
			continue
		}
		// the only way to identify a route
		// is comparing it's comment
		if c == comment {
			route = r.Map
			break
		}
	}
	return route, nil
}
