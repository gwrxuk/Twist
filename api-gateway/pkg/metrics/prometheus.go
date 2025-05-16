package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusClient is a client for Prometheus metrics
type PrometheusClient struct {
	requestsTotal       *prometheus.CounterVec
	requestDuration     *prometheus.HistogramVec
	activeConnections   prometheus.Gauge
	databaseConnections *prometheus.GaugeVec
	nodesTotal          prometheus.Gauge
	nodesActive         prometheus.Gauge
}

// NewPrometheusClient creates a new Prometheus metrics client
func NewPrometheusClient() *PrometheusClient {
	// Initialize metrics
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	activeConnections := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_connections",
			Help: "Current number of active HTTP connections",
		},
	)

	databaseConnections := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "database_connections",
			Help: "Current number of database connections",
		},
		[]string{"database"},
	)

	nodesTotal := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "blockchain_nodes_total",
			Help: "Total number of blockchain nodes managed",
		},
	)

	nodesActive := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "blockchain_nodes_active",
			Help: "Number of active blockchain nodes",
		},
	)

	// Register metrics
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(activeConnections)
	prometheus.MustRegister(databaseConnections)
	prometheus.MustRegister(nodesTotal)
	prometheus.MustRegister(nodesActive)

	return &PrometheusClient{
		requestsTotal:       requestsTotal,
		requestDuration:     requestDuration,
		activeConnections:   activeConnections,
		databaseConnections: databaseConnections,
		nodesTotal:          nodesTotal,
		nodesActive:         nodesActive,
	}
}

// RecordRequest records an HTTP request
func (p *PrometheusClient) RecordRequest(path, method string, status int) {
	p.requestsTotal.WithLabelValues(path, method, strconv.Itoa(status)).Inc()
}

// RecordRequestDuration records the duration of an HTTP request
func (p *PrometheusClient) RecordRequestDuration(path, method string, durationSeconds float64) {
	p.requestDuration.WithLabelValues(path, method).Observe(durationSeconds)
}

// SetActiveConnections sets the number of active HTTP connections
func (p *PrometheusClient) SetActiveConnections(count int) {
	p.activeConnections.Set(float64(count))
}

// SetDatabaseConnections sets the number of database connections
func (p *PrometheusClient) SetDatabaseConnections(database string, count int) {
	p.databaseConnections.WithLabelValues(database).Set(float64(count))
}

// SetNodesTotal sets the total number of blockchain nodes
func (p *PrometheusClient) SetNodesTotal(count int) {
	p.nodesTotal.Set(float64(count))
}

// SetNodesActive sets the number of active blockchain nodes
func (p *PrometheusClient) SetNodesActive(count int) {
	p.nodesActive.Set(float64(count))
}

// Handler returns the HTTP handler for Prometheus metrics
func (p *PrometheusClient) Handler() http.Handler {
	return promhttp.Handler()
}
