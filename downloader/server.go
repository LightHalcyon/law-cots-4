package main

import (
	"log"
	"strings"
	"fmt"
	"strconv"

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

func initCh(url string, vhost string, exchangeName string, exchangeType string, queueName string) (*mq.Channel, error) {
	ch, err = mq.InitMQ(url, vhost)
	if err != nil {
		return ch, err
	}

	err = ch.ExcDeclare(exchangeName, exchangeType)
	if err != nil {
		return ch, err
	}

	return ch, nil
}

func main() {
	// url := "amqp://" + os.Getenv("UNAME") + ":" + os.Getenv("PW") + "@" + os.Getenv("URL") + ":" + os.Getenv("PORT") + "/"
	url = "amqp://1406568753:167664@152.118.148.103:5672/"
	// vhost := os.Getenv("VHOST")
	vhost = "1406568753"
	// exchangeName := os.Getenv("EXCNAME")
	exchangeName = "1406568753-front"
	exchangeType = "direct"
	exchangeName1 := "1406568753-dl1"
	exchangeName2 := "1406568753-compress"
	exchangeType1 := "fanout"

	ch, err := initCh(url, vhost, exchangeName, exchangeType, "urlpass")
	if err != nil {
		panic(err)
	}

	ch1, err := initCh(url, vhost, exchangeName1, exchangeType, "dlstatus")
	if err != nil {
		panic(err)
	}

	ch2, err := initCh(url, vhost, exchangeName2, exchangeType1, "compresspass")
	if err != nil {
		panic(err)
	}

	err = ch.QueueDeclare("urlpass")
	if err != nil {
		panic(err)
	}
	
	for i := 0; i < 10; i++  {
		err = ch1.QueueDeclare(joint("dlstatus", fmt.Sprint(strconv.Itoa(i))))
		if err != nil {
			panic(err)
		}
	}

	err = ch2.QueueDeclare("compresspass")
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

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			in := strings.Split(string(d.Body), " ")
			arr := strings.Split(in[0], ";")
			log.Println(string(d.Body))
			dl(arr, ch1, ch2, in[1])
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}