CREATE TABLE IF NOT EXISTS ports (
    id          uuid             NOT NULL,
    port_code   varchar(4)       NOT NULL,
    port_type   varchar(1)       NOT NULL,
    port_name   varchar(100)     NOT NULL,
    created_at  timestamp        NOT NULL DEFAULT NOW(),
    updated_at  timestamp        NOT NULL DEFAULT NOW(),
    CONSTRAINT  port_pkey        PRIMARY KEY(id),
    CONSTRAINT  unique_port_code UNIQUE (port_code),
    CONSTRAINT  chk_port_type    CHECK (port_type IN ('A', 'S', 'B'))
);
COMMENT ON TABLE ports IS 'Stores ports identification data';
COMMENT ON COLUMN ports.port_code IS 'The 4 digit port code';
COMMENT ON COLUMN ports.port_type IS 'The 1 char port type A, S, or B';
COMMENT ON COLUMN ports.port_name IS 'The name of the port';
ALTER TABLE mto_service_items ADD COLUMN IF NOT EXISTS poe_location_id uuid;
ALTER TABLE mto_service_items ADD CONSTRAINT fk_poe_location_id FOREIGN KEY (poe_location_id) REFERENCES ports (id);
ALTER TABLE mto_service_items ADD COLUMN IF NOT EXISTS pod_location_id uuid;
ALTER TABLE mto_service_items ADD CONSTRAINT fk_pod_location_id FOREIGN KEY (pod_location_id) REFERENCES ports (id);
COMMENT ON COLUMN mto_service_items.poe_location_id IS 'Stores the POE location id for port of embarkation';
COMMENT ON COLUMN mto_service_items.pod_location_id IS 'Stores the POD location id for port of debarkation';