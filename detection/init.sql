-- ADR-P Database Schema

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    verdict TEXT,
    score FLOAT,
    matched_rule_id INT,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS detection_rules (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    condition_json JSONB NOT NULL,
    action TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_timestamp ON events(timestamp);
CREATE INDEX idx_events_verdict ON events(verdict);

INSERT INTO detection_rules (name, condition_json, action, is_active) VALUES
('rate_limit', '{"window_sec": 60, "threshold": 10}', 'BLOCK', true),
('burst_detection', '{"window_sec": 10, "threshold": 5}', 'WARN', true),
('spam_keywords', '{"keywords": ["spam", "scam", "free money"], "match_count": 1}', 'BLOCK', true)
ON CONFLICT DO NOTHING;
