# Docker Compose - Calendar Booking Application

Полный стек приложения календаря бронирований в Docker Compose.

## Компоненты

1. **PostgreSQL** - База данных (порт 5432)
2. **Go Backend** - API сервер (порт 8080)
3. **React Frontend** - Веб-интерфейс (порт 80)

## Быстрый старт

### Запуск всех сервисов

```bash
docker-compose up -d
```

### Проверка статуса

```bash
docker-compose ps
```

### Просмотр логов

```bash
# Все логи
docker-compose logs -f

# Только backend
docker-compose logs -f backend

# Только база данных
docker-compose logs -f postgres

# Только фронтенд
docker-compose logs -f frontend
```

### Остановка

```bash
docker-compose down
```

### Остановка с удалением данных

```bash
docker-compose down -v
```

## Доступные адреса

- **Фронтенд**: http://localhost:80
- **Бэкенд API**: http://localhost:8080
- **База данных**: localhost:5432

## Проверка работоспособности

```bash
# Backend health check
curl http://localhost:8080/health

# Frontend health check
curl http://localhost:80/health
```

## Конфигурация

### Backend (переменные окружения)

```env
SERVER_PORT=8080
APP_ENV=production
DATABASE_URL=postgres://postgres:postgres@postgres:5432/booking_db?sslmode=disable
DB_MAX_CONNS=10
DB_MIN_CONNS=2
DB_MAX_CONN_LIFETIME_HOURS=1
DB_MAX_CONN_IDLE_TIME_MIN=30
```

### Frontend (build args)

```env
VITE_API_BASE_URL=http://localhost:8080
```

### Database

```env
POSTGRES_DB=booking_db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
```

## Архитектура

```
                  ┌─────────────┐
                  │   Фронтенд   │
                  │   (Nginx)    │
                  │  Порт 80     │
                  └──────┬───────┘
                         │
                         │ HTTP
                         │
                  ┌──────▼───────┐
                  │   Бэкенд     │
                  │    (Go)      │
                  │  Порт 8080   │
                  └──────┬───────┘
                         │
                         │ SQL
                         │
                  ┌──────▼───────┐
                  │ PostgreSQL   │
                  │  Порт 5432   │
                  └──────────────┘
```

## Миграции базы данных

Миграции автоматически применяются при первом запуске PostgreSQL из директории:
- `backend/migrations/`

## Сеть

Все сервисы находятся в одной Docker сети `booking_network` и могут общаться друг с другом по именам сервисов:
- Фронтенд → Бэкенд: `http://backend:8080`
- Бэкенд → БД: `postgres:5432`

## Тома

- `postgres_data` - Постоянное хранилище данных PostgreSQL

## Пересборка

```bash
# Пересобрать все образы
docker-compose build --no-cache

# Пересобрать только бэкенд
docker-compose build backend

# Пересобрать только фронтенд
docker-compose build frontend
```

## Разработка

Для разработки можно запускать сервисы отдельно:

```bash
# Только база данных
docker-compose up -d postgres

# База + бэкенд
docker-compose up -d postgres backend

# Все сервисы
docker-compose up -d
```

## Устранение неполадок

### База данных не запускается

```bash
docker-compose logs postgres
```

### Бэкенд не может подключиться к БД

```bash
# Проверить readiness базы
docker-compose exec postgres pg_isready -U postgres

# Перезапустить бэкенд
docker-compose restart backend
```

### Фронтенд не может подключиться к бэкенду

Убедитесь, что VITE_API_BASE_URL правильно настроен в docker-compose.yml и пересоберите фронтенд:

```bash
docker-compose build frontend
docker-compose up -d frontend
```

## Production использование

Для продакшена рекомендуется:

1. Изменить пароли базы данных
2. Использовать secrets для чувствительных данных
3. Настроить reverse proxy (nginx/traefik)
4. Добавить SSL/TLS сертификаты
5. Настроить мониторинг и алерты
6. Настроить резервное копирование БД

---

**Версия**: 1.0
**Дата**: 6 апреля 2026
