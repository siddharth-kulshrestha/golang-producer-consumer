package main

import (
	"./producer"
	"fmt"
)

func main() {
	fmt.Println("Example of producer and consumer")

	producer.New()
	producer := producer.New(source, msgQueue)
	consumer := consumer.New(producer, msgQueue, consumerOp)

	consumer.consume()
}
