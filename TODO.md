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

3. **Интегрироваться с github.com/urfave/cli/v2:**

> Или с аналгами типа viper, но мне нравится больше эта поэтому будем работать с ней :)

```go
// Следующий код нужно поместить в pkg/signal/signal.go

package signal

import (
   "os"
   "os/signal"
   "syscall"
)

// Notifier function returns functions for waiting and interrupting the signal for graceful termination
// of the application and the ability to clean up all resources
//
// ```go
//
//	   await, stop := signal.Notifier(func() {
//			 stdout.Info("received a system signal to shut down API server, start the shutdown process..")
//	   })
//
//	   go func() {
//	      if err := someService.Run(ctx); err != nil {
//	         stop(err) // Gracefully closing the application cleaning up the weight of resources
//	      }
//	   }()
//
//	   ...some code
//
//	   if err := await(); err != nil {
//	      some code for handle error
//	   }
//
// ```
func Notifier(on ...func()) (wait func() error, stop func(err ...error)) {
   sig := make(chan os.Signal, 1)
   signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
   err := make(chan error, 1)

   wait = func() error {
      <-sig
      if len(on) > 0 {
         on[0]()
      }
      select {
      case e := <-err:
         return e
      default:
         return nil
      }
   }

   stop = func(e ...error) {
      if len(e) > 0 {
         err <- e[0]
      }
      sig <- syscall.SIGINT
   }

   return wait, stop
}
```

```go
// Следущий код нужно поместить в cmd/crud/main.go, только не все копировать

package main

// тут импорты
// тут еще хендлеры
// ...
// ...

var database *sql.DB

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

   db, err := sql.Open("mysql", "xxxx:xxxxxx@tcp(xxx.xx.xx.xx)/pavel")
   if err != nil {
      stop(err)
   }

   if err := db.PingContext(ctx); err != nil {
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

      http.Handle("/", router)

      if err := http.ListenAndServe(":8181", nil); err != nil {
		  stop(err)
      }
   }()
   
   return await()
}
```