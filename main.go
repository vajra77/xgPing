package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
	"xgPing/probe"
)

func parseUrl(url string) ([]*probe.Peer, error) {
	fmt.Printf("url: %s\n", url)
	result := make([]*probe.Peer, 0)
	result = append(result, probe.NewPeer("Namex", "namex", "193.201.28.100",
		"2001:7f8:10::2:4796"))
	result = append(result, probe.NewPeer("AS112", "as112", "193.201.28.112",
		"2001:7f8:10::112"))
	return result, nil
}

func main() {

	// parse command line arguments
	url := flag.String("json", "", "JSON IXP-F File")
	count := flag.Int("count", 10, "Number of ICMP pings to send")
	flag.Parse()

	// retrieve peers from json file
	peers, _ := parseUrl(*url)

	for {
		// main peers loop
		wg := sync.WaitGroup{}

		for _, peer := range peers {
			wg.Add(1)
			go peer.Ping(*count, &wg)
		}

		wg.Wait()

		time.Sleep(30 * time.Second)

		for _, peer := range peers {
			last := peer.LastSample()
			fmt.Printf("=== Last Statistics for peer: %s (%s) ===\n", peer.Name(), peer.V4Address())
			fmt.Printf("RTT (ms) min: %.2f, max: %.2f, avg: %.2f, dev: %.2f | LOSS: %.2f\n",
				last.Min(), last.Max(), last.Avg(), last.StdDev(), last.Loss())
		}
	}
}
