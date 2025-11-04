package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"xgPing/probe"
)

func main() {

	// parse command line arguments
	input := flag.String("input", "peers.csv", "Input source for importing peers' data; may be either an url pointing to a JSON resource or local path of CSV file")
	output := flag.String("output", "result.csv", "Output filename")
	ixpId := flag.Int("ixp-id", 1, "Import peers from IXP ID")
	vlanId := flag.Int("vlan-id", 1, "Import peers from VLAN ID")
	count := flag.Int("count", 10, "Number of ICMP pings to send")
	ttl := flag.Int("ttl", 1, "TTL for ICMP pings")

	flag.Parse()

	var peers []*probe.Peer
	var err error

	if strings.Contains(*input, "http") {
		peers, err = ImportJSONPeers(*input, *ixpId, *vlanId)
	} else {
		peers, err = ImportCSVPeers(*input)
	}
	if err != nil {
		fmt.Printf("Unable to import peers from input: %s\n", err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}

	for _, peer := range peers {
		wg.Add(1)
		fmt.Printf("Pinging peer: %s\n", peer.Name())
		go peer.Ping(*count, *ttl, &wg)
		time.Sleep(300 * time.Millisecond)
	}

	wg.Wait()

	outFile, err := os.Create(*output)

	if err != nil {
		fmt.Printf("Unable to create output file: %s", err)
		os.Exit(1)
	}

	_, err = fmt.Fprintln(outFile, "PEER;NODE;IPV4;COUNT;RTT_MIN;RTT_MAX;RTT_AVG;RTT_STDEV;LOSS;STATUS")

	for _, peer := range peers {
		last := peer.LastSample()
		status := "OK"
		if last.Loss() > 0 {
			status = "WARN"
		}
		_, err = fmt.Fprintf(outFile, "%s;%s;%s;%d;%.2f;%.2f;%.2f;%.2f;%.2f;%s\n", peer.Name(), peer.Node(), peer.V4Address(), *count,
			last.Min(), last.Max(), last.Avg(), last.StdDev(), last.Loss(), status)
	}
}
