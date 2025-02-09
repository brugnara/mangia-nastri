package main

import (
	"fmt"
	"net/http"
	"sync"

	"mangia_nastri/commander"
	"mangia_nastri/conf"
	"mangia_nastri/logger"
	"mangia_nastri/proxy"
)

var log = logger.New("mangianastri")
var config = conf.New("./conf.yaml")

func main() {
	var wg sync.WaitGroup
	log.Info("Starting MangiaNastri...")

	cmd := commander.New(config.Commander.Port, log.CloneWithPrefix("commander"))

	// wait until commander is ready
	<-cmd.Ready

	for index, p := range config.Proxy {
		wg.Add(1)
		go func(p conf.Proxy, mux *http.ServeMux, log logger.Logger) {
			log.Info("Starting proxy", "name", p.Name)

			prx := proxy.New(&p, log)
			mux.Handle("/", prx)

			server := &http.Server{
				Addr:    fmt.Sprintf(":%d", p.Port),
				Handler: mux,
			}

			cmd.Subscribe(prx.Action)

			wg.Done()

			log.Info(server.ListenAndServe())

		}(p, http.NewServeMux(), log.CloneWithPrefix(fmt.Sprintf("#%d", index)))
	}

	// wait until all proxies are ready
	wg.Wait()

	log.Info("MangiaNastri is ready to rock...")
	// loop forever
	select {}
}
