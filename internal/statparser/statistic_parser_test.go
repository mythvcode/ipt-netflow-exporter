package statparser

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/mythvcode/ipt-netflow-exporter/internal/logger"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.SetDefaultDiscardLogger()
}

const fileContent = `
inBitRate    1
inPacketRate 2
inFlows      3
inPackets    4
inBytes      562004
hashMetric   1.03
hashMemory   2560
hashFlows    3
hashPackets  973
hashBytes    5620
dropPackets  4
dropBytes    5
outByteRate  6
outFlows     15894
outPackets   105
outBytes     1551
lostFlows    7
lostPackets  1
lostBytes    9
errTotal     10
cpu0 1 2 3 4 1.35 5 6 7 8 9 10
cpu1 1 2 3 4 1.35 5 6 7 8 9 10
cpu2 1 2 3 4 1.35 5 6 7 8 9 10
sock0 127.0.0.1:1234 1 2 3 4 5 263 6 7
sock1 127.0.0.1:5555 1 2 3 4 5 263 6 7
sndbufPeak   0
`

func setReadFileFunc(t *testing.T, metrics string, err error) {
	t.Helper()
	prev := readFile

	readFile = func(string) ([]byte, error) {
		return []byte(metrics), err
	}
	t.Cleanup(func() { readFile = prev })
}

func testCPUStat(t *testing.T, index int, cpuStat CPUStat) {
	t.Helper()
	require.Equal(t, fmt.Sprintf("cpu%d", index), cpuStat.CPU)
	require.Equal(t, uint64(1), cpuStat.CPUInPacketRate)
	require.Equal(t, uint64(2), cpuStat.CPUInFlows)
	require.Equal(t, uint64(3), cpuStat.CPUInPackets)
	require.Equal(t, uint64(4), cpuStat.CPUInBytes)
	require.InEpsilon(t, float64(1.35), cpuStat.CPUHashMetric, 0.0001)
	require.Equal(t, uint64(5), cpuStat.CPUDropPackets)
	require.Equal(t, uint64(6), cpuStat.CPUuDropBytes)
	require.Equal(t, uint64(7), cpuStat.CPUErrTrunc)
	require.Equal(t, uint64(8), cpuStat.CPUErrFrag)
	require.Equal(t, uint64(9), cpuStat.CPUErrAlloc)
	require.Equal(t, uint64(10), cpuStat.CPUErrMaxflows)
}

func testSockStat(t *testing.T, index int, sockStat NFSockEntry) {
	t.Helper()
	var sockDST string
	if index == 0 {
		sockDST = "127.0.0.1:1234"
	} else {
		sockDST = "127.0.0.1:5555"
	}
	require.Equal(t, fmt.Sprintf("sock%d", index), sockStat.SockName)
	require.Equal(t, sockDST, sockStat.SockDestination)
	require.Equal(t, uint32(1), sockStat.SockActive)
	require.Equal(t, uint32(2), sockStat.SockErrConnect)
	require.Equal(t, uint32(3), sockStat.SockErrFull)
	require.Equal(t, uint32(4), sockStat.SockErrCberr)
	require.Equal(t, uint32(5), sockStat.SockErrOther)
	require.Equal(t, uint32(263), sockStat.SockSndbuf)
	require.Equal(t, uint32(6), sockStat.SockSndbufFill)
	require.Equal(t, uint32(7), sockStat.SockSndbufPeak)
}

func testDefaults(t *testing.T, stat Statistics) {
	t.Helper()
	require.Equal(t, uint64(1), stat.InBitRate)
	require.Equal(t, uint64(2), stat.InPacketRate)
	require.Equal(t, uint64(3), stat.InFlows)
	require.Equal(t, uint64(4), stat.InPackets)
	require.Equal(t, uint64(562004), stat.InBytes)
	require.InEpsilon(t, float64(1.03), stat.HashMetric, 0.0001)
	require.Equal(t, uint64(2560), stat.HashMemory)
	require.Equal(t, uint64(3), stat.HashFlows)
	require.Equal(t, uint64(973), stat.HashPackets)
	require.Equal(t, uint64(5620), stat.HashBytes)
	require.Equal(t, uint64(4), stat.DropPackets)
	require.Equal(t, uint64(5), stat.DropBytes)
	require.Equal(t, uint64(6), stat.OutByteRate)
	require.Equal(t, uint64(15894), stat.OutFlows)
	require.Equal(t, uint64(105), stat.OutPackets)
	require.Equal(t, uint64(1551), stat.OutBytes)
	require.Equal(t, uint64(7), stat.LostFlows)
	require.Equal(t, uint64(1), stat.LostPackets)
	require.Equal(t, uint64(9), stat.LostBytes)
	require.Equal(t, uint64(10), stat.ErrTotal)
	require.Equal(t, uint64(0), stat.SndbufPeak)
	require.Len(t, stat.CPUStatList, 3)
	require.Len(t, stat.SockStatList, 2)
	for index, cpuStat := range stat.CPUStatList {
		testCPUStat(t, index, cpuStat)
	}

	for index, sockStat := range stat.SockStatList {
		testSockStat(t, index, sockStat)
	}
}

func TestReadStatistics(t *testing.T) {
	setReadFileFunc(t, fileContent, nil)
	statCollector := New("test_path")
	stat, err := statCollector.CollectAndMarshal()
	require.NoError(t, err)
	testDefaults(t, stat)
}

func TestReadFileError(t *testing.T) {
	setReadFileFunc(t, fileContent, errors.New("test_error"))
	statCollector := New("test_path")
	_, err := statCollector.CollectAndMarshal()
	require.Error(t, err)
	require.Equal(t, "test_error", err.Error())
}

func TestParseError(t *testing.T) {
	metrics := "inBitRate    1.2"
	setReadFileFunc(t, metrics, nil)
	statCollector := New("test_path")
	_, err := statCollector.CollectAndMarshal()
	require.Error(t, err)
	require.Contains(t, err.Error(), "strconv.ParseUint:")
}

func TestParseFloatError(t *testing.T) {
	metrics := "hashMetric   test"
	setReadFileFunc(t, metrics, nil)
	statCollector := New("test_path")
	_, err := statCollector.CollectAndMarshal()
	require.Error(t, err)
	require.Contains(t, err.Error(), "strconv.ParseFloat")
}

func TestParseErrorCpuMetricCount(t *testing.T) {
	metrics := "cpu0 1 2 3 4 1.35 5 6 7\ncpu1 1 2 3 4 1.35 5 6 7 8 9 10"
	setReadFileFunc(t, metrics, nil)
	statCollector := New("test_path")
	stat, err := statCollector.CollectAndMarshal()
	require.NoError(t, err)
	require.Len(t, stat.CPUStatList, 1)
	testCPUStat(t, 1, stat.CPUStatList[0])
}

func TestParseUint32(t *testing.T) {
	_, err := getValueByType(reflect.ValueOf(uint32(0)), "5000000000")
	require.Error(t, err)
	require.Contains(t, err.Error(), "strconv.ParseUint:")
}

func TestParseUnsupportedType(t *testing.T) {
	_, err := getValueByType(reflect.ValueOf(uint8(0)), "1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported type")
}

func TestUnsupportedMetric(t *testing.T) {
	err := setValueByName(&Statistics{}, "not_exist", "123")
	require.Error(t, err)
	require.ErrorIs(t, err, errNotFoundField)
}
