CREATE TABLE permissions (
	id uuid PRIMARY KEY,
	permission_type text,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL
);

CREATE TABLE role_permissions (
	permission_id uuid NOT NULL,
	role_id uuid NOT NULL,
	FOREIGN KEY (permission_id) REFERENCES permissions (id),
	FOREIGN KEY (role_id) REFERENCES roles (id)
);

-- Add permissions
INSERT INTO
	permissions
VALUES
	(
		uuid_generate_v4(),
		'edit.shipment',
		now(),
		now()
	),
	(
		uuid_generate_v4(),
		'edit.orders',
		now(),
		now()
	),
	(
		uuid_generate_v4(),
		'edit.max_billable_weight',
		now(),
		now()
	),
	(
		uuid_generate_v4(),
		'edit.financial_review_flag',
		now(),
		now()
	);

-- Assign role permissions
INSERT INTO
	role_permissions (permission_id, role_id)
SELECT
	p.id,
	r.id
FROM
	permissions p,
	roles r
WHERE
	(
		p.permission_type = 'edit.shipment'
		AND (
			r.role_type = 'transportation_ordering_officer'
			OR r.role_type = 'services_counselor' -- OR r.role_type = 'transportation_invoicing_officer'
		)
	)
	OR (
		p.permission_type = 'edit.orders'
		AND (
			r.role_type = 'transportation_ordering_officer'
			OR r.role_type = 'services_counselor' --OR r.role_type = 'transportation_invoicing_officer'
		)
	)
	OR (
		p.permission_type = 'edit.max_billable_weight'
		AND (
			r.role_type = 'transportation_ordering_officer'
			OR r.role_type = 'services_counselor' --OR r.role_type = 'transportation_invoicing_officer'
		)
	)
	OR (
		p.permission_type = 'edit.financial_review_flag'
		AND (
			r.role_type = 'transportation_ordering_officer'
			OR r.role_type = 'services_counselor' --OR r.role_type = 'transportation_invoicing_officer'
		)
	)
