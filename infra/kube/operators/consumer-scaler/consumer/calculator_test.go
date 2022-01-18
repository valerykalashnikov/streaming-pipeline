package consumer_test

import (
	"testing"

	"github.com/adjust/rmq"
	"github.com/valerykalashnikov/streaming-pipeline/infra/kube/operators/consumer-scaler/consumer"
)

type MockRMQConnection struct {
}

func (*MockRMQConnection) CollectStats(queueName []string) (stats rmq.Stats) {
	queueStat := rmq.QueueStat{
		ReadyCount: 1000,
	}
	stats = rmq.Stats{
		QueueStats: rmq.QueueStats{},
	}
	stats.QueueStats["queue1"] = queueStat
	return stats
}
func (m *MockRMQConnection) GetOpenQueues() []string {
	return []string{"queue1"}
}

func (*MockRMQConnection) OpenQueue(name string) rmq.Queue   { panic("error: not supported") }
func (*MockRMQConnection) StopAllConsuming() <-chan struct{} { panic("error: not supported") }

func TestCalculator(t *testing.T) {
	conn := &MockRMQConnection{}
	calc := consumer.NewCalculator("queue1", conn)
	expected := 5
	actual := calc.ReplicaSize()
	if expected != int(actual) {
		t.Errorf("replica size get error: expected %v, actual %v", expected, actual)
	}
}
