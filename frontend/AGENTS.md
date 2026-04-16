# FRONTEND KNOWLEDGE BASE

## OVERVIEW
React 19 + Vite + TypeScript SPA using Mantine, Vitest, and Playwright. Supports guest booking flows and owner management flows against either Prism mock API or the real backend.

## STRUCTURE
```text
frontend/
├── src/api/              # Axios client
├── src/components/       # shared and owner-facing UI
├── src/pages/            # guest/owner route-level screens
├── src/test/             # Vitest setup + MSW mocks
├── e2e/                  # Playwright suite; see child AGENTS
├── .env.development      # Prism mock API base URL
├── .env.production       # real backend base URL
└── mock-api.sh           # Prism launcher using TypeSpec output
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| App bootstrap | `src/main.tsx` | routing, providers |
| Routing shell | `src/App.tsx` | top-level app composition |
| Guest booking flow | `src/pages/guest/BookingPage.tsx` | largest guest UI module |
| Owner CRUD flow | `src/components/owner/EventTypeManagement.tsx` | create/edit/delete/generate |
| Owner bookings | `src/components/owner/BookingsList.tsx` | paginated owner list |
| API mocks | `src/test/setup.ts`, `src/test/mocks.ts` | Vitest + MSW |

## CONVENTIONS
- Prefer repo-local scripts or root Make targets; package scripts are in `package.json`.
- Mock-first local development is normal: `./frontend/mock-api.sh` reads `../typespec/tsp-output/schema/openapi.yaml`.
- `dist/`, `coverage/`, `playwright-report/`, and `test-results/` are outputs, not source.
- Frontend validation and error handling are already covered by both Vitest and Playwright; extend those instead of inventing new test harnesses.
- Mantine and route/page split are the dominant UI organization patterns.

## ANTI-PATTERNS
- Do not hand-edit generated build/report artifacts.
- Do not hardcode a backend URL in source when env files already define `VITE_API_BASE_URL` flows.
- Do not duplicate API shapes if TypeSpec or shared frontend types already define them.
- Do not treat placeholder E2E features as implemented; cookie banner, i18n, and theme switching are still marked incomplete.

## COMMANDS
```bash
npm run dev
npm run build
npm run lint
npm test -- --run
npm run test:coverage -- --run
make frontend-e2e
```

## NOTES
- `frontend/e2e/` has distinct page-object and browser-runner rules; use its child AGENTS for test-only work.
- The frontend README is much richer than the root README and is usually the best local orientation doc.
