from prometheus_client import Counter, Histogram, Gauge

EVENTS_CONSUMED = Counter(
    "adr_events_consumed_total",
    "Total events consumed from Redis Stream",
)

EVENTS_PROCESSED = Counter(
    "adr_events_processed_total",
    "Total events processed successfully",
    ["verdict"],
)

EVENTS_FAILED = Counter(
    "adr_events_failed_total",
    "Total events that failed processing",
)

RULE_MATCHES = Counter(
    "adr_rule_triggers_total",
    "Total rule matches",
    ["rule_name", "action"],
)

ML_SCORE = Histogram(
    "adr_anomaly_score",
    "ML anomaly score distribution",
    buckets=[-1.0, -0.8, -0.6, -0.5, -0.4, -0.2, 0.0, 0.2],
)

ML_TRAINED = Gauge(
    "adr_ml_model_trained",
    "Whether ML model is trained (1=yes, 0=no)",
)

ML_SAMPLES = Gauge(
    "adr_ml_training_samples",
    "Number of samples used in last ML training",
)

REDIS_CONNECTED = Gauge(
    "adr_redis_connected",
    "Redis connection status (1=connected, 0=disconnected)",
)

DB_CONNECTED = Gauge(
    "adr_db_connected",
    "PostgreSQL connection status (1=connected, 0=disconnected)",
)
