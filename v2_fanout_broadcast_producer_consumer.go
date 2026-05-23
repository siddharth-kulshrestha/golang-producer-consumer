package main

import (
	"fmt"
	"sync"
)

/*
Problem:
Build a notification service where a single producer publishes notifications
(e.g., "Flash Sale Started"), and all subscribed consumers should receive the message.

There are 3 subscribers:
- Email Service
- SMS Service
- Push Notification Service

Every notification must be delivered to ALL consumers concurrently.
*/

type NotificationProducer struct {
	ch    chan<- string
	errCh chan error
}

func (np *NotificationProducer) Produce() []string {
	return []string{"Notification1", "Notification2", "Notification3", "Notification4"}
}

func (np *NotificationProducer) Send() {}

func NewNotificationProducer(ch chan<- string, errCh chan error) Producer[string] {
	return &NotificationProducer{
		ch:    ch,
		errCh: errCh,
	}
}

type NotificationWorker struct {
	id    int
	ch    <-chan string
	errCh chan error
}

func (nw *NotificationWorker) ExtractAndProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case msg, ok := <-nw.ch:
			if !ok {
				print(fmt.Sprintln(nw.id, ": ", "channel is closed"))
				return
			}
			print(fmt.Sprintln(nw.id, ": ", msg))
		case err := <-nw.errCh:
			if err != nil {
				print(fmt.Sprintln(nw.id, ": ", err))
				close(nw.errCh)
				return
			}
		}
	}
}

func (nw *NotificationWorker) Process(data string) (string, error) {}

func NewNotificationWorker(id int, ch <-chan string, err chan error) Worker[string] {
	return &NotificationWorker{
		id:    id,
		ch:    ch,
		errCh: err,
	}
}

func ExmapleNotificationBroadCast() {

}
