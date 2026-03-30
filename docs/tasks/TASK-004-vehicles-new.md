---
id: TASK-004
title: "/vehicles/new — Add vehicle page"
phase: "Phase 1 - Core Loop"
status: todo
priority: high
assigned_to: [frontend]
depends_on: []
created_at: 2026-03-27
---

## Summary

Standalone page to add a new vehicle to the organization's registry. This form is also reused as an inline/modal component during inspection creation (TASK-002). Inspectors and managers use this to register vehicles before or during inspections.

## Acceptance criteria

- [ ] Form fields: VIN (17 chars, validated), license plate, brand, model, body type, fuel type, transmission, odometer reading, color, engine power
- [ ] VIN validation: exactly 17 alphanumeric characters
- [ ] License plate: sanitized (uppercase, trimmed)
- [ ] Body type dropdown: sedan, hatchback, wagon, SUV, coupe, convertible, van, pickup, other
- [ ] Fuel type dropdown: petrol, diesel, electric, hybrid, plugin_hybrid, lpg, cng
- [ ] Transmission dropdown: manual, automatic
- [ ] On submit, calls create vehicle API
- [ ] Success: redirect to /vehicles/[id] (standalone) or return vehicle to caller (inline mode)
- [ ] Error: show "vehicle with this VIN or license plate already exists" on 409 conflict
- [ ] Loading state on submit button

## API endpoints

- `POST /api/v1/assets/vehicles` -- create vehicle
  - Body: `{ vin, license_plate, brand, model, body_type, fuel_type, transmission, odometer_reading, color, engine_power }`
  - Response: created asset object

## UI components

- Form (React Hook Form + Zod)
- Input (exists) for text fields
- Select (exists) for dropdowns
- Button (exists)
- Card (exists) for form container

## Technical notes

- Create `web/src/app/(app)/vehicles/new/page.tsx`
- Create `web/src/features/vehicles/components/vehicle-form.tsx` -- reusable form component
- Create `web/src/features/vehicles/hooks/use-create-vehicle.ts`
- Create `web/src/features/vehicles/api.ts`
- Create `web/src/features/vehicles/types.ts`
- Create `web/src/features/vehicles/schemas.ts` (Zod schema for validation)
- The form component should accept an `onSuccess` callback prop for inline mode (returns created vehicle) vs standalone mode (redirects)
