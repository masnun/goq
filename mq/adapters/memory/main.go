package memory

func New() *InMemoryMQ {
	return &InMemoryMQ{Channels: make(map[string]chan string)}
}

type InMemoryMQ struct {
	Channels map[string]chan string
}

func (mq *InMemoryMQ) Channel(channel string) chan string {
	targetChannel, ok := mq.Channels[channel]
	if !ok {
		targetChannel = make(chan string, 10)
		mq.Channels[channel] = targetChannel
	}

	return targetChannel
}
