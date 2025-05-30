package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Pacman29/observability/metrics"
	prometheus2 "github.com/Pacman29/observability/metrics/prometheus"
)

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()

	m := metrics.New(prometheus2.NewPrometheusDriver(prometheus.DefaultRegisterer))
	ctx = m.WithTags(ctx, map[string]string{
		"app": "test_app",
	})

	m.Counter(m.WithTag(ctx, "label", "key"), "test_metric", 1)

	mux.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics:                   true,
			EnableOpenMetricsTextCreatedSamples: true,
		}))

	mux.HandleFunc("/count", func(writer http.ResponseWriter, request *http.Request) {
		m.Counter(request.Context(), "test_metric", 10)
	})

	server := &http.Server{Addr: "0.0.0.0:8080", Handler: mux}

	if err := server.ListenAndServe(); err != nil {
		slog.Error(err.Error())
	}
}
