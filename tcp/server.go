package tcp

import (
	"net"
)

type Server struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

func getAddr(ip string, port int) *net.UDPAddr {
	return &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
}

func NewServer(ip string, port int) (*Server, error) {
	addr := getAddr(ip, port)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		conn: conn,
		addr: addr,
	}, nil
}

func (s *Server) FetchClient() (*Connection, error) {
	buf := make([]byte, 2048)
	_, remoteaddr, err := s.conn.ReadFromUDP(buf)
	if err != nil {
		return nil, nil, err
	}
	return buf, &Connection{conn: s.conn, addr: remoteaddr}, nil
}
