package sc

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

//Server is the representing struct of a universal server.
type Server struct {
	port           int
	defaultTimeout time.Duration

	//maxClients <= 0 means no restriction in client count.
	defaultMaxReadBuffer, maxClients, curClients int64
	sigchan                                      chan struct{}
	proto                                        protocol
}

//NewServer is the constructor for a server.
//A defaultTimeout of 0 means no timeouts.
//Timeouts have to be implemented by the programmer in the handle method that is passed to server.Start.
//For limited reading use conn.LimitedRead in the handle method.
//A maxClients value of 0 or lower causes the server to accept all incoming connections.
func NewServer(port int, defaultTimeout time.Duration, defaultMaxReadBuffer, maxClients int64, proto protocol) *Server {
	sigchan := make(chan struct{})
	return &Server{port: port,
		defaultTimeout:       defaultTimeout,
		defaultMaxReadBuffer: defaultMaxReadBuffer,
		maxClients:           maxClients,
		sigchan:              sigchan,
		proto:                proto}
}

//NewTCPServer is the constructor for a server with the protocol prefilled.
func NewTCPServer(port int, defaultTimeout time.Duration, defaultMaxReadBuffer, maxClients int64) *Server {
	sigchan := make(chan struct{})
	return &Server{port: port,
		defaultTimeout:       defaultTimeout,
		defaultMaxReadBuffer: defaultMaxReadBuffer,
		maxClients:           maxClients,
		sigchan:              sigchan,
		proto:                tcp}
}

//NewUDPServer is the constructor for a server with the protocol prefilled.
func NewUDPServer(port int, defaultTimeout time.Duration, defaultMaxReadBuffer, maxClients int64) *Server {
	sigchan := make(chan struct{})
	return &Server{port: port,
		defaultTimeout:       defaultTimeout,
		defaultMaxReadBuffer: defaultMaxReadBuffer,
		maxClients:           maxClients,
		sigchan:              sigchan,
		proto:                udp}
}

//CurClients is the getter for server.curClients.
func (server Server) CurClients() int64 {
	return server.curClients
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
//The handle function has to handle the close of the passed connection itself.
func (server *Server) Start(handle func(*Conn, ...interface{}), a ...interface{}) *sync.WaitGroup {
	var serverWaitGroup, connWaitGroup sync.WaitGroup
	serverWaitGroup.Add(1)

	//Shutdown Routine.
	go func() {
		defer serverWaitGroup.Done()
		<-server.sigchan
		cleanup(&connWaitGroup)
	}()

	serverWaitGroup.Add(1)
	go server.listenAndServe(&serverWaitGroup, &connWaitGroup, handle, a...)
	return &serverWaitGroup
}

//Stop triggers the shut down of the server by closing the signal channel and triggering the cleanup.
func (server *Server) Stop() {
	close(server.sigchan)
}

//Sigchan gets the receiving part of the servers signal channel.
func (server *Server) Sigchan() <-chan struct{} {
	return server.sigchan
}

//listenAndServe boots the server. Is designed to be called into a go routine.
//connWaitGroup manages all instances of handle and thus all clients.
func (server *Server) listenAndServe(serverWaitGroup, connWaitGroup *sync.WaitGroup, handle func(*Conn, ...interface{}), a ...interface{}) {
	log.Println("Starting service ...")
	defer serverWaitGroup.Done()

	closeFlag := false
	connChan := make(chan *net.Conn, 10)

	serverSocket, err := net.Listen(server.proto.String(), fmt.Sprintf(":%s", strconv.Itoa(server.port)))
	if err != nil {
		log.Printf("Failed at establishing serverSocket: %s\n", err.Error())
		return
	}

	defer func() {
		err = serverSocket.Close()
		if err != nil {
			if closeFlag {
				return
			}
			log.Println(err.Error())
		}
	}()

	serverWaitGroup.Add(1)
	go server.listen(&serverSocket, connChan, &closeFlag, serverWaitGroup)
	log.Println("Service started successfully!")

fl:
	for {
		select {
		case <-server.sigchan:
			err = serverSocket.Close()
			if err != nil {
				log.Println(err.Error())
			}
			closeFlag = true
			break fl

		case netConn := <-connChan:
			conn := NewConn(*netConn, server.defaultTimeout, server.defaultMaxReadBuffer)

			connWaitGroup.Add(1)
			atomic.AddInt64(&server.curClients, 1)

			go func() {
				defer connWaitGroup.Done()
				defer atomic.AddInt64(&server.curClients, -1)
				handle(conn, a...)
			}()
		}
	}
}

func (server *Server) listen(socket *net.Listener, connChan chan *net.Conn, closeFlag *bool, serverWaitGroup *sync.WaitGroup) {
	defer serverWaitGroup.Done()
	defer close(connChan)

fl:
	for {
		select {
		case <-server.sigchan:
			break fl
		default:
			if server.maxClients > 0 && server.curClients >= server.maxClients {
				continue
			}

			netConn, err := (*socket).Accept()
			if err != nil {
				if *closeFlag {
					continue
				}
				log.Println("Failed at accepting new connection: " + err.Error())
				continue
			}

			connChan <- &netConn
		}
	}
}

func cleanup(connWaitGroup *sync.WaitGroup) {
	log.Println("Caught shutdown signal. Shutting down server ...")
	connWaitGroup.Wait()
	log.Println("Successfully shut down server!")
}
