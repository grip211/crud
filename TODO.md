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

3. **Интегрироваться с github.com/gofiber/fiber/v2:**

- создаем файл pkg/signal/signal.go

```go
package signal

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	
	"github.com/gofiber/fiber/v2"
)

const (
	ListenerTCP = iota + 1
	ListenerUDS
)

func Listener(ctx context.Context, listener int, uds, tcp string) (net.Listener, error) {
	if listener == ListenerUDS {
		defer maybeChmodSocket(ctx, uds)
		ln, err := listenToUnix(uds)

		return ln, err
	}
	if !strings.Contains(tcp, ":") {
		tcp = ":" + tcp
	}

	ln, err := net.Listen(fiber.NetworkTCP4, tcp)
	if err != nil {
		return nil, err
	}

	return ln, nil
}

func maybeChmodSocket(c context.Context, sock string) {
	// on Linux and similar systems, there may be problems with the rights to the UDS socket
	go func() {
		ctx, cancel := context.WithTimeout(c, 500*time.Millisecond)
		defer cancel()

		var tryCount uint

		fmt.Println("run chmod")
		defer fmt.Println("stop chmod")

		for {
			select {
			case <-ctx.Done():
				fmt.Println("context is canceled")
				return
			case <-time.After(time.Millisecond * 100):
				fmt.Println(fmt.Sprintf("loop %d for chmod unix socket (%s)", tryCount, sock))

				if err := os.Chmod(sock, 0o666); err != nil {
					fmt.Println(err)
					continue
				}

				_, err := os.Stat(sock)
				// if the file exists and it already has permissions
				if err == nil {
					fmt.Println(fmt.Sprintf("unix socket (%s) is ready for listen", sock))
					return
				}

				tryCount++
				if tryCount > 5 {
					fmt.Println("try count is outside for chmod unix socket")
					return
				}
			}
		}
	}()

	_ = os.Chmod(sock, 0o666)
}

func listenToUnix(bind string) (net.Listener, error) {
	_, err := os.Stat(bind)
	if err == nil {
		err = os.Remove(bind)
		if err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return net.Listen("unix", bind)
}
```

- далее заменяем стандарный сервер net/http в функции Main

```go
// Create a new engine by passing the template folder
// and template extension using <engine>.New(dir, ext string)
engine := html.New("./templates", ".html")

// We also support the http.FileSystem interface
// See examples below to load templates from embedded files
engine := html.NewFileSystem(http.Dir("./templates"), ".html")

// Reload the templates on each render, good for development
engine.Reload(true) // Optional. Default: false

// Debug will print each template that is parsed, good for debugging
engine.Debug(true) // Optional. Default: false

// Layout defines the variable name that is used to yield templates within layouts
engine.Layout("embed") // Optional. Default: "embed"

// Delims sets the action delimiters to the specified strings
engine.Delims("{{", "}}") // Optional. Default: engine delimiters

// After you created your engine, you can pass it to Fiber's Views Engine
server := fiber.New(fiber.Config{
    Views: engine,
})

// To render a template, you can call the ctx.Render function
// Render(tmpl string, values interface{}, layout ...string)
server.Get("/", buildIndexHandler(repo))

go func() {
		ln, err := signal.Listener(
			appContext,
			1, "/tmp/crud.sock", ":8181",
		)

		if err != nil {
			stop(err)
			return
		}

		if err := server.Listener(ln); err != nil {
			stop(err)
		}
}()
```

- заменяем хендлеры

```go
func buildIndexHandler(repo *repository.Repo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		products, err := repo.Read(r.Context())
		if err != nil {
			log.Println(err)
			return
		}

        return c.Render("index", fiber.Map{
            "Products": products,
        })
	}
}
```

- в html index надо будет заменить строку 13

```html
{{range . }}
```

на

```html
{{range .Products }}
```

- остальные хедлеры по тому же типу, но надо переделать парсинг форм как вот тут https://docs.gofiber.io/api/ctx/#bodyparser

например, редактирование

```go
type EditForm struct {
	ID      string `form:"name"`
    Model   string `form:"name"`
    Company string `form:"pass"`
    Quantity string `from:"quantity"`
	Price    string `form:"price"`
}

server.Get("/edit/:id", buildEditHandler(repo))

...


func buildEditHandler(repo *repository.Repo) fiber.Handler {
return func(ctx *fiber.Ctx) error {
    id := ctx.Params("id") // "56"
    iid, err := strconv.Atoi(id)
    if err != nil {
        return err
    }

    edit := &EditForm{}
if err := ctx.BodyParser(edit); err != nil {
    return err
}

updateCommand, err := commands.NewUpdateCommand(edit.ID, ...)
if err != nil {
    log.Println(err)
    return err
}

// ...
}
}
```