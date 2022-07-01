package otelsqlcommenter_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ucarion/sqlcommenter"
	"github.com/ucarion/sqlcommenter/otelsqlcommenter"
	"go.opentelemetry.io/otel/trace"
)

func TestNewContext(t *testing.T) {
	ctx := context.Background()

	var state trace.TraceState
	state, _ = state.Insert("k", "v")

	ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1},
		SpanID:     trace.SpanID{1},
		TraceFlags: 0,
		TraceState: state,
		Remote:     false,
	}))

	ctx = otelsqlcommenter.NewContext(ctx)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}

	// sqlcommenter requires a driver.ContextDriver, but sqlmock only provides a
	// driver.Driver, so we wrap it with fakeDriverContext.
	sql.Register("sqlcommenter", sqlcommenter.Driver(fakeDriverContext{db.Driver()}))
	db, err = sql.Open("sqlcommenter", "sqlmock_db_0")
	if err != nil {
		t.Fatalf("open sqlmock: %v", err)
	}

	mock.ExpectQuery("select 1 /*traceparent='00-01000000000000000000000000000000-0100000000000000-00',tracestate='k=v'*/").WillReturnRows()
	if _, err := db.QueryContext(ctx, "select 1"); err != nil {
		t.Fatalf("db: %v", err)
	}
}

// fakeDriverContext makes a driver.DriverContext out of any driver.Driver.
type fakeDriverContext struct {
	driver.Driver
}

func (f fakeDriverContext) OpenConnector(name string) (driver.Connector, error) {
	conn, err := f.Driver.Open(name)
	if err != nil {
		return nil, err
	}

	return fakeConnector{f, conn}, nil
}

type fakeConnector struct {
	driver driver.Driver
	conn   driver.Conn
}

func (f fakeConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return f.conn, nil
}

func (f fakeConnector) Driver() driver.Driver {
	return f.driver
}
