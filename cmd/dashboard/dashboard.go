package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adjust/rmq/v4"
)

type DashboarHandler struct {
	connection rmq.Connection
}

func NewDashboardHandler(connection rmq.Connection) *DashboarHandler {
	return &DashboarHandler{connection: connection}
}

func (dashboard *DashboarHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	layout := request.FormValue("layout")
	refresh := request.FormValue("refresh")

	queues, err := dashboard.connection.GetOpenQueues()
	if err != nil {
		log.Fatal(err)
	}

	stats, err := dashboard.connection.CollectStats(queues)
	if err != nil {
		log.Fatal()
	}

	fmt.Fprint(writer, stats.GetHtml(layout, refresh))
}
