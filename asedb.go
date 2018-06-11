package asedb

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

// DriverName is the driver name to use with sql.Open for ase databases.
const DriverName = "asedb"

var ase = &aseDriver{}

func init() {
	sql.Register(DriverName, ase)
}

//  check if driver implements all required interfaces
var (
	_ driver.Driver = (*aseDriver)(nil)
)

type aseDriver struct {
}

func (ase *aseDriver) Open(dsn string) (driver.Conn, error) {
	conn, err := NewDSNConnector(dsn)
	if err != nil {
		return nil, err
	}
	return conn.Connect(context.Background())
}

//  check if conn implements all required interfaces
var (
	_ driver.Conn           = (*conn)(nil)
	_ driver.QueryerContext = (*conn)(nil)
)

// database connection
type conn struct {
}

func newConn(ctx context.Context, c *Connector) (driver.Conn, error) {
	return &conn{}, nil
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	panic("deprecated")
}

func (c *conn) Close() error {
	//c.session.Close()
	return nil
}

func (c *conn) Begin() (driver.Tx, error) {
	panic("deprecated")
}

// QueryContext implements the database/sql/driver/QueryerContext interface.
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	//todo: implement me
	return nil, nil
}
