package sc

import (
	"io"
	"net"
	"time"
)

//Conn wraps net.Conn and implements timeouts and limited reading of conn.
type Conn struct {
	net.Conn
	timeout       time.Duration
	maxReadBuffer int64
}

//Timeout is the getter of type Conn.timeout
func (c Conn) Timeout() time.Duration {
	return c.timeout
}

//NewConn is the constructor for the Conn struct. The timeout of conn has to be set manually.
//The default 0 means no timeouts. The handle function for Conn has to use net.Conn.SetDeadline for timeouts.
func NewConn(conn net.Conn, timeout time.Duration, maxReadBuffer int64) *Conn {
	return &Conn{conn, timeout, maxReadBuffer}
}

//LimitedRead wraps the standard call to Read in a LimitReader.
//Returns the the amount of bytes read, which is the lower number of len(b) and maxReadBuffer.
func (c Conn) LimitedRead(b []byte) (int, error) {
	r := io.LimitReader(c.Conn, c.maxReadBuffer)
	return r.Read(b)
}
