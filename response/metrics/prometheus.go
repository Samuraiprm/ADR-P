package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	VerdictsConsumed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "adr_verdicts_consumed_total",
			Help: "Total number of verdicts consumed from Redis Stream",
		},
	)

	VerdictsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "adr_verdicts_processed_total",
			Help: "Total number of verdicts processed by action type",
		},
		[]string{"action"},
	)

	VerdictsFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "adr_verdicts_failed_total",
			Help: "Total number of failed verdict processing attempts",
		},
	)

	TelegramMessagesSent = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "adr_telegram_messages_sent_total",
			Help: "Total number of Telegram messages sent",
		},
	)

	TelegramCallbacksHandled = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "adr_telegram_callbacks_handled_total",
			Help: "Total number of Telegram callback queries handled",
		},
	)

	RulesCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "adr_rules_created_total",
			Help: "Total number of detection rules created",
		},
	)

	DBOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "adr_db_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "status"},
	)

	RedisHealth = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "adr_response_redis_healthy",
			Help: "Redis connection health (1 = healthy, 0 = unhealthy)",
		},
	)

	PostgresHealth = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "adr_response_postgres_healthy",
			Help: "PostgreSQL connection health (1 = healthy, 0 = unhealthy)",
		},
	)
)

func init() {
	prometheus.MustRegister(
		VerdictsConsumed,
		VerdictsProcessed,
		VerdictsFailed,
		TelegramMessagesSent,
		TelegramCallbacksHandled,
		RulesCreated,
		DBOperations,
		RedisHealth,
		PostgresHealth,
	)
}
