---
id: TASK-002
title: "/inspections/new — Create inspection page"
phase: "Phase 1 - Core Loop"
status: todo
priority: high
assigned_to: [frontend]
depends_on: [TASK-001]
created_at: 2026-03-27
---

## Summary

The create inspection page allows inspectors and managers to start a new inspection. The user selects an existing vehicle (or quick-adds a new one inline) and enters the contract number. This is the entry point to the core inspection workflow. Speed is critical -- inspectors create 10-20 inspections per day.

## Acceptance criteria

- [ ] Form with fields: vehicle selector, contract number
- [ ] Vehicle selector searches existing vehicles by VIN, license plate, or brand+model
- [ ] "Quick add vehicle" button opens an inline form (or modal) to create a vehicle without leaving the page (see TASK-006)
- [ ] After vehicle is created inline, it is auto-selected in the selector
- [ ] Contract number is a required text field
- [ ] On submit, calls create inspection API, then redirects to /inspections/[id]
- [ ] Validation: contract number required, vehicle required
- [ ] Error handling: show API errors (e.g., duplicate contract)
- [ ] Loading state on submit button

## API endpoints

- `POST /api/v1/inspections` -- create inspection
  - Body: `{ asset_id: string, contract_number: string }`
  - Response: created inspection object
- `POST /api/v1/assets/vehicles` -- quick-add vehicle (inline)

## UI components

- Form (React Hook Form + Zod)
- Input (exists) for contract number
- Select or Combobox (shadcn/ui -- needs to be added) for vehicle selection
- Dialog (shadcn/ui -- needs to be added) for quick-add vehicle modal
- Button (exists)

## Technical notes

- Create `web/src/app/(app)/inspections/new/page.tsx`
- Create `web/src/features/inspections/components/create-inspection-form.tsx`
- Create `web/src/features/inspections/hooks/use-create-inspection.ts`
- Need a vehicle list/search hook -- backend currently has no GET /api/v1/assets list endpoint. Options: (a) add backend endpoint, or (b) store recently created vehicles client-side. Recommend adding GET /api/v1/assets?type=vehicle to backend.
- The inline vehicle creation reuses the form from TASK-006
