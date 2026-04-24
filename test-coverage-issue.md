## Test Coverage Gaps - Prioritized

**Методология приоритизации:**
1. Количество зависимостей (dependents) - вызывающие модули
2. Частота недавних изменений
3. Сложность модуля

---

### Высокий приоритет

| # | Модуль/Функция | Файлы | Зависимости | Сложность | Текущий статус |
|---|----------------|-------|------------|-----------|----------------|
| 1 | `TimeSlotHandler.List` | backend/internal/handlers/time_slot_handler.go:40 | Сайт, гостевое бронирование | Высокая (пагинация, фильтры) | Только happy path |
| 2 | `TimeSlotHandler.GenerateSlots` | backend/internal/handlers/time_slot_handler.go:104 | Генерация слотов | Высокая | Нет тестов |
| 3 | `EventTypeHandler.Delete` | backend/internal/handlers/event_type_handler.go:110 | CRUD | Средняя | Нет тестов |
| 4 | `helpers.go` - `validationMessage` | backend/internal/handlers/helpers.go:71 | Все хендлеры | Низкая | Нет тестов |

---

### Средний приоритет

| # | Модуль/Функция | Файлы | Зависимости | Сложность | Текущий статус |
|---|----------------|-------|------------|-----------|----------------|
| 5 | `Conflict` helper | backend/internal/handlers/helpers.go:51 | PublicBookingHandler | Средняя | Нет тестов |
| 6 | `InvalidTime` helper | backend/internal/handlers/helpers.go:46 | PublicBookingHandler | Средняя | Нет тестов |
| 7 | `PublicEventTypeHandler.List` | backend/internal/handlers/public_event_type_handler.go:41 | Гостевое API | Средняя | Нет тестов |
| 8 | `EventTypeHandler.List` | backend/internal/handlers/event_type_handler.go:66 | Owner API | Средняя | Нет тестов |
| 9 | `EventTypeHandler.GetByID` | backend/internal/handlers/event_type_handler.go:53 | Owner API | Средняя | Нет тестов |
| 10 | `BookingHandler.Cancel` edge cases | backend/internal/handlers/booking_handler.go:83 | Owner API | Средняя | Не полностью |

---

### Низкий приоритет

| # | Модуль/Функция | Файлы | Зависимости | Сложность | Текущий статус |
|---|----------------|-------|------------|-----------|----------------|
| 11 | `lowerFirst` helper | backend/internal/handlers/helpers.go:64 | Валидация | Низкая | Нет тестов |
| 12 | `errorAsValidation` helper | backend/internal/handlers/helpers.go:60 | Валидация | Низкая | Нет тестов |
| 13 | `middleware.ErrorHandler` error types | backend/internal/middleware/middleware.go:11 | Все ошибки | Низкая | Не полностью |
| 14 | Валидация времени (dateFrom/dateTo parsing) | time_slot_handler.go, booking_handler.go | Пагинация | Средняя | Недостаточно |

---

### Frontend (низкий приоритет, CI уже покрывает)

| # | Модуль/Функция | Зависимости | Статус |
|---|----------------|------------|--------|
| 15 | `api/client.ts` error interceptors | Все API вызовы | Mini tests exist |
| 16 | `utils/validation.ts` | Формы | Mini tests exist |

---

### Рекомендуемые действия

1. **Срочно**: Добавить тесты для `TimeSlotHandler.List` и `GenerateSlots` - критичный функционал
2. **Срочно**: покрыть helper функции `validationMessage`, `Conflict`, `InvalidTime` 
3. Все `EventTypeHandler` методы кроме Create/Update
4. Edge cases для `BookingHandler.Cancel`

---
*Пояснение: приоритет основан на количествеdependents (вызовов из других модулей) и бизнес-критичности*