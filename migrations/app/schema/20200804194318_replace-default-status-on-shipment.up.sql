ALTER TABLE mto_shipments ALTER COLUMN status DROP DEFAULT;

ALTER TABLE mto_shipments ALTER COLUMN status SET DEFAULT 'DRAFT';