# Project Context: AI for Developers — Calendar Booking App

## Project Overview

This is a **calendar booking application** being developed as part of the Hexlet "AI for Developers" course (project #386). The project is currently in the **specification/planning phase** — no source code has been implemented yet.

The application enables a calendar owner to define event types and manage bookings, while guests can browse available event types, view a calendar, and book free time slots.

## Planned Architecture

| Layer | Technology |
|-------|-----------|
| Backend | Go |
| Frontend | React SPA |
| Database | PostgreSQL |
| Deployment | Docker Compose, deployed via SSH |
| API Design | TypeSpec (generates OpenAPI spec) |

## Domain Entities

1. **Owner** — A single predefined profile (no registration/auth). Creates event types and views upcoming meetings.
2. **Event Type** — Defined by `id`, name, description, and duration in minutes.
3. **Slot** — A time slot derived from an event type's duration.
4. **Booking** — A reservation made by a guest for a specific slot.

## User Roles & Capabilities

### Owner
- Create event types (id, name, description, duration)
- View a single list of all upcoming bookings across all event types

### Guest
- View a public page listing available event types (name, description, duration)
- Select an event type, open a calendar, and choose a free slot
- Create a booking for the selected slot

### Business Rule
- **No overlapping bookings**: Two bookings cannot occupy the same time slot, even if they are different event types.

## Development Approach

- **API First**: Define the API contract (via TypeSpec) before implementation.
- **TDD (Test-Driven Development)**: Red → Green → Refactor cycle.
- **No authentication**: The owner is a single pre-configured profile; guests book without accounts.
- **No global dependencies**: All tools must be run via `npx` or npm scripts, never install globally with `npm install -g`. See [RULES.md](./RULES.md) for details.

## Key Files

| File | Description |
|------|-------------|
| `SPEC.MD` | Project specification — describes domain, roles, rules, tech stack, and tasks |
| `RULES.md` | **Development rules** — includes prohibition on global dependencies (use npx, not npm install -g) |
| `README.md` | Minimal README with Hexlet CI status badge |
| `.pre-commit-config.yaml` | Pre-commit hooks (end-of-file-fixer, trailing-whitespace) |
| `.github/workflows/hexlet-check.yml` | Hexlet automated test runner (do not modify) |

## Open Tasks (from SPEC.MD)

- [ ] Define domain entities: owner, event type, slot, booking, and guest scenario.
- [ ] Formulate tasks for an AI agent and prepare a TypeSpec specification that fixes the API contract.
- [ ] Verify the specification covers owner and guest scenarios, including the slot occupancy rule.

## Current State

**No implementation code exists yet.** The repository contains only project documentation and CI configuration. The next step is to produce a TypeSpec API specification before writing any Go or frontend code.

@SPEC.MD
