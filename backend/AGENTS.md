# BACKEND KNOWLEDGE BASE

## OVERVIEW
Go 1.25 API service implementing the calendar-booking domain against PostgreSQL, organized around `cmd/`, `internal/`, and `migrations/`.

## STRUCTURE
```text
backend/
├── cmd/server/        # app bootstrap
├── internal/config/   # env/config loading
├── internal/db/       # connection and migration helpers
├── internal/handlers/ # HTTP translation layer
├── internal/models/   # domain/data structs
├── internal/repositories/ # persistence
├── internal/services/ # business rules
├── migrations/        # SQL schema files
└── scripts/           # local API smoke helpers
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| App startup | `cmd/server/main.go` | entrypoint |
| Business rules | `internal/services/` | booking conflicts, slot generation |
| DB queries | `internal/repositories/` | SQL / persistence logic |
| HTTP contracts | `internal/handlers/` | request validation and status codes |
| Migrations | `migrations/` | mounted by root compose into Postgres init |
| Config | `internal/config/` and `.env` | local env wiring |

## CONVENTIONS
- Keep layering clean: handlers translate HTTP, services enforce rules, repositories talk to DB.
- Preserve the no-overlapping-bookings invariant in service/repository changes.
- Use project Make targets when possible: `make backend-build`, `make backend-run`, `make backend-test`, `make backend-lint`, `make backend-fmt`.
- Backend CI uses `go test ./... -v -race -coverprofile=coverage.out`, `go vet`, `gofmt -l .`, and `go build -v ./...`.
- Current CI coverage floor is 30%; do not lower it casually.

## ANTI-PATTERNS
- Do not put business rules only in handlers.
- Do not bypass migrations with ad-hoc schema drift; schema starts from `migrations/`.
- Do not treat `README` testing text as current truth; real Go tests already exist across handlers/services/repositories.
- Do not commit `.env`, `coverage.out`, or generated coverage HTML.

## COMMANDS
```bash
make backend-db-up
make backend-run
make backend-test
make backend-test-coverage
make backend-lint
make backend-fmt
```

## NOTES
- Local DB compose lives here, but repo-wide stack wiring lives in root `docker-compose.yml`.
- If backend guidance must become more granular later, the next reasonable split is `internal/handlers`, `internal/services`, and `internal/repositories`.
