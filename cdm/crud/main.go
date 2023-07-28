package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/grip211/crud/pkg/apperror"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/urfave/cli/v2"

	"github.com/grip211/crud/pkg/commands"
	"github.com/grip211/crud/pkg/database"
	"github.com/grip211/crud/pkg/database/mysql"
	"github.com/grip211/crud/pkg/repository"
	"github.com/grip211/crud/pkg/signal"
)

// удаление наименований
func buildRestDeleteHandler(repo *repository.Repo) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		command, err := commands.NewDeleteCommand(id)
		if err != nil {
			// убрать после того как добавишь обработку ошибок в ErrorHandler
			return err
		}

		_, err = repo.Delete(ctx.Context(), command)
		if err != nil {
			// убрать после того как добавишь обработку ошибок в ErrorHandler
			return err
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": 301})
	}
}

func buildRestEditPageHandler(repo *repository.Repo) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")

		iid, err := strconv.Atoi(id)
		if err != nil {
			// убрать после того как добавишь обработку ошибок в ErrorHandler
			return err
		}

		prod, err := repo.ReadOneWithFeatures(ctx.Context(), iid)
		if err != nil {
			return ctx.Status(http.StatusNotFound).SendString("NotFound")
		}

		return ctx.JSON(prod)
	}
}

type EditForm struct {
	ID       string `form:"id"`
	Model    string `form:"model"`
	Company  string `form:"company"`
	Quantity string `from:"quantity"`
	Price    string `form:"price"`

	CPU     string `form:"cpu"`
	Memory  string `form:"memory"`
	Display string `form:"display"`
	Camera  string `form:"camera"`
}

// получаем измененные данные и сохраняем их в БД
func buildRestEditHandler(repo *repository.Repo) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		edit := &EditForm{}
		if err := ctx.BodyParser(edit); err != nil {
			return err
		}

		updateCommand, err := commands.NewUpdateCommand(
			edit.ID,
			edit.Model,
			edit.Company,
			edit.Quantity,
			edit.Price,
			edit.CPU,
			edit.Memory,
			edit.Display,
			edit.Camera,
		)
		if err != nil {
			// убрать после того как добавишь обработку ошибок в ErrorHandler
			return err
		}

		err = repo.Update(ctx.Context(), updateCommand)
		if err != nil {
			// убрать после того как добавишь обработку ошибок в ErrorHandler
			return err
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": 301})
	}
}

type CreatForm struct {
	ID       string `form:"id"`
	Model    string `form:"model"`
	Company  string `form:"company"`
	Quantity string `from:"quantity"`
	Price    string `form:"price"`

	CPU     string `form:"cpu"`
	Memory  string `form:"memory"`
	Display string `form:"display"`
	Camera  string `form:"camera"`
}

func buildRestCreateHandler(repo *repository.Repo) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Method() == "POST" {
			creat := &CreatForm{}
			if err := ctx.BodyParser(creat); err != nil {
				return err
			}

			createCommand, err := commands.NewCreteCommand(
				creat.Model,
				creat.Company,
				creat.Quantity,
				creat.Price,
				creat.CPU,
				creat.Memory,
				creat.Display,
				creat.Camera,
			)
			if err != nil {
				// убрать после того как добавишь обработку ошибок в ErrorHandler
				return apperror.ErrEndFound
			}

			_, err = repo.Create(ctx.Context(), createCommand)
			if err != nil {
				// убрать после того как добавишь обработку ошибок в ErrorHandler
				return err
			}
			return ctx.Redirect("/", 301)
		}

		return ctx.JSON("create")
	}
}

/*
func buildIndexHandler(repo *repository.Repo) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		products, err := repo.Read(ctx.Context())
		if err != nil {
			// убрать после того как добавишь обработку ошибок в ErrorHandler
			return apperror.ErrEndFound
		}

		return ctx.Render("index", fiber.Map{
			"Products": products,
		})
	}
}

*/

func buildRestIndexHandler(repo *repository.Repo) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		products, err := repo.Read(ctx.Context())
		if err != nil {
			return err
		}
		return ctx.JSON(products)
	}
}

func buildRestFeatureHandler(repo *repository.Repo) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")

		iid, err := strconv.Atoi(id)
		if err != nil {
			// убрать после того как добавишь обработку ошибок в ErrorHandler
			return err
		}

		product, err := repo.ReadOneWithFeatures(ctx.Context(), iid)
		if err != nil {
			// убрать после того как добавишь обработку ошибок в ErrorHandler
			return err
		}

		return ctx.JSON(product)
	}
}

func main() {
	application := &cli.App{
		Flags:  []cli.Flag{},
		Action: Main,
	}
	if err := application.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Main(ctx *cli.Context) error {
	appContext, cancel := context.WithCancel(ctx.Context)
	defer func() {
		cancel()
		<-time.After(time.Second * 1)
	}()

	await, stop := signal.Notifier(func() {
		fmt.Println("received a system signal, start shutdown process..")
	})

	conn, err := mysql.New(appContext, &database.Opt{
		Host:               os.Getenv("DB_Host"),
		User:               os.Getenv("DB_USER"),
		Password:           os.Getenv("DB_PASS"),
		Name:               os.Getenv("DB_NAME"),
		Dialect:            "mysql",
		MaxConnMaxLifetime: time.Minute * 5,
		MaxOpenConns:       10,
		MaxIdleConns:       9,
		Debug:              true,
	})
	if err != nil {
		return err
	}

	repo := repository.New(conn)

	go func() {
		engine := html.New("./templates", ".html")

		server := fiber.New(fiber.Config{
			Views: engine,
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				// тут будут печататься все ошибки с хедлеров
				fmt.Println(err)

				// показываем страницу ошибки
				return ctx.Render("error", nil)
			},
		})

		v1 := server.Group("/api/v1")
		v1.Get("/", buildRestIndexHandler(repo))
		v1.Post("/create", buildRestCreateHandler(repo))
		v1.Post("/edit/:id", buildRestEditHandler(repo))
		v1.Delete("/delete/:id", buildRestDeleteHandler(repo))
		v1.Get("/feature/:id", buildRestFeatureHandler(repo))

		server.Get("/", buildRestIndexHandler(repo))
		server.Get("/create", buildRestCreateHandler(repo))
		server.Get("/delete/:id", buildRestDeleteHandler(repo))
		server.Get("/edit/:id", buildRestEditPageHandler(repo))
		server.Get("/feature/:id", buildRestFeatureHandler(repo))

		server.Post("/edit/:id?", buildRestEditHandler(repo))
		server.Post("/create", buildRestCreateHandler(repo))

		ln, err := signal.Listener(appContext, 1, "/tmp/crud.sock", ":8181")
		if err != nil {
			stop(err)
			return
		}

		if err := server.Listener(ln); err != nil {
			stop(err)
		}
	}()
	return await()
}
