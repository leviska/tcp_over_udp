package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/leviska/tcp-over-udp/udp"
	"go.uber.org/zap"
)

var (
	port   = flag.Int("port", 9000, "port to run the server on/to connect client to")
	ip     = flag.String("ip", "127.0.0.1", "ip to run the server on/to connect client to")
	mode   = flag.String("mode", "client", "(client|server) which mode to run")
	Logger *zap.SugaredLogger
)

func runServer() {
	Logger.Info("running server")
	server, err := udp.NewServer(udp.NewAddr(*ip, *port))
	if err != nil {
		Logger.Panic(err)
	}
	Logger.Infof("listening on %s:%d", *ip, *port)

	for {
		conn, err := server.Connection()
		go func() {
			if err != nil {
				Logger.Error(err)
				return
			}
			defer conn.Close()
			Logger.Infof("new connection from %q", conn.Addr())
			for {
				buf, err := conn.Receive()
				if err != nil {
					Logger.Error(err)
					return
				}
				Logger.Infof("message from %q: %q", conn.Addr(), string(buf))
				err = conn.Send(buf)
				if err != nil {
					Logger.Error(err)
					return
				}
			}
		}()
	}
}

func runClient() {
	fmt.Println("running client")
	fmt.Println("enter text to send:")
	conn, err := udp.NewConnection(udp.NewAddr(*ip, *port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
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
		fmt.Printf("got message %q\n", string(buf))
	}
}

func main() {
	flag.Parse()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	Logger = logger.Sugar()

	if *mode == "server" {
		runServer()
	} else if *mode == "client" {
		runClient()
	} else {
		flag.PrintDefaults()
	}
}
