package connectivity

import (
	"fmt"
	"log"
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

// PingHost will send a specified amount of pings to the specified IP address or URL and return
// the average round trip time. In case the IP address is unreachable, an error will be returned
// after the provided timeout.
// In order to execute this command you might need elevated privileges on Linux.
// See: https://github.com/sparrc/go-ping for more details.
func PingHost(address string, timeout time.Duration, samples int) (PingResult, error) {
	timeoutC := time.NewTimer(timeout).C

	pr := PingResult{
		Address: address,
		RTT:     0,
		Failed:  true,
	}

	pinger, err := goping.NewPinger(address)
	if err != nil {
		return pr, err
	}

	pinger.SetPrivileged(true)

	result := make(chan (*goping.Statistics))

	go func() {
		pinger.Count = samples
		pinger.Run() // blocks until finished
		result <- pinger.Statistics()
	}()

	select {
	case <-timeoutC:
		return pr, fmt.Errorf("no reply received from %s after %v", address, timeout)
	case s := <-result:
		pr.RTT = s.AvgRtt
		pr.Failed = false
	}

	return pr, nil
}

// PingHosts will send a specified amount of pings to the provided list of IP addresses or URLs
// and return the average round trip time. In case the IP address is unreachable, an error will
// be returned after the provided timeout.
// In order to execute this command you might need elevated privileges on Linux.
// See: https://github.com/sparrc/go-ping for more details.
func PingHosts(addresses []string, timeout time.Duration, samples int) PingResults {

	resultCh := make(chan PingResult)

	wg := &sync.WaitGroup{}

	for _, addr := range addresses {
		wg.Add(1)
		go pingAsync(addr, wg, resultCh, timeout, samples)
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

func pingAsync(address string, wg *sync.WaitGroup, resCh chan<- PingResult, timeout time.Duration, samples int) {
	defer wg.Done()
	var res PingResult
	res, err := PingHost(address, timeout, samples)
	if err != nil {
		log.Println(err)
	}
	resCh <- res
}

func (r PingResult) String() string {
	if r.Failed {
		return fmt.Sprintf("%s: failed", r.Address)
	}
	return fmt.Sprintf("%s: %v", r.Address, r.RTT)
}
