ALTER TABLE roles
    ADD COLUMN role_name VARCHAR(255);

UPDATE
    roles
SET
    role_name = 'Customer'
WHERE
    role_type = 'customer';

UPDATE
    roles
SET
    role_name = 'Transporation Ordering Officer'
WHERE
    role_type = 'transportation_ordering_officer';

UPDATE
    roles
SET
    role_name = 'Transporation Invoicing Officer'
WHERE
    role_type = 'transportation_invoicing_officer';

UPDATE
    roles
SET
    role_name = 'Contracting Officer'
WHERE
    role_type = 'contracting_officer';

UPDATE
    roles
SET
    role_name = 'PPP Office User'
WHERE
    role_type = 'ppm_office_users';

ALTER TABLE roles
    ALTER COLUMN role_name SET NOT NULL;

