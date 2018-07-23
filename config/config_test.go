package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		yaml     string
		isErr    bool
		expected *Config
	}{
		{
			yaml: `
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
`,
			isErr: false,
			expected: &Config{
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
			},
		},
		{
			yaml: `
servers:
  -
    name: "server-1"
    host: "127.0.0.1"
    port: "1111"

  -
    name: "server-1"
    host: "127.0.0.1"
    port: "2222"
`,
			isErr:    true,
			expected: nil,
		},
		{
			yaml: `
servers:
  -
    name: "server-1"
    host: "127.0.0.1"
    port: "1111"

  -
    name: "server-2"
    host: "127.0.0.1"
    port: "1111"
`,
			isErr:    true,
			expected: nil,
		},
	}

	fname := "test.yaml"
	for i, test := range tests {
		func() {
			err := createFile(fname, []byte(test.yaml))
			if err != nil {
				t.Errorf("tests[%d] - createfile error: %v", i, err)
			}
			defer deleteFile(fname)

			conf, err := LoadConfig(fname)
			isErr := !(err == nil)

			if test.isErr != isErr {
				t.Errorf("tests[%d] - LoadConfig is wrong. isErr expected: %v, got: %v", i, test.isErr, isErr)
			}

			if !reflect.DeepEqual(test.expected, conf) {
				t.Errorf("tests[%d] - LoadConfig is wrong. expected: %v, got: %v", i, test.expected, conf)
			}
		}()
	}
}

func createFile(fname string, b []byte) error {
	err := ioutil.WriteFile(fname, b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func deleteFile(fname string) {
	os.Remove(fname)
}
