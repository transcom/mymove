-- TODO: This is not written to be a zero-downtime migration, but prod doesn't have the
--  ability to write to customers yet so there should be no prod downtime. This will
--  need more thought if we want to try to make these changes without any downtime in
--  any environment.

-- Lock the customers table so we don't have any customer changes while doing this
LOCK TABLE customers IN SHARE MODE;

-- Copy current customers over to service_members
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
       -- ignoring destination_address_id from customers as its never set except from testdatagen
FROM customers;

-- Migrate move_orders.customer_id FK reference from customers to service_members
ALTER TABLE move_orders
    DROP CONSTRAINT move_orders_customer_id_fkey; -- FK name not given when FK created, so assuming default.
ALTER TABLE move_orders
    ADD CONSTRAINT move_orders_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES service_members (id);

-- Delete all records from customers
-- TODO: Change this to a DROP TABLE once remaining customers refs are removed from authghc.go
DELETE
FROM customers;
