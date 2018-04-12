package tcp

import (
	"log"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//Client is the implementation of a tcp Client that is meant to connect to the tcp server
type Client struct {
	remotePort           int
	remoteAddr           net.IP
	defaultTimeout       time.Duration
	defaultMaxReadBuffer int64
}

//NewClient is the constructor for client
func NewClient(remotePort int, remoteAddr net.IP, defaultTimeout time.Duration, defaultMaxReadBuffer int64) *Client {
	return &Client{remotePort, remoteAddr, defaultTimeout, defaultMaxReadBuffer}
}

//Connect is the exported api for the connect method. Is run in own routine.
//After the spawned routine ends, that is when the passed handle func returns waitgroup.Done is called.
//For using the built in timeout, look at net.Conn.SetDeadline .
func (client *Client) Connect(clientWaitGroup *sync.WaitGroup, handle func(*Conn, ...interface{}), a ...interface{}) {
	go client.connect(clientWaitGroup, handle, a)
}

func (client *Client) connect(clientWaitGroup *sync.WaitGroup, handle func(*Conn, ...interface{}), a ...interface{}) {
	defer clientWaitGroup.Done()

	addr := client.remoteAddr.String()
	if client.remoteAddr.To16() != nil {
		addr = "[" + addr + "]"
	}

	addr = addr + ":" + strconv.Itoa(client.remotePort)
	netConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Establishing a conn with [%s] failed: %s", addr, err)
		runtime.Goexit()
	}

	conn := newConn(netConn, client.defaultTimeout, client.defaultMaxReadBuffer)
	defer conn.Close()

	handle(conn, a)
}
