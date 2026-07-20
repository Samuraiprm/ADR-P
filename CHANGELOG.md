# Changelog

All notable changes to ADR-P project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2026-07-20

### Added
- Prometheus metrics for response service (9 metrics: verdicts, telegram, rules, DB, health)
- TelegramCallback route registered (`POST /api/v1/telegram/callback`)
- `updated_at` audit trail column on events table
- DBOperations metrics on all database queries
- RedisHealth/PostgresHealth gauges in `/healthz` endpoint
- README with architecture diagram, API examples, and project structure
- CI/CD: added response test job, bumped Go version to 1.24

### Fixed
- Verdict string mismatch: BLOCKED→BLOCK, WARNED→WARN (synced with detection service)
- Nginx redundant `/api/v1/telegram/` location removed
- Go 1.22→1.24 in both services (required for `t.Context()`)
- Errcheck warnings in ingestion and response services

## [0.2.0] - 2026-07-12

### Added
- Nginx reverse proxy
- GitHub Actions CI/CD
- Detection Prometheus metrics
- ML retraining loop
- Seed detection rules (rate_limit, burst_detection, spam_keywords)
- Dashboard rewritten with real metrics

### Fixed
- Detection asyncpg driver → sync postgresql
- Logger Unicode rune status → strconv.Itoa
- Detection consumer per-message error handling

## [0.1.0] - 2026-07-10

### Added
- Go Ingestion Service (Stage 1)
- Python Detection Service (Stage 2)
- Go Response Service (Stage 3)
- Telegram Bot integration with inline keyboards
- Grafana dashboard (11 panels)
- Prometheus scrape configuration
- Docker Compose orchestration (7 services)
