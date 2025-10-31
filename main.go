package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
	"xgPing/probe"
)

func importDummyPeers() ([]*probe.Peer, error) {
	result := make([]*probe.Peer, 0)
	result = append(result, probe.NewPeer("Namex", "namex", "193.201.28.100",
		"2001:7f8:10::2:4796"))
	result = append(result, probe.NewPeer("AS112", "as112", "193.201.28.112",
		"2001:7f8:10::112"))
	return result, nil
}

func importCSVPeers(filename string) ([]*probe.Peer, error) {
	result := make([]*probe.Peer, 0)

	file, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return result, err
	}
	for _, r := range records {
		peer := probe.NewPeer(r[0], r[1], r[2], r[3])
		result = append(result, peer)
	}

	return result, nil
}

func main() {

	// parse command line arguments
	csvFile := flag.String("csv", "peers.csv", "Peer list in CSV format")
	count := flag.Int("count", 10, "Number of ICMP pings to send")
	flag.Parse()

	// retrieve peers from json file
	peers, err := importCSVPeers(*csvFile)
	if err != nil {
		fmt.Printf("Unable to import CSV file: %s", err)
		os.Exit(1)
	}

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
			fmt.Printf("RTT min: %.2f, max: %.2f, avg: %.2f, dev: %.2f ms | LOSS: %.2f %%\n",
				last.Min(), last.Max(), last.Avg(), last.StdDev(), last.Loss())
		}
	}
}
