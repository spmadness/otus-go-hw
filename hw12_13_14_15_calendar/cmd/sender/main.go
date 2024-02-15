package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/broker"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config_sender.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := NewConfig(configFile)
	log := logger.New(config.Logger.Level)

	b := broker.New(config.BrokerConnectionString())
	err := b.Open()
	if err != nil {
		log.Error("failed to open broker connection: " + err.Error())
		return
	}

	log.Info("sender is running...")

	msgChan, err := b.ConsumeMessage(config.Broker.Queue)
	if err != nil {
		log.Error("Failed to start consuming messages: " + err.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		for msg := range msgChan {
			log.Info(fmt.Sprintf("%s Received message: %s",
				time.Now().Format("2006-01-02 15:04:05"), msg.Body))
		}
	}()

	<-ctx.Done()

	err = b.Close()
	if err != nil {
		log.Error("failed to close broker: " + err.Error())
	}

	log.Info("shutting down sender...")
}
