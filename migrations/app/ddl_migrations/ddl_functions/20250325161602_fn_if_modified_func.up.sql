-- cam b-22911. When inserting Beth's new DB proc we add a seq_num which can't be NULL
-- as it's a new serial column. The old if_modified_func trigger is inserting NULL
-- in the seq_num, so this migration is the exact same if_modified_func code,
-- but just omitting inserting a null seq_num.
-- Also adjusted the row declaration, I was getting a bunch of syntax
-- errors before wrapping it in plpgsql and manually declaring the insert values
CREATE OR REPLACE FUNCTION public.if_modified_func()
 RETURNS trigger
 LANGUAGE plpgsql
 SECURITY DEFINER
 SET search_path TO 'pg_catalog', 'public'
AS $function$
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
		FALSE,                                        -- statement_only
		NULL										  --seq_num
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

		IF jsonb_exists(j_old, 'id') THEN
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
		IF jsonb_exists(j_old, 'id') THEN
			audit_row.object_id = j_old->>'id';
		END IF;

		audit_row.old_data = j_old - excluded_cols;
	ELSIF (TG_OP = 'INSERT' AND TG_LEVEL = 'ROW') THEN
		j_new := row_to_json(NEW)::jsonb;
		IF jsonb_exists(j_new, 'id') THEN
			audit_row.object_id = j_new->>'id';
		END IF;
		audit_row.changed_data = j_new - excluded_cols;
	ELSIF (TG_LEVEL = 'STATEMENT' AND TG_OP IN ('INSERT','UPDATE','DELETE','TRUNCATE')) THEN
		audit_row.statement_only = 't';
	ELSE
		RAISE EXCEPTION '[if_modified_func] - Trigger func added as trigger for unhandled case: %, %',TG_OP, TG_LEVEL;
		RETURN NULL;
	END IF;

	INSERT INTO audit_history
		(id, schema_name, table_name, relid, object_id, session_userid, event_name, action_tstamp_tx, action_tstamp_stm, action_tstamp_clk, transaction_id, client_query, "action", old_data, changed_data, statement_only)
	VALUES(
		audit_row.id, audit_row.schema_name, audit_row.table_name, audit_row.relid, audit_row.object_id,
		audit_row.session_userid, audit_row.event_name, audit_row.action_tstamp_tx, audit_row.action_tstamp_stm,
		audit_row.action_tstamp_clk, audit_row.transaction_id, audit_row.client_query, audit_row.action,
		audit_row.old_data, audit_row.changed_data, audit_row.statement_only);

	RETURN NULL;
END;
$function$
;
