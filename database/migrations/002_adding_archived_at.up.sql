ALTER TABLE IF EXISTS todo
    ADD COLUMN IF NOT EXISTS archived_at timestamptz;

ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS archived_at timestamptz;

ALTER TABLE IF EXISTS user_session
    ADD COLUMN IF NOT EXISTS archived_at timestamptz;