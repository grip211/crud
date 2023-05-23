package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/grip211/crud/pkg/signal"
	"github.com/urfave/cli/v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type Product struct {
	Id      int
	Model   string
	Company string
	Price   int
}

var database *sql.DB

// удаление наименований
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	_, err := database.Exec("delete from productdb.Products where id = ?", id)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)

}

// возвращаем пользователю страницу для редактирования объекта
func EditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := database.QueryRow("select * from productdb.Products where id = ?", id)
	prod := Product{}
	err := row.Scan(&prod.Id, &prod.Model, &prod.Company, &prod.Price)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		tmpl, _ := template.ParseFiles("templates/edit.html")
		tmpl.Execute(w, prod)
	}
}

// получаем измененные данные и сохраняем их в БД
func EditHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	model := r.FormValue("model")
	company := r.FormValue("company")
	price := r.FormValue("price")

	_, err = database.Exec("update productdb.Products set model=?, company=?, price = ? where id = ?",
		model, company, price, id)

	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		model := r.FormValue("model")
		company := r.FormValue("company")
		price := r.FormValue("price")

		_, err = database.Exec("insert into productdb.Products (model, company, price) values (?, ?, ?)",
			model, company, price)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "templates/create.html")
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("select * from productdb.Products")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	products := []Product{}

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, products)
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
		Flags: []cli.Flag{
			// тут будет потом
		},
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

	db, err := sql.Open("mysql", getConnectionString())
	if err != nil {
		stop(err)
	}

	if err := db.PingContext(appContext); err != nil {
		stop(err)
	}

	database = db
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/", IndexHandler)
		router.HandleFunc("/create", CreateHandler)
		router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("GET")
		router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")
		router.HandleFunc("/delete/{id:[0-9]+}", DeleteHandler)

		http.Handle("/", router)

		if err := http.ListenAndServe(":8181", nil); err != nil {
			stop(err)
		}
	}()

	return await()
}
