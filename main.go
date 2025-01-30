package main

import (
	"net/http"

	"mangia_nastri/conf"
	"mangia_nastri/logger"
	"mangia_nastri/proxy"
)

var log = logger.New("main")
var config = conf.New("./conf.yaml")

func main() {
	http.Handle("/", proxy.New(config))

	log.Print("\nMangia_nastri is ready to rock!!!\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
