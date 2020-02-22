package util

import (
	"fskcallerid/module"
	"github.com/fishedee/app/config"
	"github.com/fishedee/app/ioc"
	"github.com/fishedee/app/log"
	"github.com/fishedee/app/queue"
)

func MustRegisterIoc(obj interface{}) {
	ioc.MustRegisterIoc(obj)
}

func MustInvokeIoc(invoker interface{}) {
	ioc.MustInvokeIoc(invoker)
}

type Log = log.Log

type Config = config.Config

type Queue = queue.Queue

type Serial = module.Serial

type Contact = module.Contact

type Subscriber = module.Subscriber

type ConfigFile string

func NewConfig(configFile ConfigFile) config.Config {
	config, err := config.NewConfig("ini", string(configFile))
	if err != nil {
		panic(err)
	}
	return config
}

func NewLog(config config.Config) log.Log {
	logConfig := log.LogConfig{}
	config.MustBind("log", &logConfig)
	log, err := log.NewLog(logConfig)
	if err != nil {
		panic(err)
	}
	return log
}

func NewQueue(config config.Config, log log.Log) queue.Queue {
	queueConfig := queue.QueueConfig{}
	config.MustBind("queue", &queueConfig)
	queue, err := queue.NewQueue(log, queueConfig)
	if err != nil {
		panic(err)
	}
	return queue
}

func NewSerial(config config.Config, log log.Log) *module.Serial {
	serialConfig := module.SerialConfig{}
	config.MustBind("serial", &serialConfig)
	serial, err := module.NewSerial(log, serialConfig)
	if err != nil {
		panic(err)
	}
	return serial
}

func NewSubscriber(config config.Config, log log.Log) *module.Subscriber {
	subscriberConfig := module.SubscriberConfig{}
	config.MustBind("subscriber", &subscriberConfig)
	subscriber, err := module.NewSubscriber(log, subscriberConfig)
	if err != nil {
		panic(err)
	}
	return subscriber
}

func NewContact(config config.Config, log log.Log, queue queue.Queue) *module.Contact {
	contactConfig := module.ContactConfig{}
	config.MustBind("contact", &contactConfig)
	contact, err := module.NewContact(log, queue, contactConfig)
	if err != nil {
		panic(err)
	}
	return contact
}

func init() {
	ioc.MustRegisterIoc(NewConfig)
	ioc.MustRegisterIoc(NewLog)
	ioc.MustRegisterIoc(NewQueue)
	ioc.MustRegisterIoc(NewSerial)
	ioc.MustRegisterIoc(NewSubscriber)
	ioc.MustRegisterIoc(NewContact)
}
