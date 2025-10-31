package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sync"
	"xgPing/probe"
)

func importCSVPeers(filename string) ([]*probe.Peer, error) {
	result := make([]*probe.Peer, 0)

	file, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("[W] Error closing file: %s\n", err)
		}
	}(file)

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

	wg := sync.WaitGroup{}

	for _, peer := range peers {
		wg.Add(1)
		go peer.Ping(*count, &wg)
	}

	wg.Wait()
	fmt.Println("PEER;NODE;IPV4;COUNT;RTT_MIN;RTT_MAX;RTT_AVG;RTT_STDEV;LOSS;STATUS")
	for _, peer := range peers {
		last := peer.LastSample()
		status := "OK"
		if last.Loss() > 0 {
			status = "WARN"
		}
		fmt.Printf("%s;%s;%s;%d;%.2f;%.2f;%.2f;%.2f;%.2f;%s\n", peer.Name(), peer.Node(), peer.V4Address(), *count,
			last.Min(), last.Max(), last.Avg(), last.StdDev(), last.Loss(), status)
	}
}
