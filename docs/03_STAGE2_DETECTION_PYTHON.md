# Этап 2: Движок детекции (Python)

## Цель
Реализовать гибридную систему детекции: правила (SQL/Python) + ML-аномалии.

## Функциональные требования
1. **Consumer:** Чтение из Redis Stream `abuse:events:raw` с consumer group.
2. **Rule Engine:**
   - Загрузка правил из БД (таблица `detection_rules`).
   - Применение правил к батчу событий.
   - Пример правила: "N событий от одного user_id за T секунд".
3. **ML Detector:**
   - Периодическое (раз в 5 мин) обучение/инференс модели Isolation Forest на последних N событиях.
   - Скоринг аномальности для событий, не попавших под правила.
4. **Enrichment:** Обогащение событий контекстом из БД (история жалоб, возраст аккаунта).
5. **Output:** Запись вердикта (`PASS`, `WARN`, `BLOCK`) в Redis Stream `abuse:verdicts`.

## Модель данных (PostgreSQL)
```sql
CREATE TABLE events (
    id UUID PRIMARY KEY,
    user_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    verdict TEXT,
    score FLOAT,
    matched_rule_id INT
);

CREATE TABLE detection_rules (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    condition_json JSONB NOT NULL, -- e.g. {"window_sec": 60, "threshold": 10}
    action TEXT NOT NULL, -- BLOCK, WARN
    is_active BOOLEAN DEFAULT TRUE
);
```
