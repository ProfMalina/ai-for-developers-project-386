# Docker Compose - Полный стек приложения

## Что было создано

Создан файл `docker-compose.yml` в корне проекта для запуска полного стека приложения:

### Компоненты

1. **PostgreSQL** (порт 5432)
   - Версия: PostgreSQL 15 Alpine
   - Автоматическое применение миграций при старте
   - Постоянное хранилище через Docker volume
   - Health check для проверки готовности

2. **Go Backend** (порт 8080)
   - Собран из `backend/Dockerfile`
   - Подключается к PostgreSQL
   - Автоматический запуск миграций
   - Health check endpoint `/health`

3. **React Frontend** (порт 80)
   - Собран из `frontend/Dockerfile`
   - Nginx для раздачи статических файлов
   - Настроен на API бэкенда
   - SPA routing

## Как запустить

### Быстрый старт

```bash
# Запустить все сервисы
docker-compose up -d

# Проверить статус
docker-compose ps

# Посмотреть логи
docker-compose logs -f
```

### Остановка

```bash
# Остановить все сервисы
docker-compose down

# Остановить с удалением данных
docker-compose down -v
```

### Makefile команды

```bash
make docker-up           # Запустить все сервисы
make docker-down         # Остановить все сервисы
make docker-down-v       # Остановить с удалением volumes
make docker-logs         # Просмотр логов
make docker-logs-backend # Логи бэкенда
make docker-logs-frontend # Логи фронтенда
make docker-logs-db      # Логи базы данных
make docker-build        # Пересобрать все образы
make docker-build-backend # Пересобрать бэкенд
make docker-build-frontend # Пересобрать фронтенд
make docker-rebuild      # Остановить, пересобрать, запустить
make docker-ps           # Статус сервисов
```

## Доступные адреса

- **Фронтенд**: http://localhost:80
- **Бэкенд API**: http://localhost:8080
- **База данных**: localhost:5432

## Проверка работоспособности

```bash
# Проверить бэкенд
curl http://localhost:8080/health

# Проверить фронтенд
curl http://localhost:80/health
```

## Дополнительные файлы

Созданные файлы для поддержки Docker Compose:

1. **`docker-compose.yml`** (корень) - Главный файл конфигурации
2. **`frontend/Dockerfile`** - Docker образ для фронтенда
3. **`frontend/nginx.conf`** - Конфигурация Nginx для фронтенда
4. **`frontend/.dockerignore`** - Исключения для Docker (фронтенд)
5. **`backend/.dockerignore`** - Исключения для Docker (бэкенд)
6. **`.gitignore`** (корень) - Git ignore правила
7. **`DOCKER-COMPOSE.md`** - Подробная документация

## Архитектура

```
┌─────────────────────────────────────────┐
│         Docker Compose Stack            │
│                                         │
│  ┌──────────────┐                      │
│  │   Frontend   │  Port 80              │
│  │   (Nginx)    │  http://localhost:80  │
│  └──────┬───────┘                      │
│         │                               │
│         │ HTTP                          │
│         ▼                               │
│  ┌──────────────┐                      │
│  │   Backend    │  Port 8080            │
│  │     (Go)     │  http://localhost:8080│
│  └──────┬───────┘                      │
│         │                               │
│         │ SQL                           │
│         ▼                               │
│  ┌──────────────┐                      │
│  │  PostgreSQL  │  Port 5432            │
│  │              │  localhost:5432       │
│  └──────────────┘                      │
│                                         │
└─────────────────────────────────────────┘
```

## Сеть

Все сервисы в одной Docker сети `booking_network`:
- Фронтенд → Бэкенд: `http://backend:8080`
- Бэкенд → БД: `postgres:5432`

## Миграции

Миграции автоматически применяются при первом запуске PostgreSQL из:
- `backend/migrations/001_initial_schema.up.sql`

## Переменные окружения

### Backend
```env
SERVER_PORT=8080
APP_ENV=production
DATABASE_URL=postgres://postgres:postgres@postgres:5432/booking_db?sslmode=disable
```

### Frontend
```env
VITE_API_BASE_URL=http://localhost:8080
```

### Database
```env
POSTGRES_DB=booking_db
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
```

## Устранение неполадок

### База не запускается
```bash
docker-compose logs postgres
```

### Бэкенд не подключается к БД
```bash
docker-compose restart backend
```

### Фронтенд не подключается к бэкенду
```bash
docker-compose build frontend
docker-compose up -d frontend
```

## Production рекомендации

1. ✅ Изменить пароли БД
2. ✅ Использовать Docker secrets
3. ✅ Настроить reverse proxy
4. ✅ Добавить SSL/TLS
5. ✅ Настроить мониторинг
6. ✅ Настроить бэкапы БД

---

**Статус**: ✅ ГОТОВО
**Дата**: 6 апреля 2026
**Результат**: Полный стек приложения готов к запуску через Docker Compose
