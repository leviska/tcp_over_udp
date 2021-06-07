package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/leviska/tcp-over-udp/mock"
	"github.com/leviska/tcp-over-udp/tcp"
	"github.com/leviska/tcp-over-udp/udp"
	"github.com/leviska/tcp-over-udp/util"
)

func runTCPServer() {
	util.Logger.Info("running server")
	server, err := tcp.NewServer(udp.NewAddr(*ip, *port), mock.NewPipeline(mock.NewLogging()))
	if err != nil {
		util.Logger.Panic(err)
	}
	defer server.Close()
	util.Logger.Infof("listening on %s:%d", *ip, *port)

	server.Run()
	clients := server.Connections()

	for c := range clients {
		go func(c *tcp.Connection) {
			for !c.IsClosed() {
				d, err := c.Receive()
				if err != nil {
					util.Logger.Error(err)
					break
				}
				err = c.Send(d)
				if err != nil {
					util.Logger.Error(err)
				}
			}
		}(c)
	}
}

func runTCPClient() {
	fmt.Println("running client")
	fmt.Println("enter text to send:")
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", *ip, *port))
	if err != nil {
		panic(err)
	}
	conn, err := tcp.NewConnection(addr, mock.NewPipeline(mock.NewLogging()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	
	var file io.WriteCloser
	if *output != "" {
		var err error
		file, err = os.Create(*output)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := conn.Send(scanner.Bytes())
		if err != nil {
			panic("couldn't send message, error: " + err.Error())
		}
		buf, err := conn.Receive()
		if err != nil {
			panic("expected message from server, got error: " + err.Error())
		}
		fmt.Printf("got message %q\n", util.Truncate(string(buf), 40))
		if file != nil {
			file.Write(buf)
			file.Write([]byte{'\n'})
		}
	}
}
