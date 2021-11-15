package main

import (
	"net/http"

	"github.com/adjust/rmq/v4"
	"github.com/valerykalashnikov/streaming-pipeline/log"
)

func main() {
	connection, err := rmq.OpenConnection("dashboard", "tcp", "localhost:6379", 2, nil)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/overview", NewDashboardHandler(connection))
	log.Info("Dashboard is rendered on http://localhost:3333/overview")
	if err := http.ListenAndServe(":3333", nil); err != nil {
		panic(err)
	}
}
