package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/grip211/crud/pkg/commands"
	"github.com/grip211/crud/pkg/database"
	"github.com/grip211/crud/pkg/database/mysql"
	"github.com/grip211/crud/pkg/repository"
	"github.com/grip211/crud/pkg/signal"
	"github.com/urfave/cli/v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// удаление наименований
func buildDeleteHandler(repo *repository.Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		command, err := commands.NewDeleteCommand(id)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = repo.Delete(r.Context(), command)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	}
}

// возвращаем пользователю страницу для редактирования объекта
func buildEditPageHandler(repo *repository.Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		iid, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println(err)
			return
		}

		prod, err := repo.ReadOne(r.Context(), iid)

		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(404), http.StatusNotFound)
		} else {
			tmpl, _ := template.ParseFiles("templates/edit.html")
			tmpl.Execute(w, prod)
		}
	}
}

// получаем измененные данные и сохраняем их в БД
func buildEditHandler(repo *repository.Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			return
		}
		id := r.FormValue("id")
		model := r.FormValue("model")
		company := r.FormValue("company")
		quantity := r.FormValue("quantity")
		price := r.FormValue("price")

		updateCommand, err := commands.NewUpdateCommand(id, model, company, quantity, price)
		if err != nil {
			log.Println(err)
			return
		}

		err = repo.Update(r.Context(), updateCommand)
		if err != nil {
			log.Println(err)
			return
		}
		http.Redirect(w, r, "/", 301)
	}
}

func buildCreateHandler(repo *repository.Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				return
			}
			model := r.FormValue("model")
			company := r.FormValue("company")
			quantity := r.FormValue("quantity")
			price := r.FormValue("price")

			createCommand, err := commands.NewCreteCommand(model, company, quantity, price)
			if err != nil {
				log.Println(err)
				return
			}

			err = repo.Create(r.Context(), createCommand)
			if err != nil {
				log.Println(err)
				return
			}
			http.Redirect(w, r, "/", 301)
		} else {
			http.ServeFile(w, r, "templates/create.html")
		}
	}
}

func buildIndexHandler(repo *repository.Repo) func(w http.ResponseWriter, r *http.Request) {
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

func getConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
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
		router := mux.NewRouter()
		router.HandleFunc("/", buildIndexHandler(repo))
		router.HandleFunc("/create", buildCreateHandler(repo))
		router.HandleFunc("/edit/{id:[0-9]+}", buildEditPageHandler(repo)).Methods("GET")
		router.HandleFunc("/edit/{id:[0-9]+}", buildEditHandler(repo)).Methods("POST")
		router.HandleFunc("/delete/{id:[0-9]+}", buildDeleteHandler(repo))

		http.Handle("/", router)

		if err := http.ListenAndServe(":8181", nil); err != nil {
			stop(err)
		}
	}()

	return await()
}
