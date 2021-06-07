package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/leviska/tcp-over-udp/udp"
	"github.com/leviska/tcp-over-udp/util"
)

func runUDPServer() {
	util.Logger.Info("running server")
	server, err := udp.NewServer(udp.NewAddr(*ip, *port))
	if err != nil {
		util.Logger.Panic(err)
	}
	util.Logger.Infof("listening on %s:%d", *ip, *port)

	for {
		conn, err := server.Connection()
		go func() {
			if err != nil {
				util.Logger.Error(err)
				return
			}
			defer conn.Close()
			util.Logger.Infof("new connection from %q", conn.Addr())
			for {
				buf, err := conn.Receive()
				if err != nil {
					util.Logger.Error(err)
					return
				}
				util.Logger.Infof("message from %q: %q", conn.Addr(), string(buf))
				err = conn.Send(buf)
				if err != nil {
					util.Logger.Error(err)
					return
				}
			}
		}()
		break
	}
}

func runUDPClient() {
	fmt.Println("running client")
	fmt.Println("enter text to send:")
	conn, err := udp.NewConnection(udp.NewAddr(*ip, *port))
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
		if err != nil {
			panic("couldn't send message, error: " + err.Error())
		}
	}
}
