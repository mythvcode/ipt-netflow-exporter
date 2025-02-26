package exporter

import (
	"github.com/mythvcode/ipt-netflow-exporter/internal/statparser"
	"github.com/prometheus/client_golang/prometheus"
)

const cpuLabel = "cpu"

type CPUMetrics struct {
	cpuInPacketRate prometheus.GaugeVec
	cpuInFlows      prometheus.CounterVec
	cpuInPackets    prometheus.CounterVec
	cpuInBytes      prometheus.CounterVec
	cpuHashMetric   prometheus.GaugeVec
	cpuDropPackets  prometheus.CounterVec
	cpuDropBytes    prometheus.CounterVec
	cpuErrTrunc     prometheus.CounterVec
	cpuErrFrag      prometheus.CounterVec
	cpuErrAlloc     prometheus.CounterVec
	cpuErrMaxFlows  prometheus.CounterVec
}

func NewCPUMetrics() *CPUMetrics {
	return &CPUMetrics{
		cpuInPacketRate: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_in_packet_rate",
				Help:      "Incoming packets per second for this cpu.",
			}, []string{cpuLabel},
		),
		cpuInFlows: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_in_flows",
				Help:      "Flows metered on this cpu.",
			}, []string{cpuLabel},
		),
		cpuInPackets: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_in_packets",
				Help:      "Packets metered for cpu.",
			}, []string{cpuLabel},
		),
		cpuInBytes: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_in_bytes",
				Help:      "Bytes metered on this cpu.",
			}, []string{cpuLabel},
		),
		cpuHashMetric: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_hash_metric",
				Help:      "Measure of performance of hash table on this cpu.",
			}, []string{cpuLabel},
		),
		cpuDropPackets: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_drop_packets",
				Help:      "Packets dropped by metering process on this cpu.",
			}, []string{cpuLabel},
		),
		cpuDropBytes: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_drop_bytes",
				Help:      "Bytes in cpu_drop_packets for this cpu.",
			}, []string{cpuLabel},
		),
		cpuErrTrunc: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_err_trunc",
				Help:      "Truncated packets dropped for this cpu.",
			}, []string{cpuLabel},
		),
		cpuErrFrag: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_err_flag",
				Help:      "Fragmented packets dropped for this cpu.",
			}, []string{cpuLabel},
		),
		cpuErrAlloc: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_err_alloc",
				Help:      "Packets dropped due to memory allocation errors.",
			}, []string{cpuLabel},
		),
		cpuErrMaxFlows: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "cpu_err_max_flows",
				Help:      "Packets dropped due to maxflows limit being reached.",
			}, []string{cpuLabel},
		),
	}
}

func (c *CPUMetrics) metricList() []iptNetFlowMetric {
	return []iptNetFlowMetric{
		c.cpuInPacketRate,
		c.cpuInFlows,
		c.cpuInPackets,
		c.cpuInBytes,
		c.cpuHashMetric,
		c.cpuDropPackets,
		c.cpuDropBytes,
		c.cpuErrTrunc,
		c.cpuErrFrag,
		c.cpuErrAlloc,
		c.cpuErrMaxFlows,
	}
}

func (c *CPUMetrics) reset() {
	for _, metric := range c.metricList() {
		metric.Reset()
	}
}

func (c *CPUMetrics) updateValues(stat *statparser.Statistics) {
	c.reset()

	for _, cpuStat := range stat.CPUStatList {
		c.cpuInPacketRate.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Set(float64(cpuStat.CPUInPacketRate))
		c.cpuInFlows.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Add(float64(cpuStat.CPUInFlows))
		c.cpuInPackets.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Add(float64(cpuStat.CPUInPackets))
		c.cpuInBytes.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Add(float64(cpuStat.CPUInBytes))
		c.cpuHashMetric.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Set(float64(cpuStat.CPUHashMetric))
		c.cpuDropPackets.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Add(float64(cpuStat.CPUDropPackets))
		c.cpuDropBytes.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Add(float64(cpuStat.CPUuDropBytes))
		c.cpuErrTrunc.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Add(float64(cpuStat.CPUErrTrunc))
		c.cpuErrFrag.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Add(float64(cpuStat.CPUErrFrag))
		c.cpuErrAlloc.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Add(float64(cpuStat.CPUErrAlloc))
		c.cpuErrMaxFlows.With(prometheus.Labels{cpuLabel: cpuStat.CPU}).Add(float64(cpuStat.CPUErrMaxflows))
	}
}

func (c *CPUMetrics) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metricList() {
		metric.Describe(ch)
	}
}

func (c *CPUMetrics) Collect(metricChan chan<- prometheus.Metric) {
	for _, metric := range c.metricList() {
		metric.Collect(metricChan)
	}
}
