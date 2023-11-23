package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const defaultConfigPath = "./configs/dev.yaml"

type (
	Config struct {
		HTTP           HTTPConfig     `yaml:"http"`
		Postgres       PostgresConfig `yaml:"postgres"`
		InsertTestData bool           `yaml:"insertTestData"`
	}

	HTTPConfig struct {
		Host           string        `yaml:"host"`
		Port           string        `yaml:"port"`
		WriteTimeOut   time.Duration `yaml:"write-timeout"`
		ReadTimeOut    time.Duration `yaml:"read-timeout"`
		MaxHeaderBytes int           `yaml:"maxHeaderBytes"`
	}

	PostgresConfig struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
		Database string `yaml:"database"`
		SSLMode  string `yaml:"ssl-mode"`
	}
)

func InitConfig(configPath string) *Config {

	if configPath == "" {
		configPath = defaultConfigPath
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
