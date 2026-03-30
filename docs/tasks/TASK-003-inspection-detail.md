---
id: TASK-003
title: "/inspections/[id] — Inspection detail page"
phase: "Phase 1 - Core Loop"
status: todo
priority: high
assigned_to: [frontend]
depends_on: [TASK-001]
created_at: 2026-03-27
---

## Summary

The inspection detail page is the most complex and important page in Veylo. It serves three roles:
- **Inspector** (ENTRY stage): adds findings/damages with location, type, description
- **Evaluator** (EVALUATION stage): assesses each finding with severity, repair method, cost breakdown
- **Manager** (REVIEW stage): reviews all findings, approves or returns the inspection

This single page handles the entire inspection lifecycle. Findings are managed inline, not on a separate page.

## Acceptance criteria

### Header section
- [ ] Shows vehicle info: brand, model, license plate, VIN
- [ ] Shows inspection metadata: contract number, created date, current status + system stage
- [ ] Status transition buttons shown based on allowed transitions from workflow
- [ ] Clicking a transition button calls the transition API and updates the UI

### Findings section (ENTRY stage -- Inspector view)
- [ ] List of existing findings displayed as cards
- [ ] "Add finding" button opens inline form or panel
- [ ] Finding form: finding type (dropdown), description (text), location (body area dropdown + optional coordinates)
- [ ] Each finding card shows: type, description, body area
- [ ] Inspector cannot set severity or costs (those fields are hidden in ENTRY)

### Assessment section (EVALUATION stage -- Evaluator view)
- [ ] Each finding has an "Assess" action
- [ ] Assessment form: severity (ACCEPTED / NOT_ACCEPTED / INSURANCE_EVENT), repair method (REPAIR / REPLACEMENT / CLEANING / POLISHING / NO_ACTION), cost breakdown (parts, labor, paint, other -- all in cents)
- [ ] Total cost auto-calculated from breakdown
- [ ] After assessment, finding card shows severity badge + total cost

### Review section (REVIEW stage -- Manager view)
- [ ] All findings displayed with their assessments (read-only)
- [ ] Total cost summary across all findings
- [ ] Approve or return buttons (status transitions)

### Report section (FINAL stage)
- [ ] "Download PDF" button visible when report is available
- [ ] Links to GET /api/v1/inspections/{id}/report

### General
- [ ] Loading skeleton while data fetches
- [ ] 404 handling if inspection not found
- [ ] Optimistic updates for transitions

## API endpoints

- `GET /api/v1/inspections/{id}` -- get inspection detail
- `POST /api/v1/inspections/{id}/transitions` -- transition status. Body: `{ status: string }`
- `GET /api/v1/inspections/{id}/report` -- get PDF report URL
- `POST /api/v1/inspections/{inspectionID}/findings` -- create finding. Body: `{ finding_type, description, location: { body_area, coordinate_x, coordinate_y } }`
- `GET /api/v1/inspections/{inspectionID}/findings` -- list findings
- `PUT /api/v1/inspections/{inspectionID}/findings/{id}/assessment` -- assess finding. Body: `{ severity, repair_method, cost_breakdown: { parts, labor, paint, other } }`

## UI components

- Card (exists) for finding cards
- Badge (needs to be added) for severity and status
- Dialog or Sheet (needs to be added) for finding form
- Select (exists) for dropdowns
- Input (exists) for cost fields
- Button (exists) for transitions
- Tabs (needs to be added) -- optional, to separate findings/summary/audit

## Technical notes

- Create `web/src/app/(app)/inspections/[id]/page.tsx`
- Create `web/src/features/inspections/components/inspection-detail.tsx`
- Create `web/src/features/inspections/hooks/use-inspection.ts` (single inspection)
- Create `web/src/features/findings/` feature folder with api, types, hooks, components
- The page should be role-aware: show/hide actions based on current user role + inspection stage
- Cost inputs should accept euros (display) but send cents (API). Multiply by 100 on submit.
- Body area values: front_bumper, rear_bumper, hood, trunk, roof, left_front_door, left_rear_door, right_front_door, right_rear_door, left_front_fender, left_rear_fender, right_front_fender, right_rear_fender, windshield, rear_window, left_mirror, right_mirror, left_headlight, right_headlight, left_taillight, right_taillight, wheel_fl, wheel_fr, wheel_rl, wheel_rr
