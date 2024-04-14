package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type PushNotification struct {
	ID      string
	Message string
}

var successCount int
var mutex sync.Mutex

func generateNotifications(out chan<- PushNotification) {
	for {
		time.Sleep(time.Millisecond * 900)
		notification := PushNotification{
			ID:      fmt.Sprintf("%d", rand.Int63()),
			Message: fmt.Sprintf("Message %d", rand.Int()),
		}
		out <- notification
	}
}

func sendNotifications(in <-chan PushNotification, done chan<- bool) {
	for range in {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(16)+5))
		mutex.Lock()
		successCount++
		mutex.Unlock()
	}
	done <- true
}

func main() {
	notificationChannel := make(chan PushNotification, 100)
	doneChannel := make(chan bool)

	go generateNotifications(notificationChannel)
	go sendNotifications(notificationChannel, doneChannel)

	go func() {
		for {
			time.Sleep(time.Second * 10)
			mutex.Lock()
			fmt.Printf("Successful sends: %d\n", successCount)
			successCount = 0
			mutex.Unlock()
		}
	}()

	<-doneChannel
}
