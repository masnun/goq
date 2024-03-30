package mq

type MQAdapter interface {
	Channel(channel string) chan string
}

type Message interface {
	Type() string
	Marshal() (string, error)
	UnMarshal(content string) (Message, error)
}
