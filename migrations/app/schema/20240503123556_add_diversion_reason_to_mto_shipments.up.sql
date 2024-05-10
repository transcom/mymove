-- Adds new column to mto shipments table
ALTER TABLE mto_shipments
ADD COLUMN IF NOT EXISTS diversion_reason TEXT NULL;

-- Comments on new column
COMMENT ON COLUMN mto_shipments.diversion_reason IS 'Stores the reason for a requested diversion.';
