-- deliveries table will be used to handle entire and partial deliveries
-- foreign keys to the shipment and service_items tables
CREATE TABLE deliveries
(
	id uuid PRIMARY KEY NOT NULL,
	mto_shipment_id uuid REFERENCES mto_shipments(id) NOT NULL,
    mto_service_item_id uuid REFERENCES mto_service_items(id) NOT NULL,
    delivery_weight_pounds integer NOT NULL,
	created_at timestamp WITH TIME ZONE NOT NULL,
	updated_at timestamp WITH TIME ZONE NOT NULL
);

-- adding indexes on these columns to accelerate database queries that involve these columns
CREATE INDEX ON deliveries (mto_shipment_id);
CREATE INDEX ON deliveries (mto_service_item_id);