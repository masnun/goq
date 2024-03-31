package main

import (
	"github.com/masnun/goq/cmd/internal"
	"github.com/masnun/goq/mq"
	"github.com/masnun/goq/mq/adapters/redis"
	"time"
)

func main() {
	redisAdapter := redis.New("redis://127.0.0.1:6379/0")
	userQueue := mq.New(redisAdapter, "user")
	postQueue := mq.New(redisAdapter, "post")

	for i := 1; i < 100; i++ {

		time.Sleep(2 * time.Second)

		userQueue.Publish(internal.User{
			Name:  "Masnun",
			Email: "masnun@gmail.com",
			ID:    i,
		})

		postQueue.Publish(&internal.Post{
			Title:   "Example Title",
			Summary: "Example Summary",
			ID:      i,
		})

	}

	<-time.After(3 * time.Minute)

}
