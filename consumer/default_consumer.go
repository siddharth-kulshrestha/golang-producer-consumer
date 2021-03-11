package consumer

import (
	"fmt"
	"sync"

	"../packet"
)

type DefaultConsumer struct {
	WaitGroup    *sync.WaitGroup
	IpChannel    chan packet.Packet
	ReqHandlers  int
	QuitChannel  chan bool
	ConsumerFunc ConsumerFunc
}

// NewDefaNewDefaultConsumer will return consumer where we can configure reqHandlerCount, maxPackets, and consumerFunc
func NewDefaultConsumer(reqHandlerCount int, maxPackets int, consumerFunc ConsumerFunc) Consumer {
	ipChannel := make(chan packet.Packet, maxPackets)
	quitChannel := make(chan bool)
	var wg sync.WaitGroup
	dc := &DefaultConsumer{
		WaitGroup:    &wg,
		IpChannel:    ipChannel,
		QuitChannel:  quitChannel,
		ReqHandlers:  reqHandlerCount,
		ConsumerFunc: consumerFunc,
	}
	for i := 0; i < reqHandlerCount; i++ {
		go dc.handlePackets()
	}
	return dc
}

// Consume will send packet to input channel
func (dc *DefaultConsumer) Consume(packet packet.Packet) {
	dc.WaitGroup.Add(1)
	dc.IpChannel <- packet
}

// Wait will be used by client to wait for consumer to finish its job
func (dc *DefaultConsumer) Wait() {
	dc.WaitGroup.Wait()
	for i := 0; i < dc.ReqHandlers; i++ {
		dc.QuitChannel <- true
	}
}

// handlePackets will handle packets coming on IpChannel and quitChannel for exiting the function
func (dc *DefaultConsumer) handlePackets() {
	for {
		select {
		case packet := <-dc.IpChannel:
			err := dc.ConsumerFunc(packet)
			if err != nil {
				fmt.Println("Error occured while consuming: ", err.Error())
			}
			dc.WaitGroup.Done()
		case <-dc.QuitChannel:
			return
		}
	}
}
