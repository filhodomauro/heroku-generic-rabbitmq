package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

const rmqCONSUMERVARIABLE = "RABBITMQ_BIGWIG_RX_URL"
const rmqPRODUCERVARIABLE = "RABBITMQ_BIGWIG_TX_URL"

func main() {

}

func ConfigureRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.GET("/publish", publish)
	return router
}

func publish(c *gin.Context) {

	var message PublishMessage
	err := c.ShouldBindJSON(&message)
	checkError(err, "Error to bind JSON")

	conn, err := amqp.Dial(os.Getenv(rmqPRODUCERVARIABLE))
	checkError(err, "Error to connect MQ")
	defer conn.Close()

	channel, err := conn.Channel()
	checkError(err, "Error to open a channel")
	channel.Close()

	queue, err := queueDeclare(channel, message.Queue)
	checkError(err, fmt.Sprintf("Erro to declare Queue: %v", message.Queue))

	err = publishOnMQ(channel, queue.Name, message)
	checkError(err, fmt.Sprintf("Error to publish message to queue: %v", queue.Name))
}

func checkError(err error, message string) {
	if err != nil {
		message := fmt.Sprintf("Error=> %v - %v", message, err)
		log.Fatal(message)
		panic(message)
	}
}

func queueDeclare(channel *amqp.Channel, name string) (amqp.Queue, error) {
	q, err := channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return q, err
}

func publishOnMQ(channel *amqp.Channel, queueName string, message PublishMessage) (err error) {
	channel.Publish(
		"", queueName, false, false, 
		amqp.Publishing{
			ContentType: "application/json",
			Body: 
		},
	)
	return err
}

type PublishMessage struct {
	Queue string `json:"queue"`
	CallbackUrl string `json:"callback_url"`
}
