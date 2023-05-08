package sc

import (
	"log"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/beeemT/Packages/netutil"
)

//Client is the implementation of a tcp Client that is meant to connect to the tcp server.
type Client struct {
	remoteAddr           net.IP
	remotePort           int
	defaultTimeout       time.Duration
	defaultMaxReadBuffer int64
	proto                protocol
}

//NewClient is the constructor for a networking client
//A defaultTimeout of 0 means that the connection does not time out.
//If the connection uses a limited read or not has to be decided in the passed handle method.
func NewClient(remoteAddr net.IP, remotePort int, defaultTimeout time.Duration, defaultMaxReadBuffer int64, proto protocol) *Client {
	return &Client{remoteAddr: remoteAddr,
		remotePort:           remotePort,
		defaultTimeout:       defaultTimeout,
		defaultMaxReadBuffer: defaultMaxReadBuffer,
		proto:                proto}
}

//NewTCPClient is the constructor for a networking client with the protocol field prefilled.
func NewTCPClient(remoteAddr net.IP, remotePort int, defaultTimeout time.Duration, defaultMaxReadBuffer int64) *Client {
	return &Client{remoteAddr: remoteAddr,
		remotePort:           remotePort,
		defaultTimeout:       defaultTimeout,
		defaultMaxReadBuffer: defaultMaxReadBuffer,
		proto:                tcp}
}

//NewUDPClient is the constructor for a networking client with the protocol field prefilled.
func NewUDPClient(remoteAddr net.IP, remotePort int, defaultTimeout time.Duration, defaultMaxReadBuffer int64) *Client {
	return &Client{remoteAddr: remoteAddr,
		remotePort:           remotePort,
		defaultTimeout:       defaultTimeout,
		defaultMaxReadBuffer: defaultMaxReadBuffer,
		proto:                udp}
}

//Connect is the exported api for the connect method. Is run in its' own routine.
//After the spawned routine ends, that is when the passed handle func returns, waitgroup.Done is called on the returned waitgroup.
//For using the built in timeout, look at net.Conn.SetDeadline .
//The opened connection is not automatically closed. This has to be part of the passed handle function.
func (client *Client) Connect(handle func(*Conn, ...interface{}), a ...interface{}) *sync.WaitGroup {
	var clientWaitGroup sync.WaitGroup
	clientWaitGroup.Add(1)
	go client.connect(&clientWaitGroup, handle, a...)
	return &clientWaitGroup
}

func (client *Client) connect(clientWaitGroup *sync.WaitGroup, handle func(*Conn, ...interface{}), a ...interface{}) {
	defer clientWaitGroup.Done()

	addr := netutil.BuildIPAddressString(client.remoteAddr, client.remotePort)
	netConn, err := net.Dial(client.proto.String(), addr)
	if err != nil {
		log.Printf("Establishing a conn with [%s over %s] failed: %s", addr, client.proto, err)
		runtime.Goexit()
	}

	conn := NewConn(netConn, client.defaultTimeout, client.defaultMaxReadBuffer)

	handle(conn, a...)
}
