package config

import (
	"fmt"
	"os"
)

func readEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %q not found", key)
	}
	return value, nil
}

type ServerConfig struct {
	Host string
	Port string
}

func newServerConfig() (*ServerConfig, error) {
	host, err := readEnv("SERVER_HOST")
	if err != nil {
		return nil, err
	}

	port, err := readEnv("SERVER_PORT")
	if err != nil {
		return nil, err
	}

	config := &ServerConfig{
		Host: host,
		Port: port,
	}

	return config, nil
}

func (c ServerConfig) Address() string {
	return c.Host + ":" + c.Port
}

type Config struct {
	Server *ServerConfig
}

func NewConfig() (*Config, error) {
	serverConfig, err := newServerConfig()
	if err != nil {
		return nil, fmt.Errorf("server config: %w", err)
	}

	config := &Config{
		Server: serverConfig,
	}

	return config, nil
}
