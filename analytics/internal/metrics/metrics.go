package analytics_metrics

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

	KafkaMessagesConsumedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_consumed_total",
			Help: "Total number of messages consumed from Kafka",
		},
	)

	ClickhouseBatchInsertDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "clickhouse_batch_insert_duration_seconds",
			Help:    "Duration of ClickHouse batch inserts in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(grpcRequestsTotal)
	prometheus.MustRegister(grpcRequestDuration)
	prometheus.MustRegister(KafkaMessagesConsumedTotal)
	prometheus.MustRegister(ClickhouseBatchInsertDuration)
}
