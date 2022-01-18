package consumer

import (
	"math"

	"github.com/adjust/rmq"
)

const handlesInAPeriod = 200

type Calculator struct {
	connection rmq.Connection
	queueName  string
}

func NewCalculator(queueName string, connection rmq.Connection) *Calculator {
	return &Calculator{queueName: queueName, connection: connection}
}

func (c *Calculator) ReplicaSize() int32 {
	queues := c.connection.GetOpenQueues()
	stats := c.connection.CollectStats(queues)
	queueStat := stats.QueueStats[c.queueName]
	sizeFloat := queueStat.ReadyCount / handlesInAPeriod
	size := int32(math.Ceil(float64(sizeFloat)))
	return size
}
