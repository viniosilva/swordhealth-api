package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type MySQLConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

type NotificationService struct {
	SummaryLength int `mapstructure:"summary_length"`
}

type Service struct {
	Notification NotificationService `mapstructure:"notification"`
}

type Config struct {
	Server  ServerConfig `mapstructure:"server"`
	MySQL   MySQLConfig  `mapstructure:"mysql"`
	Service Service      `mapstructure:"service"`
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
