# PROJECT KNOWLEDGE BASE

**Generated:** 2026-04-16T20:02:01Z
**Commit:** b490b7b
**Branch:** master

## OVERVIEW
Calendar booking monorepo with four real working areas: Go backend, React frontend, Playwright E2E, and TypeSpec API contract. The old planning-only description is stale; implemented code, tests, Docker, and CI all exist.

## STRUCTURE
```text
./
├── backend/        # Go API, DB access, migrations, service logic
├── frontend/       # React SPA, Vitest, Playwright config
├── typespec/       # Contract-first API spec, OpenAPI generation
├── .github/        # CI workflows; hexlet-check is protected
├── Makefile        # Canonical command entrypoint across areas
├── SPEC.MD         # Product/process rules and high-level requirements
└── docker-compose.yml  # Full-stack local/runtime wiring
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Repo workflow / commands | `Makefile` | Prefer root `make` targets over ad-hoc commands |
| Product rules / backlog | `SPEC.MD` | API-first, TDD, commit hygiene, booking constraints |
| Backend runtime entry | `backend/cmd/server/main.go` | Go app bootstrap |
| Frontend runtime entry | `frontend/src/main.tsx` | React app bootstrap |
| API contract | `typespec/main.tsp` | Upstream of backend/frontend integration |
| Local stack wiring | `docker-compose.yml` | Frontend → backend → postgres |
| Repo-wide CI rule | `.github/workflows/hexlet-check.yml` | Do not edit/delete |

## KEY MODULES
| Path | Role |
|------|------|
| `backend/internal/handlers` | HTTP layer, request/response translation |
| `backend/internal/services` | Business rules, slot/booking invariants |
| `backend/internal/repositories` | PostgreSQL access |
| `frontend/src/pages/guest/BookingPage.tsx` | Largest guest booking flow UI |
| `frontend/src/components/owner/EventTypeManagement.tsx` | Owner CRUD / slot generation UI |
| `frontend/e2e/specs` | Browser flows and API-integration checks |
| `typespec/tspconfig.yaml` | OpenAPI 3.1 emitter config |

## CONVENTIONS
- Treat this as a mixed monorepo, not a single-app repo.
- Root `Makefile` is the canonical command hub.
- Answer the user in Russian unless they explicitly request another language.
- API-first still applies: update `typespec/` before changing API behavior in app layers.
- TDD is expected, but current enforcement is uneven: backend CI is stricter than frontend CI.
- `frontend/mock-api.sh` depends on generated OpenAPI at `typespec/tsp-output/schema/openapi.yaml`.
- Root Docker Compose mounts `backend/migrations` into Postgres init.

## ANTI-PATTERNS (THIS PROJECT)
- Do not edit/delete `.github/workflows/hexlet-check.yml` or rename the repository.
- Do not use global npm installs; prefer `npx`, npm scripts, or `make`.
- Do not commit generated/dependency/secret paths named in project docs: `node_modules/`, `vendor/`, `tsp-output/`, `.env`.
- Do not assume `AGENTS.md` or `README` status text is current without checking live files; older docs still say “no implementation code exists.”
- Do not count report/build artifacts as real module complexity when exploring: `frontend/playwright-report/`, `frontend/test-results/`, `frontend/coverage/`, `frontend/dist/`, `typespec/tsp-output/`.

## UNIQUE STYLES
- Root docs are thinner than child docs; operational detail lives mostly in `backend/README.md`, `frontend/README.md`, and `frontend/e2e/README.md`.
- CI is split by area: backend workflow, Playwright workflow, plus repo-wide Hexlet workflow.
- No root workspace manager; each area owns its own toolchain.

## COMMANDS
```bash
make help
make compile
make openapi
make fronttest
make backtest
make frontend-e2e
make docker-up
make docker-down
make alltest
```

## NOTES
- `RULES.md` is referenced by older guidance but does not exist in the repo.
- Pre-commit only fixes EOF/trailing whitespace; it does not run tests.
- Worth keeping child AGENTS only at real boundaries: `backend/`, `frontend/`, `frontend/e2e/`, `typespec/`.
