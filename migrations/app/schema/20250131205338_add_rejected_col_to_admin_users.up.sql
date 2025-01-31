-- Adds new columns to office_users table
ALTER TABLE public.office_users
ADD COLUMN IF NOT EXISTS rejected_on timestamptz;

-- Comments on new columns
COMMENT on COLUMN office_users.rejected_on IS 'Date requested office users were rejected.';
