-- Update shipment factor for boat tow away and haul away shipments
UPDATE re_shipment_type_prices AS rstp SET factor = 35.33 FROM re_services AS rs WHERE rs.id = rstp.service_id AND rs.code = 'DBTF';
UPDATE re_shipment_type_prices AS rstp SET factor = 45.77 FROM re_services AS rs WHERE rs.id = rstp.service_id AND rs.code = 'DBHF';