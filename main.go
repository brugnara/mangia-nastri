package main

import (
	"net/http"

	"github.com/brugnara/mangia-nastri/logger"
	"github.com/brugnara/mangia-nastri/proxy"
)

var log = logger.New()

func main() {
	http.Handle("/", proxy.New())

	log.Print("\nMangia_nastri is ready to rock.\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
