package mock

import (
	"github.com/leviska/tcp-over-udp/tcp"
	"github.com/leviska/tcp-over-udp/udp"
)

type StaticFailureConn struct {
	BaseConn
	badPrd  int
	goodPrd int
	cur     int
}

func NewStaticFailure(badPrd int, goodPrd int) *StaticFailureConn {
	return &StaticFailureConn{badPrd: badPrd, goodPrd: goodPrd}
}

func (c *StaticFailureConn) Clone() tcp.Converter {
	res := *c
	return &res
}

func (c *StaticFailureConn) Convert(conn udp.Connector) udp.Connector {
	c.Conn = conn
	return c
}

func (c *StaticFailureConn) Receive() ([]byte, error) {
	return c.Conn.Receive()
}

func (c *StaticFailureConn) Send(d []byte) error {
	defer func() {
		c.cur++
		if c.cur >= c.badPrd+c.goodPrd {
			c.cur = 0
		}
	}()

	if c.cur < c.badPrd {
		return c.genError(d)
	}

	return c.Conn.Send(d)
}
