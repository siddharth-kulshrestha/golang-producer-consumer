package producer

import (
	"../consumer"
	"../packet"
)

func init() {
	New = NewDefaultProducer
}

type ProducerFunc func(args ...interface{}) ([]packet.Packet, bool, []interface{})

var New func(consumer.Consumer, bool, ProducerFunc) Producer

type Producer interface {
	Produce(initialArgs []interface{})
	Wait()
}
