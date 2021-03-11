package producer

import (
	"../consumer"
	"../packet"
	"log"
	"sync"
)

type DefaultProducer struct {
	ShouldContinueOnError bool
	WaitGroup             *sync.WaitGroup
	ProducerFunc          ProducerFunc
	Consumer              consumer.Consumer
}

func NewDefaultProducer(cnsmr consumer.Consumer, shouldContinueOnError bool, producerFunc ProducerFunc) Producer {
	var wg sync.WaitGroup
	return &DefaultProducer{
		ShouldContinueOnError: shouldContinueOnError,
		WaitGroup:             &wg,
		Consumer:              cnsmr,
		ProducerFunc:          producerFunc,
	}
}

// Wait will be used by client to wait for producer to complete its job
func (d *DefaultProducer) Wait() {
	d.WaitGroup.Wait()
}

// fetchAndReturn calls producer function again and again until producer func returns shouldContinue as false
func (d *DefaultProducer) fetchAndReturn(initialArgs []interface{}) {
	var packets []packet.Packet
	var shouldContinue bool
	var args []interface{}
	args = initialArgs

	for {
		packets, shouldContinue, args = d.ProducerFunc(args...)
		for _, packet := range packets {
			d.Consumer.Consume(packet)
			if packet.Err != nil {
				log.Fatalf("error occured while producing data.")

				if d.ShouldContinueOnError {
					continue
				}
				d.WaitGroup.Done()
				return
			}
		}
		if !shouldContinue {
			d.WaitGroup.Done()
			return
		}
	}
}

// Produce concurrently calls producer function after adding waitgroup.
func (d *DefaultProducer) Produce(initialArgs []interface{}) {
	log.Println("Producer started to produce data.")
	d.WaitGroup.Add(1)
	go d.fetchAndReturn(initialArgs)
}
