CREATE TYPE office_user_status AS enum (
    'APPROVED',
	'REQUESTED',
    'REJECTED'
);

-- Adds new columns to office_users table
-- making status default of APPROVED to ensure past users using MM are already approved
ALTER TABLE office_users
ADD COLUMN IF NOT EXISTS status office_user_status DEFAULT 'APPROVED' NULL,
ADD COLUMN IF NOT EXISTS edipi TEXT UNIQUE DEFAULT NULL,
ADD COLUMN IF NOT EXISTS other_unique_id TEXT UNIQUE DEFAULT NULL,
ADD COLUMN IF NOT EXISTS rejection_reason TEXT DEFAULT NULL;

-- Comments on new columns
COMMENT on COLUMN office_users.status IS 'Status of an office user account';
COMMENT on COLUMN office_users.edipi IS 'DoD ID or EDIPI of office user';
COMMENT on COLUMN office_users.other_unique_id IS 'Other unique id for PIV or ECA cert users';
COMMENT on COLUMN office_users.rejection_reason IS 'Rejection reason when account request is rejected by an admin';
