package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Version string
	DBDriver string
	DBHost string
	DBUsername string
	DBPassword string
	DBName string
	DBCharset string
}

// Read and parse the configuration file
func (c *Config) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
