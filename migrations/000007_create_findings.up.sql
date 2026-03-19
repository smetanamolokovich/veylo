CREATE TABLE findings (
    id              VARCHAR(26) PRIMARY KEY,
    inspection_id   VARCHAR(26) NOT NULL REFERENCES inspections(id) ON DELETE CASCADE,
    organization_id VARCHAR(26) NOT NULL,
    body_area       VARCHAR(100) NOT NULL DEFAULT '',
    coordinate_x    DOUBLE PRECISION NOT NULL DEFAULT 0,
    coordinate_y    DOUBLE PRECISION NOT NULL DEFAULT 0,
    type            VARCHAR(100) NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    images          TEXT[] NOT NULL DEFAULT '{}',
    severity        VARCHAR(50),
    repair_method   VARCHAR(50),
    cost_parts      INT NOT NULL DEFAULT 0,
    cost_labor      INT NOT NULL DEFAULT 0,
    cost_paint      INT NOT NULL DEFAULT 0,
    cost_other      INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_findings_inspection_id ON findings (inspection_id);
CREATE INDEX idx_findings_organization_id ON findings (organization_id);
