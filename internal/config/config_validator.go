package config

import (
	"fmt"
	"net"
	"slices"
)

type validateFunction func(c *Config) error

var logLevels = []string{"debug", "info", "warning", "error"}

var logFormats = []string{"text", "json"}

var validatorList = []validateFunction{
	validateLogLevel,
	validatePort,
	validateIP,
	validateLogFormat,
}

func ValidateConfig(cfg Config) (Config, error) {
	for _, vFunc := range validatorList {
		if err := vFunc(&cfg); err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}

func validateLogLevel(cfg *Config) error {
	if !slices.Contains(logLevels, cfg.Logger.Level) {
		return fmt.Errorf("error incorrect log level %s", cfg.Logger.Level)
	}

	return nil
}

func validatePort(cfg *Config) error {
	if cfg.Exporter.ServerPort < 1 || cfg.Exporter.ServerPort > 65535 {
		return fmt.Errorf("error incorrect port number %d", cfg.Exporter.ServerPort)
	}

	return nil
}

func validateIP(cfg *Config) error {
	if cfg.Exporter.ServerAddress != "localhost" && net.ParseIP(cfg.Exporter.ServerAddress) == nil {
		return fmt.Errorf("error incorrect ip address %s", cfg.Exporter.ServerAddress)
	}

	return nil
}

func validateLogFormat(cfg *Config) error {
	if !slices.Contains(logFormats, cfg.Logger.Format) {
		return fmt.Errorf("error incorrect log format %s", cfg.Logger.Format)
	}

	return nil
}
