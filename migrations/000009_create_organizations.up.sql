CREATE TABLE organizations (
    id         TEXT        PRIMARY KEY,
    name       TEXT        NOT NULL,
    vertical   TEXT        NOT NULL CHECK (vertical IN ('VEHICLE', 'PROPERTY')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
