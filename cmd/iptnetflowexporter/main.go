package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mythvcode/ipt-netflow-exporter/internal/config"
	"github.com/mythvcode/ipt-netflow-exporter/internal/exporter"
	"github.com/mythvcode/ipt-netflow-exporter/internal/logger"
	"github.com/mythvcode/ipt-netflow-exporter/internal/statparser"
)

var cfgPath string

func init() {
	flag.StringVar(&cfgPath, "config", "", "Path to config file")
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	flag.Parse()
	cfg, err := config.ReadConfig(cfgPath)
	if err != nil {
		if cfgPath != "" {
			logger.Default().Errorf("Error read config file from file %s: %s", cfgPath, err.Error())
		} else {
			logger.Default().Errorf("Error read config: %s", err.Error())
		}
		os.Exit(1)
	}

	if err := logger.Init(cfg.Logger.File, cfg.Logger.Level, cfg.Logger.Format); err != nil {
		logger.Default().Errorf("error init logger %s", err.Error())
		os.Exit(1)
	}
	stat := statparser.New(cfg.Exporter.IPTNetFlowStatFile)
	exporter, err := exporter.New(cfg.Exporter, stat)
	if err != nil {
		logger.GetLogger().Errorf("error init exporter %s", err.Error())
	}
	started := make(chan error)
	go func() {
		started <- exporter.Start()
	}()
	go func() {
		if err := <-started; err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.GetLogger().Errorf("Error start exporter: %s", err.Error())
			os.Exit(1)
		}
	}()

	defer exporter.Stop()
	<-sigs
}
