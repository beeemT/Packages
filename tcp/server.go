package tcp

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

//Server is the representing struct of a universal tcpserver.
type Server struct {
	port           int
	defaultTimeout time.Duration

	//maxClients <= 0 means no restriction in client count.
	defaultMaxReadBuffer, maxClients, curClients int64
	sigchan                                      chan struct{}
}

//NewServer is the constructor for tcp.server.
//A defaultTimeout of 0 means no timeouts.
//Timeouts have to be implemented by the programmer in the handle method that is passed to tcpserver.Start.
//For limited reading use conn.LimitedRead in the handle method.
//A maxClients value of 0 or lower causes the server to accept all incoming connections.
func NewServer(port int, defaultTimeout time.Duration, defaultMaxReadBuffer, maxClients int64) *Server {
	sigchan := make(chan struct{})
	return &Server{port, defaultTimeout, defaultMaxReadBuffer, maxClients, 0, sigchan}
}

//CurClients is the getter for server.curClients.
func (server Server) CurClients() int64 {
	return server.curClients
}

//CurClientsP is the getter for *server.curClients.
//The resulting pointer is not safe for concurrent read or write.
//Modifying CurClients through this pointer might cause the server to malfunction.
func (server Server) CurClientsP() *int64 {
	return &server.curClients
}

//MaxClients is the getter for server.maxClients.
func (server Server) MaxClients() int64 {
	return server.maxClients
}

//SetMaxClients is the setter for server.maxClients.
func (server *Server) SetMaxClients(maxClients int64) {
	server.maxClients = maxClients
}

//Start boots the server. The server waits for calling s.Stop() for a graceful shut down.
//Start returns the waitGroup for the server so the caller can wait for the server to finish.
//This method is mainly for specific server implementations and should not be called to start specific servers but by a specific Start implementation.
func (server *Server) Start(handle func(*Conn, ...interface{}), a ...interface{}) *sync.WaitGroup {
	var serverWaitGroup, connWaitGroup sync.WaitGroup
	serverWaitGroup.Add(1)

	//Shutdown Routine.
	go func() {
		defer serverWaitGroup.Done()
		<-server.sigchan
		cleanup(&connWaitGroup)
	}()

	go server.listenAndServe(&serverWaitGroup, &connWaitGroup, handle, a)
	return &serverWaitGroup
}

//Stop triggers the shut down of the server by closing the signal channel and triggering the cleanup.
func (server *Server) Stop() {
	close(server.sigchan)
}

//listenAndServe boots the server. Is designed to be called into a go routine.
//serverWaitGroup is returned upon start of the server so the caller can wait for the shutdown of the server.
//connWaitGroup manages all instances of handle and thus all clients.
func (server *Server) listenAndServe(serverWaitGroup, connWaitGroup *sync.WaitGroup, handle func(*Conn, ...interface{}), a ...interface{}) error {
	fmt.Println("Starting service ...")

	serverSocket, err := net.Listen("tcp", fmt.Sprintf(":%s", strconv.Itoa(server.port)))
	if err != nil {
		return errors.New("Failed at establishing serverSocket: " + err.Error())
	}
	defer serverSocket.Close()

	fmt.Println("Service started successfully!")

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
			conn := NewConn(netconn, server.defaultTimeout, server.defaultMaxReadBuffer)

			connWaitGroup.Add(1)
			atomic.AddInt64(&server.curClients, 1)

			go func() {
				defer connWaitGroup.Done()
				defer atomic.AddInt64(&server.curClients, -1)
				//TODO: Remove view for handle function on server.curClients ?
				handle(conn, a)
			}()

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
