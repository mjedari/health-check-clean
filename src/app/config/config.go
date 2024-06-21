package config

import (
	"strings"
	"time"
)

var Config Configuration

type Server struct {
	Host        string
	Port        string
	Environment string
}

type MySQL struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

type Redis struct {
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

type Service struct {
	Cache string
}

type Configuration struct {
	Server  Server
	MySQL   MySQL
	Redis   Redis
	Webhook Webhook
	Service Service
}

func (c *Configuration) IsProduction() bool {
	return strings.ToLower(c.Server.Environment) == "production" ||
		strings.ToLower(c.Server.Environment) == "prod"
}
