-- Add external_crate to mto_service_items
ALTER TABLE mto_service_items ADD COLUMN IF NOT EXISTS external_crate bool NULL;
COMMENT ON COLUMN mto_service_items.external_crate IS 'Boolean value indicating whether the international crate is externally crated.';

-- removing 'International crating - standalone' (ICRTSA) from the tables
delete from service_params sp
where service_id in (select id from re_services where code in ('ICRTSA'));

delete from re_intl_accessorial_prices reiap
where service_id in (select id from re_services where code in ('ICRTSA'));

delete from re_services rs
where code in ('ICRTSA');