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

--update current contract to old dates
update re_contract_years set start_date = start_date - interval '3 years', end_date = end_date - interval '3 years'
 where id = 'fe25a8b9-d6d8-4182-9193-3730c0700925'
   and start_date > '2015-06-01';

update re_contract_years set start_date = start_date - interval '3 years', end_date = end_date - interval '3 years'
 where id = '2ed876b2-603b-4930-a026-b4570cdd8515'
   and start_date > '2016-06-01';
  
update re_contract_years set start_date = start_date - interval '3 years', end_date = end_date - interval '3 years'
 where id = 'f81c12ef-696a-45c8-ac1f-df4583eee359'
   and start_date > '2017-06-01';
  
update re_contract_years set start_date = start_date - interval '3 years', end_date = end_date - interval '3 years'
 where id = '9ac23b92-742a-4d17-9a50-61ae4cf3f3e3'
   and start_date > '2018-06-01';
  
update re_contract_years set start_date = start_date - interval '3 years', end_date = end_date - interval '3 years'
 where id = 'a6457124-b8d5-4b75-95bd-fd402974a043'
   and start_date > '2019-06-01';
  
update re_contract_years set start_date = start_date - interval '3 years', end_date = end_date - interval '3 years'
 where id = '2546c4de-25dd-499d-a6da-61e4de247649'
   and start_date > '2020-06-01';
  
update re_contract_years set start_date = start_date - interval '3 years', end_date = end_date - interval '3 years'
 where id = '53b69fab-c9d6-4900-9998-03b4980519ba'
   and start_date > '2021-06-01';

update re_contract_years set start_date = start_date - interval '3 years', end_date = end_date - interval '3 years'
 where id = '741f66ee-34c6-4388-b6ce-c48b6597b6e3'
   and start_date > '2022-06-01';

INSERT INTO re_contract_years (id,contract_id,"name",start_date,end_date,escalation,escalation_compounded,created_at,updated_at) VALUES
	 ('3b9d9f6f-a341-4e8f-bfbc-5bb8bf3e36c0'::uuid,'070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'Base Period Year 1','2024-06-01','2025-05-31',1.00000,1.00000,'2025-01-21 12:27:36.412','2025-01-21 12:27:36.412'),
	 ('e3d300d2-e3c9-46d0-aec1-a3bb8f52f37c'::uuid,'070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'Base Period Year 2','2025-06-01','2026-05-31',1.02060,1.02060,'2025-01-21 12:27:36.412','2025-01-21 12:27:36.412'),
	 ('f803879d-7cf5-435e-abd7-2a43a9c10213'::uuid,'070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'Base Period Year 3','2026-06-01','2027-05-31',1.01970,1.04071,'2025-01-21 12:27:36.412','2025-01-21 12:27:36.412'),
	 ('6700aeee-8dc7-40cd-b291-19717d270263'::uuid,'070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'Option Period 1','2027-06-01','2028-05-31',1.02140,1.06298,'2025-01-21 12:27:36.412','2025-01-21 12:27:36.412'),
	 ('be1316ba-ecd6-454d-9244-fa48f9fcea88'::uuid,'070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'Option Period 2','2028-06-01','2029-05-31',1.02110,1.08541,'2025-01-21 12:27:36.412','2025-01-21 12:27:36.412'),
	 ('3b2bed81-a9ba-44d2-a087-042cf2a71054'::uuid,'070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'Award Term 1','2029-06-01','2030-05-31',1.01990,1.10701,'2025-01-21 12:27:36.412','2025-01-21 12:27:36.412'),
	 ('ee01ccfc-c1ba-477b-b771-be1accff922d'::uuid,'070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'Award Term 2','2030-06-01','2031-05-31',1.01940,1.12848,'2025-01-21 12:27:36.412','2025-01-21 12:27:36.412'),
	 ('821c967d-5a70-4f9e-873b-21f34366794a'::uuid,'070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'Option Period 3','2031-06-01','2032-05-31',1.02020,1.15128,'2025-01-21 12:27:36.412','2025-01-21 12:27:36.412');

