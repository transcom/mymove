ALTER TABLE shipment_line_items RENAME TO shipment_accessorials;

ALTER TABLE shipment_line_items
RENAME COLUMN tariff400ng_item_id TO accessorial_id ;