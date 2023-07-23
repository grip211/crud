package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	builder "github.com/doug-martin/goqu/v9"

	// nolint:revive // it's OK
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/grip211/crud/pkg/database"
)

// тут реализуем структуру, инкапсулирующая внутри себя коннекты с базой данных

type ConnectionPool struct {
	db *builder.Database
}

func (c *ConnectionPool) Builder() *builder.Database {
	return c.db
}

func New(ctx context.Context, opt *database.Opt) (*ConnectionPool, error) {
	opt.UnwrapOrPanic()

	db, err := sql.Open(opt.Dialect, opt.ConnectionString())
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = db.PingContext(pingCtx); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(opt.MaxIdleConns)
	db.SetMaxOpenConns(opt.MaxOpenConns)
	db.SetConnMaxLifetime(opt.MaxConnMaxLifetime)

	dialect := builder.Dialect(opt.Dialect)
	connect := &ConnectionPool{
		db: dialect.DB(db),
	}

	if opt.Debug {
		logger := &database.Logger{}
		logger.SetCallback(func(format string, v ...interface{}) {
			fmt.Println(v)
		})
		connect.db.Logger(logger)
	}

	return connect, nil
}
