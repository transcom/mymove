ALTER TABLE shipment_accessorials RENAME TO shipment_line_items;

ALTER TABLE shipment_line_items
RENAME COLUMN accessorial_id TO tariff400ng_item_id;