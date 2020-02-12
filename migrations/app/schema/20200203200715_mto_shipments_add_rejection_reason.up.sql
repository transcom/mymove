-- Column add
ALTER TABLE mto_shipments
ADD COLUMN rejection_reason text;

-- Column comment
COMMENT ON COLUMN mto_shipments.rejection_reason IS 'TOO will provide a reason if they reject the shipment';