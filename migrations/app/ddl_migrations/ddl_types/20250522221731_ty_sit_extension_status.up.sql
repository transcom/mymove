--B-23206   Konstance Haffaney    Add REMOVED status

ALTER TYPE sit_extension_status
ADD VALUE IF NOT EXISTS 'REMOVED';

COMMENT ON COLUMN sit_extensions.status IS 'Status of this SIT Extension (Pending, Approved, Removed, or Denied).';