ALTER TABLE client_certs
    ADD COLUMN IF NOT EXISTS allow_space_force_orders_read boolean DEFAULT false NOT NULL,
    ADD COLUMN IF NOT EXISTS allow_space_force_orders_write boolean DEFAULT false NOT NULL;

COMMENT ON COLUMN client_certs.allow_space_force_orders_read IS 'Indicates whether or not the cert grants view-only access to Space Force orders';
COMMENT ON COLUMN client_certs.allow_space_force_orders_read IS 'Indicates whether or not the cert grants edit access to Space Force orders';