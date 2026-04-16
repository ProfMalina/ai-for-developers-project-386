# TYPESPEC KNOWLEDGE BASE

## OVERVIEW
Contract-first API specification for the calendar-booking system. This directory is the upstream source for endpoint shape, validation, and generated OpenAPI consumed by frontend mocking and backend alignment.

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Main contract | `main.tsp` | service, models, routes, validation |
| Emitter config | `tspconfig.yaml` | OpenAPI 3.1 YAML to `tsp-output/schema` |
| Scripts | `package.json` | compile/watch/format commands |

## CONVENTIONS
- Change API shape here first, then update backend/frontend to match.
- `main.tsp` is source of truth; `tsp-output/` is generated output.
- Keep validation constraints explicit in the spec: lengths, min/max values, email/date formats, timezone fields.
- Frontend mock API depends on generated `tsp-output/schema/openapi.yaml`.

## ANTI-PATTERNS
- Do not hand-edit `tsp-output/`.
- Do not bypass compile/format checks after modifying the spec.
- Do not drift endpoint behavior between TypeSpec and implementation layers.

## COMMANDS
```bash
npx tsp compile .
npx tsp fmt . --check
npx tsp fmt .
make compile
make fmt-check
make openapi
```

## NOTES
- Current emitter target is OpenAPI 3.1 YAML under `{output-dir}/schema`.
- `main.tsp` is one of the largest hand-written files in the repo; read it before changing API-facing behavior elsewhere.
