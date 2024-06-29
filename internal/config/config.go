package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
)

type limiter struct {
	Enabled bool    `env:"RATE_LIMIT_ENABLED,default=true"`
	Rps     float64 `env:"RATE_LIMIT_RPS,default=2"`
	Burst   int     `env:"RATE_LIMIT_BURST,default=4"`
}

type Conf struct {
	Names   names
	Db      dbConf
	Server  serverConf
	Limiter limiter
}

type dbConf struct {
	DbName string `env:"DB_NAME,default=./data/data.db"`
}

type names struct {
	AppName          string `env:"APP_NAME,default=Tars.Run"`
	TwitterAccount   string `env:"TWITTER_ACCOUNT,default=#"`
	GithubAccount    string `env:"GITHUB_ACCOUNT,default=#"`
	PlausibleAccount string `env:"PLAUSIBLE_ACCOUNT,default="`
}

type serverConf struct {
	Domain       string        `env:"DOMAIN,default=http://localhost:1987"`
	LogLevel     string        `env:"LOG_LEVEL,default=info"`
	Port         int           `env:"SERVER_PORT,default=1987"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,default=5s"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,default=10s"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,default=120s"`
	LogConcise   bool          `env:"LOG_CONCISE,default=true"`
	LogJson      bool          `env:"LOG_JSON,default=false"`
}

// AppConfig Setup and install the applications' configuration environment variables
func AppConfig() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &c
}
