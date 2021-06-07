package tcp_test

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/leviska/tcp-over-udp/mock"
	"github.com/leviska/tcp-over-udp/tcp"
	"github.com/leviska/tcp-over-udp/udp"
	"github.com/leviska/tcp-over-udp/util"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	util.SetupLogger()
	tcp.RetryTimeout = time.Millisecond * 10 // for faster tests
	os.Exit(m.Run())
}

func runEcho(t *testing.T, conv tcp.Converter) (*tcp.Server, chan error) {
	addr := udp.NewAddr(udp.IP, udp.PORT)
	server, err := tcp.NewServer(addr, conv)
	assert.NoError(t, err)
	server.Run()
	clients := server.Connections()
	errChan := make(chan error)
	go func() {
		wait := sync.WaitGroup{}
		for c := range clients {
			wait.Add(1)
			go func(c *tcp.Connection) {
				defer wait.Done()
				for !c.IsClosed() {
					d, err := c.Receive()
					if err != nil {
						if !server.IsClosed() {
							errChan <- err
						}
						break
					}
					err = c.Send(d)
					if err != nil {
						if !server.IsClosed() {
							errChan <- err
						}
					}
				}
			}(c)
		}
		wait.Wait()
		close(errChan)
	}()
	time.Sleep(time.Millisecond * 100)
	return server, errChan
}

func noErrors(t *testing.T, errs chan error) {
	_, ok := <-errs
	assert.False(t, ok)
}

func runClient(t *testing.T, conv tcp.Converter) *tcp.Connection {
	conn, err := tcp.NewConnection(udp.NewAddr(udp.IP, udp.PORT), conv)
	assert.NoError(t, err)
	return conn
}

func sendEcho(t *testing.T, client *tcp.Connection, message string) {
	err := client.Send([]byte(message))
	assert.NoError(t, err)
	res, err := client.Receive()
	assert.NoError(t, err)
	if !assert.Equal(t, message, string(res)) {
		panic("aa")
	}
}

func TestEcho(t *testing.T) {
	server, errs := runEcho(t, nil)

	client := runClient(t, nil)

	sendEcho(t, client, "hello")
	sendEcho(t, client, "")
	sendEcho(t, client, strings.Repeat("a", tcp.MaxPacketSize))
	sendEcho(t, client, `{"glossary":{"title":"example glossary","GlossDiv":{"title":"S","GlossList":{"GlossEntry":{"ID":"SGML","SortAs":"SGML","GlossTerm":"Standard Generalized Markup Language","Acronym":"SGML","Abbrev":"ISO 8879:1986","GlossDef":{"para":"A meta-markup language, used to create markup languages such as DocBook.","GlossSeeAlso":["GML","XML"]},"GlossSee":"markup"}}}}}`)

	client.Close()
	server.Close()
	noErrors(t, errs)
}

func TestMany(t *testing.T) {
	server, errs := runEcho(t, nil)
	client := runClient(t, nil)

	defer func() {
		client.Close()
		server.Close()
		noErrors(t, errs)
	}()

	for i := 0; i < 10000; i++ {
		sendEcho(t, client, util.RandomString(1024))
	}
}

func TestEveryOther(t *testing.T) {
	pipeline := mock.NewPipeline(mock.NewStaticFailure(1, 1))
	//pipeline.Convs = append(pipeline.Convs, mock.NewLogging())
	server, errs := runEcho(t, pipeline)
	client := runClient(t, pipeline)

	defer func() {
		client.Close()
		server.Close()
		noErrors(t, errs)
	}()

	for i := 0; i < 30; i++ {
		sendEcho(t, client, util.RandomString(10))
	}
}

func TestRandom(t *testing.T) {
	pipeline := mock.NewPipeline(mock.NewRandomFailure(70.0))
	//pipeline.Convs = append(pipeline.Convs, mock.NewLogging())
	server, errs := runEcho(t, pipeline)
	client := runClient(t, pipeline)

	defer func() {
		client.Close()
		server.Close()
		noErrors(t, errs)
	}()

	for i := 0; i < 30; i++ {
		sendEcho(t, client, util.RandomString(10))
	}
}

func TestRandomEveryOther(t *testing.T) {
	pipeline := mock.NewPipeline(mock.NewRandomFailure(50.0))
	//pipeline.Convs = append(pipeline.Convs, mock.NewLogging())
	server, errs := runEcho(t, pipeline)
	client := runClient(t, pipeline)

	defer func() {
		client.Close()
		server.Close()
		noErrors(t, errs)
	}()

	for i := 0; i < 30; i++ {
		sendEcho(t, client, util.RandomString(10))
	}
}

func TestMixing(t *testing.T) {
	pipeline := mock.NewPipeline()
	//pipeline.Convs = append(pipeline.Convs, mock.NewLogging())
	pipeline.Convs = append(pipeline.Convs, mock.NewMixingFailure(10))
	server, errs := runEcho(t, pipeline)
	client := runClient(t, pipeline)

	defer func() {
		client.Close()
		server.Close()
		noErrors(t, errs)
	}()

	for i := 0; i < 10; i++ {
		sendEcho(t, client, fmt.Sprint(i))
	}
}
