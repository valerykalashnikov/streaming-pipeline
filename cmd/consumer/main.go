package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/valerykalashnikov/streaming-pipeline/db"
	"github.com/valerykalashnikov/streaming-pipeline/log"
)

const (
	prefetchLimit = 1000
	pollDuration  = 100 * time.Millisecond
	numConsumers  = 5

	reportBatchSize = 10000
	consumeDuration = time.Millisecond
	shouldLog       = false
)

func main() {
	errChan := make(chan error, 10)
	go logErrors(errChan)

	db, err := db.GetDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Successfully connected to database")

	connection, err := createRMQConnection(errChan)
	if err != nil {
		log.Fatal(err)
	}

	queue, err := connection.OpenQueue("data")
	if err != nil {
		log.Fatal(err)
	}

	if err := queue.StartConsuming(prefetchLimit, pollDuration); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < numConsumers; i++ {
		name := fmt.Sprintf("consumer %d", i)
		if _, err := queue.AddConsumer(name, NewConsumer(db, i)); err != nil {
			log.Fatal(err)
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	defer signal.Stop(signals)

	<-signals // wait for signal
	go func() {
		<-signals // hard exit on second signal (in case shutdown gets stuck)
		os.Exit(1)
	}()

	<-connection.StopAllConsuming() // wait for all Consume() calls to finish
}

func createRMQConnection(errChan chan error) (rmq.Connection, error) {
	rand.Seed(time.Now().UnixNano())
	// generate the postfix with a lenght of 5 bytes
	b := make([]byte, 5)
	rand.Read(b)

	connName := "consumer" + fmt.Sprintf("%x", b)[:5]
	return rmq.OpenConnection(connName, "tcp", "localhost:6379", 2, errChan)
}

func logErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {
		case *rmq.HeartbeatError:
			if err.Count == rmq.HeartbeatErrorLimit {
				log.Error("heartbeat error (limit): ", err)
			} else {
				log.Error("heartbeat error: ", err)
			}
		case *rmq.ConsumeError:
			log.Error("consume error: ", err)
		case *rmq.DeliveryError:
			log.Error("delivery error: ", err.Delivery, err)
		default:
			log.Error("other error: ", err)
		}
	}
}
