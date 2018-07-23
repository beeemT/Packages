package tcp

import (
	"log"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/grekhor/Packages/netutil"
)

//Client is the implementation of a tcp Client that is meant to connect to the tcp server.
type Client struct {
	remoteAddr           net.IP
	remotePort           int
	defaultTimeout       time.Duration
	defaultMaxReadBuffer int64
}

//NewClient is the constructor for a tcp client.
//A defaultTimeout of 0 means that the connection does not time out.
//If the connection uses a limited read or not has to be decided in the passed handle method.
func NewClient(remoteAddr net.IP, remotePort int, defaultTimeout time.Duration, defaultMaxReadBuffer int64) *Client {
	return &Client{remoteAddr, remotePort, defaultTimeout, defaultMaxReadBuffer}
}

//Connect is the exported api for the connect method. Is run in its' own routine.
//After the spawned routine ends, that is when the passed handle func returns, waitgroup.Done is called on the returned waitgroup.
//For using the built in timeout, look at net.Conn.SetDeadline .
func (client *Client) Connect(handle func(*Conn, ...interface{}), a ...interface{}) *sync.WaitGroup {
	var clientWaitGroup sync.WaitGroup
	clientWaitGroup.Add(1)
	go client.connect(&clientWaitGroup, handle, a...)
	return &clientWaitGroup
}

func (client *Client) connect(clientWaitGroup *sync.WaitGroup, handle func(*Conn, ...interface{}), a ...interface{}) {
	defer clientWaitGroup.Done()

	addr := client.remoteAddr.String()
	if netutil.IsIPv6(client.remoteAddr) {
		addr = "[" + addr + "]"
	}

	addr = addr + ":" + strconv.Itoa(client.remotePort)
	netConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Establishing a conn with [%s] failed: %s", addr, err)
		runtime.Goexit()
	}

	conn := NewConn(netConn, client.defaultTimeout, client.defaultMaxReadBuffer)
	defer conn.Close()

	handle(conn, a...)
}
