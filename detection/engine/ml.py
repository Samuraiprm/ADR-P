import numpy as np
from sklearn.ensemble import IsolationForest
from typing import List, Optional
import structlog
import re
from datetime import datetime

logger = structlog.get_logger()


class MLDetector:
    def __init__(self):
        self.model: Optional[IsolationForest] = None
        self.is_trained = False

    def train(self, events: List[dict]):
        if len(events) < 50:
            logger.info("insufficient_data_for_training", count=len(events))
            return

        features = self._extract_features(events)
        if features is None or len(features) == 0:
            return

        self.model = IsolationForest(
            n_estimators=100,
            contamination=0.1,
            random_state=42,
            n_jobs=1,
        )
        self.model.fit(features)
        self.is_trained = True
        logger.info("model_trained", samples=len(features))

    def predict(self, event: dict) -> float:
        if not self.is_trained or self.model is None:
            return 0.0

        features = self._extract_single(event)
        if features is None:
            return 0.0

        score = self.model.score_samples([features])[0]
        return float(score)

    def _extract_features(self, events: List[dict]) -> Optional[np.ndarray]:
        try:
            feature_list = []
            for event in events:
                feat = self._extract_single(event)
                if feat is not None:
                    feature_list.append(feat)

            if not feature_list:
                return None

            return np.array(feature_list)
        except Exception as e:
            logger.error("feature_extraction_error", error=str(e))
            return None

    def _extract_single(self, event: dict) -> Optional[np.ndarray]:
        try:
            payload = event.get("payload", {})
            text = payload.get("text", payload.get("message", payload.get("content", "")))
            if not text:
                text = str(payload)

            hour = 0
            ts = event.get("timestamp")
            if isinstance(ts, (int, float)):
                hour = datetime.fromtimestamp(ts).hour
            elif isinstance(ts, str):
                try:
                    hour = datetime.fromisoformat(ts.replace("Z", "+00:00")).hour
                except Exception:
                    hour = 0

            source_map = {"telegram": 0, "vk": 1, "discord": 2, "web": 3}
            source_val = source_map.get(event.get("source", ""), 4)

            type_map = {"message": 0, "comment": 1, "post": 2, "reply": 3}
            type_val = type_map.get(event.get("event_type", ""), 4)

            return np.array([
                float(len(text)),
                float(len(text.split())),
                float(len(re.findall(r'https?://\S+', text))),
                float(sum(1 for c in text if c.isupper()) / max(len(text), 1)),
                float(hour),
                float(source_val),
                float(type_val),
                float(len(set(text.lower().split()))),
            ])
        except (ValueError, TypeError):
            return None
