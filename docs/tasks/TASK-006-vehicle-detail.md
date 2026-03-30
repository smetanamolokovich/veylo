---
id: TASK-006
title: "/vehicles/[id] — Vehicle detail page"
phase: "Phase 2 - Data Entry"
status: todo
priority: medium
assigned_to: [frontend]
depends_on: [TASK-005]
created_at: 2026-03-27
---

## Summary

Vehicle detail page shows all information about a specific vehicle and its inspection history. Managers use this to review a vehicle's condition over time. Inspectors use it to check past inspections before starting a new one.

## Acceptance criteria

- [ ] Displays all vehicle fields: VIN, license plate, brand, model, body type, fuel type, transmission, odometer, color, engine power
- [ ] Shows inspection history for this vehicle (list of past inspections with status, date, inspector)
- [ ] Each inspection in the history is a clickable link to /inspections/[id]
- [ ] "Start new inspection" button creates a new inspection pre-filled with this vehicle
- [ ] 404 handling if vehicle not found
- [ ] Loading skeleton while data fetches

## API endpoints

- `GET /api/v1/assets/{id}` -- get vehicle detail
- `GET /api/v1/inspections?asset_id={id}` -- list inspections for this vehicle (may need backend filter support)

## UI components

- Card (exists) for vehicle info sections
- Table or list for inspection history
- Button (exists) for actions
- Badge for inspection statuses

## Technical notes

- Create `web/src/app/(app)/vehicles/[id]/page.tsx`
- Create `web/src/features/vehicles/components/vehicle-detail.tsx`
- Create `web/src/features/vehicles/hooks/use-vehicle.ts`
- **Backend note:** The inspections list endpoint may need a filter by asset_id. Check if ListInspectionsRequest supports this; if not, backend enhancement needed.
