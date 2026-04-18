# Issue #8 gap audit matrix

## Audit basis

This matrix is based on live repository files and command output from the isolated worktree at `/home/admn/git/ai-for-developers-project-386/.worktrees/issue-8-gap-closure`.

Observed repo baseline used for this report:

- `SPEC.MD` is present and defines the product requirements for guest flow, owner flow, slot conflict rule, pagination, slot generation, i18n, theme support, and cookie notice.
- `frontend/package.json`, `frontend/playwright.config.ts`, `frontend/e2e/README.md`, `frontend/e2e/specs/*.ts`, `Makefile`, and `.github/workflows/playwright.yml` are present and were audited directly.
- The expected `docs/superpowers/plans/2026-04-18-issue-8-gap-closure.md` and `docs/superpowers/specs/2026-04-18-issue-8-playwright-e2e-design.md` files are not present in this worktree. The root listing contains no `docs/` directory before this report file was created.

Status vocabulary used in this matrix is limited to `implemented`, `partial`, `blocked-by-product`, `blocked-by-test-data`, and `missing-now`.

## Execution scope

1. Implement now:

- Owner cancellation happy-path coverage: `frontend/e2e/specs/owner-flow.spec.ts` does not currently contain a passing owner booking cancellation happy-path, while `frontend/e2e/README.md` still advertises cancellation coverage as a scenario and `api-integration.spec.ts` only covers the cancellation error path.
- Owner bookings pagination coverage: `SPEC.MD` requires pagination for bookings, and the owner/common specs mock `pagination` objects, but the current suite only provides weak evidence beyond bookings-list visibility and does not provide strong current passing proof of real pagination behavior.
- README/config/workflow truthfulness fixes: the audited docs and runner metadata overstate current readiness relative to the observed baseline, because `frontend/e2e/README.md` presents broad feature coverage while `make frontend-e2e-chromium` currently fails with `22 failed`, `19 passed`, `2 skipped`; truthfulness must be verified against `frontend/playwright.config.ts`, `Makefile`, and `.github/workflows/playwright.yml`.

2. Keep out of scope for issue closure unless already implemented in app:

- Cookie consent UI: `SPEC.MD` requires cookie warnings, but `frontend/e2e/README.md` says the cookie banner is not yet implemented and current common tests only verify graceful loading.
- RU/EN switcher: `SPEC.MD` requires Russian and English selection, but `frontend/e2e/specs/common.spec.ts` skips English switching and persistence and the README says language switching is not yet implemented.
- Explicit light/dark toggle UI: current common coverage verifies system color scheme handling, but `frontend/e2e/README.md` says theme switching is not yet implemented.

3. Verify before closure:

- Chromium Playwright subset: re-run `make frontend-e2e-chromium` and confirm improvement over the current failing baseline of `22 failed`, `19 passed`, `2 skipped`.
- Frontend build: run the frontend production build and confirm it succeeds before any issue-closure claim.

## Inventory command evidence

### 1. Exact inventory command: suite listing

Command:

```bash
cd frontend && npm run test:e2e -- --list
```

Observed output summary:

- `playwright test --list` runs through the package script `test:e2e`
- Total listed inventory: `129 tests in 4 files`
- Browsers listed: Chromium, Firefox, WebKit
- Files listed: `api-integration.spec.ts`, `common.spec.ts`, `guest-flow.spec.ts`, `owner-flow.spec.ts`

### 2. Exact inventory command: Makefile recipes

Command:

```bash
make -n frontend-e2e frontend-e2e-ui frontend-e2e-headed frontend-e2e-debug frontend-e2e-chromium frontend-e2e-report
```

Observed output:

```text
cd frontend && npx playwright test --workers=1
cd frontend && npx playwright test --ui
cd frontend && npx playwright test --headed
cd frontend && npx playwright test --debug
cd frontend && npx playwright test --project=chromium --workers=1
cd frontend && npx playwright show-report
```

### 3. Supporting inventory context from audited files

- `frontend/package.json` defines `test:e2e`, `test:e2e:ui`, `test:e2e:headed`, and `test:e2e:debug`.
- `Makefile` exposes `frontend-e2e`, `frontend-e2e-ui`, `frontend-e2e-headed`, `frontend-e2e-debug`, `frontend-e2e-chromium`, and `frontend-e2e-report`.
- `frontend/playwright.config.ts` defines Chromium, Firefox, and WebKit projects and uses a Vite `webServer` with `baseURL` `http://localhost:5173`.

### 4. Current Chromium baseline

Command:

```bash
make frontend-e2e-chromium
```

Observed output summary:

- Result: failing baseline
- Totals: `22 failed`, `19 passed`, `2 skipped`
- Confirmed failing areas:
  - Guest flow: all 7 Chromium tests currently fail
  - Owner flow: 5 Chromium tests currently fail, 2 pass
  - API integration: multiple failure-handling tests currently fail
  - Common features: mobile responsive test currently fails

Representative failing assertions from the command output:

- `guest-flow.spec.ts` expects event text and booking headings that are not found in the live UI
- `owner-flow.spec.ts` waits for owner form controls such as duration and slot interval fields that are not found
- `api-integration.spec.ts` has failing error-handling expectations for 201, 409, 400, 500, network failure, API unavailable, malformed JSON, owner creation error, and mocked 500 cases
- `common.spec.ts` mobile responsive test expects `article` cards, but the Chromium run reports `Received: 0`

## Requirement to repository matrix

| Issue requirement | Current file(s) | Status | Evidence | Action |
| --- | --- | --- | --- | --- |
| Playwright suite inventory exists and is wired through the required package script | `frontend/package.json`, `frontend/e2e/specs/*.ts` | implemented | `frontend/package.json` defines `test:e2e`; the exact inventory command `cd frontend && npm run test:e2e -- --list` reports `129 tests in 4 files` across Chromium, Firefox, and WebKit. | Keep the existing script entrypoint and use it as the baseline inventory command for issue closure. |
| Root Makefile exposes the required Playwright execution targets | `Makefile` | implemented | The exact inventory command `make -n frontend-e2e frontend-e2e-ui frontend-e2e-headed frontend-e2e-debug frontend-e2e-chromium frontend-e2e-report` expands to `npx playwright test --workers=1`, `--ui`, `--headed`, `--debug`, `--project=chromium --workers=1`, and `show-report`. | Do not add new command surfaces; Task 2 should treat these recipes as the command source of truth when reconciling docs and report language. |
| Playwright browser/config harness is present for Chromium, Firefox, and WebKit | `frontend/playwright.config.ts` | implemented | `frontend/playwright.config.ts` defines `chromium`, `firefox`, and `webkit` projects and a `webServer` running Vite at `http://localhost:5173`. The package-script inventory command also confirms three browser projects in the listed suite. | Treat harness setup as done; closure work belongs in reliability and spec fit. |
| README coverage claims are not fully truthful relative to the observed suite and baseline | `frontend/e2e/README.md`, `frontend/e2e/specs/*.ts` | partial | `frontend/e2e/README.md` says `**Total: 138 tests**`, but the exact inventory command reports `129 tests in 4 files`. The README also presents several scenarios with ✅ status while the current Chromium baseline still fails across guest flow, owner flow, API handling, and one common responsive test. | Task 2 should update README totals and coverage language so they match the observed suite inventory and clearly distinguish implemented coverage from currently failing coverage. |
| Playwright config comments are not fully truthful relative to actual runtime behavior | `frontend/playwright.config.ts`, `.github/workflows/playwright.yml` | partial | The comment above `webServer` says to start the dev server manually and that CI starts the server separately, but `frontend/playwright.config.ts` still declares `webServer.command = 'npx vite'` and the CI workflow simply runs `npx playwright test`, which relies on Playwright config behavior. | Task 2 should rewrite the config comments to match the actual `webServer` behavior and the current CI wiring instead of describing an older manual-server model. |
| Guest flow scenarios required for issue #8 are present and currently executable as tests | `frontend/e2e/specs/guest-flow.spec.ts` | partial | The suite inventory lists 7 Chromium guest-flow tests. The current `make frontend-e2e-chromium` baseline shows all 7 guest-flow tests failing, including list view, booking page navigation, successful booking, unavailable slots, and validation checks. | Fix guest-flow selector/data mismatches before using this area as closure evidence. |
| Owner flow scenarios required for issue #8 are present and currently executable as tests | `frontend/e2e/specs/owner-flow.spec.ts` | partial | The suite inventory lists owner dashboard, create, edit, delete, slot generation, bookings list, and validation scenarios. The current Chromium baseline fails create, edit, delete, slot generation, and validation, while dashboard and bookings list pass. | Fix owner CRUD and slot-generation interaction coverage before closure. |
| Owner cancellation happy-path coverage is traceable only to helper support, not to a runnable happy-path spec | `frontend/e2e/pages/OwnerDashboard.ts`, `frontend/e2e/specs/owner-flow.spec.ts`, `frontend/e2e/specs/api-integration.spec.ts`, `frontend/e2e/README.md` | missing-now | `OwnerDashboard.ts` already provides `cancelBooking()` helper support, and the README lists cancel booking as a scenario that requires test data. However, the current suite inventory for `owner-flow.spec.ts` has no owner cancellation happy-path test, and `api-integration.spec.ts` only covers booking cancellation error handling rather than successful cancellation. | Task 3 should add an explicit owner cancellation happy-path to `owner-flow.spec.ts`: seed or mock a visible booking, reuse `OwnerDashboard.cancelBooking()`, mock a successful `DELETE /api/bookings/:id` plus refreshed bookings list, and assert both success feedback and removal or status change of the cancelled booking. |
| Common feature coverage exists for navigation, responsiveness, theme detection, and Russian default content | `frontend/e2e/specs/common.spec.ts` | partial | The suite inventory lists navigation, responsive, cookie, theme, i18n, and date/time cases. The current Chromium baseline still fails the mobile responsive test, while some navigation and common checks pass. | Keep this area in closure verification, especially the currently failing mobile responsive scenario. |
| API mocking and API failure coverage exists | `frontend/e2e/specs/api-integration.spec.ts` | partial | The suite inventory lists success, 409, 400, 404, 500, network failure, API unavailable, malformed JSON, missing fields, owner errors, and mocked API tests. The current Chromium baseline fails many of these API-handling scenarios. | Stabilize existing API tests rather than adding a new test category. |
| Slot conflict rule is covered in tests | `SPEC.MD`, `frontend/e2e/specs/api-integration.spec.ts`, `frontend/e2e/README.md` | blocked-by-test-data | `SPEC.MD` requires one booking per slot; `api-integration.spec.ts` has a `409` conflict scenario; `frontend/e2e/README.md` still describes conflict coverage as requiring setup. The current Chromium run also fails the conflict test. | Keep the rule in scope, but treat it as needing reliable test-data/setup evidence before using it for closure. |
| Owner bookings pagination is represented in helpers and docs, but not as a demonstrated happy-path flow | `SPEC.MD`, `frontend/e2e/pages/OwnerDashboard.ts`, `frontend/e2e/specs/owner-flow.spec.ts`, `frontend/e2e/README.md` | partial | `SPEC.MD` requires pagination for bookings, `frontend/e2e/README.md` claims bookings pagination coverage, and `OwnerDashboard.ts` already has `expectPaginationVisible()` and `goToBookingsPage()` helpers. The current `owner-flow.spec.ts` inventory, however, only shows `should view bookings list`; it does not click a page control, assert a second-page dataset, or verify that pagination controls are visible and wired. | Task 3 should add an owner bookings pagination happy-path that mocks at least two bookings pages, asserts pagination controls are shown, clicks the next page or page number, and verifies that page-2 booking rows replace or extend the page-1 dataset so the behavior is observable rather than implied by static mock shapes. |
| Slot auto-generation with configurable interval is represented in coverage | `SPEC.MD`, `frontend/e2e/specs/owner-flow.spec.ts` | partial | `SPEC.MD` requires auto-generation with 15 or 30 minute intervals. `owner-flow.spec.ts` covers slot generation and interval input, but the current Chromium run times out waiting for the slot interval field. | Keep slot generation in implement-now work because the scenario exists but is broken now. |
| Language switching between Russian and English is available | `SPEC.MD`, `frontend/e2e/specs/common.spec.ts`, `frontend/e2e/README.md` | blocked-by-product | `SPEC.MD` requires language selection. `common.spec.ts` only verifies Russian default text and explicitly skips English switching and persistence; `frontend/e2e/README.md` says language switching is not yet implemented. | Keep out of closure scope unless the app already implements language switching. |
| Theme switching is available to the user beyond system preference detection | `SPEC.MD`, `frontend/e2e/specs/common.spec.ts`, `frontend/e2e/README.md` | blocked-by-product | `SPEC.MD` requires theme selection with system default. Current tests only verify dark/light system preference handling, and the README says theme switching is not yet implemented. | Keep out of closure scope unless manual theme switching already exists in the app. |
| Cookie consent warning exists and is testable | `SPEC.MD`, `frontend/e2e/specs/common.spec.ts`, `frontend/e2e/README.md` | missing-now | `SPEC.MD` requires cookie warnings. `common.spec.ts` only checks graceful first visit and `acceptCookies()` without a real banner, and the README says the cookie banner is not yet implemented. | Do not count cookie consent as closed in issue #8. |
| CI workflow wiring for Playwright exists and should be described as existing infrastructure, not new scope | `.github/workflows/playwright.yml`, `frontend/e2e/README.md` | implemented | The workflow runs on frontend changes for push and pull request to `main` and `master`, installs dependencies and browsers, executes `npx playwright test`, and uploads `playwright-report` and `test-results` artifacts. | Task 2 should keep CI wiring descriptions aligned with the existing workflow and avoid implying that issue #8 needs a new workflow rather than truthful reporting plus suite stabilization. |

## Bottom line

The repository already has the Playwright harness, Makefile entrypoints, config, and CI wiring needed for issue #8, so Tasks 2–4 should focus on gap closure rather than new infrastructure.

Task 2 priority is truthfulness: update README coverage claims and config comments so they match the real suite inventory, current command surface, and existing CI workflow. Task 3 priority is executable owner coverage: add a real owner cancellation happy-path and observable bookings pagination happy-path on top of the helpers that already exist. Task 4 priority is verification: confirm the Chromium subset and frontend build after those changes. Product-gated items such as cookie consent UI, RU/EN switching, and an explicit light/dark toggle remain secondary and out of scope unless the app already implements them.

## Final closure decision

- [x] Playwright setup exists and runs
- [x] Guest booking flow is covered
- [x] Owner management flow is covered
- [x] API failure handling is covered
- [x] CI integration exists
- [x] Remaining unchecked items are blocked by product functionality, not missing E2E infrastructure

Fresh `make frontend-e2e-chromium` evidence now shows `42 passed`, `3 skipped`, `0 failed`, which closes the earlier gap from the `22 failed`, `19 passed`, `2 skipped` baseline. The only remaining non-running checks are product-gated skips for cookie banner acceptance and RU/EN language switching persistence, not unresolved E2E infrastructure gaps. Together with the passing frontend build, the current worktree evidence supports treating issue #8 as ready for closure.

## Recommended issue action

Close issue #8 if all boxes above are checked.
Otherwise, split remaining blockers into follow-up product issues and keep #8 open only for true E2E work.

Based on the fresh Task 5 evidence, the first path now applies: issue #8 is ready to close. The Chromium suite is green for all runnable tests, the build passes, and the remaining skipped checks are explicitly tied to product functionality that the frontend still does not implement.

## Issue closure summary draft

```md
Issue #8 is ready for closure based on the final verification in the gap-closure worktree.

Fresh verification shows:

- `make frontend-e2e-chromium`: `42 passed`, `3 skipped`, `0 failed`.
- `cd frontend && npm run build`: passed, with only the existing Vite chunk-size warning.
- The intended issue-8 gap work is now in place: truthfulness fixes landed in the E2E docs/config, guest-flow is covered through stable mocked guest endpoints, and owner-flow now includes runnable coverage for CRUD, slot generation, bookings pagination, and booking cancellation happy path.

The remaining skipped checks are product-gated areas the frontend still does not implement yet: cookie banner acceptance and RU/EN language-switch persistence. They are no longer unresolved E2E harness failures.

Recommended action: close issue #8 and track the product-gated placeholders separately if maintainers want dedicated follow-up issues.
```
