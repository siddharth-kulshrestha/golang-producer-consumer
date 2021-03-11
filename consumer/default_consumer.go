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

func (dc *DefaultConsumer) Consume(packet packet.Packet) {
	dc.WaitGroup.Add(1)
	dc.IpChannel <- packet
}

func (dc *DefaultConsumer) Wait() {
	dc.WaitGroup.Wait()
	for i := 0; i < dc.ReqHandlers; i++ {
		dc.QuitChannel <- true
	}
}

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
