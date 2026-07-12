import uvicorn
import threading
import time
from fastapi import FastAPI
from fastapi.responses import PlainTextResponse
from contextlib import asynccontextmanager
from prometheus_client import generate_latest, CONTENT_TYPE_LATEST
from consumer import DetectionConsumer
from db.session import engine, SessionLocal
from db.models import Base, Event
from engine.ml import MLDetector
from config import settings
from metrics import REDIS_CONNECTED, DB_CONNECTED
import structlog

logger = structlog.get_logger()

consumer = DetectionConsumer()


def retrain_loop():
    while consumer.running:
        time.sleep(settings.ML_RETRAIN_INTERVAL_SEC)
        if not consumer.running:
            break
        consumer.retrain_ml()
        logger.info("ml_retrained")


@asynccontextmanager
async def lifespan(app: FastAPI):
    Base.metadata.create_all(bind=engine)

    consumer_thread = threading.Thread(target=consumer.start, daemon=True)
    consumer_thread.start()

    retrain_thread = threading.Thread(target=retrain_loop, daemon=True)
    retrain_thread.start()

    yield

    consumer.stop()


app = FastAPI(title="ADR-P Detection Engine", lifespan=lifespan)


@app.get("/health")
def health():
    checks = {}

    try:
        from redis.client import get_redis
        r = get_redis()
        r.ping()
        REDIS_CONNECTED.set(1)
        checks["redis"] = "connected"
    except Exception:
        REDIS_CONNECTED.set(0)
        checks["redis"] = "disconnected"

    try:
        with engine.connect() as conn:
            conn.execute(engine.dialect.statement_compiler(engine.dialect, None))
        DB_CONNECTED.set(1)
        checks["postgres"] = "connected"
    except Exception:
        DB_CONNECTED.set(0)
        checks["postgres"] = "disconnected"

    status = "healthy" if all(v == "connected" for v in checks.values()) else "degraded"
    return {"status": status, **checks}


@app.get("/metrics")
def metrics():
    return PlainTextResponse(
        generate_latest(),
        media_type=CONTENT_TYPE_LATEST,
    )


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8001)
