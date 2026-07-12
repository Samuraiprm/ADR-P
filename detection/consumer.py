import json
import time
from datetime import datetime
from typing import Optional
from redis.client import get_redis
from db.session import SessionLocal
from db.models import Event
from engine.rules import RuleEngine
from engine.ml import MLDetector
from metrics import (
    EVENTS_CONSUMED, EVENTS_PROCESSED, EVENTS_FAILED,
    RULE_MATCHES, ML_SCORE, ML_TRAINED, ML_SAMPLES,
)
import structlog

logger = structlog.get_logger()

STREAM_NAME = "abuse:events:raw"
VERDICT_STREAM = "abuse:verdicts"
GROUP_NAME = "detection"
CONSUMER_NAME = "detector-1"


class DetectionConsumer:
    def __init__(self):
        self.redis = get_redis()
        self.rule_engine = RuleEngine()
        self.ml_detector = MLDetector()
        self.running = False

    def start(self):
        self.running = True
        logger.info("consumer_started")

        try:
            self.redis.xgroup_create(STREAM_NAME, GROUP_NAME, id="0", mkstream=True)
        except Exception:
            pass

        while self.running:
            try:
                self._process_batch()
            except Exception as e:
                logger.error("processing_error", error=str(e))
                time.sleep(1)

    def _process_batch(self):
        results = self.redis.xreadgroup(
            GROUP_NAME,
            CONSUMER_NAME,
            {STREAM_NAME: ">"},
            count=100,
            block=1000,
        )

        if not results:
            return

        with SessionLocal() as db:
            self.rule_engine.load_rules(db)

            for stream, messages in results:
                for msg_id, data in messages:
                    EVENTS_CONSUMED.inc()
                    try:
                        self._handle_event(db, msg_id, data)
                    except Exception as e:
                        EVENTS_FAILED.inc()
                        logger.error("event_handle_error", msg_id=msg_id, error=str(e))
                        try:
                            self.redis.xack(STREAM_NAME, GROUP_NAME, msg_id)
                        except Exception:
                            pass

    def _handle_event(self, db, msg_id: str, data: dict):
        event = {
            "id": data.get("id"),
            "user_id": data.get("user_meta", {}).get("user_id", "unknown"),
            "event_type": data.get("event_type"),
            "source": data.get("source"),
            "payload": json.loads(data.get("payload", "{}")) if isinstance(data.get("payload"), str) else data.get("payload", {}),
        }

        verdict = self._detect(event)

        db_event = Event(
            id=event["id"],
            user_id=event["user_id"],
            event_type=event["event_type"],
            timestamp=datetime.now(),
            verdict=verdict["action"],
            score=verdict["score"],
            matched_rule_id=verdict.get("rule_id"),
        )
        db.add(db_event)
        db.commit()

        self.redis.xadd(VERDICT_STREAM, {
            "event_id": event["id"],
            "verdict": verdict["action"],
            "score": str(verdict["score"]),
        })

        self.redis.xack(STREAM_NAME, GROUP_NAME, msg_id)

        EVENTS_PROCESSED.labels(verdict=verdict["action"]).inc()
        ML_SCORE.observe(verdict["score"])

        logger.info("event_processed",
                     event_id=event["id"],
                     verdict=verdict["action"])

    def _detect(self, event: dict) -> dict:
        rule_action, rule_name = self.rule_engine.evaluate(event)
        if rule_action:
            RULE_MATCHES.labels(rule_name=rule_name or "unknown", action=rule_action).inc()
            return {"action": rule_action, "score": 1.0}

        ml_score = self.ml_detector.predict(event)
        if ml_score < -0.5:
            return {"action": "WARN", "score": ml_score}

        return {"action": "PASS", "score": ml_score}

    def retrain_ml(self):
        try:
            from db.session import SessionLocal as SL
            from db.models import Event as Ev
            from config import settings

            with SL() as db:
                rows = db.query(Ev).order_by(Ev.timestamp.desc()).limit(settings.ML_WINDOW_SIZE).all()
                events = [
                    {
                        "user_id": r.user_id,
                        "event_type": r.event_type,
                        "payload": {},
                    }
                    for r in rows
                ]
                self.ml_detector.train(events)
                ML_TRAINED.set(1 if self.ml_detector.is_trained else 0)
                ML_SAMPLES.set(len(events))
        except Exception as e:
            logger.error("retrain_error", error=str(e))

    def stop(self):
        self.running = False
        logger.info("consumer_stopped")
