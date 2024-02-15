package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Broker MessageBrokerConf
	Logger LoggerConf
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

	return config
}

func (c Config) BrokerConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		c.Broker.User, c.Broker.Password, c.Broker.Host, c.Broker.Port)
}
