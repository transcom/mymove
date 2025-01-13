ALTER TABLE entitlements
ADD column weight_restriction int;

COMMENT ON COLUMN entitlements.weight_restriction IS 'The weight restricted on the move to a particular location';
