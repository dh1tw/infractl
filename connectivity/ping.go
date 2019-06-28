package connectivity

import (
	"fmt"
	"time"

	goping "github.com/sparrc/go-ping"
)

// Ping will send a single ping to the specified IP address or URL and return
// the round trip time. In case the IP address is unreachable, an error
// will be returned after 2 seconds.
// In order to execute this command you might need elevated privileges on Linux.
// See: https://github.com/sparrc/go-ping
func Ping(address string) (time.Duration, error) {
	timeout := time.NewTimer(time.Second * 2).C
	avgPing := time.Second * 2

	pinger, err := goping.NewPinger(address)
	if err != nil {
		return avgPing, err
	}

	result := make(chan (*goping.Statistics))

	go func() {
		pinger.Count = 1
		pinger.Run() // blocks until finished
		result <- pinger.Statistics()
	}()

	select {
	case <-timeout:
		return avgPing, fmt.Errorf("no reply received from %s", address)
	case s := <-result:
		avgPing = s.AvgRtt
	}

	return avgPing, nil
}
