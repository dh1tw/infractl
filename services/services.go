package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-systemd/dbus"
)

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

	timeout := time.NewTimer(time.Second * 5).C

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
