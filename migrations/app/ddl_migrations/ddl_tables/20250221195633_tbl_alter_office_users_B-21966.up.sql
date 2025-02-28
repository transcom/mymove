--B-21966   Konstance Haffaney   Add column to office_users

ALTER TABLE public.office_users
ADD COLUMN IF NOT EXISTS rejected_on timestamptz;

COMMENT on COLUMN office_users.rejected_on IS 'Date requested office users were rejected.';