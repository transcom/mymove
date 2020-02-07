CREATE TYPE order_types AS ENUM (
    'GHC',
    'NTS'
    );

ALTER TABLE move_orders
	ADD COLUMN order_type order_types,
	ADD COLUMN order_type_detail text,
	ADD COLUMN date_issued date,
	ADD COLUMN report_by_date date;
