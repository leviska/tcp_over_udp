package main

import (
	"bufio"
	"fmt"
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
			util.Logger.Info("got client")
			for !c.IsClosed() {
				d, err := c.Receive()
				util.Logger.Info("got mesage: ", string(d))
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
	conn, err := tcp.NewConnection(udp.NewAddr(*ip, *port), mock.NewPipeline(mock.NewLogging()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	go func() {
		for {
			buf, err := conn.Receive()
			if err != nil {
				panic("expected message from server, got error: " + err.Error())
			}
			fmt.Printf("got message %q\n", string(buf))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := conn.Send(scanner.Bytes())
		fmt.Println("sent message")
		if err != nil {
			panic("couldn't send message, error: " + err.Error())
		}
	}
}
