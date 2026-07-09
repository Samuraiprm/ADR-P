# ADR-P - Abuse Detection & Response Platform

Microservice platform for abuse detection and automated response.

## Architecture

- **Ingestion Service** (Go) - Webhook receiver and event queue
- **Detection Engine** (Python) - ML-based anomaly detection
- **Response Service** (Go) - Automated response actions
- **Grafana Dashboard** - Real-time monitoring

## Quick Start

```bash
# Start services
docker-compose up -d

# Test webhook
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{"source":"telegram","event_type":"message","timestamp":1234567890,"payload":{"text":"test"}}'
```

## Documentation

- [Architecture & Stack](docs/01_ARCHITECTURE_AND_STACK.md)
- [Stage 1: Ingestion](docs/02_STAGE1_INGESTION_GO.md)
- [Stage 2: Detection](docs/03_STAGE2_DETECTION_PYTHON.md)
- [Stage 3: Response](docs/04_STAGE3_RESPONSE_AND_VISUALIZATION.md)

## License

MIT
