package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/valerykalashnikov/streaming-pipeline/log"
)

type Consumer struct {
	name   string
	count  int
	before time.Time
	db     *sql.DB
}

func NewConsumer(db *sql.DB, tag int) *Consumer {
	return &Consumer{
		name:   fmt.Sprintf("consumer%d", tag),
		count:  0,
		before: time.Now(),
		db:     db,
	}
}

func (consumer *Consumer) Consume(delivery rmq.Delivery) {
	payload := delivery.Payload()
	log.Debug("start consume %s", payload)
	time.Sleep(consumeDuration)

	consumer.count++
	if consumer.count%reportBatchSize == 0 {
		duration := time.Since(consumer.before)
		consumer.before = time.Now()
		perSecond := time.Second / (duration / reportBatchSize)
		log.Info("%s consumed %d %s %d", consumer.name, consumer.count, payload, perSecond)
	}

	if consumer.count%reportBatchSize > 0 {
		if err := delivery.Ack(); err != nil {
			log.Info("failed to ack %s: %s", payload, err)
		} else {
			record, err := NewRecordFromString(payload)
			if err != nil {
				log.Error(fmt.Sprintf("canot create record from %s", payload))
			}
			err = SaveUser(consumer.db, record)
			if err != nil {
				log.Error(fmt.Sprintf("canot save record %s, record %v", err, record))
			}
		}
	} else { // reject one per batch
		if err := delivery.Reject(); err != nil {
			log.Info("failed to reject %s: %s", payload, err)
		} else {
			log.Info("rejected %s", payload)
		}
	}
}
