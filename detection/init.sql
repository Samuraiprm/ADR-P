-- ADR-P Database Schema

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    verdict TEXT,
    score FLOAT,
    matched_rule_id INT
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
