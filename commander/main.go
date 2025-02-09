package commander

import (
	"fmt"
	"mangia_nastri/logger"
	"net/http"
)

type Action int

const (
	DO_RECORD Action = iota
	DO_NOT_RECORD
)

type Command struct {
	port int
	subs []chan<- Action

	Ready chan bool
}

func (c *Command) Subscribe(action chan<- Action) {
	c.subs = append(c.subs, action)
}

func New(port int, log logger.Logger) *Command {
	actionChannel := make(chan Action)
	commander := &Command{
		port:  port,
		Ready: make(chan bool),
		subs:  make([]chan<- Action, 0),
	}

	go func() {
		for a := range actionChannel {
			log.Info(fmt.Sprintf("Propagating action to %d subscribers", len(commander.subs)))
			for _, sub := range commander.subs {
				sub <- a
			}
		}
	}()

	go func(port int, mux *http.ServeMux, log logger.Logger) {
		log.Info("Commander is starting on", "port", port)

		mux.Handle("/do-record", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("Received do-record command")
			actionChannel <- DO_RECORD
		}))

		mux.Handle("/do-not-record", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("Received do-not-record command")
			actionChannel <- DO_NOT_RECORD
		}))

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		}

		commander.Ready <- true

		log.Info(server.ListenAndServe())
	}(port, http.NewServeMux(), log)

	return commander
}
