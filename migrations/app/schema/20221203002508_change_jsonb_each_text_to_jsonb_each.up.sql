-- POP RAW MIGRATION --
-- audit history table
-- inspired by
-- https://github.com/2ndQuadrant/audit-trigger/blob/master/audit.sql
-- and
-- http://8kb.co.uk/blog/2015/01/19/copying-pavel-stehules-simple-history-table-but-with-the-jsonb-type/

CREATE TABLE IF NOT EXISTS audit_history (
											 id uuid primary key,
											 schema_name text not null,
											 table_name text not null,
											 relid oid not null,
											 object_id uuid,
											 session_userid uuid,
											 event_name text,
											 action_tstamp_tx TIMESTAMP WITH TIME ZONE NOT NULL,
											 action_tstamp_stm TIMESTAMP WITH TIME ZONE NOT NULL,
											 action_tstamp_clk TIMESTAMP WITH TIME ZONE NOT NULL,
											 transaction_id bigint,
											 client_query text,
											 action TEXT NOT NULL CHECK (action IN ('INSERT','DELETE','UPDATE','TRUNCATE')),
											 old_data jsonb,
											 changed_data jsonb,
											 statement_only boolean not null
);

COMMENT ON TABLE audit_history IS 'History of auditable actions on audited tables, from if_modified_func()';
COMMENT ON COLUMN audit_history.id IS 'Unique identifier for each auditable event';
COMMENT ON COLUMN audit_history.schema_name IS 'Name of audited table that this event is in';
COMMENT ON COLUMN audit_history.table_name IS 'Non-schema-qualified table name of table event occured in';
COMMENT ON COLUMN audit_history.relid IS 'Table OID. Changes with drop/create. Get with ''tablename''::regclass';
COMMENT ON COLUMN audit_history.object_id IS 'if the changed data has an id column';
COMMENT ON COLUMN audit_history.session_userid IS 'id of user whose statement caused the audited event';
COMMENT ON COLUMN audit_history.event_name IS 'name of event that caused the audited event';
COMMENT ON COLUMN audit_history.action_tstamp_tx IS 'Transaction start timestamp for tx in which audited event occurred';
COMMENT ON COLUMN audit_history.action_tstamp_stm IS 'Statement start timestamp for tx in which audited event occurred';
COMMENT ON COLUMN audit_history.action_tstamp_clk IS 'Wall clock time at which audited event''s trigger call occurred';
COMMENT ON COLUMN audit_history.transaction_id IS 'Identifier of transaction that made the change. May wrap, but unique paired with action_tstamp_tx.';
COMMENT ON COLUMN audit_history.action IS 'Action type';
COMMENT ON COLUMN audit_history.old_data IS 'Record value. Null for statement-level trigger. For INSERT this is NULL. For DELETE and UPDATE it is the old state of the record stored in json.';
COMMENT ON COLUMN audit_history.changed_data IS 'New values of fields changed by INSERT AND UPDATE. Null except for row-level INSERT and UPDATE events.';
COMMENT ON COLUMN audit_history.statement_only IS 'TRUE if audit event is from an FOR EACH STATEMENT trigger, FALSE for FOR EACH ROW';

CREATE INDEX IF NOT EXISTS audit_history_relid_idx ON audit_history(relid);
CREATE INDEX IF NOT EXISTS audit_history_action_tstamp_tx_stm_idx ON audit_history(action_tstamp_stm);
CREATE INDEX IF NOT EXISTS audit_history_action_idx ON audit_history(action);
CREATE INDEX IF NOT EXISTS audit_history_table_name_idx ON audit_history(table_name);

CREATE OR REPLACE FUNCTION if_modified_func() RETURNS TRIGGER AS $body$
DECLARE
	audit_row audit_history;
	include_values boolean;
	log_diffs boolean;
	j_old jsonb;
	j_new jsonb;
	j_diff jsonb;
	excluded_cols text[] = ARRAY[]::text[];
	-- do NOT require these setting to exist
	_user_id text;
	_event_name text;
BEGIN
	IF TG_WHEN <> 'AFTER' THEN
		RAISE EXCEPTION 'if_modified_func() may only run as an AFTER trigger';
	END IF;

	_event_name := current_setting('audit.current_event_name', true);

	BEGIN
		_user_id := current_setting('audit.current_user_id', true)::uuid;
	EXCEPTION WHEN OTHERS THEN
		_user_id := NULL;
	END;

	audit_row = ROW(
		uuid_generate_v4(),                           -- id
		TG_TABLE_SCHEMA::text,                        -- schema_name
		TG_TABLE_NAME::text,                          -- table_name
		TG_RELID,                                     -- relation OID for much quicker searches
		NULL,                                         -- object id
		_user_id,                                     -- session user_id
		_event_name,                                  -- session event_name
		current_timestamp,                            -- action_tstamp_tx
		statement_timestamp(),                        -- action_tstamp_stm
		clock_timestamp(),                            -- action_tstamp_clk
		txid_current(),                               -- transaction ID
		current_query(),                              -- top-level query or queries if multistatement from client
		TG_OP,                                        -- action
		NULL, NULL,                                   -- old_data, changed_data
		FALSE                                           -- statement_only
		);


	IF NOT TG_ARGV[0]::boolean IS DISTINCT FROM 'f'::boolean THEN
		audit_row.client_query = NULL;
	END IF;

	IF TG_ARGV[1] IS NOT NULL THEN
		excluded_cols = TG_ARGV[1]::text[];
	END IF;

	IF (TG_OP = 'UPDATE' AND TG_LEVEL = 'ROW') THEN
		j_old := row_to_json(OLD)::jsonb;
		j_new := row_to_json(NEW)::jsonb;

		IF j_old ? 'id' THEN
			audit_row.object_id = j_old->>'id';
		END IF;

		audit_row.old_data = j_old - excluded_cols;
		-- inspired by https://stackoverflow.com/a/55852047
		j_diff := (SELECT json_object_agg(COALESCE(oldkv.key, newkv.key), newkv.value)
				   FROM jsonb_each(j_old) oldkv
				   FULL OUTER JOIN jsonb_each(j_new) newkv ON newkv.key = oldkv.key
				   WHERE newkv.value IS DISTINCT FROM oldkv.value);
		audit_row.changed_data = j_diff - excluded_cols;

		-- No fields were changed, skip creating audit record
		IF audit_row.changed_data = jsonb('{}') OR audit_row.changed_data IS NULL THEN
			RETURN NULL;
		END IF;
	ELSIF (TG_OP = 'DELETE' AND TG_LEVEL = 'ROW') THEN
		j_old := row_to_json(OLD)::jsonb;
		IF j_old ? 'id' THEN
			audit_row.object_id = j_old->>'id';
		END IF;

		audit_row.old_data = j_old - excluded_cols;
	ELSIF (TG_OP = 'INSERT' AND TG_LEVEL = 'ROW') THEN
		j_new := row_to_json(NEW)::jsonb;
		IF j_new ? 'id' THEN
			audit_row.object_id = j_new->>'id';
		END IF;
		audit_row.changed_data = j_new - excluded_cols;
	ELSIF (TG_LEVEL = 'STATEMENT' AND TG_OP IN ('INSERT','UPDATE','DELETE','TRUNCATE')) THEN
		audit_row.statement_only = 't';
	ELSE
		RAISE EXCEPTION '[if_modified_func] - Trigger func added as trigger for unhandled case: %, %',TG_OP, TG_LEVEL;
		RETURN NULL;
	END IF;
	INSERT INTO audit_history VALUES (audit_row.*);
	RETURN NULL;
END;
$body$
	LANGUAGE plpgsql
	SECURITY DEFINER
	SET search_path = pg_catalog, public;


COMMENT ON FUNCTION if_modified_func() IS $body$
Track changes to a table at the statement and/or row level.

Optional parameters to trigger in CREATE TRIGGER call:

param 0: boolean, whether to log the query text. Default 't'.

param 1: text[], columns to ignore in updates. Default [].

         Updates to ignored cols are omitted from changed_data.

         Updates with only ignored cols changed or have no changes are not inserted
         into the audit log.

         Almost all the processing work is still done for updates
         that ignored. If you need to save the load, you need to use
         WHEN clause on the trigger instead.

         No warning or error is issued if ignored_cols contains columns
         that do not exist in the target table. This lets you specify
         a standard set of ignored columns.

There is no parameter to disable logging of values. Add this trigger as
a 'FOR EACH STATEMENT' rather than 'FOR EACH ROW' trigger if you do not
want to log row values.

Note that the user name logged is the login role for the session. The audit trigger
cannot obtain the active role because it is reset by the SECURITY DEFINER invocation
of the audit trigger its self.
$body$;
