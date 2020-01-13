ALTER TABLE re_shipment_type_prices RENAME COLUMN factor_hundredths TO factor;
ALTER TABLE re_shipment_type_prices ALTER COLUMN factor TYPE numeric(6,5);
