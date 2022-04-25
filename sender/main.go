package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	// Define RabbitMQ server URL
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new rabbitMQ connection
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	/*
		Let's start by opening a channel to our RabbitMQ
		instance over the connection we have already
		established
	*/
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to
	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	// Create a new fiber instance
	app := fiber.New()

	// Add middleware
	app.Use(
		logger.New(), // add simple logger
	)

	// Add route
	app.Get("/send", func(ctx *fiber.Ctx) error {
		// create a message to publish
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(ctx.Query("msg")),
		}

		// Attempt to publish a message to the queue
		if err := channelRabbitMQ.Publish(
			"",
			"QueueService1",
			false,
			false,
			message,
		); err != nil {
			return err
		}
		return nil
	})

	// Start Fiber API server
	log.Fatal(app.Listen(":3000"))
}
