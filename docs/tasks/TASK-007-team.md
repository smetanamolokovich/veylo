---
id: TASK-007
title: "/team — Team management page"
phase: "Phase 2 - Data Entry"
status: todo
priority: medium
assigned_to: [frontend, backend]
depends_on: []
created_at: 2026-03-27
---

## Summary

Team management page for admins and managers. Shows all members of the organization with their roles. Allows inviting new team members via email. This is essential for multi-user organizations -- a leasing company might have 5-20 inspectors, 2-5 evaluators, and 1-3 managers.

## Acceptance criteria

- [ ] List of all team members with: name, email, role, status (active/pending)
- [ ] "Invite member" button opens a form with email + role selector
- [ ] Role options: INSPECTOR, EVALUATOR, MANAGER, ADMIN
- [ ] Pending invitations shown with "pending" badge
- [ ] Success message after sending invitation
- [ ] Error handling: duplicate email, invalid email
- [ ] Only ADMIN and MANAGER roles can access this page

## API endpoints

- `POST /api/v1/organizations/me/invitations` -- send invitation. Body: `{ email, role }`
- **MISSING:** `GET /api/v1/organizations/me/members` -- list team members. This endpoint does NOT exist yet. Backend work required.
- **MISSING:** `GET /api/v1/organizations/me/invitations` -- list pending invitations.

## UI components

- Table for members list
- Dialog for invite form
- Select (exists) for role picker
- Input (exists) for email
- Badge for role and status display
- Button (exists)

## Technical notes

- Create `web/src/app/(app)/team/page.tsx`
- Create `web/src/features/team/` feature folder
- Reuse existing `web/src/features/invitations/hooks/use-invite-user.ts`
- **Backend dependency:** Need list members and list invitations endpoints. Currently only POST invitation exists.
- Access control: redirect non-admin/manager users away from this page
