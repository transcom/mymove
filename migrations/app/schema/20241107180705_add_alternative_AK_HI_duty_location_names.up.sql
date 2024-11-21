SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

-- drop to allow duplicate names for some of the one to many mappings
DROP INDEX IF EXISTS duty_location_names_name_idx;

-- Delete duty locations. It will be moved to duty_location_names and mapped to official duty location.
-- Anaktuvuk, AK 99721
delete from public.duty_locations where id = 'e38e7a71-f29b-497e-a4bd-2635d0c1d3a7';
-- Arctic Vlg, AK 99722
delete from public.duty_locations where id = '0543eb47-c657-4c9d-a18d-7c471820da3f';
-- Brevig Msn, AK 99785
delete from public.duty_locations where id = 'c4ed9f97-e146-4e3c-a0ba-68100c475f74';
-- Chignik Lagn, AK 99565
delete from public.duty_locations where id = '6ee652b0-588e-43a5-9aef-b2d57fb095b1';
-- Cooper Lndg, AK 99572
delete from public.duty_locations where id = 'dbaffeb2-5f52-423c-adb6-1778ee348952';
-- Delta Jct, AK 99731
delete from public.duty_locations where id = 'ffd5c239-5ae4-4a10-a299-aac021c239f5';
-- Delta Jct, AK 99737
delete from public.duty_locations where id = '3fb5efd3-61e7-4a39-8946-0fa9c638e7ca';
-- Denali Park, AK 99755
delete from public.duty_locations where id = '21afbdcc-4eb4-43a8-be16-6b095e15c4b7';
-- JBER, AK 99506
delete from public.duty_locations where id = 'b77ed728-43cd-48f9-9ce0-b634118cca2d';
-- Ft Richardson, AK 99505
delete from public.duty_locations where id = 'dc1a0648-fc2a-48a2-9274-475d4d774ec5';
-- JBER, AK 99505
delete from public.duty_locations where id = '2bb0c7af-e66f-450d-aa95-1b3fbeeeceec';
-- Lk Minchumina, AK 99757
delete from public.duty_locations where id = '2a3e2d6d-f42f-4e4f-adae-6ade1bd0e81a';
-- Ltl Diomede, AK 99762
delete from public.duty_locations where id = '24f0dc19-1942-4ba4-bd6b-8dcb5c7b1a02';
-- Manley Spgs, AK 99756
delete from public.duty_locations where id = 'be9c43f8-be66-45a0-861d-5138833642e3';
-- Mountain Vlg, AK 99632
delete from public.duty_locations where id = '09805568-5a29-496a-85cd-d5a9b1a88014';
-- Prt Alexander, AK 99836
delete from public.duty_locations where id = '978c037a-e5be-4074-a12d-120234fbd04d';
-- Russian Msn, AK 99657
delete from public.duty_locations where id = '137c41e5-6f60-4675-ac30-909a2ad112d4';
-- St George Is, AK 99591
delete from public.duty_locations where id = 'e386f82e-ea13-446b-b10c-58cb2757d5d4';
-- St Paul, AK 99660
delete from public.duty_locations where id = 'd1d6c91d-aafc-4260-94ba-00251969f096';
-- St Paul Isle, AK 99591
delete from public.duty_locations where id = '5838da02-9b7f-4135-a809-434eddf93478';
-- St Paul Isle, AK 99660
delete from public.duty_locations where id = '54d6938d-31c2-4c1f-b473-3e9d4ef908a4';
-- Stevens Vlg, AK 99774
delete from public.duty_locations where id = 'b514f591-be16-4f1e-bb4d-b770b2a3921a';
-- Tenakee Spgs, AK 99841
delete from public.duty_locations where id = '8c0e8b6f-5c13-4587-891a-31b14d3c3620';
-- White Mtn, AK 99784
delete from public.duty_locations where id = '14849ebf-1004-4535-af3a-67f01bb57ea2';
-- Camp Smith, HI 96861
delete from public.duty_locations where id = '933fdf6f-e664-4ad0-aa68-f139f4950994';
-- Hi Natl Park, HI 96718
delete from public.duty_locations where id = '62333230-cb17-405b-b11e-97bb08284d28';
-- JBPHH, HI 96853
delete from public.duty_locations where id = '7487a873-86a2-49b3-8d93-9658ab303c21';
-- Joint Base Pearl Hbr Hickam, HI 96853
delete from public.duty_locations where id = '8d4e3c05-9f1f-4b23-a26c-3c971967fde8';
-- JBPHH, HI 96860
delete from public.duty_locations where id = '106cf7a8-8019-4c44-96aa-e678781e9d8c';
-- Joint Base Pearl Hbr Hickam, HI 96860
delete from public.duty_locations where id = 'd6684c82-07da-48e5-a43a-c3a4b18cc103';
-- MCBH K Bay, HI 96863
delete from public.duty_locations where id = 'd62e79c7-39c2-434e-baf9-27b93e06326b';
-- Schofield, HI 96857
delete from public.duty_locations where id = '6b511cd1-042a-4616-8137-b3a50feac1b1';
-- Tripler AMC, HI 96859
delete from public.duty_locations where id = 'dcdae05f-4ef8-4ed3-9af4-dc5acc8dd550';
-- Wheeler Army Airfield, HI 96854
delete from public.duty_locations where id = 'd2186bf5-c704-4236-848f-e47f9aa56de5';

-- Delete addresses associated to deleted duty locations
delete from public.addresses where id = '6cafd1f0-211d-4ad1-8837-d6dde20710ed';
delete from public.addresses where id = '37672ec6-8646-42c0-b0ea-99cb30f1385c';
delete from public.addresses where id = '8f424337-3d0a-489b-8947-b0aa36345a2f';
delete from public.addresses where id = '84fcc19b-901b-4fe4-91a0-de3de505f845';
delete from public.addresses where id = 'f17dee0e-9580-4a83-abc1-76ed4330b395';
delete from public.addresses where id = '80e27e0c-6f7e-42b0-87d7-76cf64316a6c';
delete from public.addresses where id = 'fc771d16-88cb-4141-8e21-229ccb0759cd';
delete from public.addresses where id = '79f2623f-71e6-4a4b-a52b-b18cbe905b94';
delete from public.addresses where id = 'a6c902b6-2156-402b-b592-a79b3ceb09cb';
delete from public.addresses where id = 'ffd60e16-098e-4f3a-80a6-a35b581268e2';
delete from public.addresses where id = '0db73725-706d-478f-8fc7-e32a37bf1b78';
delete from public.addresses where id = '66a201ea-17dd-4259-a3b2-8d3ddd759f20';
delete from public.addresses where id = '9345f198-4294-4b18-a889-bcfb95d84a7b';
delete from public.addresses where id = 'bafad92f-c054-4370-88b2-8c91aaffed2f';
delete from public.addresses where id = '785624fe-de0a-44d9-8f42-22c5c5b7c24e';
delete from public.addresses where id = '3adedae3-737d-4854-a1c2-2e188f3702b4';
delete from public.addresses where id = 'b0815286-7e97-4184-8d85-650c9fc98f74';
delete from public.addresses where id = 'd58c8365-5e13-4c41-bd9c-2e9d4b7b88fd';
delete from public.addresses where id = '439e8368-ba8e-46c6-98a5-f143b5f5c65b';
delete from public.addresses where id = '5de93ffd-fb91-4134-8b67-1a156b0ea4be';
delete from public.addresses where id = 'aceda6bf-5dfd-4eae-ac08-f56e39b7cda3';
delete from public.addresses where id = '4a30c211-7b2e-4af9-94e6-a7c889715e2e';
delete from public.addresses where id = 'db1aed10-a3ff-458f-8562-c00ab3915634';
delete from public.addresses where id = 'e52b6f0b-35b4-4ef1-9d48-cdc2974899ae';
delete from public.addresses where id = '3cec34ca-30c1-43a9-8f08-af1cbad67848';
delete from public.addresses where id = '7277c530-d662-4728-94b2-3b0768acd000';
delete from public.addresses where id = 'ff80201a-5826-42d6-a037-2017e21d2351';
delete from public.addresses where id = 'a870bf4b-5cb5-4c22-8fd8-2450cefe6223';
delete from public.addresses where id = '85431b39-45b8-48e2-aea9-e05000f95dc5';
delete from public.addresses where id = '68ffbb20-5bf7-4ee1-8a12-dbe128a71788';
delete from public.addresses where id = '0f7767b4-082e-47e0-be51-da39e6c4cb41';
delete from public.addresses where id = 'e9c13c27-a7f2-4848-87be-c60440c4a6cc';
delete from public.addresses where id = '017475dd-edee-41a6-b388-bcab75ef959e';
delete from public.addresses where id = '4d301ae9-cbd2-49dd-a8c2-ab277b9d623b';

-- Insert alternative duty location names and map to official duty location.
COPY public.duty_location_names (id, name, duty_location_id, created_at, updated_at) FROM stdin;
88809570-59e5-4152-b24b-a15595091058	Anaktuvuk	1507ad69-7d74-44a3-8359-0e14f12e0a2b	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
406db356-e0b6-4207-80ba-ba89c8174dff	Arctic Vlg	75e3c0c2-b6a3-4100-9396-918fdbdf1018	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
2d4f1c7d-c45b-4645-8b0a-d7985446827f	Brevig Msn	0f9cc96f-5fb9-40e7-9bd2-51537246a30f	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
65caded0-3538-4950-a22b-800e59b23fab	Chignik Lagn	14502920-775c-4541-86a1-4538ca7764d3	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
957676a2-dd23-4848-b696-5a2403f6a45c	Cooper Lndg	6518cdd7-a8ff-44a1-a726-d0c599030f8e	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
5460d08b-4e80-43df-8109-0e6d118691c7	Delta Jct	2467ac5b-13c0-4233-90f5-2d4c2d149596	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
246c7156-00d9-4e8b-8bce-f0aa695b9fcf	Delta Jct	2f1bf750-2bc9-4431-ad7a-67171a19c00e	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
ed9eb336-2340-44f0-8d6f-c428479f361f	Denali Park	d681547f-2367-4d5a-b45f-a085ca40b916	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
9035901a-da9b-4bdc-9652-d0a5b42f8d78	JBER	39a83615-2e2c-414e-a6d1-f6431279d8e7	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
b2a40180-8870-4eff-b528-f52b85e18b00	Ft Richardson	81a07647-984c-4b09-bc3c-6204d93c7096	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
41d6e3e8-162b-4f61-9fd9-56d43c282639	JBER	81a07647-984c-4b09-bc3c-6204d93c7096	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
7273b107-15fc-4e68-98b7-fe0441c22dc5	Lk Minchumina	af16181e-4d54-4c4a-9c59-774c650b4965	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
4e3ef36e-fea4-4d0c-84b6-f47941182509	Ltl Diomede	25eed791-1ecd-4418-bd3f-0912687d51ce	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
12c93c45-ebe7-4cdf-84c3-f0fc8798396d	Manley Spgs	4467be39-20c7-4b3c-a2aa-92c4adf1ff3b	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
14f84739-6d12-4a00-898a-229fda13001f	Mountain Vlg	fee7e6fc-8f10-46e5-96b1-20219f92bd6b	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
2f945c65-6ba6-4498-bf21-fbda776b9f70	Prt Alexander	95f23779-67f3-46bd-8a7e-4760a39373c7	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
756772e2-f41c-4a2d-a629-39f6dc9161d7	Russian Msn	2ea7f3d5-6551-4686-8413-a505f419f234	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
232457e8-7627-4e94-8d0f-58e27f42ab17	St George Is	a007afc5-362d-4867-b3dd-702ecde9ce4f	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
f2fca271-7d4b-4b4c-9482-8fafce5940f2	St Paul	4a1959e3-ab4e-4a75-8bf8-daa5e0b87a8d	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
8a63e079-d5df-4612-bd2a-4acc8836de65	St Paul Isle	38a49808-5a75-4b02-826c-f509b1b74333	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
a5a32b6e-cc36-4b9b-b478-c8c31c387f08	St Paul Isle	93b4580d-a196-4555-b07b-f9bcc6fb4a41	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
1cfa9836-ab72-4b8f-a458-65e81ea36649	Stevens Vlg	69081a90-e959-4a39-81f4-5d8299b31a97	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
10808880-d69b-406e-8292-cccda1d42986	Tenakee Spgs	fa636bdd-e0a4-425d-99fe-ffb246d7db95	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
a0b387a8-a055-4e2d-97c2-5549ffc3bd74	White Mtn	711ec6f0-1d87-4e0e-b444-bf52b4595050	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
83f1c1c6-9c94-4755-8c38-10ff6804b18f	Camp Smith	0387e3eb-8bb5-4c61-958a-24bf6fceeb79	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
9c308098-7219-4666-89be-57946ff37084	Hi Natl Park	66e8d5a0-9a82-42d2-9bf3-e5ebea944bcd	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
8ceb152f-61f6-410d-bebf-32692b7b56e5	JBPHH	e4c5ee6d-acb4-4324-abda-86e7daf10b32	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
0dd79853-d9ca-4087-9c9d-e7ac0cfda3f5	Joint Base Pearl Hbr Hickam	e4c5ee6d-acb4-4324-abda-86e7daf10b32	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
8d327970-ce28-4050-ae15-22f6e6b0fa5c	JBPHH	2bee231c-db48-4096-a918-816fcbd12d0d	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
2755ded0-1b2e-48c4-aec2-fc4a57321421	Joint Base Pearl Hbr Hickam	2bee231c-db48-4096-a918-816fcbd12d0d	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
383fc3b5-0fc6-4760-ab94-06120561ef65	MCBH K Bay	b0553a92-9daa-4a25-b38d-70e7d3773b2e	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
990d5d04-d01d-41d3-b682-c4e9b2258761	Schofield	8b0dd123-cfad-4f74-af36-a71380f7cc2f	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
a54d9f31-692c-4855-a970-917870594fb6	Tripler AMC	704ba475-e7b5-4e26-8dcf-b73e97f64fdf	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
463a4572-1778-49b1-bd55-559353b9ac2f	Wheeler Army Airfield	3b223bec-58df-4921-ad6d-402798c71976	2024-11-07 21:52:01.047685	2024-11-07 21:52:01.047685
\.
