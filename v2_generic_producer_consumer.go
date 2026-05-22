/*
Problem: Write a program where a single producer distributes 10 tasks (integers 1–10) across a
worker pool of 3 concurrent worker goroutines.
Each worker should process the task (e.g., multiply by 2),
and a collector should aggregate all results into a single slice.
*/

package main

import (
	"fmt"
	"sync"
)

// Producer
type Producer[T any] interface {
	Produce() []T
	Send()
}

// Worker
type Worker[T any] interface {
	ExtractAndProcess(*sync.WaitGroup)
	Process(data T) (T, error)
}

type NumberWorker struct {
	ch    <-chan int
	errCh chan error
	id    int
}

func NewNumberWorker(id int, ch <-chan int, errChan chan error) Worker[int] {
	return &NumberWorker{
		ch:    ch,
		errCh: errChan,
		id:    id,
	}
}

func (nw *NumberWorker) Process(data int) (int, error) {
	return data * 2, nil

}

func (nw *NumberWorker) ExtractAndProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case v, ok := <-nw.ch:
			if !ok {
				fmt.Println("Channel is closed")
				return
			}
			ret, err := nw.Process(v)
			if err != nil {
				fmt.Println("error occured while processing: ")
				fmt.Println(err)
			}
			fmt.Println(ret)
		case err := <-nw.errCh:
			if err != nil {
				fmt.Println(err)
				close(nw.errCh)
				return
			}
		}
	}
}

type NumberProducer struct {
	ch chan<- int
}

func NewNumberProducer(ch chan<- int) Producer[int] {
	return &NumberProducer{
		ch: ch,
	}
}

func (np *NumberProducer) Produce() []int {
	sl := make([]int, 0, 10)
	for i := 1; i <= 10; i += 1 {
		sl = append(sl, i)
	}
	return sl
}

func (np *NumberProducer) Send() {
	// go func() {
	for _, v := range np.Produce() {
		np.ch <- v
	}
	close(np.ch)
	// }()
}

func ExampleNumberProducerConsumer() {
	ch := make(chan int, 10)
	errCh := make(chan error)
	p := NewNumberProducer(ch)
	p.Send()

	wg := sync.WaitGroup{}

	pool := make([]Worker[int], 3)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		pool[i] = NewNumberWorker(i, ch, errCh)
	}

	for _, w := range pool {
		go w.ExtractAndProcess(&wg)
	}

	wg.Wait()
}

func main() {
	ExampleNumberProducerConsumer()
}
