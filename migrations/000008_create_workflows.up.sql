CREATE TABLE workflows (
    id              UUID PRIMARY KEY,
    organization_id UUID NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workflow_statuses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID        NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    stage       TEXT        NOT NULL CHECK (stage IN ('ENTRY', 'EVALUATION', 'REVIEW', 'FINAL')),
    is_initial  BOOLEAN     NOT NULL DEFAULT FALSE,
    UNIQUE (workflow_id, name)
);

CREATE TABLE workflow_transitions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    from_status TEXT NOT NULL,
    to_status   TEXT NOT NULL,
    UNIQUE (workflow_id, from_status, to_status)
);
