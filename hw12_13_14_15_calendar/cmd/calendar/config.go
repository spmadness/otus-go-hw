package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
	Mode     string
}

type ServerConf struct {
	Host string
	Port int
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

func (c Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s",
		c.Storage.Host, c.Storage.Port, c.Storage.User, c.Storage.Password)
}
