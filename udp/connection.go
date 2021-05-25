package udp

import "net"

type Connection struct {
	conn   *net.UDPConn
	addr   *net.UDPAddr
	buf    []byte
	client bool
	hasBuf bool
}

func emptyConnection() *Connection {
	return &Connection{
		buf: make([]byte, 2048),
	}
}

func NewConnection(addr *net.UDPAddr) (*Connection, error) {
	// create raw udp socket without connection
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	//conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	res := emptyConnection()
	res.conn = conn
	res.addr = addr
	res.client = true
	return res, nil
}

func (c *Connection) Close() {
	if !c.client {
		return
	}
	c.conn.Close()
}

func (c *Connection) Send(data []byte) error {
	_, err := c.conn.WriteToUDP(data, c.addr)
	return err
}

func (c *Connection) Receive() ([]byte, error) {
	if c.hasBuf {
		c.hasBuf = false
		return c.buf, nil
	}
	c.buf = c.buf[:cap(c.buf)]
	n, _, err := c.conn.ReadFromUDP(c.buf)
	c.buf = c.buf[:n]
	if err != nil {
		return nil, err
	}
	return c.buf, nil
}

func (c *Connection) Addr() *net.UDPAddr {
	return c.addr
}
