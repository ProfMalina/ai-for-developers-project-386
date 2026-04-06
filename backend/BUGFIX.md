# Исправление: "Не удалось сохранить тип встречи"

## Проблема

При попытке создать тип события (event type) бэкенд возвращал ошибку, так как использовал захардкоженный `ownerID = "00000000-0000-0000-0000-000000000001"`, которого не существовало в базе данных. Из-за foreign key constraint PostgreSQL отклонял запрос.

## Причина

- Бэкенд использовал фиксированный `ownerID` для всех запросов
- Этот владелец не создавался автоматически при инициализации базы данных
- Foreign key constraint `event_types.owner_id -> owners.id` блокировал вставку

## Решение

### 1. Добавлен владелец по умолчанию в миграцию

**Файл**: `backend/migrations/001_initial_schema.up.sql`

```sql
-- Insert default owner for development
INSERT INTO owners (id, name, email, timezone, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Default Owner',
    'owner@example.com',
    'Europe/Moscow',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;
```

### 2. Добавлена функция SeedDefaultOwner

**Файл**: `backend/internal/db/database.go`

- Проверяет существование владельца по умолчанию
- Создаёт его, если не существует
- Устанавливает `db.DefaultOwnerID` для использования в handlers

### 3. Обновлены handlers

**Файлы**:
- `backend/internal/handlers/event_type_handler.go`
- `backend/internal/handlers/time_slot_handler.go`
- `backend/internal/handlers/public_event_type_handler.go`

Заменён захардкоженный ID на `db.DefaultOwnerID`:

```go
ownerID := c.GetString("ownerID")
if ownerID == "" {
    ownerID = db.DefaultOwnerID  // Вместо "00000000-0000-0000-0000-000000000001"
}
```

### 4. Вызов seed при старте

**Файл**: `backend/cmd/server/main.go`

```go
// Seed default owner for development
if err := db.SeedDefaultOwner(ctx); err != nil {
    log.Printf("Warning: Failed to seed default owner: %v", err)
}
```

## Как проверить

### Вариант 1: Docker Compose (рекомендуется)

```bash
# Остановить и пересобрать
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d

# Проверить создание типа события
curl -X POST http://localhost:8080/api/event-types \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Консультация",
    "description": "Индивидуальная консультация",
    "durationMinutes": 30
  }'
```

### Вариант 2: Локальный запуск

```bash
# 1. Запустить БД
make backend-db-up

# 2. Пересобрать и запустить бэкенд
cd backend
go build -o server ./cmd/server
./server

# 3. В другом терминале запустить тесты
./backend/scripts/test-api.sh
```

## Ожидаемый результат

```json
{
  "id": "uuid-generated",
  "ownerId": "00000000-0000-0000-0000-000000000001",
  "name": "Консультация",
  "description": "Индивидуальная консультация",
  "durationMinutes": 30,
  "isActive": true,
  "createdAt": "2026-04-06T...",
  "updatedAt": "2026-04-06T..."
}
```

## Статус

✅ **ИСПРАВЛЕНО** - Типы событий теперь создаются успешно

## Дата

6 апреля 2026
