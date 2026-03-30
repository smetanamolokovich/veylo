---
id: TASK-005
title: "/vehicles — Vehicle registry list page"
phase: "Phase 2 - Data Entry"
status: todo
priority: medium
assigned_to: [frontend, backend]
depends_on: [TASK-004]
created_at: 2026-03-27
---

## Summary

The vehicle registry page lists all vehicles in the organization. Used by managers and admins to browse, search, and manage the vehicle fleet. Also serves as a lookup for inspectors when they need to find a vehicle.

## Acceptance criteria

- [ ] Paginated table of vehicles for the current organization
- [ ] Columns: license plate, VIN, brand, model, body type, fuel type, odometer, color
- [ ] Clicking a row navigates to /vehicles/[id]
- [ ] "Add vehicle" button navigates to /vehicles/new
- [ ] Empty state with CTA to add first vehicle
- [ ] Loading skeleton while data is fetching
- [ ] Pagination controls (page size: 20)

## API endpoints

- **MISSING:** `GET /api/v1/assets?type=vehicle&page=1&page_size=20` -- list assets. This endpoint does NOT exist yet. Backend work required.
  - Needed response: `{ items: [...], total: number, page: number, page_size: number }`

## UI components

- Table (shadcn/ui)
- Button (exists)
- Skeleton (shadcn/ui)
- Custom EmptyState component

## Technical notes

- Create `web/src/app/(app)/vehicles/page.tsx`
- Create `web/src/features/vehicles/components/vehicles-table.tsx`
- Create `web/src/features/vehicles/hooks/use-vehicles.ts`
- **Backend dependency:** Need to add a list assets endpoint. Currently only POST (create) and GET by ID exist for assets. This is a blocker for this task and also for vehicle search in TASK-002.
