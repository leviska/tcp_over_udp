package tcp

import (
	"fmt"
	"net"
	"sync"

	"github.com/leviska/tcp-over-udp/udp"
)

const ( // this is an enum in go
	DATA = iota
	ACK
)

const (
	MaxPacketSize = 8192 - 16
)

type infoConn struct {
	udpConn    udp.Connector
	recvSeqNum uint64
	sendSeqNum uint64
	closed     bool
	mu         sync.Mutex
	wait       sync.WaitGroup
	notify     chan struct{}
}

type Connection struct {
	recv *recvConn
	send *sendConn
	info *infoConn
}

func newConnection(udpConn udp.Connector) *Connection {
	info := &infoConn{
		udpConn: udpConn,
		notify:  make(chan struct{}),
	}
	return &Connection{
		info: info,
		recv: &recvConn{
			info:     info,
			packets:  make(chan []byte, 2),
			recvChan: make(chan []byte, 2),
		},
		send: &sendConn{
			info:     info,
			sendChan: make(chan []byte),
			doneChan: make(chan error),
		},
	}
}

func NewConnection(addr *net.UDPAddr, converter Converter) (*Connection, error) {
	var udpConn udp.Connector
	var err error
	udpConn, err = udp.NewConnection(addr)
	if err != nil {
		return nil, err
	}
	if converter != nil {
		udpConn = converter.Convert(udpConn)
	}
	conn := newConnection(udpConn)
	conn.run()
	return conn, nil
}

func (c *Connection) run() {
	c.info.wait.Add(2)
	go c.recv.run()
	go c.send.run()
}

func (c *Connection) Receive() ([]byte, error) {
	return c.recv.receive()
}

func (c *Connection) Send(data []byte) error {
	if data == nil {
		return fmt.Errorf("can't sent nil data")
	}
	return c.send.send(data)
}

func (c *Connection) Close() {
	c.info.mu.Lock()
	if c.info.closed {
		c.info.mu.Unlock()
		return
	}
	c.info.closed = true
	c.info.mu.Unlock()
	close(c.send.sendChan)
	close(c.send.doneChan)
	close(c.recv.packets)
	close(c.recv.recvChan)
	close(c.info.notify)
	c.info.wait.Wait()
}

func (c *Connection) IsClosed() bool {
	return c.info.isClosed()
}

func (c *infoConn) isClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}
