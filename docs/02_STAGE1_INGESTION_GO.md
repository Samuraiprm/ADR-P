# Этап 1: Сервис приёма событий (Go)

## Цель
Создать отказоустойчивый сервис приёма вебхуков, который валидирует входящие данные и складывает их в очередь Redis.

## Функциональные требования
1. **Endpoint:** `POST /api/v1/events`
   - Принимает JSON с полями: `source`, `event_type`, `timestamp`, `user_meta`, `payload`.
   - Валидация схемы через `go-playground/validator`.
2. **Аутентификация:** Проверка HMAC-SHA256 подписи в заголовке `X-Webhook-Signature`.
3. **Буферизация:** Запись валидных событий в Redis Stream `abuse:events:raw`.
4. **Метрики:** Экспорт `/metrics` для Prometheus (counters: received, validated, dropped, queue_size).
5. **Healthcheck:** `GET /healthz` (проверка Redis + DB).

## Нефункциональные требования
- RPS целевой: 1000+ на одном инстансе.
- Graceful shutdown: ожидание завершения записи в Redis при SIGTERM.
- Структурированные логи (JSON) с trace_id.

## Acceptance Criteria
- [ ] При невалидной подписи возвращает 401, событие не попадает в Redis.
- [ ] При недоступности Redis сервис возвращает 503 и пишет ошибку в лог.
- [ ] Нагрузка в 500 RPS не вызывает рост памяти > 200MB.
- [ ] Есть unit-тесты на валидацию и интеграционный тест с Testcontainers.
