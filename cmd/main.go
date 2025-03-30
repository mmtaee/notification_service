package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"notification/internal/consumers"
	"notification/pkg/configs"
	"notification/pkg/rabbitmq"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	queueName := "notification"
	workers := runtime.NumCPU()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context.WithValue(ctx, "queueName", queueName)
	ctx = context.WithValue(ctx, "workers", workers)

	configs.LoadConfig()
	rbtmq := rabbitmq.NewRabbitMQ()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go consumers.Consumer(ctx, rbtmq)

	sig := <-quit
	log.Println("Got signal:", sig)
	rbtmq.CloseConnection()
	log.Println("shutting down")
}

// {"message":"aaaa","provider":{"extra":{"template":"falogin"},"name":"kavenegar"},"to":"+989125573688"}
// {"message":"aaaa","provider":{"extra":{"template":"falogin"},"name":"ycloud"},"to":"+16315551111"}
