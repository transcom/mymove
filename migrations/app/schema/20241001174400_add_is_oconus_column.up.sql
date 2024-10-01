-- adding column for conus or oconus address
ALTER TABLE addresses
ADD COLUMN IF NOT EXISTS is_oconus boolean;

-- column comments
COMMENT ON COLUMN addresses.is_oconus IS 'Indicates whether address is CONUS (false) or OCONUS (true)';