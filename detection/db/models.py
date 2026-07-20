import uuid
from datetime import datetime
from sqlalchemy import Column, String, Float, Integer, DateTime, Text, JSON, Boolean
from sqlalchemy.dialects.postgresql import UUID
from db.session import Base


class Event(Base):
    __tablename__ = "events"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(String, nullable=False)
    event_type = Column(String, nullable=False)
    timestamp = Column(DateTime(timezone=True), nullable=False)
    verdict = Column(String)
    score = Column(Float)
    matched_rule_id = Column(Integer)
    updated_at = Column(DateTime(timezone=True), default=datetime.utcnow)


class DetectionRule(Base):
    __tablename__ = "detection_rules"

    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String, nullable=False)
    condition_json = Column(JSON, nullable=False)
    action = Column(String, nullable=False)
    is_active = Column(Boolean, default=True)
