package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	ReceivedEvents = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "adr_received_events_total",
			Help: "Total number of received events",
		},
	)

	ValidatedEvents = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "adr_validated_events_total",
			Help: "Total number of validated events",
		},
	)

	DroppedEvents = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "adr_dropped_events_total",
			Help: "Total number of dropped events",
		},
	)

	HTTPRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "adr_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "status"},
	)

	HTTPDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "adr_http_duration_seconds",
			Help:    "HTTP request duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(
		ReceivedEvents,
		ValidatedEvents,
		DroppedEvents,
		HTTPRequests,
		HTTPDuration,
	)
}
