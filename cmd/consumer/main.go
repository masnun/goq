package main

import (
	"fmt"
	"github.com/masnun/goq/cmd/internal"
	"github.com/masnun/goq/mq"
	"github.com/masnun/goq/mq/adapters/redis"
	"github.com/masnun/goq/worker"
	"time"
)

func PrintUser(info worker.WorkerInfo, message mq.Message) error {

	user := message.(internal.User)
	fmt.Printf("[WorkerID %d] [Queue: %s] Name: %s Email: %s Serial: %d \n", info.ID, info.QueueName, user.Name, user.Email, user.ID)
	return nil
}

func PrintPost(info worker.WorkerInfo, message mq.Message) error {
	post := message.(*internal.Post)
	fmt.Printf("[WorkerID %d] [Queue: %s] Title: %s Summary: %s Serial: %d \n", info.ID, info.QueueName, post.Title, post.Summary, post.ID)
	return nil
}

func main() {

	go ConsumeUserQueue()
	go ConsumePostQueue()

	<-time.After(30 * time.Minute)

}

func ConsumeUserQueue() {
	redisAdapter := redis.New("redis://127.0.0.1:6379/0")
	queue := mq.New(redisAdapter, "user")

	user := internal.User{}

	w := worker.New(queue, user, PrintUser, 3)
	w.Start()

}

func ConsumePostQueue() {
	redisAdapter := redis.New("redis://127.0.0.1:6379/0")
	queue := mq.New(redisAdapter, "post")

	post := &internal.Post{}

	w := worker.New(queue, post, PrintPost, 1)
	w.Start()

}
