-- "Removing" diacritics is a surprisingly awkward problem
-- Diacritics are not just decorative, they are an important part of a person's name.
-- However, when searching for names, we should not expect office users to be able to
-- type diacritics/get the right ones, so we want to omit them for search, so that
-- we don't make it harder to search for names that have diacritics.
-- Postgres has an unaccent function that does exactly what we want,
-- BUT it is stable, not immutable, so we cannot use it in an index.
-- Many people suggest wrapping it in a function that is declared immutable
--   eg: https://stackoverflow.com/a/11007216
-- The other main option is to implement our own function to remove accents, which is
-- what I've tentatively gone with here, although it also feels like a bad option.
-- This approach to removing diacritics should work in the majority of cases that we will
-- see, but is certainly not going to work in all cases. The nice thing about it
-- is that it is actually immutable in the postgres sense and will not break
-- our index if we update the database.


CREATE OR REPLACE FUNCTION lower_unaccent(text text)
	RETURNS text
	LANGUAGE sql
	IMMUTABLE STRICT
AS
$function$
SELECT TRANSLATE(
		   LOWER(text),
		   'àáâãäåèéêëìíîïòóôõöùúûüñ',
		   'aaaaaaeeeeiiiiooooouuuun');
$function$;

COMMENT ON FUNCTION lower_unaccent(text text) IS 'Lowercase and remove common accents';

CREATE OR REPLACE FUNCTION searchable_full_name(first_name text, last_name text)
	RETURNS text
	LANGUAGE sql
	IMMUTABLE STRICT
AS
$function$
	-- CONCAT_WS is immutable when given only text arguments
SELECT lower_unaccent(CONCAT_WS(' ', first_name, last_name));
$function$;

COMMENT ON FUNCTION searchable_full_name(first_name text, last_name text) IS 'Prepares a first/last name for search by lowercasing and removing diacritics';

CREATE INDEX service_members_searchable_full_name_trgm_idx
	ON service_members
		USING gin (searchable_full_name(first_name, last_name) gin_trgm_ops);
