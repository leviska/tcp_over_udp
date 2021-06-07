package udp

import (
	"fmt"
	"net"
	"time"
)

type Connector interface {
	Receive() ([]byte, error)
	Send([]byte) error
	Addr() *net.UDPAddr
	IsClient() bool
	Close()
}

type Connection struct {
	conn   *net.UDPConn
	addr   *net.UDPAddr
	buf    []byte
	client bool
	hasBuf bool
}

const (
	MaxPacketSize = 8192
	TimeoutSecs   = 1
)

func newConnection() *Connection {
	return &Connection{
		buf: make([]byte, MaxPacketSize),
	}
}

func NewConnection(addr *net.UDPAddr) (*Connection, error) {
	// create raw udp socket without connection
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	//conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	res := newConnection()
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
	if !c.client {
		return nil, fmt.Errorf("can't receive from server connection")
	}
	c.buf = c.buf[:cap(c.buf)]
	c.conn.SetDeadline(time.Now().Add(time.Second * TimeoutSecs))
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

func (c *Connection) IsClient() bool {
	return c.client
}
