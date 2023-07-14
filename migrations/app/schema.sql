--
-- PostgreSQL database dump
--

-- Dumped from database version 12.13
-- Dumped by pg_dump version 12.13

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: btree_gist; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS btree_gist WITH SCHEMA public;


--
-- Name: EXTENSION btree_gist; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION btree_gist IS 'support for indexing common datatypes in GiST';


--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


--
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


--
-- Name: unaccent; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS unaccent WITH SCHEMA public;


--
-- Name: EXTENSION unaccent; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION unaccent IS 'text search dictionary that removes accents';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: admin_role; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.admin_role AS ENUM (
    'SYSTEM_ADMIN',
    'PROGRAM_ADMIN'
);


ALTER TYPE public.admin_role OWNER TO postgres;

--
-- Name: category_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.category_type AS ENUM (
    'Pre-Move Services',
    'Physical Move Services',
    'Liability'
);


ALTER TYPE public.category_type OWNER TO postgres;

--
-- Name: customer_contact_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.customer_contact_type AS ENUM (
    'FIRST',
    'SECOND'
);


ALTER TYPE public.customer_contact_type OWNER TO postgres;

--
-- Name: destination_address_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.destination_address_type AS ENUM (
    'HOME_OF_RECORD',
    'HOME_OF_SELECTION',
    'PLACE_ENTERED_ACTIVE_DUTY',
    'OTHER_THAN_AUTHORIZED'
);


ALTER TYPE public.destination_address_type OWNER TO postgres;

--
-- Name: TYPE destination_address_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TYPE public.destination_address_type IS 'List of possible destination address types';


--
-- Name: dimension_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.dimension_type AS ENUM (
    'ITEM',
    'CRATE'
);


ALTER TYPE public.dimension_type OWNER TO postgres;

--
-- Name: edi_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.edi_type AS ENUM (
    '997',
    '824',
    '810',
    '858'
);


ALTER TYPE public.edi_type OWNER TO postgres;

--
-- Name: evaluation_location_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.evaluation_location_type AS ENUM (
    'ORIGIN',
    'DESTINATION',
    'OTHER'
);


ALTER TYPE public.evaluation_location_type OWNER TO postgres;

--
-- Name: evaluation_report_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.evaluation_report_type AS ENUM (
    'SHIPMENT',
    'COUNSELING'
);


ALTER TYPE public.evaluation_report_type OWNER TO postgres;

--
-- Name: ghc_approval_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.ghc_approval_status AS ENUM (
    'APPROVED',
    'DRAFT',
    'REJECTED'
);


ALTER TYPE public.ghc_approval_status OWNER TO postgres;

--
-- Name: inspection_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.inspection_type AS ENUM (
    'DATA_REVIEW',
    'PHYSICAL',
    'VIRTUAL'
);


ALTER TYPE public.inspection_type OWNER TO postgres;

--
-- Name: loa_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.loa_type AS ENUM (
    'HHG',
    'NTS'
);


ALTER TYPE public.loa_type OWNER TO postgres;

--
-- Name: move_task_order_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.move_task_order_type AS ENUM (
    'prime',
    'non_temporary_storage'
);


ALTER TYPE public.move_task_order_type OWNER TO postgres;

--
-- Name: moving_expense_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.moving_expense_type AS ENUM (
    'CONTRACTED_EXPENSE',
    'OIL',
    'PACKING_MATERIALS',
    'RENTAL_EQUIPMENT',
    'STORAGE',
    'TOLLS',
    'WEIGHING_FEE',
    'OTHER'
);


ALTER TYPE public.moving_expense_type OWNER TO postgres;

--
-- Name: mto_agents_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.mto_agents_type AS ENUM (
    'RELEASING_AGENT',
    'RECEIVING_AGENT'
);


ALTER TYPE public.mto_agents_type OWNER TO postgres;

--
-- Name: mto_shipment_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.mto_shipment_status AS ENUM (
    'SUBMITTED',
    'APPROVED',
    'REJECTED',
    'DRAFT',
    'CANCELLATION_REQUESTED',
    'CANCELED',
    'DIVERSION_REQUESTED'
);


ALTER TYPE public.mto_shipment_status OWNER TO postgres;

--
-- Name: mto_shipment_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.mto_shipment_type AS ENUM (
    'HHG',
    'INTERNATIONAL_HHG',
    'INTERNATIONAL_UB',
    'HHG_INTO_NTS_DOMESTIC',
    'HHG_OUTOF_NTS_DOMESTIC',
    'MOTORHOME',
    'BOAT_HAUL_AWAY',
    'BOAT_TOW_AWAY',
    'PPM'
);


ALTER TYPE public.mto_shipment_type OWNER TO postgres;

--
-- Name: order_types; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.order_types AS ENUM (
    'GHC',
    'NTS'
);


ALTER TYPE public.order_types OWNER TO postgres;

--
-- Name: payment_request_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.payment_request_status AS ENUM (
    'PENDING',
    'REVIEWED',
    'SENT_TO_GEX',
    'RECEIVED_BY_GEX',
    'PAID',
    'REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED',
    'EDI_ERROR',
    'DEPRECATED'
);


ALTER TYPE public.payment_request_status OWNER TO postgres;

--
-- Name: payment_service_item_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.payment_service_item_status AS ENUM (
    'REQUESTED',
    'APPROVED',
    'DENIED',
    'SENT_TO_GEX',
    'PAID'
);


ALTER TYPE public.payment_service_item_status OWNER TO postgres;

--
-- Name: ppm_advance_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.ppm_advance_status AS ENUM (
    'APPROVED',
    'EDITED',
    'REJECTED'
);


ALTER TYPE public.ppm_advance_status OWNER TO postgres;

--
-- Name: ppm_document_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.ppm_document_status AS ENUM (
    'APPROVED',
    'EXCLUDED',
    'REJECTED'
);


ALTER TYPE public.ppm_document_status OWNER TO postgres;

--
-- Name: ppm_shipment_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.ppm_shipment_status AS ENUM (
    'DRAFT',
    'SUBMITTED',
    'WAITING_ON_CUSTOMER',
    'NEEDS_ADVANCE_APPROVAL',
    'NEEDS_PAYMENT_APPROVAL',
    'PAYMENT_APPROVED'
);


ALTER TYPE public.ppm_shipment_status OWNER TO postgres;

--
-- Name: progear_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.progear_status AS ENUM (
    'YES',
    'NO',
    'NOT SURE'
);


ALTER TYPE public.progear_status OWNER TO postgres;

--
-- Name: reweigh_requester; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.reweigh_requester AS ENUM (
    'CUSTOMER',
    'PRIME',
    'SYSTEM',
    'TOO'
);


ALTER TYPE public.reweigh_requester OWNER TO postgres;

--
-- Name: service_item_param_origin; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.service_item_param_origin AS ENUM (
    'PRIME',
    'SYSTEM',
    'PRICER',
    'PAYMENT_REQUEST'
);


ALTER TYPE public.service_item_param_origin OWNER TO postgres;

--
-- Name: service_item_param_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.service_item_param_type AS ENUM (
    'STRING',
    'DATE',
    'INTEGER',
    'DECIMAL',
    'TIMESTAMP',
    'PaymentServiceItemUUID',
    'BOOLEAN'
);


ALTER TYPE public.service_item_param_type OWNER TO postgres;

--
-- Name: service_item_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.service_item_status AS ENUM (
    'SUBMITTED',
    'APPROVED',
    'REJECTED'
);


ALTER TYPE public.service_item_status OWNER TO postgres;

--
-- Name: shipment_address_update_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.shipment_address_update_status AS ENUM (
    'REQUESTED',
    'REJECTED',
    'APPROVED'
);


ALTER TYPE public.shipment_address_update_status OWNER TO postgres;

--
-- Name: sit_address_update_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.sit_address_update_status AS ENUM (
    'REQUESTED',
    'REJECTED',
    'APPROVED'
);


ALTER TYPE public.sit_address_update_status OWNER TO postgres;

--
-- Name: sit_extension_request_reason; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.sit_extension_request_reason AS ENUM (
    'SERIOUS_ILLNESS_MEMBER',
    'SERIOUS_ILLNESS_DEPENDENT',
    'IMPENDING_ASSIGNEMENT',
    'DIRECTED_TEMPORARY_DUTY',
    'NONAVAILABILITY_OF_CIVILIAN_HOUSING',
    'AWAITING_COMPLETION_OF_RESIDENCE',
    'OTHER'
);


ALTER TYPE public.sit_extension_request_reason OWNER TO postgres;

--
-- Name: TYPE sit_extension_request_reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TYPE public.sit_extension_request_reason IS 'List of reasons a SIT extension can be requested for';


--
-- Name: sit_extension_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.sit_extension_status AS ENUM (
    'PENDING',
    'APPROVED',
    'DENIED'
);


ALTER TYPE public.sit_extension_status OWNER TO postgres;

--
-- Name: TYPE sit_extension_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TYPE public.sit_extension_status IS 'List of possible statuses for a SIT Extension';


--
-- Name: sit_location_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.sit_location_type AS ENUM (
    'ORIGIN',
    'DESTINATION'
);


ALTER TYPE public.sit_location_type OWNER TO postgres;

--
-- Name: TYPE sit_location_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TYPE public.sit_location_type IS 'The type of location for the PPM''s SIT.';


--
-- Name: sub_category_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.sub_category_type AS ENUM (
    'Customer Support',
    'Counseling',
    'Weight Estimate',
    'Additional Services',
    'Inventory & Documentation',
    'Packing/Unpacking',
    'Shipment Schedule',
    'Shipment Weights',
    'Storage',
    'Workforce/Sub-Contractor Management',
    'Loss & Damage',
    'Inconvenience & Hardship Claims'
);


ALTER TYPE public.sub_category_type OWNER TO postgres;

--
-- Name: upload_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.upload_type AS ENUM (
    'PRIME',
    'USER'
);


ALTER TYPE public.upload_type OWNER TO postgres;

--
-- Name: webhook_notifications_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.webhook_notifications_status AS ENUM (
    'PENDING',
    'SENT',
    'FAILED',
    'FAILING',
    'SKIPPED'
);


ALTER TYPE public.webhook_notifications_status OWNER TO postgres;

--
-- Name: webhook_subscriptions_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.webhook_subscriptions_status AS ENUM (
    'ACTIVE',
    'DISABLED',
    'FAILING'
);


ALTER TYPE public.webhook_subscriptions_status OWNER TO postgres;

--
-- Name: add_audit_history_table(regclass); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.add_audit_history_table(target_table regclass) RETURNS void
    LANGUAGE sql
    AS $_$
SELECT add_audit_history_table($1, BOOLEAN 't', BOOLEAN 't');
$_$;


ALTER FUNCTION public.add_audit_history_table(target_table regclass) OWNER TO postgres;

--
-- Name: FUNCTION add_audit_history_table(target_table regclass); Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON FUNCTION public.add_audit_history_table(target_table regclass) IS '
Add auditing support to the given table. Row-level changes will be logged with full client query text. No cols are ignored.
';


--
-- Name: add_audit_history_table(regclass, boolean, boolean); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean) RETURNS void
    LANGUAGE sql
    AS $_$
SELECT add_audit_history_table($1, $2, $3, ARRAY[]::text[]);
$_$;


ALTER FUNCTION public.add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean) OWNER TO postgres;

--
-- Name: add_audit_history_table(regclass, boolean, boolean, text[]); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean, ignored_cols text[]) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE
  stm_targets text = 'INSERT OR UPDATE OR DELETE OR TRUNCATE';
  _q_txt text;
  _ignored_cols_snip text = '';
BEGIN
    EXECUTE 'DROP TRIGGER IF EXISTS audit_trigger_row ON ' || target_table;
    EXECUTE 'DROP TRIGGER IF EXISTS audit_trigger_stm ON ' || target_table;

    IF audit_rows THEN
        IF array_length(ignored_cols,1) > 0 THEN
            _ignored_cols_snip = ', ' || quote_literal(ignored_cols);
        END IF;
        _q_txt = 'CREATE TRIGGER audit_trigger_row AFTER INSERT OR UPDATE OR DELETE ON ' ||
                 target_table ||
                 ' FOR EACH ROW EXECUTE PROCEDURE if_modified_func(' ||
                 quote_literal(audit_query_text) || _ignored_cols_snip || ');';
        RAISE NOTICE '%',_q_txt;
        EXECUTE _q_txt;
        stm_targets = 'TRUNCATE';
    ELSE
    END IF;

    _q_txt = 'CREATE TRIGGER audit_trigger_stm AFTER ' || stm_targets || ' ON ' ||
             target_table ||
             ' FOR EACH STATEMENT EXECUTE PROCEDURE if_modified_func('||
             quote_literal(audit_query_text) || ');';
    RAISE NOTICE '%',_q_txt;
    EXECUTE _q_txt;

END;
$$;


ALTER FUNCTION public.add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean, ignored_cols text[]) OWNER TO postgres;

--
-- Name: FUNCTION add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean, ignored_cols text[]); Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON FUNCTION public.add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean, ignored_cols text[]) IS '
Add auditing support to a table.

Arguments:
   target_table:     Table name, schema qualified if not on search_path
   audit_rows:       Record each row change, or only audit at a statement level
   audit_query_text: Record the text of the client query that triggered the audit event?
   ignored_cols:     Columns to exclude from update diffs, ignore updates that change only ignored cols.
';


--
-- Name: f_unaccent(text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.f_unaccent(text) RETURNS text
    LANGUAGE sql IMMUTABLE STRICT PARALLEL SAFE
    AS $_$
SELECT public.unaccent('public.unaccent', $1)  -- schema-qualify function and dictionary
$_$;


ALTER FUNCTION public.f_unaccent(text) OWNER TO postgres;

--
-- Name: FUNCTION f_unaccent(text); Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON FUNCTION public.f_unaccent(text) IS 'Wrapper around unaccent that is marked as immutable so it can be used in indexes';


--
-- Name: if_modified_func(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.if_modified_func() RETURNS trigger
    LANGUAGE plpgsql SECURITY DEFINER
    SET search_path TO 'pg_catalog', 'public'
    AS $$
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
$$;


ALTER FUNCTION public.if_modified_func() OWNER TO postgres;

--
-- Name: FUNCTION if_modified_func(); Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON FUNCTION public.if_modified_func() IS '
Track changes to a table at the statement and/or row level.

Optional parameters to trigger in CREATE TRIGGER call:

param 0: boolean, whether to log the query text. Default ''t''.

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
a ''FOR EACH STATEMENT'' rather than ''FOR EACH ROW'' trigger if you do not
want to log row values.

Note that the user name logged is the login role for the session. The audit trigger
cannot obtain the active role because it is reset by the SECURITY DEFINER invocation
of the audit trigger its self.
';


--
-- Name: searchable_full_name(text, text); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.searchable_full_name(first_name text, last_name text) RETURNS text
    LANGUAGE sql IMMUTABLE STRICT
    AS $$
        -- CONCAT_WS is immutable when given only text arguments
SELECT public.f_unaccent(LOWER(CONCAT_WS(' ', first_name, last_name)));
$$;


ALTER FUNCTION public.searchable_full_name(first_name text, last_name text) OWNER TO postgres;

--
-- Name: FUNCTION searchable_full_name(first_name text, last_name text); Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON FUNCTION public.searchable_full_name(first_name text, last_name text) IS 'Prepares a first/last name for search by lowercasing and removing diacritics';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: addresses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.addresses (
    id uuid NOT NULL,
    street_address_1 character varying(255) NOT NULL,
    street_address_2 character varying(255),
    city character varying(255) NOT NULL,
    state character varying(255) NOT NULL,
    postal_code character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    street_address_3 character varying(255),
    country character varying(255) DEFAULT 'United States'::character varying
);


ALTER TABLE public.addresses OWNER TO postgres;

--
-- Name: TABLE addresses; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.addresses IS 'Holds all address information';


--
-- Name: COLUMN addresses.street_address_1; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.addresses.street_address_1 IS 'First street address value for address record';


--
-- Name: COLUMN addresses.street_address_2; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.addresses.street_address_2 IS 'Second street address value for address record';


--
-- Name: COLUMN addresses.city; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.addresses.city IS 'City value for address record';


--
-- Name: COLUMN addresses.state; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.addresses.state IS 'State value for address record';


--
-- Name: COLUMN addresses.postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.addresses.postal_code IS 'Postal code value for address record';


--
-- Name: COLUMN addresses.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.addresses.created_at IS 'Date & time the address was created';


--
-- Name: COLUMN addresses.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.addresses.updated_at IS 'Date & time the address was last updated';


--
-- Name: COLUMN addresses.street_address_3; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.addresses.street_address_3 IS 'Third street address value for address record';


--
-- Name: COLUMN addresses.country; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.addresses.country IS 'Country address value for address record';


--
-- Name: admin_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.admin_users (
    id uuid NOT NULL,
    user_id uuid,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    organization_id uuid,
    role public.admin_role NOT NULL,
    email character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    active boolean DEFAULT false NOT NULL
);


ALTER TABLE public.admin_users OWNER TO postgres;

--
-- Name: TABLE admin_users; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.admin_users IS 'Holds all users who have access to the admin interface, where one can perform CRUD operations on entities such as office users and admin users. Individual authenticated sessions can also be revoked via the admin interface.';


--
-- Name: COLUMN admin_users.user_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_users.user_id IS 'The foreign key that points to the user id in the users table';


--
-- Name: COLUMN admin_users.first_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_users.first_name IS 'The first name of the admin user';


--
-- Name: COLUMN admin_users.last_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_users.last_name IS 'The last name of the admin user';


--
-- Name: COLUMN admin_users.organization_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_users.organization_id IS 'The foreign key that points to the organization id in the organizations table. Truss admin users belong to the Truss organization.';


--
-- Name: COLUMN admin_users.role; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_users.role IS 'An enum with two possible values: SYSTEM_ADMIN or PROGRAM_ADMIN. Note that PROGRAM_ADMIN is no longer used and there is a JIRA story to remove it.';


--
-- Name: COLUMN admin_users.email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_users.email IS 'The email of the admin user';


--
-- Name: COLUMN admin_users.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_users.created_at IS 'Date & time the admin user was created';


--
-- Name: COLUMN admin_users.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_users.updated_at IS 'Date & time the admin user was updated';


--
-- Name: COLUMN admin_users.active; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.admin_users.active IS 'A boolean that determines whether or not an admin user is active. Users that are not active are not allowed to access the admin site. See https://github.com/transcom/mymove/wiki/create-or-deactivate-users.';


--
-- Name: archived_access_codes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.archived_access_codes (
    id uuid NOT NULL,
    service_member_id uuid,
    code text NOT NULL,
    move_type text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    claimed_at timestamp with time zone
);


ALTER TABLE public.archived_access_codes OWNER TO postgres;

--
-- Name: archived_move_documents; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.archived_move_documents (
    id uuid NOT NULL,
    move_id uuid NOT NULL,
    document_id uuid NOT NULL,
    move_document_type character varying(255) NOT NULL,
    status character varying(255) NOT NULL,
    notes text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    title character varying(255) NOT NULL,
    personally_procured_move_id uuid,
    deleted_at timestamp with time zone
);


ALTER TABLE public.archived_move_documents OWNER TO postgres;

--
-- Name: archived_moving_expense_documents; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.archived_moving_expense_documents (
    id uuid NOT NULL,
    move_document_id uuid NOT NULL,
    moving_expense_type character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    requested_amount_cents integer NOT NULL,
    payment_method character varying(255) NOT NULL,
    receipt_missing boolean DEFAULT false NOT NULL,
    storage_start_date date,
    storage_end_date date,
    deleted_at timestamp with time zone
);


ALTER TABLE public.archived_moving_expense_documents OWNER TO postgres;

--
-- Name: archived_personally_procured_moves; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.archived_personally_procured_moves (
    id uuid NOT NULL,
    move_id uuid NOT NULL,
    size character varying(255),
    weight_estimate integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    pickup_postal_code character varying(255),
    additional_pickup_postal_code character varying(255),
    destination_postal_code character varying(255),
    days_in_storage integer,
    status character varying(255) DEFAULT 'DRAFT'::character varying NOT NULL,
    has_additional_postal_code boolean,
    has_sit boolean,
    has_requested_advance boolean DEFAULT false NOT NULL,
    advance_id uuid,
    estimated_storage_reimbursement character varying(255),
    mileage integer,
    planned_sit_max integer,
    sit_max integer,
    incentive_estimate_min integer,
    incentive_estimate_max integer,
    advance_worksheet_id uuid,
    net_weight integer,
    original_move_date date,
    actual_move_date date,
    total_sit_cost integer,
    submit_date timestamp with time zone,
    approve_date timestamp with time zone,
    reviewed_date timestamp with time zone,
    has_pro_gear public.progear_status,
    has_pro_gear_over_thousand public.progear_status
);


ALTER TABLE public.archived_personally_procured_moves OWNER TO postgres;

--
-- Name: archived_reimbursements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.archived_reimbursements (
    id uuid NOT NULL,
    requested_amount integer NOT NULL,
    method_of_receipt character varying(255) NOT NULL,
    status character varying(255) NOT NULL,
    requested_date date,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.archived_reimbursements OWNER TO postgres;

--
-- Name: TABLE archived_reimbursements; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.archived_reimbursements IS 'Holds information about reimbursements to a customer';


--
-- Name: COLUMN archived_reimbursements.requested_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.archived_reimbursements.requested_amount IS 'The reimbursement amount the customer is requesting in cents';


--
-- Name: COLUMN archived_reimbursements.method_of_receipt; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.archived_reimbursements.method_of_receipt IS 'The way the customer wants to be reimbursed: OTHER (any other payment type other than GTCC), GTCC (Govt travel charge card)';


--
-- Name: COLUMN archived_reimbursements.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.archived_reimbursements.status IS 'The current status of the reimbursement: DRAFT, REQUESTED, APPROVED, REJECTED, PAID';


--
-- Name: COLUMN archived_reimbursements.requested_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.archived_reimbursements.requested_date IS 'Date the reimbursement was requested';


--
-- Name: COLUMN archived_reimbursements.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.archived_reimbursements.created_at IS 'Date & time the reimbursement was created';


--
-- Name: COLUMN archived_reimbursements.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.archived_reimbursements.updated_at IS 'Date & time the reimbursement was last updated';


--
-- Name: archived_signed_certifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.archived_signed_certifications (
    id uuid NOT NULL,
    submitting_user_id uuid NOT NULL,
    move_id uuid NOT NULL,
    certification_text text NOT NULL,
    signature text NOT NULL,
    date timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    personally_procured_move_id uuid,
    certification_type text
);


ALTER TABLE public.archived_signed_certifications OWNER TO postgres;

--
-- Name: archived_weight_ticket_set_documents; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.archived_weight_ticket_set_documents (
    id uuid NOT NULL,
    weight_ticket_set_type text NOT NULL,
    vehicle_nickname text,
    move_document_id uuid NOT NULL,
    empty_weight integer,
    empty_weight_ticket_missing boolean NOT NULL,
    full_weight integer,
    full_weight_ticket_missing boolean NOT NULL,
    weight_ticket_date date,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    trailer_ownership_missing boolean NOT NULL,
    deleted_at timestamp with time zone,
    vehicle_make text,
    vehicle_model text
);


ALTER TABLE public.archived_weight_ticket_set_documents OWNER TO postgres;

--
-- Name: audit_history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.audit_history (
    id uuid NOT NULL,
    schema_name text NOT NULL,
    table_name text NOT NULL,
    relid oid NOT NULL,
    object_id uuid,
    session_userid uuid,
    event_name text,
    action_tstamp_tx timestamp with time zone NOT NULL,
    action_tstamp_stm timestamp with time zone NOT NULL,
    action_tstamp_clk timestamp with time zone NOT NULL,
    transaction_id bigint,
    client_query text,
    action text NOT NULL,
    old_data jsonb,
    changed_data jsonb,
    statement_only boolean NOT NULL,
    CONSTRAINT audit_history_action_check CHECK ((action = ANY (ARRAY['INSERT'::text, 'DELETE'::text, 'UPDATE'::text, 'TRUNCATE'::text])))
);


ALTER TABLE public.audit_history OWNER TO postgres;

--
-- Name: TABLE audit_history; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.audit_history IS 'History of auditable actions on audited tables, from if_modified_func()';


--
-- Name: COLUMN audit_history.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.id IS 'Unique identifier for each auditable event';


--
-- Name: COLUMN audit_history.schema_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.schema_name IS 'Name of audited table that this event is in';


--
-- Name: COLUMN audit_history.table_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.table_name IS 'Non-schema-qualified table name of table event occured in';


--
-- Name: COLUMN audit_history.relid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.relid IS 'Table OID. Changes with drop/create. Get with ''tablename''::regclass';


--
-- Name: COLUMN audit_history.object_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.object_id IS 'if the changed data has an id column';


--
-- Name: COLUMN audit_history.session_userid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.session_userid IS 'id of user whose statement caused the audited event';


--
-- Name: COLUMN audit_history.event_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.event_name IS 'name of event that caused the audited event';


--
-- Name: COLUMN audit_history.action_tstamp_tx; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.action_tstamp_tx IS 'Transaction start timestamp for tx in which audited event occurred';


--
-- Name: COLUMN audit_history.action_tstamp_stm; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.action_tstamp_stm IS 'Statement start timestamp for tx in which audited event occurred';


--
-- Name: COLUMN audit_history.action_tstamp_clk; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.action_tstamp_clk IS 'Wall clock time at which audited event''s trigger call occurred';


--
-- Name: COLUMN audit_history.transaction_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.transaction_id IS 'Identifier of transaction that made the change. May wrap, but unique paired with action_tstamp_tx.';


--
-- Name: COLUMN audit_history.client_query; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.client_query IS 'Text of the client query that triggered the audit event';


--
-- Name: COLUMN audit_history.action; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.action IS 'Action type';


--
-- Name: COLUMN audit_history.old_data; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.old_data IS 'Record value. Null for statement-level trigger. For INSERT this is NULL. For DELETE and UPDATE it is the old state of the record stored in json.';


--
-- Name: COLUMN audit_history.changed_data; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.changed_data IS 'New values of fields changed by INSERT AND UPDATE. Null except for row-level INSERT and UPDATE events.';


--
-- Name: COLUMN audit_history.statement_only; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.audit_history.statement_only IS 'TRUE if audit event is from an FOR EACH STATEMENT trigger, FALSE for FOR EACH ROW';


--
-- Name: audit_history_tableslist; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.audit_history_tableslist AS
 SELECT DISTINCT triggers.trigger_schema AS schema,
    triggers.event_object_table AS auditedtable
   FROM information_schema.triggers
  WHERE ((triggers.trigger_name)::text = ANY (ARRAY['audit_trigger_row'::text, 'audit_trigger_stm'::text]))
  ORDER BY triggers.trigger_schema, triggers.event_object_table;


ALTER TABLE public.audit_history_tableslist OWNER TO postgres;

--
-- Name: VIEW audit_history_tableslist; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON VIEW public.audit_history_tableslist IS '
View showing all tables with auditing set up. Ordered by schema, then table.
';


--
-- Name: backup_contacts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.backup_contacts (
    id uuid NOT NULL,
    service_member_id uuid NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    phone character varying(255),
    permission character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.backup_contacts OWNER TO postgres;

--
-- Name: TABLE backup_contacts; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.backup_contacts IS 'Holds all information regarding a backup contact for the customer';


--
-- Name: COLUMN backup_contacts.service_member_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.backup_contacts.service_member_id IS 'A foreign key that points to the service_members table';


--
-- Name: COLUMN backup_contacts.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.backup_contacts.name IS 'The name of the backup contact';


--
-- Name: COLUMN backup_contacts.email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.backup_contacts.email IS 'The email of the backup contact';


--
-- Name: COLUMN backup_contacts.phone; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.backup_contacts.phone IS 'The phone number of the backup contact';


--
-- Name: COLUMN backup_contacts.permission; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.backup_contacts.permission IS 'An enum with 3 possible values: None, View, Edit. Meanings: None: can contact only, View: can view all move details, Edit: can view and edit all move details';


--
-- Name: COLUMN backup_contacts.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.backup_contacts.created_at IS 'Date & time the backup contacts was created';


--
-- Name: COLUMN backup_contacts.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.backup_contacts.updated_at IS 'Date & time the backup contacts was last updated';


--
-- Name: client_certs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.client_certs (
    id uuid NOT NULL,
    sha256_digest character(64) NOT NULL,
    subject text NOT NULL,
    allow_orders_api boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    allow_air_force_orders_read boolean DEFAULT false NOT NULL,
    allow_air_force_orders_write boolean DEFAULT false NOT NULL,
    allow_army_orders_read boolean DEFAULT false NOT NULL,
    allow_army_orders_write boolean DEFAULT false NOT NULL,
    allow_coast_guard_orders_read boolean DEFAULT false NOT NULL,
    allow_coast_guard_orders_write boolean DEFAULT false NOT NULL,
    allow_marine_corps_orders_read boolean DEFAULT false NOT NULL,
    allow_marine_corps_orders_write boolean DEFAULT false NOT NULL,
    allow_navy_orders_read boolean DEFAULT false NOT NULL,
    allow_navy_orders_write boolean DEFAULT false NOT NULL,
    allow_prime boolean DEFAULT false NOT NULL,
    user_id uuid
);


ALTER TABLE public.client_certs OWNER TO postgres;

--
-- Name: TABLE client_certs; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.client_certs IS 'Holds the SSL/TLS certificates authorized for MilMove and indicates to which parts of the app they have access.';


--
-- Name: COLUMN client_certs.sha256_digest; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.sha256_digest IS 'The encrypted signature of the certificate';


--
-- Name: COLUMN client_certs.subject; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.subject IS 'The entity the certificate belongs to';


--
-- Name: COLUMN client_certs.allow_orders_api; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_orders_api IS 'Indicates whether or not the cert grants access to the Orders API';


--
-- Name: COLUMN client_certs.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.created_at IS 'Date & time the client cert was created';


--
-- Name: COLUMN client_certs.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.updated_at IS 'Date & time the client cert was last updated';


--
-- Name: COLUMN client_certs.allow_air_force_orders_read; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_air_force_orders_read IS 'Indicates whether or not the cert grants view-only access to Air Force orders';


--
-- Name: COLUMN client_certs.allow_air_force_orders_write; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_air_force_orders_write IS 'Indicates whether or not the cert grants edit access to Air Force orders';


--
-- Name: COLUMN client_certs.allow_army_orders_read; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_army_orders_read IS 'Indicates whether or not the cert grants view-only access to Army orders';


--
-- Name: COLUMN client_certs.allow_army_orders_write; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_army_orders_write IS 'Indicates whether or not the cert grants edit access to Army orders';


--
-- Name: COLUMN client_certs.allow_coast_guard_orders_read; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_coast_guard_orders_read IS 'Indicates whether or not the cert grants view-only access to Coast Guard orders';


--
-- Name: COLUMN client_certs.allow_coast_guard_orders_write; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_coast_guard_orders_write IS 'Indicates whether or not the cert grants edit access to Coast Guard orders';


--
-- Name: COLUMN client_certs.allow_marine_corps_orders_read; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_marine_corps_orders_read IS 'Indicates whether or not the cert grants view-only access to Marine Corps orders';


--
-- Name: COLUMN client_certs.allow_marine_corps_orders_write; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_marine_corps_orders_write IS 'Indicates whether or not the cert grants edit access to Marine Corps orders';


--
-- Name: COLUMN client_certs.allow_navy_orders_read; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_navy_orders_read IS 'Indicates whether or not the cert grants view-only access to Navy orders';


--
-- Name: COLUMN client_certs.allow_navy_orders_write; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_navy_orders_write IS 'Indicates whether or not the cert grants edit access to Navy orders';


--
-- Name: COLUMN client_certs.allow_prime; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.allow_prime IS 'Indicates whether or not the cert grants access to the Prime API';


--
-- Name: COLUMN client_certs.user_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.client_certs.user_id IS 'Associate a user with each client cert; initially designed to identify a prime "user" from http requests';


--
-- Name: contractors; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.contractors (
    id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    name character varying(80) NOT NULL,
    contract_number character varying(80) NOT NULL,
    type character varying(80) NOT NULL
);


ALTER TABLE public.contractors OWNER TO postgres;

--
-- Name: TABLE contractors; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.contractors IS 'Holds all contractors who handle moves. There is only one active contractor per type at a time, though we do not yet have a way to identify that.';


--
-- Name: COLUMN contractors.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.contractors.created_at IS 'Date & time the contractor was created';


--
-- Name: COLUMN contractors.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.contractors.updated_at IS 'Date & time the contractor was updated';


--
-- Name: COLUMN contractors.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.contractors.name IS 'The name of the contractor';


--
-- Name: COLUMN contractors.contract_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.contractors.contract_number IS 'The government-issued contract number for the contractor.';


--
-- Name: COLUMN contractors.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.contractors.type IS 'A string to represent the type of contractor. Examples are Prime and NTS.';


--
-- Name: customer_support_remarks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.customer_support_remarks (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    content text NOT NULL,
    office_user_id uuid NOT NULL,
    move_id uuid NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.customer_support_remarks OWNER TO postgres;

--
-- Name: TABLE customer_support_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.customer_support_remarks IS 'Store remarks from office users pertaining to moves.';


--
-- Name: COLUMN customer_support_remarks.content; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.customer_support_remarks.content IS 'Text content of the customer support remark written by an office user.';


--
-- Name: COLUMN customer_support_remarks.office_user_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.customer_support_remarks.office_user_id IS 'The office_user who authored the customer support remark.';


--
-- Name: COLUMN customer_support_remarks.move_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.customer_support_remarks.move_id IS 'The move the customer support remark is associated with.';


--
-- Name: COLUMN customer_support_remarks.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.customer_support_remarks.deleted_at IS 'Date & time that the customer support remark was deleted';


--
-- Name: distance_calculations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.distance_calculations (
    id uuid NOT NULL,
    origin_address_id uuid NOT NULL,
    destination_address_id uuid NOT NULL,
    distance_miles integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.distance_calculations OWNER TO postgres;

--
-- Name: TABLE distance_calculations; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.distance_calculations IS 'Represents a distance calculation in miles between an origin and destination address.';


--
-- Name: COLUMN distance_calculations.origin_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.distance_calculations.origin_address_id IS 'Represents the origin address as a foreign key to the addresses table.';


--
-- Name: COLUMN distance_calculations.destination_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.distance_calculations.destination_address_id IS 'Represents the destination address as a foreign key to the addresses table.';


--
-- Name: COLUMN distance_calculations.distance_miles; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.distance_calculations.distance_miles IS 'The distance in miles between the origin and destination address.';


--
-- Name: COLUMN distance_calculations.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.distance_calculations.created_at IS 'Date & time the distance_calculation was created';


--
-- Name: COLUMN distance_calculations.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.distance_calculations.updated_at IS 'Date & time the distance_calculation was updated';


--
-- Name: documents; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.documents (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    service_member_id uuid NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.documents OWNER TO postgres;

--
-- Name: TABLE documents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.documents IS 'Holds information about uploaded documents';


--
-- Name: COLUMN documents.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.documents.created_at IS 'Date & time the document was created';


--
-- Name: COLUMN documents.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.documents.updated_at IS 'Date & time the document was last updated';


--
-- Name: COLUMN documents.service_member_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.documents.service_member_id IS 'A foreign key that points to the service_members table';


--
-- Name: COLUMN documents.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.documents.deleted_at IS 'Date & time document was deleted';


--
-- Name: duty_location_names; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.duty_location_names (
    id uuid NOT NULL,
    name text NOT NULL,
    duty_location_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.duty_location_names OWNER TO postgres;

--
-- Name: TABLE duty_location_names; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.duty_location_names IS 'Holds information regarding alternate names for a duty station (used for duty station lookups)';


--
-- Name: COLUMN duty_location_names.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_location_names.name IS 'Any alternate name for a duty station other than the official name (common names, abbreviations, etc)';


--
-- Name: COLUMN duty_location_names.duty_location_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_location_names.duty_location_id IS 'A foreign key that points to the duty stations table';


--
-- Name: COLUMN duty_location_names.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_location_names.created_at IS 'Date & time the duty station name was created';


--
-- Name: COLUMN duty_location_names.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_location_names.updated_at IS 'Date & time the duty station name was last updated';


--
-- Name: duty_locations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.duty_locations (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    affiliation character varying(255),
    address_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    transportation_office_id uuid,
    provides_services_counseling boolean DEFAULT false NOT NULL
);


ALTER TABLE public.duty_locations OWNER TO postgres;

--
-- Name: TABLE duty_locations; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.duty_locations IS 'Holds information about the duty stations';


--
-- Name: COLUMN duty_locations.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_locations.name IS 'The name of the duty station';


--
-- Name: COLUMN duty_locations.affiliation; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_locations.affiliation IS 'The affiliation of the duty station (Army, Air Force, Navy, Marines, Coast Guard';


--
-- Name: COLUMN duty_locations.address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_locations.address_id IS 'A foreign key that points to the address table';


--
-- Name: COLUMN duty_locations.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_locations.created_at IS 'Date & time the duty station was created';


--
-- Name: COLUMN duty_locations.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_locations.updated_at IS 'Date & time the duty station was last updated';


--
-- Name: COLUMN duty_locations.transportation_office_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_locations.transportation_office_id IS 'A foreign key that points to the transportation_offices table';


--
-- Name: COLUMN duty_locations.provides_services_counseling; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.duty_locations.provides_services_counseling IS 'Indicates whether a duty station provides services counseling or not';


--
-- Name: edi_errors; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.edi_errors (
    id uuid NOT NULL,
    payment_request_id uuid NOT NULL,
    interchange_control_number_id uuid,
    code character varying,
    description character varying,
    edi_type character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.edi_errors OWNER TO postgres;

--
-- Name: TABLE edi_errors; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.edi_errors IS 'Stores errors when sending an EDI 858 or stores errors reported from EDI responses (997 & 824)';


--
-- Name: COLUMN edi_errors.payment_request_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_errors.payment_request_id IS 'Payment Request ID associated with this error';


--
-- Name: COLUMN edi_errors.interchange_control_number_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_errors.interchange_control_number_id IS 'ID for payment_request_to_interchange_control_numbers associated with this error. This will identify the ICN for the payment request.';


--
-- Name: COLUMN edi_errors.code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_errors.code IS 'Reported code from syncada for the EDI error encountered';


--
-- Name: COLUMN edi_errors.description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_errors.description IS 'Description of the error. Can be used with the edi_errors.code.';


--
-- Name: COLUMN edi_errors.edi_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_errors.edi_type IS 'Type of EDI reporting or causing the issue. Can be EDI 997, 824, and 858.';


--
-- Name: edi_processings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.edi_processings (
    id uuid NOT NULL,
    edi_type public.edi_type NOT NULL,
    num_edis_processed integer NOT NULL,
    process_started_at timestamp without time zone NOT NULL,
    process_ended_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.edi_processings OWNER TO postgres;

--
-- Name: TABLE edi_processings; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.edi_processings IS 'Stores metrics for the processing of EDIs.';


--
-- Name: COLUMN edi_processings.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_processings.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN edi_processings.edi_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_processings.edi_type IS 'The type of EDI being processed (e.g., "858", "998", etc.).';


--
-- Name: COLUMN edi_processings.num_edis_processed; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_processings.num_edis_processed IS 'The number of successfully processed EDIs of the given type.';


--
-- Name: COLUMN edi_processings.process_started_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_processings.process_started_at IS 'Timestamp when this processing started.';


--
-- Name: COLUMN edi_processings.process_ended_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_processings.process_ended_at IS 'Timestamp when this processing ended.';


--
-- Name: COLUMN edi_processings.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_processings.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN edi_processings.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.edi_processings.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: electronic_orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.electronic_orders (
    id uuid NOT NULL,
    orders_number character varying(255) NOT NULL,
    edipi character varying(255) NOT NULL,
    issuer character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.electronic_orders OWNER TO postgres;

--
-- Name: TABLE electronic_orders; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.electronic_orders IS 'Represents the electronic move orders issued by a particular branch of the military';


--
-- Name: COLUMN electronic_orders.orders_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders.orders_number IS 'A (generally) unique number identifying the orders, corresponding to the ORDERS number (Army), the CT SDN (Navy, Marines), the SPECIAL ORDER NO (Air Force), the Travel Order No. (Coast Guard), or the Travel Authorization Number (Civilian)';


--
-- Name: COLUMN electronic_orders.edipi; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders.edipi IS 'Electronic Data Interchange Personal Identifier, the 10 digit DoD ID Number of the service member';


--
-- Name: COLUMN electronic_orders.issuer; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders.issuer IS 'The organization that issued the orders (army, navy, etc.)';


--
-- Name: COLUMN electronic_orders.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders.created_at IS 'Date & time the electronic order was created';


--
-- Name: COLUMN electronic_orders.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders.updated_at IS 'Date & time the electronic order was last updated';


--
-- Name: electronic_orders_revisions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.electronic_orders_revisions (
    id uuid NOT NULL,
    electronic_order_id uuid NOT NULL,
    seq_num integer DEFAULT 0 NOT NULL,
    given_name character varying(255) NOT NULL,
    middle_name character varying(255),
    family_name character varying(255) NOT NULL,
    name_suffix character varying(255),
    affiliation character varying(255) NOT NULL,
    paygrade character varying(255) NOT NULL,
    title character varying(255),
    status character varying(255) NOT NULL,
    date_issued timestamp without time zone NOT NULL,
    no_cost_move boolean DEFAULT false NOT NULL,
    tdy_en_route boolean DEFAULT false NOT NULL,
    tour_type character varying(255) DEFAULT 'accompanied'::character varying NOT NULL,
    orders_type character varying(255) NOT NULL,
    has_dependents boolean NOT NULL,
    losing_uic character varying(255),
    losing_unit_name character varying(255),
    losing_unit_city character varying(255),
    losing_unit_locality character varying(255),
    losing_unit_country character varying(255),
    losing_unit_postal_code character varying(255),
    gaining_uic character varying(255),
    gaining_unit_name character varying(255),
    gaining_unit_city character varying(255),
    gaining_unit_locality character varying(255),
    gaining_unit_country character varying(255),
    gaining_unit_postal_code character varying(255),
    report_no_earlier_than timestamp without time zone,
    report_no_later_than timestamp without time zone,
    hhg_tac character varying(255),
    hhg_sdn character varying(255),
    hhg_loa character varying(255),
    nts_tac character varying(255),
    nts_sdn character varying(255),
    nts_loa character varying(255),
    pov_shipment_tac character varying(255),
    pov_shipment_sdn character varying(255),
    pov_shipment_loa character varying(255),
    pov_storage_tac character varying(255),
    pov_storage_sdn character varying(255),
    pov_storage_loa character varying(255),
    ub_tac character varying(255),
    ub_sdn character varying(255),
    ub_loa character varying(255),
    comments text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.electronic_orders_revisions OWNER TO postgres;

--
-- Name: TABLE electronic_orders_revisions; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.electronic_orders_revisions IS 'Represents revisions or edits to an issued set of electronic move orders';


--
-- Name: COLUMN electronic_orders_revisions.electronic_order_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.electronic_order_id IS 'The UUID of the electronic orders being revised';


--
-- Name: COLUMN electronic_orders_revisions.seq_num; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.seq_num IS 'An integer representing the sequence number for this revision. As orders are amended, the revision with the highest sequence number is considered the current, authoritative version of the orders, even if date_issued is earlier.';


--
-- Name: COLUMN electronic_orders_revisions.given_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.given_name IS 'The first name of the service member';


--
-- Name: COLUMN electronic_orders_revisions.middle_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.middle_name IS 'The middle name or initial of the service member';


--
-- Name: COLUMN electronic_orders_revisions.family_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.family_name IS 'The last name of the service member';


--
-- Name: COLUMN electronic_orders_revisions.name_suffix; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.name_suffix IS 'The suffix of the service member''s name, if any (Jr. Sr., III, etc.)';


--
-- Name: COLUMN electronic_orders_revisions.affiliation; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.affiliation IS 'The service member''s affiliated military branch (army, navy, etc.)';


--
-- Name: COLUMN electronic_orders_revisions.paygrade; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.paygrade IS 'The DoD paygrade or rank of the service member';


--
-- Name: COLUMN electronic_orders_revisions.title; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.title IS 'The preferred form of address for the service member. Used mainly for ranks that have multiple possible titles.';


--
-- Name: COLUMN electronic_orders_revisions.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.status IS 'Indicates whether these Orders are authorized, RFO (Request For Orders), or canceled';


--
-- Name: COLUMN electronic_orders_revisions.date_issued; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.date_issued IS 'The date and time thath these orders were cut. If omitted, the current date and time will be used.';


--
-- Name: COLUMN electronic_orders_revisions.no_cost_move; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.no_cost_move IS 'If true, indicates that these orders do not authorize any move expenses. If false, these orders are a PCS and should authorize expenses.';


--
-- Name: COLUMN electronic_orders_revisions.tdy_en_route; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.tdy_en_route IS 'TDY (Temporary Duty Yonder) en-route. If omitted, assume false.';


--
-- Name: COLUMN electronic_orders_revisions.tour_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.tour_type IS 'Accompanied or Unaccompanied - indicates whether or not dependents are authorized to accompany the service member on the move. If omitted, assume accompanied.';


--
-- Name: COLUMN electronic_orders_revisions.orders_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.orders_type IS 'The type of orders for this move (joining the military, retirement, training, etc.)';


--
-- Name: COLUMN electronic_orders_revisions.has_dependents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.has_dependents IS 'Indicates whether or not the service member has any dependents (spouse, children, etc.)';


--
-- Name: COLUMN electronic_orders_revisions.losing_uic; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.losing_uic IS 'The Unit Identification Code for the unit the service member is moving away from. A six character code that identifies each DoD entity.';


--
-- Name: COLUMN electronic_orders_revisions.losing_unit_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.losing_unit_name IS 'The human-readable name of the losing unit';


--
-- Name: COLUMN electronic_orders_revisions.losing_unit_city; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.losing_unit_city IS 'The city of the losing unit. May be FPO or APO for OCONUS commands.';


--
-- Name: COLUMN electronic_orders_revisions.losing_unit_locality; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.losing_unit_locality IS 'The locality of the losing unit. Will be the state for US units.';


--
-- Name: COLUMN electronic_orders_revisions.losing_unit_country; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.losing_unit_country IS 'The ISO 3166-1 alpha-2 country code for the losing unit. If blank, but city, locality, or postal_code are not blank, assume US';


--
-- Name: COLUMN electronic_orders_revisions.losing_unit_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.losing_unit_postal_code IS 'The postal code of the losing unit. Will be the ZIP code for US units.';


--
-- Name: COLUMN electronic_orders_revisions.gaining_uic; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.gaining_uic IS 'The Unit Identification Code for the unit the service member is moving to. May be blank if these are separation orders.';


--
-- Name: COLUMN electronic_orders_revisions.gaining_unit_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.gaining_unit_name IS 'The human-readable name of the gaining unit';


--
-- Name: COLUMN electronic_orders_revisions.gaining_unit_city; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.gaining_unit_city IS 'The city of the gaining unit. May be FPO or APO for OCONUS commands.';


--
-- Name: COLUMN electronic_orders_revisions.gaining_unit_locality; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.gaining_unit_locality IS 'The locality of the gaining unit. Will be the state for US units.';


--
-- Name: COLUMN electronic_orders_revisions.gaining_unit_country; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.gaining_unit_country IS 'The ISO 3166-1 alpha-2 country code for the gaining unit. If blank, but city, locality, or postal_code are not blank, assume US';


--
-- Name: COLUMN electronic_orders_revisions.gaining_unit_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.gaining_unit_postal_code IS 'The postal code of the gaining unit. Will be the ZIP code for US units.';


--
-- Name: COLUMN electronic_orders_revisions.report_no_earlier_than; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.report_no_earlier_than IS 'Earliest date that the service member is allowed to report for duty at the new duty station. If omitted, the member is allowed to report as early as desired.';


--
-- Name: COLUMN electronic_orders_revisions.report_no_later_than; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.report_no_later_than IS 'Latest date that the service member is allowed to report for duty at the new duty station. Should be included for most Orders types, but can be missing for Separation / Retirement Orders.';


--
-- Name: COLUMN electronic_orders_revisions.hhg_tac; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.hhg_tac IS 'Transportation Account Code. Used for accounting purposes in HHG expenses.';


--
-- Name: COLUMN electronic_orders_revisions.hhg_sdn; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.hhg_sdn IS 'Standard Document Number. Used for routing money for an HHG expense.';


--
-- Name: COLUMN electronic_orders_revisions.hhg_loa; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.hhg_loa IS 'The full Line of Accounting. Required if there is no TAC. Used for HHG expenses.';


--
-- Name: COLUMN electronic_orders_revisions.nts_tac; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.nts_tac IS 'Transportation Account Code. Used for accounting purposes in NTS expenses.';


--
-- Name: COLUMN electronic_orders_revisions.nts_sdn; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.nts_sdn IS 'Standard Document Number. Used for routing money for an NTS expense.';


--
-- Name: COLUMN electronic_orders_revisions.nts_loa; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.nts_loa IS 'The full Line of Accounting. Required if there is no TAC. Used for NTS expenses.';


--
-- Name: COLUMN electronic_orders_revisions.pov_shipment_tac; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.pov_shipment_tac IS 'Transportation Account Code. Used for accounting purposes in POV shipment expenses.';


--
-- Name: COLUMN electronic_orders_revisions.pov_shipment_sdn; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.pov_shipment_sdn IS 'Standard Document Number. Used for routing money for a POV shipment expense.';


--
-- Name: COLUMN electronic_orders_revisions.pov_shipment_loa; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.pov_shipment_loa IS 'The full Line of Accounting. Required if there is no TAC. Used for POV shipment expenses.';


--
-- Name: COLUMN electronic_orders_revisions.pov_storage_tac; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.pov_storage_tac IS 'Transportation Account Code. Used for accounting purposes in POV storage expenses.';


--
-- Name: COLUMN electronic_orders_revisions.pov_storage_sdn; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.pov_storage_sdn IS 'Standard Document Number. Used for routing money for a POV storage expense.';


--
-- Name: COLUMN electronic_orders_revisions.pov_storage_loa; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.pov_storage_loa IS 'The full Line of Accounting. Required if there is no TAC. Used for POV storage expenses.';


--
-- Name: COLUMN electronic_orders_revisions.ub_tac; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.ub_tac IS 'Transportation Account Code. Used for accounting purposes in UB expenses.';


--
-- Name: COLUMN electronic_orders_revisions.ub_sdn; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.ub_sdn IS 'Standard Document Number. Used for routing money for a UB expense.';


--
-- Name: COLUMN electronic_orders_revisions.ub_loa; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.ub_loa IS 'The full Line of Accounting. Required if there is no TAC. Used for UB expenses.';


--
-- Name: COLUMN electronic_orders_revisions.comments; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.comments IS 'Free-form text that may or may not contain information relevant to moving';


--
-- Name: COLUMN electronic_orders_revisions.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.created_at IS 'Date & time the revision for the electronic order was created';


--
-- Name: COLUMN electronic_orders_revisions.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.electronic_orders_revisions.updated_at IS 'Date & time the revision for the electronic order was last updated';


--
-- Name: entitlements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.entitlements (
    id uuid NOT NULL,
    dependents_authorized boolean,
    total_dependents integer,
    non_temporary_storage boolean,
    privately_owned_vehicle boolean,
    storage_in_transit integer,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    authorized_weight integer,
    required_medical_equipment_weight integer DEFAULT 0 NOT NULL,
    organizational_clothing_and_individual_equipment boolean DEFAULT false NOT NULL,
    pro_gear_weight integer DEFAULT 0 NOT NULL,
    pro_gear_weight_spouse integer DEFAULT 0 NOT NULL,
    CONSTRAINT entitlements_pro_gear_weight_check CHECK (((pro_gear_weight >= 0) AND (pro_gear_weight <= 2000))),
    CONSTRAINT entitlements_pro_gear_weight_spouse_check CHECK (((pro_gear_weight_spouse >= 0) AND (pro_gear_weight_spouse <= 500)))
);


ALTER TABLE public.entitlements OWNER TO postgres;

--
-- Name: TABLE entitlements; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.entitlements IS 'Service members are entitled to have the government pay to move a certain amount of weight, based on their rank, whether or not they have dependents, and whether their destination is CONUS or OCONUS. "Entitlements" is an older term, and the services now call these "allowances".';


--
-- Name: COLUMN entitlements.dependents_authorized; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.dependents_authorized IS 'A yes/no field reflecting whether dependents are authorized on the customer''s move orders';


--
-- Name: COLUMN entitlements.total_dependents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.total_dependents IS 'An integer reflecting the total number of dependents that are authorized on the customer''s move orders. For UB shipments, if dependents are authorized, each dependent adds to the customer''s weight allowance. (Note that the exact amount depends on the dependent''s age.)';


--
-- Name: COLUMN entitlements.non_temporary_storage; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.non_temporary_storage IS 'A yes/no field reflecting whether the customer is requesting a Non Temporary Storage shipment';


--
-- Name: COLUMN entitlements.privately_owned_vehicle; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.privately_owned_vehicle IS 'A yes/no field reflecting whether the customer has a privately owned vehicle that will need to be shipped as part of their move';


--
-- Name: COLUMN entitlements.storage_in_transit; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.storage_in_transit IS 'The maximum number of days of storage in transit allowed by the customer''s move orders';


--
-- Name: COLUMN entitlements.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.created_at IS 'Date & time the entitlement was created';


--
-- Name: COLUMN entitlements.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.updated_at IS 'Date & time the entitlement was last updated';


--
-- Name: COLUMN entitlements.authorized_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.authorized_weight IS 'The maximum number of pounds the Prime contractor is authorized to move for the customer';


--
-- Name: COLUMN entitlements.required_medical_equipment_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.required_medical_equipment_weight IS 'The RME (required medical equipment) weight in pounds. A Service member or a dependent who is entitled to, and
receiving, medical care authorized by 10 U.S.C. §1071-§1110. may ship medical equipment necessary for such care. The medical equipment may be shipped in the same way as HHG, but has no weight limit.
The weight of authorized medical equipment is not included in the maximum authorized HHG weight
allowance.
1. Required medical equipment does not include a modified personally owned vehicle.
2. For medical equipment to qualify for shipment under this paragraph, an appropriate.
Uniformed Services healthcare provider must certify that the equipment is necessary for medical
treatment of the Service member or the dependent who is authorized medical care under.';


--
-- Name: COLUMN entitlements.organizational_clothing_and_individual_equipment; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.organizational_clothing_and_individual_equipment IS 'A yes/no field reflecting whether the customer has OCIE (organizational clothing and individual equipment) that will need to be shipped as part of their move. Government property issued to the Service
member or employee by an Agency or Service for official use. A term specific to the Army and not other services.';


--
-- Name: COLUMN entitlements.pro_gear_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.pro_gear_weight IS 'This is equipment a member needs for the performance of official duties at the next or a later destination. Members are given a weight allowance for progear that is separate from their normal weight allowance.';


--
-- Name: COLUMN entitlements.pro_gear_weight_spouse; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.entitlements.pro_gear_weight_spouse IS 'This is equipment a member''s spouse needs for the performance of official duties at the next or a later destination. Members are given a weight allowance for progear that is separate from their normal weight allowance.';


--
-- Name: evaluation_reports; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.evaluation_reports (
    id uuid NOT NULL,
    office_user_id uuid NOT NULL,
    move_id uuid NOT NULL,
    shipment_id uuid,
    type public.evaluation_report_type NOT NULL,
    inspection_date date,
    inspection_type public.inspection_type,
    location public.evaluation_location_type,
    location_description text,
    violations_observed boolean,
    remarks text,
    submitted_at timestamp with time zone,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    completion_time interval GENERATED ALWAYS AS ((submitted_at - created_at)) STORED,
    serious_incident boolean,
    serious_incident_desc text,
    observed_claims_response_date date,
    observed_pickup_date date,
    observed_pickup_spread_start_date date,
    observed_pickup_spread_end_date date,
    observed_delivery_date date,
    observed_shipment_delivery_date date,
    observed_shipment_physical_pickup_date date,
    time_depart time without time zone,
    eval_start time without time zone,
    eval_end time without time zone
);


ALTER TABLE public.evaluation_reports OWNER TO postgres;

--
-- Name: TABLE evaluation_reports; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.evaluation_reports IS 'Contains QAE evaluation reports. There are two kinds of reports: shipment and counseling. You can tell them apart based on whether shipment_id is NULL.';


--
-- Name: COLUMN evaluation_reports.office_user_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.office_user_id IS 'The office_user who authored the evaluation report.';


--
-- Name: COLUMN evaluation_reports.move_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.move_id IS 'Move that the report is associated with';


--
-- Name: COLUMN evaluation_reports.shipment_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.shipment_id IS 'If present, indicates the shipment that this report is based on. NULL if this is not a shipment report.';


--
-- Name: COLUMN evaluation_reports.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.type IS 'Indicates type of report. Either counseling or shipment';


--
-- Name: COLUMN evaluation_reports.inspection_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.inspection_date IS 'date of inspection';


--
-- Name: COLUMN evaluation_reports.inspection_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.inspection_type IS 'Indicates the type of evaluation that is being described by this report. Either physical, virtual, or data review';


--
-- Name: COLUMN evaluation_reports.location; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.location IS 'Indicates whether the inspection was performed at the origin or destination of the shipment. If OTHER is selected, location_description should contain a description of the alternative location';


--
-- Name: COLUMN evaluation_reports.location_description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.location_description IS 'If the inspection was performed at a location other than the origin or destination of the shipment, this field contains a description of the location';


--
-- Name: COLUMN evaluation_reports.violations_observed; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.violations_observed IS 'True if any PWS violations were observed during the inspection';


--
-- Name: COLUMN evaluation_reports.remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.remarks IS 'Free text field for the evaluator''s notes about the inspection';


--
-- Name: COLUMN evaluation_reports.submitted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.submitted_at IS 'Time when the report was submitted. If NULL, then the report is still a draft';


--
-- Name: COLUMN evaluation_reports.completion_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.completion_time IS 'The time it took for an evaluation report to go from a created state to a submitted state.';


--
-- Name: COLUMN evaluation_reports.serious_incident; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.serious_incident IS 'Indicates is a serious incident was found';


--
-- Name: COLUMN evaluation_reports.serious_incident_desc; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.serious_incident_desc IS 'Text field for the description of the serious incident';


--
-- Name: COLUMN evaluation_reports.observed_claims_response_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.observed_claims_response_date IS 'Date of observed claims response ';


--
-- Name: COLUMN evaluation_reports.observed_pickup_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.observed_pickup_date IS 'Date of observed pickup';


--
-- Name: COLUMN evaluation_reports.observed_pickup_spread_start_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.observed_pickup_spread_start_date IS 'Start date of observed pickup spread';


--
-- Name: COLUMN evaluation_reports.observed_pickup_spread_end_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.observed_pickup_spread_end_date IS 'End date of observed pickup spread';


--
-- Name: COLUMN evaluation_reports.observed_delivery_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.observed_delivery_date IS 'Observed delivery date';


--
-- Name: COLUMN evaluation_reports.observed_shipment_delivery_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.observed_shipment_delivery_date IS 'Indicates shipment delivery date was different from scheduled';


--
-- Name: COLUMN evaluation_reports.observed_shipment_physical_pickup_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.observed_shipment_physical_pickup_date IS 'Indicates shipment pickup date was different from scheduled';


--
-- Name: COLUMN evaluation_reports.time_depart; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.time_depart IS 'Time departed for the evaluation, recorded in 24 hour format without timezone info';


--
-- Name: COLUMN evaluation_reports.eval_start; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.eval_start IS 'Time evaluation started, recorded in 24 hour format without timezone info';


--
-- Name: COLUMN evaluation_reports.eval_end; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.evaluation_reports.eval_end IS 'Time evaluation ended, recorded in 24 hour format without timezone info';


--
-- Name: fuel_eia_diesel_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.fuel_eia_diesel_prices (
    id uuid NOT NULL,
    pub_date date NOT NULL,
    rate_start_date date NOT NULL,
    rate_end_date date NOT NULL,
    eia_price_per_gallon_millicents integer NOT NULL,
    baseline_rate integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT fuel_eia_diesel_prices_baseline_rate_check CHECK ((baseline_rate > '-1'::integer)),
    CONSTRAINT fuel_eia_diesel_prices_baseline_rate_check1 CHECK ((baseline_rate < 101))
);


ALTER TABLE public.fuel_eia_diesel_prices OWNER TO postgres;

--
-- Name: TABLE fuel_eia_diesel_prices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.fuel_eia_diesel_prices IS 'Stores SDDC Fuel Surcharge rate information; used by pre-GHC HHG moves.';


--
-- Name: COLUMN fuel_eia_diesel_prices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.fuel_eia_diesel_prices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN fuel_eia_diesel_prices.pub_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.fuel_eia_diesel_prices.pub_date IS 'The date this rate was published.';


--
-- Name: COLUMN fuel_eia_diesel_prices.rate_start_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.fuel_eia_diesel_prices.rate_start_date IS 'The start date that this rate is applicable (inclusive).';


--
-- Name: COLUMN fuel_eia_diesel_prices.rate_end_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.fuel_eia_diesel_prices.rate_end_date IS 'The end date that this rate is applicable (inclusive).';


--
-- Name: COLUMN fuel_eia_diesel_prices.eia_price_per_gallon_millicents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.fuel_eia_diesel_prices.eia_price_per_gallon_millicents IS 'The national average price per gallon in millicents for this period as determined by the EIA (Energy Information Administration).';


--
-- Name: COLUMN fuel_eia_diesel_prices.baseline_rate; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.fuel_eia_diesel_prices.baseline_rate IS 'The calculated baseline fuel surcharge rate in cents for this period.';


--
-- Name: COLUMN fuel_eia_diesel_prices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.fuel_eia_diesel_prices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN fuel_eia_diesel_prices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.fuel_eia_diesel_prices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: ghc_diesel_fuel_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ghc_diesel_fuel_prices (
    id uuid NOT NULL,
    fuel_price_in_millicents integer NOT NULL,
    publication_date date NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.ghc_diesel_fuel_prices OWNER TO postgres;

--
-- Name: TABLE ghc_diesel_fuel_prices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.ghc_diesel_fuel_prices IS 'Represents the weekly average diesel fuel price; used in GHC pricing.';


--
-- Name: COLUMN ghc_diesel_fuel_prices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ghc_diesel_fuel_prices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN ghc_diesel_fuel_prices.fuel_price_in_millicents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ghc_diesel_fuel_prices.fuel_price_in_millicents IS 'The national average price per gallon in millicents for the week following the publication date as determined by the EIA (Energy Information Administration).';


--
-- Name: COLUMN ghc_diesel_fuel_prices.publication_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ghc_diesel_fuel_prices.publication_date IS 'The date this rate was published.';


--
-- Name: COLUMN ghc_diesel_fuel_prices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ghc_diesel_fuel_prices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN ghc_diesel_fuel_prices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ghc_diesel_fuel_prices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: ghc_domestic_transit_times; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ghc_domestic_transit_times (
    id uuid NOT NULL,
    max_days_transit_time integer NOT NULL,
    weight_lbs_lower integer NOT NULL,
    weight_lbs_upper integer NOT NULL,
    distance_miles_lower integer NOT NULL,
    distance_miles_upper integer NOT NULL
);


ALTER TABLE public.ghc_domestic_transit_times OWNER TO postgres;

--
-- Name: TABLE ghc_domestic_transit_times; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.ghc_domestic_transit_times IS 'Allows calculation of the maximum transit time based on the distance and weight ranges.';


--
-- Name: COLUMN ghc_domestic_transit_times.max_days_transit_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ghc_domestic_transit_times.max_days_transit_time IS 'The max transit time for the corresponding weight and distance ranges defined via the _lower and _upper columns.';


--
-- Name: COLUMN ghc_domestic_transit_times.weight_lbs_lower; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ghc_domestic_transit_times.weight_lbs_lower IS 'The minimum weight in the range.';


--
-- Name: COLUMN ghc_domestic_transit_times.weight_lbs_upper; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ghc_domestic_transit_times.weight_lbs_upper IS 'The maximum weight in the range. If 0 (zero), there is no upper bound';


--
-- Name: COLUMN ghc_domestic_transit_times.distance_miles_lower; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ghc_domestic_transit_times.distance_miles_lower IS 'The minimum distance in the range.';


--
-- Name: COLUMN ghc_domestic_transit_times.distance_miles_upper; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ghc_domestic_transit_times.distance_miles_upper IS 'The maximum distance in the range.';


--
-- Name: interchange_control_number; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.interchange_control_number
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    MAXVALUE 999999999
    CACHE 1
    CYCLE;


ALTER TABLE public.interchange_control_number OWNER TO postgres;

--
-- Name: invoice_number_trackers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.invoice_number_trackers (
    standard_carrier_alpha_code text NOT NULL,
    year integer NOT NULL,
    sequence_number integer NOT NULL
);


ALTER TABLE public.invoice_number_trackers OWNER TO postgres;

--
-- Name: TABLE invoice_number_trackers; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.invoice_number_trackers IS 'Tracks latest sequence numbers in SCAC/year groupings; this sequence number is part of the generated invoice number.';


--
-- Name: COLUMN invoice_number_trackers.standard_carrier_alpha_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoice_number_trackers.standard_carrier_alpha_code IS 'The associated SCAC for this sequence number (see the transportation_service_providers table).';


--
-- Name: COLUMN invoice_number_trackers.year; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoice_number_trackers.year IS 'The associated year for this sequence number.';


--
-- Name: COLUMN invoice_number_trackers.sequence_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoice_number_trackers.sequence_number IS 'The last used sequence number for the given SCAC/year.';


--
-- Name: invoices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.invoices (
    id uuid NOT NULL,
    status character varying(255) NOT NULL,
    invoiced_date timestamp with time zone NOT NULL,
    invoice_number character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    approver_id uuid NOT NULL,
    user_uploads_id uuid
);


ALTER TABLE public.invoices OWNER TO postgres;

--
-- Name: TABLE invoices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.invoices IS 'Represents an invoice sent to GEX; only used by pre-GHC HHG moves at the moment.';


--
-- Name: COLUMN invoices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN invoices.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoices.status IS 'Status of this invoice; options are DRAFT, IN_PROCESS, SUBMITTED, SUBMISSION_FAILURE, UPDATE_FAILURE.';


--
-- Name: COLUMN invoices.invoiced_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoices.invoiced_date IS 'Timestamp when this invoice was sent to GEX.';


--
-- Name: COLUMN invoices.invoice_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoices.invoice_number IS 'A unique invoice number. Format is SCAC + two digit year + sequence number (with a suffix of -01, -02, etc. appended for subsequent invoices on the same shipment).';


--
-- Name: COLUMN invoices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN invoices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: COLUMN invoices.approver_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoices.approver_id IS 'The office user that approved this invoice.';


--
-- Name: COLUMN invoices.user_uploads_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.invoices.user_uploads_id IS 'The associated uploads used as justification for this invoice.';


--
-- Name: jppso_region_state_assignments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.jppso_region_state_assignments (
    id uuid NOT NULL,
    jppso_region_id uuid,
    state_name text NOT NULL,
    state_abbreviation text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.jppso_region_state_assignments OWNER TO postgres;

--
-- Name: TABLE jppso_region_state_assignments; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.jppso_region_state_assignments IS 'Maps US states to JPPSO regions. This table is not currently used, but will be soon in order to associate a TOO with a specific JPPSO, which will allow the TOO to filter the list of moves by region.';


--
-- Name: COLUMN jppso_region_state_assignments.jppso_region_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.jppso_region_state_assignments.jppso_region_id IS 'The JPPSO region this state is part of. A foreign key to the jppso_regions table.';


--
-- Name: COLUMN jppso_region_state_assignments.state_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.jppso_region_state_assignments.state_name IS 'The full capitalized US state name.';


--
-- Name: COLUMN jppso_region_state_assignments.state_abbreviation; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.jppso_region_state_assignments.state_abbreviation IS 'The two-letter state abbreviation.';


--
-- Name: COLUMN jppso_region_state_assignments.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.jppso_region_state_assignments.created_at IS 'Date & time the jppso_region_state_assignment was created.';


--
-- Name: COLUMN jppso_region_state_assignments.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.jppso_region_state_assignments.updated_at IS 'Date & time the jppso_region_state_assignment was updated.';


--
-- Name: jppso_regions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.jppso_regions (
    id uuid NOT NULL,
    code text NOT NULL,
    name text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.jppso_regions OWNER TO postgres;

--
-- Name: TABLE jppso_regions; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.jppso_regions IS 'Holds all JPPSO region names and codes. This is used to map states to regions. This table is not currently used, but will be soon in order to associate a TOO with a specific JPPSO, which will allow the TOO to filter the list of moves by region.';


--
-- Name: COLUMN jppso_regions.code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.jppso_regions.code IS 'The 4-character code for the region.';


--
-- Name: COLUMN jppso_regions.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.jppso_regions.name IS 'The human-readable name of the region.';


--
-- Name: COLUMN jppso_regions.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.jppso_regions.created_at IS 'Date & time the jppso_region was created.';


--
-- Name: COLUMN jppso_regions.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.jppso_regions.updated_at IS 'Date & time the jppso_region was updated.';


--
-- Name: mto_shipments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mto_shipments (
    id uuid NOT NULL,
    move_id uuid,
    scheduled_pickup_date date,
    requested_pickup_date date,
    customer_remarks text,
    pickup_address_id uuid,
    destination_address_id uuid,
    secondary_pickup_address_id uuid,
    secondary_delivery_address_id uuid,
    prime_estimated_weight integer,
    prime_estimated_weight_recorded_date date,
    prime_actual_weight integer,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    shipment_type public.mto_shipment_type DEFAULT 'HHG'::public.mto_shipment_type NOT NULL,
    status public.mto_shipment_status DEFAULT 'DRAFT'::public.mto_shipment_status,
    rejection_reason text,
    actual_pickup_date date,
    approved_date date,
    first_available_delivery_date date,
    required_delivery_date date,
    days_in_storage integer,
    requested_delivery_date date,
    distance integer,
    diversion boolean DEFAULT false NOT NULL,
    counselor_remarks text,
    deleted_at timestamp with time zone,
    billable_weight_cap integer,
    billable_weight_justification text,
    sit_days_allowance integer,
    uses_external_vendor boolean DEFAULT false NOT NULL,
    storage_facility_id uuid,
    service_order_number character varying(255),
    tac_type public.loa_type,
    sac_type public.loa_type,
    nts_recorded_weight integer,
    destination_address_type public.destination_address_type,
    scheduled_delivery_date date,
    actual_delivery_date date,
    has_secondary_pickup_address boolean,
    has_secondary_delivery_address boolean
);


ALTER TABLE public.mto_shipments OWNER TO postgres;

--
-- Name: TABLE mto_shipments; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.mto_shipments IS 'A move task order (MTO) shipment for a specific MTO.';


--
-- Name: COLUMN mto_shipments.move_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.move_id IS 'The UUID of the move this shipment is for';


--
-- Name: COLUMN mto_shipments.scheduled_pickup_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.scheduled_pickup_date IS 'The pickup date the Prime contractor schedules for a shipment after consultation with the customer';


--
-- Name: COLUMN mto_shipments.requested_pickup_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.requested_pickup_date IS 'The date the customer is requesting that a given shipment be picked up.';


--
-- Name: COLUMN mto_shipments.customer_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.customer_remarks IS 'The remarks field is where the customer can describe special circumstances for their shipment, in order to inform the Prime contractor of any unique shipping and handling needs.';


--
-- Name: COLUMN mto_shipments.pickup_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.pickup_address_id IS 'The customer''s pickup address for a shipment';


--
-- Name: COLUMN mto_shipments.destination_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.destination_address_id IS 'The customer''s destination address for a shipment';


--
-- Name: COLUMN mto_shipments.secondary_pickup_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.secondary_pickup_address_id IS 'The secondary pickup address for this shipment';


--
-- Name: COLUMN mto_shipments.secondary_delivery_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.secondary_delivery_address_id IS 'The secondary delivery address for this shipment';


--
-- Name: COLUMN mto_shipments.prime_estimated_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.prime_estimated_weight IS 'The estimated weight of a shipment, provided by the Prime contractor after they survey a customer''s shipment';


--
-- Name: COLUMN mto_shipments.prime_estimated_weight_recorded_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.prime_estimated_weight_recorded_date IS 'Date when the Prime contractor records the shipment''s estimated weight';


--
-- Name: COLUMN mto_shipments.prime_actual_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.prime_actual_weight IS 'The actual weight of a shipment, provided by the Prime contractor after they pack, pickup, and weigh a customer''s shipment';


--
-- Name: COLUMN mto_shipments.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.created_at IS 'Date & time the shipment was created';


--
-- Name: COLUMN mto_shipments.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.updated_at IS 'Date & time the shipment was last updated';


--
-- Name: COLUMN mto_shipments.shipment_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.shipment_type IS 'The type of shipment. The list includes:
1. Personally procured move (PPM)
2. Household goods move (HHG)
3. Non-temporary storage (NTS)';


--
-- Name: COLUMN mto_shipments.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.status IS 'The status of a shipment.';


--
-- Name: COLUMN mto_shipments.rejection_reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.rejection_reason IS 'Not currently used, until the "reject" or "cancel" a shipment feature is implemented. When the Transportation Ordering Officer rejects or cancels a shipment, they will explain why';


--
-- Name: COLUMN mto_shipments.actual_pickup_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.actual_pickup_date IS 'The actual pickup date when the Prime contractor picks up the customer''s shipment';


--
-- Name: COLUMN mto_shipments.approved_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.approved_date IS 'The date when the Transportation Ordering Officer approves the shipment, and it is added to the Move Task Order for the Prime contractor';


--
-- Name: COLUMN mto_shipments.first_available_delivery_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.first_available_delivery_date IS 'Date the Prime provides to the customer so the customer can plan their own travel accordingly. We need to collect the FADD on the MTO so there is a record of what the Prime said they told the customer in case a situation arises in which the customer is unavailable to receive delivery of a shipment and the Prime wants to put the shipment in SIT.';


--
-- Name: COLUMN mto_shipments.required_delivery_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.required_delivery_date IS 'Latest date the Prime can deliver a customer''s shipment without violating the contract. RDD is the last date in the spread of available dates calculated from the scheduled pickup date.';


--
-- Name: COLUMN mto_shipments.days_in_storage; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.days_in_storage IS 'Related specifically to SIT. Total number of days a shipment was in temporary storage, determined after it comes out of SIT.';


--
-- Name: COLUMN mto_shipments.requested_delivery_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.requested_delivery_date IS 'Entered by the customer. Available delivery dates and required delivery date are calculated based on this date. Not at all a guarantee that this is the date the Prime will deliver the shipment.';


--
-- Name: COLUMN mto_shipments.distance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.distance IS 'Distance the shipment traveled, in miles';


--
-- Name: COLUMN mto_shipments.diversion; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.diversion IS 'Indicate if the shipment is part of a diversion';


--
-- Name: COLUMN mto_shipments.counselor_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.counselor_remarks IS 'Remarks service counselor has on the MTO Shipment';


--
-- Name: COLUMN mto_shipments.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.deleted_at IS 'Indicates whether the shipment has been soft deleted or not, and when it was soft deleted.';


--
-- Name: COLUMN mto_shipments.billable_weight_cap; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.billable_weight_cap IS 'The billable weight cap that the TIO can set per shipment that affects pricing';


--
-- Name: COLUMN mto_shipments.billable_weight_justification; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.billable_weight_justification IS 'The reasoning for why the TIO has set the billable_weight_cap to the chosen value';


--
-- Name: COLUMN mto_shipments.sit_days_allowance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.sit_days_allowance IS 'Total number of SIT days allowed for this shipment, including any sit extensions that have been approved';


--
-- Name: COLUMN mto_shipments.uses_external_vendor; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.uses_external_vendor IS 'Whether this shipment is handled by an external vendor, or by the prime';


--
-- Name: COLUMN mto_shipments.storage_facility_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.storage_facility_id IS 'The storage facility for an NTS shipment where items are stored';


--
-- Name: COLUMN mto_shipments.service_order_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.service_order_number IS 'The order number for a shipment in TOPS';


--
-- Name: COLUMN mto_shipments.tac_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.tac_type IS 'Indicates which type of TAC code to use for the shipment';


--
-- Name: COLUMN mto_shipments.sac_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.sac_type IS 'Indicates which type of SAC code to use for the shipment';


--
-- Name: COLUMN mto_shipments.nts_recorded_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.nts_recorded_weight IS 'Previously recorded weight used for NTS shipment';


--
-- Name: COLUMN mto_shipments.destination_address_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.destination_address_type IS 'Type of destination address location for retirees and separatees';


--
-- Name: COLUMN mto_shipments.scheduled_delivery_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.scheduled_delivery_date IS 'The delivery date the Prime contractor schedules for a shipment after consultation with the customer';


--
-- Name: COLUMN mto_shipments.actual_delivery_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.actual_delivery_date IS 'The actual date that the shipment was delivered to the destination address by the Prime';


--
-- Name: COLUMN mto_shipments.has_secondary_pickup_address; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.has_secondary_pickup_address IS 'False if the shipment does not have a secondary pickup address. This column exists to make it possible to tell whether a shipment update should delete an address or not modify it.';


--
-- Name: COLUMN mto_shipments.has_secondary_delivery_address; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_shipments.has_secondary_delivery_address IS 'False if the shipment does not have a secondary delivery address. This column exists to make it possible to tell whether a shipment update should delete an address or not modify it.';


--
-- Name: postal_code_to_gblocs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.postal_code_to_gblocs (
    postal_code character varying(5) NOT NULL,
    gbloc character varying(4) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    id uuid NOT NULL
);


ALTER TABLE public.postal_code_to_gblocs OWNER TO postgres;

--
-- Name: TABLE postal_code_to_gblocs; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.postal_code_to_gblocs IS 'This table is used to look up which GBLOC to use for a given postal code. Shipments from postal codes that are not in this table will not be supported, so it will need to be updated occasionally as new codes are added';


--
-- Name: COLUMN postal_code_to_gblocs.postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.postal_code_to_gblocs.postal_code IS 'A United States Postal Code, also known as ZIP code';


--
-- Name: COLUMN postal_code_to_gblocs.gbloc; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.postal_code_to_gblocs.gbloc IS 'GBLOC (Government Bill of Lading Office Code) used for a particular postal code';


--
-- Name: ppm_shipments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ppm_shipments (
    id uuid NOT NULL,
    shipment_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    status public.ppm_shipment_status NOT NULL,
    expected_departure_date date NOT NULL,
    actual_move_date date,
    submitted_at timestamp with time zone,
    reviewed_at timestamp with time zone,
    approved_at timestamp with time zone,
    pickup_postal_code character varying NOT NULL,
    secondary_pickup_postal_code character varying,
    destination_postal_code character varying NOT NULL,
    secondary_destination_postal_code character varying,
    sit_expected boolean DEFAULT false NOT NULL,
    estimated_weight integer,
    has_pro_gear boolean,
    pro_gear_weight integer,
    spouse_pro_gear_weight integer,
    estimated_incentive integer,
    deleted_at timestamp with time zone,
    sit_location public.sit_location_type,
    sit_estimated_weight integer,
    sit_estimated_entry_date date,
    sit_estimated_departure_date date,
    sit_estimated_cost integer,
    actual_pickup_postal_code character varying,
    actual_destination_postal_code character varying,
    has_requested_advance boolean,
    advance_amount_requested integer,
    has_received_advance boolean,
    advance_amount_received integer,
    advance_status public.ppm_advance_status,
    w2_address_id uuid,
    final_incentive integer,
    aoa_packet_id uuid,
    payment_packet_id uuid
);


ALTER TABLE public.ppm_shipments OWNER TO postgres;

--
-- Name: TABLE ppm_shipments; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.ppm_shipments IS 'Stores all PPM shipments, and their details.';


--
-- Name: COLUMN ppm_shipments.shipment_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.shipment_id IS 'MTO shipment ID associated with this PPM shipment.';


--
-- Name: COLUMN ppm_shipments.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.status IS 'Status of the PPM shipment.';


--
-- Name: COLUMN ppm_shipments.expected_departure_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.expected_departure_date IS 'Expected date this PPM shipment begins.';


--
-- Name: COLUMN ppm_shipments.actual_move_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.actual_move_date IS 'Actual date of the move associated with this PPM shipment.';


--
-- Name: COLUMN ppm_shipments.submitted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.submitted_at IS 'Date that PPM shipment information was submitted.';


--
-- Name: COLUMN ppm_shipments.reviewed_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.reviewed_at IS 'Date that PPM shipment information was reviewed.';


--
-- Name: COLUMN ppm_shipments.approved_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.approved_at IS 'Date that PPM shipment information was approved.';


--
-- Name: COLUMN ppm_shipments.pickup_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.pickup_postal_code IS 'Postal code where PPM begins.';


--
-- Name: COLUMN ppm_shipments.secondary_pickup_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.secondary_pickup_postal_code IS 'Secondary postal code where PPM shipment is to be picked up.';


--
-- Name: COLUMN ppm_shipments.destination_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.destination_postal_code IS 'Destination postal code for PPM shipment.';


--
-- Name: COLUMN ppm_shipments.secondary_destination_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.secondary_destination_postal_code IS 'Secondary destination postal code for PPM shipment.';


--
-- Name: COLUMN ppm_shipments.sit_expected; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.sit_expected IS 'Indicate if SIT is expected for PPM shipment.';


--
-- Name: COLUMN ppm_shipments.estimated_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.estimated_weight IS 'Estimated weight of PPM shipment.';


--
-- Name: COLUMN ppm_shipments.has_pro_gear; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.has_pro_gear IS 'Indicate if PPM shipment has pro gear.';


--
-- Name: COLUMN ppm_shipments.pro_gear_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.pro_gear_weight IS 'Indicate weight of PPM shipment pro gear.';


--
-- Name: COLUMN ppm_shipments.spouse_pro_gear_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.spouse_pro_gear_weight IS 'Indicate weight of PPM shipment spouse pro gear.';


--
-- Name: COLUMN ppm_shipments.estimated_incentive; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.estimated_incentive IS 'Estimated incentive associated with PPM shipment.';


--
-- Name: COLUMN ppm_shipments.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.deleted_at IS 'Indicates whether the ppm shipment has been soft deleted or not, and when it was soft deleted.';


--
-- Name: COLUMN ppm_shipments.sit_location; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.sit_location IS 'Records whether the PPM''s SIT is at the origin or destination.';


--
-- Name: COLUMN ppm_shipments.sit_estimated_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.sit_estimated_weight IS 'The estimated weight of the PPM''s SIT.';


--
-- Name: COLUMN ppm_shipments.sit_estimated_entry_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.sit_estimated_entry_date IS 'The estimated date the PPM''s items will go into SIT.';


--
-- Name: COLUMN ppm_shipments.sit_estimated_departure_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.sit_estimated_departure_date IS 'The estimated date the PPM''s items will come out of SIT.';


--
-- Name: COLUMN ppm_shipments.sit_estimated_cost; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.sit_estimated_cost IS 'The estimated cost (in cents) of the PPM''s SIT.';


--
-- Name: COLUMN ppm_shipments.actual_pickup_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.actual_pickup_postal_code IS 'Tracks the actual postal code where the PPM shipment began.';


--
-- Name: COLUMN ppm_shipments.actual_destination_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.actual_destination_postal_code IS 'Tracks the actual destination postal code for PPM shipment once customer moved the shipment.';


--
-- Name: COLUMN ppm_shipments.has_requested_advance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.has_requested_advance IS 'Indicates if a customer requested an advance for their PPM shipment.';


--
-- Name: COLUMN ppm_shipments.advance_amount_requested; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.advance_amount_requested IS 'Tracks the amount a customer requested for their advance; amount should be a percentage of estimated incentive.';


--
-- Name: COLUMN ppm_shipments.has_received_advance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.has_received_advance IS 'Indicates if a customer actually received an advance for their PPM shipment.';


--
-- Name: COLUMN ppm_shipments.advance_amount_received; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.advance_amount_received IS 'Tracks the amount a customer received for their advance; amount should be a percentage of estimated incentive.';


--
-- Name: COLUMN ppm_shipments.advance_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.advance_status IS 'An indicator that an office user has denied, approved or edited the requested advance';


--
-- Name: COLUMN ppm_shipments.w2_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.w2_address_id IS 'Address for where a customer received their W2 tax form';


--
-- Name: COLUMN ppm_shipments.final_incentive; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.final_incentive IS 'The final calculated incentive for the PPM shipment. This does not include SIT as it is a reimbursement.';


--
-- Name: COLUMN ppm_shipments.aoa_packet_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.aoa_packet_id IS 'The ID of the document that is associated with the upload containing the generated AOA packet for this PPM Shipment.';


--
-- Name: COLUMN ppm_shipments.payment_packet_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.ppm_shipments.payment_packet_id IS 'The ID of the document that is associated with the upload containing the generated payment packet for this PPM Shipment.';


--
-- Name: move_to_gbloc; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.move_to_gbloc AS
 SELECT DISTINCT ON (sh.move_id) sh.move_id,
    COALESCE(pctg.gbloc, pctg_ppm.gbloc) AS gbloc
   FROM ((public.mto_shipments sh
     LEFT JOIN ( SELECT a.id AS address_id,
            pctg_1.gbloc
           FROM (public.addresses a
             JOIN public.postal_code_to_gblocs pctg_1 ON (((a.postal_code)::text = (pctg_1.postal_code)::text)))) pctg ON ((pctg.address_id = sh.pickup_address_id)))
     LEFT JOIN ( SELECT ppm.shipment_id,
            pctg_1.gbloc
           FROM (public.ppm_shipments ppm
             JOIN public.postal_code_to_gblocs pctg_1 ON (((ppm.pickup_postal_code)::text = (pctg_1.postal_code)::text)))) pctg_ppm ON ((pctg_ppm.shipment_id = sh.id)))
  WHERE (sh.deleted_at IS NULL)
  ORDER BY sh.move_id, sh.created_at;


ALTER TABLE public.move_to_gbloc OWNER TO postgres;

--
-- Name: moves; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.moves (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    orders_id uuid NOT NULL,
    status character varying(255) DEFAULT 'DRAFT'::character varying NOT NULL,
    locator character(6) NOT NULL,
    cancel_reason character varying(255),
    show boolean DEFAULT true NOT NULL,
    contractor_id uuid,
    available_to_prime_at timestamp with time zone,
    ppm_type character varying(10),
    ppm_estimated_weight integer,
    reference_id character varying(255),
    submitted_at timestamp without time zone,
    service_counseling_completed_at timestamp with time zone,
    excess_weight_qualified_at timestamp with time zone,
    excess_weight_upload_id uuid,
    excess_weight_acknowledged_at timestamp with time zone,
    tio_remarks text,
    billable_weights_reviewed_at timestamp with time zone,
    financial_review_flag boolean DEFAULT false NOT NULL,
    financial_review_remarks text,
    financial_review_flag_set_at timestamp with time zone,
    prime_counseling_completed_at timestamp with time zone,
    closeout_office_id uuid,
    approvals_requested_at timestamp with time zone
);


ALTER TABLE public.moves OWNER TO postgres;

--
-- Name: TABLE moves; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.moves IS 'Contains all the information on the Move and Move Task Order (MTO). There is one MTO per a customer''s move.';


--
-- Name: COLUMN moves.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.created_at IS 'Date & time the Move was created';


--
-- Name: COLUMN moves.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.updated_at IS 'Date & time the Move was last updated';


--
-- Name: COLUMN moves.orders_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.orders_id IS 'Unique identifier for the orders issued for this move.';


--
-- Name: COLUMN moves.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.status IS 'The current status of the Move. Allowed values are:
DRAFT,
SUBMITTED,
APPROVED,
APPROVALS REQUESTED,
CANCELED,
NEEDS SERVICE COUNSELING,
SERVICE COUNSELING COMPLETED.
For more details about the lifecycle of a move and its statuses, check out this
Miro board: https://miro.com/app/board/o9J_krR2Tt8=/';


--
-- Name: COLUMN moves.locator; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.locator IS 'A 6-digit alphanumeric value that is a sharable, human-readable identifier for a move (so it could be disclosed to support staff, for instance).';


--
-- Name: COLUMN moves.cancel_reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.cancel_reason IS 'A string to explain why a move was canceled.';


--
-- Name: COLUMN moves.show; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.show IS 'A boolean that allows admin users to prevent a move from showing up in the TxO queue. This came out of a HackerOne engagement where hundreds of fake moves were created.';


--
-- Name: COLUMN moves.contractor_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.contractor_id IS 'Unique identifier for the prime contractor.';


--
-- Name: COLUMN moves.available_to_prime_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.available_to_prime_at IS 'Date & time the TOO made the MTO available to the prime contractor.';


--
-- Name: COLUMN moves.ppm_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.ppm_type IS 'Identifies whether a move is a full PPM or a partial PPM — the customer moving everything or only some things. This field is set by the Prime.';


--
-- Name: COLUMN moves.ppm_estimated_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.ppm_estimated_weight IS 'Estimated weight of the part of a customer''s belongings that they will move in a PPM. Unit is pounds. Customer does the estimation for PPMs. This field is set by the Prime.';


--
-- Name: COLUMN moves.reference_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.reference_id IS 'A unique identifier for an MTO (which also serves as the prefix for payment request numbers) in `dddd-dddd` format. There is still an ongoing discussion as to whether or not we need this `reference_id` in addition to the unique `locator` identifier.';


--
-- Name: COLUMN moves.service_counseling_completed_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.service_counseling_completed_at IS 'The timestamp when service counseling was completed.';


--
-- Name: COLUMN moves.excess_weight_qualified_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.excess_weight_qualified_at IS 'The date and time the sum of all the move''s shipments met the excess weight qualification threshold';


--
-- Name: COLUMN moves.excess_weight_upload_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.excess_weight_upload_id IS 'An uploaded document by the movers proving that the customer has been counseled about excess weight';


--
-- Name: COLUMN moves.excess_weight_acknowledged_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.excess_weight_acknowledged_at IS 'The date and time the TOO dismissed the risk of excess weight alert or updated the max billable weight.';


--
-- Name: COLUMN moves.tio_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.tio_remarks IS 'Remarks a TIO has on a move';


--
-- Name: COLUMN moves.billable_weights_reviewed_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.billable_weights_reviewed_at IS 'The date and time the TIO reviewed the billable weight for a move and its shipments.';


--
-- Name: COLUMN moves.financial_review_flag; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.financial_review_flag IS 'This flag is set by office users when they believe a move may incur excess costs to the customer and should have Finance Office review. The government will query this field from the data warehouse, so changes to it may require coordination.';


--
-- Name: COLUMN moves.financial_review_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.financial_review_remarks IS 'Reason provided by an office user for requesting financial review. The government will query this field from the data warehouse, so changes to it may require coordination.';


--
-- Name: COLUMN moves.financial_review_flag_set_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.financial_review_flag_set_at IS 'Time that financial review was requested at';


--
-- Name: COLUMN moves.prime_counseling_completed_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.prime_counseling_completed_at IS 'The timestamp when prime counseling was completed.';


--
-- Name: COLUMN moves.closeout_office_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.closeout_office_id IS 'The ID of the associated transportation office that is the closeout office for a move.';


--
-- Name: COLUMN moves.approvals_requested_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moves.approvals_requested_at IS 'The timestamp when a service item was added that made the move to a status of APPROVALS REQUESTED.';


--
-- Name: moving_expenses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.moving_expenses (
    id uuid NOT NULL,
    ppm_shipment_id uuid NOT NULL,
    document_id uuid NOT NULL,
    moving_expense_type public.moving_expense_type,
    description character varying,
    paid_with_gtcc boolean,
    amount integer,
    missing_receipt boolean,
    status public.ppm_document_status,
    reason character varying,
    sit_start_date date,
    sit_end_date date,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.moving_expenses OWNER TO postgres;

--
-- Name: TABLE moving_expenses; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.moving_expenses IS 'Stores expense doc and information associated with a trip for a PPM shipment.';


--
-- Name: COLUMN moving_expenses.ppm_shipment_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.ppm_shipment_id IS 'The ID of the PPM shipment that this expense is for.';


--
-- Name: COLUMN moving_expenses.document_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.document_id IS 'The ID of the document that is associated with the user uploads containing the moving expense receipt.';


--
-- Name: COLUMN moving_expenses.moving_expense_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.moving_expense_type IS 'Identifies the type of expense this is.';


--
-- Name: COLUMN moving_expenses.description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.description IS 'Stores a description of the expense.';


--
-- Name: COLUMN moving_expenses.paid_with_gtcc; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.paid_with_gtcc IS 'Indicates if the customer paid using a Government Travel Charge Card (GTCC).';


--
-- Name: COLUMN moving_expenses.amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.amount IS 'Stores the cost of the expense.';


--
-- Name: COLUMN moving_expenses.missing_receipt; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.missing_receipt IS 'Indicates if the customer is missing the receipt for their expense.';


--
-- Name: COLUMN moving_expenses.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.status IS 'Status of the expense, e.g. APPROVED.';


--
-- Name: COLUMN moving_expenses.reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.reason IS 'Contains the reason an expense is excluded or rejected; otherwise null.';


--
-- Name: COLUMN moving_expenses.sit_start_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.sit_start_date IS 'If this is a STORAGE expense, this indicates the date storage began.';


--
-- Name: COLUMN moving_expenses.sit_end_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.moving_expenses.sit_end_date IS 'If this is a STORAGE expense, this indicates the date storage ended.';


--
-- Name: mto_agents; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mto_agents (
    id uuid NOT NULL,
    mto_shipment_id uuid,
    agent_type public.mto_agents_type,
    first_name text,
    last_name text,
    email text,
    phone text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.mto_agents OWNER TO postgres;

--
-- Name: TABLE mto_agents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.mto_agents IS 'An agent is someone who can interact with movers on a customer''s behalf. There are receiving agents — people who can accept delivery at a location when the customer is not there. And releasing agents — people who can authorize a pickup from a location when the customer is not there.
Agents are assigned per shipment, not per move. The same person may be an agent for multiple shipments. An agent is not a requirement for a shipment.';


--
-- Name: COLUMN mto_agents.mto_shipment_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_agents.mto_shipment_id IS 'The shipment this particular agent applies to — a receiving agent for one shipment is not necessarily an agent for other shipments.';


--
-- Name: COLUMN mto_agents.agent_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_agents.agent_type IS 'Either RELEASING agent, or RECEIVING agent. Someone who can authorize a pickup, or who can authorize a delivery.';


--
-- Name: COLUMN mto_agents.first_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_agents.first_name IS 'First name of the agent, not the customer.';


--
-- Name: COLUMN mto_agents.last_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_agents.last_name IS 'Last name of the agent, not the customer.';


--
-- Name: COLUMN mto_agents.email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_agents.email IS 'Email contact for the agent.';


--
-- Name: COLUMN mto_agents.phone; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_agents.phone IS 'Phone number for the agent.';


--
-- Name: COLUMN mto_agents.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_agents.created_at IS 'Date & time the agent was created';


--
-- Name: COLUMN mto_agents.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_agents.updated_at IS 'Date & time the agent was last updated';


--
-- Name: COLUMN mto_agents.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_agents.deleted_at IS 'Indicates whether the mto agent has been soft deleted or not, and when it was soft deleted.';


--
-- Name: mto_service_item_customer_contacts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mto_service_item_customer_contacts (
    id uuid NOT NULL,
    type public.customer_contact_type NOT NULL,
    time_military text NOT NULL,
    first_available_delivery_date timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    date_of_contact timestamp with time zone NOT NULL
);


ALTER TABLE public.mto_service_item_customer_contacts OWNER TO postgres;

--
-- Name: TABLE mto_service_item_customer_contacts; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.mto_service_item_customer_contacts IS 'Holds the data for when the Prime contacted the customer to deliver their shipment but were unable to do so. Used to justify the Prime putting the shipment into a SIT facility.';


--
-- Name: COLUMN mto_service_item_customer_contacts.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_customer_contacts.type IS 'Either the FIRST or SECOND attempt at contacting the customer for delivery';


--
-- Name: COLUMN mto_service_item_customer_contacts.time_military; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_customer_contacts.time_military IS 'The time of attempted contact with the customer by the prime, in military format (HHMMZ), corresponding to the date_of_contact column';


--
-- Name: COLUMN mto_service_item_customer_contacts.first_available_delivery_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_customer_contacts.first_available_delivery_date IS 'The date when the Prime attempted to deliver the shipment';


--
-- Name: COLUMN mto_service_item_customer_contacts.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_customer_contacts.created_at IS 'Date & time the customer contact was created';


--
-- Name: COLUMN mto_service_item_customer_contacts.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_customer_contacts.updated_at IS 'Date & time the customer contact was last updated';


--
-- Name: COLUMN mto_service_item_customer_contacts.date_of_contact; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_customer_contacts.date_of_contact IS 'The date of attempted contact with the customer by the prime corresponding to the time_military column';


--
-- Name: mto_service_item_dimensions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mto_service_item_dimensions (
    id uuid NOT NULL,
    mto_service_item_id uuid NOT NULL,
    type public.dimension_type NOT NULL,
    length_thousandth_inches integer NOT NULL,
    height_thousandth_inches integer NOT NULL,
    width_thousandth_inches integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


ALTER TABLE public.mto_service_item_dimensions OWNER TO postgres;

--
-- Name: TABLE mto_service_item_dimensions; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.mto_service_item_dimensions IS 'The dimensions of a particular object within a particular MTO.';


--
-- Name: COLUMN mto_service_item_dimensions.mto_service_item_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_dimensions.mto_service_item_id IS 'The UUID of the service item these dimensions are associated with';


--
-- Name: COLUMN mto_service_item_dimensions.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_dimensions.type IS 'Identifies if the dimensions apply to the item being crated, or to the crate itself. (ITEM or CRATE)';


--
-- Name: COLUMN mto_service_item_dimensions.length_thousandth_inches; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_dimensions.length_thousandth_inches IS 'Length in thousandth inches. 1000 thou = 1 inch.';


--
-- Name: COLUMN mto_service_item_dimensions.height_thousandth_inches; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_dimensions.height_thousandth_inches IS 'Height in thousandth inches. 1000 thou = 1 inch.';


--
-- Name: COLUMN mto_service_item_dimensions.width_thousandth_inches; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_dimensions.width_thousandth_inches IS 'Width in thousandth inches. 1000 thou = 1 inch.';


--
-- Name: COLUMN mto_service_item_dimensions.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_dimensions.created_at IS 'Date & time the service item dimension was created';


--
-- Name: COLUMN mto_service_item_dimensions.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_item_dimensions.updated_at IS 'Date & time the service item dimension was last updated';


--
-- Name: mto_service_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.mto_service_items (
    id uuid NOT NULL,
    move_id uuid NOT NULL,
    mto_shipment_id uuid,
    re_service_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    reason text,
    pickup_postal_code text,
    description text,
    status public.service_item_status DEFAULT 'SUBMITTED'::public.service_item_status NOT NULL,
    rejection_reason text,
    approved_at timestamp without time zone,
    rejected_at timestamp without time zone,
    sit_postal_code text,
    sit_entry_date date,
    sit_departure_date date,
    sit_destination_final_address_id uuid,
    sit_origin_hhg_original_address_id uuid,
    sit_origin_hhg_actual_address_id uuid,
    estimated_weight integer,
    actual_weight integer,
    sit_destination_original_address_id uuid
);


ALTER TABLE public.mto_service_items OWNER TO postgres;

--
-- Name: TABLE mto_service_items; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.mto_service_items IS 'Service items associated with a particular MTO and shipment.';


--
-- Name: COLUMN mto_service_items.move_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.move_id IS 'The UUID of the move this service item is for';


--
-- Name: COLUMN mto_service_items.mto_shipment_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.mto_shipment_id IS 'The UUID of the shipment this service item is for - optional';


--
-- Name: COLUMN mto_service_items.re_service_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.re_service_id IS 'The UUID of the service code for this service item';


--
-- Name: COLUMN mto_service_items.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.created_at IS 'Date & time the service item was created';


--
-- Name: COLUMN mto_service_items.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.updated_at IS 'Date & time the service item was last updated';


--
-- Name: COLUMN mto_service_items.reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.reason IS 'A reason why this particular service item is justified. TXOs would use the information here to accept or reject a service item (crating, shuttling, etc.).';


--
-- Name: COLUMN mto_service_items.pickup_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.pickup_postal_code IS 'ZIP for the location where the shipment pickup is taking place.';


--
-- Name: COLUMN mto_service_items.description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.description IS 'Description of the item that the service item applies to. If it''s a request for crating, for example, this describes the item being crated (piano, moose head, etc.).';


--
-- Name: COLUMN mto_service_items.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.status IS 'The status of this service item in the review process. Can be:
1. SUBMITTED
2. APPROVED
3. REJECTED';


--
-- Name: COLUMN mto_service_items.rejection_reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.rejection_reason IS 'The reason why the TOO might have rejected this service item request';


--
-- Name: COLUMN mto_service_items.approved_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.approved_at IS 'Date & time the service item was marked as APPROVED';


--
-- Name: COLUMN mto_service_items.rejected_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.rejected_at IS 'Date & time the service item was marked as REJECTED';


--
-- Name: COLUMN mto_service_items.sit_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.sit_postal_code IS 'The postal code for the origin SIT facility where the Prime stores the shipment, used in pricing.';


--
-- Name: COLUMN mto_service_items.sit_entry_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.sit_entry_date IS 'The date when the Prime contractor places the shipment into a SIT facility. Relevant for DOFSIT and DDFSIT service items.';


--
-- Name: COLUMN mto_service_items.sit_departure_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.sit_departure_date IS 'The date when the Prime contractor removes the item from the SIT facility. Relevant for DOPSIT and DDDSIT service items.';


--
-- Name: COLUMN mto_service_items.sit_destination_final_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.sit_destination_final_address_id IS 'Final delivery address for Destination SIT';


--
-- Name: COLUMN mto_service_items.sit_origin_hhg_original_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.sit_origin_hhg_original_address_id IS 'HHG Original pickup address, using Origin SIT';


--
-- Name: COLUMN mto_service_items.sit_origin_hhg_actual_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.sit_origin_hhg_actual_address_id IS 'HHG (new) Actual pickup address, using Origin SIT';


--
-- Name: COLUMN mto_service_items.estimated_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.estimated_weight IS 'An estimate of how much weight from a shipment will be included in a shuttling (DDSHUT & DOSHUT) service item.';


--
-- Name: COLUMN mto_service_items.actual_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.actual_weight IS 'Provided by the movers, based on weight tickets. Relevant for shuttling (DDSHUT & DOSHUT) service items.';


--
-- Name: COLUMN mto_service_items.sit_destination_original_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.mto_service_items.sit_destination_original_address_id IS 'This is to capture the first sit destination address. Once this is captured, the initial address cannot be changed. Any subsequent updates to the sit destination address should be set by the sit_destination_final_address_id';


--
-- Name: notifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notifications (
    id uuid NOT NULL,
    service_member_id uuid NOT NULL,
    ses_message_id text NOT NULL,
    notification_type text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.notifications OWNER TO postgres;

--
-- Name: TABLE notifications; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.notifications IS 'Holds information about the notifications (emails) sent to customers';


--
-- Name: COLUMN notifications.service_member_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.notifications.service_member_id IS 'A foreign key that points to the service_members table';


--
-- Name: COLUMN notifications.ses_message_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.notifications.ses_message_id IS 'Uuid returned after a successful sent email message';


--
-- Name: COLUMN notifications.notification_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.notifications.notification_type IS 'The type of notification sent to the customer including: move approved, move canceled, move reviewed, move submitted, and payment reminder';


--
-- Name: COLUMN notifications.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.notifications.created_at IS 'Date & time the notification was created';


--
-- Name: office_emails; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.office_emails (
    id uuid NOT NULL,
    transportation_office_id uuid NOT NULL,
    email text NOT NULL,
    label text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.office_emails OWNER TO postgres;

--
-- Name: TABLE office_emails; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.office_emails IS 'Stores email addresses for the Transportation Offices.';


--
-- Name: COLUMN office_emails.transportation_office_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_emails.transportation_office_id IS 'A foreign key to the transportation_offices table.';


--
-- Name: COLUMN office_emails.email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_emails.email IS 'The email address for the transportation office.';


--
-- Name: COLUMN office_emails.label; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_emails.label IS 'The department the email gets sent to. For example, ''Customer Service''';


--
-- Name: COLUMN office_emails.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_emails.created_at IS 'Date & time the office_email was created.';


--
-- Name: COLUMN office_emails.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_emails.updated_at IS 'Date & time the office_email was updated.';


--
-- Name: office_phone_lines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.office_phone_lines (
    id uuid NOT NULL,
    transportation_office_id uuid NOT NULL,
    number text NOT NULL,
    label text,
    is_dsn_number boolean DEFAULT false NOT NULL,
    type text DEFAULT 'voice'::text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.office_phone_lines OWNER TO postgres;

--
-- Name: TABLE office_phone_lines; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.office_phone_lines IS 'Stores phone numbers for the Transportation Offices.';


--
-- Name: COLUMN office_phone_lines.transportation_office_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_phone_lines.transportation_office_id IS 'A foreign key to the transportation_offices table.';


--
-- Name: COLUMN office_phone_lines.number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_phone_lines.number IS 'The phone number for the transportation office.';


--
-- Name: COLUMN office_phone_lines.label; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_phone_lines.label IS 'This field is not populated locally. It''s not clear how it differs from type';


--
-- Name: COLUMN office_phone_lines.is_dsn_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_phone_lines.is_dsn_number IS 'A boolean that represents whether or not this number is a Defense Switched Network number. Defaults to false.';


--
-- Name: COLUMN office_phone_lines.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_phone_lines.type IS 'The kind of phone line, such as ''voice'' or ''fax''';


--
-- Name: COLUMN office_phone_lines.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_phone_lines.created_at IS 'Date & time the office_phone_line was created.';


--
-- Name: COLUMN office_phone_lines.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_phone_lines.updated_at IS 'Date & time the office_phone_line was updated.';


--
-- Name: office_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.office_users (
    id uuid NOT NULL,
    user_id uuid,
    last_name text NOT NULL,
    first_name text NOT NULL,
    middle_initials text,
    email text NOT NULL,
    telephone text NOT NULL,
    transportation_office_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    active boolean DEFAULT false NOT NULL
);


ALTER TABLE public.office_users OWNER TO postgres;

--
-- Name: TABLE office_users; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.office_users IS 'Holds all users who have access to the office site.';


--
-- Name: COLUMN office_users.user_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_users.user_id IS 'The foreign key that points to the user id in the users table. This gets populated when the user first signs in via login.gov, which then creates the user in the users table, and the link is then made in this table.';


--
-- Name: COLUMN office_users.last_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_users.last_name IS 'The last name of the office user.';


--
-- Name: COLUMN office_users.first_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_users.first_name IS 'The first name of the office user.';


--
-- Name: COLUMN office_users.middle_initials; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_users.middle_initials IS 'The middle initials of the office user.';


--
-- Name: COLUMN office_users.email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_users.email IS 'The email of the office user. This will match their login_gov_email in the users table.';


--
-- Name: COLUMN office_users.telephone; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_users.telephone IS 'The phone number of the office user.';


--
-- Name: COLUMN office_users.transportation_office_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_users.transportation_office_id IS 'The id of the transportation office the office user is assigned to.';


--
-- Name: COLUMN office_users.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_users.created_at IS 'Date & time the office user was created.';


--
-- Name: COLUMN office_users.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_users.updated_at IS 'Date & time the office user was updated.';


--
-- Name: COLUMN office_users.active; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.office_users.active IS 'A boolean that determines whether or not an office user is active. Users that are not active are not allowed to access the office site. See https://github.com/transcom/mymove/wiki/create-or-deactivate-users.';


--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id uuid NOT NULL,
    service_member_id uuid NOT NULL,
    issue_date date NOT NULL,
    report_by_date date NOT NULL,
    orders_type character varying(255) NOT NULL,
    has_dependents boolean NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    uploaded_orders_id uuid NOT NULL,
    orders_number character varying(255),
    orders_type_detail character varying(255),
    status character varying(255) DEFAULT 'DRAFT'::character varying NOT NULL,
    tac character varying(255),
    department_indicator character varying(255),
    spouse_has_pro_gear boolean DEFAULT false NOT NULL,
    sac text,
    grade text,
    entitlement_id uuid,
    uploaded_amended_orders_id uuid,
    amended_orders_acknowledged_at timestamp without time zone,
    nts_tac character varying(255),
    nts_sac character varying(255),
    origin_duty_location_id uuid,
    new_duty_location_id uuid NOT NULL,
    gbloc character varying,
    supply_and_services_cost_estimate text NOT NULL,
    packing_and_shipping_instructions text NOT NULL,
    method_of_payment text NOT NULL,
    naics text NOT NULL
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: TABLE orders; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.orders IS 'A customer''s move is initiated or changed based on orders issued to them by their service. Details change based on the service, but for MilMove purposes the orders will include what type of orders they are, where the customer is going, when the customer needs to get there, and other info.';


--
-- Name: COLUMN orders.service_member_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.service_member_id IS 'Unique identifier for the customer — the person who has the orders and is moving.';


--
-- Name: COLUMN orders.issue_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.issue_date IS 'Date on which the customer''s orders were issued by their branch of service.';


--
-- Name: COLUMN orders.report_by_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.report_by_date IS 'Date by which the customer must report to their new duty station or assignment.';


--
-- Name: COLUMN orders.orders_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.orders_type IS 'MilMove supports 4 orders types: Permanent change of station (PCS), local move, retirement orders, and separation orders.
In general, the moving process starts with the job/travel orders a customer receives from their service. In the orders, information describing rank, the duration of job/training, and their assigned location will determine if their entire dependent family can come, what the customer is allowed to bring, and how those items will arrive to their new location.';


--
-- Name: COLUMN orders.has_dependents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.has_dependents IS 'Does the customer''s orders include any dependents$1';


--
-- Name: COLUMN orders.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.created_at IS 'Date & time the orders were created';


--
-- Name: COLUMN orders.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.updated_at IS 'Date & time the orders were last updated';


--
-- Name: COLUMN orders.uploaded_orders_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.uploaded_orders_id IS 'A foreign key that points to the document table';


--
-- Name: COLUMN orders.orders_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.orders_number IS 'Information found on the customer''s orders assigned by their service that uniquely identifies the document. Entered in MilMove by the counselor or TOO.';


--
-- Name: COLUMN orders.orders_type_detail; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.orders_type_detail IS 'Selected from a drop-down list. Includes more specific info about the kind of move orders the customer received.
List includes:
- Shipment of HHG permitted - PCS with TDY en route
- Shipment of HHG restricted or prohibited
- HHG restricted area-HHG prohibited
- Course of instruction 20 weeks or more
- Shipment of HHG prohibited but authorized within 20 weeks
- Delayed approval 20 weeks or more';


--
-- Name: COLUMN orders.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.status IS 'The status of the orders. Allowed values are DRAFT, SUBMITTED, APPROVED, CANCELED.';


--
-- Name: COLUMN orders.tac; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.tac IS '(For HHG shipments) Lines of accounting are specified on the customer''s move orders, issued by their branch of service, and indicate the exact accounting codes the service will use to pay for the move. The Transportation Ordering Officer adds this information to the MTO.';


--
-- Name: COLUMN orders.department_indicator; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.department_indicator IS 'Name of the service branch. NAVY_AND_MARINES, ARMY, AIR_FORCE, COAST_GUARD';


--
-- Name: COLUMN orders.spouse_has_pro_gear; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.spouse_has_pro_gear IS 'Does the spouse have any pro-gear';


--
-- Name: COLUMN orders.sac; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.sac IS '(For HHG shipments) Shipment Account Classification - used for accounting';


--
-- Name: COLUMN orders.grade; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.grade IS 'Customer''s rank. Should be found on their orders. Entered by the customer from a drop-down list. Includes "civilian employee"';


--
-- Name: COLUMN orders.entitlement_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.entitlement_id IS 'A foreign key that points to the entitlements table';


--
-- Name: COLUMN orders.uploaded_amended_orders_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.uploaded_amended_orders_id IS 'A foreign key that points to the document table for referencing amended orders';


--
-- Name: COLUMN orders.amended_orders_acknowledged_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.amended_orders_acknowledged_at IS 'A timestamp that captures when new amended orders are reviewed after a move was previously approved with original orders';


--
-- Name: COLUMN orders.nts_tac; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.nts_tac IS '(For NTS shipments) Lines of accounting are specified on the customer''s move orders, issued by their branch of service, and indicate the exact accounting codes the service will use to pay for the move. The Transportation Ordering Officer adds this information to the MTO.';


--
-- Name: COLUMN orders.nts_sac; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.nts_sac IS '(For NTS shipments) Shipment Account Classification - used for accounting';


--
-- Name: COLUMN orders.origin_duty_location_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.origin_duty_location_id IS 'Unique identifier for the duty location the customer is moving from. Not the same as the text version of the name.';


--
-- Name: COLUMN orders.new_duty_location_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.new_duty_location_id IS 'Unique identifier for the duty location the customer is being assigned to. Not the same as the text version of the name.';


--
-- Name: COLUMN orders.gbloc; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.gbloc IS 'Services Counselor office users from transportation offices in this GBLOC will see these orders in their queue.';


--
-- Name: COLUMN orders.supply_and_services_cost_estimate; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.supply_and_services_cost_estimate IS 'Context for what the costs are based on.';


--
-- Name: COLUMN orders.packing_and_shipping_instructions; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.packing_and_shipping_instructions IS 'Context for where instructions can be found.';


--
-- Name: COLUMN orders.method_of_payment; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.method_of_payment IS 'Context regarding how the payment will occur.';


--
-- Name: COLUMN orders.naics; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.naics IS 'North American Industry Classification System Code.';


--
-- Name: organizations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.organizations (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    poc_email character varying(255),
    poc_phone character varying(255),
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.organizations OWNER TO postgres;

--
-- Name: TABLE organizations; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.organizations IS 'Holds all organizations that admin users belong to.';


--
-- Name: COLUMN organizations.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.organizations.name IS 'The organization name.';


--
-- Name: COLUMN organizations.poc_email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.organizations.poc_email IS 'The email of the organization''s point of contact.';


--
-- Name: COLUMN organizations.poc_phone; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.organizations.poc_phone IS 'The phone number of the organization''s point of contact.';


--
-- Name: COLUMN organizations.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.organizations.created_at IS 'Date & time the organization was created.';


--
-- Name: COLUMN organizations.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.organizations.updated_at IS 'Date & time the organization was updated.';


--
-- Name: payment_request_to_interchange_control_numbers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_request_to_interchange_control_numbers (
    id uuid NOT NULL,
    payment_request_id uuid NOT NULL,
    interchange_control_number integer NOT NULL,
    edi_type public.edi_type NOT NULL
);


ALTER TABLE public.payment_request_to_interchange_control_numbers OWNER TO postgres;

--
-- Name: COLUMN payment_request_to_interchange_control_numbers.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_request_to_interchange_control_numbers.id IS 'The id of this record';


--
-- Name: COLUMN payment_request_to_interchange_control_numbers.payment_request_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_request_to_interchange_control_numbers.payment_request_id IS 'The id of the associated payment request';


--
-- Name: COLUMN payment_request_to_interchange_control_numbers.interchange_control_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_request_to_interchange_control_numbers.interchange_control_number IS 'The interchange control number (ICN) generated in the out going EDI 858 invoice';


--
-- Name: COLUMN payment_request_to_interchange_control_numbers.edi_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_request_to_interchange_control_numbers.edi_type IS 'EDI Type of the EDI associated with the interchange control number';


--
-- Name: payment_requests; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_requests (
    id uuid NOT NULL,
    is_final boolean NOT NULL,
    rejection_reason character varying(255),
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    move_id uuid NOT NULL,
    status public.payment_request_status DEFAULT 'PENDING'::public.payment_request_status NOT NULL,
    requested_at timestamp without time zone DEFAULT now() NOT NULL,
    reviewed_at timestamp without time zone,
    sent_to_gex_at timestamp without time zone,
    received_by_gex_at timestamp without time zone,
    paid_at timestamp without time zone,
    payment_request_number text NOT NULL,
    sequence_number integer NOT NULL,
    recalculation_of_payment_request_id uuid
);


ALTER TABLE public.payment_requests OWNER TO postgres;

--
-- Name: TABLE payment_requests; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.payment_requests IS 'Represents a payment request from the GHC prime contractor.';


--
-- Name: COLUMN payment_requests.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN payment_requests.is_final; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.is_final IS 'True if this is the final payment request for the move task order (MTO).';


--
-- Name: COLUMN payment_requests.rejection_reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.rejection_reason IS 'The reason the payment request was rejected (if it was rejected).';


--
-- Name: COLUMN payment_requests.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN payment_requests.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: COLUMN payment_requests.move_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.move_id IS 'The associated move for the payment request.';


--
-- Name: COLUMN payment_requests.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.status IS 'Track the status of the payment request through the system. PENDING by default at creation. Options: PENDING, REVIEWED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID, EDI_ERROR, DEPRECATED';


--
-- Name: COLUMN payment_requests.requested_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.requested_at IS 'Timestamp when the payment request was requested.';


--
-- Name: COLUMN payment_requests.reviewed_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.reviewed_at IS 'Timestamp when the payment request was reviewed.';


--
-- Name: COLUMN payment_requests.sent_to_gex_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.sent_to_gex_at IS 'Timestamp when the payment request was sent to GEX.';


--
-- Name: COLUMN payment_requests.received_by_gex_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.received_by_gex_at IS 'Timestamp when the payment request was received by GEX.';


--
-- Name: COLUMN payment_requests.paid_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.paid_at IS 'Timestamp when the payment request was paid.';


--
-- Name: COLUMN payment_requests.payment_request_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.payment_request_number IS 'A human-readable identifier for the payment request; format is <reference_id>-<sequence_number>.';


--
-- Name: COLUMN payment_requests.sequence_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.sequence_number IS 'The sequence number of this payment request for the associated move (the first payment request would be 1).';


--
-- Name: COLUMN payment_requests.recalculation_of_payment_request_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_requests.recalculation_of_payment_request_id IS 'Link to the older payment request that was recalculated to form this payment request (if applicable).';


--
-- Name: payment_service_item_params; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_service_item_params (
    id uuid NOT NULL,
    payment_service_item_id uuid NOT NULL,
    service_item_param_key_id uuid NOT NULL,
    value character varying(80) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.payment_service_item_params OWNER TO postgres;

--
-- Name: TABLE payment_service_item_params; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.payment_service_item_params IS 'Represents the parameters (key/value pairs) for a given service item in a payment request.';


--
-- Name: COLUMN payment_service_item_params.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_item_params.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN payment_service_item_params.payment_service_item_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_item_params.payment_service_item_id IS 'The associated service item in the payment request.';


--
-- Name: COLUMN payment_service_item_params.service_item_param_key_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_item_params.service_item_param_key_id IS 'The key for this parameter.';


--
-- Name: COLUMN payment_service_item_params.value; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_item_params.value IS 'The value for this parameter.';


--
-- Name: COLUMN payment_service_item_params.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_item_params.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN payment_service_item_params.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_item_params.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: payment_service_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payment_service_items (
    id uuid NOT NULL,
    payment_request_id uuid NOT NULL,
    status public.payment_service_item_status NOT NULL,
    price_cents integer,
    rejection_reason text,
    requested_at timestamp without time zone NOT NULL,
    approved_at timestamp without time zone,
    denied_at timestamp without time zone,
    sent_to_gex_at timestamp without time zone,
    paid_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    mto_service_item_id uuid NOT NULL,
    reference_id character varying(255) NOT NULL
);


ALTER TABLE public.payment_service_items OWNER TO postgres;

--
-- Name: TABLE payment_service_items; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.payment_service_items IS 'Represents the service items associated with a given payment request.';


--
-- Name: COLUMN payment_service_items.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN payment_service_items.payment_request_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.payment_request_id IS 'The associated payment request.';


--
-- Name: COLUMN payment_service_items.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.status IS 'The payment status of this service item; options are REQUESTED, APPROVED, DENIED, SENT_TO_GEX, PAID.';


--
-- Name: COLUMN payment_service_items.price_cents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.price_cents IS 'The calculated price in cents for this service item (as determined by the GHC rate engine).';


--
-- Name: COLUMN payment_service_items.rejection_reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.rejection_reason IS 'The reason payment for a service item was rejected (if it was rejected).';


--
-- Name: COLUMN payment_service_items.requested_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.requested_at IS 'Timestamp when payment for the service item was requested.';


--
-- Name: COLUMN payment_service_items.approved_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.approved_at IS 'Timestamp when payment for the service item was approved.';


--
-- Name: COLUMN payment_service_items.denied_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.denied_at IS 'Timestamp when payment for the service item was denied.';


--
-- Name: COLUMN payment_service_items.sent_to_gex_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.sent_to_gex_at IS 'Timestamp when payment for the service item was sent to GEX.';


--
-- Name: COLUMN payment_service_items.paid_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.paid_at IS 'Timestamp when payment for the service item was paid.';


--
-- Name: COLUMN payment_service_items.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN payment_service_items.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: COLUMN payment_service_items.mto_service_item_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.mto_service_item_id IS 'The associated MTO service item for which payment is requested.';


--
-- Name: COLUMN payment_service_items.reference_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payment_service_items.reference_id IS 'Shorter ID (used by EDI) to uniquely identify this payment service item. Format is the MTO reference ID, followed by a dash, followed by enough of the payment service item ID (without dashes) to make it unique.';


--
-- Name: personally_procured_moves; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.personally_procured_moves (
    id uuid NOT NULL,
    move_id uuid NOT NULL,
    size character varying(255),
    weight_estimate integer,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    pickup_postal_code character varying(255),
    additional_pickup_postal_code character varying(255),
    destination_postal_code character varying(255),
    days_in_storage integer,
    status character varying(255) DEFAULT 'DRAFT'::character varying NOT NULL,
    has_additional_postal_code boolean,
    has_sit boolean,
    has_requested_advance boolean DEFAULT false NOT NULL,
    advance_id uuid,
    estimated_storage_reimbursement character varying(255),
    mileage integer,
    planned_sit_max integer,
    sit_max integer,
    incentive_estimate_min integer,
    incentive_estimate_max integer,
    advance_worksheet_id uuid,
    net_weight integer,
    original_move_date date,
    actual_move_date date,
    total_sit_cost integer,
    submit_date timestamp with time zone,
    approve_date timestamp with time zone,
    reviewed_date timestamp with time zone,
    has_pro_gear public.progear_status,
    has_pro_gear_over_thousand public.progear_status
);


ALTER TABLE public.personally_procured_moves OWNER TO postgres;

--
-- Name: TABLE personally_procured_moves; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.personally_procured_moves IS 'Holds information about the personally procured moves - moves when customers move themselves';


--
-- Name: COLUMN personally_procured_moves.move_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.move_id IS 'A foreign key that points to the moves table';


--
-- Name: COLUMN personally_procured_moves.size; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.size IS 'The size of a move: Large, Medium, Small';


--
-- Name: COLUMN personally_procured_moves.weight_estimate; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.weight_estimate IS 'The estimated weight the customer think they will move';


--
-- Name: COLUMN personally_procured_moves.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.created_at IS 'Date & time the personally procured move was created';


--
-- Name: COLUMN personally_procured_moves.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.updated_at IS 'Date & time the personally procured move was last updated';


--
-- Name: COLUMN personally_procured_moves.pickup_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.pickup_postal_code IS 'The pickup (origin) zip entered during the PPM setup process. This zip is used for pricing';


--
-- Name: COLUMN personally_procured_moves.additional_pickup_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.additional_pickup_postal_code IS 'An additional zipcode if the customer needs to pick up items from another location - an office perhaps';


--
-- Name: COLUMN personally_procured_moves.destination_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.destination_postal_code IS 'The destination zipcode, which is currently the zip of the destination duty station. This zip is used for pricing';


--
-- Name: COLUMN personally_procured_moves.days_in_storage; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.days_in_storage IS 'Number of days that a customer will put their things in temporary storage - max of 90 days';


--
-- Name: COLUMN personally_procured_moves.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.status IS 'The status of the personally procured move. Values can be: DRAFT, SUBMITTED, APPROVED, COMPLETED, CANCELED, PAYMENT_REQUESTED';


--
-- Name: COLUMN personally_procured_moves.has_additional_postal_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.has_additional_postal_code IS 'A boolean to determine if the user will have an additional postal code';


--
-- Name: COLUMN personally_procured_moves.has_sit; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.has_sit IS 'A boolean to determine if the user wants to use storage in transit';


--
-- Name: COLUMN personally_procured_moves.has_requested_advance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.has_requested_advance IS 'A Boolean to determine if the requested an advance';


--
-- Name: COLUMN personally_procured_moves.advance_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.advance_id IS 'A foreign key that points to the reimbursements table';


--
-- Name: COLUMN personally_procured_moves.estimated_storage_reimbursement; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.estimated_storage_reimbursement IS 'The estimated value of the SIT reimbursements from the rate engine';


--
-- Name: COLUMN personally_procured_moves.mileage; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.mileage IS 'The mileage between the pickup postal code and destination postal code';


--
-- Name: COLUMN personally_procured_moves.planned_sit_max; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.planned_sit_max IS 'The maximum SIT reimbursement for the planned SIT duration';


--
-- Name: COLUMN personally_procured_moves.sit_max; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.sit_max IS 'Maximum SIT reimbursement for maximum SIT duration. Typically 90 days';


--
-- Name: COLUMN personally_procured_moves.incentive_estimate_min; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.incentive_estimate_min IS 'The minimum of the estimate range returned from  the rate engine';


--
-- Name: COLUMN personally_procured_moves.incentive_estimate_max; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.incentive_estimate_max IS 'The maximum of the estimate range returned from the rate engine';


--
-- Name: COLUMN personally_procured_moves.advance_worksheet_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.advance_worksheet_id IS 'A foreign key that points to the documents table';


--
-- Name: COLUMN personally_procured_moves.net_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.net_weight IS 'Total weight moved (actual). This number is the sum of (total weight - empty weight) for all weight tickets.';


--
-- Name: COLUMN personally_procured_moves.original_move_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.original_move_date IS 'The date the customer plans to move';


--
-- Name: COLUMN personally_procured_moves.actual_move_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.actual_move_date IS 'The actual date the customer moved';


--
-- Name: COLUMN personally_procured_moves.total_sit_cost; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.total_sit_cost IS 'The total cost of SIT returned from rate engine';


--
-- Name: COLUMN personally_procured_moves.submit_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.submit_date IS 'Date & time the customer submitted the PPM';


--
-- Name: COLUMN personally_procured_moves.approve_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.approve_date IS 'Date & time the office user approved a customer''s PPM';


--
-- Name: COLUMN personally_procured_moves.reviewed_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.reviewed_date IS 'Date & time the office user reviewed weight tickets and expenses entered by the customer';


--
-- Name: COLUMN personally_procured_moves.has_pro_gear; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.has_pro_gear IS 'A boolean to indicate if the customer says they have pro-gear';


--
-- Name: COLUMN personally_procured_moves.has_pro_gear_over_thousand; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.personally_procured_moves.has_pro_gear_over_thousand IS 'Does the customer have pro-gear that weighs over 1000 lbs$1 If so, that is handled differently and may require a visit from the PPO office';


--
-- Name: prime_uploads; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.prime_uploads (
    id uuid NOT NULL,
    proof_of_service_docs_id uuid NOT NULL,
    contractor_id uuid NOT NULL,
    upload_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.prime_uploads OWNER TO postgres;

--
-- Name: TABLE prime_uploads; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.prime_uploads IS 'Represents uploads made by the GHC prime contractor.';


--
-- Name: COLUMN prime_uploads.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.prime_uploads.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN prime_uploads.proof_of_service_docs_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.prime_uploads.proof_of_service_docs_id IS 'The associated set of proof of service documents this upload belongs to.';


--
-- Name: COLUMN prime_uploads.contractor_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.prime_uploads.contractor_id IS 'The associated contractor for this upload.';


--
-- Name: COLUMN prime_uploads.upload_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.prime_uploads.upload_id IS 'The associated set of metadata for this upload.';


--
-- Name: COLUMN prime_uploads.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.prime_uploads.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN prime_uploads.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.prime_uploads.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: COLUMN prime_uploads.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.prime_uploads.deleted_at IS 'Timestamp when the upload was deleted.';


--
-- Name: progear_weight_tickets; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.progear_weight_tickets (
    id uuid NOT NULL,
    ppm_shipment_id uuid NOT NULL,
    belongs_to_self boolean,
    description character varying,
    has_weight_tickets boolean,
    weight integer,
    document_id uuid NOT NULL,
    status public.ppm_document_status,
    reason character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.progear_weight_tickets OWNER TO postgres;

--
-- Name: TABLE progear_weight_tickets; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.progear_weight_tickets IS 'Stores pro-gear associated information and weight docs for a PPM shipment.';


--
-- Name: COLUMN progear_weight_tickets.ppm_shipment_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.progear_weight_tickets.ppm_shipment_id IS 'The ID of the PPM shipment that this pro-gear information relates to.';


--
-- Name: COLUMN progear_weight_tickets.belongs_to_self; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.progear_weight_tickets.belongs_to_self IS 'Indicates if this information is for the customer''s own progear, otherwise, it''s the spouse''s.';


--
-- Name: COLUMN progear_weight_tickets.description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.progear_weight_tickets.description IS 'Stores a description of the pro-gear that was moved.';


--
-- Name: COLUMN progear_weight_tickets.has_weight_tickets; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.progear_weight_tickets.has_weight_tickets IS 'Indicates if the user has a weight ticket for their pro-gear, otherwise they have a constructed weight.';


--
-- Name: COLUMN progear_weight_tickets.weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.progear_weight_tickets.weight IS 'Stores the weight of the the pro-gear in pounds.';


--
-- Name: COLUMN progear_weight_tickets.document_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.progear_weight_tickets.document_id IS 'The ID of the document that is associated with the user uploads containing the pro-gear weight.';


--
-- Name: COLUMN progear_weight_tickets.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.progear_weight_tickets.status IS 'Status of the expense, e.g. APPROVED.';


--
-- Name: COLUMN progear_weight_tickets.reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.progear_weight_tickets.reason IS 'Contains the reason an expense is excluded or rejected; otherwise null.';


--
-- Name: proof_of_service_docs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.proof_of_service_docs (
    id uuid NOT NULL,
    payment_request_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.proof_of_service_docs OWNER TO postgres;

--
-- Name: TABLE proof_of_service_docs; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.proof_of_service_docs IS 'Ties together a set of uploads as proof of service documents.';


--
-- Name: COLUMN proof_of_service_docs.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.proof_of_service_docs.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN proof_of_service_docs.payment_request_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.proof_of_service_docs.payment_request_id IS 'The associated payment request that these proof of service documents support.';


--
-- Name: COLUMN proof_of_service_docs.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.proof_of_service_docs.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN proof_of_service_docs.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.proof_of_service_docs.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: pws_violations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pws_violations (
    id uuid NOT NULL,
    display_order integer,
    paragraph_number text,
    title text,
    category text,
    sub_category public.sub_category_type,
    requirement_summary text,
    requirement_statement text,
    is_kpi boolean,
    additional_data_elem text
);


ALTER TABLE public.pws_violations OWNER TO postgres;

--
-- Name: TABLE pws_violations; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.pws_violations IS 'Contains PWS violations used in the QAE evaluation reports.';


--
-- Name: COLUMN pws_violations.paragraph_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.pws_violations.paragraph_number IS 'Paragraph number the violation relates to (1.2.3.4.5)';


--
-- Name: COLUMN pws_violations.title; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.pws_violations.title IS 'Paragraph title of the violation';


--
-- Name: COLUMN pws_violations.category; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.pws_violations.category IS 'Top level category of the violation (Pre-Move Services, Physical Move Services, etc.)';


--
-- Name: COLUMN pws_violations.sub_category; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.pws_violations.sub_category IS 'Sub-category of the violation (Customer Support, Counseling, etc.)';


--
-- Name: COLUMN pws_violations.requirement_summary; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.pws_violations.requirement_summary IS 'PWS Requirement Summary';


--
-- Name: COLUMN pws_violations.requirement_statement; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.pws_violations.requirement_statement IS 'Requirement Statement in PWS';


--
-- Name: COLUMN pws_violations.is_kpi; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.pws_violations.is_kpi IS 'Whether the violation is a Key Performance Indicator (KPI)';


--
-- Name: COLUMN pws_violations.additional_data_elem; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.pws_violations.additional_data_elem IS 'Desired Additional Data Element (KPIs Only)';


--
-- Name: re_contract_years; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_contract_years (
    id uuid NOT NULL,
    contract_id uuid NOT NULL,
    name character varying(80) NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    escalation numeric(6,5) NOT NULL,
    escalation_compounded numeric(6,5) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.re_contract_years OWNER TO postgres;

--
-- Name: TABLE re_contract_years; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_contract_years IS 'Represents the "years" included in a GHC pricing contract (see sheet 5b).';


--
-- Name: COLUMN re_contract_years.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contract_years.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_contract_years.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contract_years.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: COLUMN re_contract_years.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contract_years.name IS 'The name of this contract year (e.g., "Base Period Year 1").';


--
-- Name: COLUMN re_contract_years.start_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contract_years.start_date IS 'The start date for this contract year (inclusive).';


--
-- Name: COLUMN re_contract_years.end_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contract_years.end_date IS 'The end date for this contract year (inclusive).';


--
-- Name: COLUMN re_contract_years.escalation; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contract_years.escalation IS 'The escalation factor for this specific contract year.';


--
-- Name: COLUMN re_contract_years.escalation_compounded; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contract_years.escalation_compounded IS 'The compounded escalation factor after applying previous year''s escalations.';


--
-- Name: COLUMN re_contract_years.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contract_years.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_contract_years.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contract_years.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_contracts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_contracts (
    id uuid NOT NULL,
    code character varying(80) NOT NULL,
    name character varying(80) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.re_contracts OWNER TO postgres;

--
-- Name: TABLE re_contracts; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_contracts IS 'Represents a GHC pricing contract; helps to tie together all data in that contract.';


--
-- Name: COLUMN re_contracts.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contracts.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_contracts.code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contracts.code IS 'A short, human-readable code that uniquely identifies a contract.';


--
-- Name: COLUMN re_contracts.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contracts.name IS 'A longer, more descriptive name for the contract.';


--
-- Name: COLUMN re_contracts.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contracts.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_contracts.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_contracts.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_domestic_accessorial_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_domestic_accessorial_prices (
    id uuid NOT NULL,
    contract_id uuid NOT NULL,
    service_id uuid NOT NULL,
    services_schedule integer NOT NULL,
    per_unit_cents integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT re_domestic_accessorial_prices_services_schedule_check CHECK (((services_schedule >= 1) AND (services_schedule <= 3)))
);


ALTER TABLE public.re_domestic_accessorial_prices OWNER TO postgres;

--
-- Name: TABLE re_domestic_accessorial_prices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_domestic_accessorial_prices IS 'Stores baseline prices for domestic accessorials for a GHC pricing contract (see sheet 5a).';


--
-- Name: COLUMN re_domestic_accessorial_prices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_accessorial_prices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_domestic_accessorial_prices.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_accessorial_prices.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: COLUMN re_domestic_accessorial_prices.service_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_accessorial_prices.service_id IS 'The associated service being priced.';


--
-- Name: COLUMN re_domestic_accessorial_prices.services_schedule; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_accessorial_prices.services_schedule IS 'The services schedule (1, 2, or 3, based on location) for this price.';


--
-- Name: COLUMN re_domestic_accessorial_prices.per_unit_cents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_accessorial_prices.per_unit_cents IS 'The price in cents, per unit of measure, for the service.';


--
-- Name: COLUMN re_domestic_accessorial_prices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_accessorial_prices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_domestic_accessorial_prices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_accessorial_prices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_domestic_linehaul_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_domestic_linehaul_prices (
    id uuid NOT NULL,
    contract_id uuid NOT NULL,
    weight_lower integer NOT NULL,
    weight_upper integer NOT NULL,
    miles_lower integer NOT NULL,
    miles_upper integer NOT NULL,
    is_peak_period boolean NOT NULL,
    domestic_service_area_id uuid NOT NULL,
    price_millicents integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.re_domestic_linehaul_prices OWNER TO postgres;

--
-- Name: TABLE re_domestic_linehaul_prices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_domestic_linehaul_prices IS 'Stores baseline prices for domestic linehaul for a GHC pricing contract (see sheet 2a).';


--
-- Name: COLUMN re_domestic_linehaul_prices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_domestic_linehaul_prices.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: COLUMN re_domestic_linehaul_prices.weight_lower; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.weight_lower IS 'The lower bound of shipment weight (inclusive) for this price.';


--
-- Name: COLUMN re_domestic_linehaul_prices.weight_upper; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.weight_upper IS 'The upper bound of shipment weight (inclusive) for this price.';


--
-- Name: COLUMN re_domestic_linehaul_prices.miles_lower; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.miles_lower IS 'The lower bound of miles traveled (inclusive) for this price.';


--
-- Name: COLUMN re_domestic_linehaul_prices.miles_upper; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.miles_upper IS 'The upper bound of miles traveled (inclusive) for this price.';


--
-- Name: COLUMN re_domestic_linehaul_prices.is_peak_period; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.is_peak_period IS 'Is this a peak period move$1  Peak is May 15-Sept 30.';


--
-- Name: COLUMN re_domestic_linehaul_prices.domestic_service_area_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.domestic_service_area_id IS 'The domestic service area (based on zip3) for this price.';


--
-- Name: COLUMN re_domestic_linehaul_prices.price_millicents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.price_millicents IS 'The price in millicents per hundred weight (CWT) per mile.';


--
-- Name: COLUMN re_domestic_linehaul_prices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_domestic_linehaul_prices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_linehaul_prices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_domestic_other_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_domestic_other_prices (
    id uuid NOT NULL,
    contract_id uuid NOT NULL,
    service_id uuid NOT NULL,
    is_peak_period boolean NOT NULL,
    schedule integer NOT NULL,
    price_cents integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT re_domestic_other_prices_schedule_check CHECK (((schedule >= 1) AND (schedule <= 3)))
);


ALTER TABLE public.re_domestic_other_prices OWNER TO postgres;

--
-- Name: TABLE re_domestic_other_prices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_domestic_other_prices IS 'Stores baseline prices for other domestic services for a GHC pricing contract (see sheet 2c).';


--
-- Name: COLUMN re_domestic_other_prices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_other_prices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_domestic_other_prices.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_other_prices.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: COLUMN re_domestic_other_prices.service_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_other_prices.service_id IS 'The associated service being priced.';


--
-- Name: COLUMN re_domestic_other_prices.is_peak_period; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_other_prices.is_peak_period IS 'Is this a peak period move$1  Peak is May 15-Sept 30.';


--
-- Name: COLUMN re_domestic_other_prices.schedule; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_other_prices.schedule IS 'The services schedule (1, 2, or 3, based on location) for this price.';


--
-- Name: COLUMN re_domestic_other_prices.price_cents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_other_prices.price_cents IS 'The price in cents per hundred weight (CWT).';


--
-- Name: COLUMN re_domestic_other_prices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_other_prices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_domestic_other_prices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_other_prices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_domestic_service_area_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_domestic_service_area_prices (
    id uuid NOT NULL,
    contract_id uuid NOT NULL,
    service_id uuid NOT NULL,
    is_peak_period boolean NOT NULL,
    domestic_service_area_id uuid NOT NULL,
    price_cents integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.re_domestic_service_area_prices OWNER TO postgres;

--
-- Name: TABLE re_domestic_service_area_prices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_domestic_service_area_prices IS 'Stores baseline prices for services within a domestic service area for a GHC pricing contract (see sheet 2b).';


--
-- Name: COLUMN re_domestic_service_area_prices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_area_prices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_domestic_service_area_prices.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_area_prices.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: COLUMN re_domestic_service_area_prices.service_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_area_prices.service_id IS 'The associated service being priced.';


--
-- Name: COLUMN re_domestic_service_area_prices.is_peak_period; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_area_prices.is_peak_period IS 'Is this a peak period move$1  Peak is May 15-Sept 30.';


--
-- Name: COLUMN re_domestic_service_area_prices.domestic_service_area_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_area_prices.domestic_service_area_id IS 'The domestic service area (based on zip3) for this price.';


--
-- Name: COLUMN re_domestic_service_area_prices.price_cents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_area_prices.price_cents IS 'The price in cents. Some services are per hundred weight (CWT) per mile while others are just per hundred weight. See pricing template for details.';


--
-- Name: COLUMN re_domestic_service_area_prices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_area_prices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_domestic_service_area_prices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_area_prices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_domestic_service_areas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_domestic_service_areas (
    id uuid NOT NULL,
    service_area character varying(80) NOT NULL,
    services_schedule integer NOT NULL,
    sit_pd_schedule integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    contract_id uuid NOT NULL,
    CONSTRAINT re_domestic_service_areas_services_schedule_check CHECK (((services_schedule >= 1) AND (services_schedule <= 3))),
    CONSTRAINT re_domestic_service_areas_sit_pd_schedule_check CHECK (((sit_pd_schedule >= 1) AND (sit_pd_schedule <= 3)))
);


ALTER TABLE public.re_domestic_service_areas OWNER TO postgres;

--
-- Name: TABLE re_domestic_service_areas; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_domestic_service_areas IS 'Represents the domestic service areas defined in a GHC pricing contract (see sheet 1b).';


--
-- Name: COLUMN re_domestic_service_areas.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_areas.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_domestic_service_areas.service_area; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_areas.service_area IS 'A 3-digit code uniquely identifying a service area (e.g., 004, 344).';


--
-- Name: COLUMN re_domestic_service_areas.services_schedule; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_areas.services_schedule IS 'The services schedule (1, 2, or 3) for this service area.';


--
-- Name: COLUMN re_domestic_service_areas.sit_pd_schedule; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_areas.sit_pd_schedule IS 'The SIT (Storage In Transit) pickup/delivery schedule (1, 2, or 3) for this service area.';


--
-- Name: COLUMN re_domestic_service_areas.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_areas.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_domestic_service_areas.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_areas.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: COLUMN re_domestic_service_areas.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_domestic_service_areas.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: re_intl_accessorial_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_intl_accessorial_prices (
    id uuid NOT NULL,
    contract_id uuid NOT NULL,
    service_id uuid NOT NULL,
    market character varying(1) NOT NULL,
    per_unit_cents integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT re_intl_accessorial_prices_market_check CHECK (((market)::text = ANY (ARRAY[('C'::character varying)::text, ('O'::character varying)::text])))
);


ALTER TABLE public.re_intl_accessorial_prices OWNER TO postgres;

--
-- Name: TABLE re_intl_accessorial_prices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_intl_accessorial_prices IS 'Stores baseline prices for international accessorials for a GHC pricing contract (see sheet 5a).';


--
-- Name: COLUMN re_intl_accessorial_prices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_accessorial_prices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_intl_accessorial_prices.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_accessorial_prices.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: COLUMN re_intl_accessorial_prices.service_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_accessorial_prices.service_id IS 'The associated service being priced.';


--
-- Name: COLUMN re_intl_accessorial_prices.market; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_accessorial_prices.market IS 'The market (CONUS or OCONUS) for this price.';


--
-- Name: COLUMN re_intl_accessorial_prices.per_unit_cents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_accessorial_prices.per_unit_cents IS 'The price in cents, per unit of measure, for the service.';


--
-- Name: COLUMN re_intl_accessorial_prices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_accessorial_prices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_intl_accessorial_prices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_accessorial_prices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_intl_other_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_intl_other_prices (
    id uuid NOT NULL,
    contract_id uuid NOT NULL,
    service_id uuid NOT NULL,
    is_peak_period boolean NOT NULL,
    rate_area_id uuid NOT NULL,
    per_unit_cents integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.re_intl_other_prices OWNER TO postgres;

--
-- Name: TABLE re_intl_other_prices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_intl_other_prices IS 'Stores baseline prices for other international services for a GHC pricing contract (see sheet 3d).';


--
-- Name: COLUMN re_intl_other_prices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_other_prices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_intl_other_prices.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_other_prices.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: COLUMN re_intl_other_prices.service_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_other_prices.service_id IS 'The associated service being priced.';


--
-- Name: COLUMN re_intl_other_prices.is_peak_period; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_other_prices.is_peak_period IS 'Is this a peak period move$1  Peak is May 15-Sept 30.';


--
-- Name: COLUMN re_intl_other_prices.rate_area_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_other_prices.rate_area_id IS 'The rate area (based on location) for this price.';


--
-- Name: COLUMN re_intl_other_prices.per_unit_cents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_other_prices.per_unit_cents IS 'The price in cents. Some services are per hundred weight (CWT); others are per hundred weight per mile. See pricing template for details.';


--
-- Name: COLUMN re_intl_other_prices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_other_prices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_intl_other_prices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_other_prices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_intl_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_intl_prices (
    id uuid NOT NULL,
    contract_id uuid NOT NULL,
    service_id uuid NOT NULL,
    is_peak_period boolean NOT NULL,
    origin_rate_area_id uuid NOT NULL,
    destination_rate_area_id uuid NOT NULL,
    per_unit_cents integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.re_intl_prices OWNER TO postgres;

--
-- Name: TABLE re_intl_prices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_intl_prices IS 'Stores baseline prices for international services for a GHC pricing contract (see sheets 3a-3c).';


--
-- Name: COLUMN re_intl_prices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_prices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_intl_prices.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_prices.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: COLUMN re_intl_prices.service_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_prices.service_id IS 'The associated service being priced.';


--
-- Name: COLUMN re_intl_prices.is_peak_period; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_prices.is_peak_period IS 'Is this a peak period move$1  Peak is May 15-Sept 30.';


--
-- Name: COLUMN re_intl_prices.origin_rate_area_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_prices.origin_rate_area_id IS 'The origin rate area (based on location) for this price.';


--
-- Name: COLUMN re_intl_prices.destination_rate_area_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_prices.destination_rate_area_id IS 'The destination rate area (based on location) for this price.';


--
-- Name: COLUMN re_intl_prices.per_unit_cents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_prices.per_unit_cents IS 'The price in cents per hundred weight (CWT).';


--
-- Name: COLUMN re_intl_prices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_prices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_intl_prices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_intl_prices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_rate_areas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_rate_areas (
    id uuid NOT NULL,
    is_oconus boolean NOT NULL,
    code character varying(20) NOT NULL,
    name character varying(80) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    contract_id uuid NOT NULL
);


ALTER TABLE public.re_rate_areas OWNER TO postgres;

--
-- Name: TABLE re_rate_areas; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_rate_areas IS 'Represents the rate areas defined in a GHC pricing contract (see sheets 3a-3e).';


--
-- Name: COLUMN re_rate_areas.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_rate_areas.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_rate_areas.is_oconus; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_rate_areas.is_oconus IS 'Is this rate area for an OCONUS location$1';


--
-- Name: COLUMN re_rate_areas.code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_rate_areas.code IS 'A short alphanumeric code uniquely identifying a rate area (e.g., AR, GR29, US13).';


--
-- Name: COLUMN re_rate_areas.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_rate_areas.name IS 'A descriptive name for the rate area.';


--
-- Name: COLUMN re_rate_areas.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_rate_areas.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_rate_areas.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_rate_areas.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: COLUMN re_rate_areas.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_rate_areas.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: re_services; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_services (
    id uuid NOT NULL,
    code character varying(20) NOT NULL,
    name character varying(80) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    priority integer DEFAULT 99 NOT NULL
);


ALTER TABLE public.re_services OWNER TO postgres;

--
-- Name: TABLE re_services; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_services IS 'Represents the move-related services that are included in a GHC pricing contract.';


--
-- Name: COLUMN re_services.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_services.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_services.code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_services.code IS 'A short alphabetical code uniquely identifying a service (e.g., DLH, FSC)';


--
-- Name: COLUMN re_services.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_services.name IS 'A descriptive name for the service.';


--
-- Name: COLUMN re_services.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_services.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_services.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_services.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: COLUMN re_services.priority; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_services.priority IS 'The priority of this service in a payment request; a lower number indicates a higher priority (i.e., should be priced first).';


--
-- Name: re_shipment_type_prices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_shipment_type_prices (
    id uuid NOT NULL,
    contract_id uuid NOT NULL,
    service_id uuid NOT NULL,
    market character varying(1) NOT NULL,
    factor numeric(4,2) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT re_shipment_type_prices_market_check CHECK (((market)::text = ANY (ARRAY[('C'::character varying)::text, ('O'::character varying)::text])))
);


ALTER TABLE public.re_shipment_type_prices OWNER TO postgres;

--
-- Name: TABLE re_shipment_type_prices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_shipment_type_prices IS 'Stores baseline prices for services associated with a shipment type for a GHC pricing contract (see sheet 5a).';


--
-- Name: COLUMN re_shipment_type_prices.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_shipment_type_prices.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_shipment_type_prices.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_shipment_type_prices.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: COLUMN re_shipment_type_prices.service_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_shipment_type_prices.service_id IS 'The associated service being priced.';


--
-- Name: COLUMN re_shipment_type_prices.market; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_shipment_type_prices.market IS 'The market (CONUS or OCONUS) for this price.';


--
-- Name: COLUMN re_shipment_type_prices.factor; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_shipment_type_prices.factor IS 'The price factor. Other domestic/international prices are multiplied by this factor. See pricing template for details.';


--
-- Name: COLUMN re_shipment_type_prices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_shipment_type_prices.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_shipment_type_prices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_shipment_type_prices.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_task_order_fees; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_task_order_fees (
    id uuid NOT NULL,
    contract_year_id uuid NOT NULL,
    service_id uuid NOT NULL,
    price_cents integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.re_task_order_fees OWNER TO postgres;

--
-- Name: TABLE re_task_order_fees; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_task_order_fees IS 'Stores prices for services associated with a task order for a GHC pricing contract (see sheet 4a).';


--
-- Name: COLUMN re_task_order_fees.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_task_order_fees.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_task_order_fees.contract_year_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_task_order_fees.contract_year_id IS 'The associated contract year.';


--
-- Name: COLUMN re_task_order_fees.service_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_task_order_fees.service_id IS 'The associated service being priced.';


--
-- Name: COLUMN re_task_order_fees.price_cents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_task_order_fees.price_cents IS 'The price in cents per task order. Note that price escalations do not apply. See pricing template for details.';


--
-- Name: COLUMN re_task_order_fees.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_task_order_fees.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_task_order_fees.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_task_order_fees.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: re_zip3s; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_zip3s (
    id uuid NOT NULL,
    zip3 character varying(3) NOT NULL,
    domestic_service_area_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    contract_id uuid NOT NULL,
    rate_area_id uuid,
    has_multiple_rate_areas boolean DEFAULT false NOT NULL,
    base_point_city character varying(80) NOT NULL,
    state character varying(80) NOT NULL
);


ALTER TABLE public.re_zip3s OWNER TO postgres;

--
-- Name: TABLE re_zip3s; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_zip3s IS 'Represents the zip3s defined in a GHC pricing contract (see sheet 1b) along with their associated service/rate areas.';


--
-- Name: COLUMN re_zip3s.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip3s.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_zip3s.zip3; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip3s.zip3 IS 'The first three digits of a zip code.';


--
-- Name: COLUMN re_zip3s.domestic_service_area_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip3s.domestic_service_area_id IS 'The associated domestic service area for this zip3.';


--
-- Name: COLUMN re_zip3s.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip3s.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_zip3s.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip3s.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: COLUMN re_zip3s.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip3s.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: COLUMN re_zip3s.rate_area_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip3s.rate_area_id IS 'The associated rate area for this zip3.';


--
-- Name: COLUMN re_zip3s.has_multiple_rate_areas; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip3s.has_multiple_rate_areas IS 'True if this zip3 has multiple rate areas within it; if true, see the re_zip5_rate_areas table to determine the rate area.';


--
-- Name: COLUMN re_zip3s.base_point_city; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip3s.base_point_city IS 'The name of the base point (primary) city associated with this zip3.';


--
-- Name: COLUMN re_zip3s.state; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip3s.state IS 'The state for the base point city.';


--
-- Name: re_zip5_rate_areas; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.re_zip5_rate_areas (
    id uuid NOT NULL,
    rate_area_id uuid NOT NULL,
    zip5 text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    contract_id uuid NOT NULL
);


ALTER TABLE public.re_zip5_rate_areas OWNER TO postgres;

--
-- Name: TABLE re_zip5_rate_areas; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.re_zip5_rate_areas IS 'Given a zip3 that has multiple rate areas, this table will associate the more-specific zip5 in that zip3 with a rate area.';


--
-- Name: COLUMN re_zip5_rate_areas.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip5_rate_areas.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN re_zip5_rate_areas.rate_area_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip5_rate_areas.rate_area_id IS 'The associated rate area for this zip5.';


--
-- Name: COLUMN re_zip5_rate_areas.zip5; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip5_rate_areas.zip5 IS 'The full five-digit zip code.';


--
-- Name: COLUMN re_zip5_rate_areas.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip5_rate_areas.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN re_zip5_rate_areas.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip5_rate_areas.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: COLUMN re_zip5_rate_areas.contract_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.re_zip5_rate_areas.contract_id IS 'The associated GHC pricing contract.';


--
-- Name: report_violations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.report_violations (
    id uuid NOT NULL,
    report_id uuid NOT NULL,
    violation_id uuid NOT NULL
);


ALTER TABLE public.report_violations OWNER TO postgres;

--
-- Name: TABLE report_violations; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.report_violations IS 'Associates PWS Violations with QAE evaluation report.';


--
-- Name: COLUMN report_violations.report_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.report_violations.report_id IS 'Report ID of the report that violations will be assiocated with.';


--
-- Name: COLUMN report_violations.violation_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.report_violations.violation_id IS 'Violation ID of the violation that will be assiocated to a report.';


--
-- Name: reweighs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reweighs (
    id uuid NOT NULL,
    shipment_id uuid NOT NULL,
    requested_at timestamp with time zone NOT NULL,
    requested_by public.reweigh_requester NOT NULL,
    weight integer,
    verification_reason text,
    verification_provided_at timestamp with time zone,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.reweighs OWNER TO postgres;

--
-- Name: TABLE reweighs; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.reweighs IS 'A reweigh represents a request from different users or the system for a shipment to be reweighed by the movers';


--
-- Name: COLUMN reweighs.shipment_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.reweighs.shipment_id IS 'A foreign key that points to the mto_shipments table for which shipment is being reweighed. There should only be one reweigh request per shipment';


--
-- Name: COLUMN reweighs.requested_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.reweighs.requested_at IS 'The date and time when the reweigh request was initiated';


--
-- Name: COLUMN reweighs.requested_by; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.reweighs.requested_by IS 'The type of user who requested the reweigh, including automated requests determined by the milmove system';


--
-- Name: COLUMN reweighs.weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.reweighs.weight IS 'The reweighed weight in pounds (lbs) of the shipment submitted by the movers';


--
-- Name: COLUMN reweighs.verification_reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.reweighs.verification_reason IS 'If a reweigh was requested but was not able to be performed the movers can provide an explanation';


--
-- Name: COLUMN reweighs.verification_provided_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.reweighs.verification_provided_at IS 'The date and time when the verification_reason value was added';


--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id uuid NOT NULL,
    role_type text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    role_name character varying(255) NOT NULL
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- Name: TABLE roles; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.roles IS 'Holds all roles that users can have.';


--
-- Name: COLUMN roles.role_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.roles.role_type IS 'These are the names of the roles in snake case.';


--
-- Name: COLUMN roles.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.roles.created_at IS 'Date & time the role was created.';


--
-- Name: COLUMN roles.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.roles.updated_at IS 'Date & time the role was updated.';


--
-- Name: COLUMN roles.role_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.roles.role_name IS 'The reader-friendly capitalized name of the role.';


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: TABLE schema_migration; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.schema_migration IS 'Stores the version (a date stamp in our case) for the database migrations that have been applied to this database.';


--
-- Name: COLUMN schema_migration.version; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.schema_migration.version IS 'A unique version string for the migration; derived from the first part of the migration filename.';


--
-- Name: service_item_param_keys; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_item_param_keys (
    id uuid NOT NULL,
    key character varying(80) NOT NULL,
    description character varying(255) NOT NULL,
    type public.service_item_param_type NOT NULL,
    origin public.service_item_param_origin NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.service_item_param_keys OWNER TO postgres;

--
-- Name: TABLE service_item_param_keys; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.service_item_param_keys IS 'Represents the keys for parameters that can be associated to a move-related service.';


--
-- Name: COLUMN service_item_param_keys.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_item_param_keys.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN service_item_param_keys.key; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_item_param_keys.key IS 'A short, human-readable string for the parameter.';


--
-- Name: COLUMN service_item_param_keys.description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_item_param_keys.description IS 'A descriptive name for the parameter.';


--
-- Name: COLUMN service_item_param_keys.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_item_param_keys.type IS 'The type of the value associated with this key; options are STRING, DATE, INTEGER, DECIMAL, TIMESTAMP, PaymentServiceItemUUID.';


--
-- Name: COLUMN service_item_param_keys.origin; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_item_param_keys.origin IS 'Where values for this key originate; options are PRIME (the GHC prime contractor provides) or SYSTEM (the system determines the value).';


--
-- Name: COLUMN service_item_param_keys.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_item_param_keys.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN service_item_param_keys.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_item_param_keys.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: service_items_customer_contacts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_items_customer_contacts (
    id uuid NOT NULL,
    mtoservice_item_id uuid NOT NULL,
    mtoservice_item_customer_contact_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.service_items_customer_contacts OWNER TO postgres;

--
-- Name: service_members; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_members (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    edipi text,
    affiliation text,
    rank text,
    first_name text,
    middle_name text,
    last_name text,
    suffix text,
    telephone text,
    secondary_telephone text,
    personal_email text,
    phone_is_preferred boolean,
    email_is_preferred boolean,
    residential_address_id uuid,
    backup_mailing_address_id uuid,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    duty_location_id uuid
);


ALTER TABLE public.service_members OWNER TO postgres;

--
-- Name: TABLE service_members; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.service_members IS 'Holds information about a customer';


--
-- Name: COLUMN service_members.user_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.user_id IS 'A foreign key that points to the users table';


--
-- Name: COLUMN service_members.edipi; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.edipi IS 'The customer''s Department of Defense ID number, which is used as their unique ID.';


--
-- Name: COLUMN service_members.affiliation; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.affiliation IS 'The customer''s branch of service';


--
-- Name: COLUMN service_members.rank; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.rank IS 'The customer''s rank';


--
-- Name: COLUMN service_members.first_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.first_name IS 'The customer''s first name';


--
-- Name: COLUMN service_members.middle_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.middle_name IS 'The customer''s middle name';


--
-- Name: COLUMN service_members.last_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.last_name IS 'The customer''s last name';


--
-- Name: COLUMN service_members.suffix; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.suffix IS 'The customer''s suffix';


--
-- Name: COLUMN service_members.telephone; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.telephone IS 'The customer''s phone number';


--
-- Name: COLUMN service_members.secondary_telephone; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.secondary_telephone IS 'The customer''s secondary phone number';


--
-- Name: COLUMN service_members.personal_email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.personal_email IS 'The customer''s email address';


--
-- Name: COLUMN service_members.phone_is_preferred; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.phone_is_preferred IS 'Does the customer prefer a phone call';


--
-- Name: COLUMN service_members.email_is_preferred; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.email_is_preferred IS 'Does the customer prefer an email';


--
-- Name: COLUMN service_members.residential_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.residential_address_id IS 'A foreign key that points to the addresses table - containing the customer''s residential address';


--
-- Name: COLUMN service_members.backup_mailing_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.backup_mailing_address_id IS 'A foreign key that points to the addresses table - containing the customer''s backup mailing address';


--
-- Name: COLUMN service_members.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.created_at IS 'Date & time the customer was created';


--
-- Name: COLUMN service_members.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.updated_at IS 'Date & time the customer was last updated';


--
-- Name: COLUMN service_members.duty_location_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_members.duty_location_id IS 'A foreign key that points to the duty location table - containing the customer''s current duty location';


--
-- Name: service_params; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_params (
    id uuid NOT NULL,
    service_id uuid NOT NULL,
    service_item_param_key_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    is_optional boolean DEFAULT false NOT NULL
);


ALTER TABLE public.service_params OWNER TO postgres;

--
-- Name: TABLE service_params; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.service_params IS 'Associates services with their expected input parameter keys.';


--
-- Name: COLUMN service_params.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_params.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN service_params.service_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_params.service_id IS 'The associated service.';


--
-- Name: COLUMN service_params.service_item_param_key_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_params.service_item_param_key_id IS 'The associated key.';


--
-- Name: COLUMN service_params.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_params.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN service_params.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_params.updated_at IS 'Timestamp when the record was last updated.';


--
-- Name: COLUMN service_params.is_optional; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_params.is_optional IS 'True if this parameter is optional for this service item.';


--
-- Name: service_request_document_uploads; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_request_document_uploads (
    id uuid NOT NULL,
    service_request_documents_id uuid NOT NULL,
    contractor_id uuid NOT NULL,
    upload_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.service_request_document_uploads OWNER TO postgres;

--
-- Name: TABLE service_request_document_uploads; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.service_request_document_uploads IS 'Stores uploads from the Prime that represent proof of a service item request';


--
-- Name: COLUMN service_request_document_uploads.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_request_document_uploads.id IS 'uuid that represents this entity';


--
-- Name: COLUMN service_request_document_uploads.service_request_documents_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_request_document_uploads.service_request_documents_id IS 'uuid that represents the associated service request document';


--
-- Name: COLUMN service_request_document_uploads.contractor_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_request_document_uploads.contractor_id IS 'uuid that represents the contractor who provided the upload';


--
-- Name: COLUMN service_request_document_uploads.upload_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_request_document_uploads.upload_id IS 'Foreign key of the uploads table';


--
-- Name: service_request_documents; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_request_documents (
    id uuid NOT NULL,
    mto_service_item_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.service_request_documents OWNER TO postgres;

--
-- Name: TABLE service_request_documents; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.service_request_documents IS 'Associates uploads from the Prime that represent proof of a service item request';


--
-- Name: COLUMN service_request_documents.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_request_documents.id IS 'uuid that represents this entity';


--
-- Name: COLUMN service_request_documents.mto_service_item_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.service_request_documents.mto_service_item_id IS 'Foreign key of the mto_service_items table';


--
-- Name: shipment_address_updates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.shipment_address_updates (
    id uuid NOT NULL,
    shipment_id uuid NOT NULL,
    original_address_id uuid NOT NULL,
    new_address_id uuid NOT NULL,
    contractor_remarks text NOT NULL,
    status public.shipment_address_update_status NOT NULL,
    office_remarks text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.shipment_address_updates OWNER TO postgres;

--
-- Name: COLUMN shipment_address_updates.shipment_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.shipment_address_updates.shipment_id IS 'The MTO Shipment ID associated with this address update request';


--
-- Name: COLUMN shipment_address_updates.original_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.shipment_address_updates.original_address_id IS 'Original address that was approved for the shipment';


--
-- Name: COLUMN shipment_address_updates.new_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.shipment_address_updates.new_address_id IS 'New address being requested';


--
-- Name: COLUMN shipment_address_updates.contractor_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.shipment_address_updates.contractor_remarks IS 'Reason contractor is requesting change to an address that was previously approved';


--
-- Name: COLUMN shipment_address_updates.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.shipment_address_updates.status IS 'REQUESTED (must be reviewed by TOO), APPROVED (auto-approved, or approved by TOO), or REJECTED (rejected by TOO)';


--
-- Name: COLUMN shipment_address_updates.office_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.shipment_address_updates.office_remarks IS 'Remarks from office user who reviewed the request';


--
-- Name: signed_certifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.signed_certifications (
    id uuid NOT NULL,
    submitting_user_id uuid NOT NULL,
    move_id uuid NOT NULL,
    certification_text text NOT NULL,
    signature text NOT NULL,
    date timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    personally_procured_move_id uuid,
    certification_type text,
    ppm_id uuid
);


ALTER TABLE public.signed_certifications OWNER TO postgres;

--
-- Name: TABLE signed_certifications; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.signed_certifications IS 'Holds information about when the customer signed the certificate';


--
-- Name: COLUMN signed_certifications.submitting_user_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.signed_certifications.submitting_user_id IS 'A foreign key that points to the users table';


--
-- Name: COLUMN signed_certifications.move_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.signed_certifications.move_id IS 'A foreign key that points to the moves table';


--
-- Name: COLUMN signed_certifications.certification_text; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.signed_certifications.certification_text IS 'The legalese text the customer agrees to. Value is hard coded and stored in: src/scenes/Legalese/legaleseText.js -> ppmPaymentLegal';


--
-- Name: COLUMN signed_certifications.signature; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.signed_certifications.signature IS 'Currently hard coded to, CHECKBOX, coming from the frontend';


--
-- Name: COLUMN signed_certifications.date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.signed_certifications.date IS 'Date & time the customer signed';


--
-- Name: COLUMN signed_certifications.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.signed_certifications.created_at IS 'Date & time the notification was created';


--
-- Name: COLUMN signed_certifications.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.signed_certifications.updated_at IS 'Date & time the notification was last updated';


--
-- Name: COLUMN signed_certifications.personally_procured_move_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.signed_certifications.personally_procured_move_id IS 'A foreign key that points to the personally_procured_moves table';


--
-- Name: COLUMN signed_certifications.certification_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.signed_certifications.certification_type IS 'A certification type: PPM, PPM_PAYMENT, HHG';


--
-- Name: COLUMN signed_certifications.ppm_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.signed_certifications.ppm_id IS 'PPM Shipment ID to associate the signed certificate to';


--
-- Name: sit_address_updates; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sit_address_updates (
    id uuid NOT NULL,
    mto_service_item_id uuid NOT NULL,
    old_address_id uuid NOT NULL,
    new_address_id uuid NOT NULL,
    status public.sit_address_update_status NOT NULL,
    distance integer NOT NULL,
    contractor_remarks text,
    office_remarks text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.sit_address_updates OWNER TO postgres;

--
-- Name: TABLE sit_address_updates; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.sit_address_updates IS 'Stores SIT destination address change requests for approval/rejection.';


--
-- Name: COLUMN sit_address_updates.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_address_updates.id IS 'uuid that represents this entity';


--
-- Name: COLUMN sit_address_updates.mto_service_item_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_address_updates.mto_service_item_id IS 'Foreign key of the mto_service_items table';


--
-- Name: COLUMN sit_address_updates.old_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_address_updates.old_address_id IS 'Foreign key of addresses. Old address that will be replaced.';


--
-- Name: COLUMN sit_address_updates.new_address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_address_updates.new_address_id IS 'Foreign key of addresses. New address that will replace the old address';


--
-- Name: COLUMN sit_address_updates.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_address_updates.status IS 'Current status of this request. Possible enum status(es): REQUESTED - Prime made this request and distance is greater than 50 miles, REJECTED - TXO rejected this request, APPROVED - TXO approved this request';


--
-- Name: COLUMN sit_address_updates.distance; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_address_updates.distance IS 'The distance in miles between the old address and the new address. This is calculated and stored using the address zip codes.';


--
-- Name: COLUMN sit_address_updates.contractor_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_address_updates.contractor_remarks IS 'Contractor remarks for the SIT address change. Eg: "Customer reached out to me this week & let me know they want to move closer to family."';


--
-- Name: COLUMN sit_address_updates.office_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_address_updates.office_remarks IS 'TXO remarks for the SIT address change.';


--
-- Name: sit_extensions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sit_extensions (
    id uuid NOT NULL,
    mto_shipment_id uuid NOT NULL,
    request_reason public.sit_extension_request_reason NOT NULL,
    contractor_remarks character varying,
    requested_days integer NOT NULL,
    status public.sit_extension_status NOT NULL,
    approved_days integer,
    decision_date timestamp without time zone,
    office_remarks character varying,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.sit_extensions OWNER TO postgres;

--
-- Name: TABLE sit_extensions; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.sit_extensions IS 'Stores all the updates to SIT Durations that have been requested, and their details. Formerly known as SIT Extensions, SITDurationUpdates can include both increases and decreases to a SIT Duration.';


--
-- Name: COLUMN sit_extensions.mto_shipment_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_extensions.mto_shipment_id IS 'The MTO Shipment ID associated with this SIT Duration Update.';


--
-- Name: COLUMN sit_extensions.request_reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_extensions.request_reason IS 'One of a limited set of contractual reasons an Update to the SIT Duration can be requested.';


--
-- Name: COLUMN sit_extensions.contractor_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_extensions.contractor_remarks IS 'Free form remarks from the contractor about this request to update the SIT Duration.';


--
-- Name: COLUMN sit_extensions.requested_days; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_extensions.requested_days IS 'Number of requested days to extend the SIT allowance by.';


--
-- Name: COLUMN sit_extensions.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_extensions.status IS 'Status of this SIT Duration Update (Pending, Approved, or Denied).';


--
-- Name: COLUMN sit_extensions.approved_days; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_extensions.approved_days IS 'The number of days by which to update the SIT allowance. This number can be positive (increasing the SIT allowance) or negative (decreasing the SIT allowance)';


--
-- Name: COLUMN sit_extensions.decision_date; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_extensions.decision_date IS 'The date on which this request for a SIT Duration Update was approved or denied.';


--
-- Name: COLUMN sit_extensions.office_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.sit_extensions.office_remarks IS 'Any comments from the TOO on the approval or denial of this request.';


--
-- Name: storage_facilities; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.storage_facilities (
    id uuid NOT NULL,
    facility_name character varying(255) NOT NULL,
    address_id uuid NOT NULL,
    lot_number character varying(255),
    phone character varying(255),
    email character varying(255),
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.storage_facilities OWNER TO postgres;

--
-- Name: TABLE storage_facilities; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.storage_facilities IS 'Storage facilities for NTS and NTS-Release shipments';


--
-- Name: COLUMN storage_facilities.facility_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.storage_facilities.facility_name IS 'Name of storage facility';


--
-- Name: COLUMN storage_facilities.address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.storage_facilities.address_id IS 'The address of the storage facility';


--
-- Name: COLUMN storage_facilities.lot_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.storage_facilities.lot_number IS 'Lot number where goods are stored within the storage facility';


--
-- Name: COLUMN storage_facilities.phone; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.storage_facilities.phone IS 'Phone number for contacting storage facility';


--
-- Name: COLUMN storage_facilities.email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.storage_facilities.email IS 'Email address for contacting storage facility';


--
-- Name: COLUMN storage_facilities.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.storage_facilities.created_at IS 'Date & time the storage facility was created';


--
-- Name: COLUMN storage_facilities.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.storage_facilities.updated_at IS 'Date & time the storage facility was updated';


--
-- Name: COLUMN storage_facilities.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.storage_facilities.deleted_at IS 'Indicates whether the storage facility has been soft deleted or not, and when it was soft deleted.';


--
-- Name: transportation_accounting_codes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transportation_accounting_codes (
    id uuid NOT NULL,
    tac character varying(4) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.transportation_accounting_codes OWNER TO postgres;

--
-- Name: COLUMN transportation_accounting_codes.tac; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_accounting_codes.tac IS 'A 4-digit alphanumeric transportation accounting code used to look up long lines of accounting.  These values are sourced from TGET.';


--
-- Name: transportation_offices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.transportation_offices (
    id uuid NOT NULL,
    shipping_office_id uuid,
    name text NOT NULL,
    address_id uuid NOT NULL,
    latitude real NOT NULL,
    longitude real NOT NULL,
    hours text,
    services text,
    note text,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    gbloc character varying(255) DEFAULT 'XXXX'::character varying NOT NULL,
    provides_ppm_closeout boolean DEFAULT false NOT NULL
);


ALTER TABLE public.transportation_offices OWNER TO postgres;

--
-- Name: TABLE transportation_offices; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.transportation_offices IS 'Holds all known transportation offices where office users are assigned.';


--
-- Name: COLUMN transportation_offices.shipping_office_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.shipping_office_id IS 'This is a foreign key that points back to this table. This does not seem right and will be removed in a separate cleanup PR.';


--
-- Name: COLUMN transportation_offices.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.name IS 'The name of the transportation office.';


--
-- Name: COLUMN transportation_offices.address_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.address_id IS 'The id of the transportation office''s address from the addresses table.';


--
-- Name: COLUMN transportation_offices.latitude; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.latitude IS 'The latitude of the transportation office.';


--
-- Name: COLUMN transportation_offices.longitude; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.longitude IS 'The longitude of the transportation office.';


--
-- Name: COLUMN transportation_offices.hours; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.hours IS 'The hours of operation in freeform text format.';


--
-- Name: COLUMN transportation_offices.services; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.services IS 'The various services offered in freeform text format.';


--
-- Name: COLUMN transportation_offices.note; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.note IS 'Unclear what this field is used for. It is not populated locally.';


--
-- Name: COLUMN transportation_offices.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.created_at IS 'Date & time the transportation_office was created.';


--
-- Name: COLUMN transportation_offices.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.updated_at IS 'Date & time the transportation_office was updated.';


--
-- Name: COLUMN transportation_offices.gbloc; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.gbloc IS 'A 4-character code representing the geographical area this transportation office is part of. This maps to the code field in the jppso_regions table.';


--
-- Name: COLUMN transportation_offices.provides_ppm_closeout; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.transportation_offices.provides_ppm_closeout IS 'Indicates whether a transportation office provides ppm closeout or not. It is used by Army and Air Force service members';


--
-- Name: uploads; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.uploads (
    id uuid NOT NULL,
    filename text NOT NULL,
    bytes bigint NOT NULL,
    content_type text NOT NULL,
    checksum text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    storage_key character varying(1024) NOT NULL,
    deleted_at timestamp with time zone,
    upload_type public.upload_type NOT NULL
);


ALTER TABLE public.uploads OWNER TO postgres;

--
-- Name: TABLE uploads; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.uploads IS 'Holds information regarding files that are uploaded';


--
-- Name: COLUMN uploads.filename; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.uploads.filename IS 'The filename of the upload';


--
-- Name: COLUMN uploads.bytes; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.uploads.bytes IS 'The number of bytes of the upload';


--
-- Name: COLUMN uploads.content_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.uploads.content_type IS 'The mime type of the upload';


--
-- Name: COLUMN uploads.checksum; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.uploads.checksum IS 'A checksum value of the upload';


--
-- Name: COLUMN uploads.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.uploads.created_at IS 'Date & time the uploads was created';


--
-- Name: COLUMN uploads.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.uploads.updated_at IS 'Date & time the uploads was last updated';


--
-- Name: COLUMN uploads.storage_key; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.uploads.storage_key IS 'The resulting path to where the upload is on S3';


--
-- Name: COLUMN uploads.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.uploads.deleted_at IS 'Date & time of when the uploads was deleted';


--
-- Name: COLUMN uploads.upload_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.uploads.upload_type IS 'Who created the upload: USER, PRIME';


--
-- Name: user_uploads; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.user_uploads (
    id uuid NOT NULL,
    document_id uuid,
    uploader_id uuid NOT NULL,
    upload_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp with time zone
);


ALTER TABLE public.user_uploads OWNER TO postgres;

--
-- Name: TABLE user_uploads; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.user_uploads IS 'Holds information that joins the uploads to the corresponding documents and users';


--
-- Name: COLUMN user_uploads.document_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_uploads.document_id IS 'A foreign key that points to the documents table';


--
-- Name: COLUMN user_uploads.uploader_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_uploads.uploader_id IS 'A foreign key that points to the users table';


--
-- Name: COLUMN user_uploads.upload_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_uploads.upload_id IS 'A foreign key that points to the uploaded table';


--
-- Name: COLUMN user_uploads.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_uploads.created_at IS 'Date & time the user uploads was created';


--
-- Name: COLUMN user_uploads.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_uploads.updated_at IS 'Date & time the user uploads was last updated';


--
-- Name: COLUMN user_uploads.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.user_uploads.deleted_at IS 'Date & time the user uploads was deleted';


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    login_gov_uuid uuid,
    login_gov_email text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    active boolean DEFAULT false NOT NULL,
    current_mil_session_id text DEFAULT ''::text,
    current_admin_session_id text DEFAULT ''::text,
    current_office_session_id text DEFAULT ''::text
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: TABLE users; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.users IS 'Holds all users. Anyone who signs in to any of the mymove apps is automatically created in this table after signing in with login.gov.';


--
-- Name: COLUMN users.login_gov_uuid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.login_gov_uuid IS 'The login.gov uuid of the user.';


--
-- Name: COLUMN users.login_gov_email; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.login_gov_email IS 'The login.gov email of the user.';


--
-- Name: COLUMN users.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.created_at IS 'Date & time the user was created.';


--
-- Name: COLUMN users.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.updated_at IS 'Date & time the user was updated.';


--
-- Name: COLUMN users.active; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.active IS 'A boolean that determines whether or not a user is active. Users that are not active are not allowed to access the mymove apps. See https://github.com/transcom/mymove/wiki/create-or-deactivate-users.';


--
-- Name: COLUMN users.current_mil_session_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.current_mil_session_id IS 'This field gets populated when a user signs into the mil app. The string matches the session id stored in Redis. It is used to allow an admin user to revoke the session if necessary.';


--
-- Name: COLUMN users.current_admin_session_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.current_admin_session_id IS 'This field gets populated when a user signs into the admin app. The string matches the session id stored in Redis. It is used to allow an admin user to revoke the session if necessary.';


--
-- Name: COLUMN users.current_office_session_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.current_office_session_id IS 'This field gets populated when a user signs into the office app. The string matches the session id stored in Redis. It is used to allow an admin user to revoke the session if necessary.';


--
-- Name: users_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users_roles (
    user_id uuid NOT NULL,
    role_id uuid NOT NULL,
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone
);


ALTER TABLE public.users_roles OWNER TO postgres;

--
-- Name: TABLE users_roles; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.users_roles IS 'A join table between users and roles to identify which users have which roles.';


--
-- Name: COLUMN users_roles.user_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users_roles.user_id IS 'The id of the user being referenced.';


--
-- Name: COLUMN users_roles.role_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users_roles.role_id IS 'The id of the role being referenced.';


--
-- Name: COLUMN users_roles.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users_roles.created_at IS 'Date & time the users_roles was created.';


--
-- Name: COLUMN users_roles.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users_roles.updated_at IS 'Date & time the users_roles was updated.';


--
-- Name: COLUMN users_roles.deleted_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users_roles.deleted_at IS 'Date & time the users_roles was deleted.';


--
-- Name: webhook_notifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.webhook_notifications (
    id uuid NOT NULL,
    event_key text NOT NULL,
    trace_id uuid,
    move_id uuid,
    object_id uuid,
    payload json NOT NULL,
    status public.webhook_notifications_status DEFAULT 'PENDING'::public.webhook_notifications_status NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    first_attempted_at timestamp without time zone
);


ALTER TABLE public.webhook_notifications OWNER TO postgres;

--
-- Name: TABLE webhook_notifications; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.webhook_notifications IS 'Represents the notifications that will be sent to an external client about changes in our database. Used to notify the Prime when Prime-available moves and related objects have been updated.';


--
-- Name: COLUMN webhook_notifications.event_key; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_notifications.event_key IS 'A string used to identify which object this notification pertains to (PaymentRequest, MTOShipment, etc.), and how it was modified (PaymentRequest.Create, MTOShipment.Update, etc.)';


--
-- Name: COLUMN webhook_notifications.trace_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_notifications.trace_id IS 'The UUID representing the specific transaction this notification represents';


--
-- Name: COLUMN webhook_notifications.move_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_notifications.move_id IS 'The UUID for the move affected by this change';


--
-- Name: COLUMN webhook_notifications.object_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_notifications.object_id IS 'The UUID for the specific object that was modified';


--
-- Name: COLUMN webhook_notifications.payload; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_notifications.payload IS 'A JSON payload containing the updates for this object';


--
-- Name: COLUMN webhook_notifications.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_notifications.status IS 'The status of this notification. Can be:
1. PENDING
2. SENT
3. FAILED';


--
-- Name: COLUMN webhook_notifications.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_notifications.created_at IS 'Date & time the webhook notification was created';


--
-- Name: COLUMN webhook_notifications.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_notifications.updated_at IS 'Date & time the webhook notification was last updated';


--
-- Name: webhook_subscriptions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.webhook_subscriptions (
    id uuid NOT NULL,
    subscriber_id uuid NOT NULL,
    status public.webhook_subscriptions_status DEFAULT 'ACTIVE'::public.webhook_subscriptions_status,
    event_key text NOT NULL,
    callback_url text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    severity integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.webhook_subscriptions OWNER TO postgres;

--
-- Name: TABLE webhook_subscriptions; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.webhook_subscriptions IS 'Represents subscribers who expect certain notifications to be pushed to their servers. Used for the Prime and Prime-related events specifically.';


--
-- Name: COLUMN webhook_subscriptions.subscriber_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_subscriptions.subscriber_id IS 'The UUID of the Prime contractor this subscription belongs to';


--
-- Name: COLUMN webhook_subscriptions.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_subscriptions.status IS 'The status of this subscription. Can be:
1. ACTIVE
2. DISABLED';


--
-- Name: COLUMN webhook_subscriptions.event_key; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_subscriptions.event_key IS 'A string used to represent which events this subscriber expects to be notified about. Corresponds to the possible `event_key` values in `webhook_notifications`';


--
-- Name: COLUMN webhook_subscriptions.callback_url; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_subscriptions.callback_url IS 'The URL to which the notifications for this subscription will be pushed to';


--
-- Name: COLUMN webhook_subscriptions.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_subscriptions.created_at IS 'Date & time the webhook subscription was created';


--
-- Name: COLUMN webhook_subscriptions.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.webhook_subscriptions.updated_at IS 'Date & time the webhook subscription was last updated';


--
-- Name: weight_tickets; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.weight_tickets (
    id uuid NOT NULL,
    ppm_shipment_id uuid NOT NULL,
    vehicle_description character varying,
    empty_weight integer,
    missing_empty_weight_ticket boolean,
    empty_document_id uuid NOT NULL,
    full_weight integer,
    missing_full_weight_ticket boolean,
    full_document_id uuid NOT NULL,
    owns_trailer boolean,
    trailer_meets_criteria boolean,
    proof_of_trailer_ownership_document_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp with time zone,
    status public.ppm_document_status,
    reason character varying,
    adjusted_net_weight integer,
    net_weight_remarks character varying
);


ALTER TABLE public.weight_tickets OWNER TO postgres;

--
-- Name: TABLE weight_tickets; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.weight_tickets IS 'Stores weight ticket docs associated with a trip for a PPM shipment.';


--
-- Name: COLUMN weight_tickets.ppm_shipment_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.ppm_shipment_id IS 'The ID of the PPM shipment that this set of weight tickets is for.';


--
-- Name: COLUMN weight_tickets.vehicle_description; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.vehicle_description IS 'Stores a description of the vehicle used for the trip. E.g. make/model, type of truck/van, etc.';


--
-- Name: COLUMN weight_tickets.empty_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.empty_weight IS 'Stores the weight of the vehicle when empty.';


--
-- Name: COLUMN weight_tickets.missing_empty_weight_ticket; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.missing_empty_weight_ticket IS 'Indicates if the customer is missing a weight ticket for the vehicle weight.';


--
-- Name: COLUMN weight_tickets.empty_document_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.empty_document_id IS 'The ID of the document that is associated with the user uploads containing the full vehicle weight.';


--
-- Name: COLUMN weight_tickets.full_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.full_weight IS 'Stores the weight of the vehicle when full.';


--
-- Name: COLUMN weight_tickets.missing_full_weight_ticket; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.missing_full_weight_ticket IS 'Indicates if the customer is missing a weight ticket for the vehicle weight.';


--
-- Name: COLUMN weight_tickets.owns_trailer; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.owns_trailer IS 'Indicates if the customer used a trailer they own for the move.';


--
-- Name: COLUMN weight_tickets.trailer_meets_criteria; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.trailer_meets_criteria IS 'Indicates if the trailer that the customer used meets all the criteria to be claimable.';


--
-- Name: COLUMN weight_tickets.proof_of_trailer_ownership_document_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.proof_of_trailer_ownership_document_id IS 'The ID of the document that is associated with the user uploads containing the proof of trailer ownership.';


--
-- Name: COLUMN weight_tickets.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.status IS 'Status of the weight ticket, e.g. APPROVED.';


--
-- Name: COLUMN weight_tickets.reason; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.reason IS 'Contains the reason a weight ticket is excluded or rejected; otherwise null.';


--
-- Name: COLUMN weight_tickets.adjusted_net_weight; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.adjusted_net_weight IS 'Stores the net weight of the vehicle';


--
-- Name: COLUMN weight_tickets.net_weight_remarks; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.weight_tickets.net_weight_remarks IS 'Stores remarks explaining any edits made to the net weight';


--
-- Name: zip3_distances; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.zip3_distances (
    id uuid NOT NULL,
    from_zip3 character(3) NOT NULL,
    to_zip3 character(3) NOT NULL,
    distance_miles integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    CONSTRAINT zip3_distances_ordering CHECK ((from_zip3 < to_zip3))
);


ALTER TABLE public.zip3_distances OWNER TO postgres;

--
-- Name: TABLE zip3_distances; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.zip3_distances IS 'Stores the distances between zip3 pairs; there should only be one record for any zip3 pair, with from_zip3 always alphabetically before to_zip3.';


--
-- Name: COLUMN zip3_distances.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.zip3_distances.id IS 'UUID that uniquely identifies the record.';


--
-- Name: COLUMN zip3_distances.from_zip3; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.zip3_distances.from_zip3 IS 'The starting zip3; this should always be alphabetically before to_zip3.';


--
-- Name: COLUMN zip3_distances.to_zip3; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.zip3_distances.to_zip3 IS 'The ending zip3; this should always be alphabetically after from_zip3.';


--
-- Name: COLUMN zip3_distances.distance_miles; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.zip3_distances.distance_miles IS 'The distance in miles between from_zip3 and to_zip3.';


--
-- Name: COLUMN zip3_distances.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.zip3_distances.created_at IS 'Timestamp when the record was first created.';


--
-- Name: COLUMN zip3_distances.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.zip3_distances.updated_at IS 'Timestamp when the record was first updated.';


--
-- Name: CONSTRAINT zip3_distances_ordering ON zip3_distances; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON CONSTRAINT zip3_distances_ordering ON public.zip3_distances IS 'Ensures that from_zip3 is always alphabetically before to_zip3.';


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: admin_users admin_users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_pkey PRIMARY KEY (id);


--
-- Name: archived_access_codes archived_access_codes_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_access_codes
    ADD CONSTRAINT archived_access_codes_code_key UNIQUE (code);


--
-- Name: archived_access_codes archived_access_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_access_codes
    ADD CONSTRAINT archived_access_codes_pkey PRIMARY KEY (id);


--
-- Name: archived_move_documents archived_move_documents_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_move_documents
    ADD CONSTRAINT archived_move_documents_pkey PRIMARY KEY (id);


--
-- Name: archived_moving_expense_documents archived_moving_expense_documents_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_moving_expense_documents
    ADD CONSTRAINT archived_moving_expense_documents_pkey PRIMARY KEY (id);


--
-- Name: archived_personally_procured_moves archived_personally_procured_moves_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_personally_procured_moves
    ADD CONSTRAINT archived_personally_procured_moves_pkey PRIMARY KEY (id);


--
-- Name: archived_signed_certifications archived_signed_certifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_signed_certifications
    ADD CONSTRAINT archived_signed_certifications_pkey PRIMARY KEY (id);


--
-- Name: archived_weight_ticket_set_documents archived_weight_ticket_set_documents_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_weight_ticket_set_documents
    ADD CONSTRAINT archived_weight_ticket_set_documents_pkey PRIMARY KEY (id);


--
-- Name: audit_history audit_history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.audit_history
    ADD CONSTRAINT audit_history_pkey PRIMARY KEY (id);


--
-- Name: backup_contacts backup_contacts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.backup_contacts
    ADD CONSTRAINT backup_contacts_pkey PRIMARY KEY (id);


--
-- Name: client_certs client_certs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.client_certs
    ADD CONSTRAINT client_certs_pkey PRIMARY KEY (id);


--
-- Name: client_certs client_certs_sha256_digest_idx; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.client_certs
    ADD CONSTRAINT client_certs_sha256_digest_idx UNIQUE (sha256_digest);


--
-- Name: client_certs client_certs_subject_idx; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.client_certs
    ADD CONSTRAINT client_certs_subject_idx UNIQUE (subject);


--
-- Name: users constraint_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT constraint_name UNIQUE (login_gov_uuid);


--
-- Name: contractors contractors_contract_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.contractors
    ADD CONSTRAINT contractors_contract_number_key UNIQUE (contract_number);


--
-- Name: contractors contractors_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.contractors
    ADD CONSTRAINT contractors_pkey PRIMARY KEY (id);


--
-- Name: customer_support_remarks customer_support_remarks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customer_support_remarks
    ADD CONSTRAINT customer_support_remarks_pkey PRIMARY KEY (id);


--
-- Name: distance_calculations distance_calculations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.distance_calculations
    ADD CONSTRAINT distance_calculations_pkey PRIMARY KEY (id);


--
-- Name: documents documents_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_pkey PRIMARY KEY (id);


--
-- Name: duty_location_names duty_location_names_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.duty_location_names
    ADD CONSTRAINT duty_location_names_pkey PRIMARY KEY (id);


--
-- Name: duty_locations duty_locations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.duty_locations
    ADD CONSTRAINT duty_locations_pkey PRIMARY KEY (id);


--
-- Name: edi_errors edi_errors_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.edi_errors
    ADD CONSTRAINT edi_errors_pkey PRIMARY KEY (id);


--
-- Name: edi_processings edi_processings_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.edi_processings
    ADD CONSTRAINT edi_processings_key PRIMARY KEY (id);


--
-- Name: electronic_orders electronic_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.electronic_orders
    ADD CONSTRAINT electronic_orders_pkey PRIMARY KEY (id);


--
-- Name: electronic_orders_revisions electronic_orders_revisions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.electronic_orders_revisions
    ADD CONSTRAINT electronic_orders_revisions_pkey PRIMARY KEY (id);


--
-- Name: entitlements entitlements_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.entitlements
    ADD CONSTRAINT entitlements_pkey PRIMARY KEY (id);


--
-- Name: evaluation_reports evaluation_reports_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.evaluation_reports
    ADD CONSTRAINT evaluation_reports_pkey PRIMARY KEY (id);


--
-- Name: fuel_eia_diesel_prices fuel_eia_diesel_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fuel_eia_diesel_prices
    ADD CONSTRAINT fuel_eia_diesel_prices_pkey PRIMARY KEY (id);


--
-- Name: ghc_diesel_fuel_prices ghc_diesel_fuel_prices_publication_date_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ghc_diesel_fuel_prices
    ADD CONSTRAINT ghc_diesel_fuel_prices_publication_date_key UNIQUE (publication_date);


--
-- Name: ghc_domestic_transit_times ghc_domestic_transit_times_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ghc_domestic_transit_times
    ADD CONSTRAINT ghc_domestic_transit_times_pkey PRIMARY KEY (id);


--
-- Name: payment_request_to_interchange_control_numbers interchange_control_number_payment_request_id_edi_type_uniq_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_request_to_interchange_control_numbers
    ADD CONSTRAINT interchange_control_number_payment_request_id_edi_type_uniq_key UNIQUE (interchange_control_number, payment_request_id, edi_type);


--
-- Name: invoice_number_trackers invoice_number_trackers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoice_number_trackers
    ADD CONSTRAINT invoice_number_trackers_pkey PRIMARY KEY (standard_carrier_alpha_code, year);


--
-- Name: invoices invoices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_pkey PRIMARY KEY (id);


--
-- Name: jppso_region_state_assignments jppso_region_state_assignments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jppso_region_state_assignments
    ADD CONSTRAINT jppso_region_state_assignments_pkey PRIMARY KEY (id);


--
-- Name: jppso_region_state_assignments jppso_region_state_assignments_state_abbreviation_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jppso_region_state_assignments
    ADD CONSTRAINT jppso_region_state_assignments_state_abbreviation_key UNIQUE (state_abbreviation);


--
-- Name: jppso_region_state_assignments jppso_region_state_assignments_state_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jppso_region_state_assignments
    ADD CONSTRAINT jppso_region_state_assignments_state_name_key UNIQUE (state_name);


--
-- Name: jppso_regions jppso_regions_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jppso_regions
    ADD CONSTRAINT jppso_regions_code_key UNIQUE (code);


--
-- Name: jppso_regions jppso_regions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jppso_regions
    ADD CONSTRAINT jppso_regions_pkey PRIMARY KEY (id);


--
-- Name: moves moves_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.moves
    ADD CONSTRAINT moves_pkey PRIMARY KEY (id);


--
-- Name: moves moves_reference_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.moves
    ADD CONSTRAINT moves_reference_id_key UNIQUE (reference_id);


--
-- Name: moving_expenses moving_expenses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.moving_expenses
    ADD CONSTRAINT moving_expenses_pkey PRIMARY KEY (id);


--
-- Name: mto_agents mto_agents_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_agents
    ADD CONSTRAINT mto_agents_pkey PRIMARY KEY (id);


--
-- Name: mto_service_item_customer_contacts mto_service_item_customer_contacts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_item_customer_contacts
    ADD CONSTRAINT mto_service_item_customer_contacts_pkey PRIMARY KEY (id);


--
-- Name: mto_service_item_dimensions mto_service_item_dimensions_mto_service_item_id_type_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_item_dimensions
    ADD CONSTRAINT mto_service_item_dimensions_mto_service_item_id_type_key UNIQUE (mto_service_item_id, type);


--
-- Name: mto_service_item_dimensions mto_service_item_dimensions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_item_dimensions
    ADD CONSTRAINT mto_service_item_dimensions_pkey PRIMARY KEY (id);


--
-- Name: mto_service_items mto_service_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_items
    ADD CONSTRAINT mto_service_items_pkey PRIMARY KEY (id);


--
-- Name: mto_shipments mto_shipments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_shipments
    ADD CONSTRAINT mto_shipments_pkey PRIMARY KEY (id);


--
-- Name: fuel_eia_diesel_prices no_overlapping_rates; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fuel_eia_diesel_prices
    ADD CONSTRAINT no_overlapping_rates EXCLUDE USING gist (daterange(rate_start_date, rate_end_date, '[]'::text) WITH &&);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: office_emails office_emails_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.office_emails
    ADD CONSTRAINT office_emails_pkey PRIMARY KEY (id);


--
-- Name: office_phone_lines office_phone_lines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.office_phone_lines
    ADD CONSTRAINT office_phone_lines_pkey PRIMARY KEY (id);


--
-- Name: office_users office_users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.office_users
    ADD CONSTRAINT office_users_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: organizations organizations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.organizations
    ADD CONSTRAINT organizations_pkey PRIMARY KEY (id);


--
-- Name: payment_request_to_interchange_control_numbers payment_request_to_interchange_control_numbers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_request_to_interchange_control_numbers
    ADD CONSTRAINT payment_request_to_interchange_control_numbers_pkey PRIMARY KEY (id);


--
-- Name: payment_requests payment_requests_payment_request_number_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_requests
    ADD CONSTRAINT payment_requests_payment_request_number_unique_key UNIQUE (payment_request_number);


--
-- Name: payment_requests payment_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_requests
    ADD CONSTRAINT payment_requests_pkey PRIMARY KEY (id);


--
-- Name: payment_requests payment_requests_sequence_number_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_requests
    ADD CONSTRAINT payment_requests_sequence_number_unique_key UNIQUE (move_id, sequence_number);


--
-- Name: payment_service_item_params payment_service_item_params_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_service_item_params
    ADD CONSTRAINT payment_service_item_params_pkey PRIMARY KEY (id);


--
-- Name: payment_service_item_params payment_service_item_params_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_service_item_params
    ADD CONSTRAINT payment_service_item_params_unique_key UNIQUE (payment_service_item_id, service_item_param_key_id);


--
-- Name: payment_service_items payment_service_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_service_items
    ADD CONSTRAINT payment_service_items_pkey PRIMARY KEY (id);


--
-- Name: payment_service_items payment_service_items_reference_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_service_items
    ADD CONSTRAINT payment_service_items_reference_id_key UNIQUE (reference_id);


--
-- Name: personally_procured_moves personally_procured_moves_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.personally_procured_moves
    ADD CONSTRAINT personally_procured_moves_pkey PRIMARY KEY (id);


--
-- Name: postal_code_to_gblocs postal_code_to_gblocs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.postal_code_to_gblocs
    ADD CONSTRAINT postal_code_to_gblocs_pkey PRIMARY KEY (id);


--
-- Name: postal_code_to_gblocs postal_code_to_gblocs_unique_postal_code; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.postal_code_to_gblocs
    ADD CONSTRAINT postal_code_to_gblocs_unique_postal_code UNIQUE (postal_code);


--
-- Name: ppm_shipments ppm_shipments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ppm_shipments
    ADD CONSTRAINT ppm_shipments_pkey PRIMARY KEY (id);


--
-- Name: prime_uploads prime_uploads_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.prime_uploads
    ADD CONSTRAINT prime_uploads_pkey PRIMARY KEY (id);


--
-- Name: progear_weight_tickets progear_weight_tickets_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.progear_weight_tickets
    ADD CONSTRAINT progear_weight_tickets_pkey PRIMARY KEY (id);


--
-- Name: proof_of_service_docs proof_of_service_docs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.proof_of_service_docs
    ADD CONSTRAINT proof_of_service_docs_pkey PRIMARY KEY (id);


--
-- Name: pws_violations pws_violations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pws_violations
    ADD CONSTRAINT pws_violations_pkey PRIMARY KEY (id);


--
-- Name: re_contract_years re_contract_years_daterange_excl; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_contract_years
    ADD CONSTRAINT re_contract_years_daterange_excl EXCLUDE USING gist (daterange(start_date, end_date, '[]'::text) WITH &&);


--
-- Name: re_contract_years re_contract_years_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_contract_years
    ADD CONSTRAINT re_contract_years_pkey PRIMARY KEY (id);


--
-- Name: re_contracts re_contracts_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_contracts
    ADD CONSTRAINT re_contracts_code_key UNIQUE (code);


--
-- Name: re_contracts re_contracts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_contracts
    ADD CONSTRAINT re_contracts_pkey PRIMARY KEY (id);


--
-- Name: re_domestic_accessorial_prices re_domestic_accessorial_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_accessorial_prices
    ADD CONSTRAINT re_domestic_accessorial_prices_pkey PRIMARY KEY (id);


--
-- Name: re_domestic_accessorial_prices re_domestic_accessorial_prices_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_accessorial_prices
    ADD CONSTRAINT re_domestic_accessorial_prices_unique_key UNIQUE (contract_id, service_id, services_schedule);


--
-- Name: re_domestic_linehaul_prices re_domestic_linehaul_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_linehaul_prices
    ADD CONSTRAINT re_domestic_linehaul_prices_pkey PRIMARY KEY (id);


--
-- Name: re_domestic_linehaul_prices re_domestic_linehaul_prices_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_linehaul_prices
    ADD CONSTRAINT re_domestic_linehaul_prices_unique_key UNIQUE (contract_id, weight_lower, weight_upper, miles_lower, miles_upper, is_peak_period, domestic_service_area_id);


--
-- Name: re_domestic_other_prices re_domestic_other_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_other_prices
    ADD CONSTRAINT re_domestic_other_prices_pkey PRIMARY KEY (id);


--
-- Name: re_domestic_other_prices re_domestic_other_prices_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_other_prices
    ADD CONSTRAINT re_domestic_other_prices_unique_key UNIQUE (contract_id, service_id, is_peak_period, schedule);


--
-- Name: re_domestic_service_area_prices re_domestic_service_area_prices_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_service_area_prices
    ADD CONSTRAINT re_domestic_service_area_prices_unique_key UNIQUE (contract_id, service_id, is_peak_period, domestic_service_area_id);


--
-- Name: re_domestic_service_areas re_domestic_service_areas_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_service_areas
    ADD CONSTRAINT re_domestic_service_areas_pkey PRIMARY KEY (id);


--
-- Name: re_domestic_service_areas re_domestic_service_areas_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_service_areas
    ADD CONSTRAINT re_domestic_service_areas_unique_key UNIQUE (contract_id, service_area);


--
-- Name: re_domestic_service_area_prices re_domestic_services_area_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_service_area_prices
    ADD CONSTRAINT re_domestic_services_area_prices_pkey PRIMARY KEY (id);


--
-- Name: re_intl_accessorial_prices re_intl_accessorial_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_accessorial_prices
    ADD CONSTRAINT re_intl_accessorial_prices_pkey PRIMARY KEY (id);


--
-- Name: re_intl_accessorial_prices re_intl_accessorial_prices_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_accessorial_prices
    ADD CONSTRAINT re_intl_accessorial_prices_unique_key UNIQUE (contract_id, service_id, market);


--
-- Name: re_intl_other_prices re_intl_other_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_other_prices
    ADD CONSTRAINT re_intl_other_prices_pkey PRIMARY KEY (id);


--
-- Name: re_intl_other_prices re_intl_other_prices_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_other_prices
    ADD CONSTRAINT re_intl_other_prices_unique_key UNIQUE (contract_id, service_id, is_peak_period, rate_area_id);


--
-- Name: re_intl_prices re_intl_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_prices
    ADD CONSTRAINT re_intl_prices_pkey PRIMARY KEY (id);


--
-- Name: re_intl_prices re_intl_prices_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_prices
    ADD CONSTRAINT re_intl_prices_unique_key UNIQUE (contract_id, service_id, is_peak_period, origin_rate_area_id, destination_rate_area_id);


--
-- Name: re_rate_areas re_rate_areas_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_rate_areas
    ADD CONSTRAINT re_rate_areas_pkey PRIMARY KEY (id);


--
-- Name: re_rate_areas re_rate_areas_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_rate_areas
    ADD CONSTRAINT re_rate_areas_unique_key UNIQUE (contract_id, code);


--
-- Name: re_services re_services_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_services
    ADD CONSTRAINT re_services_code_key UNIQUE (code);


--
-- Name: re_services re_services_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_services
    ADD CONSTRAINT re_services_pkey PRIMARY KEY (id);


--
-- Name: re_shipment_type_prices re_shipment_type_prices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_shipment_type_prices
    ADD CONSTRAINT re_shipment_type_prices_pkey PRIMARY KEY (id);


--
-- Name: re_shipment_type_prices re_shipment_type_prices_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_shipment_type_prices
    ADD CONSTRAINT re_shipment_type_prices_unique_key UNIQUE (contract_id, service_id, market);


--
-- Name: re_task_order_fees re_task_order_fees_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_task_order_fees
    ADD CONSTRAINT re_task_order_fees_pkey PRIMARY KEY (id);


--
-- Name: re_task_order_fees re_task_order_fees_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_task_order_fees
    ADD CONSTRAINT re_task_order_fees_unique_key UNIQUE (contract_year_id, service_id);


--
-- Name: re_zip3s re_zip3s_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_zip3s
    ADD CONSTRAINT re_zip3s_pkey PRIMARY KEY (id);


--
-- Name: re_zip3s re_zip3s_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_zip3s
    ADD CONSTRAINT re_zip3s_unique_key UNIQUE (contract_id, zip3);


--
-- Name: re_zip5_rate_areas re_zip5_rate_areas_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_zip5_rate_areas
    ADD CONSTRAINT re_zip5_rate_areas_pkey PRIMARY KEY (id);


--
-- Name: re_zip5_rate_areas re_zip5_rate_areas_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_zip5_rate_areas
    ADD CONSTRAINT re_zip5_rate_areas_unique_key UNIQUE (contract_id, zip5);


--
-- Name: archived_reimbursements reimbursements_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_reimbursements
    ADD CONSTRAINT reimbursements_pkey PRIMARY KEY (id);


--
-- Name: report_violations report_violations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.report_violations
    ADD CONSTRAINT report_violations_pkey PRIMARY KEY (id);


--
-- Name: report_violations report_violations_report_id_violation_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.report_violations
    ADD CONSTRAINT report_violations_report_id_violation_id_key UNIQUE (report_id, violation_id);


--
-- Name: reweighs reweighs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reweighs
    ADD CONSTRAINT reweighs_pkey PRIMARY KEY (id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: schema_migration schema_migration_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migration
    ADD CONSTRAINT schema_migration_pkey PRIMARY KEY (version);


--
-- Name: service_item_param_keys service_item_param_keys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_item_param_keys
    ADD CONSTRAINT service_item_param_keys_pkey PRIMARY KEY (id);


--
-- Name: service_item_param_keys service_item_param_keys_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_item_param_keys
    ADD CONSTRAINT service_item_param_keys_unique_key UNIQUE (key);


--
-- Name: service_items_customer_contacts service_items_customer_contacts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_items_customer_contacts
    ADD CONSTRAINT service_items_customer_contacts_pkey PRIMARY KEY (id);


--
-- Name: service_members service_members_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_members
    ADD CONSTRAINT service_members_pkey PRIMARY KEY (id);


--
-- Name: service_params service_params_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_params
    ADD CONSTRAINT service_params_pkey PRIMARY KEY (id);


--
-- Name: service_params service_params_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_params
    ADD CONSTRAINT service_params_unique_key UNIQUE (service_id, service_item_param_key_id);


--
-- Name: service_request_document_uploads service_request_document_uploads_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_request_document_uploads
    ADD CONSTRAINT service_request_document_uploads_pkey PRIMARY KEY (id);


--
-- Name: service_request_documents service_request_documents_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_request_documents
    ADD CONSTRAINT service_request_documents_pkey PRIMARY KEY (id);


--
-- Name: service_request_document_uploads service_request_documents_uploads_unique_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_request_document_uploads
    ADD CONSTRAINT service_request_documents_uploads_unique_key UNIQUE (upload_id);


--
-- Name: shipment_address_updates shipment_address_updates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.shipment_address_updates
    ADD CONSTRAINT shipment_address_updates_pkey PRIMARY KEY (id);


--
-- Name: signed_certifications signed_certifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.signed_certifications
    ADD CONSTRAINT signed_certifications_pkey PRIMARY KEY (id);


--
-- Name: sit_address_updates sit_address_updates_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sit_address_updates
    ADD CONSTRAINT sit_address_updates_pkey PRIMARY KEY (id);


--
-- Name: sit_extensions sit_extensions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sit_extensions
    ADD CONSTRAINT sit_extensions_pkey PRIMARY KEY (id);


--
-- Name: storage_facilities storage_facilities_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.storage_facilities
    ADD CONSTRAINT storage_facilities_pkey PRIMARY KEY (id);


--
-- Name: transportation_accounting_codes transportation_accounting_codes_tac_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transportation_accounting_codes
    ADD CONSTRAINT transportation_accounting_codes_tac_key UNIQUE (tac);


--
-- Name: transportation_offices transportation_offices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transportation_offices
    ADD CONSTRAINT transportation_offices_pkey PRIMARY KEY (id);


--
-- Name: contractors unique_contractors_type; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.contractors
    ADD CONSTRAINT unique_contractors_type UNIQUE (type);


--
-- Name: invoices unique_invoice_number; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT unique_invoice_number UNIQUE (invoice_number);


--
-- Name: duty_locations unique_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.duty_locations
    ADD CONSTRAINT unique_name UNIQUE (name);


--
-- Name: roles unique_roles; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT unique_roles UNIQUE (role_type);


--
-- Name: uploads uploads_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.uploads
    ADD CONSTRAINT uploads_pkey PRIMARY KEY (id);


--
-- Name: user_uploads user_uploads_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_uploads
    ADD CONSTRAINT user_uploads_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users_roles users_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_roles
    ADD CONSTRAINT users_roles_pkey PRIMARY KEY (id);


--
-- Name: webhook_notifications webhook_notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.webhook_notifications
    ADD CONSTRAINT webhook_notifications_pkey PRIMARY KEY (id);


--
-- Name: webhook_subscriptions webhook_subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.webhook_subscriptions
    ADD CONSTRAINT webhook_subscriptions_pkey PRIMARY KEY (id);


--
-- Name: weight_tickets weight_tickets_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weight_tickets
    ADD CONSTRAINT weight_tickets_pkey PRIMARY KEY (id);


--
-- Name: zip3_distances zip3_distances_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.zip3_distances
    ADD CONSTRAINT zip3_distances_pkey PRIMARY KEY (id);


--
-- Name: admin_users_active_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX admin_users_active_idx ON public.admin_users USING btree (active);


--
-- Name: admin_users_email_uniq_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX admin_users_email_uniq_idx ON public.admin_users USING btree (email);


--
-- Name: admin_users_organization_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX admin_users_organization_id_idx ON public.admin_users USING btree (organization_id);


--
-- Name: admin_users_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX admin_users_user_id_idx ON public.admin_users USING btree (user_id);


--
-- Name: archived_move_documents_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_move_documents_deleted_at_idx ON public.archived_move_documents USING btree (deleted_at);


--
-- Name: archived_move_documents_document_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_move_documents_document_id_idx ON public.archived_move_documents USING btree (document_id);


--
-- Name: archived_move_documents_move_id_document_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX archived_move_documents_move_id_document_id_idx ON public.archived_move_documents USING btree (move_id, document_id);


--
-- Name: archived_move_documents_personally_procured_move_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_move_documents_personally_procured_move_id_idx ON public.archived_move_documents USING btree (personally_procured_move_id);


--
-- Name: archived_moving_expense_documents_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_moving_expense_documents_deleted_at_idx ON public.archived_moving_expense_documents USING btree (deleted_at);


--
-- Name: archived_moving_expense_documents_move_document_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_moving_expense_documents_move_document_id_idx ON public.archived_moving_expense_documents USING btree (move_document_id);


--
-- Name: archived_personally_procured_moves_advance_worksheet_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_personally_procured_moves_advance_worksheet_id_idx ON public.archived_personally_procured_moves USING btree (advance_worksheet_id);


--
-- Name: archived_personally_procured_moves_move_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_personally_procured_moves_move_id_idx ON public.archived_personally_procured_moves USING btree (move_id);


--
-- Name: archived_personally_procured_moves_original_move_date_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_personally_procured_moves_original_move_date_idx ON public.archived_personally_procured_moves USING btree (original_move_date);


--
-- Name: archived_personally_procured_moves_reviewed_date_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_personally_procured_moves_reviewed_date_idx ON public.archived_personally_procured_moves USING btree (reviewed_date);


--
-- Name: archived_signed_certifications_move_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_signed_certifications_move_id_idx ON public.archived_signed_certifications USING btree (move_id);


--
-- Name: archived_signed_certifications_submitting_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_signed_certifications_submitting_user_id_idx ON public.archived_signed_certifications USING btree (submitting_user_id);


--
-- Name: archived_weight_ticket_set_documents_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX archived_weight_ticket_set_documents_deleted_at_idx ON public.archived_weight_ticket_set_documents USING btree (deleted_at);


--
-- Name: audit_history_action_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_history_action_idx ON public.audit_history USING btree (action);


--
-- Name: audit_history_action_tstamp_tx_stm_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_history_action_tstamp_tx_stm_idx ON public.audit_history USING btree (action_tstamp_stm);


--
-- Name: audit_history_relid_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_history_relid_idx ON public.audit_history USING btree (relid);


--
-- Name: audit_history_table_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX audit_history_table_name_idx ON public.audit_history USING btree (table_name);


--
-- Name: backup_contacts_email_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX backup_contacts_email_idx ON public.backup_contacts USING btree (email);


--
-- Name: customer_support_remarks_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX customer_support_remarks_deleted_at_idx ON public.customer_support_remarks USING btree (deleted_at);


--
-- Name: distance_calculations_destination_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX distance_calculations_destination_address_id_idx ON public.distance_calculations USING btree (destination_address_id);


--
-- Name: distance_calculations_origin_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX distance_calculations_origin_address_id_idx ON public.distance_calculations USING btree (origin_address_id);


--
-- Name: documents_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX documents_deleted_at_idx ON public.documents USING btree (deleted_at);


--
-- Name: documents_service_member_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX documents_service_member_id_idx ON public.documents USING btree (service_member_id);


--
-- Name: duty_location_names_duty_location_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX duty_location_names_duty_location_id_idx ON public.duty_location_names USING btree (duty_location_id);


--
-- Name: duty_location_names_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX duty_location_names_name_idx ON public.duty_location_names USING btree (name);


--
-- Name: duty_location_names_name_trgm_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX duty_location_names_name_trgm_idx ON public.duty_location_names USING gin (name public.gin_trgm_ops);


--
-- Name: duty_locations_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX duty_locations_address_id_idx ON public.duty_locations USING btree (address_id);


--
-- Name: duty_locations_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX duty_locations_name_idx ON public.duty_locations USING btree (name);


--
-- Name: duty_locations_name_trgm_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX duty_locations_name_trgm_idx ON public.duty_locations USING gin (name public.gin_trgm_ops);


--
-- Name: duty_locations_transportation_office_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX duty_locations_transportation_office_id_idx ON public.duty_locations USING btree (transportation_office_id);


--
-- Name: edi_errors_interchange_control_number_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX edi_errors_interchange_control_number_id_idx ON public.edi_errors USING btree (interchange_control_number_id);


--
-- Name: edi_errors_payment_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX edi_errors_payment_request_id_idx ON public.edi_errors USING btree (payment_request_id);


--
-- Name: electronic_orders_index_by_issuer_and_edipi; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX electronic_orders_index_by_issuer_and_edipi ON public.electronic_orders USING btree (issuer, edipi);


--
-- Name: electronic_orders_index_by_issuer_and_orders_number; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX electronic_orders_index_by_issuer_and_orders_number ON public.electronic_orders USING btree (issuer, orders_number);


--
-- Name: electronic_orders_revisions_electronic_order_id_seq_num_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX electronic_orders_revisions_electronic_order_id_seq_num_idx ON public.electronic_orders_revisions USING btree (electronic_order_id, seq_num);


--
-- Name: evaluation_reports_move_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX evaluation_reports_move_id_idx ON public.evaluation_reports USING btree (move_id);


--
-- Name: evaluation_reports_office_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX evaluation_reports_office_user_id_idx ON public.evaluation_reports USING btree (office_user_id);


--
-- Name: evaluation_reports_shipment_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX evaluation_reports_shipment_id_idx ON public.evaluation_reports USING btree (shipment_id);


--
-- Name: evaluation_reports_submitted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX evaluation_reports_submitted_at_idx ON public.evaluation_reports USING btree (submitted_at);


--
-- Name: invoices_approver_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX invoices_approver_id_idx ON public.invoices USING btree (approver_id);


--
-- Name: jppso_regions_code_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX jppso_regions_code_idx ON public.jppso_regions USING btree (code);


--
-- Name: jppso_state_assignments_state_abbreviation_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX jppso_state_assignments_state_abbreviation_idx ON public.jppso_region_state_assignments USING btree (state_abbreviation);


--
-- Name: jppso_state_assignments_state_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX jppso_state_assignments_state_name_idx ON public.jppso_region_state_assignments USING btree (state_name);


--
-- Name: moves_available_to_prime_and_show_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX moves_available_to_prime_and_show_idx ON public.moves USING btree (show, available_to_prime_at);


--
-- Name: moves_locator_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX moves_locator_idx ON public.moves USING btree (locator);


--
-- Name: moves_orders_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX moves_orders_id_idx ON public.moves USING btree (orders_id);


--
-- Name: moves_status_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX moves_status_idx ON public.moves USING btree (status);


--
-- Name: moves_submitted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX moves_submitted_at_idx ON public.moves USING btree (submitted_at);


--
-- Name: moving_expenses_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX moving_expenses_deleted_at_idx ON public.moving_expenses USING btree (deleted_at);


--
-- Name: moving_expenses_ppm_shipment_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX moving_expenses_ppm_shipment_id_idx ON public.moving_expenses USING hash (ppm_shipment_id);


--
-- Name: mto_service_items_move_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_service_items_move_id_idx ON public.mto_service_items USING btree (move_id);


--
-- Name: mto_service_items_mto_shipment_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_service_items_mto_shipment_id_idx ON public.mto_service_items USING btree (mto_shipment_id);


--
-- Name: mto_service_items_re_service_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_service_items_re_service_id_idx ON public.mto_service_items USING btree (re_service_id);


--
-- Name: mto_service_items_sit_origin_hhg_actual_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_service_items_sit_origin_hhg_actual_address_id_idx ON public.mto_service_items USING btree (sit_origin_hhg_actual_address_id);


--
-- Name: mto_service_items_sit_origin_hhg_original_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_service_items_sit_origin_hhg_original_address_id_idx ON public.mto_service_items USING btree (sit_origin_hhg_original_address_id);


--
-- Name: mto_shipments_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_shipments_deleted_at_idx ON public.mto_shipments USING btree (deleted_at);


--
-- Name: mto_shipments_destination_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_shipments_destination_address_id_idx ON public.mto_shipments USING btree (destination_address_id);


--
-- Name: mto_shipments_move_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_shipments_move_id_idx ON public.mto_shipments USING btree (move_id);


--
-- Name: mto_shipments_pickup_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_shipments_pickup_address_id_idx ON public.mto_shipments USING btree (pickup_address_id);


--
-- Name: mto_shipments_requested_pickup_date_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_shipments_requested_pickup_date_idx ON public.mto_shipments USING btree (requested_pickup_date);


--
-- Name: mto_shipments_secondary_delivery_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_shipments_secondary_delivery_address_id_idx ON public.mto_shipments USING btree (secondary_delivery_address_id);


--
-- Name: mto_shipments_secondary_pickup_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_shipments_secondary_pickup_address_id_idx ON public.mto_shipments USING btree (secondary_pickup_address_id);


--
-- Name: mto_shipments_storage_facility_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX mto_shipments_storage_facility_id_idx ON public.mto_shipments USING btree (storage_facility_id);


--
-- Name: notifications_notification_type_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX notifications_notification_type_idx ON public.notifications USING btree (notification_type);


--
-- Name: notifications_service_member_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX notifications_service_member_id_idx ON public.notifications USING btree (service_member_id);


--
-- Name: office_emails_transportation_office_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX office_emails_transportation_office_id_idx ON public.office_emails USING btree (transportation_office_id);


--
-- Name: office_phone_lines_is_dsn_number_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX office_phone_lines_is_dsn_number_idx ON public.office_phone_lines USING btree (is_dsn_number);


--
-- Name: office_phone_lines_transportation_office_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX office_phone_lines_transportation_office_id_idx ON public.office_phone_lines USING btree (transportation_office_id);


--
-- Name: office_users_active_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX office_users_active_idx ON public.office_users USING btree (active);


--
-- Name: office_users_email_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX office_users_email_idx ON public.office_users USING btree (email);


--
-- Name: office_users_email_trgrm_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX office_users_email_trgrm_idx ON public.office_users USING gin (email public.gin_trgm_ops);


--
-- Name: office_users_first_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX office_users_first_name_idx ON public.office_users USING btree (first_name);


--
-- Name: office_users_last_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX office_users_last_name_idx ON public.office_users USING btree (last_name);


--
-- Name: office_users_transportation_office_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX office_users_transportation_office_id_idx ON public.office_users USING btree (transportation_office_id);


--
-- Name: office_users_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX office_users_user_id_idx ON public.office_users USING btree (user_id);


--
-- Name: orders_entitlement_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX orders_entitlement_id_idx ON public.orders USING btree (entitlement_id);


--
-- Name: orders_gbloc_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX orders_gbloc_idx ON public.orders USING btree (gbloc);


--
-- Name: orders_new_duty_location_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX orders_new_duty_location_id_idx ON public.orders USING btree (new_duty_location_id);


--
-- Name: orders_origin_duty_location_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX orders_origin_duty_location_id_idx ON public.orders USING btree (origin_duty_location_id);


--
-- Name: orders_service_member_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX orders_service_member_id_idx ON public.orders USING btree (service_member_id);


--
-- Name: orders_uploaded_amended_orders_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX orders_uploaded_amended_orders_id_idx ON public.orders USING btree (uploaded_amended_orders_id);


--
-- Name: orders_uploaded_orders_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX orders_uploaded_orders_id_idx ON public.orders USING btree (uploaded_orders_id);


--
-- Name: payment_requests_created_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX payment_requests_created_at_idx ON public.payment_requests USING btree (created_at);


--
-- Name: payment_requests_move_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX payment_requests_move_id_idx ON public.payment_requests USING btree (move_id);


--
-- Name: payment_requests_recalculation_of_payment_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX payment_requests_recalculation_of_payment_request_id_idx ON public.payment_requests USING btree (recalculation_of_payment_request_id);


--
-- Name: payment_requests_status_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX payment_requests_status_idx ON public.payment_requests USING btree (status);


--
-- Name: payment_service_items_mto_service_item_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX payment_service_items_mto_service_item_id_idx ON public.payment_service_items USING btree (mto_service_item_id);


--
-- Name: payment_service_items_payment_request_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX payment_service_items_payment_request_id_idx ON public.payment_service_items USING btree (payment_request_id);


--
-- Name: personally_procured_moves_advance_worksheet_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX personally_procured_moves_advance_worksheet_id_idx ON public.personally_procured_moves USING btree (advance_worksheet_id);


--
-- Name: personally_procured_moves_move_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX personally_procured_moves_move_id_idx ON public.personally_procured_moves USING btree (move_id);


--
-- Name: personally_procured_moves_original_move_date_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX personally_procured_moves_original_move_date_idx ON public.personally_procured_moves USING btree (original_move_date);


--
-- Name: personally_procured_moves_reviewed_date_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX personally_procured_moves_reviewed_date_idx ON public.personally_procured_moves USING btree (reviewed_date);


--
-- Name: ppm_shipments_w2_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX ppm_shipments_w2_address_id_idx ON public.ppm_shipments USING btree (w2_address_id);


--
-- Name: prime_uploads_contractor_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX prime_uploads_contractor_id_idx ON public.prime_uploads USING btree (contractor_id);


--
-- Name: prime_uploads_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX prime_uploads_deleted_at_idx ON public.prime_uploads USING btree (deleted_at);


--
-- Name: prime_uploads_proof_of_service_docs_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX prime_uploads_proof_of_service_docs_id_idx ON public.prime_uploads USING btree (proof_of_service_docs_id);


--
-- Name: progear_weight_tickets_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX progear_weight_tickets_deleted_at_idx ON public.progear_weight_tickets USING btree (deleted_at);


--
-- Name: progear_weight_tickets_ppm_shipment_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX progear_weight_tickets_ppm_shipment_id_idx ON public.progear_weight_tickets USING hash (ppm_shipment_id);


--
-- Name: pws_violations_display_order_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX pws_violations_display_order_idx ON public.pws_violations USING btree (display_order);


--
-- Name: re_zip3s_rate_area_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX re_zip3s_rate_area_id_idx ON public.re_zip3s USING btree (rate_area_id);


--
-- Name: re_zip5_rate_areas_rate_area_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX re_zip5_rate_areas_rate_area_id_idx ON public.re_zip5_rate_areas USING btree (rate_area_id);


--
-- Name: report_violations_report_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX report_violations_report_id_idx ON public.report_violations USING btree (report_id);


--
-- Name: reweighs_shipment_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX reweighs_shipment_id_idx ON public.reweighs USING btree (shipment_id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: service_members_affiliation_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX service_members_affiliation_idx ON public.service_members USING btree (affiliation);


--
-- Name: service_members_backup_mailing_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX service_members_backup_mailing_address_id_idx ON public.service_members USING btree (backup_mailing_address_id);


--
-- Name: service_members_duty_location_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX service_members_duty_location_id_idx ON public.service_members USING btree (duty_location_id);


--
-- Name: service_members_edipi_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX service_members_edipi_idx ON public.service_members USING btree (edipi);


--
-- Name: service_members_last_name_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX service_members_last_name_idx ON public.service_members USING btree (last_name text_pattern_ops);


--
-- Name: service_members_personal_email_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX service_members_personal_email_idx ON public.service_members USING btree (personal_email);


--
-- Name: service_members_residential_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX service_members_residential_address_id_idx ON public.service_members USING btree (residential_address_id);


--
-- Name: service_members_searchable_full_name_trgm_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX service_members_searchable_full_name_trgm_idx ON public.service_members USING gin (public.searchable_full_name(first_name, last_name) public.gin_trgm_ops);


--
-- Name: service_members_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX service_members_user_id_idx ON public.service_members USING btree (user_id);


--
-- Name: service_members_user_id_uniq_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX service_members_user_id_uniq_idx ON public.service_members USING btree (user_id);


--
-- Name: shipment_address_updates_shipment_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX shipment_address_updates_shipment_id_idx ON public.shipment_address_updates USING btree (shipment_id);


--
-- Name: signed_certifications_move_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX signed_certifications_move_id_idx ON public.signed_certifications USING btree (move_id);


--
-- Name: signed_certifications_ppm_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX signed_certifications_ppm_id_idx ON public.signed_certifications USING btree (ppm_id);


--
-- Name: signed_certifications_submitting_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX signed_certifications_submitting_user_id_idx ON public.signed_certifications USING btree (submitting_user_id);


--
-- Name: sit_extensions_mto_shipment_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX sit_extensions_mto_shipment_id_idx ON public.sit_extensions USING btree (mto_shipment_id);


--
-- Name: storage_facilities_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX storage_facilities_address_id_idx ON public.storage_facilities USING btree (address_id);


--
-- Name: transportation_offices_address_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX transportation_offices_address_id_idx ON public.transportation_offices USING btree (address_id);


--
-- Name: transportation_offices_gbloc_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX transportation_offices_gbloc_idx ON public.transportation_offices USING btree (gbloc);


--
-- Name: transportation_offices_name_trgm_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX transportation_offices_name_trgm_idx ON public.transportation_offices USING gin (name public.gin_trgm_ops);


--
-- Name: uploads_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX uploads_deleted_at_idx ON public.uploads USING btree (deleted_at);


--
-- Name: user_uploads_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX user_uploads_deleted_at_idx ON public.user_uploads USING btree (deleted_at);


--
-- Name: user_uploads_document_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX user_uploads_document_id_idx ON public.user_uploads USING btree (document_id);


--
-- Name: user_uploads_uploader_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX user_uploads_uploader_id_idx ON public.user_uploads USING btree (uploader_id);


--
-- Name: users_active_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_active_idx ON public.users USING btree (active);


--
-- Name: users_current_admin_session_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_current_admin_session_id_idx ON public.users USING btree (current_admin_session_id);


--
-- Name: users_current_mil_session_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_current_mil_session_id_idx ON public.users USING btree (current_mil_session_id);


--
-- Name: users_current_office_session_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_current_office_session_id_idx ON public.users USING btree (current_office_session_id);


--
-- Name: users_login_gov_email_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_login_gov_email_idx ON public.users USING btree (login_gov_email);


--
-- Name: users_roles_not_deleted_partial_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_roles_not_deleted_partial_idx ON public.users_roles USING btree (deleted_at) WHERE (deleted_at IS NULL);


--
-- Name: INDEX users_roles_not_deleted_partial_idx; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON INDEX public.users_roles_not_deleted_partial_idx IS 'indexes users_roles that are not deleted';


--
-- Name: users_roles_user_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX users_roles_user_id_idx ON public.users_roles USING btree (user_id);


--
-- Name: webhook_notifications_move_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX webhook_notifications_move_id_idx ON public.webhook_notifications USING btree (move_id);


--
-- Name: webhook_notifications_unsent; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX webhook_notifications_unsent ON public.webhook_notifications USING btree (created_at) WHERE ((status <> 'SENT'::public.webhook_notifications_status) AND (status <> 'SKIPPED'::public.webhook_notifications_status));


--
-- Name: weight_tickets_deleted_at_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX weight_tickets_deleted_at_idx ON public.weight_tickets USING btree (deleted_at);


--
-- Name: weight_tickets_ppm_shipment_id_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX weight_tickets_ppm_shipment_id_idx ON public.weight_tickets USING hash (ppm_shipment_id);


--
-- Name: zip3_distances_unique_zip3s; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX zip3_distances_unique_zip3s ON public.zip3_distances USING btree (from_zip3, to_zip3);


--
-- Name: INDEX zip3_distances_unique_zip3s; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON INDEX public.zip3_distances_unique_zip3s IS 'Ensures that we have only one entry for a from/to zip3 pair.';


--
-- Name: addresses audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.addresses FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at}');


--
-- Name: backup_contacts audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.backup_contacts FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,service_member_id}');


--
-- Name: entitlements audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.entitlements FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at}');


--
-- Name: moves audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.moves FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,orders_id,contractor_id,excess_weight_upload_id,selected_move_type}');


--
-- Name: mto_agents audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.mto_agents FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,mto_shipment_id}');


--
-- Name: mto_service_item_customer_contacts audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.mto_service_item_customer_contacts FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,mto_service_item_id}');


--
-- Name: mto_service_item_dimensions audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.mto_service_item_dimensions FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,mto_service_item_id}');


--
-- Name: mto_service_items audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.mto_service_items FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,move_id,mto_shipment_id,re_service_id,sit_destination_final_address_id,sit_origin_hhg_original_address_id,sit_origin_hhg_actual_address_id}');


--
-- Name: mto_shipments audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.mto_shipments FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,destination_address_id,secondary_pickup_address_id,secondary_delivery_address_id,pickup_address_id,move_id,storage_facility_id}');


--
-- Name: orders audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.orders FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,service_member_id,uploaded_orders_id,entitlement_id,uploaded_amended_orders_id,grade}');


--
-- Name: payment_requests audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.payment_requests FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,move_id}');


--
-- Name: proof_of_service_docs audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.proof_of_service_docs FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,payment_request_id}');


--
-- Name: reweighs audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.reweighs FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,shipment_id}');


--
-- Name: service_members audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.service_members FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,user_id,residential_address_id,backup_mailing_address_id}');


--
-- Name: storage_facilities audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.storage_facilities FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,address_id}');


--
-- Name: user_uploads audit_trigger_row; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_row AFTER INSERT OR DELETE OR UPDATE ON public.user_uploads FOR EACH ROW EXECUTE FUNCTION public.if_modified_func('true', '{created_at,updated_at,document_id,uploader_id,upload_id}');


--
-- Name: addresses audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.addresses FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: backup_contacts audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.backup_contacts FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: entitlements audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.entitlements FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: moves audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.moves FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: mto_agents audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.mto_agents FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: mto_service_item_customer_contacts audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.mto_service_item_customer_contacts FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: mto_service_item_dimensions audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.mto_service_item_dimensions FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: mto_service_items audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.mto_service_items FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: mto_shipments audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.mto_shipments FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: orders audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.orders FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: payment_requests audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.payment_requests FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: proof_of_service_docs audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.proof_of_service_docs FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: reweighs audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.reweighs FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: service_members audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.service_members FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: storage_facilities audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.storage_facilities FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: user_uploads audit_trigger_stm; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER audit_trigger_stm AFTER TRUNCATE ON public.user_uploads FOR EACH STATEMENT EXECUTE FUNCTION public.if_modified_func('true');


--
-- Name: admin_users admin_users_organizations_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_organizations_id_fk FOREIGN KEY (organization_id) REFERENCES public.organizations(id);


--
-- Name: admin_users admin_users_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_users_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: archived_access_codes archived_access_codes_service_member_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_access_codes
    ADD CONSTRAINT archived_access_codes_service_member_id FOREIGN KEY (service_member_id) REFERENCES public.service_members(id);


--
-- Name: archived_move_documents archived_move_documents_document_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_move_documents
    ADD CONSTRAINT archived_move_documents_document_id FOREIGN KEY (document_id) REFERENCES public.documents(id);


--
-- Name: archived_move_documents archived_move_documents_move_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_move_documents
    ADD CONSTRAINT archived_move_documents_move_id FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: archived_move_documents archived_move_documents_personally_procured_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_move_documents
    ADD CONSTRAINT archived_move_documents_personally_procured_move_id_fkey FOREIGN KEY (personally_procured_move_id) REFERENCES public.archived_personally_procured_moves(id);


--
-- Name: archived_moving_expense_documents archived_moving_expense_documents_move_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_moving_expense_documents
    ADD CONSTRAINT archived_moving_expense_documents_move_document_id_fkey FOREIGN KEY (move_document_id) REFERENCES public.archived_move_documents(id);


--
-- Name: archived_personally_procured_moves archived_personally_procured_moves_advance_worksheet_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_personally_procured_moves
    ADD CONSTRAINT archived_personally_procured_moves_advance_worksheet_id FOREIGN KEY (advance_worksheet_id) REFERENCES public.documents(id);


--
-- Name: archived_personally_procured_moves archived_personally_procured_moves_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_personally_procured_moves
    ADD CONSTRAINT archived_personally_procured_moves_move_id_fkey FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: archived_signed_certifications archived_signed_certifications_move_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_signed_certifications
    ADD CONSTRAINT archived_signed_certifications_move_id FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: archived_signed_certifications archived_signed_certifications_personally_procured_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_signed_certifications
    ADD CONSTRAINT archived_signed_certifications_personally_procured_move_id_fkey FOREIGN KEY (personally_procured_move_id) REFERENCES public.archived_personally_procured_moves(id);


--
-- Name: archived_signed_certifications archived_signed_certifications_submitting_user_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_signed_certifications
    ADD CONSTRAINT archived_signed_certifications_submitting_user_id FOREIGN KEY (submitting_user_id) REFERENCES public.users(id);


--
-- Name: archived_weight_ticket_set_documents archived_weight_ticket_set_documents_move_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.archived_weight_ticket_set_documents
    ADD CONSTRAINT archived_weight_ticket_set_documents_move_document_id_fkey FOREIGN KEY (move_document_id) REFERENCES public.archived_move_documents(id);


--
-- Name: client_certs client_certs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.client_certs
    ADD CONSTRAINT client_certs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: distance_calculations distance_calculations_destination_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.distance_calculations
    ADD CONSTRAINT distance_calculations_destination_address_id_fkey FOREIGN KEY (destination_address_id) REFERENCES public.addresses(id);


--
-- Name: distance_calculations distance_calculations_origin_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.distance_calculations
    ADD CONSTRAINT distance_calculations_origin_address_id_fkey FOREIGN KEY (origin_address_id) REFERENCES public.addresses(id);


--
-- Name: documents documents_service_members_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_service_members_id_fk FOREIGN KEY (service_member_id) REFERENCES public.service_members(id);


--
-- Name: duty_location_names duty_location_names_duty_location_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.duty_location_names
    ADD CONSTRAINT duty_location_names_duty_location_id_fkey FOREIGN KEY (duty_location_id) REFERENCES public.duty_locations(id);


--
-- Name: duty_locations duty_locations_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.duty_locations
    ADD CONSTRAINT duty_locations_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id);


--
-- Name: duty_locations duty_locations_transportation_offices_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.duty_locations
    ADD CONSTRAINT duty_locations_transportation_offices_id_fk FOREIGN KEY (transportation_office_id) REFERENCES public.transportation_offices(id);


--
-- Name: edi_errors edi_errors_icn_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.edi_errors
    ADD CONSTRAINT edi_errors_icn_id_fkey FOREIGN KEY (interchange_control_number_id) REFERENCES public.payment_request_to_interchange_control_numbers(id);


--
-- Name: edi_errors edi_errors_payment_request_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.edi_errors
    ADD CONSTRAINT edi_errors_payment_request_id_fkey FOREIGN KEY (payment_request_id) REFERENCES public.payment_requests(id);


--
-- Name: electronic_orders_revisions electronic_orders_revisions_electronic_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.electronic_orders_revisions
    ADD CONSTRAINT electronic_orders_revisions_electronic_order_id_fkey FOREIGN KEY (electronic_order_id) REFERENCES public.electronic_orders(id) ON DELETE CASCADE;


--
-- Name: evaluation_reports evaluation_reports_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.evaluation_reports
    ADD CONSTRAINT evaluation_reports_move_id_fkey FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: evaluation_reports evaluation_reports_office_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.evaluation_reports
    ADD CONSTRAINT evaluation_reports_office_user_id_fkey FOREIGN KEY (office_user_id) REFERENCES public.office_users(id);


--
-- Name: evaluation_reports evaluation_reports_shipment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.evaluation_reports
    ADD CONSTRAINT evaluation_reports_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.mto_shipments(id);


--
-- Name: customer_support_remarks fk_moves; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customer_support_remarks
    ADD CONSTRAINT fk_moves FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: customer_support_remarks fk_office_users; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.customer_support_remarks
    ADD CONSTRAINT fk_office_users FOREIGN KEY (office_user_id) REFERENCES public.office_users(id);


--
-- Name: invoices invoices_office_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_office_users_id_fk FOREIGN KEY (approver_id) REFERENCES public.office_users(id);


--
-- Name: invoices invoices_user_uploads_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.invoices
    ADD CONSTRAINT invoices_user_uploads_id_fkey FOREIGN KEY (user_uploads_id) REFERENCES public.user_uploads(id) ON DELETE RESTRICT;


--
-- Name: jppso_region_state_assignments jppso_region_state_assignments_jppso_region_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.jppso_region_state_assignments
    ADD CONSTRAINT jppso_region_state_assignments_jppso_region_id_fkey FOREIGN KEY (jppso_region_id) REFERENCES public.jppso_regions(id);


--
-- Name: moves moves_closeout_office_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.moves
    ADD CONSTRAINT moves_closeout_office_id_fkey FOREIGN KEY (closeout_office_id) REFERENCES public.transportation_offices(id);


--
-- Name: moves moves_contractor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.moves
    ADD CONSTRAINT moves_contractor_id_fkey FOREIGN KEY (contractor_id) REFERENCES public.contractors(id);


--
-- Name: moves moves_excess_weight_upload_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.moves
    ADD CONSTRAINT moves_excess_weight_upload_id_fkey FOREIGN KEY (excess_weight_upload_id) REFERENCES public.uploads(id);


--
-- Name: moves moves_orders_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.moves
    ADD CONSTRAINT moves_orders_id_fk FOREIGN KEY (orders_id) REFERENCES public.orders(id);


--
-- Name: moving_expenses moving_expenses_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.moving_expenses
    ADD CONSTRAINT moving_expenses_document_id_fkey FOREIGN KEY (document_id) REFERENCES public.documents(id);


--
-- Name: moving_expenses moving_expenses_ppm_shipments_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.moving_expenses
    ADD CONSTRAINT moving_expenses_ppm_shipments_id_fkey FOREIGN KEY (ppm_shipment_id) REFERENCES public.ppm_shipments(id);


--
-- Name: mto_agents mto_agents_mto_shipment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_agents
    ADD CONSTRAINT mto_agents_mto_shipment_id_fkey FOREIGN KEY (mto_shipment_id) REFERENCES public.mto_shipments(id);


--
-- Name: mto_service_item_dimensions mto_service_item_dimensions_mto_service_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_item_dimensions
    ADD CONSTRAINT mto_service_item_dimensions_mto_service_item_id_fkey FOREIGN KEY (mto_service_item_id) REFERENCES public.mto_service_items(id) ON DELETE CASCADE;


--
-- Name: mto_service_items mto_service_items_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_items
    ADD CONSTRAINT mto_service_items_move_id_fkey FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: mto_service_items mto_service_items_mto_shipment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_items
    ADD CONSTRAINT mto_service_items_mto_shipment_id_fkey FOREIGN KEY (mto_shipment_id) REFERENCES public.mto_shipments(id);


--
-- Name: mto_service_items mto_service_items_re_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_items
    ADD CONSTRAINT mto_service_items_re_service_id_fkey FOREIGN KEY (re_service_id) REFERENCES public.re_services(id);


--
-- Name: mto_service_items mto_service_items_sit_destination_final_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_items
    ADD CONSTRAINT mto_service_items_sit_destination_final_address_id_fkey FOREIGN KEY (sit_destination_final_address_id) REFERENCES public.addresses(id);


--
-- Name: mto_service_items mto_service_items_sit_destination_original_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_items
    ADD CONSTRAINT mto_service_items_sit_destination_original_address_id_fkey FOREIGN KEY (sit_destination_original_address_id) REFERENCES public.addresses(id);


--
-- Name: mto_service_items mto_service_items_sit_origin_hhg_actual_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_items
    ADD CONSTRAINT mto_service_items_sit_origin_hhg_actual_address_id_fkey FOREIGN KEY (sit_origin_hhg_actual_address_id) REFERENCES public.addresses(id);


--
-- Name: mto_service_items mto_service_items_sit_origin_hhg_original_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_service_items
    ADD CONSTRAINT mto_service_items_sit_origin_hhg_original_address_id_fkey FOREIGN KEY (sit_origin_hhg_original_address_id) REFERENCES public.addresses(id);


--
-- Name: mto_shipments mto_shipments_destination_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_shipments
    ADD CONSTRAINT mto_shipments_destination_address_id_fkey FOREIGN KEY (destination_address_id) REFERENCES public.addresses(id);


--
-- Name: mto_shipments mto_shipments_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_shipments
    ADD CONSTRAINT mto_shipments_move_id_fkey FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: mto_shipments mto_shipments_pickup_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_shipments
    ADD CONSTRAINT mto_shipments_pickup_address_id_fkey FOREIGN KEY (pickup_address_id) REFERENCES public.addresses(id);


--
-- Name: mto_shipments mto_shipments_secondary_delivery_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_shipments
    ADD CONSTRAINT mto_shipments_secondary_delivery_address_id_fkey FOREIGN KEY (secondary_delivery_address_id) REFERENCES public.addresses(id);


--
-- Name: mto_shipments mto_shipments_secondary_pickup_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_shipments
    ADD CONSTRAINT mto_shipments_secondary_pickup_address_id_fkey FOREIGN KEY (secondary_pickup_address_id) REFERENCES public.addresses(id);


--
-- Name: mto_shipments mto_shipments_storage_facility_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.mto_shipments
    ADD CONSTRAINT mto_shipments_storage_facility_id_fkey FOREIGN KEY (storage_facility_id) REFERENCES public.storage_facilities(id);


--
-- Name: office_emails office_emails_transportation_office_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.office_emails
    ADD CONSTRAINT office_emails_transportation_office_id_fkey FOREIGN KEY (transportation_office_id) REFERENCES public.transportation_offices(id);


--
-- Name: office_phone_lines office_phone_lines_transportation_office_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.office_phone_lines
    ADD CONSTRAINT office_phone_lines_transportation_office_id_fkey FOREIGN KEY (transportation_office_id) REFERENCES public.transportation_offices(id);


--
-- Name: office_users office_users_transportation_office_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.office_users
    ADD CONSTRAINT office_users_transportation_office_id_fkey FOREIGN KEY (transportation_office_id) REFERENCES public.transportation_offices(id);


--
-- Name: office_users office_users_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.office_users
    ADD CONSTRAINT office_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: orders orders_documents_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_documents_id_fk FOREIGN KEY (uploaded_orders_id) REFERENCES public.documents(id);


--
-- Name: orders orders_entitlement_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_entitlement_id_fkey FOREIGN KEY (entitlement_id) REFERENCES public.entitlements(id);


--
-- Name: orders orders_new_duty_location_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_new_duty_location_id_fkey FOREIGN KEY (new_duty_location_id) REFERENCES public.duty_locations(id);


--
-- Name: orders orders_origin_duty_location_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_origin_duty_location_id_fkey FOREIGN KEY (origin_duty_location_id) REFERENCES public.duty_locations(id);


--
-- Name: orders orders_service_member_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_service_member_id_fkey FOREIGN KEY (service_member_id) REFERENCES public.service_members(id) ON DELETE CASCADE;


--
-- Name: orders orders_uploaded_amended_orders_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_uploaded_amended_orders_id_fkey FOREIGN KEY (uploaded_amended_orders_id) REFERENCES public.documents(id);


--
-- Name: payment_request_to_interchange_control_numbers payment_request_to_icns_payment_request_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_request_to_interchange_control_numbers
    ADD CONSTRAINT payment_request_to_icns_payment_request_id_fkey FOREIGN KEY (payment_request_id) REFERENCES public.payment_requests(id);


--
-- Name: payment_requests payment_requests_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_requests
    ADD CONSTRAINT payment_requests_move_id_fkey FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: payment_requests payment_requests_recalculation_of_payment_request_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_requests
    ADD CONSTRAINT payment_requests_recalculation_of_payment_request_id_fkey FOREIGN KEY (recalculation_of_payment_request_id) REFERENCES public.payment_requests(id);


--
-- Name: payment_service_item_params payment_service_item_params_payment_service_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_service_item_params
    ADD CONSTRAINT payment_service_item_params_payment_service_item_id_fkey FOREIGN KEY (payment_service_item_id) REFERENCES public.payment_service_items(id);


--
-- Name: payment_service_item_params payment_service_item_params_service_item_param_key_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_service_item_params
    ADD CONSTRAINT payment_service_item_params_service_item_param_key_id_fkey FOREIGN KEY (service_item_param_key_id) REFERENCES public.service_item_param_keys(id);


--
-- Name: payment_service_items payment_service_items_mto_service_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_service_items
    ADD CONSTRAINT payment_service_items_mto_service_item_id_fkey FOREIGN KEY (mto_service_item_id) REFERENCES public.mto_service_items(id);


--
-- Name: payment_service_items payment_service_items_payment_request_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payment_service_items
    ADD CONSTRAINT payment_service_items_payment_request_id_fkey FOREIGN KEY (payment_request_id) REFERENCES public.payment_requests(id);


--
-- Name: personally_procured_moves personally_procured_moves_documents_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.personally_procured_moves
    ADD CONSTRAINT personally_procured_moves_documents_id_fk FOREIGN KEY (advance_worksheet_id) REFERENCES public.documents(id);


--
-- Name: personally_procured_moves personally_procured_moves_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.personally_procured_moves
    ADD CONSTRAINT personally_procured_moves_move_id_fkey FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: signed_certifications ppm_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.signed_certifications
    ADD CONSTRAINT ppm_id_fkey FOREIGN KEY (ppm_id) REFERENCES public.ppm_shipments(id);


--
-- Name: ppm_shipments ppm_shipment_mto_shipment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ppm_shipments
    ADD CONSTRAINT ppm_shipment_mto_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.mto_shipments(id);


--
-- Name: ppm_shipments ppm_shipments_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ppm_shipments
    ADD CONSTRAINT ppm_shipments_address_id_fkey FOREIGN KEY (w2_address_id) REFERENCES public.addresses(id);


--
-- Name: ppm_shipments ppm_shipments_aoa_packet_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ppm_shipments
    ADD CONSTRAINT ppm_shipments_aoa_packet_id_fkey FOREIGN KEY (aoa_packet_id) REFERENCES public.documents(id);


--
-- Name: ppm_shipments ppm_shipments_payment_packet_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ppm_shipments
    ADD CONSTRAINT ppm_shipments_payment_packet_id_fkey FOREIGN KEY (payment_packet_id) REFERENCES public.documents(id);


--
-- Name: prime_uploads prime_uploads_contractor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.prime_uploads
    ADD CONSTRAINT prime_uploads_contractor_id_fkey FOREIGN KEY (contractor_id) REFERENCES public.contractors(id);


--
-- Name: prime_uploads prime_uploads_proof_of_service_docs_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.prime_uploads
    ADD CONSTRAINT prime_uploads_proof_of_service_docs_id_fkey FOREIGN KEY (proof_of_service_docs_id) REFERENCES public.proof_of_service_docs(id);


--
-- Name: prime_uploads prime_uploads_uploads_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.prime_uploads
    ADD CONSTRAINT prime_uploads_uploads_id_fkey FOREIGN KEY (upload_id) REFERENCES public.uploads(id) ON DELETE RESTRICT;


--
-- Name: progear_weight_tickets progear_weight_tickets_full_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.progear_weight_tickets
    ADD CONSTRAINT progear_weight_tickets_full_document_id_fkey FOREIGN KEY (document_id) REFERENCES public.documents(id);


--
-- Name: progear_weight_tickets progear_weight_tickets_ppm_shipment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.progear_weight_tickets
    ADD CONSTRAINT progear_weight_tickets_ppm_shipment_id_fkey FOREIGN KEY (ppm_shipment_id) REFERENCES public.ppm_shipments(id);


--
-- Name: proof_of_service_docs proof_of_service_docs_payment_request_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.proof_of_service_docs
    ADD CONSTRAINT proof_of_service_docs_payment_request_id_fkey FOREIGN KEY (payment_request_id) REFERENCES public.payment_requests(id);


--
-- Name: re_contract_years re_contract_years_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_contract_years
    ADD CONSTRAINT re_contract_years_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_domestic_accessorial_prices re_domestic_accessorial_prices_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_accessorial_prices
    ADD CONSTRAINT re_domestic_accessorial_prices_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_domestic_accessorial_prices re_domestic_accessorial_prices_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_accessorial_prices
    ADD CONSTRAINT re_domestic_accessorial_prices_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.re_services(id);


--
-- Name: re_domestic_linehaul_prices re_domestic_linehaul_prices_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_linehaul_prices
    ADD CONSTRAINT re_domestic_linehaul_prices_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_domestic_linehaul_prices re_domestic_linehaul_prices_domestic_service_area_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_linehaul_prices
    ADD CONSTRAINT re_domestic_linehaul_prices_domestic_service_area_id_fkey FOREIGN KEY (domestic_service_area_id) REFERENCES public.re_domestic_service_areas(id);


--
-- Name: re_domestic_other_prices re_domestic_other_prices_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_other_prices
    ADD CONSTRAINT re_domestic_other_prices_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_domestic_other_prices re_domestic_other_prices_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_other_prices
    ADD CONSTRAINT re_domestic_other_prices_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.re_services(id);


--
-- Name: re_domestic_service_area_prices re_domestic_service_area_prices_domestic_service_area_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_service_area_prices
    ADD CONSTRAINT re_domestic_service_area_prices_domestic_service_area_id_fkey FOREIGN KEY (domestic_service_area_id) REFERENCES public.re_domestic_service_areas(id);


--
-- Name: re_domestic_service_area_prices re_domestic_service_area_prices_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_service_area_prices
    ADD CONSTRAINT re_domestic_service_area_prices_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.re_services(id);


--
-- Name: re_domestic_service_areas re_domestic_service_areas_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_service_areas
    ADD CONSTRAINT re_domestic_service_areas_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_domestic_service_area_prices re_domestic_services_area_prices_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_domestic_service_area_prices
    ADD CONSTRAINT re_domestic_services_area_prices_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_intl_accessorial_prices re_intl_accessorial_prices_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_accessorial_prices
    ADD CONSTRAINT re_intl_accessorial_prices_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_intl_accessorial_prices re_intl_accessorial_prices_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_accessorial_prices
    ADD CONSTRAINT re_intl_accessorial_prices_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.re_services(id);


--
-- Name: re_intl_other_prices re_intl_other_prices_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_other_prices
    ADD CONSTRAINT re_intl_other_prices_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_intl_other_prices re_intl_other_prices_rate_area_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_other_prices
    ADD CONSTRAINT re_intl_other_prices_rate_area_id_fkey FOREIGN KEY (rate_area_id) REFERENCES public.re_rate_areas(id);


--
-- Name: re_intl_other_prices re_intl_other_prices_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_other_prices
    ADD CONSTRAINT re_intl_other_prices_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.re_services(id);


--
-- Name: re_intl_prices re_intl_prices_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_prices
    ADD CONSTRAINT re_intl_prices_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_intl_prices re_intl_prices_destination_rate_area_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_prices
    ADD CONSTRAINT re_intl_prices_destination_rate_area_id_fkey FOREIGN KEY (destination_rate_area_id) REFERENCES public.re_rate_areas(id);


--
-- Name: re_intl_prices re_intl_prices_origin_rate_area_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_prices
    ADD CONSTRAINT re_intl_prices_origin_rate_area_id_fkey FOREIGN KEY (origin_rate_area_id) REFERENCES public.re_rate_areas(id);


--
-- Name: re_intl_prices re_intl_prices_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_intl_prices
    ADD CONSTRAINT re_intl_prices_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.re_services(id);


--
-- Name: re_rate_areas re_rate_areas_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_rate_areas
    ADD CONSTRAINT re_rate_areas_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_shipment_type_prices re_shipment_type_prices_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_shipment_type_prices
    ADD CONSTRAINT re_shipment_type_prices_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_shipment_type_prices re_shipment_type_prices_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_shipment_type_prices
    ADD CONSTRAINT re_shipment_type_prices_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.re_services(id);


--
-- Name: re_task_order_fees re_task_order_fees_contract_year_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_task_order_fees
    ADD CONSTRAINT re_task_order_fees_contract_year_id_fkey FOREIGN KEY (contract_year_id) REFERENCES public.re_contract_years(id);


--
-- Name: re_task_order_fees re_task_order_fees_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_task_order_fees
    ADD CONSTRAINT re_task_order_fees_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.re_services(id);


--
-- Name: re_zip3s re_zip3s_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_zip3s
    ADD CONSTRAINT re_zip3s_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_zip3s re_zip3s_domestic_service_area_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_zip3s
    ADD CONSTRAINT re_zip3s_domestic_service_area_id_fkey FOREIGN KEY (domestic_service_area_id) REFERENCES public.re_domestic_service_areas(id);


--
-- Name: re_zip3s re_zip3s_rate_area_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_zip3s
    ADD CONSTRAINT re_zip3s_rate_area_id_fkey FOREIGN KEY (rate_area_id) REFERENCES public.re_rate_areas(id);


--
-- Name: re_zip5_rate_areas re_zip5_rate_areas_contract_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_zip5_rate_areas
    ADD CONSTRAINT re_zip5_rate_areas_contract_id_fkey FOREIGN KEY (contract_id) REFERENCES public.re_contracts(id);


--
-- Name: re_zip5_rate_areas re_zip5_rate_areas_rate_area_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.re_zip5_rate_areas
    ADD CONSTRAINT re_zip5_rate_areas_rate_area_id_fkey FOREIGN KEY (rate_area_id) REFERENCES public.re_rate_areas(id);


--
-- Name: report_violations report_violations_report_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.report_violations
    ADD CONSTRAINT report_violations_report_id_fkey FOREIGN KEY (report_id) REFERENCES public.evaluation_reports(id);


--
-- Name: report_violations report_violations_violation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.report_violations
    ADD CONSTRAINT report_violations_violation_id_fkey FOREIGN KEY (violation_id) REFERENCES public.pws_violations(id);


--
-- Name: reweighs reweighs_shipment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reweighs
    ADD CONSTRAINT reweighs_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.mto_shipments(id);


--
-- Name: service_items_customer_contacts service_items_customer_contac_mtoservice_item_customer_con_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_items_customer_contacts
    ADD CONSTRAINT service_items_customer_contac_mtoservice_item_customer_con_fkey FOREIGN KEY (mtoservice_item_customer_contact_id) REFERENCES public.mto_service_item_customer_contacts(id);


--
-- Name: service_items_customer_contacts service_items_customer_contacts_mtoservice_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_items_customer_contacts
    ADD CONSTRAINT service_items_customer_contacts_mtoservice_item_id_fkey FOREIGN KEY (mtoservice_item_id) REFERENCES public.mto_service_items(id);


--
-- Name: notifications service_member_id___fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT service_member_id___fk FOREIGN KEY (service_member_id) REFERENCES public.service_members(id);


--
-- Name: service_members service_members_backup_mailing_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_members
    ADD CONSTRAINT service_members_backup_mailing_address_id_fkey FOREIGN KEY (backup_mailing_address_id) REFERENCES public.addresses(id);


--
-- Name: service_members service_members_duty_location_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_members
    ADD CONSTRAINT service_members_duty_location_id_fkey FOREIGN KEY (duty_location_id) REFERENCES public.duty_locations(id) ON DELETE SET NULL;


--
-- Name: service_members service_members_residential_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_members
    ADD CONSTRAINT service_members_residential_address_id_fkey FOREIGN KEY (residential_address_id) REFERENCES public.addresses(id);


--
-- Name: service_members service_members_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_members
    ADD CONSTRAINT service_members_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: service_params service_params_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_params
    ADD CONSTRAINT service_params_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.re_services(id);


--
-- Name: service_params service_params_service_item_param_key_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_params
    ADD CONSTRAINT service_params_service_item_param_key_id_fkey FOREIGN KEY (service_item_param_key_id) REFERENCES public.service_item_param_keys(id);


--
-- Name: service_request_document_uploads service_request_documents_contractor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_request_document_uploads
    ADD CONSTRAINT service_request_documents_contractor_id_fkey FOREIGN KEY (contractor_id) REFERENCES public.contractors(id);


--
-- Name: service_request_documents service_request_documents_mto_service_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_request_documents
    ADD CONSTRAINT service_request_documents_mto_service_item_id_fkey FOREIGN KEY (mto_service_item_id) REFERENCES public.mto_service_items(id);


--
-- Name: service_request_document_uploads service_request_documents_service_request_documents_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_request_document_uploads
    ADD CONSTRAINT service_request_documents_service_request_documents_id_fkey FOREIGN KEY (service_request_documents_id) REFERENCES public.service_request_documents(id);


--
-- Name: service_request_document_uploads service_request_documents_uploads_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_request_document_uploads
    ADD CONSTRAINT service_request_documents_uploads_id_fkey FOREIGN KEY (upload_id) REFERENCES public.uploads(id);


--
-- Name: shipment_address_updates shipment_address_updates_shipment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.shipment_address_updates
    ADD CONSTRAINT shipment_address_updates_shipment_id_fkey FOREIGN KEY (shipment_id) REFERENCES public.mto_shipments(id);


--
-- Name: signed_certifications signed_certifications_moves_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.signed_certifications
    ADD CONSTRAINT signed_certifications_moves_id_fk FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: signed_certifications signed_certifications_personally_procured_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.signed_certifications
    ADD CONSTRAINT signed_certifications_personally_procured_move_id_fkey FOREIGN KEY (personally_procured_move_id) REFERENCES public.personally_procured_moves(id);


--
-- Name: signed_certifications signed_certifications_submitting_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.signed_certifications
    ADD CONSTRAINT signed_certifications_submitting_user_id_fkey FOREIGN KEY (submitting_user_id) REFERENCES public.users(id);


--
-- Name: sit_address_updates sit_address_updates_mto_service_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sit_address_updates
    ADD CONSTRAINT sit_address_updates_mto_service_item_id_fkey FOREIGN KEY (mto_service_item_id) REFERENCES public.mto_service_items(id);


--
-- Name: sit_address_updates sit_address_updates_new_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sit_address_updates
    ADD CONSTRAINT sit_address_updates_new_address_id_fkey FOREIGN KEY (new_address_id) REFERENCES public.addresses(id);


--
-- Name: sit_address_updates sit_address_updates_old_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sit_address_updates
    ADD CONSTRAINT sit_address_updates_old_address_id_fkey FOREIGN KEY (old_address_id) REFERENCES public.addresses(id);


--
-- Name: sit_extensions sit_extensions_mto_shipment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sit_extensions
    ADD CONSTRAINT sit_extensions_mto_shipment_id_fkey FOREIGN KEY (mto_shipment_id) REFERENCES public.mto_shipments(id);


--
-- Name: storage_facilities storage_facilities_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.storage_facilities
    ADD CONSTRAINT storage_facilities_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id);


--
-- Name: transportation_offices transportation_offices_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transportation_offices
    ADD CONSTRAINT transportation_offices_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.addresses(id);


--
-- Name: transportation_offices transportation_offices_shipping_office_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.transportation_offices
    ADD CONSTRAINT transportation_offices_shipping_office_id_fkey FOREIGN KEY (shipping_office_id) REFERENCES public.transportation_offices(id);


--
-- Name: user_uploads user_uploads_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_uploads
    ADD CONSTRAINT user_uploads_document_id_fkey FOREIGN KEY (document_id) REFERENCES public.documents(id);


--
-- Name: user_uploads user_uploads_uploader_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_uploads
    ADD CONSTRAINT user_uploads_uploader_id_fkey FOREIGN KEY (uploader_id) REFERENCES public.users(id);


--
-- Name: user_uploads user_uploads_uploads_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.user_uploads
    ADD CONSTRAINT user_uploads_uploads_id_fkey FOREIGN KEY (upload_id) REFERENCES public.uploads(id) ON DELETE RESTRICT;


--
-- Name: users_roles users_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_roles
    ADD CONSTRAINT users_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: users_roles users_roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_roles
    ADD CONSTRAINT users_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: webhook_notifications webhook_notifications_move_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.webhook_notifications
    ADD CONSTRAINT webhook_notifications_move_id_fkey FOREIGN KEY (move_id) REFERENCES public.moves(id);


--
-- Name: webhook_subscriptions webhook_subscriptions_subscriber_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.webhook_subscriptions
    ADD CONSTRAINT webhook_subscriptions_subscriber_id_fkey FOREIGN KEY (subscriber_id) REFERENCES public.contractors(id);


--
-- Name: weight_tickets weight_tickets_empty_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weight_tickets
    ADD CONSTRAINT weight_tickets_empty_document_id_fkey FOREIGN KEY (empty_document_id) REFERENCES public.documents(id);


--
-- Name: weight_tickets weight_tickets_full_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weight_tickets
    ADD CONSTRAINT weight_tickets_full_document_id_fkey FOREIGN KEY (full_document_id) REFERENCES public.documents(id);


--
-- Name: weight_tickets weight_tickets_ppm_shipments_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weight_tickets
    ADD CONSTRAINT weight_tickets_ppm_shipments_id_fkey FOREIGN KEY (ppm_shipment_id) REFERENCES public.ppm_shipments(id);


--
-- Name: weight_tickets weight_tickets_proof_of_trailer_ownership_document_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weight_tickets
    ADD CONSTRAINT weight_tickets_proof_of_trailer_ownership_document_id_fkey FOREIGN KEY (proof_of_trailer_ownership_document_id) REFERENCES public.documents(id);


--
-- Name: FUNCTION gbtreekey16_in(cstring); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbtreekey16_in(cstring) TO master;
GRANT ALL ON FUNCTION public.gbtreekey16_in(cstring) TO ecs_user;


--
-- Name: FUNCTION gbtreekey16_out(public.gbtreekey16); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbtreekey16_out(public.gbtreekey16) TO master;
GRANT ALL ON FUNCTION public.gbtreekey16_out(public.gbtreekey16) TO ecs_user;


--
-- Name: FUNCTION gbtreekey32_in(cstring); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbtreekey32_in(cstring) TO master;
GRANT ALL ON FUNCTION public.gbtreekey32_in(cstring) TO ecs_user;


--
-- Name: FUNCTION gbtreekey32_out(public.gbtreekey32); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbtreekey32_out(public.gbtreekey32) TO master;
GRANT ALL ON FUNCTION public.gbtreekey32_out(public.gbtreekey32) TO ecs_user;


--
-- Name: FUNCTION gbtreekey4_in(cstring); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbtreekey4_in(cstring) TO master;
GRANT ALL ON FUNCTION public.gbtreekey4_in(cstring) TO ecs_user;


--
-- Name: FUNCTION gbtreekey4_out(public.gbtreekey4); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbtreekey4_out(public.gbtreekey4) TO master;
GRANT ALL ON FUNCTION public.gbtreekey4_out(public.gbtreekey4) TO ecs_user;


--
-- Name: FUNCTION gbtreekey8_in(cstring); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbtreekey8_in(cstring) TO master;
GRANT ALL ON FUNCTION public.gbtreekey8_in(cstring) TO ecs_user;


--
-- Name: FUNCTION gbtreekey8_out(public.gbtreekey8); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbtreekey8_out(public.gbtreekey8) TO master;
GRANT ALL ON FUNCTION public.gbtreekey8_out(public.gbtreekey8) TO ecs_user;


--
-- Name: FUNCTION gbtreekey_var_in(cstring); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbtreekey_var_in(cstring) TO master;
GRANT ALL ON FUNCTION public.gbtreekey_var_in(cstring) TO ecs_user;


--
-- Name: FUNCTION gbtreekey_var_out(public.gbtreekey_var); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbtreekey_var_out(public.gbtreekey_var) TO master;
GRANT ALL ON FUNCTION public.gbtreekey_var_out(public.gbtreekey_var) TO ecs_user;


--
-- Name: FUNCTION gtrgm_in(cstring); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gtrgm_in(cstring) TO master;
GRANT ALL ON FUNCTION public.gtrgm_in(cstring) TO ecs_user;


--
-- Name: FUNCTION gtrgm_out(public.gtrgm); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gtrgm_out(public.gtrgm) TO master;
GRANT ALL ON FUNCTION public.gtrgm_out(public.gtrgm) TO ecs_user;


--
-- Name: FUNCTION add_audit_history_table(target_table regclass); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.add_audit_history_table(target_table regclass) TO master;
GRANT ALL ON FUNCTION public.add_audit_history_table(target_table regclass) TO ecs_user;


--
-- Name: FUNCTION add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean) TO master;
GRANT ALL ON FUNCTION public.add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean) TO ecs_user;


--
-- Name: FUNCTION add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean, ignored_cols text[]); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean, ignored_cols text[]) TO master;
GRANT ALL ON FUNCTION public.add_audit_history_table(target_table regclass, audit_rows boolean, audit_query_text boolean, ignored_cols text[]) TO ecs_user;


--
-- Name: FUNCTION cash_dist(money, money); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.cash_dist(money, money) TO master;
GRANT ALL ON FUNCTION public.cash_dist(money, money) TO ecs_user;


--
-- Name: FUNCTION date_dist(date, date); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.date_dist(date, date) TO master;
GRANT ALL ON FUNCTION public.date_dist(date, date) TO ecs_user;


--
-- Name: FUNCTION f_unaccent(text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.f_unaccent(text) TO master;
GRANT ALL ON FUNCTION public.f_unaccent(text) TO ecs_user;


--
-- Name: FUNCTION float4_dist(real, real); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.float4_dist(real, real) TO master;
GRANT ALL ON FUNCTION public.float4_dist(real, real) TO ecs_user;


--
-- Name: FUNCTION float8_dist(double precision, double precision); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.float8_dist(double precision, double precision) TO master;
GRANT ALL ON FUNCTION public.float8_dist(double precision, double precision) TO ecs_user;


--
-- Name: FUNCTION gbt_bit_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bit_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bit_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bit_consistent(internal, bit, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bit_consistent(internal, bit, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bit_consistent(internal, bit, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bit_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bit_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bit_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bit_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bit_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bit_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bit_same(public.gbtreekey_var, public.gbtreekey_var, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bit_same(public.gbtreekey_var, public.gbtreekey_var, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bit_same(public.gbtreekey_var, public.gbtreekey_var, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bit_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bit_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bit_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bpchar_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bpchar_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bpchar_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bpchar_consistent(internal, character, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bpchar_consistent(internal, character, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bpchar_consistent(internal, character, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bytea_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bytea_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bytea_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bytea_consistent(internal, bytea, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bytea_consistent(internal, bytea, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bytea_consistent(internal, bytea, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bytea_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bytea_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bytea_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bytea_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bytea_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bytea_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bytea_same(public.gbtreekey_var, public.gbtreekey_var, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bytea_same(public.gbtreekey_var, public.gbtreekey_var, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bytea_same(public.gbtreekey_var, public.gbtreekey_var, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_bytea_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_bytea_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_bytea_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_cash_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_cash_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_cash_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_cash_consistent(internal, money, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_cash_consistent(internal, money, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_cash_consistent(internal, money, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_cash_distance(internal, money, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_cash_distance(internal, money, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_cash_distance(internal, money, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_cash_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_cash_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_cash_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_cash_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_cash_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_cash_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_cash_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_cash_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_cash_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_cash_same(public.gbtreekey16, public.gbtreekey16, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_cash_same(public.gbtreekey16, public.gbtreekey16, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_cash_same(public.gbtreekey16, public.gbtreekey16, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_cash_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_cash_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_cash_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_date_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_date_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_date_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_date_consistent(internal, date, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_date_consistent(internal, date, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_date_consistent(internal, date, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_date_distance(internal, date, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_date_distance(internal, date, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_date_distance(internal, date, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_date_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_date_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_date_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_date_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_date_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_date_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_date_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_date_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_date_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_date_same(public.gbtreekey8, public.gbtreekey8, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_date_same(public.gbtreekey8, public.gbtreekey8, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_date_same(public.gbtreekey8, public.gbtreekey8, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_date_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_date_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_date_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_decompress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_decompress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_decompress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_enum_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_enum_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_enum_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_enum_consistent(internal, anyenum, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_enum_consistent(internal, anyenum, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_enum_consistent(internal, anyenum, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_enum_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_enum_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_enum_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_enum_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_enum_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_enum_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_enum_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_enum_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_enum_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_enum_same(public.gbtreekey8, public.gbtreekey8, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_enum_same(public.gbtreekey8, public.gbtreekey8, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_enum_same(public.gbtreekey8, public.gbtreekey8, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_enum_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_enum_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_enum_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float4_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float4_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float4_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float4_consistent(internal, real, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float4_consistent(internal, real, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float4_consistent(internal, real, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float4_distance(internal, real, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float4_distance(internal, real, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float4_distance(internal, real, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float4_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float4_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float4_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float4_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float4_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float4_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float4_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float4_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float4_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float4_same(public.gbtreekey8, public.gbtreekey8, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float4_same(public.gbtreekey8, public.gbtreekey8, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float4_same(public.gbtreekey8, public.gbtreekey8, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float4_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float4_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float4_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float8_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float8_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float8_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float8_consistent(internal, double precision, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float8_consistent(internal, double precision, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float8_consistent(internal, double precision, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float8_distance(internal, double precision, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float8_distance(internal, double precision, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float8_distance(internal, double precision, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float8_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float8_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float8_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float8_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float8_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float8_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float8_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float8_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float8_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float8_same(public.gbtreekey16, public.gbtreekey16, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float8_same(public.gbtreekey16, public.gbtreekey16, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float8_same(public.gbtreekey16, public.gbtreekey16, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_float8_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_float8_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_float8_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_inet_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_inet_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_inet_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_inet_consistent(internal, inet, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_inet_consistent(internal, inet, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_inet_consistent(internal, inet, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_inet_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_inet_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_inet_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_inet_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_inet_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_inet_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_inet_same(public.gbtreekey16, public.gbtreekey16, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_inet_same(public.gbtreekey16, public.gbtreekey16, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_inet_same(public.gbtreekey16, public.gbtreekey16, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_inet_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_inet_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_inet_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int2_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int2_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int2_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int2_consistent(internal, smallint, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int2_consistent(internal, smallint, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int2_consistent(internal, smallint, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int2_distance(internal, smallint, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int2_distance(internal, smallint, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int2_distance(internal, smallint, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int2_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int2_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int2_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int2_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int2_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int2_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int2_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int2_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int2_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int2_same(public.gbtreekey4, public.gbtreekey4, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int2_same(public.gbtreekey4, public.gbtreekey4, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int2_same(public.gbtreekey4, public.gbtreekey4, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int2_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int2_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int2_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int4_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int4_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int4_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int4_consistent(internal, integer, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int4_consistent(internal, integer, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int4_consistent(internal, integer, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int4_distance(internal, integer, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int4_distance(internal, integer, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int4_distance(internal, integer, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int4_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int4_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int4_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int4_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int4_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int4_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int4_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int4_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int4_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int4_same(public.gbtreekey8, public.gbtreekey8, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int4_same(public.gbtreekey8, public.gbtreekey8, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int4_same(public.gbtreekey8, public.gbtreekey8, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int4_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int4_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int4_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int8_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int8_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int8_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int8_consistent(internal, bigint, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int8_consistent(internal, bigint, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int8_consistent(internal, bigint, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int8_distance(internal, bigint, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int8_distance(internal, bigint, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int8_distance(internal, bigint, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int8_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int8_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int8_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int8_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int8_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int8_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int8_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int8_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int8_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int8_same(public.gbtreekey16, public.gbtreekey16, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int8_same(public.gbtreekey16, public.gbtreekey16, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int8_same(public.gbtreekey16, public.gbtreekey16, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_int8_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_int8_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_int8_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_intv_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_intv_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_intv_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_intv_consistent(internal, interval, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_intv_consistent(internal, interval, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_intv_consistent(internal, interval, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_intv_decompress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_intv_decompress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_intv_decompress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_intv_distance(internal, interval, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_intv_distance(internal, interval, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_intv_distance(internal, interval, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_intv_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_intv_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_intv_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_intv_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_intv_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_intv_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_intv_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_intv_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_intv_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_intv_same(public.gbtreekey32, public.gbtreekey32, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_intv_same(public.gbtreekey32, public.gbtreekey32, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_intv_same(public.gbtreekey32, public.gbtreekey32, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_intv_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_intv_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_intv_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad8_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad8_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad8_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad8_consistent(internal, macaddr8, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad8_consistent(internal, macaddr8, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad8_consistent(internal, macaddr8, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad8_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad8_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad8_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad8_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad8_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad8_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad8_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad8_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad8_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad8_same(public.gbtreekey16, public.gbtreekey16, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad8_same(public.gbtreekey16, public.gbtreekey16, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad8_same(public.gbtreekey16, public.gbtreekey16, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad8_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad8_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad8_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad_consistent(internal, macaddr, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad_consistent(internal, macaddr, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad_consistent(internal, macaddr, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad_same(public.gbtreekey16, public.gbtreekey16, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad_same(public.gbtreekey16, public.gbtreekey16, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad_same(public.gbtreekey16, public.gbtreekey16, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_macad_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_macad_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_macad_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_numeric_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_numeric_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_numeric_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_numeric_consistent(internal, numeric, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_numeric_consistent(internal, numeric, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_numeric_consistent(internal, numeric, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_numeric_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_numeric_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_numeric_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_numeric_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_numeric_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_numeric_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_numeric_same(public.gbtreekey_var, public.gbtreekey_var, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_numeric_same(public.gbtreekey_var, public.gbtreekey_var, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_numeric_same(public.gbtreekey_var, public.gbtreekey_var, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_numeric_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_numeric_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_numeric_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_oid_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_oid_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_oid_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_oid_consistent(internal, oid, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_oid_consistent(internal, oid, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_oid_consistent(internal, oid, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_oid_distance(internal, oid, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_oid_distance(internal, oid, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_oid_distance(internal, oid, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_oid_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_oid_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_oid_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_oid_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_oid_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_oid_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_oid_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_oid_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_oid_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_oid_same(public.gbtreekey8, public.gbtreekey8, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_oid_same(public.gbtreekey8, public.gbtreekey8, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_oid_same(public.gbtreekey8, public.gbtreekey8, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_oid_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_oid_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_oid_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_text_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_text_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_text_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_text_consistent(internal, text, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_text_consistent(internal, text, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_text_consistent(internal, text, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_text_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_text_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_text_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_text_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_text_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_text_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_text_same(public.gbtreekey_var, public.gbtreekey_var, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_text_same(public.gbtreekey_var, public.gbtreekey_var, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_text_same(public.gbtreekey_var, public.gbtreekey_var, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_text_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_text_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_text_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_time_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_time_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_time_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_time_consistent(internal, time without time zone, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_time_consistent(internal, time without time zone, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_time_consistent(internal, time without time zone, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_time_distance(internal, time without time zone, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_time_distance(internal, time without time zone, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_time_distance(internal, time without time zone, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_time_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_time_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_time_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_time_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_time_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_time_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_time_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_time_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_time_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_time_same(public.gbtreekey16, public.gbtreekey16, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_time_same(public.gbtreekey16, public.gbtreekey16, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_time_same(public.gbtreekey16, public.gbtreekey16, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_time_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_time_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_time_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_timetz_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_timetz_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_timetz_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_timetz_consistent(internal, time with time zone, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_timetz_consistent(internal, time with time zone, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_timetz_consistent(internal, time with time zone, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_ts_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_ts_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_ts_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_ts_consistent(internal, timestamp without time zone, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_ts_consistent(internal, timestamp without time zone, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_ts_consistent(internal, timestamp without time zone, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_ts_distance(internal, timestamp without time zone, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_ts_distance(internal, timestamp without time zone, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_ts_distance(internal, timestamp without time zone, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_ts_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_ts_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_ts_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_ts_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_ts_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_ts_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_ts_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_ts_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_ts_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_ts_same(public.gbtreekey16, public.gbtreekey16, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_ts_same(public.gbtreekey16, public.gbtreekey16, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_ts_same(public.gbtreekey16, public.gbtreekey16, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_ts_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_ts_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_ts_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_tstz_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_tstz_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_tstz_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_tstz_consistent(internal, timestamp with time zone, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_tstz_consistent(internal, timestamp with time zone, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_tstz_consistent(internal, timestamp with time zone, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_tstz_distance(internal, timestamp with time zone, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_tstz_distance(internal, timestamp with time zone, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_tstz_distance(internal, timestamp with time zone, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_uuid_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_uuid_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_uuid_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_uuid_consistent(internal, uuid, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_uuid_consistent(internal, uuid, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_uuid_consistent(internal, uuid, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_uuid_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_uuid_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_uuid_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_uuid_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_uuid_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_uuid_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_uuid_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_uuid_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_uuid_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_uuid_same(public.gbtreekey32, public.gbtreekey32, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_uuid_same(public.gbtreekey32, public.gbtreekey32, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_uuid_same(public.gbtreekey32, public.gbtreekey32, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_uuid_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_uuid_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gbt_uuid_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gbt_var_decompress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_var_decompress(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_var_decompress(internal) TO ecs_user;


--
-- Name: FUNCTION gbt_var_fetch(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gbt_var_fetch(internal) TO master;
GRANT ALL ON FUNCTION public.gbt_var_fetch(internal) TO ecs_user;


--
-- Name: FUNCTION gin_extract_query_trgm(text, internal, smallint, internal, internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gin_extract_query_trgm(text, internal, smallint, internal, internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gin_extract_query_trgm(text, internal, smallint, internal, internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gin_extract_value_trgm(text, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gin_extract_value_trgm(text, internal) TO master;
GRANT ALL ON FUNCTION public.gin_extract_value_trgm(text, internal) TO ecs_user;


--
-- Name: FUNCTION gin_trgm_consistent(internal, smallint, text, integer, internal, internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gin_trgm_consistent(internal, smallint, text, integer, internal, internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gin_trgm_consistent(internal, smallint, text, integer, internal, internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gin_trgm_triconsistent(internal, smallint, text, integer, internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gin_trgm_triconsistent(internal, smallint, text, integer, internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gin_trgm_triconsistent(internal, smallint, text, integer, internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gtrgm_compress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gtrgm_compress(internal) TO master;
GRANT ALL ON FUNCTION public.gtrgm_compress(internal) TO ecs_user;


--
-- Name: FUNCTION gtrgm_consistent(internal, text, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gtrgm_consistent(internal, text, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gtrgm_consistent(internal, text, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gtrgm_decompress(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gtrgm_decompress(internal) TO master;
GRANT ALL ON FUNCTION public.gtrgm_decompress(internal) TO ecs_user;


--
-- Name: FUNCTION gtrgm_distance(internal, text, smallint, oid, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gtrgm_distance(internal, text, smallint, oid, internal) TO master;
GRANT ALL ON FUNCTION public.gtrgm_distance(internal, text, smallint, oid, internal) TO ecs_user;


--
-- Name: FUNCTION gtrgm_penalty(internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gtrgm_penalty(internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.gtrgm_penalty(internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION gtrgm_picksplit(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gtrgm_picksplit(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gtrgm_picksplit(internal, internal) TO ecs_user;


--
-- Name: FUNCTION gtrgm_same(public.gtrgm, public.gtrgm, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gtrgm_same(public.gtrgm, public.gtrgm, internal) TO master;
GRANT ALL ON FUNCTION public.gtrgm_same(public.gtrgm, public.gtrgm, internal) TO ecs_user;


--
-- Name: FUNCTION gtrgm_union(internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.gtrgm_union(internal, internal) TO master;
GRANT ALL ON FUNCTION public.gtrgm_union(internal, internal) TO ecs_user;


--
-- Name: FUNCTION if_modified_func(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.if_modified_func() TO master;
GRANT ALL ON FUNCTION public.if_modified_func() TO ecs_user;


--
-- Name: FUNCTION int2_dist(smallint, smallint); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.int2_dist(smallint, smallint) TO master;
GRANT ALL ON FUNCTION public.int2_dist(smallint, smallint) TO ecs_user;


--
-- Name: FUNCTION int4_dist(integer, integer); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.int4_dist(integer, integer) TO master;
GRANT ALL ON FUNCTION public.int4_dist(integer, integer) TO ecs_user;


--
-- Name: FUNCTION int8_dist(bigint, bigint); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.int8_dist(bigint, bigint) TO master;
GRANT ALL ON FUNCTION public.int8_dist(bigint, bigint) TO ecs_user;


--
-- Name: FUNCTION interval_dist(interval, interval); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.interval_dist(interval, interval) TO master;
GRANT ALL ON FUNCTION public.interval_dist(interval, interval) TO ecs_user;


--
-- Name: FUNCTION oid_dist(oid, oid); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.oid_dist(oid, oid) TO master;
GRANT ALL ON FUNCTION public.oid_dist(oid, oid) TO ecs_user;


--
-- Name: FUNCTION searchable_full_name(first_name text, last_name text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.searchable_full_name(first_name text, last_name text) TO master;
GRANT ALL ON FUNCTION public.searchable_full_name(first_name text, last_name text) TO ecs_user;


--
-- Name: FUNCTION set_limit(real); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.set_limit(real) TO master;
GRANT ALL ON FUNCTION public.set_limit(real) TO ecs_user;


--
-- Name: FUNCTION show_limit(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.show_limit() TO master;
GRANT ALL ON FUNCTION public.show_limit() TO ecs_user;


--
-- Name: FUNCTION show_trgm(text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.show_trgm(text) TO master;
GRANT ALL ON FUNCTION public.show_trgm(text) TO ecs_user;


--
-- Name: FUNCTION similarity(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.similarity(text, text) TO master;
GRANT ALL ON FUNCTION public.similarity(text, text) TO ecs_user;


--
-- Name: FUNCTION similarity_dist(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.similarity_dist(text, text) TO master;
GRANT ALL ON FUNCTION public.similarity_dist(text, text) TO ecs_user;


--
-- Name: FUNCTION similarity_op(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.similarity_op(text, text) TO master;
GRANT ALL ON FUNCTION public.similarity_op(text, text) TO ecs_user;


--
-- Name: FUNCTION strict_word_similarity(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.strict_word_similarity(text, text) TO master;
GRANT ALL ON FUNCTION public.strict_word_similarity(text, text) TO ecs_user;


--
-- Name: FUNCTION strict_word_similarity_commutator_op(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.strict_word_similarity_commutator_op(text, text) TO master;
GRANT ALL ON FUNCTION public.strict_word_similarity_commutator_op(text, text) TO ecs_user;


--
-- Name: FUNCTION strict_word_similarity_dist_commutator_op(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.strict_word_similarity_dist_commutator_op(text, text) TO master;
GRANT ALL ON FUNCTION public.strict_word_similarity_dist_commutator_op(text, text) TO ecs_user;


--
-- Name: FUNCTION strict_word_similarity_dist_op(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.strict_word_similarity_dist_op(text, text) TO master;
GRANT ALL ON FUNCTION public.strict_word_similarity_dist_op(text, text) TO ecs_user;


--
-- Name: FUNCTION strict_word_similarity_op(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.strict_word_similarity_op(text, text) TO master;
GRANT ALL ON FUNCTION public.strict_word_similarity_op(text, text) TO ecs_user;


--
-- Name: FUNCTION time_dist(time without time zone, time without time zone); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.time_dist(time without time zone, time without time zone) TO master;
GRANT ALL ON FUNCTION public.time_dist(time without time zone, time without time zone) TO ecs_user;


--
-- Name: FUNCTION ts_dist(timestamp without time zone, timestamp without time zone); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.ts_dist(timestamp without time zone, timestamp without time zone) TO master;
GRANT ALL ON FUNCTION public.ts_dist(timestamp without time zone, timestamp without time zone) TO ecs_user;


--
-- Name: FUNCTION tstz_dist(timestamp with time zone, timestamp with time zone); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.tstz_dist(timestamp with time zone, timestamp with time zone) TO master;
GRANT ALL ON FUNCTION public.tstz_dist(timestamp with time zone, timestamp with time zone) TO ecs_user;


--
-- Name: FUNCTION unaccent(text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.unaccent(text) TO master;
GRANT ALL ON FUNCTION public.unaccent(text) TO ecs_user;


--
-- Name: FUNCTION unaccent(regdictionary, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.unaccent(regdictionary, text) TO master;
GRANT ALL ON FUNCTION public.unaccent(regdictionary, text) TO ecs_user;


--
-- Name: FUNCTION unaccent_init(internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.unaccent_init(internal) TO master;
GRANT ALL ON FUNCTION public.unaccent_init(internal) TO ecs_user;


--
-- Name: FUNCTION unaccent_lexize(internal, internal, internal, internal); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.unaccent_lexize(internal, internal, internal, internal) TO master;
GRANT ALL ON FUNCTION public.unaccent_lexize(internal, internal, internal, internal) TO ecs_user;


--
-- Name: FUNCTION uuid_generate_v1(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_generate_v1() TO master;
GRANT ALL ON FUNCTION public.uuid_generate_v1() TO ecs_user;


--
-- Name: FUNCTION uuid_generate_v1mc(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_generate_v1mc() TO master;
GRANT ALL ON FUNCTION public.uuid_generate_v1mc() TO ecs_user;


--
-- Name: FUNCTION uuid_generate_v3(namespace uuid, name text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_generate_v3(namespace uuid, name text) TO master;
GRANT ALL ON FUNCTION public.uuid_generate_v3(namespace uuid, name text) TO ecs_user;


--
-- Name: FUNCTION uuid_generate_v4(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_generate_v4() TO master;
GRANT ALL ON FUNCTION public.uuid_generate_v4() TO ecs_user;


--
-- Name: FUNCTION uuid_generate_v5(namespace uuid, name text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_generate_v5(namespace uuid, name text) TO master;
GRANT ALL ON FUNCTION public.uuid_generate_v5(namespace uuid, name text) TO ecs_user;


--
-- Name: FUNCTION uuid_nil(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_nil() TO master;
GRANT ALL ON FUNCTION public.uuid_nil() TO ecs_user;


--
-- Name: FUNCTION uuid_ns_dns(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_ns_dns() TO master;
GRANT ALL ON FUNCTION public.uuid_ns_dns() TO ecs_user;


--
-- Name: FUNCTION uuid_ns_oid(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_ns_oid() TO master;
GRANT ALL ON FUNCTION public.uuid_ns_oid() TO ecs_user;


--
-- Name: FUNCTION uuid_ns_url(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_ns_url() TO master;
GRANT ALL ON FUNCTION public.uuid_ns_url() TO ecs_user;


--
-- Name: FUNCTION uuid_ns_x500(); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.uuid_ns_x500() TO master;
GRANT ALL ON FUNCTION public.uuid_ns_x500() TO ecs_user;


--
-- Name: FUNCTION word_similarity(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.word_similarity(text, text) TO master;
GRANT ALL ON FUNCTION public.word_similarity(text, text) TO ecs_user;


--
-- Name: FUNCTION word_similarity_commutator_op(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.word_similarity_commutator_op(text, text) TO master;
GRANT ALL ON FUNCTION public.word_similarity_commutator_op(text, text) TO ecs_user;


--
-- Name: FUNCTION word_similarity_dist_commutator_op(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.word_similarity_dist_commutator_op(text, text) TO master;
GRANT ALL ON FUNCTION public.word_similarity_dist_commutator_op(text, text) TO ecs_user;


--
-- Name: FUNCTION word_similarity_dist_op(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.word_similarity_dist_op(text, text) TO master;
GRANT ALL ON FUNCTION public.word_similarity_dist_op(text, text) TO ecs_user;


--
-- Name: FUNCTION word_similarity_op(text, text); Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON FUNCTION public.word_similarity_op(text, text) TO master;
GRANT ALL ON FUNCTION public.word_similarity_op(text, text) TO ecs_user;


--
-- Name: TABLE addresses; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.addresses TO master;
GRANT ALL ON TABLE public.addresses TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.addresses TO crud;


--
-- Name: TABLE admin_users; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.admin_users TO master;
GRANT ALL ON TABLE public.admin_users TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.admin_users TO crud;


--
-- Name: TABLE archived_access_codes; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.archived_access_codes TO master;
GRANT ALL ON TABLE public.archived_access_codes TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.archived_access_codes TO crud;


--
-- Name: TABLE archived_move_documents; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.archived_move_documents TO master;
GRANT ALL ON TABLE public.archived_move_documents TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.archived_move_documents TO crud;


--
-- Name: TABLE archived_moving_expense_documents; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.archived_moving_expense_documents TO master;
GRANT ALL ON TABLE public.archived_moving_expense_documents TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.archived_moving_expense_documents TO crud;


--
-- Name: TABLE archived_personally_procured_moves; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.archived_personally_procured_moves TO master;
GRANT ALL ON TABLE public.archived_personally_procured_moves TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.archived_personally_procured_moves TO crud;


--
-- Name: TABLE archived_reimbursements; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.archived_reimbursements TO master;
GRANT ALL ON TABLE public.archived_reimbursements TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.archived_reimbursements TO crud;


--
-- Name: TABLE archived_signed_certifications; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.archived_signed_certifications TO master;
GRANT ALL ON TABLE public.archived_signed_certifications TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.archived_signed_certifications TO crud;


--
-- Name: TABLE archived_weight_ticket_set_documents; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.archived_weight_ticket_set_documents TO master;
GRANT ALL ON TABLE public.archived_weight_ticket_set_documents TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.archived_weight_ticket_set_documents TO crud;


--
-- Name: TABLE audit_history; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.audit_history TO master;
GRANT ALL ON TABLE public.audit_history TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.audit_history TO crud;


--
-- Name: TABLE audit_history_tableslist; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.audit_history_tableslist TO master;
GRANT ALL ON TABLE public.audit_history_tableslist TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.audit_history_tableslist TO crud;


--
-- Name: TABLE backup_contacts; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.backup_contacts TO master;
GRANT ALL ON TABLE public.backup_contacts TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.backup_contacts TO crud;


--
-- Name: TABLE client_certs; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.client_certs TO master;
GRANT ALL ON TABLE public.client_certs TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.client_certs TO crud;


--
-- Name: TABLE contractors; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.contractors TO master;
GRANT ALL ON TABLE public.contractors TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.contractors TO crud;


--
-- Name: TABLE customer_support_remarks; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.customer_support_remarks TO master;
GRANT ALL ON TABLE public.customer_support_remarks TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.customer_support_remarks TO crud;


--
-- Name: TABLE distance_calculations; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.distance_calculations TO master;
GRANT ALL ON TABLE public.distance_calculations TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.distance_calculations TO crud;


--
-- Name: TABLE documents; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.documents TO master;
GRANT ALL ON TABLE public.documents TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.documents TO crud;


--
-- Name: TABLE duty_location_names; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.duty_location_names TO master;
GRANT ALL ON TABLE public.duty_location_names TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.duty_location_names TO crud;


--
-- Name: TABLE duty_locations; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.duty_locations TO master;
GRANT ALL ON TABLE public.duty_locations TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.duty_locations TO crud;


--
-- Name: TABLE edi_errors; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.edi_errors TO master;
GRANT ALL ON TABLE public.edi_errors TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.edi_errors TO crud;


--
-- Name: TABLE edi_processings; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.edi_processings TO master;
GRANT ALL ON TABLE public.edi_processings TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.edi_processings TO crud;


--
-- Name: TABLE electronic_orders; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.electronic_orders TO master;
GRANT ALL ON TABLE public.electronic_orders TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.electronic_orders TO crud;


--
-- Name: TABLE electronic_orders_revisions; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.electronic_orders_revisions TO master;
GRANT ALL ON TABLE public.electronic_orders_revisions TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.electronic_orders_revisions TO crud;


--
-- Name: TABLE entitlements; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.entitlements TO master;
GRANT ALL ON TABLE public.entitlements TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.entitlements TO crud;


--
-- Name: TABLE evaluation_reports; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.evaluation_reports TO master;
GRANT ALL ON TABLE public.evaluation_reports TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.evaluation_reports TO crud;


--
-- Name: TABLE fuel_eia_diesel_prices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.fuel_eia_diesel_prices TO master;
GRANT ALL ON TABLE public.fuel_eia_diesel_prices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.fuel_eia_diesel_prices TO crud;


--
-- Name: TABLE ghc_diesel_fuel_prices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.ghc_diesel_fuel_prices TO master;
GRANT ALL ON TABLE public.ghc_diesel_fuel_prices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.ghc_diesel_fuel_prices TO crud;


--
-- Name: TABLE ghc_domestic_transit_times; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.ghc_domestic_transit_times TO master;
GRANT ALL ON TABLE public.ghc_domestic_transit_times TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.ghc_domestic_transit_times TO crud;


--
-- Name: SEQUENCE interchange_control_number; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON SEQUENCE public.interchange_control_number TO master;
GRANT ALL ON SEQUENCE public.interchange_control_number TO ecs_user;
GRANT USAGE,UPDATE ON SEQUENCE public.interchange_control_number TO crud;


--
-- Name: TABLE invoice_number_trackers; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.invoice_number_trackers TO master;
GRANT ALL ON TABLE public.invoice_number_trackers TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.invoice_number_trackers TO crud;


--
-- Name: TABLE invoices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.invoices TO master;
GRANT ALL ON TABLE public.invoices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.invoices TO crud;


--
-- Name: TABLE jppso_region_state_assignments; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.jppso_region_state_assignments TO master;
GRANT ALL ON TABLE public.jppso_region_state_assignments TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.jppso_region_state_assignments TO crud;


--
-- Name: TABLE jppso_regions; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.jppso_regions TO master;
GRANT ALL ON TABLE public.jppso_regions TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.jppso_regions TO crud;


--
-- Name: TABLE mto_shipments; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.mto_shipments TO master;
GRANT ALL ON TABLE public.mto_shipments TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.mto_shipments TO crud;


--
-- Name: TABLE postal_code_to_gblocs; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.postal_code_to_gblocs TO master;
GRANT ALL ON TABLE public.postal_code_to_gblocs TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.postal_code_to_gblocs TO crud;


--
-- Name: TABLE ppm_shipments; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.ppm_shipments TO master;
GRANT ALL ON TABLE public.ppm_shipments TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.ppm_shipments TO crud;


--
-- Name: TABLE move_to_gbloc; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.move_to_gbloc TO master;
GRANT ALL ON TABLE public.move_to_gbloc TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.move_to_gbloc TO crud;


--
-- Name: TABLE moves; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.moves TO master;
GRANT ALL ON TABLE public.moves TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.moves TO crud;


--
-- Name: TABLE moving_expenses; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.moving_expenses TO master;
GRANT ALL ON TABLE public.moving_expenses TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.moving_expenses TO crud;


--
-- Name: TABLE mto_agents; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.mto_agents TO master;
GRANT ALL ON TABLE public.mto_agents TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.mto_agents TO crud;


--
-- Name: TABLE mto_service_item_customer_contacts; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.mto_service_item_customer_contacts TO master;
GRANT ALL ON TABLE public.mto_service_item_customer_contacts TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.mto_service_item_customer_contacts TO crud;


--
-- Name: TABLE mto_service_item_dimensions; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.mto_service_item_dimensions TO master;
GRANT ALL ON TABLE public.mto_service_item_dimensions TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.mto_service_item_dimensions TO crud;


--
-- Name: TABLE mto_service_items; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.mto_service_items TO master;
GRANT ALL ON TABLE public.mto_service_items TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.mto_service_items TO crud;


--
-- Name: TABLE notifications; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.notifications TO master;
GRANT ALL ON TABLE public.notifications TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.notifications TO crud;


--
-- Name: TABLE office_emails; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.office_emails TO master;
GRANT ALL ON TABLE public.office_emails TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.office_emails TO crud;


--
-- Name: TABLE office_phone_lines; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.office_phone_lines TO master;
GRANT ALL ON TABLE public.office_phone_lines TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.office_phone_lines TO crud;


--
-- Name: TABLE office_users; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.office_users TO master;
GRANT ALL ON TABLE public.office_users TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.office_users TO crud;


--
-- Name: TABLE orders; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.orders TO master;
GRANT ALL ON TABLE public.orders TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.orders TO crud;


--
-- Name: TABLE organizations; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.organizations TO master;
GRANT ALL ON TABLE public.organizations TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.organizations TO crud;


--
-- Name: TABLE payment_request_to_interchange_control_numbers; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.payment_request_to_interchange_control_numbers TO master;
GRANT ALL ON TABLE public.payment_request_to_interchange_control_numbers TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.payment_request_to_interchange_control_numbers TO crud;


--
-- Name: TABLE payment_requests; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.payment_requests TO master;
GRANT ALL ON TABLE public.payment_requests TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.payment_requests TO crud;


--
-- Name: TABLE payment_service_item_params; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.payment_service_item_params TO master;
GRANT ALL ON TABLE public.payment_service_item_params TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.payment_service_item_params TO crud;


--
-- Name: TABLE payment_service_items; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.payment_service_items TO master;
GRANT ALL ON TABLE public.payment_service_items TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.payment_service_items TO crud;


--
-- Name: TABLE personally_procured_moves; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.personally_procured_moves TO master;
GRANT ALL ON TABLE public.personally_procured_moves TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.personally_procured_moves TO crud;


--
-- Name: TABLE prime_uploads; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.prime_uploads TO master;
GRANT ALL ON TABLE public.prime_uploads TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.prime_uploads TO crud;


--
-- Name: TABLE progear_weight_tickets; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.progear_weight_tickets TO master;
GRANT ALL ON TABLE public.progear_weight_tickets TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.progear_weight_tickets TO crud;


--
-- Name: TABLE proof_of_service_docs; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.proof_of_service_docs TO master;
GRANT ALL ON TABLE public.proof_of_service_docs TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.proof_of_service_docs TO crud;


--
-- Name: TABLE pws_violations; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.pws_violations TO master;
GRANT ALL ON TABLE public.pws_violations TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.pws_violations TO crud;


--
-- Name: TABLE re_contract_years; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_contract_years TO master;
GRANT ALL ON TABLE public.re_contract_years TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_contract_years TO crud;


--
-- Name: TABLE re_contracts; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_contracts TO master;
GRANT ALL ON TABLE public.re_contracts TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_contracts TO crud;


--
-- Name: TABLE re_domestic_accessorial_prices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_domestic_accessorial_prices TO master;
GRANT ALL ON TABLE public.re_domestic_accessorial_prices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_domestic_accessorial_prices TO crud;


--
-- Name: TABLE re_domestic_linehaul_prices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_domestic_linehaul_prices TO master;
GRANT ALL ON TABLE public.re_domestic_linehaul_prices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_domestic_linehaul_prices TO crud;


--
-- Name: TABLE re_domestic_other_prices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_domestic_other_prices TO master;
GRANT ALL ON TABLE public.re_domestic_other_prices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_domestic_other_prices TO crud;


--
-- Name: TABLE re_domestic_service_area_prices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_domestic_service_area_prices TO master;
GRANT ALL ON TABLE public.re_domestic_service_area_prices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_domestic_service_area_prices TO crud;


--
-- Name: TABLE re_domestic_service_areas; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_domestic_service_areas TO master;
GRANT ALL ON TABLE public.re_domestic_service_areas TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_domestic_service_areas TO crud;


--
-- Name: TABLE re_intl_accessorial_prices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_intl_accessorial_prices TO master;
GRANT ALL ON TABLE public.re_intl_accessorial_prices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_intl_accessorial_prices TO crud;


--
-- Name: TABLE re_intl_other_prices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_intl_other_prices TO master;
GRANT ALL ON TABLE public.re_intl_other_prices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_intl_other_prices TO crud;


--
-- Name: TABLE re_intl_prices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_intl_prices TO master;
GRANT ALL ON TABLE public.re_intl_prices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_intl_prices TO crud;


--
-- Name: TABLE re_rate_areas; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_rate_areas TO master;
GRANT ALL ON TABLE public.re_rate_areas TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_rate_areas TO crud;


--
-- Name: TABLE re_services; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_services TO master;
GRANT ALL ON TABLE public.re_services TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_services TO crud;


--
-- Name: TABLE re_shipment_type_prices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_shipment_type_prices TO master;
GRANT ALL ON TABLE public.re_shipment_type_prices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_shipment_type_prices TO crud;


--
-- Name: TABLE re_task_order_fees; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_task_order_fees TO master;
GRANT ALL ON TABLE public.re_task_order_fees TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_task_order_fees TO crud;


--
-- Name: TABLE re_zip3s; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_zip3s TO master;
GRANT ALL ON TABLE public.re_zip3s TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_zip3s TO crud;


--
-- Name: TABLE re_zip5_rate_areas; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.re_zip5_rate_areas TO master;
GRANT ALL ON TABLE public.re_zip5_rate_areas TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.re_zip5_rate_areas TO crud;


--
-- Name: TABLE report_violations; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.report_violations TO master;
GRANT ALL ON TABLE public.report_violations TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.report_violations TO crud;


--
-- Name: TABLE reweighs; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.reweighs TO master;
GRANT ALL ON TABLE public.reweighs TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.reweighs TO crud;


--
-- Name: TABLE roles; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.roles TO master;
GRANT ALL ON TABLE public.roles TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.roles TO crud;


--
-- Name: TABLE schema_migration; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.schema_migration TO master;
GRANT ALL ON TABLE public.schema_migration TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.schema_migration TO crud;


--
-- Name: TABLE service_item_param_keys; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.service_item_param_keys TO master;
GRANT ALL ON TABLE public.service_item_param_keys TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.service_item_param_keys TO crud;


--
-- Name: TABLE service_items_customer_contacts; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.service_items_customer_contacts TO master;
GRANT ALL ON TABLE public.service_items_customer_contacts TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.service_items_customer_contacts TO crud;


--
-- Name: TABLE service_members; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.service_members TO master;
GRANT ALL ON TABLE public.service_members TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.service_members TO crud;


--
-- Name: TABLE service_params; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.service_params TO master;
GRANT ALL ON TABLE public.service_params TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.service_params TO crud;


--
-- Name: TABLE service_request_document_uploads; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.service_request_document_uploads TO master;
GRANT ALL ON TABLE public.service_request_document_uploads TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.service_request_document_uploads TO crud;


--
-- Name: TABLE service_request_documents; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.service_request_documents TO master;
GRANT ALL ON TABLE public.service_request_documents TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.service_request_documents TO crud;


--
-- Name: TABLE shipment_address_updates; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.shipment_address_updates TO master;
GRANT ALL ON TABLE public.shipment_address_updates TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.shipment_address_updates TO crud;


--
-- Name: TABLE signed_certifications; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.signed_certifications TO master;
GRANT ALL ON TABLE public.signed_certifications TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.signed_certifications TO crud;


--
-- Name: TABLE sit_address_updates; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.sit_address_updates TO master;
GRANT ALL ON TABLE public.sit_address_updates TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.sit_address_updates TO crud;


--
-- Name: TABLE sit_extensions; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.sit_extensions TO master;
GRANT ALL ON TABLE public.sit_extensions TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.sit_extensions TO crud;


--
-- Name: TABLE storage_facilities; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.storage_facilities TO master;
GRANT ALL ON TABLE public.storage_facilities TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.storage_facilities TO crud;


--
-- Name: TABLE transportation_accounting_codes; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.transportation_accounting_codes TO master;
GRANT ALL ON TABLE public.transportation_accounting_codes TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.transportation_accounting_codes TO crud;


--
-- Name: TABLE transportation_offices; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.transportation_offices TO master;
GRANT ALL ON TABLE public.transportation_offices TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.transportation_offices TO crud;


--
-- Name: TABLE uploads; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.uploads TO master;
GRANT ALL ON TABLE public.uploads TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.uploads TO crud;


--
-- Name: TABLE user_uploads; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.user_uploads TO master;
GRANT ALL ON TABLE public.user_uploads TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.user_uploads TO crud;


--
-- Name: TABLE users; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.users TO master;
GRANT ALL ON TABLE public.users TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.users TO crud;


--
-- Name: TABLE users_roles; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.users_roles TO master;
GRANT ALL ON TABLE public.users_roles TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.users_roles TO crud;


--
-- Name: TABLE webhook_notifications; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.webhook_notifications TO master;
GRANT ALL ON TABLE public.webhook_notifications TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.webhook_notifications TO crud;


--
-- Name: TABLE webhook_subscriptions; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.webhook_subscriptions TO master;
GRANT ALL ON TABLE public.webhook_subscriptions TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.webhook_subscriptions TO crud;


--
-- Name: TABLE weight_tickets; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.weight_tickets TO master;
GRANT ALL ON TABLE public.weight_tickets TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.weight_tickets TO crud;


--
-- Name: TABLE zip3_distances; Type: ACL; Schema: public; Owner: postgres
--

GRANT ALL ON TABLE public.zip3_distances TO master;
GRANT ALL ON TABLE public.zip3_distances TO ecs_user;
GRANT SELECT,INSERT,DELETE,UPDATE ON TABLE public.zip3_distances TO crud;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: -; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON SEQUENCES  TO master;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON SEQUENCES  TO ecs_user;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT USAGE,UPDATE ON SEQUENCES  TO crud;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: -; Owner: master
--

ALTER DEFAULT PRIVILEGES FOR ROLE master GRANT ALL ON SEQUENCES  TO ecs_user;


--
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: -; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON FUNCTIONS  TO master;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON FUNCTIONS  TO ecs_user;


--
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: -; Owner: master
--

ALTER DEFAULT PRIVILEGES FOR ROLE master GRANT ALL ON FUNCTIONS  TO ecs_user;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: -; Owner: postgres
--

ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON TABLES  TO master;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL ON TABLES  TO ecs_user;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT SELECT,INSERT,DELETE,UPDATE ON TABLES  TO crud;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: -; Owner: master
--

ALTER DEFAULT PRIVILEGES FOR ROLE master GRANT ALL ON TABLES  TO ecs_user;


--
-- PostgreSQL database dump complete
--

