#!/bin/bash

# Тестирование API бэкенда
# Запускать после запуска бэкенда: make backend-run или docker-compose up

BASE_URL="http://localhost:8080"

echo "======================================"
echo "Тестирование API бэкенда"
echo "======================================"
echo ""

# 1. Проверка health
echo "1. Проверка работоспособности (GET /health)"
curl -s -X GET "$BASE_URL/health" | python3 -m json.tool 2>/dev/null || curl -s -X GET "$BASE_URL/health"
echo ""
echo ""

# 2. Создание типа события
echo "2. Создание типа события (POST /api/event-types)"
CREATE_EVENT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/event-types" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Консультация",
    "description": "Индивидуальная консультация",
    "durationMinutes": 30
  }')

echo "$CREATE_EVENT_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$CREATE_EVENT_RESPONSE"
echo ""

# Извлекаем ID типа события
EVENT_TYPE_ID=$(echo "$CREATE_EVENT_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$EVENT_TYPE_ID" ]; then
  echo "❌ Не удалось создать тип события"
  exit 1
fi

echo "✅ Тип события создан с ID: $EVENT_TYPE_ID"
echo ""

# 3. Получение типа события
echo "3. Получение типа события (GET /api/event-types/$EVENT_TYPE_ID)"
curl -s -X GET "$BASE_URL/api/event-types/$EVENT_TYPE_ID" | python3 -m json.tool 2>/dev/null || curl -s -X GET "$BASE_URL/api/event-types/$EVENT_TYPE_ID"
echo ""
echo ""

# 4. Список типов событий
echo "4. Список типов событий (GET /api/event-types)"
curl -s -X GET "$BASE_URL/api/event-types" | python3 -m json.tool 2>/dev/null || curl -s -X GET "$BASE_URL/api/event-types"
echo ""
echo ""

# 5. Генерация слотов
echo "5. Генерация временных слотов (POST /api/event-types/$EVENT_TYPE_ID/slots/generate)"
GEN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/event-types/$EVENT_TYPE_ID/slots/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "workingHoursStart": "09:00",
    "workingHoursEnd": "17:00",
    "intervalMinutes": 30,
    "daysOfWeek": [1, 2, 3, 4, 5],
    "dateFrom": "2026-04-07",
    "dateTo": "2026-04-10"
  }')

echo "$GEN_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$GEN_RESPONSE"
echo ""

# 6. Список слотов
echo "6. Список слотов (GET /api/slots)"
curl -s -X GET "$BASE_URL/api/slots?eventTypeId=$EVENT_TYPE_ID" | python3 -m json.tool 2>/dev/null || curl -s -X GET "$BASE_URL/api/slots?eventTypeId=$EVENT_TYPE_ID"
echo ""
echo ""

# 7. Публичный список типов событий
echo "7. Публичный список типов событий (GET /api/public/event-types)"
curl -s -X GET "$BASE_URL/api/public/event-types" | python3 -m json.tool 2>/dev/null || curl -s -X GET "$BASE_URL/api/public/event-types"
echo ""
echo ""

# 8. Обновление типа события
echo "8. Обновление типа события (PATCH /api/event-types/$EVENT_TYPE_ID)"
curl -s -X PATCH "$BASE_URL/api/event-types/$EVENT_TYPE_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Консультация (обновлено)",
    "durationMinutes": 45
  }' | python3 -m json.tool 2>/dev/null || curl -s -X PATCH "$BASE_URL/api/event-types/$EVENT_TYPE_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Консультация (обновлено)",
    "durationMinutes": 45
  }'
echo ""
echo ""

echo "======================================"
echo "✅ Тестирование завершено"
echo "======================================"
