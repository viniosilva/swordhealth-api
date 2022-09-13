package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string
	Port string
}

type MySQLConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type Config struct {
	Server ServerConfig
	MySQL  MySQLConfig
}

func LoadConfig() Config {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	var configuration Config
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	configuration.MySQL.Password = os.Getenv("MYSQL_PASSWORD")

	return configuration
}
