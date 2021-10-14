package pre_refactor

import (
	"github.com/joeshaw/envdecode"
	"log"
	"time"
)

type serverConf struct {
	Port           int           `env:"SERVER_PORT,required"`
	TimeoutRead    time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite   time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle    time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Domain         string        `env:"DOMAIN,required"`          // server's domain
	AllowedOrigins string        `env:"ALLOWED_ORIGINS,required"` // server's domain
}

type Conf struct {
	Debug  bool `env:"DEBUG,required"`
	Server serverConf
	Db     dbConf
	Local  bool `env:"LOCAL,required"`
}

type dbConf struct {
	DbName string `env:"DB_NAME,required"`
}

// AppConfig Setup and install the applications' configuration environment variables
func AppConfig() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &c
}
