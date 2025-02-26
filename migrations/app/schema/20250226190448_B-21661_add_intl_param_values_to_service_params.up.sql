-- B-21661 Elizabeth Perkins Add IHUBUPK, IHUBPK, and UBP params

-- inserting params for IHUBUPK
INSERT INTO service_params (id,service_id,service_item_param_key_id,created_at,updated_at,is_optional) VALUES
     ('34f3df3f-7beb-4947-80b6-4e861614588d'::uuid,'f2739142-97d1-40f3-a8f4-6a9daf390806','a1d31d35-c87d-4a7d-b0b8-8b2646b96e43','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false), -- ContractCode
     ('e50a24bf-aa10-4dd3-9307-f510b60b9eba'::uuid,'f2739142-97d1-40f3-a8f4-6a9daf390806','597bb77e-0ce7-4ba2-9624-24300962625f','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false), -- PerUnitCents
     ('3a7ac1a1-4fba-4d6f-8e46-99aca0d47f8a'::uuid,'f2739142-97d1-40f3-a8f4-6a9daf390806','95ee2e21-b232-4d74-9ec5-218564a8a8b9','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false) on conflict do nothing; -- IsPeak

-- inserting params for IHUBPK
INSERT INTO service_params (id,service_id,service_item_param_key_id,created_at,updated_at,is_optional) VALUES
     ('72f74e74-3206-4658-912e-800050740714'::uuid,'ae84d292-f885-4138-86e2-b451855ffbf2','a1d31d35-c87d-4a7d-b0b8-8b2646b96e43','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false), -- ContractCode
     ('467b3fd0-fb05-4cf0-8740-4dd8952f10cf'::uuid,'ae84d292-f885-4138-86e2-b451855ffbf2','597bb77e-0ce7-4ba2-9624-24300962625f','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false), -- PerUnitCents
     ('4619a245-8cf4-4ddb-b901-a764b3fd2c83'::uuid,'ae84d292-f885-4138-86e2-b451855ffbf2','95ee2e21-b232-4d74-9ec5-218564a8a8b9','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false) on conflict do nothing; -- IsPeak

-- inserting params for UBP
INSERT INTO service_params (id,service_id,service_item_param_key_id,created_at,updated_at,is_optional) VALUES
     ('ebe4a0c7-c8fa-4653-9d42-c0c383689540'::uuid,'fafad6c2-6037-4a95-af2f-d9861ba7db8e','a1d31d35-c87d-4a7d-b0b8-8b2646b96e43','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false), -- ContractCode
     ('9e2cecb9-b948-4090-bfd3-c5160f1bb7ea'::uuid,'fafad6c2-6037-4a95-af2f-d9861ba7db8e','597bb77e-0ce7-4ba2-9624-24300962625f','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false), -- PerUnitCents
     ('fd9d312b-ab24-4b92-8c1c-5ed53548cbe0'::uuid,'fafad6c2-6037-4a95-af2f-d9861ba7db8e','95ee2e21-b232-4d74-9ec5-218564a8a8b9','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false), -- IsPeak
     ('48e36c9c-d856-42c6-a27c-652ad105c6b9'::uuid,'fafad6c2-6037-4a95-af2f-d9861ba7db8e','5335e243-ab5b-4906-b84f-bd8c35ba64b3','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false), -- ReferenceDate
     ('696b8068-53cd-45c7-a3fb-da608c86daa6'::uuid,'fafad6c2-6037-4a95-af2f-d9861ba7db8e','b9739817-6408-4829-8719-1e26f8a9ceb3','2024-12-26 15:55:50.041957','2024-12-26 15:55:50.041957',false) on conflict do nothing; -- WeightBilled