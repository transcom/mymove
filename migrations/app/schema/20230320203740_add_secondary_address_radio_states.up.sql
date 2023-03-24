alter table mto_shipments
	add column has_secondary_pickup_address bool,
    add column has_secondary_delivery_address bool;
