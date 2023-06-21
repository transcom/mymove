ALTER TABLE orders
ADD supply_and_services_cost_estimate TEXT,
ADD packing_and_shipping_instructions TEXT,
ADD method_of_payment TEXT,
ADD naics TEXT;

-- Column comments
COMMENT ON COLUMN orders.supply_and_services_cost_estimate IS 'Context for what the costs are based on.';
COMMENT ON COLUMN orders.packing_and_shipping_instructions IS 'Context for where instructions can be found.';
COMMENT ON COLUMN orders.method_of_payment IS 'Context regarding how the payment will occur.';
COMMENT ON COLUMN orders.naics IS 'North American Industry Classification System Code.';
