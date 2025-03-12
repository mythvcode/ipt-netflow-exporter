package exporter

import (
	"github.com/mythvcode/ipt-netflow-exporter/internal/statparser"
	"github.com/prometheus/client_golang/prometheus"
)

type CommonMetrics struct {
	inBitRate    *prometheus.GaugeVec
	inPacketRate *prometheus.GaugeVec
	inFlows      *prometheus.CounterVec
	inPackets    *prometheus.CounterVec
	inBytes      *prometheus.CounterVec
	hashMetric   *prometheus.GaugeVec
	hashMemory   *prometheus.GaugeVec
	hashFlows    *prometheus.GaugeVec
	hashPackets  *prometheus.GaugeVec
	hashBytes    *prometheus.GaugeVec
	dropPackets  *prometheus.CounterVec
	dropBytes    *prometheus.CounterVec
	outByteRate  *prometheus.GaugeVec
	outFlows     *prometheus.CounterVec
	outPackets   *prometheus.CounterVec
	outBytes     *prometheus.CounterVec
	lostFlows    *prometheus.CounterVec
	lostPackets  *prometheus.CounterVec
	lostBytes    *prometheus.CounterVec
	errTotal     *prometheus.CounterVec
	sndbufPeak   *prometheus.CounterVec
}

func newCommonMetricsCollector() *CommonMetrics {
	metrics := CommonMetrics{
		inBitRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "in_bit_rate",
				Help:      "Total incoming bits per second.",
			}, []string{},
		),
		inPacketRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "in_packet_rate",
				Help:      "Total incoming packets per second.",
			}, []string{},
		),
		inFlows: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "in_flows",
				Help:      "Total observed (metered) flow.",
			}, []string{},
		),
		inPackets: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "in_packets",
				Help:      "Total metered packets. Not counting dropped packets.",
			}, []string{},
		),
		inBytes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "in_bytes",
				Help:      "Total metered bytes in inPackets.",
			}, []string{},
		),
		hashMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "hash_metrics",
				Help:      "Measure of performance of hash table. When optimal should attract to 1.0, when non-optimal will be highly above of 1.",
			}, []string{},
		),
		hashMemory: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "hash_memory",
				Help:      "How much system memory is used by the hash table.",
			}, []string{},
		),
		hashFlows: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "hash_flows",
				Help:      "Flows currently residing in the hash table and not exported yet.",
			}, []string{},
		),
		hashPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "hash_packets",
				Help:      "Packets in flows currently residing in the hash table.",
			}, []string{},
		),
		hashBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "hash_bytes",
				Help:      "Bytes in flows currently residing in the hash table.",
			}, []string{},
		),
		dropPackets: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "drop_packets",
				Help:      "Total packets dropped by metering process.",
			}, []string{},
		),
		dropBytes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "drop_bytes",
				Help:      "Total bytes in packets dropped by metering process.",
			}, []string{},
		),
		outByteRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "out_byte_rate",
				Help:      "Total exporter output bytes per second.",
			}, []string{},
		),
		outFlows: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "out_flows",
				Help:      "Total exported flow data records.",
			}, []string{},
		),
		outPackets: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "out_packets",
				Help:      "Total exported packets of netflow stream itself.",
			}, []string{},
		),
		outBytes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "out_bytes",
				Help:      "Total exported bytes of netflow stream itself.",
			}, []string{},
		),
		lostFlows: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "lost_flows",
				Help:      "Total of accounted flows that are lost by exporting process due to socket errors. This value will not include asynchronous errors (cberr), these will be counted in err_total.",
			}, []string{},
		),
		lostPackets: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "lost_packets",
				Help:      "Total metered packets lost by exporting process. See lost_flows for details.",
			}, []string{},
		),
		lostBytes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "lost_bytes",
				Help:      "Total bytes in packets lost by exporting process. See lost_flows for details.",
			}, []string{},
		),
		errTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "lost_total",
				Help:      "Total exporting sockets errors (including cberr).",
			}, []string{},
		),
		sndbufPeak: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "sndbuf_peak",
				Help:      "Global maximum value of socket sndbuf. Sort of outputqueue length.",
			}, []string{},
		),
	}

	return &metrics
}

func (c *CommonMetrics) metricList() []iptNetFlowMetric {
	return []iptNetFlowMetric{
		c.inBitRate,
		c.inPacketRate,
		c.inFlows,
		c.inPackets,
		c.inBytes,
		c.hashMetric,
		c.hashMemory,
		c.hashFlows,
		c.hashPackets,
		c.hashBytes,
		c.dropPackets,
		c.dropBytes,
		c.outByteRate,
		c.outFlows,
		c.outPackets,
		c.outBytes,
		c.lostFlows,
		c.lostPackets,
		c.lostBytes,
		c.errTotal,
		c.sndbufPeak,
	}
}

func (c *CommonMetrics) reset() {
	for _, metric := range c.metricList() {
		metric.Reset()
	}
}

func (c *CommonMetrics) updateValues(stat *statparser.Statistics) {
	c.reset()

	c.inBitRate.With(prometheus.Labels{}).Set(float64(stat.InBitRate))
	c.inPacketRate.With(prometheus.Labels{}).Set(float64(stat.InPacketRate))
	c.inFlows.With(prometheus.Labels{}).Add(float64(stat.InBitRate))
	c.inPackets.With(prometheus.Labels{}).Add(float64(stat.InPackets))
	c.inBytes.With(prometheus.Labels{}).Add(float64(stat.InBytes))
	c.hashMetric.With(prometheus.Labels{}).Set(stat.HashMetric)
	c.hashMemory.With(prometheus.Labels{}).Set(float64(stat.HashMemory))
	c.hashFlows.With(prometheus.Labels{}).Set(float64(stat.HashFlows))
	c.hashPackets.With(prometheus.Labels{}).Set(float64(stat.HashPackets))
	c.hashBytes.With(prometheus.Labels{}).Set(float64(stat.HashBytes))
	c.dropPackets.With(prometheus.Labels{}).Add(float64(stat.DropPackets))
	c.dropBytes.With(prometheus.Labels{}).Add(float64(stat.DropBytes))
	c.outByteRate.With(prometheus.Labels{}).Set(float64(stat.OutByteRate))
	c.outFlows.With(prometheus.Labels{}).Add(float64(stat.OutFlows))
	c.outPackets.With(prometheus.Labels{}).Add(float64(stat.OutPackets))
	c.outBytes.With(prometheus.Labels{}).Add(float64(stat.OutBytes))
	c.lostFlows.With(prometheus.Labels{}).Add(float64(stat.LostFlows))
	c.lostPackets.With(prometheus.Labels{}).Add(float64(stat.LostPackets))
	c.lostBytes.With(prometheus.Labels{}).Add(float64(stat.LostBytes))
	c.errTotal.With(prometheus.Labels{}).Add(float64(stat.ErrTotal))
	c.sndbufPeak.With(prometheus.Labels{}).Add(float64(stat.ErrTotal))
}

func (c *CommonMetrics) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metricList() {
		metric.Describe(ch)
	}
}

func (c *CommonMetrics) Collect(metricChan chan<- prometheus.Metric) {
	for _, metric := range c.metricList() {
		metric.Collect(metricChan)
	}
}
