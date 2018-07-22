package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Config -
type Config struct {
	Servers Servers `yaml:"servers"`
}

// Server -
type Server struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Servers -
type Servers []Server

// LoadConfig -
func LoadConfig(fname string) (*Config, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	var conf Config
	err = yaml.NewDecoder(f).Decode(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
