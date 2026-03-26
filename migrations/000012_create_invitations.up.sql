CREATE TABLE invitations (
    id              VARCHAR(26) PRIMARY KEY,
    organization_id VARCHAR(26) NOT NULL REFERENCES organizations(id),
    email           VARCHAR(255) NOT NULL,
    role            VARCHAR(50) NOT NULL,
    token           VARCHAR(64) NOT NULL UNIQUE,
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    expires_at      TIMESTAMPTZ NOT NULL,
    used_at         TIMESTAMPTZ,
    created_by      VARCHAR(26) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invitations_org_id ON invitations (organization_id);
CREATE INDEX idx_invitations_token ON invitations (token);
CREATE UNIQUE INDEX idx_invitations_org_email_pending ON invitations (organization_id, email)
    WHERE status = 'PENDING';
