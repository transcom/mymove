ALTER TABLE users
    ADD COLUMN current_session_id TEXT DEFAULT '',
	ADD COLUMN unique_session_id TEXT DEFAULT '';
