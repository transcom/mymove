-- Removes hackerone users from tsp and office users tables.

DELETE FROM public.office_users where email like '%(@hackerone.com|@wearehackerone.com|@managed.hackerone.com)';
