package config

import (
	"os"

	"github.com/pkg/errors"

	yaml "gopkg.in/yaml.v2"
)

// ErrDuplicateServer -
var ErrDuplicateServer = errors.New("duplicate server")

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

// Validation -
func (ss Servers) Validation() error {
	m := make(map[string]map[string]int)

	for _, server := range ss {
		if m[server.Name] == nil {
			m[server.Name] = make(map[string]int)
		}

		m[server.Name][server.Name]++
		m[server.Name][server.Host+":"+server.Port]++
	}

	for sname, mm := range m {
		for _, cnt := range mm {
			if cnt > 1 {
				return errors.WithMessage(ErrDuplicateServer, sname)
			}
		}
	}

	return nil
}

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

	err = conf.Servers.Validation()
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
