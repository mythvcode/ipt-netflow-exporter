package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/creasty/defaults"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger   Logger   `env:", prefix=EXPORTER_" yaml:"logger"`
	Exporter Exporter `env:", prefix=EXPORTER_" yaml:"exporter"`
}

type Logger struct {
	Format string `default:"json"  env:"LOG_FORMAT" yaml:"format"`
	Level  string `default:"debug" env:"LOG_LEVEL"  yaml:"level"`
	File   string `default:""      env:"LOG_FILE"   yaml:"file"`
}

type Exporter struct {
	ServerAddress        string `default:"localhost"                       env:"HOST"                   yaml:"server_address"`
	ServerPort           int    `default:"8080"                            env:"PORT"                   yaml:"server_port"`
	RequestTimeout       int    `default:"10"                              env:"REQUEST_TIMEOUT"        yaml:"request_timeout"`
	TelemetryPath        string `default:"/metrics"                        env:"TELEMETRY_PATH"         yaml:"telemetry_path"`
	IPTNetFlowStatFile   string `default:"/proc/net/stat/ipt_netflow_snmp" env:"IPT_NETFLOW_STAT"       yaml:"ipt_netflow_stat"`
	EnableRuntimeMetrics bool   `default:"false"                           env:"ENABLE_RUNTIME_METRICS" yaml:"enable_runtime_metrics"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := defaults.Set(c); err != nil {
		return err
	}
	type plain Config

	return unmarshal((*plain)(c))
}

func ReadEnv(cfg Config) (Config, error) {
	err := envconfig.Process(
		context.Background(),
		&envconfig.Config{DefaultOverwrite: true, Target: &cfg},
	)

	return cfg, err
}

func ReadConfig(file string) (Config, error) {
	var err error
	var cfg Config
	if file == "" {
		cfg, err = ReadEnv(getDefault())
	} else {
		cfg, err = loadFromFile(file)
	}
	if err != nil {
		return cfg, err
	}

	return ValidateConfig(cfg)
}

func loadFromFile(file string) (Config, error) {
	configBytes, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return Config{}, fmt.Errorf("unable to read config: %w", err)
	}

	return loadFromBytes(configBytes)
}

func getDefault() Config {
	res, _ := loadFromBytes([]byte{})

	return res
}

func loadFromBytes(data []byte) (Config, error) {
	var config Config

	// make empty config for defaults package to call function UnmarshalYAML
	if len(data) == 0 {
		data = []byte("exporter:")
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return config, nil
}
