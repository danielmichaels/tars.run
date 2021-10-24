package config

import (
	"github.com/joeshaw/envdecode"
	"log"
	"time"
)

type limiter struct {
	Enabled bool    `env:"RATE_LIMIT_ENABLED,default=true"`
	Rps     float64 `env:"RATE_LIMIT_RPS,default=2"`
	Burst   int     `env:"RATE_LIMIT_BURST,default=4"`
}

type Conf struct {
	Debug   bool `env:"DEBUG,default=false"`
	Server  serverConf
	Db      dbConf
	Limiter limiter
}

type dbConf struct {
	DbName string `env:"DB_NAME,default=./data/shorty.db"`
}

type serverConf struct {
	Port           int           `env:"SERVER_PORT,default=1987"`
	TimeoutRead    time.Duration `env:"SERVER_TIMEOUT_READ,default=5s"`
	TimeoutWrite   time.Duration `env:"SERVER_TIMEOUT_WRITE,default=10s"`
	TimeoutIdle    time.Duration `env:"SERVER_TIMEOUT_IDLE,default=120s"`
	ApiDomain      string        `env:"API_DOMAIN,default=http://localhost:1987"`      // server's domain
	FrontendDomain string        `env:"FRONTEND_DOMAIN,default=http://localhost:1988"` // server's domain
	AllowedOrigins []string      `env:"ALLOWED_ORIGINS,default=http://localhost:1988"` // server's domain
}

// AppConfig Setup and install the applications' configuration environment variables
func AppConfig() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &c
}
