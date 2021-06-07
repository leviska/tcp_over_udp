package udp

import (
	"net"
	"time"
)

const (
	IP   = "127.0.0.1"
	PORT = 1337
)

type Server struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

func NewAddr(ip string, port int) *net.UDPAddr {
	return &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
}

func NewServer(addr *net.UDPAddr) (*Server, error) {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		conn: conn,
		addr: addr,
	}, nil
}

func (s *Server) Connection() (Connector, error) {
	c := newConnection()
	s.conn.SetDeadline(time.Now().Add(time.Second * TimeoutSecs))
	n, remoteaddr, err := s.conn.ReadFromUDP(c.buf)
	if err != nil {
		return nil, err
	}
	c.conn = s.conn
	c.addr = remoteaddr
	c.buf = c.buf[:n]
	c.hasBuf = true
	return c, nil
}

func (s *Server) Close() {
	s.conn.Close()
}
