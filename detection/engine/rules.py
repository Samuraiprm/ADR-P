from datetime import datetime
from typing import List
from db.models import DetectionRule
from redis.client import get_redis
import json


class RuleEngine:
    def __init__(self):
        self.rules: List[DetectionRule] = []

    def load_rules(self, db_session):
        self.rules = db_session.query(DetectionRule).filter(
            DetectionRule.is_active == True
        ).all()

    def evaluate(self, event: dict) -> str | None:
        for rule in self.rules:
            if self._matches_rule(event, rule):
                return rule.action
        return None

    def _matches_rule(self, event: dict, rule: DetectionRule) -> bool:
        condition = rule.condition_json

        if rule.name == "rate_limit":
            return self._check_rate_limit(event, condition)

        return False

    def _check_rate_limit(self, event: dict, condition: dict) -> bool:
        window_sec = condition.get("window_sec", 60)
        threshold = condition.get("threshold", 10)

        r = get_redis()
        key = f"rate:{event.get('user_id')}:{event.get('event_type')}"

        now = datetime.now().timestamp()
        pipe = r.pipeline()
        pipe.zremrangebyscore(key, 0, now - window_sec)
        pipe.zadd(key, {str(now): now})
        pipe.zcard(key)
        pipe.expire(key, window_sec)
        results = pipe.execute()

        count = results[2]
        return count > threshold
