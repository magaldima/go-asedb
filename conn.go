package asedb

import (
	"context"
	"crypto/tls"
	"database/sql/driver"
	"net/url"
	"sync"
	"time"
)

const (
	minFetchSize = 1000
	minTimeout   = 60 * time.Second
)

// A Connector represents an ase driver in a fixed configuration
type Connector struct {
	mu                       sync.RWMutex
	host, username, password string
	locale                   string
	bufferSize, fetchSize    int
	timeout                  time.Duration
	tlsConfig                *tls.Config
}

func newConnector() *Connector {
	return &Connector{
		fetchSize: 0,
		timeout:   10 * time.Second,
	}
}

// NewBasicAuthConnector creates a connector for basic authentication.
func NewBasicAuthConnector(host, username, password string) *Connector {
	c := newConnector()
	c.host = host
	c.username = username
	c.password = password
	return c
}

// NewDSNConnector creates a connector from a data source name.
func NewDSNConnector(dsn string) (*Connector, error) {
	c := newConnector()

	url, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	c.host = url.Host
	if url.User != nil {
		c.username = url.User.Username()
		c.password, _ = url.User.Password()
	}

	//todo: add url query logic

	return c, nil
}

// Host returns the host of the connector.
func (c *Connector) Host() string {
	return c.host
}

// Username returns the username of the connector.
func (c *Connector) Username() string {
	return c.username
}

// Password returns the password of the connector.
func (c *Connector) Password() string {
	return c.password
}

// Locale returns the locale of the connector.
func (c *Connector) Locale() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.locale
}

/*
SetLocale sets the locale of the connector.
For more information please see DSNLocale.
*/
func (c *Connector) SetLocale(locale string) {
	c.mu.Lock()
	c.locale = locale
	c.mu.Unlock()
}

// FetchSize returns the fetchSize of the connector.
func (c *Connector) FetchSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.fetchSize
}

/*
SetFetchSize sets the fetchSize of the connector.
For more information please see DSNFetchSize.
*/
func (c *Connector) SetFetchSize(fetchSize int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if fetchSize < minFetchSize {
		fetchSize = minFetchSize
	}
	c.fetchSize = fetchSize
	return nil
}

// Timeout returns the timeout of the connector.
func (c *Connector) Timeout() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.timeout
}

/*
SetTimeout sets the timeout of the connector.
For more information please see DSNTimeout.
*/
func (c *Connector) SetTimeout(timeout time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if timeout < minTimeout {
		timeout = minTimeout
	}
	c.timeout = timeout
	return nil
}

// TLSConfig returns the TLS configuration of the connector.
func (c *Connector) TLSConfig() *tls.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tlsConfig
}

// SetTLSConfig sets the TLS configuration of the connector.
func (c *Connector) SetTLSConfig(tlsConfig *tls.Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tlsConfig = tlsConfig
	return nil
}

// Connect implements the database/sql/driver/Connector interface.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	return newConn(ctx, c)
}

// Driver implements the database/sql/driver/Connector interface.
func (c *Connector) Driver() driver.Driver {
	return ase
}
