package exporter

import (
	"github.com/mythvcode/ipt-netflow-exporter/internal/statparser"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	socketNameLabel = "socket"
	socketDstLabel  = "destination"
)

type SockMetrics struct {
	sockActive     *prometheus.CounterVec
	sockErrConnect *prometheus.CounterVec
	sockErrFull    *prometheus.CounterVec
	sockErrCberr   *prometheus.CounterVec
	sockErrOther   *prometheus.CounterVec
	sockSndbuf     *prometheus.GaugeVec
	sockSndbufFill *prometheus.GaugeVec
	sockSndbufPeak *prometheus.GaugeVec
}

func newSocketMetrics() *SockMetrics {
	return &SockMetrics{
		sockActive: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "socket_active",
				Help:      "Connection state of this socket.",
			}, []string{socketNameLabel, socketDstLabel},
		),
		sockErrConnect: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "socket_error_connect",
				Help: "Connections attempt count. High value usually mean " +
					"that network is not set up properly, or module is loaded " +
					"before network is up, in this case it is not dangerous" +
					"and should be ignored.",
			}, []string{socketNameLabel, socketDstLabel},
		),
		sockErrFull: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "socket_error_full",
				Help:      "Socket full errors on this socket. Usually mean sndbuf value is too small.",
			}, []string{socketNameLabel, socketDstLabel},
		),
		sockErrCberr: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "socket_error_cberr",
				Help: "Asynchronous callback errors on this socket. Usually mean " +
					"that there is 'connection refused' errors on UDP socket " +
					"reported via ICMP messages.",
			}, []string{socketNameLabel, socketDstLabel},
		),
		sockErrOther: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: metricsNamespace,
				Name:      "socket_error_other",
				Help:      "All other possible errors on this socket.",
			}, []string{socketNameLabel, socketDstLabel},
		),
		sockSndbuf: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "socket_snd_buf",
				Help:      "Sndbuf value for this socket. Higher value allows accommodate (exporting) traffic bursts.",
			}, []string{socketNameLabel, socketDstLabel},
		),
		sockSndbufFill: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "socket_snd_buf_fill",
				Help: "Amount of data currently in socket buffers. When this value " +
					"will reach size sndbuf, packet loss will occur.",
			}, []string{socketNameLabel, socketDstLabel},
		),
		sockSndbufPeak: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: metricsNamespace,
				Name:      "socket_snd_buf_peak",
				Help: "Historical peak amount of data in socket buffers. Useful to " +
					"evaluate sndbuf size, because sockSndbufFill is transient.",
			}, []string{socketNameLabel, socketDstLabel},
		),
	}
}

func (c *SockMetrics) metricList() []iptNetFlowMetric {
	return []iptNetFlowMetric{
		c.sockActive,
		c.sockErrConnect,
		c.sockErrFull,
		c.sockErrCberr,
		c.sockErrOther,
		c.sockSndbuf,
		c.sockSndbufFill,
		c.sockSndbufPeak,
	}
}

func (c *SockMetrics) reset() {
	for _, metric := range c.metricList() {
		metric.Reset()
	}
}

func (c *SockMetrics) updateValues(stat *statparser.Statistics) {
	c.reset()

	for _, sockStat := range stat.SockStatList {
		c.sockActive.With(
			prometheus.Labels{
				socketNameLabel: sockStat.SockName,
				socketDstLabel:  sockStat.SockDestination,
			}).Add(float64(sockStat.SockActive))
		c.sockErrConnect.With(
			prometheus.Labels{
				socketNameLabel: sockStat.SockName,
				socketDstLabel:  sockStat.SockDestination,
			}).Add(float64(sockStat.SockErrConnect))
		c.sockErrFull.With(
			prometheus.Labels{
				socketNameLabel: sockStat.SockName,
				socketDstLabel:  sockStat.SockDestination,
			}).Add(float64(sockStat.SockErrFull))
		c.sockErrCberr.With(
			prometheus.Labels{
				socketNameLabel: sockStat.SockName,
				socketDstLabel:  sockStat.SockDestination,
			}).Add(float64(sockStat.SockErrCberr))
		c.sockErrOther.With(
			prometheus.Labels{
				socketNameLabel: sockStat.SockName,
				socketDstLabel:  sockStat.SockDestination,
			}).Add(float64(sockStat.SockErrOther))
		c.sockSndbuf.With(
			prometheus.Labels{
				socketNameLabel: sockStat.SockName,
				socketDstLabel:  sockStat.SockDestination,
			}).Set(float64(sockStat.SockSndbuf))
		c.sockSndbufFill.With(
			prometheus.Labels{
				socketNameLabel: sockStat.SockName,
				socketDstLabel:  sockStat.SockDestination,
			}).Set(float64(sockStat.SockSndbufFill))
		c.sockSndbufPeak.With(
			prometheus.Labels{
				socketNameLabel: sockStat.SockName,
				socketDstLabel:  sockStat.SockDestination,
			}).Set(float64(sockStat.SockSndbufPeak))
	}
}

func (c *SockMetrics) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metricList() {
		metric.Describe(ch)
	}
}

func (c *SockMetrics) Collect(metricChan chan<- prometheus.Metric) {
	for _, metric := range c.metricList() {
		metric.Collect(metricChan)
	}
}
