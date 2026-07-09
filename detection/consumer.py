import json
import time
from datetime import datetime
from typing import Optional
from redis.client import get_redis
from db.session import SessionLocal
from db.models import Event
from engine.rules import RuleEngine
from engine.ml import MLDetector
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
                    self._handle_event(db, msg_id, data)

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

        logger.info("event_processed",
                     event_id=event["id"],
                     verdict=verdict["action"])

    def _detect(self, event: dict) -> dict:
        rule_action = self.rule_engine.evaluate(event)
        if rule_action:
            return {"action": rule_action, "score": 1.0}

        ml_score = self.ml_detector.predict(event)
        if ml_score < -0.5:
            return {"action": "WARN", "score": ml_score}

        return {"action": "PASS", "score": ml_score}

    def stop(self):
        self.running = False
        logger.info("consumer_stopped")
