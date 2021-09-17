-- Add default value for existing HHG shipments
UPDATE mto_shipments
SET sit_days_allowance = 90
WHERE sit_days_allowance IS NULL AND shipment_type = 'HHG'
