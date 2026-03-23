CREATE TABLE reports (
    id            TEXT        PRIMARY KEY,
    inspection_id TEXT        NOT NULL UNIQUE,
    org_id        TEXT        NOT NULL,
    s3_key        TEXT        NOT NULL,
    url           TEXT        NOT NULL DEFAULT '',
    generated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
