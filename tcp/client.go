package tcp

import (
	"log"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
)

//Client is the implementation of a tcp Client that is meant to connect to the tcp client
type Client struct {
	port                 int
	remoteAddr           net.IP
	defaultTimeout       time.Duration
	defaultMaxReadBuffer int64
}

//NewClient is the constructor for client
func NewClient(port int, remoteAddr net.IP, defaultTimeout time.Duration, defaultMaxReadBuffer int64) *Client {
	return &Client{port, remoteAddr, defaultTimeout, defaultMaxReadBuffer}
}

//Connect is the exported api for the connect method. Is run in own routine.
func (client *Client) Connect(clientWaitGroup *sync.WaitGroup, handle func(*Conn, *sync.WaitGroup, ...interface{}), a ...interface{}) {
	go client.Connect(clientWaitGroup, handle, a)
}

func (client *Client) connect(clientWaitGroup *sync.WaitGroup, handle func(*Conn, []byte, ...interface{}), a ...interface{}) {
	defer clientWaitGroup.Done()

	addr := client.remoteAddr.String()
	if client.remoteAddr.To16() != nil {
		addr = "[" + addr + "]"
	}
	addr = addr + ":" + strconv.Itoa(client.port)
	netconn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Establishing a conn with [%s] failed: %s", addr, err)
		runtime.Goexit()
	}
	conn := newConn(netconn, client.defaultTimeout, client.defaultMaxReadBuffer)
	defer conn.Close()

	readBuffer := make([]byte, client.defaultMaxReadBuffer)
	handle(conn, readBuffer, a)
}
