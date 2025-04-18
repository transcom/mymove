-- B-22924  Daniel Jordan  adding sit_extensions table to move history so we can track the activity

DROP TRIGGER IF EXISTS audit_trigger_row ON public.sit_extensions;

CREATE TRIGGER audit_trigger_row
AFTER INSERT OR DELETE OR UPDATE
ON public.sit_extensions
FOR EACH ROW
EXECUTE FUNCTION if_modified_func('true', '{created_at, updated_at}');
