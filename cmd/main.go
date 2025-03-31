package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"notification/internal/consumers"
	"notification/internal/providers"
	"notification/pkg/configs"
	"notification/pkg/logger"
	"notification/pkg/rabbitmq"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	var (
		debug   bool
		workers int
	)

	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "number of workers")
	flag.Parse()

	if debug {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	logger.Init(workers)
	providers.Register()
	configs.LoadConfig()
	rabbitmq.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context.WithValue(ctx, "workers", workers)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	go consumers.Consumer(ctx)

	sig := <-quit
	logger.Warning("Received signal: %s", sig)
	rabbitmq.CloseConnection()
	logger.Close()
}

// {"data":{"code":"1234","template":"falogin","to":"+989125573688"},"provider":"kavenegar"}
// {"data":{"message":"your code is 1234","to":"+16315551111"},"provider":"ycloud"}
