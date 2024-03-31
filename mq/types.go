package mq

type MQAdapter interface {
	Channel(channel string) (chan string, chan string)
}

type Message interface {
	Marshal() (string, error)
	UnMarshal(content string) (Message, error)
}
