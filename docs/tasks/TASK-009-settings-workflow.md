---
id: TASK-009
title: "/settings/workflow — Workflow configuration page"
phase: "Phase 4 - Workflow Settings"
status: todo
priority: low
assigned_to: [frontend]
depends_on: []
created_at: 2026-03-27
---

## Summary

The workflow configuration page is Veylo's core competitive advantage. It allows admins to define custom status names, map them to system stages (ENTRY, EVALUATION, REVIEW, FINAL), set an initial status, and define allowed transitions between statuses. This is what makes Veylo sellable to any leasing/fleet company -- they can configure the workflow to match their internal processes.

## Acceptance criteria

### Status management
- [ ] List of all configured statuses with: name, description, system stage, is_initial badge
- [ ] "Add status" form: name (text), description (text), stage (select: ENTRY/EVALUATION/REVIEW/FINAL), is_initial (checkbox)
- [ ] Visual grouping of statuses by system stage (e.g., columns or sections)
- [ ] Only one status can be marked as initial

### Transition management
- [ ] List of configured transitions: from_status -> to_status
- [ ] "Add transition" form: from_status (select), to_status (select)
- [ ] Visual representation of the workflow as a flow/graph (nice-to-have: use a simple node diagram)

### General
- [ ] Only ADMIN role can access this page
- [ ] Success/error feedback on all operations
- [ ] Page loads current workflow on mount
- [ ] If no workflow exists yet, show setup wizard or "Initialize workflow" button

## API endpoints

- `GET /api/v1/workflow` -- get current workflow (statuses + transitions)
- `POST /api/v1/workflow` -- create workflow (if none exists)
- `POST /api/v1/workflow/statuses` -- add status. Body: `{ name, description, stage, is_initial }`
- `POST /api/v1/workflow/transitions` -- add transition. Body: `{ from_status, to_status }`

## UI components

- Card (exists) for status cards
- Select (exists) for stage and status dropdowns
- Input (exists) for name and description
- Button (exists)
- Badge for stage indicators and is_initial flag
- Checkbox (needs to be added) for is_initial
- Dialog for add forms (or inline forms)

## Technical notes

- Create `web/src/app/(app)/settings/workflow/page.tsx`
- Create `web/src/features/workflow/` feature folder with api, types, hooks, components
- This is the most visually complex settings page -- consider a two-column layout: statuses on the left, transitions on the right
- The workflow visualization (flow diagram) is nice-to-have for v1, can be a simple list of transitions initially
- System stages are fixed constants -- display them as non-editable reference
