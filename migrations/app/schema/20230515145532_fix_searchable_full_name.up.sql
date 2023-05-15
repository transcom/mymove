-- POP RAW MIGRATION --

-- Make sure to reference f_unaccent function by schema public,
-- otherwise it may not be found
CREATE OR REPLACE FUNCTION searchable_full_name(first_name text, last_name text)
        RETURNS text
        LANGUAGE sql
        IMMUTABLE STRICT
AS
$function$
        -- CONCAT_WS is immutable when given only text arguments
SELECT public.f_unaccent(LOWER(CONCAT_WS(' ', first_name, last_name)));
$function$;

COMMENT ON FUNCTION searchable_full_name(first_name text, last_name text) IS 'Prepares a first/last name for search by lowercasing and removing diacritics';
