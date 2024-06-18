package config

import "time"

var Config Configuration

type Server struct {
	Host        string
	Port        string
	Environment string
}

type MySQLConfig struct {
	Host string
	Port string
	User string
	Pass string
}

type RedisConfig struct {
	Host string
	Port string
	User string
	Pass string
}

type Webhook struct {
	URL     string
	Method  string
	Timeout time.Duration
}

type Configuration struct {
	Server  Server
	MySQL   MySQLConfig
	Redis   RedisConfig
	Webhook Webhook
}
