package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(conStr)
	if err != nil {
		log.Fatal("Error al establecer la conexión")
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error al crear el canal")
	}

	defer conn.Close()
	fmt.Println("Conexión establecida con RabbitMQ exitosamente")

	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}
		command := userInput[0]

		switch command {
		case "pause":
			log.Println("Sending pause message")
			if err = pubsub.PublishJSON(ch, string(routing.ExchangePerilDirect),
				string(routing.PauseKey), routing.PlayingState{IsPaused: true}); err != nil {
				log.Fatal("Error al intentar enviar el mensaje")
			}
		case "resume":
			log.Println("Sending resume message")
			if err = pubsub.PublishJSON(ch, string(routing.ExchangePerilDirect),
				string(routing.PauseKey), routing.PlayingState{IsPaused: false}); err != nil {
				log.Fatal("Error al intentar enviar el mensaje")
			}
		case "quit":
			log.Println("Saliendo...")
			return
		default:
			log.Println("Comando inválido. Inténtalo de nuevo")
		}

	}
}
