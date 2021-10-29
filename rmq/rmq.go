package rmq

import (
	"net/url"
	"os"

	daas "github.com/touno-io/goasa"

	"github.com/streadway/amqp"
)

func Connect(host string) (*amqp.Connection, *amqp.Channel) {
	if host == "" {
		host = "amqp://guest:guest@localhost:5672/"
	}

	u, err := url.Parse(host)
	if err != nil {
		daas.Fatal("URLParse:", err)
	}

	daas.Infof("RMQ::%s Connected... '%s' ", os.Getenv("RMQ_MODEL_BU"), u.Host)
	conn, err := amqp.Dial(host)
	if err != nil {
		daas.Fatal("Dial:", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		daas.Fatal("Channel:", err)
	}
	return conn, channel
}

func Consume(channel *amqp.Channel, channelName string, queueName string, routingKey string) <-chan amqp.Delivery {
	daas.Infof("RMQ::%s Consume '%s' (%s).", os.Getenv("RMQ_MODEL_BU"), channelName, queueName)

	if os.Getenv("SENTRY_ENV") == "manual" {
		Queue(channel, channelName, queueName, routingKey)
	}
	message, err := channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		daas.Fatal("Consume::", err)
	}

	return message
}

func Queue(channel *amqp.Channel, channelName string, queueName string, routingKey string) {

	if os.Getenv("SENTRY_ENV") == "manual" {
		err := channel.ExchangeDeclare(channelName, "direct", true, false, false, false, nil)
		if err != nil {
			daas.Fatal("ExchangeDeclare::", err)
		}
		_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			daas.Fatal("QueueDeclare::", err)
		}
	}

	err := channel.QueueBind(queueName, routingKey, channelName, false, nil)
	if err != nil {
		daas.Fatal("QueueBind::", err)
	}
}

func Publish(channel *amqp.Channel, channelName string, routingKey string, body []byte) {
	daas.Infof("RMQ::%s Publish '%s:%s' (%d bytes).", os.Getenv("RMQ_MODEL_BU"), channelName, routingKey, len(body))

	err := channel.Publish(channelName, routingKey, false, false, amqp.Publishing{ContentType: "application/json", Body: body})
	if err != nil {
		daas.Fatal("Publish::", err)
	}
}
