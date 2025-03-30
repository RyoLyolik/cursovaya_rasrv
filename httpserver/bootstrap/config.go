package bootstrap

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type SQLStorage struct {
	Host     string `yaml:"host" env:"PG_HOST" env-required:"true"`
	Port     string `yaml:"port" env:"PG_PORT" env-required:"true"`
	User     string `yaml:"user" env:"PG_USER" env-required:"true"`
	Password string `yaml:"password" env:"PG_PASSWORD" env-required:"true"`
	DBName   string `yaml:"dbname" env:"PG_DBNAME" env-required:"true"`
	SSLMode  string `yaml:"sslmode"`
}

type SessionStorage struct {
	Host string `yaml:"host" env:"REDIS_HOST" env-required:"true"`
	Port string `yaml:"port" env:"REDIS_PORT" env-required:"true"`
}

type LoggerConfig struct {
	Level   string `yaml:"level"`
	Handler string `yaml:"handler"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env:"HTTP_HOST" env-required:"true"`
	Port        int           `yaml:"port" env:"HTTP_PORT" env-required:"true"`
	TimeOut     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Config struct {
	Env            string         `yaml:"env" env:"ENVIRONMENT" env-required:"true"`
	SQLStorage     SQLStorage     `yaml:"sqlStorage"`
	HttpServer     HTTPServer     `yaml:"httpServer"`
	LogCfg         LoggerConfig   `yaml:"log"`
	AdminToken     string         `yaml:"adminToken" env:"ADMIN_TOKEN" env-required:"true"`
	SessionStorage SessionStorage `yaml:"sessionStorage"`
}

func MustLoad(envFile *string) *Config {
	if envFile != nil {
		godotenv.Load(*envFile)
	}
	configPath := os.Getenv("config_file")
	if configPath == "" {
		configPath = "./local.yaml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config gile does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
