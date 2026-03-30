---
id: TASK-008
title: "/dashboard — Role-aware dashboard"
phase: "Phase 3 - Dashboard"
status: todo
priority: medium
assigned_to: [frontend, backend]
depends_on: [TASK-001, TASK-003]
created_at: 2026-03-27
---

## Summary

The dashboard is the landing page after login. Currently it shows a basic stats placeholder with hardcoded dashes. The real dashboard must be role-aware: each role sees their own task queue and relevant stats. This is what makes Veylo feel like a purpose-built tool rather than a generic CRUD app.

## Acceptance criteria

### Inspector dashboard
- [ ] "My queue" section: inspections in ENTRY stage assigned to me
- [ ] Quick stats: inspections created today, inspections completed this week
- [ ] "Start new inspection" CTA button

### Evaluator dashboard
- [ ] "Awaiting assessment" section: inspections in EVALUATION stage with unassessed findings
- [ ] Quick stats: findings assessed today, pending assessments count
- [ ] Click any item to go to /inspections/[id]

### Manager dashboard
- [ ] "Pending review" section: inspections in REVIEW stage
- [ ] Quick stats: total inspections this month, completed this month, total cost this month
- [ ] "Recently completed" section: last 5 FINAL inspections

### General
- [ ] Role detected from JWT/user context
- [ ] Loading skeletons for all sections
- [ ] Empty states for each section when no items
- [ ] Responsive: works on tablet (inspectors use tablets)

## API endpoints

- `GET /api/v1/inspections?page=1&page_size=5` -- recent inspections (exists)
- **MISSING:** `GET /api/v1/inspections/stats` -- aggregated stats (counts by stage, costs). Backend enhancement needed.
- **MISSING:** Filter inspections by system_stage and assigned_user. Backend enhancement needed.

## UI components

- Card (exists) for stat cards
- Existing dashboard layout (rewrite current page.tsx)
- Button (exists) for CTAs
- Skeleton for loading states

## Technical notes

- Rewrite `web/src/app/(app)/dashboard/page.tsx` (basic version already exists)
- Create role-specific dashboard components in `web/src/features/dashboard/`
- **Backend dependencies:**
  - Need stats endpoint or add stage/assignee filters to inspections list
  - Need user role available in frontend context (from JWT or user endpoint)
- Current dashboard page already fetches inspections -- extend, do not rewrite from scratch
