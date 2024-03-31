package main

import (
	"encoding/json"
	"fmt"
	"github.com/masnun/goq/mq"
	"github.com/masnun/goq/mq/adapters/redis"
	"github.com/masnun/goq/worker"
	"time"
)

type User struct {
	Name  string
	Email string
	ID    int
}

func (u User) Marshal() (string, error) {
	j, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (u User) UnMarshal(content string) (mq.Message, error) {

	newUser := User{}
	err := json.Unmarshal([]byte(content), &newUser)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func PrintUser(info worker.WorkerInfo, message mq.Message) error {
	user := message.(User)
	fmt.Printf("[WorkerID %d] Name: %s Email: %s Serial: %d \n", info.ID, user.Name, user.Email, user.ID)
	return nil
}

func main() {
	redisAdapter := redis.New("redis://127.0.0.1:6379/0")
	queue := mq.New(redisAdapter, "test")

	user := User{
		Name:  "Masnun",
		Email: "masnun@gmail.com",
		ID:    0,
	}

	for i := 1; i < 100; i++ {

		time.Sleep(2 * time.Second)

		queue.Publish(User{
			Name:  user.Name,
			Email: user.Email,
			ID:    i,
		})
	}

	<-time.After(3 * time.Minute)

}
