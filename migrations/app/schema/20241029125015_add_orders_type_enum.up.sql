-- Add enum values for order.orders_type

UPDATE orders
SET orders_type = 'PERMANENT_CHANGE_OF_STATION'
WHERE orders_type = 'VARIOUS' OR orders_type = 'NTS' OR orders_type = 'DEPENDENT_TRAVEL' OR orders_type = 'GHC';

CREATE TYPE orders_type AS ENUM (
'PERMANENT_CHANGE_OF_STATION',
'LOCAL_MOVE',
'RETIREMENT',
'SEPARATION',
'WOUNDED_WARRIOR',
'BLUEBARK',
'SAFETY',
'TEMPORARY_DUTY');

COMMENT ON TYPE orders_type IS 'The type of orders.';
COMMENT ON COLUMN orders.orders_type IS 'MilMove supports 8 orders types: Permanent change of station (PCS), local move, retirement, separation, wounded warrior, bluebark, safety, and temporary duty (TDY).
In general, the moving process starts with the job/travel orders a customer receives from their service. In the orders, information describing rank, the duration of job/training, and their assigned location will determine if their entire dependent family can come, what the customer is allowed to bring, and how those items will arrive to their new location.';

ALTER TABLE orders
  ALTER COLUMN orders_type TYPE orders_type using orders_type::orders_type;