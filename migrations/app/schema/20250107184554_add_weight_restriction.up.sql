ALTER TABLE entitlements
ADD column  IF NOT EXISTS weight_restriction int;

COMMENT ON COLUMN entitlements.weight_restriction IS 'The weight restricted on the move to a particular location';
