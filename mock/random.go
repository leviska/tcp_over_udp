package mock

import (
	"math/rand"

	"github.com/leviska/tcp-over-udp/tcp"
	"github.com/leviska/tcp-over-udp/udp"
)

type RandomFailureConn struct {
	BaseConn
	prob float32
}

func NewRandomFailure(prob float32) *RandomFailureConn {
	return &RandomFailureConn{prob: prob}
}

func (c *RandomFailureConn) Clone() tcp.Converter {
	res := *c
	return &res
}

func (c *RandomFailureConn) Convert(conn udp.Connector) udp.Connector {
	c.Conn = conn
	return c
}

func (c *RandomFailureConn) Unlucky() bool {
	return rand.Float32()*100 > c.prob
}

func (c *RandomFailureConn) Receive() ([]byte, error) {
	return c.Conn.Receive()
}

func (c *RandomFailureConn) Send(d []byte) error {
	if c.Unlucky() {
		return c.genError(d)
	}
	err := c.Conn.Send(d)
	return err
}
