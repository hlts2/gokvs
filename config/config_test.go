package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var yamlFile = `
servers:
  - host: "127.0.0.1"
    port: "1234"

  - host: "127.0.0.1"
    port: "5678"
`

func TestLoadConfig(t *testing.T) {
	filename := "test.yaml"
	err := ioutil.WriteFile(filename, []byte(yamlFile), os.ModePerm)
	if err != nil {
		t.Errorf("WriteFile is error: %v", err)
	}
	defer os.Remove(filename)

	conf, err := LoadConfig(filename)
	if err != nil {
		t.Errorf("LoadConfig is error: %v", err)
	}

	if conf == nil {
		t.Errorf("LoadConfig is nil")
	}

	got := &Config{
		Servers: Servers{
			Server{
				Host: "127.0.0.1",
				Port: "1234",
			},
			Server{
				Host: "127.0.0.1",
				Port: "5678",
			},
		},
	}

	if !reflect.DeepEqual(conf, got) {
		t.Errorf("LoadConfig is wrong. expected: %v, got: %v", conf, got)
	}
}
