package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Use(middleware.StripSlashes)

  conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@rabbitmq:5672/", os.Getenv("MQ_USERNAME"), os.Getenv("MQ_PASSWORD")))
  if err != nil {
    panic(err)
  }

  channel, err := conn.Channel()
  defer channel.Close()

  queue, err := channel.QueueDeclare(
    "queuepods",
    false,
    false,
    false,
    false,
    nil,
  )

  r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Pong!")
  })

  http.ListenAndServe(":80", r)
}
