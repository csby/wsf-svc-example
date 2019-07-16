package main

import (
	"encoding/json"
	"fmt"
	"github.com/csby/wsf/server/configure"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	sync.RWMutex
	configure.Configure
}

func NewConfig() *Config {
	return &Config{
		Configure: configure.Configure{
			Log: configure.Log{
				Folder: "",
				Level:  "error|warning|info",
			},
			Http: configure.Http{
				Enabled: true,
				Port:    80,
			},
			Https: configure.Https{
				Enabled: false,
				Port:    443,
			},
			Root: "",
			Document: configure.Document{
				Enabled: true,
				Root:    "",
			},
			Operation: configure.Operation{
				Root: "",
				Api: configure.Api{
					Token: configure.Token{
						Expiration: 30,
					},
				},
				Users: []configure.User{
					{
						Account:  "admin",
						Password: "1",
					},
				},
				Ldap: configure.Ldap{
					Enable: false,
					Host:   "192.168.1.80",
					Port:   389,
					Base:   "dc=waf-example,dc=com",
				},
			},
		},
	}
}

func (s *Config) LoadFromFile(filePath string) error {
	s.Lock()
	defer s.Unlock()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

func (s *Config) SaveToFile(filePath string) error {
	s.Lock()
	defer s.Unlock()

	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	fileFolder := filepath.Dir(filePath)
	_, err = os.Stat(fileFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(fileFolder, 0777)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, string(bytes[:]))

	return err
}

func (s *Config) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

func (s *Config) FormatString() string {
	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes[:])
}
