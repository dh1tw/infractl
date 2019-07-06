package connectivity

import (
	"fmt"
	"sync"
	"time"

	goping "github.com/sparrc/go-ping"
)

// PingResults contains a map with the ping Result for one or more hosts
type PingResults map[string]PingResult

// PingResult is a struct containing the result of a ping to one particular host
type PingResult struct {
	Address string        `json:"address"`
	RTT     time.Duration `json:"rtt"`
	Failed  bool          `json:"failed"`
}

// PingHost will send a single ping to the specified IP address or URL and return
// the round trip time. In case the IP address is unreachable, an error
// will be returned after 2 seconds.
// In order to execute this command you might need elevated privileges on Linux.
// See: https://github.com/sparrc/go-ping
func PingHost(address string) (time.Duration, error) {
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

// PingHosts will send a single ping to the specified IP addresses or URLs.
// In case the IP address is unreachable, an error will be returned after
// 2 seconds.
// In order to execute this command you might need elevated privileges on Linux.
// See: https://github.com/sparrc/go-ping
func PingHosts(addresses []string) PingResults {

	resultCh := make(chan PingResult)

	wg := &sync.WaitGroup{}

	for _, addr := range addresses {
		wg.Add(1)
		go pingAsync(addr, wg, resultCh)
	}

	results := make(PingResults)

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for res := range resultCh {
		results[res.Address] = res
	}

	return results
}

func pingAsync(address string, wg *sync.WaitGroup, resCh chan<- PingResult) {
	defer wg.Done()
	var res PingResult
	ping, err := PingHost(address)
	if err != nil {
		res = PingResult{address, time.Second * 0, true}
	} else {
		res = PingResult{address, ping, false}
	}
	resCh <- res
}

func (r PingResult) String() string {
	if r.Failed {
		return fmt.Sprintf("%s: failed", r.Address)
	}
	return fmt.Sprintf("%s: %v", r.Address, r.RTT)
}
