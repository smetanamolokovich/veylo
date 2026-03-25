---
name: migration
description: Create the next numbered database migration file pair (.up.sql and .down.sql)
argument-hint: <migration-name>
allowed-tools: Glob, Read, Write, Bash
---

Create a new database migration for the Veylo project.

Migration name: $ARGUMENTS

## Steps

1. List existing migrations to find the highest number:
```bash
ls migrations/ | sort
```

2. Next number = highest existing number + 1, zero-padded to 6 digits (e.g. `000011`)

3. Create two files:
   - `migrations/<number>_$ARGUMENTS.up.sql`
   - `migrations/<number>_$ARGUMENTS.down.sql`

## Rules for the SQL

### .up.sql
- Use `TEXT` (not `UUID`) for ID columns that hold ULIDs
- Use `UUID` with `DEFAULT gen_random_uuid()` only for internal join-table IDs
- Always include `organization_id TEXT NOT NULL` for multi-tenant isolation
- Always include `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- Always include `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()` if the entity is mutable
- Use `ON DELETE CASCADE` for foreign keys to parent entities
- Add indexes on `organization_id` and any frequently queried columns
- Soft delete: add `deleted_at TIMESTAMPTZ` if the entity needs audit trail

### .down.sql
- Reverse every operation in .up.sql
- Drop tables in reverse order (child before parent)
- Drop indexes if explicitly created

## Example

For `/migration create_photos_table`:

**000011_create_photos_table.up.sql:**
```sql
CREATE TABLE photos (
    id              TEXT        PRIMARY KEY,
    organization_id TEXT        NOT NULL,
    inspection_id   TEXT        NOT NULL REFERENCES inspections(id) ON DELETE CASCADE,
    s3_key          TEXT        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_photos_inspection_id ON photos(inspection_id);
CREATE INDEX idx_photos_organization_id ON photos(organization_id);
```

**000011_create_photos_table.down.sql:**
```sql
DROP TABLE IF EXISTS photos;
```
