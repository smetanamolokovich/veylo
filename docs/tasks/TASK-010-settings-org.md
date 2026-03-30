---
id: TASK-010
title: "/settings — Organization settings page"
phase: "Phase 4 - Workflow Settings"
status: todo
priority: low
assigned_to: [frontend, backend]
depends_on: []
created_at: 2026-03-27
---

## Summary

Organization settings page for admins to view and edit organization details: name, vertical, and other metadata. This is a relatively simple page but necessary for complete admin functionality.

## Acceptance criteria

- [ ] Display current organization info: name, vertical, onboarding date
- [ ] Edit organization name (inline edit or form)
- [ ] Vertical displayed as read-only (changing vertical is a major operation)
- [ ] Only ADMIN role can edit; other roles see read-only view
- [ ] Success/error feedback on save
- [ ] Navigation sidebar in /settings with links to: General (this page), Workflow (/settings/workflow)

## API endpoints

- `GET /api/v1/organizations/me` -- get current organization (exists)
- **MISSING:** `PUT /api/v1/organizations/me` -- update organization. Backend enhancement needed.

## UI components

- Card (exists) for settings sections
- Input (exists) for editable fields
- Button (exists) for save
- Navigation/tabs for settings sub-pages

## Technical notes

- Create `web/src/app/(app)/settings/page.tsx`
- Create `web/src/app/(app)/settings/layout.tsx` -- shared settings layout with sidebar navigation
- Create `web/src/features/organization/components/org-settings-form.tsx`
- **Backend dependency:** Need PUT endpoint for updating organization. Currently only GET me and POST create exist.
- The settings layout should be reusable for future settings pages (billing, integrations, etc.)
