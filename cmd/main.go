package main

import (
	"flag"

	"github.com/leviska/tcp-over-udp/util"
)

var (
	port = flag.Int("port", 1337, "port to run the server on/to connect client to")
	ip   = flag.String("ip", "127.0.0.1", "ip to run the server on/to connect client to")
	mode = flag.String("mode", "client", "(client|server) which mode to run")
)

func main() {
	flag.Parse()
	util.SetupLogger()

	if *mode == "server" {
		runTCPServer()
	} else if *mode == "client" {
		runTCPClient()
	} else {
		flag.PrintDefaults()
	}
}
