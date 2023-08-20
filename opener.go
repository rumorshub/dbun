package dbun

import (
	"errors"
	"fmt"
	"sync"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/schema"
)

var _ Opener = (*BunOpener)(nil)

const (
	MySQL      Driver = "mysql"
	SQLiteShim Driver = "sqliteshim"
	Pgx        Driver = "pgx"
	Pg         Driver = "pg"
	AzureSQL   Driver = "azuresql"
	MsSQL      Driver = "mssql"
	SQLServer  Driver = "sqlserver"
)

var ErrDialectNotFound = errors.New("bun sql dialect not found")

type Driver string

type Opener interface {
	OpenBunDB(name string, opts ...bun.DBOption) (*bun.DB, error)
}

type BunOpener struct {
	inner SQLDBOpener
	dbs   map[string]*bun.DB
	mu    sync.Mutex
}

func NewOpener(inner SQLDBOpener) *BunOpener {
	return &BunOpener{
		inner: inner,
		dbs:   map[string]*bun.DB{},
	}
}

func (o *BunOpener) OpenBunDB(name string, opts ...bun.DBOption) (*bun.DB, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if db, ok := o.dbs[name]; ok {
		return db, nil
	}

	sqldb, driverName, err := o.inner.OpenDB(name)
	if err != nil {
		return nil, err
	}

	d, err := dialect(driverName)
	if err != nil {
		return nil, err
	}

	o.dbs[name] = bun.NewDB(sqldb, d, opts...)

	return o.dbs[name], nil
}

func dialect(name string) (schema.Dialect, error) {
	switch Driver(name) {
	case MySQL:
		return mysqldialect.New(), nil
	case SQLiteShim:
		return sqlitedialect.New(), nil
	case AzureSQL, MsSQL, SQLServer:
		return mssqldialect.New(), nil
	case Pgx, Pg:
		return pgdialect.New(), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrDialectNotFound, name)
	}
}
