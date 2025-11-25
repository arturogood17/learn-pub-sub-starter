package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	if err = pubsub.PublishJSON(ch, string(routing.ExchangePerilDirect),
		string(routing.PauseKey), routing.PlayingState{IsPaused: true}); err != nil {
		log.Fatal("Error al intentar enviar el mensaje")
	}

	defer conn.Close()
	fmt.Println("Conexión establecida con RabbitMQ exitosamente")
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	_, ok := <-sigs
	if ok {
		fmt.Println()
		fmt.Println("Cerrando el programa")
	}
}
