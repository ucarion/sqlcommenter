package sqlcommenter

import (
	"context"
	"database/sql/driver"
)

// Driver wraps a database/sql/driver.Driver so that all attributes set via
// NewContext are appended to queries.
//
// driver must be a driver.DriverContext. Any driver.Conn returned from driver
// must be a driver.ConnPrepareContext, driver.QueryerContext, and
// driver.ExecerContext. Otherwise, the returned driver will panic when used.
func Driver(driver driver.Driver) driver.Driver {
	return commentDriver{driver}
}

type commentDriver struct {
	driver.Driver
}

func (c commentDriver) OpenConnector(name string) (driver.Connector, error) {
	cc, err := c.Driver.(driver.DriverContext).OpenConnector(name)
	return connector{c, cc}, err
}

func (c commentDriver) Open(name string) (driver.Conn, error) {
	cc, err := c.Driver.Open(name)
	return conn{cc}, err
}

type connector struct {
	commentDriver
	driver.Connector
}

func (c connector) Connect(ctx context.Context) (driver.Conn, error) {
	cc, err := c.Connector.Connect(ctx)
	return conn{cc}, err
}

func (c connector) Driver() driver.Driver {
	return c.commentDriver
}

type conn struct {
	driver.Conn
}

func (c conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	query = appendComment(attrsFromContext(ctx), query)
	return c.Conn.(driver.ConnPrepareContext).PrepareContext(ctx, query)
}

func (c conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	query = appendComment(attrsFromContext(ctx), query)
	return c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
}

func (c conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	query = appendComment(attrsFromContext(ctx), query)
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}
