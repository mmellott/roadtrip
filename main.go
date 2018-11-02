package main

import (
	"flag"
)

type Config struct {
	udp     bool
	address string
	port    string
	size    int
	tick    int
}

func main() {
	var listen bool
	flag.BoolVar(&listen, "listen", false, "operate as echo server")

	var config Config
	flag.BoolVar(&config.udp, "udp", false, "dial UDP instead of TCP")
	flag.StringVar(&config.address, "address", "localhost", "address to dial")
	flag.StringVar(&config.port, "port", "5950", "port to listen on or dial")
	flag.IntVar(&config.size, "size", 1024, "size of packets")

	flag.Parse()

	if listen {
		server(&config)
	} else {
		client(&config)
	}
}
