package commander

import (
	"fmt"
	"mangia_nastri/logger"
	"net/http"
	"strings"
)

type Action int

const (
	DO_RECORD Action = iota
	DO_NOT_RECORD
)

type Command struct {
	port int
	subs []internalSubscriber

	Ready chan bool
}

type internalSubscriber struct {
	action chan<- Action
	proxy  string
}
type internalCommand struct {
	action Action
	proxy  string
}

func (c *Command) Subscribe(proxyName string, action chan<- Action) {
	c.subs = append(c.subs, internalSubscriber{
		action: action,
		proxy:  proxyName,
	})
}

func New(port int, log logger.Logger) *Command {
	actionChannel := make(chan internalCommand)
	commander := &Command{
		port:  port,
		Ready: make(chan bool),
		subs:  make([]internalSubscriber, 0),
	}

	go func() {
		for a := range actionChannel {
			var count = 0
			for _, sub := range commander.subs {
				var propagate = false

				if sub.proxy == a.proxy || a.proxy == "*" {
					propagate = true
				}

				if !propagate && strings.HasSuffix(a.proxy, "*") {
					if strings.HasPrefix(sub.proxy, strings.TrimSuffix(a.proxy, "*")) {
						propagate = true
					}
				}

				if !propagate && strings.HasPrefix(a.proxy, "*") {
					if strings.HasSuffix(sub.proxy, strings.TrimPrefix(a.proxy, "*")) {
						propagate = true
					}
				}

				if propagate {
					sub.action <- a.action
					count++
				}
			}
			log.Info(fmt.Sprintf("Propagated action to %d/%d subscribers", count, len(commander.subs)))

		}
	}()

	go func(port int, mux *http.ServeMux, log logger.Logger) {
		log.Info("Commander is starting on", "port", port)

		mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("Received command")
			// command received should be in the format
			// /<proxy-name>/<command>
			// or glob pattern like
			// /<proxy-name>*/<command>
			// or
			// /*<proxy-name-suffix>/<command>

			// extract the proxy name and the command
			log.Info("Received command", "url", r.URL.Path)
			var proxyName string
			var command Action

			parts := strings.Split(r.URL.Path, "/")
			if len(parts) < 3 {
				log.Error("Invalid command", "url", r.URL.Path)
				http.Error(w, "Invalid command", http.StatusBadRequest)
				return
			}

			proxyName = parts[1]

			if parts[2] == "do-record" {
				command = DO_RECORD
			} else if parts[2] == "do-not-record" {
				command = DO_NOT_RECORD
			} else {
				log.Error("Invalid command", "url", r.URL.Path)
				http.Error(w, "Invalid command", http.StatusBadRequest)
				return
			}

			actionChannel <- internalCommand{
				action: command,
				proxy:  proxyName,
			}
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
