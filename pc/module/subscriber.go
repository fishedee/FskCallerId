package module

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/language"
	"time"
	"yinghao/push/sdk"
)

type Subscriber struct {
	log        Log
	subscriber *sdk.ForeverSubscriber
	listener   SubscriberContactUpdateListener
}

type SubscriberConfig struct {
	Addr        string `config:"addr"`
	MessageSize int    `config:"messagesize"`
}

type SubscriberContactUpdateListener func(data string)

func NewSubscriber(log Log, config SubscriberConfig) (*Subscriber, error) {
	subscriber, err := sdk.NewForeverSubscriber(log, config.Addr, config.MessageSize, time.Second*10)
	if err != nil {
		return nil, err
	}
	return &Subscriber{
		log:        log,
		subscriber: subscriber,
	}, nil
}

func (this *Subscriber) SetContactUpdateListener(listener SubscriberContactUpdateListener) {
	this.listener = listener
}

func (this *Subscriber) Run() error {
	this.subscriber.On("/contact/update", func(data string) {
		defer CatchCrash(func(e Exception) {
			this.log.Critical("subscriber crash %v", e.Error())
		})
		if this.listener != nil {
			this.listener(data)
		}
	})
	return this.subscriber.Run()
}

func (this *Subscriber) Close() {
	this.subscriber.Close()
}
