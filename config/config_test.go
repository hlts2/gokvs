package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var yamlFile = `
servers:
  -
    name: "server-1"
    host: "127.0.0.1"
    port: "1111"

  -
    name: "server-2"
    host: "127.0.0.1"
    port: "2222"

  -
    name: "server-3"
    host: "127.0.0.1"
    port: "3333"

  -
    name: "server-4"
    host: "127.0.0.1"
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
				Name: "server-1",
				Host: "127.0.0.1",
				Port: "1111",
			},
			Server{
				Name: "server-2",
				Host: "127.0.0.1",
				Port: "2222",
			},
			Server{
				Name: "server-3",
				Host: "127.0.0.1",
				Port: "3333",
			},
			Server{
				Name: "server-4",
				Host: "127.0.0.1",
				Port: "4444",
			},
		},
	}

	if !reflect.DeepEqual(conf, got) {
		t.Errorf("LoadConfig is wrong. expected: %v, got: %v", conf, got)
	}
}
