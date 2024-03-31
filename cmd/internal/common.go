package internal

import (
	"encoding/json"
	"github.com/masnun/goq/mq"
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

type Post struct {
	Title   string
	Summary string
	ID      int
}

func (p *Post) Marshal() (string, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (p *Post) UnMarshal(content string) (mq.Message, error) {

	newPost := &Post{}
	err := json.Unmarshal([]byte(content), newPost)
	if err != nil {
		return &Post{}, err
	}

	return newPost, nil
}
