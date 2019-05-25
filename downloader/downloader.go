package main

import (
	"log"

	"github.com/reznov53/law-cots2/mq"
	// "github.com/reznov53/law-cots2/download"
)

type appError struct {
	Code	int    `json:"status"`
	Message	string `json:"message"`
}

var ch *mq.Channel
var err error
// var files map[string]string
var url, vhost, exchangeName, exchangeType string

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// url := "amqp://" + os.Getenv("UNAME") + ":" + os.Getenv("PW") + "@" + os.Getenv("URL") + ":" + os.Getenv("PORT") + "/"
	url = "amqp://1406568753:167664@152.118.148.103:5672/"
	// vhost := os.Getenv("VHOST")
	vhost = "1406568753"
	// exchangeName := os.Getenv("EXCNAME")
	exchangeName = "1406568753-front"
	exchangeType = "direct"

	ch, err = mq.InitMQ(url, vhost)
	if err != nil {
		panic(err)
	}

	err = ch.ExcDeclare(exchangeName, exchangeType)
	if err != nil {
		panic(err)
	}

	err = ch.QueueDeclare("urlpass")
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Ch.Consume(
		"urlpass", // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}