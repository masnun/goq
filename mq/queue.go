package mq

type Queue struct {
	sendChannel    chan<- string
	recieveChannel <-chan string
}

func New(adapter MQAdapter, channelName string) *Queue {
	send, recieve := adapter.Channel(channelName)
	return &Queue{
		sendChannel:    send,
		recieveChannel: recieve,
	}
}

func (q *Queue) Channel() <-chan string {
	return q.recieveChannel
}

func (q *Queue) Publish(message Message) error {
	str, err := message.Marshal()
	if err != nil {
		return err
	}

	q.sendChannel <- str
	return nil
}
