package producer

import (
	"log"
	"sync"
)

type DefaultProducer struct {
	ShouldContinueOnError bool
	WaitGroup             sync.WaitGroup
}

func NewDefaultProducer(shouldContinueOnError bool, wg WaitGroup) DefaultProducer {
	return DefaultProducer{
		ShouldContinueOnError: shouldContinueOnError,
		WaitGroup:             wg,
	}
}

func (d DefaultProducer) fetchAndReturn(opChannel chan Data, producerFunc ProducerFunc, initialArgs []interface{}, quitChan chan bool) {
	var packet Data
	var shouldContinue bool
	var args []interface{}
	var err error

	for {
		packet, shouldContinue, args = producerFunc(initialArgs)
		opChannel <- packet
		if packet.Err != nil {
			log.Fatalf("error occured while producing data.")

			if d.shouldContinueOnError {
				continue
			}
			return
		}

		select {
		case <-quitChan:
			log.Println("Received request to quit producer.")
			d.WaitGroup.Done()
			return
		}
	}
}

func (DefaultProducer) Produce(opChannel chan Data, producerFunc ProducerFunc, initialArgs []interface{}, quitChan chan bool) {
	log.Println("Producer started to produce data.")
	go fetchAndReturn(opChannel, producerFunc, initialArgs, quitChan)
}
