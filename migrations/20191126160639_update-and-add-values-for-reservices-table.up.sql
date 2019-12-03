UPDATE re_services SET code = 'MS', name = 'Shipment Mgmt. Services', updated_at = now() WHERE id = '1130e612-94eb-49a7-973d-72f33685e551';
UPDATE re_services SET code = 'CS', updated_at = now() WHERE id = '9dc919da-9b66-407b-9f17-05c0f03fcb50';
UPDATE re_services SET code = 'DOP', name = 'Dom. Origin Price', updated_at = now() WHERE id = '2bc3e5cb-adef-46b1-bde9-55570bfdd43e';
UPDATE re_services SET code = 'DOFSIT', name = 'Dom. Origin 1st Day SIT', updated_at = now() WHERE id = '998beda7-e390-4a83-b15e-578a24326937';
UPDATE re_services SET code = 'DOASIT', name = 'Dom. Origin Add''l SIT', updated_at = now() WHERE id = '05eb6ff1-5cf6-4918-b887-8260dda6b9fe';
UPDATE re_services SET code = 'DOPSIT', name = 'Dom. Origin SIT Pickup', updated_at = now() WHERE id = 'd1a4f062-0ca3-4387-8f8e-3dd20493d0b7';
UPDATE re_services SET code = 'DUPK', updated_at = now() WHERE id = '15f01bc1-0754-4341-8e0f-25c8f04d5a77';
UPDATE re_services SET code = 'DUCRT', updated_at = now() WHERE id = 'fc14935b-ebd3-4df3-940b-f30e71b6a56c';
UPDATE re_services SET code = 'DOSHUT', name = 'Dom. Origin Shuttle Service', updated_at = now() WHERE id = 'd979e8af-501a-44bb-8532-2799753a5810';
UPDATE re_services SET name = 'Int''l. O->O Shipping & LH', updated_at = now() WHERE id = '56bb94cd-f160-4239-a028-b31ffc641eb7';
UPDATE re_services SET name = 'Int''l. O->O UB', updated_at = now() WHERE id = '133bc44a-2c0f-4ff9-a61d-9f2dad349d14';
UPDATE re_services SET name = 'Int''l. C->O Shipping & LH', updated_at = now() WHERE id = '07051352-4715-49b5-88e7-045b7541919d';
UPDATE re_services SET name = 'Int''l. C->O UB', updated_at = now() WHERE id = '16949697-6171-47aa-bcd4-479995cc5206';
UPDATE re_services SET name = 'Int''l. O->C Shipping & LH', updated_at = now() WHERE id = 'd0bb2cae-838a-4fc7-8efc-f7c6ad57431d';
UPDATE re_services SET name = 'Int''l. O->C UB', updated_at = now() WHERE id = '7f14f357-f2ae-42b9-b834-5ab24c2cd2af';
UPDATE re_services SET name = 'Int''l. HHG Pack', updated_at = now() WHERE id = '67ba1eaf-6ffd-49de-9a69-497be7789877';
UPDATE re_services SET code = 'IHUPK', name = 'Int''l. HHG Unpack', updated_at = now() WHERE id = '56e91c2d-015d-4243-9657-3ed34867abaa';
UPDATE re_services SET name = 'Int''l. UB Pack', updated_at = now() WHERE id = 'ae84d292-f885-4138-86e2-b451855ffbf2';
UPDATE re_services SET code = 'IUBUPK', name = 'Int''l. UB Unpack', updated_at = now() WHERE id = 'f2739142-97d1-40f3-a8f4-6a9daf390806';
UPDATE re_services SET code = 'IOFSIT', name = 'Int''l. Origin 1st Day SIT', updated_at = now() WHERE id = 'b488bf85-ea5e-49c8-ba5c-e2fa278ac806';
UPDATE re_services SET code = 'IOASIT', name = 'Int''l. Origin Add''l Day SIT', updated_at = now() WHERE id = 'bd424e45-397b-4766-9712-de4ae3a2da36';
UPDATE re_services SET code = 'IOPSIT', name = 'Int''l. Origin SIT Pickup', updated_at = now() WHERE id = '6f4f6e31-0675-4051-b659-89832259f390';
UPDATE re_services SET code = 'IUCRT', updated_at = now() WHERE id = '4132416b-b1aa-42e7-98f2-0ac0a03e8a31';
UPDATE re_services SET code = 'IOSHUT', name = 'Int''l. Origin Shuttle Service', updated_at = now() WHERE id = '624a97c5-dfbf-4da9-a6e9-526b4f95af8d';
UPDATE re_services SET name = 'NonStd. HHG', updated_at = now() WHERE id = '7e1c4f99-0054-4fac-a302-1a07a1daf58e';
UPDATE re_services SET name = 'NonStd. UB', updated_at = now() WHERE id = 'a68fa051-b09d-43f3-8290-f56fc89a4fe8';

INSERT INTO re_services
(id, code, name, created_at, updated_at)
VALUES
('50f1179a-3b72-4fa1-a951-fe5bcc70bd14', 'DDP', 'Dom. Destination Price', now(), now()),
('d0561c49-e1a9-40b8-a739-3e639a9d77af', 'DDFSIT', 'Dom. Destination 1st Day SIT', now(), now()),
('a0ead168-7469-4cb6-bc5b-2ebef5a38f92', 'DDASIT', 'Dom. Destination Add''l SIT', now(), now()),
('5c80f3b5-548e-4077-9b8e-8d0390e73668', 'DDDSIT', 'Dom. Destination SIT Delivery', now(), now()),
('556663e3-675a-4b06-8da3-e4f1e9a9d3cd', 'DDSHUT', 'Dom. Destination Shuttle Service', now(), now()),
('bd6064ca-e780-4ab4-a37b-0ae98eebb244', 'IDFSIT', 'Int''l. Destination 1st Day SIT', now(), now()),
('806c6d59-57ff-4a3f-9518-ebf29ba9cb10', 'IDASIT', 'Int''l. Destination Add''l Day SIT', now(), now()),
('28389ee1-56cf-400c-aa52-1501ecdd7c69', 'IDDSIT', 'Int''l. Destination SIT Delivery', now(), now()),
('22fc07ed-be15-4f50-b941-cbd38153b378', 'IDSHUT', 'Int''l. Destination Shuttle Service', now(), now()),
('4780b30c-e846-437a-b39a-c499a6b09872', 'FSC', 'Fuel Surcharge', now(), now()),
('dbd3a39a-6bb9-42da-b81a-9229df7019cf', 'DMHF', 'Dom. Mobile Home Factor', now(), now()),
('0e45b6f5-f2f5-4235-94e4-7b4cb899eb5d', 'DBTF', 'Dom. Tow Away Boat Factor', now(), now()),
('2471cc2d-6ed5-4ecc-9d43-db9711c8645b', 'DBHF', 'Dom. Haul Away Boat Factor', now(), now()),
('20998cfd-bfc7-410b-a3c5-d709ead4f94e', 'IBTF', 'Int’l. Tow Away Boat Factor', now(), now()),
('387b9654-5685-4ac9-b213-81962be9c145', 'IBHF', 'Int’l. Haul Away Boat Factor', now(), now()),
('3cc83af7-ecb9-4b33-bbc6-ff1459f001e2', 'DNPKF', 'Dom. NTS Packing Factor', now(), now()),
('874cb86a-bc39-4f57-a614-53ee3fcacf14', 'INPKF', 'Int’l. NTS Packing Factor', now(), now());