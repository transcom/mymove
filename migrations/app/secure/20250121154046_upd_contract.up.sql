-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.

--run re_intl_other_prices_prod.sql to populate prod data first

--add prod contract id
INSERT INTO public.re_contracts
(id, code, "name", created_at, updated_at)
VALUES('070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid, 'HTC711-22-D-R002', 'Global HHG Relocation Services', now(), now());

--update all pricing for new contract id - re_intl_prices
update re_domestic_accessorial_prices set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_domestic_linehaul_prices set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_domestic_other_prices set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_domestic_service_area_prices set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_domestic_service_areas set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_intl_accessorial_prices set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_intl_other_prices set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_intl_prices set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_rate_areas set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_shipment_type_prices set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_zip3s set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_zip5_rate_areas set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
update re_contract_years set contract_id = '070f7c82-fad0-4ae8-9a83-5de87a56472e' where contract_id = '51393fa4-b31c-40fe-bedf-b692703c46eb';

--delete old contract id
delete from re_contracts where id = '51393fa4-b31c-40fe-bedf-b692703c46eb';
