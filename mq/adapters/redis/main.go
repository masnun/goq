package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type ChannelSet struct {
	Send    chan string
	Recieve chan string
}

func New(dsn string) *RedisMQ {
	return &RedisMQ{
		Client:   createRedisClient(dsn),
		Channels: make(map[string]ChannelSet),
	}
}

type RedisMQ struct {
	Client   *redis.Client
	Channels map[string]ChannelSet
}

func (mq *RedisMQ) Channel(channel string) (chan string, chan string) {
	targetChannel, ok := mq.Channels[channel]
	if !ok {

		sendChan := make(chan string, 10)
		receieveChan := make(chan string, 10)

		targetChannel = ChannelSet{Send: sendChan, Recieve: receieveChan}
		mq.Channels[channel] = targetChannel
		ctx := context.Background()
		pubsub := mq.Client.Subscribe(ctx, channel)

		go readFromRedis(ctx, pubsub, targetChannel.Recieve)
		go writeToRedis(ctx, channel, targetChannel.Send, mq.Client)

	}

	return targetChannel.Send, targetChannel.Recieve
}

func createRedisClient(dsn string) *redis.Client {
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		panic(err)
	}

	//opts.TLSConfig = &tls.Config{
	//	InsecureSkipVerify: true,
	//}

	return redis.NewClient(opts)

}

func readFromRedis(ctx context.Context, pubsub *redis.PubSub, channel chan<- string) {

	for {
		msg, err := pubsub.ReceiveMessage(ctx) // TODO handle errors, perhaps unregister the channel

		if err != nil {
			fmt.Println("rfr err", err)
		} else {
			channel <- msg.Payload
		}

	}

}

func writeToRedis(ctx context.Context, channelName string, channel <-chan string, client *redis.Client) {

	for {
		msg, ok := <-channel
		if !ok {
			break
		}

		//fmt.Println("Writing to redis", msg)
		client.Publish(ctx, channelName, msg)

	}
}
