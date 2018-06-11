package internal

import (
	"bufio"
	"crypto/tls"
	"net"
	"sync"
	"time"
)

// SessionConn wraps the database tcp connection. It sets timeouts and handles driver ErrBadConn behavior.
type sessionConn struct {
	addr     string
	timeout  time.Duration
	conn     net.Conn
	isBad    bool  // bad connection
	badError error // error cause for session bad state
	inTx     bool  // in transaction
}

// session parameter
type sessionPrm interface {
	Host() string
	Username() string
	Password() string
	Locale() string
	FetchSize() int
	Timeout() int
	TLSConfig() *tls.Config
}

// Session represents a HDB session.
type Session struct {
	prm sessionPrm

	conn *sessionConn
	rd   *bufio.Reader
	wr   *bufio.Writer

	//serialize write request - read reply
	//supports calling session methods in go routines (driver methods with context cancellation)
	mu sync.Mutex
}
