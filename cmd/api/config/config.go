package config

import (
	"github.com/joeshaw/envdecode"
	"log"
	"time"
)

type limiter struct {
	enabled bool
	rps     float64
	burst   int
}

type smtp struct {
	host     string
	port     int
	username string
	password string
	sender   string
}

type Conf struct {
	Debug   bool `env:"DEBUG,required"`
	Server  serverConf
	Db      dbConf
	Limiter limiter
	Smtp    smtp
	Local   bool `env:"LOCAL,required"`
}

type dbConf struct {
	DbName string `env:"DB_NAME,required"`
}

type serverConf struct {
	Port           int           `env:"SERVER_PORT,required"`
	TimeoutRead    time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite   time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle    time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Domain         string        `env:"DOMAIN,required"`          // server's domain
	AllowedOrigins []string      `env:"ALLOWED_ORIGINS,required"` // server's domain
}

// AppConfig Setup and install the applications' configuration environment variables
func AppConfig() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &c
}
