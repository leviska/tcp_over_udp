package tcp

import (
	"encoding/binary"
	"fmt"
	"time"
)

const MaxTries = 10000

var RetryTimeout = time.Millisecond * 100

type sendConn struct {
	info     *infoConn
	sendChan chan []byte
	doneChan chan error
}

func newPacket(seqNum uint64, kind uint64, data []byte) []byte {
	res := make([]byte, 16, 16+len(data))
	binary.BigEndian.PutUint64(res[0:8], seqNum)
	binary.BigEndian.PutUint64(res[8:16], kind)
	res = append(res, data...)
	return res
}

func (c *sendConn) send(data []byte) error {
	c.sendChan <- data
	err, ok := <-c.doneChan
	if err == nil && !ok {
		return fmt.Errorf("can't send: %v", err)
	}
	return err
}

func (c *sendConn) timeout(seqNum uint64) {
	c.info.mu.Lock()
	defer c.info.mu.Unlock()
	curNum := c.info.sendSeqNum
	if c.info.closed || curNum > seqNum {
		return
	}
	c.info.notify <- struct{}{}
}

func (c *sendConn) run() {
	defer c.info.wait.Done()
	for !c.info.isClosed() {
		var data []byte
		var ok bool
		select {
		case data, ok = <-c.sendChan:
			if !ok {
				continue
			}
		case <-c.info.notify:
			continue
		}
		c.info.mu.Lock()
		seqNum := c.info.sendSeqNum
		c.info.mu.Unlock()
		kind := uint64(DATA)
		packet := newPacket(seqNum, kind, data)
		tries := 0
		var lastErr error
		for ; !c.info.isClosed() && tries < MaxTries; tries++ {
			if err := c.info.udpConn.Send(packet); err != nil {
				lastErr = err
				continue
			}
			go func() {
				time.Sleep(RetryTimeout)
				c.timeout(seqNum)
			}()
			<-c.info.notify
			c.info.mu.Lock()
			curNum := c.info.sendSeqNum
			c.info.mu.Unlock()
			if curNum > seqNum {
				break
			}
		}
		var err error
		if tries >= MaxTries {
			if lastErr == nil {
				err = fmt.Errorf("too many retries")
			} else {
				err = fmt.Errorf("too many retries, last error: %v", lastErr)
			}
		}
		c.info.mu.Lock()
		if !c.info.closed {
			c.doneChan <- err
		}
		c.info.mu.Unlock()
	}
}
