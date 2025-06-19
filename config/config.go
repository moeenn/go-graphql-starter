package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

const (
	defaultServerHost       = "0.0.0.0"
	defaultServerPort       = "5000"
	defaultServerTimeout    = time.Second * 5
	defaultGraphqlUrl       = "/graphql"
	defaultGraphiqlUrl      = "/graphiql"
	defaultJwtExpiryMinutes = time.Minute * 60 * 24 // 24 hours.
)

var (
	defaultGraphqlWhitelistOperations = []string{
		"login",
	}
)

func readEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %q not found", key)
	}
	return value, nil
}

type ServerConfig struct {
	Host           string
	Port           string
	Timeout        time.Duration
	GraphqlUrl     string
	GraphiqlUrl    string
	ExposeGraphiql bool
}

func newServerConfig(logger *slog.Logger) *ServerConfig {
	host, err := readEnv("SERVER_HOST")
	if err != nil {
		host = defaultServerHost
	}

	port, err := readEnv("SERVER_PORT")
	if err != nil {
		port = defaultServerPort
	}

	exposeGraphiql := true
	exposeGraphqlRaw, err := readEnv("EXPOSE_GRAPHIQL")
	if err == nil {
		exposeGraphiql = exposeGraphqlRaw == "true" || exposeGraphqlRaw == "1"
	}

	if exposeGraphiql {
		logger.Info("GraphiQL playground running", "path", defaultGraphiqlUrl)
	} else {
		logger.Info("GraphiQL playground not running")
	}

	return &ServerConfig{
		Host:           host,
		Port:           port,
		Timeout:        defaultServerTimeout,
		GraphqlUrl:     defaultGraphqlUrl,
		GraphiqlUrl:    defaultGraphiqlUrl,
		ExposeGraphiql: exposeGraphiql,
	}
}
func (c ServerConfig) Address() string {
	return c.Host + ":" + c.Port
}

type DatabaseConfig struct {
	URI string
}

func newDatabaseConfig() (*DatabaseConfig, error) {
	uri, err := readEnv("DATABASE_URI")
	if err != nil {
		return nil, err
	}

	config := &DatabaseConfig{
		URI: uri,
	}

	return config, nil
}

type AuthConfig struct {
	JwtSecret             []byte
	JwtExpiryMinutes      time.Duration
	WhiteListedOperations []string
}

func newAuthConfig(logger *slog.Logger, exposeGraphiql bool) (*AuthConfig, error) {
	secret, err := readEnv("JWT_SECRET")
	if err != nil {
		return nil, fmt.Errorf("auth config: %w", err)
	}

	var expiryMinutes time.Duration
	expiryMinutesString, err := readEnv("JWT_EXPIRY_MINUTES")
	if err != nil {
		expiryMinutes = defaultJwtExpiryMinutes
	} else {
		expiryMinutesNum, err := strconv.ParseInt(expiryMinutesString, 10, 64)
		if err != nil {
			logger.Warn(
				"env variable JWT_EXPIRY_MINUTES does not contain a valid unsigned number, using default value",
				"value", expiryMinutesString, "default", defaultJwtExpiryMinutes)

			expiryMinutes = defaultJwtExpiryMinutes
		} else {
			expiryMinutes = time.Minute * time.Duration(expiryMinutesNum)
		}
	}

	whiteListOperations := defaultGraphqlWhitelistOperations
	if exposeGraphiql {
		whiteListOperations = append(whiteListOperations, "__schema")
	}

	config := &AuthConfig{
		JwtSecret:             []byte(secret),
		JwtExpiryMinutes:      expiryMinutes,
		WhiteListedOperations: whiteListOperations,
	}

	return config, nil
}

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	Auth     *AuthConfig
}

func NewConfig(logger *slog.Logger) (*Config, error) {
	serverConfig := newServerConfig(logger)
	databaseConfig, err := newDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("database config: %w", err)
	}

	authConfig, err := newAuthConfig(logger, serverConfig.ExposeGraphiql)
	if err != nil {
		return nil, fmt.Errorf("auth config: %w", err)
	}

	config := &Config{
		Server:   serverConfig,
		Database: databaseConfig,
		Auth:     authConfig,
	}

	return config, nil
}
