/*
Problem: Write a program where a single producer distributes 10 tasks (integers 1–10) across a
worker pool of 3 concurrent worker goroutines.
Each worker should process the task (e.g., multiply by 2),
and a collector should aggregate all results into a single slice.
*/

package main

import "fmt"

// Producer
type Producer[T any] interface {
	Produce() []T
	Register()
}

// Worker
type Worker[T any] interface {
	ExtractAndProcess()
	Process(data T) (T, error)
}

type NumberWorker struct {
	ch    <-chan int
	errCh chan error
}

func NewNumberWorker(ch <-chan int, errChan chan error) Worker[int] {
	return &NumberWorker{
		ch:    ch,
		errCh: errChan,
	}
}

func (nw *NumberWorker) Process(data int) (int, error) {
	return data * 2, nil

}

func (nw *NumberWorker) ExtractAndProcess() {
	for {
		select {
		case v := <-nw.ch:
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
	sl := make([]int, 10)
	for i := 1; i <= 10; i += 1 {
		sl = append(sl, i)
	}
	return sl
}

func (np *NumberProducer) Register() {
	go func() {
		for _, v := range np.Produce() {
			np.ch <- v
		}
	}()
}
