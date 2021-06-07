package mock

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/leviska/tcp-over-udp/tcp"
	"github.com/leviska/tcp-over-udp/udp"
)

// allows to rewrite only send and receive
type BaseConn struct {
	Conn udp.Connector
}

func (c *BaseConn) Addr() *net.UDPAddr {
	return c.Conn.Addr()
}

func (c *BaseConn) IsClient() bool {
	return c.Conn.IsClient()
}

func (c *BaseConn) Close() {
	c.Conn.Close()
}

func (c *BaseConn) cKind() string {
	if c.IsClient() {
		return "client"
	} else {
		return "server"
	}
}

func (c *BaseConn) pKind(d []byte) string {
	kind := binary.BigEndian.Uint64(d[8:16])
	if kind == tcp.DATA {
		return "DATA"
	} else if kind == tcp.ACK {
		return "ACK"
	} else {
		return "UNDEFINED"
	}
}

func (c *BaseConn) genError(d []byte) error {
	return fmt.Errorf("packet (%s: %q) lost in the net", c.pKind(d), string(d[16:]))
}
