UPDATE client_certs SET allow_army_orders_read = true, allow_army_orders_write = true WHERE subject SIMILAR TO '%[Aa]rmy%';
UPDATE client_certs SET allow_marine_corps_orders_read = true, allow_marine_corps_orders_write = true WHERE subject SIMILAR TO '%(usmc|Marine Corps)%';
