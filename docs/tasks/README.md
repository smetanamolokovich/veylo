# Veylo -- Task Board

Tasks are created by the `product-manager` agent and managed by `team-lead` throughout the pipeline.

## Status legend

| Status | Meaning |
|--------|---------|
| `todo` | Defined, not started |
| `in_progress` | Agent is actively working |
| `review` | Implemented, awaiting review |
| `done` | Reviewed and accepted |
| `blocked` | Blocked by critical review issues or dependency |

## Tasks

### Phase 1 - Core Loop (HIGH priority)

| ID | Title | Status | Assigned | Depends on |
|----|-------|--------|----------|------------|
| TASK-001 | /inspections -- List page | todo | frontend | -- |
| TASK-002 | /inspections/new -- Create inspection page | todo | frontend | TASK-001 |
| TASK-003 | /inspections/[id] -- Inspection detail page | todo | frontend | TASK-001 |
| TASK-004 | /vehicles/new -- Add vehicle page | todo | frontend | -- |

### Phase 2 - Data Entry (MEDIUM priority)

| ID | Title | Status | Assigned | Depends on |
|----|-------|--------|----------|------------|
| TASK-005 | /vehicles -- Vehicle registry list | todo | frontend, backend | TASK-004 |
| TASK-006 | /vehicles/[id] -- Vehicle detail page | todo | frontend | TASK-005 |
| TASK-007 | /team -- Team management page | todo | frontend, backend | -- |

### Phase 3 - Dashboard (MEDIUM priority)

| ID | Title | Status | Assigned | Depends on |
|----|-------|--------|----------|------------|
| TASK-008 | /dashboard -- Role-aware dashboard | todo | frontend, backend | TASK-001, TASK-003 |

### Phase 4 - Workflow Settings (LOW priority)

| ID | Title | Status | Assigned | Depends on |
|----|-------|--------|----------|------------|
| TASK-009 | /settings/workflow -- Workflow configuration | todo | frontend | -- |
| TASK-010 | /settings -- Organization settings | todo | frontend, backend | -- |

## Backend Dependencies Identified

Several frontend tasks require backend endpoints that do not exist yet:

| Missing endpoint | Needed by | Description |
|-----------------|-----------|-------------|
| `GET /api/v1/assets?type=vehicle` | TASK-002, TASK-005 | List/search assets with pagination |
| `GET /api/v1/organizations/me/members` | TASK-007 | List organization team members |
| `GET /api/v1/organizations/me/invitations` | TASK-007 | List pending invitations |
| `GET /api/v1/inspections/stats` | TASK-008 | Aggregated inspection statistics |
| `GET /api/v1/inspections?stage=X&assignee=Y` | TASK-008 | Filter inspections by stage and assignee |
| `PUT /api/v1/organizations/me` | TASK-010 | Update organization details |
