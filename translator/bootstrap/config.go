package bootstrap

import (
	"log"
	"os"

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

type Websocket struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Scheme string `yaml:"scheme"`
	Path   string `yaml:"path"`
}

type LoggerConfig struct {
	Level   string `yaml:"level"`
	Handler string `yaml:"handler"`
}

type Config struct {
	SQLStorage SQLStorage   `yaml:"sqlStorage"`
	WS         Websocket    `yaml:"ws"`
	LogCfg     LoggerConfig `yaml:"log"`
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
