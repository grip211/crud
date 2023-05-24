## CRUD - simple project

### TODO:

1. **Поменять структуру проекта на следующие**
    - cmd/crud/main.go (server.go перемещается в этот файл)
    - pkg/repository/repo.go (сюда надо перенести работу с базой)
    - templates/*.html (остается также

2. **Добавить в проект следующие зависимости:**
    - `go get -u github.com/urfave/cli/v2`
    - `go get -u github.com/doug-martin/goqu/v9`
    - `go get -u github.com/gofiber/fiber/v2`

3. **Интегрироваться с github.com/doug-martin/goqu/v9:**

> Или с аналгами типа viper, но мне нравится больше эта поэтому будем работать с ней :)

1. создаем файл в pkg/database/mysql/pool.go
2. пишем код

```go
package mysql

// тут реализуем структуру, инкапсулирующая внутри себя коннекты с базой данных

import (
	"context"
	"database/sql"
	"time"

	builder "github.com/doug-martin/goqu/v9"

	// nolint:golint // it's OK
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/grip211/crud/pkg/database"
)

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
			log.Info(v)
		})
		connect.db.Logger(logger)
	}

	return connect, nil
}
```

3. создаем файл pkg/database/option.go

```go
// тут мы пишем простую обертку над параметрами подключения и их валидацией

package database

import (
	"fmt"
	"strings"
	"time"
)

type Opt struct {
	Host               string        `yaml:"host"`
	User               string        `yaml:"user"`
	Password           string        `yaml:"password"`
	Port               string        `yaml:"port"`
	Name               string        `yaml:"name"`
	Dialect            string        `yaml:"dialect"`
	Debug              bool          `yaml:"debug"`
	MaxIdleConns       int           `yaml:"max_idle_conns"`
	MaxOpenConns       int           `yaml:"max_open_conns"`
	MaxConnMaxLifetime time.Duration `yaml:"max_conn_max_lifetime"`
}

func (o *Opt) UnwrapOrPanic() {
	if o.Dialect == "" {
		o.Dialect = "mysql"
	}
	if o.Host == "" {
		o.Host = "@"
	}
	if strings.EqualFold(o.Host, "") {
		o.Host = "@"
	} else if !strings.Contains(o.Host, "@") {
		o.Host = fmt.Sprintf("@tcp(%s)", o.Host)
	}

	if o.MaxIdleConns <= 0 {
		panic("max_idle_conns must be greater than zero")
	}
	if o.MaxOpenConns <= 0 {
		panic("max_open_conns must be greater than zero")
	}
	if o.MaxConnMaxLifetime <= 0 {
		panic("max_conn_max_lifetime must be greater than zero")
	}
}

func (o *Opt) ConnectionString() string {
	return fmt.Sprintf("%s:%s%s/%s?parseTime=true", o.User, o.Password, o.Host, o.Name)
}
```

4. создаем файл pkg/database/logger.go

```go
// тут пишем простую обертку для логирования запросов

package database

type Logger struct {
	callback func(format string, v ...interface{})
}

func (l *Logger) SetCallback(callback func(format string, v ...interface{})) {
	l.callback = callback
}

// Printf данный метод реализует интефейс Logger для печати информации
func (l *Logger) Printf(format string, v ...interface{}) {
	l.callback(format, v)
}
```

5. создаем файл pkg/database/pool.go

```go
// тут мы пишем интефейс, для получения доступа к пулу коннектов базы данных
// его мы и будем использовать, а не частный случай MySQL

package database

import builder "github.com/doug-martin/goqu/v9"

type Pool interface {
	Builder() *builder.Database
}
```

6. создаем файл pkg/repository/repo.go

```go
package repository

import (
	"context"
	"database/sql"
	"time"

	builder "github.com/doug-martin/goqu/v9"

	"github.com/grip211/crud/pkg/commands"
	"github.com/grip211/crud/pkg/database"
)

var ErrNotFound = errors.New("not found")

// Repo пишем структуру которая будет реализовывать все нужные нам методы для работы с базой данных
// в нашем случае Create, Read, Update, Delete
type Repo struct {
	db database.Pool
}

func New(db database.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(ctx context.Context, command *commands.CreateCommand) error {
	_, err := r.db.Builder().
		Insert("productdb.Products").
		Rows(builder.Record{
			"model":   command.Model,
			"company": command.Company,
			"price":   command.Price,
		}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Read(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Builder().
		Select(
			builder.C("id"),
			builder.C("company"),
			builder.C("model"),
			builder.C("price"),
		).
		From("productdb.Products").
		ScanStructsContext(ctx, &products)

	if err != nil {
		return nil, err
	}
	return products, nil
}
```

7. создаем файл pkg/commands/commands.go

```go
// в этом файле будем описывать структуры для более удобной манипуляции с входными и выходными данными

package commands

type CreateCommand struct {
	Model   string
	Company string
	Price   float32
}

func NewCreateCommand(model, company, price string) (*CreateCommand, error) {
	val, err := strconv.ParseFloat(price, 32);
	if err == nil {
		return nil, err
	}

	return &CreateCommand{
		Model:   model,
		Company: company,
		Price:   float32(val),
	}, nil
}
```

8. создаем файл pkg/models/model.go

```go
// тут мы будем описывать структуры для чтения

package models

// Product перемещаем сюда нашу структуру из main
type Product struct {
	Id      int     `db:"id"`
	Model   string  `db:"model"`
	Company string  `db:"company"`
	Price   float32 `db:"price"`
}
```

10. обновляем main.go

```go
// вместо IndexHandler

package main

func buildIndexHandler(repo *repository.Repo) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := repo.Read(r.Context())
		if err != nil {
			log.Println(err)
			return
		}

		tmpl, _ := template.ParseFiles("templates/index.html")
		tmpl.Execute(w, products)
	}
}
```

далее в main.go

```go
package main

func Main(ctx *cli.Context) error {
	appContext, cancel := context.WithCancel(ctx.Context)
	defer func() {
		cancel()
		<-time.After(time.Second * 1)
	}()

	await, stop := signal.Notifier(func() {
		fmt.Println("received a system signal, start shutdown process..")
	})

	conn, err := mysql.New(b.Context(), &database.Opt{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		Dialect:  "mysql",
		Debug:    true,
	})

	repo := repository.New(conn)

	if err := db.PingContext(appContext); err != nil {
		stop(err)
	}

	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/", buildIndexHandler(repo))
		//....

		http.Handle("/", router)

		if err := http.ListenAndServe(":8181", nil); err != nil {
			stop(err)
		}
	}()

	return await()
}
```