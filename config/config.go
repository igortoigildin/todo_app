package config

import "os"

type Config struct {
	Port 	string
	DBpath 	string
}

func LoadConfig() *Config {
	cfg := Config{
		Port: 	"7540",
		DBpath:	"../scheduler.db",
	}
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		cfg.Port = envPort
	}
	if envDBpath := os.Getenv("TODO_PORT"); envDBpath != "" {
		cfg.DBpath = envDBpath
	}
	return &cfg
}