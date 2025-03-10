CREATE TABLE IF NOT EXISTS re_fsc_multipliers (
    id          uuid             NOT NULL PRIMARY KEY,
    low_weight  int              NOT NULL,
    high_weight int              NOT NULL,
    multiplier  decimal          NOT NULL,
    created_at  timestamp        NOT NULL DEFAULT NOW(),
    updated_at  timestamp        NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE re_fsc_multipliers IS 'Stores data needed to calculate FSC';
COMMENT ON COLUMN re_fsc_multipliers.low_weight IS 'The lowest weight permitted for a shipment';
COMMENT ON COLUMN re_fsc_multipliers.high_weight IS 'The highest weight permitted for a shipment';
COMMENT ON COLUMN re_fsc_multipliers.multiplier IS 'The decimal multiplier used to calculate the FSC';

INSERT INTO re_fsc_multipliers (id,low_weight,high_weight,multiplier,created_at,updated_at) VALUES
	 ('e8053bda-e19e-4343-b858-04d691f1438b',0,5000,0.000417,'2024-11-19 12:59:06.276892-06','2024-11-19 12:59:06.276892-06'),
	 ('5f2d402e-9c41-4034-bfea-46e699e2ed95',5001,10000,0.0006255,'2024-11-19 12:59:06.276892-06','2024-11-19 12:59:06.276892-06'),
	 ('0faa2887-bf9e-4fb5-9d96-42bd79d963bd',10001,24000,0.000834,'2024-11-19 12:59:06.276892-06','2024-11-19 12:59:06.276892-06'),
	 ('9bb374f7-12a2-4c51-8391-a691406b4c2c',24001,99999999,0.00139,'2024-11-19 12:59:06.276892-06','2024-11-19 12:59:06.276892-06');
