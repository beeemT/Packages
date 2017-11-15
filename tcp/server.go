package tcp

/*
Logic:
2 waitgroups : serverWaitGroup and connWaitGroup. serverWaitGroup is returned upon start of the server so the caller
can wait on termination of the server.
serverWaitGroup.Done() is called in an extra go routine that manages the shutdown and is triggered upon
closing the sigchan of the server.
that routine then waits for the connWwaitgroup to finish before exiting itself.

the second waitgroup simply contains all current instances of handle.
*/

import (
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

//Server is the representing struct of a tcpserver
type Server struct {
	port                                         int
	defaultTimeout                               time.Duration
	defaultMaxReadBuffer, maxClients, curClients int64
	sigchan                                      chan struct{}
}

//NewServer is the constructor for tcp.server
func NewServer(port int, defaultTimeout time.Duration, defaultMaxReadBuffer, maxClients int64, sigchan chan struct{}) *Server {
	return &Server{port, defaultTimeout, defaultMaxReadBuffer, maxClients, 0, sigchan}
}

//CurClients is the getter for server.curCLients
func (server Server) CurClients() int64 {
	return server.curClients
}

//SetCurClients is the setter for server.curClients
func (server *Server) SetCurClients(curClients int64) {
	server.curClients = curClients
}

//Start boots the server. the server waits for closing its sigchan chan struct{} for shutting down.
//Start returns the waitGroup for the server so the caller can wait for the server to finish
//This method is mainly for specific server implementations and should not be called to start specific servers but be called by a specific start implementation.
func (server *Server) Start(handle func(*Conn, *sync.WaitGroup, *int64, ...interface{}), a ...interface{}) *sync.WaitGroup {
	var serverWaitGroup sync.WaitGroup
	serverWaitGroup.Add(1)
	go server.listenAndServe(&serverWaitGroup, handle, a)
	return &serverWaitGroup
}

//Stop triggers the shut down of the server by closing the signal channel this triggering the cleanup
func (server *Server) Stop() {
	close(server.sigchan)
}

//listenAndServe boots the server. Is designed to be called to go routine
func (server *Server) listenAndServe(serverWaitGroup *sync.WaitGroup, handle func(*Conn, *sync.WaitGroup, *int64, ...interface{}), a ...interface{}) error {
	fmt.Println("Starting service ...")

	serverSocket, err := net.Listen("tcp", strconv.Itoa(server.port))
	if err != nil {
		return errors.New("Failed at establishing serverSocket: " + err.Error())
	}
	defer serverSocket.Close()

	fmt.Println("Service started successfully!")

	var connWaitGroup sync.WaitGroup
	go func() {
		defer serverWaitGroup.Done()
		<-server.sigchan
		cleanup(&connWaitGroup)
	}()

	for {
		if server.maxClients > 0 && server.curClients >= server.maxClients {
			continue
		}
		select {
		default:
			netconn, err := serverSocket.Accept()
			if err != nil {
				log.Println("Failed at accepting new net.Conn: " + err.Error())
			}
			conn := newConn(netconn, server.defaultTimeout, server.defaultMaxReadBuffer)
			connWaitGroup.Add(1)
			atomic.AddInt64(&server.curClients, 1)
			go handle(conn, &connWaitGroup, &server.curClients, a)
		case <-server.sigchan:
			break
		}
	}
}

func cleanup(connWaitGroup *sync.WaitGroup) {
	fmt.Println("Caught shutdown signal. Shutting down server ...")
	connWaitGroup.Wait()
	fmt.Println("Successfully shut down server!")
	runtime.Goexit()
}
