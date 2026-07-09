import numpy as np
from sklearn.ensemble import IsolationForest
from typing import List, Optional
import structlog

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
            n_jobs=-1,
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
            return np.array([
                float(payload.get("message_length", 0)),
                float(payload.get("word_count", 0)),
                float(payload.get("url_count", 0)),
                float(payload.get("caps_ratio", 0)),
                float(payload.get("hour_of_day", 0)),
            ])
        except (ValueError, TypeError):
            return None
