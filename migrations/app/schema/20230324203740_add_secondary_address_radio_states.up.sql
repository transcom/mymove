ALTER TABLE mto_shipments
	ADD COLUMN has_secondary_pickup_address   bool,
	ADD COLUMN has_secondary_delivery_address bool;

COMMENT ON COLUMN mto_shipments.has_secondary_pickup_address IS 'False if the shipment does not have a secondary pickup address. This column exists to make it possible to tell whether a shipment update should delete an address or not modify it.';
COMMENT ON COLUMN mto_shipments.has_secondary_delivery_address IS 'False if the shipment does not have a secondary delivery address. This column exists to make it possible to tell whether a shipment update should delete an address or not modify it.';
