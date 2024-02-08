package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Broker  MessageBrokerConf
	Storage StorageConf
	Logger  LoggerConf
}

type LoggerConf struct {
	Level string
}

type MessageBrokerConf struct {
	User     string
	Password string
	Host     string
	Port     int
	Queue    string
}

type StorageConf struct {
	User              string
	Password          string
	Host              string
	Port              int
	Name              string
	PollTimeSeconds   int `mapstructure:"poll_time_seconds"`
	OutdatedEventDays int `mapstructure:"outdated_event_days"`
}

func NewConfig(configFile string) Config {
	var config Config

	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config read error: %s", err)
		os.Exit(1)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("config unmarshal error: %s", err)
		os.Exit(1)
	}

	if config.Storage.PollTimeSeconds < 1 {
		fmt.Println("poll_time_seconds parameter must be > 0")
		os.Exit(1)
	}

	if config.Storage.OutdatedEventDays < 1 {
		fmt.Println("outdated_event_days parameter must be > 0")
		os.Exit(1)
	}

	return config
}

func (c Config) BrokerConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		c.Broker.User, c.Broker.Password, c.Broker.Host, c.Broker.Port)
}

func (c Config) StorageConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		c.Storage.Host, c.Storage.Port, c.Storage.User, c.Storage.Password, c.Storage.Name)
}
