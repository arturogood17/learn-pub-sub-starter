package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(conStr)
	if err != nil {
		log.Fatal("Error al establecer la conexión")
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
