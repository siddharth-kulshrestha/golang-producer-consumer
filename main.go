package main

import (
	"fmt"
	"strings"

	"./consumer"
	"./packet"
	"./producer"
)

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

	cons := consumer.New(1, 1000, tweetConsumerFunc) // first arg -> number of parallel consumers that we want to run.
	// second arg -> Max packet size for buffered channel
	// third arg -> consumer func

	prdcr := producer.New(cons, true, tweetProducerFunc) // first arg -> consumer
	// second arg -> should continue to work even if there is an error in production of any packet

	prdcr.Produce(nil) // Here we are producing the packet

	prdcr.Wait() // Waiting for all packets to get produced

	cons.Wait() // Waiting for all consumers to consume the data.

}
