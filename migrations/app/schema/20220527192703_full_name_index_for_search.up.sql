-- POP RAW MIGRATION --

-- "Removing" diacritics is a surprisingly awkward problem
-- Diacritics are not just decorative, they are an important part of a person's name.
-- However, when searching for names, we should not expect office users to be able to
-- type diacritics/get the right ones, so we want to omit them for search, so that
-- we don't make it harder to search for names that have diacritics.
-- Postgres has an unaccent function that does exactly what we want,
-- BUT it is stable, not immutable, so we cannot use it in an index.
-- Many people suggest wrapping it in a function that is declared immutable
--   eg: https://stackoverflow.com/a/11007216

CREATE EXTENSION unaccent;

CREATE OR REPLACE FUNCTION public.immutable_unaccent(regdictionary, text)
	RETURNS text LANGUAGE c IMMUTABLE PARALLEL SAFE STRICT AS
'$libdir/unaccent', 'unaccent_dict';

COMMENT ON FUNCTION immutable_unaccent(regdictionary, text) IS 'Do not use outside of the wrapper f_unnacent! This is a copy of the C unaccent function that we are marking as IMMUTABLE';

CREATE OR REPLACE FUNCTION public.f_unaccent(text)
	RETURNS text LANGUAGE sql IMMUTABLE PARALLEL SAFE STRICT AS
$func$
SELECT public.immutable_unaccent(regdictionary 'public.unaccent', $1)
$func$;

COMMENT ON FUNCTION f_unaccent(text) IS 'Wrapper around unaccent that is marked as immutable so it can be used in indexes';


CREATE OR REPLACE FUNCTION searchable_full_name(first_name text, last_name text)
	RETURNS text
	LANGUAGE sql
	IMMUTABLE STRICT
AS
$function$
	-- CONCAT_WS is immutable when given only text arguments
SELECT f_unaccent(LOWER(CONCAT_WS(' ', first_name, last_name)));
$function$;

COMMENT ON FUNCTION searchable_full_name(first_name text, last_name text) IS 'Prepares a first/last name for search by lowercasing and removing diacritics';

CREATE INDEX service_members_searchable_full_name_trgm_idx
	ON service_members
		USING gin (searchable_full_name(first_name, last_name) gin_trgm_ops);
