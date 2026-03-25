ALTER TABLE organizations
    DROP COLUMN IF EXISTS onboarding_completed_at;

ALTER TABLE users
    ALTER COLUMN organization_id SET NOT NULL;

ALTER TABLE refresh_tokens
    ALTER COLUMN organization_id SET NOT NULL;
