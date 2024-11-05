CREATE TABLE IF NOT EXISTS re_service_items
(id 			uuid		NOT NULL,
service_id      uuid        NOT NULL
    CONSTRAINT fk_re_service_items_service_id REFERENCES re_services (id),
shipment_type   mto_shipment_type	NOT NULL,
market_code     market_code_enum  	NOT NULL,
is_auto_approved   bool 		NOT NULL,
created_at		timestamp	NOT NULL DEFAULT NOW(),
updated_at		timestamp	NOT NULL DEFAULT NOW(),
CONSTRAINT re_service_items_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_service_items UNIQUE (service_id, shipment_type, market_code)
);

COMMENT ON TABLE re_service_items IS 'Associates service items to market_code and shipment_type.';
COMMENT ON COLUMN re_service_items.service_id IS 'The associated id for the service item';
COMMENT ON COLUMN re_service_items.shipment_type IS 'The type of shipment from mto_shipments table.';
COMMENT ON COLUMN re_service_items.market_code IS 'Market code indicator. i for international and d for domestic.';
COMMENT ON COLUMN re_service_items.is_auto_approved IS 'Set to true if the service item is automatically approved when assigned to the shipment.';

--Create enum type for service_location
CREATE TYPE service_location_enum AS ENUM (
    'O', --Origin
    'D', --Destination
    'B'); --Both

--Add service_location to re_services
ALTER TABLE re_services ADD COLUMN IF NOT EXISTS service_location service_location_enum;
COMMENT ON COLUMN re_services.service_location IS 'Specifies where this service item may be billed at (Origin, Destination or Both)';

--Add service_location to mto_service_items
ALTER TABLE mto_service_items ADD COLUMN IF NOT EXISTS service_location service_location_enum;
COMMENT ON COLUMN mto_service_items.service_location IS 'Specifies the responsible billing location for this service item (Origin or Destination)';

--Add new international service items
INSERT INTO public.re_services
(id, code, name, created_at, updated_at, priority, service_location)
VALUES('fafad6c2-6037-4a95-af2f-d9861ba7db8e', 'UBP', 'International UB', now(), now(), 99, 'B');

INSERT INTO public.re_services
(id, code, name, created_at, updated_at, priority, service_location)
VALUES('9f3d551a-0725-430e-897e-80ee9add3ae9', 'ISLH', 'International Shipping & Linehaul', now(), now(), 99, 'B');

INSERT INTO public.re_services
(id, code, name, created_at, updated_at, priority, service_location)
VALUES('f75758d8-2fcd-40ba-9432-3ff3032a71d1', 'POEFSC', 'International POE Fuel Surcharge', now(), now(), 99, 'O');

INSERT INTO public.re_services
(id, code, name, created_at, updated_at, priority, service_location)
VALUES('388115e8-abe9-441d-96cf-a39f24baa0a3', 'PODFSC', 'International POD Fuel Surcharge', now(), now(), 99, 'D');

INSERT INTO public.re_services
(id, code, name, created_at, updated_at, priority, service_location)
VALUES('81e29d0c-02a6-4a7a-be02-554deb3ee49e', 'IOSFSC', 'International Origin SIT Fuel Surcharge', now(), now(), 99, 'O');

INSERT INTO public.re_services
(id, code, name, created_at, updated_at, priority, service_location)
VALUES('690a5fc1-0ea5-4554-8294-a367b5daefa9', 'IDSFSC', 'International Destination SIT Fuel Surcharge', now(), now(), 99, 'D');

--update existing re_services
update re_services set service_location = 'O' where code in ('IUBPK','IHPK','IOASIT','IOFSIT','IOPSIT','IOSHUT','ICRT');
update re_services set service_location = 'D' where code in ('IUBUPK','IHUPK','IDASIT','IDFSIT','IDDSIT','IDSHUT','IUCRT');

INSERT INTO re_service_items (id,service_id,shipment_type,market_code,is_auto_approved,created_at,updated_at) VALUES
	 ('2933a8a0-5d89-4fcb-99f0-60893edf8c8a'::uuid,'f75758d8-2fcd-40ba-9432-3ff3032a71d1'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('6b2dedf0-389f-41b7-908e-22a8ecaa320c'::uuid,'388115e8-abe9-441d-96cf-a39f24baa0a3'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('df4f7cbe-10d7-4042-a380-657835d408e1'::uuid,'fafad6c2-6037-4a95-af2f-d9861ba7db8e'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('78bde6ed-9fac-442b-aa74-dfba1a537336'::uuid,'ae84d292-f885-4138-86e2-b451855ffbf2'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('c4089abd-3225-4986-a5d9-37b46ae5b645'::uuid,'f2739142-97d1-40f3-a8f4-6a9daf390806'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('ba2dc894-163d-4db9-81e3-0259dc66a70f'::uuid,'f75758d8-2fcd-40ba-9432-3ff3032a71d1'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('c361227b-1aae-49ec-bcbb-30f663d70445'::uuid,'388115e8-abe9-441d-96cf-a39f24baa0a3'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('567f3d6d-c1a4-46fe-9424-4b60f1dba939'::uuid,'9f3d551a-0725-430e-897e-80ee9add3ae9'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('2008074d-e447-4619-b2bd-fe2a38a2759e'::uuid,'67ba1eaf-6ffd-49de-9a69-497be7789877'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('daf7b541-f9d1-4b1a-998b-d530ca48da3e'::uuid,'56e91c2d-015d-4243-9657-3ed34867abaa'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('32be29e1-5269-4fd0-85b4-03270cbe6a65'::uuid,'86203d72-7f7c-49ff-82f0-5b95e4958f60'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now());
INSERT INTO re_service_items (id,service_id,shipment_type,market_code,is_auto_approved,created_at,updated_at) VALUES
	 ('0ee7f930-6421-44b0-8820-48bffc2b2605'::uuid,'806c6d59-57ff-4a3f-9518-ebf29ba9cb10'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('e555e8a5-23bf-40ea-9ee8-01479c971716'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('2705b054-3a05-492d-9262-3a6d4125c268'::uuid,'bd6064ca-e780-4ab4-a37b-0ae98eebb244'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('cbed2e49-799e-42c0-bd86-4204c3ca049d'::uuid,'22fc07ed-be15-4f50-b941-cbd38153b378'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('60be16ad-8d92-4e51-b565-7ef7904f4fe9'::uuid,'bd424e45-397b-4766-9712-de4ae3a2da36'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('c1000d89-ae2b-47bc-b3b1-63c2c24fc076'::uuid,'b488bf85-ea5e-49c8-ba5c-e2fa278ac806'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('b09f964d-fee7-4bf9-84ea-c06571b7ef51'::uuid,'6f4f6e31-0675-4051-b659-89832259f390'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('ab8981ca-b3f0-4421-8753-95860ca40cfb'::uuid,'624a97c5-dfbf-4da9-a6e9-526b4f95af8d'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('4d09972f-5fb3-44a2-bd0e-4b21f716e222'::uuid,'4132416b-b1aa-42e7-98f2-0ac0a03e8a31'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now());
INSERT INTO re_service_items (id,service_id,shipment_type,market_code,is_auto_approved,created_at,updated_at) VALUES
	 ('28487c9e-a91b-4870-b2e3-fbefb6a324e6'::uuid,'86203d72-7f7c-49ff-82f0-5b95e4958f60'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('6bfeb792-7f08-4b33-904d-f837dc69c097'::uuid,'806c6d59-57ff-4a3f-9518-ebf29ba9cb10'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('f896c537-14e8-4612-833c-794d0ada17b8'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('3d5cc86e-bd3a-49e0-a041-63e0625e690f'::uuid,'bd6064ca-e780-4ab4-a37b-0ae98eebb244'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('6d4cc684-55ba-4308-864f-76846ec66cf1'::uuid,'22fc07ed-be15-4f50-b941-cbd38153b378'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('3c3e99f3-bf63-48b9-9101-a9f8b91afea3'::uuid,'bd424e45-397b-4766-9712-de4ae3a2da36'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('047fcb9f-4b90-454b-aca4-3561b46e4de9'::uuid,'b488bf85-ea5e-49c8-ba5c-e2fa278ac806'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now());
INSERT INTO re_service_items (id,service_id,shipment_type,market_code,is_auto_approved,created_at,updated_at) VALUES
	 ('d6c8ce58-bc00-499f-9cad-9047210b8746'::uuid,'6f4f6e31-0675-4051-b659-89832259f390'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('6799f4b5-c4ba-48e4-a4bd-c6adb71b7041'::uuid,'624a97c5-dfbf-4da9-a6e9-526b4f95af8d'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('faa6dc9d-e120-4d26-a54b-b84a2b2f6d7c'::uuid,'4132416b-b1aa-42e7-98f2-0ac0a03e8a31'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",false,now(),now()),
	 ('cdc719ba-8f57-4b11-80a0-f416413b0064'::uuid,'81e29d0c-02a6-4a7a-be02-554deb3ee49e'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('b56b974c-ce60-4efc-bf79-553d8094d18a'::uuid,'690a5fc1-0ea5-4554-8294-a367b5daefa9'::uuid,'HHG'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('a5b07f41-b279-44c0-9ecd-241bd8a5c2eb'::uuid,'81e29d0c-02a6-4a7a-be02-554deb3ee49e'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now()),
	 ('39cd5798-f0b1-49a6-b9fe-0f4f4039b80b'::uuid,'690a5fc1-0ea5-4554-8294-a367b5daefa9'::uuid,'UNACCOMPANIED_BAGGAGE'::public."mto_shipment_type",'i'::public."market_code_enum",true,now(),now());

--Point existing UB service pricing to new UBP service
update re_intl_prices
set service_id = (select id from re_services where code = 'UBP')
where service_id in (select id from re_services where code in ('IOCUB','ICOUB','IOOUB','NSTUB'));

update re_intl_prices
set service_id = (select id from re_services where code = 'ISLH')
where service_id in (select id from re_services where code in ('IOCLH','ICOLH','IOOLH','NSTH'));

--Remove UB stuff from service_params as we don't need it.
delete from service_params sp
where service_id in (select id from re_services where code in ('IOCUB','ICOUB','IOOUB','NSTUB','IOCLH','ICOLH','IOOLH','NSTH'));

--Remove obsolete UB service itmes from re_services
delete from re_services rs
where code in ('IOCUB','ICOUB','IOOUB','NSTUB','IOCLH','ICOLH','IOOLH','NSTH');