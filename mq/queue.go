package mq

type Queue struct {
	sendChannel    chan<- string
	recieveChannel <-chan string
	channelName    string
}

func New(adapter MQAdapter, channelName string) *Queue {
	send, recieve := adapter.Channel(channelName)
	return &Queue{
		sendChannel:    send,
		recieveChannel: recieve,
		channelName:    channelName,
	}
}
func (q *Queue) String() string {
	return q.channelName
}

func (q *Queue) Name() string {
	return q.channelName
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
