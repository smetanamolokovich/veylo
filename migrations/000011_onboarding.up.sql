ALTER TABLE organizations
    ADD COLUMN onboarding_completed_at TIMESTAMPTZ;

ALTER TABLE users
    ALTER COLUMN organization_id DROP NOT NULL;

-- Drop and recreate the unique index; PostgreSQL treats NULLs as distinct
-- in unique indexes so (email, NULL) won't conflict with (email, NULL).
DROP INDEX IF EXISTS idx_users_email_org;
CREATE UNIQUE INDEX idx_users_email_org ON users (email, organization_id);

ALTER TABLE refresh_tokens
    ALTER COLUMN organization_id DROP NOT NULL;
