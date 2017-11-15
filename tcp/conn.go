package tcp

import (
	"io"
	"net"
	"time"
)

//Conn wraps the net.Conn and implements timeouts and limiting the read bytes of conn
type Conn struct {
	net.Conn
	timeout       time.Duration
	maxReadBuffer int64
}

//Timeout is the getter of type Conn.timeout
func (c Conn) Timeout() time.Duration {
	return c.timeout
}

//NewConn is the constructor for a fileserver.Conn . The timeout of conn has to be set manually. Default is 0.
func newConn(conn net.Conn, timeout time.Duration, maxReadBuffer int64) *Conn {
	return &Conn{conn, timeout, maxReadBuffer}
}

func (c Conn) Read(b []byte) (int, error) {
	r := io.LimitReader(c.Conn, c.maxReadBuffer)
	return r.Read(b)
}
