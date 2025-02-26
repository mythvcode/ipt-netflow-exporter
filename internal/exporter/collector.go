package exporter

import (
	"log/slog"

	"github.com/mythvcode/ipt-netflow-exporter/internal/logger"
	"github.com/mythvcode/ipt-netflow-exporter/internal/statparser"
	"github.com/prometheus/client_golang/prometheus"
)

const metricsNamespace = "ipt_netflow"

type iptNetFlowMetric interface {
	prometheus.Collector
	Reset()
}

type iptNetFlowCollectors interface {
	prometheus.Collector
	updateValues(stat *statparser.Statistics)
}

type IPTNetFlowTCollector struct {
	statParser    StatParser
	log           *logger.Logger
	commonMetrics *CommonMetrics
	cpuMetrics    *CPUMetrics
	sockMetrics   *SockMetrics
}

func (i *IPTNetFlowTCollector) Name() string {
	return "ipt-netflow-collector"
}

func newIPTNetFlowTCollector(stat StatParser) *IPTNetFlowTCollector {
	return &IPTNetFlowTCollector{
		statParser:    stat,
		log:           logger.GetLogger().With(slog.String(logger.Component, "IPTNetFlowTCollector")),
		commonMetrics: newCommonMetricsCollector(),
		cpuMetrics:    NewCPUMetrics(),
		sockMetrics:   newSocketMetrics(),
	}
}

func (i *IPTNetFlowTCollector) Initialized() bool {
	return !(i.statParser == nil && i.log != nil)
}

func (i *IPTNetFlowTCollector) collectorList() []iptNetFlowCollectors {
	return []iptNetFlowCollectors{
		i.commonMetrics,
		i.cpuMetrics,
		i.sockMetrics,
	}
}

func (i *IPTNetFlowTCollector) Collect(metricChan chan<- prometheus.Metric) {
	metrics, err := i.statParser.CollectAndMarshal()
	if err != nil {
		i.log.Errorf("error collect metrics: %s", err.Error())
		// empty metrics
		metrics = statparser.Statistics{}
	}

	collectors := i.collectorList()

	for _, collector := range collectors {
		collector.updateValues(&metrics)
	}
	for _, collector := range collectors {
		collector.Collect(metricChan)
	}
}

func (i *IPTNetFlowTCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, collector := range i.collectorList() {
		collector.Describe(ch)
	}
}
