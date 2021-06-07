package tcp

import (
	"net"
	"sync"

	"github.com/leviska/tcp-over-udp/udp"
	"go.uber.org/atomic"
)

// converts connection to test variants
type Converter interface {
	Convert(udp.Connector) udp.Connector
}

type Server struct {
	udpServer  *udp.Server
	clients    map[string]*Connection
	closed     atomic.Bool
	wait       sync.WaitGroup
	newClients chan *Connection
	converter  Converter
}

func NewServer(addr *net.UDPAddr, converter Converter) (*Server, error) {
	udpServer, err := udp.NewServer(addr)
	if err != nil {
		return nil, err
	}
	s := &Server{
		udpServer:  udpServer,
		clients:    map[string]*Connection{},
		newClients: make(chan *Connection),
		converter:  converter,
	}
	return s, nil
}

func (s *Server) Run() {
	s.wait.Add(1)
	go func() {
		defer s.wait.Done()
		for !s.closed.Load() {
			packet, err := s.udpServer.Connection()
			if err != nil {
				continue
			}
			if s.converter != nil {
				packet = s.converter.Convert(packet)
			}
			data, err := packet.Receive()
			if err != nil {
				continue
			}
			conn, ok := s.clients[packet.Addr().String()]
			if !ok {
				conn = newConnection(packet)
				s.clients[packet.Addr().String()] = conn
				s.newClients <- conn
				conn.run()
			}
			if !s.closed.Load() {
				conn.recv.packets <- data
			}
		}
	}()
}

func (s *Server) Connections() <-chan *Connection {
	return s.newClients
}

func (s *Server) Close() {
	if s.closed.Load() {
		return
	}
	close(s.newClients)
	s.udpServer.Close()
	s.closed.Store(true)
	s.wait.Wait()
	for _, conn := range s.clients {
		conn.Close()
	}
}

func (s *Server) IsClosed() bool {
	return s.closed.Load()
}
