---
id: TASK-001
title: "/inspections — List page"
phase: "Phase 1 - Core Loop"
status: todo
priority: high
assigned_to: [frontend]
depends_on: []
created_at: 2026-03-27
---

## Summary

The inspections list page is the primary working surface for all roles. Inspectors see their active inspections, evaluators find inspections awaiting assessment, managers see those pending review. This is the most critical page in the entire app — it's where users spend 80% of their time navigating to specific inspections.

## Acceptance criteria

- [ ] Page renders a paginated table of inspections for the current organization
- [ ] Each row shows: contract number, vehicle (brand + model + plate), status, system stage, created date, assigned inspector
- [ ] Status is displayed as a colored badge (color per system stage: ENTRY=blue, EVALUATION=yellow, REVIEW=orange, FINAL=green)
- [ ] Table supports pagination (page size: 20, with page controls)
- [ ] Empty state shown when no inspections exist, with CTA to create first inspection
- [ ] Clicking a row navigates to `/inspections/[id]`
- [ ] Loading skeleton shown while data is fetching
- [ ] Page title: "Inspections"

## API endpoints

- `GET /api/v1/inspections?page=1&page_size=20` — list inspections (paginated)
  - Response: `{ items: [...], total: number, page: number, page_size: number }`

## UI components

- `Table` (shadcn/ui — needs to be added)
- `Badge` (shadcn/ui — needs to be added) for status display
- `Button` for pagination controls
- `Skeleton` (shadcn/ui — needs to be added) for loading state
- Custom `EmptyState` component (reuse pattern from dashboard-empty-state)

## Technical notes

- Create `web/src/app/(app)/inspections/page.tsx`
- Create `web/src/features/inspections/components/inspections-table.tsx`
- Extend existing `web/src/features/inspections/hooks/use-inspections.ts` (already exists with basic pagination)
- Extend `web/src/features/inspections/types.ts` if response shape needs more fields
