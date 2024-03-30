package main

import (
	"encoding/json"
	"fmt"
	"github.com/masnun/goq/mq"
	"github.com/masnun/goq/mq/adapters/memory"
	"github.com/masnun/goq/worker"
	"time"
)

type User struct {
	Name  string
	Email string
	ID    int
}

func (u *User) Marshal() (string, error) {
	j, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (u *User) UnMarshal(content string) (mq.Message, error) {
	err := json.Unmarshal([]byte(content), u)
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) Type() string {
	return "user"
}

func PrintUser(info worker.WorkerInfo, message mq.Message) error {
	user := message.(*User)
	fmt.Printf("[WorkerID %d] Name: %s Email: %s Serial: %d \n", info.ID, user.Name, user.Email, user.ID)
	return nil
}

func main() {
	mq := memory.New()

	user := &User{
		Name:  "Masnun",
		Email: "masnun@gmail.com",
		ID:    0,
	}

	w := worker.New(mq, user, PrintUser, 2)
	w.Start()

	for i := 1; i < 101; i++ {
		w.Submit(&User{
			Name:  user.Name,
			Email: user.Email,
			ID:    i,
		})
	}

	w.Submit(user)

	<-time.After(3 * time.Minute)

}
