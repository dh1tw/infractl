package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-systemd/dbus"
)

type UnitStatus struct {
	Name        string // The primary unit name as string
	Description string // The human readable description string
	ActiveState string // The active state (i.e. whether the unit is currently started or not)
	LoadState   string // The load state (i.e. whether the unit file has been loaded successfully)
	SubState    string // The sub state (a more fine-grained version of the active state that is specific to the unit type, which the active state is not)
}

func Status(name ...string) ([]UnitStatus, error) {

	conn, err := dbus.New()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	for i := 0; i < len(name); i++ {
		if !strings.Contains(name[i], ".service") {
			name[i] = name[i] + ".service"
		}
	}

	stats, err := conn.ListUnitsByNames(name)
	if err != nil {
		return nil, err
	}

	res := []UnitStatus{}

	for _, s := range stats {
		us := UnitStatus{
			Name:        s.Name,
			Description: s.Description,
			ActiveState: s.ActiveState,
			LoadState:   s.LoadState,
			SubState:    s.SubState,
		}
		res = append(res, us)
	}

	return res, nil
}

func (us UnitStatus) String() string {
	return fmt.Sprintf("Name: %s\n Description: %s\n State: %s\n LoadState: %s\n SubState: %s\n", us.Name, us.Description, us.ActiveState, us.LoadState, us.SubState)
}

// Start a systemd service
func Start(name string) error {
	conn, err := dbus.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	resultCh := make(chan string)

	if !strings.Contains(name, ".service") {
		name = name + ".service"
	}

	// The mode needs to be one of replace, fail, isolate, ignore-dependencies,
	// ignore-requirements. If "replace" the call will start the unit and its
	// dependencies, possibly replacing already queued jobs that conflict with
	// this. If "fail" the call will start the unit and its dependencies,
	// but will fail if this would change an already queued job. If "isolate"
	// the call will start the unit in question and terminate all units that
	// aren't dependencies of it. If "ignore-dependencies" it will start a
	// unit but ignore all its dependencies. If "ignore-requirements" it will
	//  start a unit but only ignore the requirement dependencies. It is not
	// recommended to make use of the latter two options. Returns the newly
	// created job object.
	job, err := conn.StartUnit(name, "replace", resultCh)
	if err != nil {
		return err
	}

	if job == 0 {
		return fmt.Errorf("service start could not be scheduled")
	}

	timeout := time.NewTimer(time.Second * 20).C

	select {
	case res := <-resultCh:
		if strings.Compare(res, "done") != 0 {
			return fmt.Errorf("service %s could not be started", name)
		}
	case <-timeout:
		return fmt.Errorf("timeout, received no feedback from systemd")
	}

	return nil
}

// Stop a systemd service
func Stop(name string) error {
	conn, err := dbus.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	resultCh := make(chan string)

	if !strings.Contains(name, ".service") {
		name = name + ".service"
	}

	// The mode needs to be one of replace, fail, isolate, ignore-dependencies,
	// ignore-requirements. If "replace" the call will start the unit and its
	// dependencies, possibly replacing already queued jobs that conflict with
	// this. If "fail" the call will start the unit and its dependencies,
	// but will fail if this would change an already queued job. If "isolate"
	// the call will start the unit in question and terminate all units that
	// aren't dependencies of it. If "ignore-dependencies" it will start a
	// unit but ignore all its dependencies. If "ignore-requirements" it will
	//  start a unit but only ignore the requirement dependencies. It is not
	// recommended to make use of the latter two options. Returns the newly
	// created job object.
	job, err := conn.StopUnit(name, "replace", resultCh)
	if err != nil {
		return err
	}

	if job == 0 {
		return fmt.Errorf("service stop could not be scheduled")
	}

	timeout := time.NewTimer(time.Second * 20).C

	select {
	case res := <-resultCh:
		if strings.Compare(res, "done") != 0 {
			return fmt.Errorf("service %s could not be stopped", name)
		}
	case <-timeout:
		return fmt.Errorf("timeout, received no feedback from systemd")
	}

	return nil
}

// Restart will try to restart a systemd service
func Restart(name string) error {

	conn, err := dbus.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	resultCh := make(chan string)

	if !strings.Contains(name, ".service") {
		name = name + ".service"
	}

	// The mode needs to be one of replace, fail, isolate, ignore-dependencies,
	// ignore-requirements. If "replace" the call will start the unit and its
	// dependencies, possibly replacing already queued jobs that conflict with
	// this. If "fail" the call will start the unit and its dependencies,
	// but will fail if this would change an already queued job. If "isolate"
	// the call will start the unit in question and terminate all units that
	// aren't dependencies of it. If "ignore-dependencies" it will start a
	// unit but ignore all its dependencies. If "ignore-requirements" it will
	//  start a unit but only ignore the requirement dependencies. It is not
	// recommended to make use of the latter two options. Returns the newly
	// created job object.
	job, err := conn.ReloadOrRestartUnit(name, "replace", resultCh)
	if err != nil {
		return err
	}

	if job == 0 {
		return fmt.Errorf("service restart could not be scheduled")
	}

	timeout := time.NewTimer(time.Second * 20).C

	select {
	case res := <-resultCh:
		if strings.Compare(res, "done") != 0 {
			return fmt.Errorf("service %s could not be restarted", name)
		}
	case <-timeout:
		return fmt.Errorf("timeout, received no feedback from systemd")
	}

	return nil
}
