package tcp

import (
	"encoding/binary"
	"fmt"
)

type recvConn struct {
	info     *infoConn
	packets  chan []byte
	recvChan chan []byte
}

func (c *recvConn) receive() ([]byte, error) {
	data, ok := <-c.recvChan
	if data == nil || !ok {
		return nil, fmt.Errorf("can't reveive: connection closed")
	}
	return data, nil
}

func (c *recvConn) run() {
	defer c.info.wait.Done()
	for !c.info.isClosed() {
		var data []byte
		var err error
		if c.info.udpConn.IsClient() {
			data, err = c.info.udpConn.Receive()
		} else {
			var ok bool
			data, ok = <-c.packets
			if !ok {
				err = fmt.Errorf("channel closed")
			}
		}
		if err != nil {
			continue
		}
		seqNum := binary.BigEndian.Uint64(data[0:8])
		kind := binary.BigEndian.Uint64(data[8:16])
		switch kind {
		case DATA:
			c.info.mu.Lock()
			if seqNum == c.info.recvSeqNum {
				c.info.recvSeqNum++
				res := make([]byte, len(data)-16)
				copy(res, data[16:])
				c.recvChan <- res
			}
			c.info.mu.Unlock()
			c.info.udpConn.Send(newPacket(c.info.recvSeqNum, ACK, []byte{}))
		case ACK:
			c.info.mu.Lock()
			if !c.info.closed && seqNum > c.info.sendSeqNum {
				c.info.sendSeqNum = seqNum
				c.info.notify <- struct{}{}
			}
			c.info.mu.Unlock()
		}
	}
}
