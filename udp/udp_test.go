package udp

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/leviska/tcp-over-udp/util"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	util.SetupLogger()
	go runServer(IP, PORT)
	time.Sleep(time.Millisecond * 100)
	os.Exit(m.Run())
}

func runServer(ip string, port int) {
	util.Logger.Info("[server] ", "running server")
	server, err := NewServer(NewAddr(ip, port))
	if err != nil {
		util.Logger.Panic("[server] ", err)
	}
	util.Logger.Infof("[server] listening on %s:%d", ip, port)

	for {
		conn, err := server.Connection()
		go func() {
			if err != nil {
				util.Logger.Error("[server] ", err)
				return
			}
			defer conn.Close()
			//util.Logger.Infof("[server] new connection from %q", conn.Addr())
			buf, err := conn.Receive()
			if err != nil {
				util.Logger.Error("[server] ", err)
				return
			}
			//util.Logger.Infof("[server] message from %q: %q", conn.Addr(), util.Truncate(string(buf), 40))
			err = conn.Send(buf)
			if err != nil {
				util.Logger.Error("[server] ", err)
				return
			}
		}()
	}
}

func sendEcho(assert *assert.Assertions, conn *Connection, message string) {
	util.Logger.Infof("[client] sending message %q", util.Truncate(message, 40))
	err := conn.Send([]byte(message))
	assert.NoError(err)
	util.Logger.Infof("[client] waiting for response")
	buf, err := conn.Receive()
	assert.NoError(err)
	received := string(buf)
	util.Logger.Infof("[client] got response: %q", util.Truncate(received, 40))
	assert.Equal(message, received)
}

func TestEcho(t *testing.T) {
	assert := assert.New(t)

	conn, err := NewConnection(NewAddr(IP, PORT))
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	sendEcho(assert, conn, "hello")
	sendEcho(assert, conn, "")
	sendEcho(assert, conn, strings.Repeat("a", MaxPacketSize))
	sendEcho(assert, conn, `{"glossary":{"title":"example glossary","GlossDiv":{"title":"S","GlossList":{"GlossEntry":{"ID":"SGML","SortAs":"SGML","GlossTerm":"Standard Generalized Markup Language","Acronym":"SGML","Abbrev":"ISO 8879:1986","GlossDef":{"para":"A meta-markup language, used to create markup languages such as DocBook.","GlossSeeAlso":["GML","XML"]},"GlossSee":"markup"}}}}}`)
}

func TestBandwidth(t *testing.T) {
	assert := assert.New(t)

	conn, err := NewConnection(NewAddr(IP, PORT))
	if !assert.NoError(err) {
		return
	}
	defer conn.Close()

	message := []byte(strings.Repeat("a", 1024))
	start := time.Now()
	cnt := 0
	for ; time.Since(start) < time.Second; cnt++ {
		err := conn.Send(message)
		assert.NoError(err)
		_, err = conn.Receive()
		assert.NoError(err)
	}
	fmt.Printf("Bandwidth is ~%d kb/s", cnt)
}
