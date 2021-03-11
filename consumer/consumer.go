package consumer

import (
	"../packet"
)

func init() {
	New = NewDefaultConsumer
}

type ConsumerFunc func(packet.Packet) error

var New func(reqHandlerCount int, maxPackets int, consumerFunc ConsumerFunc) Consumer

type Consumer interface {
	Consume(packet.Packet)
	Wait()
}
