## Анализ пробелов тестового покрытия

### Методология оценки
- **Зависимости**: число других модулей, вызывающих этот
- **Недавние изменения**: частота коммитов (git log)
- **Сложность**: cyclomatic complexity, edge cases, внешние зависимости

---

## КРИТИЧЕСКИЕ (без единого теста)

### 1. `backend/internal/handlers/time_slot_handler.go` — Приоритет: 10/10

**Функции без тестов:** `List()`, `GenerateSlots()`

**Почему критично:**
- `GenerateSlots()` — 256 строк бизнес-логики: timezone parsing, date validation, slot generation loop
- Cyclomatic complexity 8+, 4 внешние зависимости
- Endpoint для создания слотов — основа приложения

---

### 2. `backend/internal/handlers/public_event_type_handler.go` — Приоритет: 9/10

**Функции без тестов:** `List()`, `GetByID()`, `GetSlots()`

**Почему критично:**
- Публичный API для гостей
- Guest flow — основа приложения
- Связность с TimeSlotService.GetAvailableSlots()

---

### 3. `backend/internal/memory/*.go` — Приоритет: 9/10

**Файлы без тестов:**
- `memory/owner_repository.go`
- `memory/time_slot_repository.go`
- `memory/booking_repository.go`
- `memory/event_type_repository.go`
- `memory/store.go`

**Почему критично:**
- Fallback при недоступности PostgreSQL
- Thread safety, in-memory state management
- Используется в `container.go:66-96` как production fallback

---

### 4. `backend/internal/repositories/slot_config_repository.go` — Приоритет: 8/10

**Функции без тестов:** `Create()`, `GetByOwnerID()`, `Delete()`

**Почему важно:**
- Сохраняет конфигурацию генерации слотов
- Зависимость от TimeSlotService.GenerateSlots()

---

## ВЫСОКИЙ ПРИОРИТЕТ (частичное покрытие)

### 5. `backend/internal/services/time_slot_service.go` — Приоритет: 8/10

**Пробелы в существующем покрытии:**
| Функция | Текущее покрытие | Пробелы |
|---------|-----------------|---------|
| `GenerateSlots()` | Минимальное | Invalid timezone, date validation edge cases |
| `GetAvailableSlots()` | Нет unit тестов | Filter logic after DB query |
| `List()` | Частично | Date validation edge cases |

---

### 6. `frontend/src/utils/validation.ts` — Приоритет: 7/10

**Почему важно:** Валидация форм, используется в EventTypeManagement, BookingPage

---

### 7. `backend/internal/handlers/helpers.go` — Приоритет: 6/10

**Функции без тестов:** `lowerFirst()`, `validationMessage()` — критичны для error handling

---

## СРЕДНИЙ ПРИОРИТЕТ

### 8. `frontend/src/api/client.ts` — Приоритет: 6/10

**Без тестов:**
- `ownerApi.updateEventType`
- `ownerApi.deleteEventType`
- `ownerApi.cancelBooking`
- `ownerApi.generateSlots`

---

### 9. `backend/internal/handlers/booking_handler.go` — Приоритет: 5/10

**Пробел:** `Delete()` — owner может удалять букинги

---

## РЕЗЮМЕ СТАТУСА

| Область | Покрытие | Критические пропуски |
|---------|----------|---------------------|
| Backend | ~65% | 4 модуля без тестов |
| Frontend | ~75% | Основные потоки покрыты |
| E2E | Хорошее | — |

---

## РЕКОМЕНДУЕМЫЙ ПЛАН

1. **Немедленно**: `time_slot_handler.go` — критический endpoint
2. **Немедленно**: `public_event_type_handler.go` — guest flow
3. **Срочно**: Memory repositories — fallback storage
4. **Срочно**: Усилить `time_slot_service_test.go`
5. **Средний**: Frontend validation, API client edge cases