package mock

import (
	"encoding/binary"
	"errors"
	"os"

	"github.com/leviska/tcp-over-udp/tcp"
	"github.com/leviska/tcp-over-udp/udp"
	"github.com/leviska/tcp-over-udp/util"
)

type LoggingConn struct {
	BaseConn
}

func NewLogging() *LoggingConn {
	return &LoggingConn{}
}

func (c *LoggingConn) Clone() tcp.Converter {
	res := *c
	return &res
}

func (c *LoggingConn) Convert(conn udp.Connector) udp.Connector {
	c.Conn = conn
	return c
}

func (c *LoggingConn) Receive() ([]byte, error) {
	d, err := c.Conn.Receive()
	if err != nil {
		if !errors.Is(err, os.ErrDeadlineExceeded) {
			util.Logger.Infof("%s logging receive error: %v", c.cKind(), err)
		}
	} else {
		seqNum := binary.BigEndian.Uint64(d[0:8])
		util.Logger.Infof("%s logging receive %s [%d]: %q", c.cKind(), c.pKind(d), seqNum, string(d[16:]))
	}
	return d, err
}

func (c *LoggingConn) Send(d []byte) error {
	err := c.Conn.Send(d)
	if err != nil {
		util.Logger.Infof("%s logging send error: %v", c.cKind(), err)
	} else {
		seqNum := binary.BigEndian.Uint64(d[0:8])
		util.Logger.Infof("%s logging send %s [%d]: %q", c.cKind(), c.pKind(d), seqNum, string(d[16:]))
	}
	return err
}
