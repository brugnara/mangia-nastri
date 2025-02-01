package main

import (
	"fmt"
	"net/http"

	"mangia_nastri/conf"
	"mangia_nastri/logger"
	"mangia_nastri/proxy"
)

var log = logger.New("mangianastri")
var config = conf.New("./conf.yaml")

func main() {
	for index, p := range config.Proxy {
		go func(p conf.Proxy, mux *http.ServeMux, log logger.Logger) {
			log.Info("Starting proxy on", "port", p.Port)

			mux.Handle("/", proxy.New(&p, log))

			server := &http.Server{
				Addr:    ":" + p.Port,
				Handler: mux,
			}

			log.Info(server.ListenAndServe())
		}(p, http.NewServeMux(), log.CloneWithPrefix(fmt.Sprintf("#%v", index)))
	}

	log.Info("MangiaNastri is ready to rock...")
	// loop forever
	select {}
}
