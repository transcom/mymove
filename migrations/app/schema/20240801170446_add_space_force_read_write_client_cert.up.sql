-- Add space force permissions
ALTER TABLE client_certs
    ADD COLUMN IF NOT EXISTS allow_space_force_orders_read boolean DEFAULT false NOT NULL,
    ADD COLUMN IF NOT EXISTS allow_space_force_orders_write boolean DEFAULT false NOT NULL;

COMMENT ON COLUMN client_certs.allow_space_force_orders_read IS 'Indicates whether or not the cert grants view-only access to Space Force orders';
COMMENT ON COLUMN client_certs.allow_space_force_orders_read IS 'Indicates whether or not the cert grants edit access to Space Force orders';

-- Update all existing client certs with Air Force permissions to have Space Force permissions
UPDATE client_certs
SET allow_space_force_orders_read = TRUE,
    allow_space_force_orders_write = TRUE
WHERE id IN ('190b1e07-eef8-445a-9696-5a2b49ee488d',
             '1e7998d0-3145-4252-b293-7a6a3d52cb32',
             'b8e767d6-fb38-4602-b7f1-6bead16ca7e1',
             'de7605ec-2edd-4252-b176-27f0d3fe4b6f',
             '9928caf4-072c-438a-8a5e-07b213bc1826',
             '03191873-8cc5-4af1-9866-b4bf12842d54',
             '102628a2-b4a2-49a8-a850-cd0e4998f846',
             '730bf94b-f54e-46c3-b125-0a572d885209',
             '9d7fa3b2-8896-4f10-a89c-e161003a3387',
             'c1ef75d3-3c99-4155-b636-518e1b2b5390',
             '1adfca38-3be5-4db0-b5a7-dfc9a739519c',
             '642cf8c6-91e0-490e-8db6-eedf5e583797',
             '5eda0836-2910-4280-8f37-f129e3859f2a',
             '2fbfb740-cc36-4eeb-a6c6-1e800c74b348',
             'debb6c7f-bfdf-4621-ade4-8def4571034a',
             'cf2b7400-0890-400d-85a2-906ce34281f3',
             '913d7488-fded-424e-9bf7-7f14dde3d596');