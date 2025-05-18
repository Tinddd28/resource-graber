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
	App        App        `yaml:"application"`
}

type Client struct {
	Host string `yaml:"host"`
}

type Network struct {
	Path    string        `yaml:"path"`
	Timeout time.Duration `yaml:"timeout"`
}

type Screenshot struct {
	Timeout     time.Duration `yaml:"timeout"`
	LoadTimeout time.Duration `yaml:"load_timeout"`
	Path        string        `yaml:"path"`
}

type App struct {
	Path string `yaml:"path"`
}

const configPath = "config/config.yaml"
const debugPath = "config.yaml"

func NewConfig() *Config {
	if _, err := os.Stat(debugPath); err != nil {
		panic(err)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(debugPath, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
