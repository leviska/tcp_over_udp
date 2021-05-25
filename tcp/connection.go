package tcp

import "net"

type Connection struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

func NewConnection(ip string, port int) (*Connection, error) {
	addr := getAddr(ip, port)
    conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	return &Connection{
		conn: conn,
		addr: addr,
	}, nil
}

func (c *Connection) Write()
