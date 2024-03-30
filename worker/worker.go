package worker

import (
	"fmt"
	"github.com/masnun/goq/mq"
	"reflect"
)

type TaskFunction func(info WorkerInfo, message mq.Message) error

type Worker struct {
	Concurrency int

	TaskFunction TaskFunction

	messageType    reflect.Type
	sendChannel    chan<- string
	recieveChannel <-chan string
}

type WorkerInfo struct {
	ID int
}

func getType(message mq.Message) reflect.Type {
	t := reflect.TypeOf(message)
	v := reflect.ValueOf(message)
	if v.Kind() == reflect.Ptr {
		t = v.Elem().Type()
	}

	return t

}

func New(adapter mq.MQAdapter, message mq.Message, taskFunction TaskFunction, concurrency int) *Worker {
	send, recieve := adapter.Channel(message.Type())
	return &Worker{
		Concurrency:    concurrency,
		TaskFunction:   taskFunction,
		messageType:    getType(message),
		sendChannel:    send,
		recieveChannel: recieve,
	}
}

func (w *Worker) Start() {
	for i := 0; i < w.Concurrency; i++ {
		info := WorkerInfo{ID: i}
		go w.ProcessTask(info)
	}
}

func (w *Worker) Submit(message mq.Message) error {
	str, err := message.Marshal()
	if err != nil {
		return err
	}

	w.sendChannel <- str
	return nil
}

func (w *Worker) ProcessTask(info WorkerInfo) error {

	for {
		rawMsg, ok := <-w.recieveChannel
		if !ok {
			break
		}

		message := reflect.New(w.messageType).Interface().(mq.Message)

		message, err := message.UnMarshal(rawMsg)
		if err != nil {
			return err
		}

		err = w.TaskFunction(info, message)

		if err != nil {
			fmt.Println("Task encountered error: ", err)
		}

	}

	return nil

}
