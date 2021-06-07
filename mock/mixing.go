package mock

import (
	"math/rand"
	"sync"

	"github.com/leviska/tcp-over-udp/tcp"
	"github.com/leviska/tcp-over-udp/udp"
)

type MixingFailureConn struct {
	BaseConn
	buf [][]byte
	mu  sync.Mutex
}

func NewMixingFailure(size int) *MixingFailureConn {
	return &MixingFailureConn{buf: make([][]byte, size)}
}

func (c *MixingFailureConn) Clone() tcp.Converter {
	res := NewMixingFailure(len(c.buf))
	return res
}

func (c *MixingFailureConn) Convert(conn udp.Connector) udp.Connector {
	c.Conn = conn
	return c
}

func (c *MixingFailureConn) Receive() ([]byte, error) {
	return c.Conn.Receive()
}

func (c *MixingFailureConn) Send(d []byte) error {
	ind := rand.Intn(len(c.buf))
	c.mu.Lock()
	v := c.buf[ind]
	c.buf[ind] = d
	c.mu.Unlock()
	if v != nil {
		return c.Conn.Send(v)
	}
	return nil
}
