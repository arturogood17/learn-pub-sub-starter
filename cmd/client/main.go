package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatal("Error al establecer la conexión")
	}

	queueName := fmt.Sprintf("%v.%v", routing.PauseKey, username)

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient)
	if err != nil {
		log.Fatal("Error al establecer la conexión")
	}

	defer conn.Close()

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	_, ok := <-sigs
	if ok {
		fmt.Println()
		fmt.Println("Cerrando el programa")
	}

	gamelogic.PrintClientHelp()
}
