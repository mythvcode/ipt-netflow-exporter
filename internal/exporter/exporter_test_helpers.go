package exporter

import (
	"testing"

	"github.com/mythvcode/ipt-netflow-exporter/internal/statparser"
)

func getTestStatistic(t *testing.T) statparser.Statistics {
	t.Helper()

	return statparser.Statistics{
		InBitRate:    1,
		InPacketRate: 2,
		InFlows:      3,
		InPackets:    4,
		InBytes:      5,
		HashMetric:   5.5,
		HashMemory:   6,
		HashFlows:    7,
		HashPackets:  8,
		HashBytes:    9,
		DropPackets:  10,
		DropBytes:    11,
		OutByteRate:  12,
		OutFlows:     13,
		OutPackets:   14,
		OutBytes:     15,
		LostFlows:    16,
		LostPackets:  17,
		LostBytes:    18,
		ErrTotal:     19,
		SndbufPeak:   20,
		CPUStatList: []statparser.CPUStat{
			{
				CPU:             "cpu0",
				CPUInPacketRate: 1,
				CPUInFlows:      2,
				CPUInPackets:    3,
				CPUInBytes:      4,
				CPUHashMetric:   5.55,
				CPUDropPackets:  6,
				CPUuDropBytes:   7,
				CPUErrTrunc:     8,
				CPUErrFrag:      9,
				CPUErrAlloc:     10,
				CPUErrMaxflows:  11,
			},
		},
		SockStatList: []statparser.NFSockEntry{
			{
				SockName:        "sock0",
				SockDestination: "localhost:1234",
				SockActive:      1,
				SockErrConnect:  2,
				SockErrFull:     3,
				SockErrCberr:    4,
				SockErrOther:    5,
				SockSndbuf:      6,
				SockSndbufFill:  7,
				SockSndbufPeak:  8,
			},
		},
	}
}

func getPromTestStat(t *testing.T) string {
	return `
	# HELP ipt_netflow_cpu_drop_bytes Bytes in cpu_drop_packets for this cpu.
	# TYPE ipt_netflow_cpu_drop_bytes counter
	ipt_netflow_cpu_drop_bytes{cpu="cpu0"} 7
	# HELP ipt_netflow_cpu_drop_packets Packets dropped by metering process on this cpu.
	# TYPE ipt_netflow_cpu_drop_packets counter
	ipt_netflow_cpu_drop_packets{cpu="cpu0"} 6
	# HELP ipt_netflow_cpu_err_alloc Packets dropped due to memory allocation errors.
	# TYPE ipt_netflow_cpu_err_alloc counter
	ipt_netflow_cpu_err_alloc{cpu="cpu0"} 10
	# HELP ipt_netflow_cpu_err_flag Fragmented packets dropped for this cpu.
	# TYPE ipt_netflow_cpu_err_flag counter
	ipt_netflow_cpu_err_flag{cpu="cpu0"} 9
	# HELP ipt_netflow_cpu_err_max_flows Packets dropped due to maxflows limit being reached.
	# TYPE ipt_netflow_cpu_err_max_flows counter
	ipt_netflow_cpu_err_max_flows{cpu="cpu0"} 11
	# HELP ipt_netflow_cpu_err_trunc Truncated packets dropped for this cpu.
	# TYPE ipt_netflow_cpu_err_trunc counter
	ipt_netflow_cpu_err_trunc{cpu="cpu0"} 8
	# HELP ipt_netflow_cpu_hash_metric Measure of performance of hash table on this cpu.
	# TYPE ipt_netflow_cpu_hash_metric gauge
	ipt_netflow_cpu_hash_metric{cpu="cpu0"} 5.55
	# HELP ipt_netflow_cpu_in_bytes Bytes metered on this cpu.
	# TYPE ipt_netflow_cpu_in_bytes counter
	ipt_netflow_cpu_in_bytes{cpu="cpu0"} 4
	# HELP ipt_netflow_cpu_in_flows Flows metered on this cpu.
	# TYPE ipt_netflow_cpu_in_flows counter
	ipt_netflow_cpu_in_flows{cpu="cpu0"} 2
	# HELP ipt_netflow_cpu_in_packet_rate Incoming packets per second for this cpu.
	# TYPE ipt_netflow_cpu_in_packet_rate gauge
	ipt_netflow_cpu_in_packet_rate{cpu="cpu0"} 1
	# HELP ipt_netflow_cpu_in_packets Packets metered for cpu.
	# TYPE ipt_netflow_cpu_in_packets counter
	ipt_netflow_cpu_in_packets{cpu="cpu0"} 3
	# HELP ipt_netflow_drop_bytes Total bytes in packets dropped by metering process.
	# TYPE ipt_netflow_drop_bytes counter
	ipt_netflow_drop_bytes 11
	# HELP ipt_netflow_drop_flows Total exported flow data records.
	# TYPE ipt_netflow_drop_flows counter
	ipt_netflow_drop_flows 13
	# HELP ipt_netflow_drop_packets Total packets dropped by metering process.
	# TYPE ipt_netflow_drop_packets counter
	ipt_netflow_drop_packets 10
	# HELP ipt_netflow_hash_bytes Bytes in flows currently residing in the hash table.
	# TYPE ipt_netflow_hash_bytes gauge
	ipt_netflow_hash_bytes 9
	# HELP ipt_netflow_hash_flows Flows currently residing in the hash table and not exported yet.
	# TYPE ipt_netflow_hash_flows gauge
	ipt_netflow_hash_flows 7
	# HELP ipt_netflow_hash_memory How much system memory is used by the hash table.
	# TYPE ipt_netflow_hash_memory gauge
	ipt_netflow_hash_memory 6
	# HELP ipt_netflow_hash_metrics Measure of performance of hash table. When optimal should attract to 1.0, when non-optimal will be highly above of 1.
	# TYPE ipt_netflow_hash_metrics gauge
	ipt_netflow_hash_metrics 5.5
	# HELP ipt_netflow_hash_packets Packets in flows currently residing in the hash table.
	# TYPE ipt_netflow_hash_packets gauge
	ipt_netflow_hash_packets 8
	# HELP ipt_netflow_in_bit_rate Total incoming bits per second.
	# TYPE ipt_netflow_in_bit_rate gauge
	ipt_netflow_in_bit_rate 1
	# HELP ipt_netflow_in_bytes Total metered bytes in inPackets.
	# TYPE ipt_netflow_in_bytes counter
	ipt_netflow_in_bytes 5
	# HELP ipt_netflow_in_flows Total observed (metered) flow.
	# TYPE ipt_netflow_in_flows counter
	ipt_netflow_in_flows 1
	# HELP ipt_netflow_in_packet_rate Total incoming packets per second.
	# TYPE ipt_netflow_in_packet_rate gauge
	ipt_netflow_in_packet_rate 2
	# HELP ipt_netflow_in_packets Total metered packets. Not counting dropped packets.
	# TYPE ipt_netflow_in_packets counter
	ipt_netflow_in_packets 4
	# HELP ipt_netflow_lost_bytes Total bytes in packets lost by exporting process. See lost_flows for details.
	# TYPE ipt_netflow_lost_bytes counter
	ipt_netflow_lost_bytes 18
	# HELP ipt_netflow_lost_flows Total of accounted flows that are lost by exporting process due to socket errors. This value will not include asynchronous errors (cberr), these will be counted in err_total.
	# TYPE ipt_netflow_lost_flows counter
	ipt_netflow_lost_flows 16
	# HELP ipt_netflow_lost_packets Total metered packets lost by exporting process. See lost_flows for details.
	# TYPE ipt_netflow_lost_packets counter
	ipt_netflow_lost_packets 17
	# HELP ipt_netflow_lost_total Total exporting sockets errors (including cberr).
	# TYPE ipt_netflow_lost_total counter
	ipt_netflow_lost_total 19
	# HELP ipt_netflow_out_byte_rate Total exporter output bytes per second.
	# TYPE ipt_netflow_out_byte_rate gauge
	ipt_netflow_out_byte_rate 12
	# HELP ipt_netflow_out_bytes Total exported bytes of netflow stream itself.
	# TYPE ipt_netflow_out_bytes counter
	ipt_netflow_out_bytes 15
	# HELP ipt_netflow_out_packets Total exported packets of netflow stream itself.
	# TYPE ipt_netflow_out_packets counter
	ipt_netflow_out_packets 14
	# HELP ipt_netflow_sndbuf_peak Global maximum value of socket sndbuf. Sort of outputqueue length.
	# TYPE ipt_netflow_sndbuf_peak counter
	ipt_netflow_sndbuf_peak 19
	# HELP ipt_netflow_socket_active Connection state of this socket.
	# TYPE ipt_netflow_socket_active counter
	ipt_netflow_socket_active{destination="localhost:1234",socket="sock0"} 1
	# HELP ipt_netflow_socket_error_cberr Asynchronous callback errors on this socket. Usually mean that there is 'connection refused' errors on UDP socket reported via ICMP messages.
	# TYPE ipt_netflow_socket_error_cberr counter
	ipt_netflow_socket_error_cberr{destination="localhost:1234",socket="sock0"} 4
	# HELP ipt_netflow_socket_error_connect Connections attempt count. High value usually mean that network is not set up properly, or module is loaded before network is up, in this case it is not dangerousand should be ignored.
	# TYPE ipt_netflow_socket_error_connect counter
	ipt_netflow_socket_error_connect{destination="localhost:1234",socket="sock0"} 2
	# HELP ipt_netflow_socket_error_full Socket full errors on this socket. Usually mean sndbuf value is too small.
	# TYPE ipt_netflow_socket_error_full counter
	ipt_netflow_socket_error_full{destination="localhost:1234",socket="sock0"} 3
	# HELP ipt_netflow_socket_error_other All other possible errors on this socket.
	# TYPE ipt_netflow_socket_error_other counter
	ipt_netflow_socket_error_other{destination="localhost:1234",socket="sock0"} 5
	# HELP ipt_netflow_socket_snd_buf Sndbuf value for this socket. Higher value allows accommodate (exporting) traffic bursts.
	# TYPE ipt_netflow_socket_snd_buf gauge
	ipt_netflow_socket_snd_buf{destination="localhost:1234",socket="sock0"} 6
	# HELP ipt_netflow_socket_snd_buf_fill Amount of data currently in socket buffers. When this value will reach size sndbuf, packet loss will occur.
	# TYPE ipt_netflow_socket_snd_buf_fill gauge
	ipt_netflow_socket_snd_buf_fill{destination="localhost:1234",socket="sock0"} 7
	# HELP ipt_netflow_socket_snd_buf_peak Historical peak amount of data in socket buffers. Useful to evaluate sndbuf size, because sockSndbufFill is transient.
	# TYPE ipt_netflow_socket_snd_buf_peak gauge
	ipt_netflow_socket_snd_buf_peak{destination="localhost:1234",socket="sock0"} 8
	`
}
