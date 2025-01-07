package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
  time.Sleep(time.Second * 5)

  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Use(middleware.StripSlashes)

  conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@rabbitmq", os.Getenv("MQ_USERNAME"), os.Getenv("MQ_PASSWORD")))
  if err != nil {
    panic(err)
  }

  channel, err := conn.Channel()
  if err != nil {
    panic(err)
  }
  defer channel.Close()

  messages, err := channel.Consume(
    "queuepods",
    "api",
    false,
    false,
    false,
    false,
    nil,
  )
  if err != nil {
    panic(err)
  }

  go func() {
    for msg := range messages {
      fmt.Println(string(msg.Body))
    }
  }()

  r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Pong!")
  })

  http.ListenAndServe(":80", r)
}
