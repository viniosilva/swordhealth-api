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
	Password string
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

type Crypto struct {
	HashKey   string
	JwtKey    string
	ExpiresIn int64 `mapstructure:"expires_in"`
}

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Crypto Crypto       `mapstructure:"crypto"`
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
	configuration.Crypto.HashKey = os.Getenv("CRYPTO_KEY")
	configuration.Crypto.JwtKey = os.Getenv("JWT_KEY")

	return configuration
}
