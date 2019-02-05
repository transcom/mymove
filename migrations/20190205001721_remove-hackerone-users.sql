-- Removes hackerone users from tsp and office users tables.

DELETE FROM public.tsp_users where email like '%hackerone%'
DELETE FROM public.office_users where email like '%hackerone%'