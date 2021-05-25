package udp

import "net"

type Connector interface {
	Receive() ([]byte, error)
	Send([]byte) error
	Addr() *net.UDPAddr
	Close()
}
