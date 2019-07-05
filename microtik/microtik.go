package microtik

import (
	"fmt"

	"gopkg.in/routeros.v2"
)

// Microtik is a struct holding the parameters and the connection to a
// microtik device.
type Microtik struct {
	*routeros.Client
	config   Config
	routeIDs map[string]string
}

// Config is a struct which contains the configuration parameters
// which are needed to connect to your microtik device
type Config struct {
	Address  string
	Port     int
	Username string
	Password string
}

// RouteResult is type used to return route results.
type RouteResult map[string]bool

// Option is a function argument type for the Microtik constructor
type Option func(m *Microtik)

// New is a constructor method and returns an initalized (but not yet connected)
// instance of a Microtik object
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

// Reset4G cuts the power of a USB LTE/4G modem connected to the routerboard
// for a period of 5 seconds.
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

// RouteStatus allows to query this status of a particular route on a
// microtik device (ip/route). The corresponing route must be registered during
// construction of the Microtik object, otherwise this method will fail.
// This method returns three parameters:
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

// SetRoute allows to set parameters for a particular route on a microtik
// device (ip/route). The corresponding route must be registered during
// construction of the Microtik object, otherwise this method will fail.
func (m *Microtik) SetRoute(name, command string) error {

	rComment, ok := m.routeIDs[name]
	if !ok {
		return fmt.Errorf("unknown route %s", name)
	}

	if err := m.connect(); err != nil {
		return err
	}

	reply, err := m.Run("/ip/route/print")
	if err != nil {
		return err
	}

	route, err := getRoute(reply, rComment)

	if route == nil {
		return fmt.Errorf("unable to find route %s", name)
	}

	nr, ok := route[".id"]
	if !ok {
		return fmt.Errorf("unable to determine the id of route %s", name)
	}

	query := []string{"/ip/route/set", "=.id=" + nr, "=" + command}

	reply, err = m.RunArgs(query)

	if err != nil {
		return err
	}

	// fmt.Println("reply:", reply)

	return nil
}

// getRoute retrieves the parameters for a route from a routeros.Reply
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
