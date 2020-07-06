-- TODO: This is not written to be a zero-downtime migration, but prod doesn't have the
--  ability to write to customers yet so there should be no prod downtime. This will
--  need more thought if we want to try to make these changes without any downtime in
--  any environment.

-- Lock the customers and service_members tables so we don't have changes while
-- doing this migration
LOCK TABLE customers, service_members IN SHARE MODE;

-- Temporary field to hold the customer_id for a row that will be excluded on
-- conflict. See the insertion below for more details.
ALTER TABLE service_members ADD COLUMN old_customer_id uuid;

-- Copy current customers over to service_members
-- If a customer points to a user_id that already exists in service_members,
-- then we exclude this customer, and temporarily store its `id` in the
-- `old_customer_id` field in service_members on the row with the matching user_id.
-- This is needed so that we can then update the move_orders table `customer_id`
-- (the one that was excluded due to the user_id conflict) to point to the
-- corresponding service_member id.
INSERT INTO service_members
(id, created_at, updated_at, user_id, edipi, first_name, last_name, affiliation, personal_email, telephone,
 residential_address_id, requires_access_code)
SELECT id,
       created_at,
       updated_at,
       user_id,
       dod_id,
       first_name,
       last_name,
       agency,
       email,
       phone,
       current_address_id,
       false -- required field, so set to false (matching what's in dev/exp/prod)
       -- ignoring destination_address_id from customers as it's never set except from testdatagen
FROM customers
ON CONFLICT (user_id) DO UPDATE SET old_customer_id = EXCLUDED.id;

-- Migrate move_orders.customer_id FK reference from customers to service_members
ALTER TABLE move_orders
    DROP CONSTRAINT move_orders_customer_id_fkey; -- FK name not given when FK created, so assuming default.

UPDATE move_orders
SET customer_id = service_members.id
FROM
  service_members
WHERE
  move_orders.customer_id = service_members.old_customer_id;

-- Even though we are now pointing to service_members, we are not renaming the
-- customer_id key because this PR does not make any API changes. Moreover, we
-- plan to eventually rename the service_members table to customers.
ALTER TABLE move_orders
    ADD CONSTRAINT move_orders_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES service_members (id);

-- Drop the field that we temporarily created for this migration.
ALTER TABLE service_members DROP COLUMN old_customer_id;
