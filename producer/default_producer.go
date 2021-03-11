package producer

import (
	"../consumer"
	"../packet"
	"fmt"
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

func (d *DefaultProducer) Wait() {
	d.WaitGroup.Wait()
}

func (d *DefaultProducer) fetchAndReturn(initialArgs []interface{}) {
	var packets []packet.Packet
	var shouldContinue bool
	var args []interface{}
	args = initialArgs
	fmt.Println(args)

	for {
		fmt.Println("Restarting the loop")
		packets, shouldContinue, args = d.ProducerFunc(args...)
		log.Println("Fetched first set of data!")
		log.Printf("Should Continue: %s\n\n", shouldContinue)
		if !shouldContinue {
			d.WaitGroup.Done()
			return
		}
		for _, packet := range packets {
			//opChannel <- packet
			d.Consumer.Consume(packet)
			if packet.Err != nil {
				log.Fatalf("error occured while producing data.")

				if d.ShouldContinueOnError {
					continue
				}
				return
			}
		}
		fmt.Println("Calling select!!")
		//select {
		//case <-quitChan:
		//log.Println("Received request to quit producer.")
		//d.WaitGroup.Done()
		//return
		//}
		fmt.Println("Reiterating!!")
	}
}

func (d *DefaultProducer) Produce(initialArgs []interface{}) {
	log.Println("Producer started to produce data.")
	d.WaitGroup.Add(1)
	go d.fetchAndReturn(initialArgs)
}
