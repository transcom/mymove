-- These column additions will handle locking a move
-- when an office user is working on them
-- Add columns if they don't exist
ALTER TABLE moves
ADD COLUMN IF NOT EXISTS locked_by UUID NULL,
ADD COLUMN IF NOT EXISTS lock_expires_at TIMESTAMP WITH TIME ZONE NULL;

-- Add foreign key constraint to office_users table
ALTER TABLE moves
ADD CONSTRAINT fk_locked_by
FOREIGN KEY (locked_by)
REFERENCES office_users(id);

-- Add comments for the columns
COMMENT ON COLUMN moves.locked_by IS 'The id of the office user that locked the move.';
COMMENT ON COLUMN moves.lock_expires_at IS 'The expiration time that a move is locked until, the default value of this will be 30 minutes from initial lock.';
