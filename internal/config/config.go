package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	Arg = "config"
)

type Config struct {
	DB      DBConfig    `json:"db"`
	Minio   MinioConfig `json:"minio"`
	Service HttpService `json:"service"`
}

type DBConfig struct {
	DSN string `json:"dsn"`
}

type MinioConfig struct {
	Endpoint        string `json:"endpoint"`
	PublicEndpoint  string `json:"public_endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

type HttpService struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

func Must(path string) Config {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("read config file: %w", err))
	}

	c := Config{}

	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Errorf("unmarshall config: %w", err))
	}

	return c
}
