package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/stgnet/pperf"
)

func main() {

	target := flag.String("target", "stg.net", "test with this server")
	server := flag.Bool("server", false, "Start server")
	port := flag.Int("port", 5048, "Port to listen/send to")
	ifname := flag.String("interface", "", "bind to interface (linux only)")
	seconds := flag.Int("seconds", 5, "Seconds to run test")

	flag.Parse()

	results := pperf.Pperf(pperf.API{
		Server:    *server,
		Target:    *target,
		Port:      *port,
		Seconds:   *seconds,
		Interface: *ifname,
	})

	pretty, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(pretty))
}
