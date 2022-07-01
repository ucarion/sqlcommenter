package sqlcommenter_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ucarion/sqlcommenter"
)

func TestDriver(t *testing.T) {
	ctx := sqlcommenter.NewContext(context.Background(), "foo", "bar")

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

	mock.ExpectQuery("select 1 /*foo='bar'*/").WillReturnRows()
	if _, err := db.QueryContext(ctx, "select 1"); err != nil {
		t.Fatalf("db: %v", err)
	}

	mock.ExpectExec("select 2 /*foo='bar'*/").WillReturnResult(driver.ResultNoRows)
	if _, err := db.ExecContext(ctx, "select 2"); err != nil {
		t.Fatalf("db: %v", err)
	}

	mock.ExpectPrepare("select 3 /*foo='bar'*/")
	if _, err := db.PrepareContext(ctx, "select 3"); err != nil {
		t.Fatalf("db: %v", err)
	}

	// this is notionally redundant, but let's confirm that things work even
	// when using txs
	mock.ExpectBegin()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("db: %v", err)
	}

	mock.ExpectQuery("select 1 /*foo='bar'*/").WillReturnRows()
	if _, err := tx.QueryContext(ctx, "select 1"); err != nil {
		t.Fatalf("db: %v", err)
	}

	mock.ExpectExec("select 2 /*foo='bar'*/").WillReturnResult(driver.ResultNoRows)
	if _, err := tx.ExecContext(ctx, "select 2"); err != nil {
		t.Fatalf("db: %v", err)
	}

	mock.ExpectPrepare("select 3 /*foo='bar'*/")
	if _, err := tx.PrepareContext(ctx, "select 3"); err != nil {
		t.Fatalf("db: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock: %v", err)
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
