package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/broker"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config_scheduler.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := NewConfig(configFile)
	log := logger.New(config.Logger.Level)

	storage := sqlstorage.New(config.StorageConnectionString())

	err := storage.Open()
	if err != nil {
		log.Error("failed to open storage connection: " + err.Error())
		return
	}

	b := broker.New(config.BrokerConnectionString())
	err = b.Open()
	if err != nil {
		log.Error("failed to open broker connection: " + err.Error())
		return
	}

	err = b.SetQueue(config.Broker.Queue)
	if err != nil {
		log.Error("failed to set queue: " + err.Error())
		return
	}

	scheduler := app.NewScheduler(storage, b, log)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go scheduler.ProcessNotifications(ctx, config.Storage.PollTimeSeconds, config.Storage.OutdatedEventDays)

	log.Info("scheduler is running...")

	<-ctx.Done()

	err = storage.Close()
	if err != nil {
		log.Error("failed to close storage: " + err.Error())
	}

	err = b.Close()
	if err != nil {
		log.Error("failed to close broker: " + err.Error())
	}

	log.Info("shutting down scheduler...")
}
