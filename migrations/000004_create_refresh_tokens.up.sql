CREATE TABLE
    refresh_tokens (
        id VARCHAR(26) PRIMARY KEY,
        user_id VARCHAR(26) NOT NULL,
        organization_id VARCHAR(26) NOT NULL,
        token_hash VARCHAR(255) NOT NULL,
        expires_at TIMESTAMPTZ NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);