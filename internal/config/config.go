package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Client     Client     `yaml:"client"`
	Network    Network    `yaml:"network"`
	Screenshot Screenshot `yaml:"screenshot"`
	App        App        `yaml:"app"`
}

type Client struct {
	Host string `yaml:"host"`
}

type Network struct {
	Path string `yaml:"path"`
}

type Screenshot struct {
	Timeout     time.Time `yaml:"timeout"`
	LoadTimeout time.Time `yaml:"load_timeout"`
	Path        string    `yaml:"path"`
}

type App struct {
	Path string `yaml:"path"`
}

const configPath = "config/config.yaml"

func NewConfig() *Config {
	if _, err := os.Stat(configPath); err != nil {
		panic(err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
