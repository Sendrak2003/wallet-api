# Wallet API

REST API для управления кошельками с поддержкой высокой нагрузки.

## Стек

- Go 1.24
- PostgreSQL 15
- Gin Framework
- Docker

## Запуск

```bash
docker-compose up --build
```

## Тесты

```bash
go test ./test/... -v
```

## Остановка

```bash
docker-compose down -v
```


## API

**POST /api/v1/wallet** - Применить операцию
```json
{
  "walletId": "11111111-1111-1111-1111-111111111111",
  "operationType": "DEPOSIT",
  "amount": 1000
}
```

**GET /api/v1/wallets/{id}** - Получить баланс

**GET /health** - Health check

**GET /swagger/index.html** - Swagger UI

## Архитектура

### Конкурентность
- SELECT FOR UPDATE для блокировки строк
- Connection Pool (100 соединений)
- Транзакции для атомарности

### Отказоустойчивость
- Graceful Shutdown (30 сек на завершение запросов)
- Health Check для мониторинга
- Timeouts (Read/Write/Idle)
- Proper error handling (4xx/5xx)

### Структура
```
cmd/server/          - Точка входа
internal/
  handler/http/      - HTTP handlers
  service/           - Бизнес-логика
  repository/        - Работа с БД
migrations/          - SQL миграции
```
