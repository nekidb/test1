package config

import (
	"bytes"
	"encoding/json"
	"io/fs"
)

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
	DB   string `json:"db"`
}

func Get(filesystem fs.FS, filename string) (*Config, error) {
	data, err := fs.ReadFile(filesystem, filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	reader := bytes.NewReader(data)
	if err := json.NewDecoder(reader).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
