ALTER TABLE users
    ADD COLUMN current_mil_session_id TEXT DEFAULT '',
    ADD COLUMN current_admin_session_id TEXT DEFAULT '',
    ADD COLUMN current_office_session_id TEXT DEFAULT '';
CREATE INDEX ON users (current_mil_session_id);
CREATE INDEX ON users (current_admin_session_id);
CREATE INDEX ON users (current_office_session_id);
