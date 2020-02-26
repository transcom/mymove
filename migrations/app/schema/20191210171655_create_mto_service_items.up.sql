CREATE TABLE mto_service_items
(
	id uuid PRIMARY KEY NOT NULL,
	move_task_order_id uuid REFERENCES move_task_orders,
	mto_shipment_id uuid REFERENCES mto_shipments,
	re_service_id uuid REFERENCES re_services,
	meta_id uuid NOT NULL,
	meta_type text NOT NULL,
	created_at timestamp WITH TIME ZONE NOT NULL,
	updated_at timestamp WITH TIME ZONE NOT NULL
);

CREATE INDEX ON mto_service_items (move_task_order_id);
CREATE INDEX ON mto_service_items (mto_shipment_id);
CREATE INDEX ON mto_service_items (re_service_id);
