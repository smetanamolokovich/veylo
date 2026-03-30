# Veylo Frontend Development Roadmap

Created: 2026-03-27

## Overview

This roadmap covers the 10 missing frontend pages needed to complete the Veylo MVP. The backend is fully built -- all core API endpoints exist. The frontend currently has auth flow (login, signup, onboarding, invitation acceptance) and a placeholder dashboard.

## Current State

### Built (frontend)
- /login, /signup -- authentication
- /onboarding -- organization creation + team invite
- /invite/[token] -- accept invitation
- /dashboard -- basic placeholder (shows total count + recent inspections)

### Built (backend API)
- Auth: signup, register, login, refresh
- Organizations: create, get me, complete onboarding, send invitation
- Invitations: get by token, accept
- Inspections: create, list (paginated), get by ID, transition status, get report
- Findings: create, list by inspection, assess
- Assets: create vehicle, get by ID
- Workflow: create, get, add status, add transition

### UI components available
- Button, Input, Label, Select, Card

## Phase Plan

### Phase 1 - Core Loop (Week 1-2) -- HIGH priority

The inspection lifecycle is the product. Without these pages, Veylo cannot be used.

| Task | Page | Effort estimate |
|------|------|----------------|
| TASK-001 | /inspections (list) | 1 day |
| TASK-002 | /inspections/new (create) | 1 day |
| TASK-003 | /inspections/[id] (detail + findings + transitions) | 3-4 days |
| TASK-004 | /vehicles/new (add vehicle, reusable form) | 0.5 day |

**Total: ~6 days**

TASK-003 is the most complex page in the entire app. It handles three different user roles and four inspection stages. Recommend breaking it into sub-tasks:
- 3a: Layout + inspection header + status display
- 3b: Findings list + create finding form (Inspector view)
- 3c: Assessment form per finding (Evaluator view)
- 3d: Review summary + transition buttons (Manager view)
- 3e: Report download (FINAL stage)

### Phase 2 - Data Entry (Week 3) -- MEDIUM priority

Supporting pages that make the app complete but are not on the critical path.

| Task | Page | Effort estimate |
|------|------|----------------|
| TASK-005 | /vehicles (list) | 1 day (+ backend: list assets endpoint) |
| TASK-006 | /vehicles/[id] (detail) | 1 day |
| TASK-007 | /team (management) | 1 day (+ backend: list members/invitations) |

**Total: ~3 days + 1 day backend**

### Phase 3 - Dashboard (Week 4) -- MEDIUM priority

Transform the placeholder dashboard into a role-aware command center.

| Task | Page | Effort estimate |
|------|------|----------------|
| TASK-008 | /dashboard (role-aware) | 2 days (+ backend: stats endpoint, filters) |

**Total: ~2 days + 1 day backend**

### Phase 4 - Workflow Settings (Week 4-5) -- LOW priority

The workflow configuration is Veylo's competitive differentiator, but it's admin-only and not needed for day-to-day operations. Most orgs will be set up with a default workflow during onboarding.

| Task | Page | Effort estimate |
|------|------|----------------|
| TASK-009 | /settings/workflow | 2 days |
| TASK-010 | /settings (org info) | 0.5 day (+ backend: update org endpoint) |

**Total: ~2.5 days + 0.5 day backend**

## Backend Work Required

6 missing endpoints identified across all phases:

1. **GET /api/v1/assets** -- List assets with pagination and type filter. Needed for vehicle list page and vehicle search in inspection creation. (Phase 1 blocker for search, Phase 2 for list page)
2. **GET /api/v1/organizations/me/members** -- List team members. (Phase 2)
3. **GET /api/v1/organizations/me/invitations** -- List pending invitations. (Phase 2)
4. **GET /api/v1/inspections/stats** -- Aggregated stats by stage. (Phase 3)
5. **GET /api/v1/inspections** with stage/assignee filters -- Enhanced list endpoint. (Phase 3)
6. **PUT /api/v1/organizations/me** -- Update org details. (Phase 4)

## New shadcn/ui Components Needed

These components need to be added to `web/src/components/ui/`:

- **Table** -- used in inspections list, vehicles list, team list
- **Badge** -- used everywhere for status, severity, role display
- **Skeleton** -- loading states
- **Dialog** -- modals for forms (add finding, invite member, quick-add vehicle)
- **Tabs** -- inspection detail sections, settings navigation
- **Checkbox** -- workflow settings (is_initial)
- **Sheet** -- side panel alternative to dialog for finding forms on mobile
- **Combobox/Autocomplete** -- vehicle search in inspection creation

## Total Effort Estimate

| Phase | Frontend | Backend | Total |
|-------|----------|---------|-------|
| Phase 1 | 6 days | 0 days | 6 days |
| Phase 2 | 3 days | 1 day | 4 days |
| Phase 3 | 2 days | 1 day | 3 days |
| Phase 4 | 2.5 days | 0.5 day | 3 days |
| **Total** | **13.5 days** | **2.5 days** | **16 days** |

With a single developer, the full frontend MVP is achievable in 3-4 weeks.
