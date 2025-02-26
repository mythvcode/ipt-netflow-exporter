package exporter

import (
	"strings"
	"testing"

	"github.com/mythvcode/ipt-netflow-exporter/internal/config"
	"github.com/mythvcode/ipt-netflow-exporter/internal/exporter/mocks"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestNewExporter(t *testing.T) {
	mock := mocks.NewMockStatParser(t)
	cfg, err := config.ReadConfig("")
	require.NoError(t, err)
	_, err = New(cfg.Exporter, mock)
	require.NoError(t, err)
}

func TestNewGetStats(t *testing.T) {
	mock := mocks.NewMockStatParser(t)
	collector := newIPTNetFlowTCollector(mock)
	mock.EXPECT().CollectAndMarshal().Return(getTestStatistic(t), nil)
	err := testutil.CollectAndCompare(collector, strings.NewReader(getPromTestStat(t)))
	require.NoError(t, err)
}
