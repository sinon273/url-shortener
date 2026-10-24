package url_metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	grpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)
	grpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	LinksCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "links_created_total",
			Help: "Total number of created short links",
		},
	)

	RedirectsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "redirects_total",
			Help: "Total number of redirects",
		},
	)

	KafkaMessagesProduced = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_produced_total",
			Help: "Total number of messages produced to Kafka",
		},
	)
)

func init() {
	prometheus.MustRegister(grpcRequestsTotal)
	prometheus.MustRegister(grpcRequestDuration)
	prometheus.MustRegister(LinksCreatedTotal)
	prometheus.MustRegister(RedirectsTotal)
	prometheus.MustRegister(KafkaMessagesProduced)
}
