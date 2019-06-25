package microtik

import (
	"fmt"

	"gopkg.in/routeros.v2"
)

type Microtik struct {
	*routeros.Client
	config Config
}

type Config struct {
	Address  string
	Port     int
	Username string
	Password string
}

func New(c Config) *Microtik {

	m := &Microtik{
		config: c,
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
