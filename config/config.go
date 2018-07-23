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
	m := make(map[string]int)

	for _, server := range ss {
		m[server.Name]++
		m[server.Host+":"+server.Port]++
	}

	for field, cnt := range m {
		if cnt > 1 {
			return errors.WithMessage(ErrDuplicateServer, ss.getName(field))
		}
	}

	return nil
}

func (ss Servers) getName(field string) string {
	for _, server := range ss {
		if server.Name == field || server.Host+":"+server.Port == field {
			return server.Name
		}
	}
	return ""
}

// GetHostAndPorts -
func (ss Servers) GetHostAndPorts() []string {
	ips := make([]string, 0, len(ss))

	for _, s := range ss {
		ips = append(ips, s.Host+":"+s.Name)
	}
	return ips
}

// GetServer -
func (ss Servers) GetServer(sname string) *Server {
	for _, server := range ss {
		if server.Name == sname {
			return &server
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
