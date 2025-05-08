ALTER TABLE sessions ADD COLUMN token VARCHAR(255) NOT NULL;

CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions (token);