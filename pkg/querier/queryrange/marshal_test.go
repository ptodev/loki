package queryrange

import (
	"testing"

	"github.com/prometheus/prometheus/promql"
	"github.com/stretchr/testify/require"

	"github.com/grafana/loki/pkg/loghttp"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/logql"
	"github.com/grafana/loki/pkg/logqlmodel"
	"github.com/grafana/loki/pkg/querier/queryrange/queryrangebase"
)

func TestResultToResponse(t *testing.T) {
	tests := []struct {
		name     string
		result   logqlmodel.Result
		response queryrangebase.Response
	}{
		{
			name: "nil matrix",
			result: logqlmodel.Result{
				Data: promql.Matrix(nil),
			},
			response: &LokiPromResponse{
				Response: &queryrangebase.PrometheusResponse{
					Status: "success",
					Data: queryrangebase.PrometheusData{
						ResultType: loghttp.ResultTypeMatrix,
						Result:     []queryrangebase.SampleStream{},
					},
				},
			},
		},
		{
			name: "empty probabilistic quantile matrix",
			result: logqlmodel.Result{
				Data: logql.ProbabilisticQuantileMatrix([]logql.ProbabilisticQuantileVector{}),
			},
			response: &QuantileSketchResponse{
				Response: &logproto.QuantileSketchMatrix{
					Values: []*logproto.QuantileSketchVector{},
				},
				Headers: []queryrangebase.PrometheusResponseHeader(nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ResultToResponse(tt.result, nil)
			require.NoError(t, err)

			require.Equal(t, tt.response, actual)
		})
	}
}

func TestResponseWrap(t *testing.T) {
	for _, tt := range []struct {
		name     string
		response queryrangebase.Response
		expected isQueryResponse_Response
	}{
		{"volume", &VolumeResponse{}, &QueryResponse_Volume{}},
		{"series", &LokiSeriesResponse{}, &QueryResponse_Series{}},
		{"label", &LokiLabelNamesResponse{}, &QueryResponse_Labels{}},
		{"stats", &IndexStatsResponse{}, &QueryResponse_Stats{}},
		{"prom", &LokiPromResponse{}, &QueryResponse_Prom{}},
		{"streams", &LokiResponse{}, &QueryResponse_Streams{}},
		{"topk", &TopKSketchesResponse{}, &QueryResponse_TopkSketches{}},
		{"quantile", &QuantileSketchResponse{}, &QueryResponse_QuantileSketches{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := QueryResponseWrap(tt.response)
			require.NoError(t, err)
			require.IsType(t, tt.expected, actual.Response)
		})
	}
}
