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
	var seenNames = make(map[string]bool)
	log.Info("Starting MangiaNastri...")

	cmd := commander.New(config.Commander.Port, log.CloneWithPrefix("commander"))

	// wait until commander is ready
	<-cmd.Ready

	for index, p := range config.Proxy {
		if _, ok := seenNames[p.Name]; ok {
			log.Error("Proxy name already used", "name", p.Name)
			panic("Proxy name must be unique")
		}
		seenNames[p.Name] = true

		wg.Add(1)
		go func(p conf.Proxy, mux *http.ServeMux, log logger.Logger) {
			log.Info("Starting proxy", "name", p.Name)

			prx := proxy.New(&p, log)
			mux.Handle("/", prx)

			server := &http.Server{
				Addr:    fmt.Sprintf(":%d", p.Port),
				Handler: mux,
			}

			cmd.Subscribe(prx.Name, prx.Action)

			log.Info("Waiting for proxy to be ready", "name", p.Name)
			<-prx.Ready
			log.Info("Proxy is ready", "name", p.Name)

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
