-- Update the mobile home factor in re_shipment_type_prices
UPDATE re_shipment_type_prices AS rstp SET factor = 33.51 FROM re_services AS rs WHERE rs.id = rstp.service_id AND rs.code = 'DMHF';