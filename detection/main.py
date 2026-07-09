import uvicorn
from fastapi import FastAPI
from contextlib import asynccontextmanager
from consumer import DetectionConsumer
from db.session import engine
from db.models import Base
import threading
import structlog

logger = structlog.get_logger()

consumer = DetectionConsumer()


@asynccontextmanager
async def lifespan(app: FastAPI):
    Base.metadata.create_all(bind=engine)
    thread = threading.Thread(target=consumer.start, daemon=True)
    thread.start()
    yield
    consumer.stop()


app = FastAPI(title="ADR-P Detection Engine", lifespan=lifespan)


@app.get("/health")
def health():
    return {"status": "healthy"}


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8001)
