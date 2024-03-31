package worker

import (
	"fmt"
	"github.com/masnun/goq/mq"
	"reflect"
)

type TaskFunction func(info WorkerInfo, message mq.Message) error

type Worker struct {
	Queue                *mq.Queue
	Concurrency          int
	TaskFunction         TaskFunction
	messageType          reflect.Type
	messageTypeIsPointer bool
}

type WorkerInfo struct {
	ID        int
	QueueName string
}

func getType(message mq.Message) (reflect.Type, bool) {
	t := reflect.TypeOf(message)
	v := reflect.ValueOf(message)
	if v.Kind() == reflect.Ptr {
		t = v.Elem().Type()
		return t, true
	}

	return t, false

}

func New(queue *mq.Queue, message mq.Message, taskFunction TaskFunction, concurrency int) *Worker {

	t, p := getType(message)

	return &Worker{
		Concurrency:          concurrency,
		TaskFunction:         taskFunction,
		Queue:                queue,
		messageType:          t,
		messageTypeIsPointer: p,
	}
}

func (w *Worker) Start() {
	for i := 0; i < w.Concurrency; i++ {
		info := WorkerInfo{ID: i, QueueName: w.Queue.Name()}
		go w.ProcessTask(info)
	}
}

func (w *Worker) ProcessTask(info WorkerInfo) error {

	channel := w.Queue.Channel()

	for {
		rawMsg, ok := <-channel
		if !ok {
			break
		}

		msgType := reflect.New(w.messageType)
		if !w.messageTypeIsPointer {
			msgType = msgType.Elem()
		}

		message := msgType.Interface().(mq.Message)
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
