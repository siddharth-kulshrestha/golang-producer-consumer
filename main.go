package main

import (
	"fmt"
	"strings"

	"./consumer"
	"./packet"
	"./producer"
)

//func mockDataProducer(args ...interface{}) ([]packet.Packet, bool, []interface{}) {
//fmt.Println(args...)
//var packetsLength int64
//packetsLength = 100
//if len(args) > 0 {
//fmt.Println("Args 0 :")
//fmt.Println(args[0])
//_, ok2 := args[0].(int)
//fmt.Println(ok2)

//if val, ok := args[0].(int64); ok {
//fmt.Println("Ok")
//packetsLength = val
//}
//}
//fmt.Println("packetsLength: ", packetsLength)
//prepender := "mock001"
//if len(args) > 1 {
//if prepe, ok := args[1].(string); ok {
//prepender = prepe
//}
//}

//op := []packet.Packet{}
//var i int64
//for i = 0; i < packetsLength; i++ {
//op = append(op, packet.Packet{
//Data: fmt.Sprintf("%s: %d", prepender, i),
//Err:  nil,
//})
//}
//if packetsLength > 1000 {
//return op, false, nil
//}
//return op, true, []interface{}{packetsLength + 100, fmt.Sprintf("%s_%d", prepender, packetsLength+100)}
//}

//func receiver(opChannel chan packet.Packet, wg *sync.WaitGroup) {
//fmt.Println("Starting receiver")
//for packet := range opChannel {
////packet := <-opChannel
//fmt.Print("Received one packet: ")
//fmt.Print(packet.Data)
//}
//wg.Done()
//}

//func consumerFunc(pk packet.Packet) error {
//fmt.Println("Received one packet: ")
//fmt.Println(pk.Data)
//return nil
//}

func tweetProducerFunc(args ...interface{}) ([]packet.Packet, bool, []interface{}) {
	mockData := []string{
		"davecheney tweets about golang",
		"beertocode does not tweet about golang",
		"ironzeb tweets about golang",
		"beertocode tweets about golang",
		"vampirewalk666 tweets about golang",
	}
	op := make([]packet.Packet, len(mockData))
	for _, tweet := range mockData {
		op = append(op, packet.New(tweet, nil))
	}
	fmt.Println(mockData)
	return op, false, nil
}

func tweetConsumerFunc(pk packet.Packet) error {
	if pk.Data == nil {
		return nil
	}
	if _, ok := pk.Data.(string); !ok {
		return fmt.Errorf("Could not parse this tweet into string: %v", pk.Data)
	}
	tweet, _ := pk.Data.(string)

	// Parsing tweet whether to check if it contains "tweets about golang" or not.
	suffix := "tweets about golang"

	if strings.HasSuffix(tweet, suffix) {
		tweetLiterals := strings.Split(tweet, " ")
		if len(tweetLiterals) > 0 {
			fmt.Println(tweetLiterals[0])
		}
	}
	return nil
}

func main() {
	fmt.Println("Example of producer and consumer")

	/*	var wg sync.WaitGroup
		prdcr := producer.New(true, &wg)

		var opChannel chan producer.Packet
		var quitChan chan bool
		//opChannel = make(chan producer.Packet, 1000)
		//quitChan = make(chan bool, 1)
		wg.Add(2)
		go receiver(opChannel, &wg)
		prdcr.Produce(opChannel, mockDataProducer, []interface{}{int64(10), "mockBySid"}, quitChan)
		wg.Wait()
		fmt.Println("Released wait group")
		//producer := producer.New(source, msgQueue)
		//consumer := consumer.New(producer, msgQueue, consumerOp)

		//consumer.consume()
	*/
	cons := consumer.New(1, 1000, tweetConsumerFunc)

	prdcr := producer.New(cons, true, tweetProducerFunc)

	prdcr.Produce(nil)

	prdcr.Wait()

	cons.Wait()

}
