package main

import (
	"net/http"

	"mangia_nastri/logger"
	"mangia_nastri/proxy"
)

var log = logger.New()

func main() {
	http.Handle("/", proxy.New())

	log.Print("\nMangia_nastri is ready to rock!!!\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
