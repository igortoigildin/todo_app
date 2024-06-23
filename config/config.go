package config

import "os"

type Config struct {
	Port 	string
	DBpath 	string
	Pass 	string
}

func LoadConfig() *Config {
	cfg := Config{
		Port: 	"7540",
		DBpath:	"../scheduler.db",
	}
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		cfg.Port = envPort
	}
	if envDBpath := os.Getenv("TODO_DBFILE"); envDBpath != "" {
		cfg.DBpath = envDBpath
	}
	if envPass := os.Getenv("TODO_PASSWORD"); envPass != "" {
		cfg.Pass = envPass
	}
	return &cfg
}