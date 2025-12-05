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

	defer conn.Close()

	pbChannel, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("Error al establecer la conexión")
	}

	queueName := fmt.Sprintf("%v.%v", routing.PauseKey, username)

	gameState := gamelogic.NewGameState(username)

	if err := pubsub.SubscribeJSON(conn, string(routing.ExchangePerilDirect),
		queueName, string(routing.PauseKey), pubsub.Transient, HandlerPause(gameState)); err != nil {
		log.Fatalln(err)
	}

	queueTopic := fmt.Sprintf("%s.%v", routing.ArmyMovesPrefix, username)
	keyTopic := fmt.Sprintf("%s.*", routing.ArmyMovesPrefix)

	if err := pubsub.SubscribeJSON(conn, string(routing.ExchangePerilTopic),
		queueTopic, keyTopic, pubsub.Transient,
		HandlerMove(gameState)); err != nil {
		log.Fatalln(err)
	}

	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}
		switch userInput[0] {
		case "spawn":
			if err := gameState.CommandSpawn(userInput); err != nil {
				log.Println(err)
			}
		case "move":
			armyMove, err := gameState.CommandMove(userInput)
			if err != nil {
				log.Println(err)
			}
			if err := pubsub.PublishJSON(pbChannel, string(routing.ExchangePerilTopic), queueTopic, armyMove); err != nil {
				log.Println(err)
			}
			log.Println(armyMove)
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Invalid command")
		}
	}
}
