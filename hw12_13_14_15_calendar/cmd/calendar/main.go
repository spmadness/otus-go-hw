package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config_calendar.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	log := logger.New(config.Logger.Level)

	storage := app.NewStorage(config.Storage.Mode, config.StorageConnectionString())

	err := storage.Open()
	if err != nil {
		log.Error("failed to open storage connection: " + err.Error())
		return
	}

	calendar := app.New(log, storage)

	serverHTTP := internalhttp.NewServer(log, calendar, config.HTTPServerAddress())
	serverGRPC := internalgrpc.NewServer(log, storage, config.Server.GRPC.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err := storage.Close()
		if err != nil {
			log.Error("failed to close storage: " + err.Error())
		}

		if err := serverHTTP.Stop(ctx); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}

		serverGRPC.Stop()
	}()

	go func() {
		serverGRPC.Start()
	}()

	log.Info("calendar is running...")

	if err := serverHTTP.Start(ctx); err != nil {
		log.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
