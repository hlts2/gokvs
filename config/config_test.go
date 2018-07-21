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
    port: "1111"

  - host: "127.0.0.1"
    port: "2222"

  - host: "127.0.0.1"
    port: "3333"

  - host: "127.0.0.1"
    port: "4444"
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
				Port: "1111",
			},
			Server{
				Host: "127.0.0.1",
				Port: "2222",
			},
			Server{
				Host: "127.0.0.1",
				Port: "3333",
			},
			Server{
				Host: "127.0.0.1",
				Port: "4444",
			},
		},
	}

	if !reflect.DeepEqual(conf, got) {
		t.Errorf("LoadConfig is wrong. expected: %v, got: %v", conf, got)
	}
}
