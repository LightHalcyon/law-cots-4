package main

import (
	"net/http"
	"log"
	// "strings"
	// "os"

	"github.com/gin-gonic/gin"
	"github.com/reznov53/law-cots2/mq"
	"github.com/gin-contrib/cors"
)

// ErrorJSON error struct to be used when error occured
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

func sendMessage(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin","*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-Routing-Key, Host")

	routingKey := c.GetHeader("X-Routing-Key")

	urlOut, found := c.GetPostForm("url")
	if !found {
		c.JSON(http.StatusBadRequest, appError{
			Code:		http.StatusBadRequest,
			Message:	"URL not found, did you guys specify the URL?",
		})
		return
	}
	err := ch.PostMessage(urlOut, routingKey)
	if err != nil {
		c.JSON(http.StatusNotFound, appError{
			Code:		http.StatusNotFound,
			Message:	"RabbitMQ Server Down/Not Found",
		})
		return
	}
	c.JSON(http.StatusOK, appError{
		Code:		http.StatusOK,
		Message:	"URL Received",
	})
	return
}

func main() {
	// url := "amqp://" + os.Getenv("UNAME") + ":" + os.Getenv("PW") + "@" + os.Getenv("URL") + ":" + os.Getenv("PORT") + "/"
	url = "amqp://1406568753:167664@152.118.148.103:5672/"
	// vhost := os.Getenv("VHOST")
	vhost = "1406568753"
	// exchangeName := os.Getenv("EXCNAME")
	exchangeName = "1406568753-front"
	exchangeType = "direct"

	ch, err = mq.InitMQ(url, vhost, exchangeName, exchangeType)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	
	r.Static("/asset", "./asset")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Upload Page",
		})
	})
	r.POST("/", sendMessage)

	conf := cors.DefaultConfig()
	conf.AllowOrigins = []string{"*"}
	conf.AddAllowHeaders("X-ROUTING-KEY")
	conf.AddAllowHeaders("Content-Type")
	conf.AddAllowHeaders("Access-Control-Allow-Origin")
	conf.AddAllowHeaders("Access-Control-Allow-Headers")
	conf.AddAllowHeaders("Access-Control-Allow-Methods")
	conf.AddAllowHeaders("Host")
	r.Use(cors.New(conf))

	r.Run("0.0.0.0:21005")
	ch.Conn.Close()
	ch.Ch.Close()
}