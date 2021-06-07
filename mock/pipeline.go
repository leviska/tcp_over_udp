package mock

import (
	"github.com/leviska/tcp-over-udp/tcp"
	"github.com/leviska/tcp-over-udp/udp"
)

type Cloner interface {
	Clone() tcp.Converter
}

type Pipeline struct {
	Convs []Cloner
}

func NewPipeline(convs ...Cloner) *Pipeline {
	return &Pipeline{Convs: convs}
}

func (c *Pipeline) Convert(conn udp.Connector) udp.Connector {
	for _, conv := range c.Convs {
		conn = conv.Clone().Convert(conn)
	}
	return conn
}
