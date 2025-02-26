package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testConfig = `
---
# Logger options.
logger:
  file: "test_file"
  level: info
  format: text
exporter:
  enable_runtime_metrics: true
  server_address: yaml_test_address
  server_port:    1010
  request_timeout: 55555
  telemetry_path: /test_conf_path
  ipt_netflow_stat: /proc/net/stat/config_stat
`

func testDefaults(t *testing.T, cfg Config) {
	t.Helper()
	require.Empty(t, cfg.Logger.File)
	require.Equal(t, "debug", cfg.Logger.Level)
	require.Equal(t, "json", cfg.Logger.Format)
	require.False(t, cfg.Exporter.EnableRuntimeMetrics)
	require.Equal(t, "/proc/net/stat/ipt_netflow_snmp", cfg.Exporter.IPTNetFlowStatFile)
	require.Equal(t, 10, cfg.Exporter.RequestTimeout)
	require.Equal(t, "localhost", cfg.Exporter.ServerAddress)
	require.Equal(t, 8080, cfg.Exporter.ServerPort)
	require.Equal(t, "/metrics", cfg.Exporter.TelemetryPath)
}

func setEnvVars(t *testing.T) {
	t.Helper()
	envVars := []struct {
		envName string
		value   string
	}{
		{
			"EXPORTER_LOG_LEVEL",
			"error",
		},
		{
			"EXPORTER_LOG_FORMAT",
			"text",
		},
		{
			"EXPORTER_LOG_FILE",
			"env_log_file",
		},
		{
			"EXPORTER_HOST",
			"1.2.3.4",
		},
		{
			"EXPORTER_PORT",
			"12345",
		},
		{
			"EXPORTER_REQUEST_TIMEOUT",
			"11111",
		},
		{
			"EXPORTER_TELEMETRY_PATH",
			"/test_path",
		},
		{
			"EXPORTER_ENABLE_RUNTIME_METRICS",
			"true",
		},
		{
			"EXPORTER_IPT_NETFLOW_STAT",
			"env_file_stat",
		},
	}
	for _, env := range envVars {
		t.Setenv(env.envName, env.value)
	}

	t.Cleanup(func() {
		for _, env := range envVars {
			require.NoError(t, os.Unsetenv(env.envName))
		}
	})
}

func TestCheckDefaults(t *testing.T) {
	cfg, err := loadFromBytes([]byte{})
	require.NoError(t, err)
	testDefaults(t, cfg)
}

func TestLoadFromEnv(t *testing.T) {
	cfg, err := ReadConfig("")
	require.NoError(t, err)
	testDefaults(t, cfg)
	setEnvVars(t)
	cfg, err = ReadConfig("")
	require.NoError(t, err)
	require.Equal(t, "env_log_file", cfg.Logger.File)
	require.Equal(t, "error", cfg.Logger.Level)
	require.Equal(t, "text", cfg.Logger.Format)
	require.True(t, cfg.Exporter.EnableRuntimeMetrics)
	require.Equal(t, "env_file_stat", cfg.Exporter.IPTNetFlowStatFile)
	require.Equal(t, 11111, cfg.Exporter.RequestTimeout)
	require.Equal(t, "1.2.3.4", cfg.Exporter.ServerAddress)
	require.Equal(t, 12345, cfg.Exporter.ServerPort)
	require.Equal(t, "/test_path", cfg.Exporter.TelemetryPath)
}

func TestLoadFromFile(t *testing.T) {
	cfg, err := loadFromBytes([]byte(testConfig))
	require.NoError(t, err)
	require.Equal(t, "test_file", cfg.Logger.File)
	require.Equal(t, "info", cfg.Logger.Level)
	require.Equal(t, "text", cfg.Logger.Format)
	require.True(t, cfg.Exporter.EnableRuntimeMetrics)
	require.Equal(t, "/proc/net/stat/config_stat", cfg.Exporter.IPTNetFlowStatFile)
	require.Equal(t, 55555, cfg.Exporter.RequestTimeout)
	require.Equal(t, "yaml_test_address", cfg.Exporter.ServerAddress)
	require.Equal(t, 1010, cfg.Exporter.ServerPort)
	require.Equal(t, "/test_conf_path", cfg.Exporter.TelemetryPath)
}

func TestValidators(t *testing.T) {
	cfg := getDefault()
	testDefaults(t, cfg)
	tCases := []struct {
		cfg   Config
		error string
	}{
		{
			cfg: func() (cfg Config) {
				cfg = getDefault()
				cfg.Logger.Format = "not_exist"

				return
			}(),
			error: "error incorrect log format not_exist",
		},
		{
			cfg: func() (cfg Config) {
				cfg = getDefault()
				cfg.Logger.Level = "error_log_level"

				return
			}(),
			error: "error incorrect log level error_log_level",
		},
		{
			cfg: func() (cfg Config) {
				cfg = getDefault()
				cfg.Exporter.ServerAddress = "error_address"

				return
			}(),
			error: "error incorrect ip address error_address",
		},
		{
			cfg: func() (cfg Config) {
				cfg = getDefault()
				cfg.Exporter.ServerPort = 1234123121

				return
			}(),
			error: "error incorrect port number 1234123121",
		},
	}

	for _, tCase := range tCases {
		_, err := ValidateConfig(tCase.cfg)
		require.Error(t, err)
		require.Equal(t, tCase.error, err.Error())
	}
}
