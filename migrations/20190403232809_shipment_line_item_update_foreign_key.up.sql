ALTER TABLE "shipment_line_items" DROP CONSTRAINT "shipment_line_items_addresses_id_fk";

ALTER TABLE "shipment_line_items" ADD FOREIGN KEY ("address_id") REFERENCES "addresses" ("id");