package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Event      string    `json:"event"`
	Attraction string    `json:"attraction"`
	Timestamp  time.Time `json:"timestamp"`
	Value      string    `json:"value"`
}

type Population struct {
	Attraction string `json:"attraction"`
	Population int    `json:"population"`
	Capacity   int    `json:"capacity"`
	State      string `json:"state"`
}

const (
	POPULATION_MONITORING string = "population_monitoring"
	QUEUE_EMPTY           string = "queue_empty"
	QUEUE_ALMOST_EMPTY    string = "queue_almost_empty"
	QUEUE_ALMOST_FULL     string = "queue_almost_full"
	QUEUE_FULL            string = "queue_full"
)

var actual_queue_populations map[string]Population = make(map[string]Population)
var message_history []Message = make([]Message, 0)

func main() {
	time.Sleep(time.Second * 5) // Waiting for RabbitMQ to be approachable

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.StripSlashes)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*", "https://*"},
		AllowedMethods: []string{"GET"},
	}))

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

	queue_max_capacity, err := strconv.Atoi(os.Getenv("QUEUE_MAX_CAPACITY"))
	if err != nil {
		panic(err)
	}

	ack_messages, err := strconv.ParseBool(os.Getenv("ACK_MESSAGES_IN_API"))
	if err != nil {
		panic(err)
	}

	go func() {
		for msg := range messages {
			var decodedMessage Message
			if err := json.Unmarshal(msg.Body, &decodedMessage); err != nil {
				slog.Error(fmt.Sprintf("Could not decode message, err: `%s`", err.Error()))
				continue
			}

			attraction, ok := actual_queue_populations[decodedMessage.Attraction]
			if !ok {
				attraction.Capacity = queue_max_capacity
				attraction.Attraction = decodedMessage.Attraction
			}

			if decodedMessage.Event == POPULATION_MONITORING {
				population, err := strconv.Atoi(decodedMessage.Value)
				if err != nil {
					slog.Error(fmt.Sprintf("Could not convert population value to int, err: `%s`", err.Error()))
				}
				attraction.Population = population
			} else {
				attraction.State = decodedMessage.Event
			}

			actual_queue_populations[decodedMessage.Attraction] = attraction

			if len(message_history) == 10000 {
				message_history = message_history[:9998]
			}

			message_history = append(message_history, decodedMessage)

			if !ack_messages {
				continue
			}

			if err := msg.Ack(false); err != nil {
				slog.Error(fmt.Sprintf("Could not acknowledge message, err: `%s`", err.Error()))
			}
		}
	}()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Pong!")
	})

	r.Get("/populations", func(w http.ResponseWriter, r *http.Request) {
		var populations []Population
		for pop := range maps.Values(actual_queue_populations) {
			populations = append(populations, pop)
		}

		raw, err := json.Marshal(populations)
		if err != nil {
			slog.Error(fmt.Sprintf("Could not serialize response into json, err: `%s`", err))
			return
		}

		fmt.Fprint(w, string(raw))
	})

	r.Get("/messages", func(w http.ResponseWriter, r *http.Request) {
		arg_attraction_names := r.URL.Query().Get("attractions")
		attraction_names := []string{}
		if strings.Contains(arg_attraction_names, ",") {
			attraction_names = strings.Split(arg_attraction_names, ",")
		} else if arg_attraction_names != "" {
			attraction_names = append(attraction_names, arg_attraction_names)
		}

		arg_event_names := r.URL.Query().Get("events")
		event_names := []string{}
		if strings.Contains(arg_event_names, ",") {
			event_names = strings.Split(arg_event_names, ",")
		} else if arg_event_names != "" {
			event_names = append(event_names, arg_event_names)
		}

		from_timestamp, _ := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

		message_history_reversed := slices.Clone(message_history)
		slices.Reverse(message_history_reversed)

		history := []Message{}
		for _, event := range message_history_reversed {
			if limit > 0 && len(history) == limit {
				break
			}

			if !from_timestamp.IsZero() && event.Timestamp.Before(from_timestamp) {
				continue
			}

			if len(attraction_names) > 0 && !slices.Contains(attraction_names, event.Attraction) {
				continue
			}

			if len(event_names) > 0 && !slices.Contains(event_names, event.Event) {
				continue
			}

			history = append(history, event)
		}
		raw, err := json.Marshal(history)
		if err != nil {
			slog.Error(fmt.Sprintf("Could not serialize response into json, err: `%s`", err))
			return
		}

		fmt.Fprint(w, string(raw))
	})

	http.ListenAndServe(":80", r)
}
