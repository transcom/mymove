--insert USMC gbloc
insert into jppso_regions values ('85b324ae-1a0d-4f7d-971f-ea509dcc73d7', 'USMC','USMC',now(),now());

--insert USMC AORs
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('a8eb35e2-275e-490b-9945-1971b954b958'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'e418be11-2b6f-4714-b026-e293528c50bd'::uuid,'MARINES',NULL,true,now(),now()),
	 ('ada1b48d-d2e0-481a-8a2e-a265a824647d'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'b41c5636-96dd-4f0d-a18e-eebb17f97ea5'::uuid,'MARINES',NULL,true,now(),now()),
	 ('588af482-7cd7-42ea-8e05-49dce645ecbe'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'02f32a6a-0338-4545-8437-059862892d2c'::uuid,'MARINES',NULL,true,now(),now()),
	 ('1ff3bed7-bbf3-432d-9da3-d76264d72913'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'0f128476-d52a-418c-8ba0-c8bfd1c32629'::uuid,'MARINES',NULL,true,now(),now()),
	 ('0853a854-98b2-4363-a2b1-db14c44dde2f'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'081f84c3-17ec-4ff6-97ce-d5c44a8e4a28'::uuid,'MARINES',NULL,true,now(),now()),
	 ('21c1ba40-2533-4196-9eb5-6ffddff3a794'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'ce7cdd91-e323-43a6-a604-daaa6bf8be06'::uuid,'MARINES',NULL,true,now(),now()),
	 ('19ca4073-8736-453e-bb0d-9b13e3b557b0'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'655050d4-711a-43b7-b06d-828ec990c35e'::uuid,'MARINES',NULL,true,now(),now()),
	 ('a9c2131e-09c3-480d-ba3e-0c144de18aa5'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'455a34af-30a8-4a98-a62b-6f40fd7f047b'::uuid,'MARINES',NULL,true,now(),now()),
	 ('649e72f8-cac8-483f-a9ed-c9659e37545b'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'1bf1daee-51bb-4c28-aac8-a126ef596486'::uuid,'MARINES',NULL,true,now(),now()),
	 ('7173f871-f948-4eed-86ae-5e977b16c426'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'15b1d852-0dde-4e1b-b3c6-e08bbc714db3'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('4f626e2a-cad0-4b4d-baa8-3101275bac23'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'d63de5bb-379b-4f3b-b4c1-554b9746d311'::uuid,'MARINES',NULL,true,now(),now()),
	 ('de59c604-9119-48fc-bf3f-883de17b7ee6'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'95d142f8-b50a-4108-b5e2-fbd3b7022d3b'::uuid,'MARINES',NULL,true,now(),now()),
	 ('1cec884c-9e34-42fe-8887-1e8b8fa1cd2e'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'922bb7da-d0e9-431c-a493-95a861dca929'::uuid,'MARINES',NULL,true,now(),now()),
	 ('0f24b453-e3bc-47ec-86b5-8d8937f65504'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'79aef0ca-9065-4c0e-a134-2889b250cc38'::uuid,'MARINES',NULL,true,now(),now()),
	 ('2d55560a-7d0a-474f-a0f9-b31c5be8d80e'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'1f181025-e410-41ac-935b-4b5147353f84'::uuid,'MARINES',NULL,true,now(),now()),
	 ('abcc37f6-9209-4639-8fa1-c8d5f6e4e77d'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'d9ede620-c3f1-4f8d-b60b-eb93664382f7'::uuid,'MARINES',NULL,true,now(),now()),
	 ('3ac990d2-5df4-4889-a943-2710f818e75a'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'bd27700b-3094-46cc-807b-f18472cbfaf0'::uuid,'MARINES',NULL,true,now(),now()),
	 ('f3084090-680f-4656-947a-eb2e773e4076'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'2f4a1689-ee65-45fa-9d0c-debd9344f8b9'::uuid,'MARINES',NULL,true,now(),now()),
	 ('a35ac50b-09a1-46ef-969e-17569717ee10'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'20c52934-2588-4d8f-b3ed-c046d771f4e9'::uuid,'MARINES',NULL,true,now(),now()),
	 ('107c1479-e6d9-44cb-8342-ac934055074d'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'15fe364b-add5-4ade-bdd9-38fe442616fb'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('843db089-67cb-463d-a255-1d198f4f7aaa'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'e5a1248e-3870-4a4c-9c9d-2056234eb64a'::uuid,'MARINES',NULL,true,now(),now()),
	 ('ac820ac9-d380-4c11-9103-172795658e1f'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'efc92db1-4b06-49ab-a295-a98f3f6c9c04'::uuid,'MARINES',NULL,true,now(),now()),
	 ('dca7ccd2-e438-4f82-8b76-9b61fbf2a593'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'13474ce5-8839-4af7-b975-c3f70ccdab7b'::uuid,'MARINES',NULL,true,now(),now()),
	 ('def1095b-2a5c-4f7c-8889-9fde15a7ec06'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'d3b05c5e-6faa-4fa9-b725-0904f8a4f3d7'::uuid,'MARINES',NULL,true,now(),now()),
	 ('d5fbc738-ee31-4c51-8fd4-bbb8db941dc1'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'e5f849fe-5672-4d0d-881c-078e56eea33d'::uuid,'MARINES',NULL,true,now(),now()),
	 ('ac8c75fe-f637-429d-80fa-913321a65372'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'fca7c60e-fbd9-4885-afdb-ae41c521b560'::uuid,'MARINES',NULL,true,now(),now()),
	 ('c33eb1f8-b0fc-4670-af5e-5bd423eca6e7'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'0d52a0f0-f39c-4d34-9387-2df45a5810d4'::uuid,'MARINES',NULL,true,now(),now()),
	 ('188ea995-b8b7-4ce0-97a9-f553f3b72c2f'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'af008e75-81d5-4211-8204-4964c78e70d9'::uuid,'MARINES',NULL,true,now(),now()),
	 ('4f13c1a6-059b-4aa2-9250-4316e60da2a7'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'dd689e55-af29-4c76-b7e0-c2429ea4833c'::uuid,'MARINES',NULL,true,now(),now()),
	 ('67e1ff6f-b79b-45cf-b250-3d4ec89bebae'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'60e08330-6869-4586-8198-35f7a4ae9ea7'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('3a7c679c-0439-4030-8f7e-f6d8e92720d7'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'f7478e79-dbfe-46c8-b337-1e7c46df79dc'::uuid,'MARINES',NULL,true,now(),now()),
	 ('db91f133-a301-4c87-af9f-6c10584063e6'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'486d7dd4-f51f-4b13-88fc-830b5b15f0a8'::uuid,'MARINES',NULL,true,now(),now()),
	 ('6b83cc22-36a9-470e-9f43-e65ade5e8a66'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'f66edf62-ba5e-47f8-8264-b6dc9e7dd9ba'::uuid,'MARINES',NULL,true,now(),now()),
	 ('3df30221-acd3-4428-890c-3a5ef5296cb1'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'0052c55a-d9d7-46b0-9328-58d3911f61b4'::uuid,'MARINES',NULL,true,now(),now()),
	 ('bae4c9c6-94d8-4ad0-bfc0-7642f5353199'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'16a51fd1-04ed-432a-a8d7-9f17c1de095d'::uuid,'MARINES',NULL,true,now(),now()),
	 ('1cd873ff-d170-4f48-8f5b-5c5146052d68'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'64b4756f-437d-4aa5-a95e-c396f0cafcbb'::uuid,'MARINES',NULL,true,now(),now()),
	 ('eb4d870b-8d66-4135-b294-8992e56ad76f'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'f2de1da5-a737-4493-b3f7-7700944a5b62'::uuid,'MARINES',NULL,true,now(),now()),
	 ('e8c0709c-1f08-4d3e-b848-72ab5b524677'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'7a9c2adb-0562-42c1-a5f1-81710dd19590'::uuid,'MARINES',NULL,true,now(),now()),
	 ('3c53df36-ee6b-4bc0-a5ca-cedc1dc3c32e'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'e713eed6-35a6-456d-bfb2-0e0646078ab8'::uuid,'MARINES',NULL,true,now(),now()),
	 ('7cbc7e6c-4da7-47a1-ac93-171b89dba1e0'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'a4594957-0ec6-4edc-85c2-68871f4f6359'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('d22cb1d2-79b4-45c9-bb6f-8acef54a67b0'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'acd4d2f9-9a49-4a73-be14-3515f19ba0d1'::uuid,'MARINES',NULL,true,now(),now()),
	 ('5b5e7a5a-d027-44f2-9b4f-f25c1a91bc00'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'c8300ab6-519c-4bef-9b81-25e7334773ca'::uuid,'MARINES',NULL,true,now(),now()),
	 ('21307477-912d-40ca-a399-6dfebcc322ea'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'67eca2ce-9d88-44ca-bb17-615f3199415c'::uuid,'MARINES',NULL,true,now(),now()),
	 ('82282efb-7fef-4a6d-a260-a18d2f21fa8d'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'149c3a94-abb1-4af0-aabf-3af019d5e243'::uuid,'MARINES',NULL,true,now(),now()),
	 ('3bdad313-d414-4842-8e6e-3675e20d78eb'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'87c09b47-058e-47ea-9684-37a8ac1f7120'::uuid,'MARINES',NULL,true,now(),now()),
	 ('db80b591-045d-4907-9b77-6694fe34e3ed'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'0ba81e85-f175-435a-a7c2-a22d0c44cc7b'::uuid,'MARINES',NULL,true,now(),now()),
	 ('05970a8a-28aa-454f-a28e-31327a2415dd'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'946256dc-1572-4201-b5ff-464da876f5ff'::uuid,'MARINES',NULL,true,now(),now()),
	 ('c1b3e0be-0463-4dfc-ab22-016037b41a05'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'5d381518-6beb-422e-8c57-4682b87ff1fe'::uuid,'MARINES',NULL,true,now(),now()),
	 ('1a92f4b0-4060-4f1b-9b45-b3fc47c5f08d'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'b07fb701-2c9b-4847-a459-76b1f47aa872'::uuid,'MARINES',NULL,true,now(),now()),
	 ('73794c53-a915-41af-b93d-5ada2e174409'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'b6c51584-f3e8-4c9e-9bdf-f7cac0433319'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('c03994c2-9e2a-4888-a5c2-c81ee05eba31'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'6097d01d-78ef-40d5-8699-7f8a8e48f4e7'::uuid,'MARINES',NULL,true,now(),now()),
	 ('ae6a55fd-2171-425a-9716-56d3fb9452e3'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'e48a3602-7fdd-4527-a8a7-f244fb331228'::uuid,'MARINES',NULL,true,now(),now()),
	 ('f052acf4-7d4a-4061-a2f4-ba19ff17ec4d'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'df47ef29-e902-4bd3-a32d-88d6187399b3'::uuid,'MARINES',NULL,true,now(),now()),
	 ('10b47428-bef7-4cfe-8564-2ccb654d514a'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'0f4eddf9-b727-4725-848e-3d9329553823'::uuid,'MARINES',NULL,true,now(),now()),
	 ('4a547501-6aae-4886-be5c-e5a0fad05441'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'d16ebe80-2195-48cf-ba61-b5efb8212754'::uuid,'MARINES',NULL,true,now(),now()),
	 ('d99cbe3a-84a0-4a2c-8e05-41ce066570ea'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'093268d0-597b-40bd-882a-8f385480bc68'::uuid,'MARINES',NULL,true,now(),now()),
	 ('ad09f13c-9dd9-468b-ab92-ad3e2d77c905'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'edb8b022-9534-44e7-87ea-e93f5451f057'::uuid,'MARINES',NULL,true,now(),now()),
	 ('362482c6-7770-49f4-86c9-89e2990e5345'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'4c62bb35-4d2a-4995-9c4f-8c665f2f6d3e'::uuid,'MARINES',NULL,true,now(),now()),
	 ('45994e2f-1aa4-4aa1-8631-a347a4463bc2'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'58ccfcc5-ded6-4f91-8cb7-8687bc56c4c6'::uuid,'MARINES',NULL,true,now(),now()),
	 ('b4fd2bf2-05b2-4429-afcb-ebbae512fd2b'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'dbc839c6-7b56-45b0-ab36-a0e77b0d538c'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('5aa5f596-a9bb-4f37-af0b-6e0cfbfbc711'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'c77e434a-1bf9-44c6-92aa-d377d72d1d44'::uuid,'MARINES',NULL,true,now(),now()),
	 ('c7bc79d4-525e-496d-93bd-2ea76b32baf4'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'89533190-e613-4c36-99df-7ee3871bb071'::uuid,'MARINES',NULL,true,now(),now()),
	 ('5bc99ccb-7010-4d2f-99a0-9277eda982ba'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'c86cb50a-99f6-41ed-8e9d-f096bd5cadca'::uuid,'MARINES',NULL,true,now(),now()),
	 ('dc97dd55-213b-4212-8c1a-aaee7b92fc55'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'b65bd7c4-6251-4682-aa07-77f761c363af'::uuid,'MARINES',NULL,true,now(),now()),
	 ('3a4daddc-a377-440f-bb5f-d33024374e3e'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'88df8618-798a-4a12-9b14-a7267a6cfd7f'::uuid,'MARINES',NULL,true,now(),now()),
	 ('e729c30a-0bec-424a-919a-45cbe31998e9'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'15894c4f-eb1a-4b6b-8b3f-e3abc4eeee8d'::uuid,'MARINES',NULL,true,now(),now()),
	 ('74b6df4f-08ea-48fb-9628-15c3f47f3a27'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'83fd3eb7-6269-46e2-84e1-f6180f40b9e8'::uuid,'MARINES',NULL,true,now(),now()),
	 ('643112a7-71a9-41d1-97f6-92aee969478a'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'b0f95e07-e3d1-4403-a402-50be9a542cf9'::uuid,'MARINES',NULL,true,now(),now()),
	 ('73b35029-077b-4d30-8f5a-34c3785f6e96'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'cef18ff2-8f48-47c0-8062-a40f9dec641c'::uuid,'MARINES',NULL,true,now(),now()),
	 ('9d8862ec-4358-497a-8154-65a83c676261'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'15eefe71-55ae-40b0-8bb5-f1952fcf45c8'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('8d5591d4-fcc4-44a9-babd-b575672ad6a9'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'d3c5d0a7-0477-44f6-8279-f6b9fa7b3436'::uuid,'MARINES',NULL,true,now(),now()),
	 ('c1057fc7-1c8b-44e2-a175-131ff0d7429f'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'8b89e7d7-cb36-40bd-ba14-699f1e4b1806'::uuid,'MARINES',NULL,true,now(),now()),
	 ('ae7949e2-9504-4874-92be-98460e8126da'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'5b833661-8fda-4f2a-9f2b-e1f4c21749a4'::uuid,'MARINES',NULL,true,now(),now()),
	 ('acf2cacd-0d4f-4de1-afdc-2eb1e72eeb80'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'38bed277-7990-460d-ad27-a69559638f42'::uuid,'MARINES',NULL,true,now(),now()),
	 ('c66c3d59-def9-4e08-8bee-5b2bacfc5cfd'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'67081e66-7e84-4bac-9c31-aa3e45bbba23'::uuid,'MARINES',NULL,true,now(),now()),
	 ('1f8cd63a-8dfe-4357-9ac5-2bef71f9d564'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'99e3f894-595f-403a-83f5-f7e035dd1f20'::uuid,'MARINES',NULL,true,now(),now()),
	 ('3e8030da-8316-4330-abea-35452e39fa61'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'1c0a5988-492f-4bc2-8409-9a143a494248'::uuid,'MARINES',NULL,true,now(),now()),
	 ('ce3c740a-b0a2-4a55-9abe-c2426fb8d821'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'63e10b54-5348-4362-932e-6613a8db4d42'::uuid,'MARINES',NULL,true,now(),now()),
	 ('9e913f22-287f-4e42-9544-46d9c6741db7'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'437e9930-0625-4448-938d-ba93a0a98ab5'::uuid,'MARINES',NULL,true,now(),now()),
	 ('d0ae21fe-07c0-40ee-9aa3-877d6d6a6bb9'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'87d3cb47-1026-4f42-bff6-899ff0fa7660'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('6f8d7c90-7682-41b7-b5cc-9d4aeb2e22ef'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'0115586e-be8c-4808-b65a-d417fad19238'::uuid,'MARINES',NULL,true,now(),now()),
	 ('cf7f39c0-8f6b-4598-85db-f465241e66f4'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'6638a77d-4b91-48f6-8241-c87fbdddacd1'::uuid,'MARINES',NULL,true,now(),now()),
	 ('296a8951-19b1-4868-9937-aef61bb73106'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'97c16bc3-e174-410b-9a09-a0db31420dbc'::uuid,'MARINES',NULL,true,now(),now()),
	 ('cdcec10e-5300-443f-9aae-7d8ce07142b0'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'3c6c5e35-1281-45fc-83ee-e6d656e155b6'::uuid,'MARINES',NULL,true,now(),now()),
	 ('e638701f-f8fc-48c0-b2e0-db134b9ece1f'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'6af248be-e5a8-49e9-a9e2-516748279ab5'::uuid,'MARINES',NULL,true,now(),now()),
	 ('ed4f4905-b1cf-4e57-beac-bc0a2d167c71'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'3580abe3-84da-4b46-af7b-d4379e6cff46'::uuid,'MARINES',NULL,true,now(),now()),
	 ('2a8b16c9-99f8-43ca-a759-c4884f8f7b24'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'f3bb397e-04e6-4b37-9bf3-b3ebab79a9b6'::uuid,'MARINES',NULL,true,now(),now()),
	 ('70530105-a4ab-4af3-ba02-9e4cf81237fa'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'f9bfe297-4ee0-4f76-a4bd-64b3a514af5d'::uuid,'MARINES',NULL,true,now(),now()),
	 ('61e17c04-ed7b-4da9-816a-1b6343086d94'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'70f11a71-667b-422d-ae37-8f25539f7782'::uuid,'MARINES',NULL,true,now(),now()),
	 ('19d5158f-96a6-48ef-8dd9-b831c582c9c4'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'baf66a9a-3c8a-49e7-83ea-841b9960e184'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('91043072-657f-4b2a-b5d1-42d8f6a7ba38'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'e96d46b1-3ab7-4b29-b798-ec3b728dd6a1'::uuid,'MARINES',NULL,true,now(),now()),
	 ('a3dc835e-3989-4476-b560-006745a884bc'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'e4171b5b-d26c-43b3-b41c-68b43bcbb079'::uuid,'MARINES',NULL,true,now(),now()),
	 ('0cf1dc14-0c9a-4262-9cab-db73c64f6e36'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'a0651bec-1258-4e36-9c76-155247e42c0a'::uuid,'MARINES',NULL,true,now(),now()),
	 ('e03b53cb-386f-4e1e-ac93-cd1a9260c6b4'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'737a8e63-af19-4902-a4b5-8f80e2268e4b'::uuid,'MARINES',NULL,true,now(),now()),
	 ('8b3f1142-f6d4-4c2d-9e75-9ab3568742f7'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'870df26e-2c50-4512-aa2e-61094cfbc3e1'::uuid,'MARINES',NULL,true,now(),now()),
	 ('5f5383b8-f29e-4b9f-8798-99575440c888'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'09dc6547-d346-40c3-93fb-fb7be1fa1b3e'::uuid,'MARINES',NULL,true,now(),now()),
	 ('1f7a198c-5f62-4f8c-bf91-19d7cdef9bae'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'43d7d081-bb32-4544-84f1-419fe0cb76e1'::uuid,'MARINES',NULL,true,now(),now()),
	 ('56c17bab-a124-4710-a047-0d67d30a9610'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'d7e57942-9c83-4138-baa0-70e8b5f08598'::uuid,'MARINES',NULL,true,now(),now()),
	 ('aad1440d-acaf-4ba8-9564-da97ab9ba651'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'c2f06691-5989-41a3-a848-19c9f0fec5df'::uuid,'MARINES',NULL,true,now(),now()),
	 ('6cc24224-2298-4023-888f-5e624e585171'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'5263a9ed-ff4d-42cc-91d5-dbdefeef54d1'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('88cd8184-7e3a-48cf-bb72-a9e1c3666cc4'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'54ab3a49-6d78-4922-bac1-94a722b9859a'::uuid,'MARINES',NULL,true,now(),now()),
	 ('11fdac1c-7ae9-49ea-bee4-90d4582f7c6d'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'3b6fe2e9-9116-4a70-8eef-96205205b0e3'::uuid,'MARINES',NULL,true,now(),now()),
	 ('83c88ffb-5182-42b3-93f4-635556d8caaf'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'ba0017c8-5d48-4efe-8802-361d2f2bc16d'::uuid,'MARINES',NULL,true,now(),now()),
	 ('802f7a5f-8ce3-4c40-975a-83cf0ec502fc'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'af65234a-8577-4f6d-a346-7d486e963287'::uuid,'MARINES',NULL,true,now(),now()),
	 ('5197eed4-f6ae-4f70-b640-b67714b73f87'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'25f62695-ba9b-40c3-a7e4-0c078731123d'::uuid,'MARINES',NULL,true,now(),now()),
	 ('50d675ab-4064-4edb-8f1d-5739b0318ed9'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'b0422b13-4afe-443e-901f-fe652cde23a4'::uuid,'MARINES',NULL,true,now(),now()),
	 ('a3733567-2b57-4f12-9390-967e04bc1453'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'ecbc1d89-9cf6-4f52-b453-3ba473a0ff4e'::uuid,'MARINES',NULL,true,now(),now()),
	 ('6a321494-cdd8-4372-90a1-6a6c67f4e220'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'4e85cb86-e9dd-4c7c-9677-e0a327ac895c'::uuid,'MARINES',NULL,true,now(),now()),
	 ('99d37d3c-be31-4d3e-9c5e-c9d816a46014'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'7f9dd10d-100e-4252-9786-706349f456ca'::uuid,'MARINES',NULL,true,now(),now()),
	 ('e9b4c884-ad7f-4c1a-8815-05353834f5c3'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'5451f8d7-60c5-4e22-bbf6-d9af8e6ace54'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('cfaa046e-6be5-4178-a451-d368317ecb86'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'9e9cba85-7c39-4836-809d-70b54baf392e'::uuid,'MARINES',NULL,true,now(),now()),
	 ('9b509efb-12d9-42aa-a85a-ffb026866b56'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'f7470914-cf48-43be-b431-a3ca2fe5b290'::uuid,'MARINES',NULL,true,now(),now()),
	 ('7e54d995-b7bd-457c-bd97-0fd76891402e'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'5a0d3cc1-b866-4bde-b67f-78d565facf3e'::uuid,'MARINES',NULL,true,now(),now()),
	 ('1357aadd-b420-42c4-8a39-bab8027fa910'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'e629b95a-ec5b-4fc4-897f-e0d1050e1ec6'::uuid,'MARINES',NULL,true,now(),now()),
	 ('56fbbf0f-e819-41dd-802d-4e677aecd1c9'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'d68a626f-935a-4eb1-ba9b-6829feeff91c'::uuid,'MARINES',NULL,true,now(),now()),
	 ('bfc2dc53-896f-4f29-92c7-9a7a392b22f2'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'ed496f6e-fd34-48d1-8586-63a1d305c49c'::uuid,'MARINES',NULL,true,now(),now()),
	 ('d015e8b1-86ea-489f-979d-458ae35ae8d6'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'153b62b2-b1b8-4b9d-afa5-53df4150aba4'::uuid,'MARINES',NULL,true,now(),now()),
	 ('1f32e712-8b5e-4ae4-b409-c5c92337aed8'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'0020441b-bc0c-436e-be05-b997ca6a853c'::uuid,'MARINES',NULL,true,now(),now()),
	 ('a8a85e00-e657-41a2-8f32-84bdd9c92ec8'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'a8ae9bb9-e9ac-49b4-81dc-336b8a0dcb54'::uuid,'MARINES',NULL,true,now(),now()),
	 ('4ad1a57f-0e9e-4405-9a08-0ffa211fc8ce'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'1aa43046-8d6b-4249-88dc-b259d86c0cb8'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('dd765820-ffa5-4673-a347-fbe3464cd2d8'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'84b50723-95fc-41d1-8115-a734c7e53f66'::uuid,'MARINES',NULL,true,now(),now()),
	 ('6eff7586-59bd-4638-809b-5cc346646dc9'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'63dc3f78-235e-4b1c-b1db-459d7f5ae25f'::uuid,'MARINES',NULL,true,now(),now()),
	 ('428b0d7a-3848-4882-a5cc-80d5ae3500d6'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'4bd4c579-c163-4b1b-925a-d852d5a12642'::uuid,'MARINES',NULL,true,now(),now()),
	 ('9f8311f3-e191-4383-8fb8-2b58cd545dd4'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'ce7f1fc9-5b94-43cb-b398-b31cb6350d6a'::uuid,'MARINES',NULL,true,now(),now()),
	 ('81e3e6fe-7db3-49ff-b18c-8b078e3d129e'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'4a6a3260-e1ec-4d78-8c07-d89f1405ca16'::uuid,'MARINES',NULL,true,now(),now()),
	 ('4cf1c40f-60f0-4a47-a351-471720ba0fd3'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'658274b5-720c-4a70-ba9d-249614d85ebc'::uuid,'MARINES',NULL,true,now(),now()),
	 ('0754978b-d11e-4f2c-b59c-3252d2735b26'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'41f0736c-6d26-4e93-9668-860e9f0e222a'::uuid,'MARINES',NULL,true,now(),now()),
	 ('c9b8305d-2e16-46a7-9b7c-b3edeb6f8e93'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'c2a8e8c3-dddc-4c0f-a23a-0b4e2c20af0d'::uuid,'MARINES',NULL,true,now(),now()),
	 ('a8768e6d-1a6d-449a-9f2e-2e198dcd6e00'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'0f1e1c87-0497-4ee2-970d-21ac2d2155db'::uuid,'MARINES',NULL,true,now(),now()),
	 ('cf292129-d543-4632-9cb3-b074279e42be'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'e40b411d-55f0-4470-83d0-0bbe11fa77dd'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('3f07cedf-ad90-465e-95a5-ce44a2f088b8'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'2632b4e5-c6cb-4e64-8924-0b7e4b1115ec'::uuid,'MARINES',NULL,true,now(),now()),
	 ('42a2d93b-9dea-4c63-b0a7-c39364aacf75'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'1336deb9-5c87-409d-8051-4ab9f211eb29'::uuid,'MARINES',NULL,true,now(),now()),
	 ('64fe0d14-98f7-4f73-9aa3-a19617b2d8c3'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'7a6d3b5b-81a6-4db5-b2ab-ecbfd6bd7941'::uuid,'MARINES',NULL,true,now(),now()),
	 ('3ab04740-9f13-47f8-a80e-b63ab5b67590'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'91b69254-5976-4839-a31d-972e9958d9cf'::uuid,'MARINES',NULL,true,now(),now()),
	 ('3be8e483-6bce-4a7d-a3bd-fc1485e79818'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'166f3629-79b9-451a-90a3-c43680929a2f'::uuid,'MARINES',NULL,true,now(),now()),
	 ('ddf69dcb-345e-47c1-a585-ce24d0854de5'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'4fe05eb4-1b1c-4d4a-a185-0b039ac64835'::uuid,'MARINES',NULL,true,now(),now()),
	 ('fe9e365f-98a5-4658-a58d-5f8279ff3e5a'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'ca8aecda-4642-45c7-96ed-309c35c4b78f'::uuid,'MARINES',NULL,true,now(),now()),
	 ('2051a441-f4d0-4b6e-8614-74761de505e6'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'e4bc9404-5466-4a41-993e-09474266afc3'::uuid,'MARINES',NULL,true,now(),now()),
	 ('6797daed-f002-431c-829b-dab7c1b16ff2'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'d4a51d90-3945-4ad3-9cba-a18d8d7b34d7'::uuid,'MARINES',NULL,true,now(),now()),
	 ('2398bf11-1986-4914-8e47-6afac423283a'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'a5a60d63-d9a8-4bde-9081-f011784b2d31'::uuid,'MARINES',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('f3a6c247-4c4c-4f45-9162-50307d4711f5'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'93842d74-1f3e-46cd-aca9-9f0dafbd20a1'::uuid,'MARINES',NULL,true,now(),now()),
	 ('0a82b196-7c24-4214-9155-05f0c5c2d7e9'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'1bc0dbda-f0ce-4b76-a551-78dbaaa9e3ec'::uuid,'MARINES',NULL,true,now(),now()),
	 ('af0a2db3-2e2e-4a78-804b-9a8b89b96e12'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'f1a7ef90-cfa6-4e0c-92c3-f8d70c07ba4d'::uuid,'MARINES',NULL,true,now(),now()),
	 ('5205a965-d424-469a-9526-17ef551685e6'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'9c5b4c4d-e05c-42ca-bd77-b61f5d8c7afc'::uuid,'MARINES',NULL,true,now(),now()),
	 ('f263fdb0-933b-42a7-925e-a9852c5804fa'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'27ec2576-78dd-4605-b1a8-0b9ca207fc26'::uuid,'MARINES',NULL,true,now(),now()),
	 ('b03c7ae3-4d6d-445c-b94a-73af723e5226'::uuid,'85b324ae-1a0d-4f7d-971f-ea509dcc73d7'::uuid,'12396ebc-59e9-430a-8475-759a38af6b7a'::uuid,'MARINES',NULL,true,now(),now());


drop view move_to_gbloc;
CREATE OR REPLACE VIEW move_to_gbloc AS
SELECT move_id, gbloc FROM (
  SELECT DISTINCT ON (sh.move_id) sh.move_id, s.affiliation,
    COALESCE(pctg.gbloc, coalesce(pctg_oconus_bos.gbloc, coalesce(pctg_oconus.gbloc, pctg_ppm.gbloc))) AS gbloc
  FROM mto_shipments sh
  JOIN moves m ON sh.move_id = m.id
  JOIN orders o on m.orders_id = o.id
  JOIN service_members s on o.service_member_id = s.id
    LEFT JOIN ( SELECT a.id AS address_id,
           pctg_1.gbloc, pctg_1.postal_code
           FROM addresses a
           JOIN postal_code_to_gblocs pctg_1 ON a.postal_code::text = pctg_1.postal_code::text) pctg ON pctg.address_id = sh.pickup_address_id
    LEFT JOIN ( SELECT ppm.shipment_id,
           pctg_1.gbloc
           FROM ppm_shipments ppm
           JOIN addresses ppm_address ON ppm.pickup_postal_address_id = ppm_address.id
           JOIN postal_code_to_gblocs pctg_1 ON ppm_address.postal_code::text = pctg_1.postal_code::text) pctg_ppm ON pctg_ppm.shipment_id = sh.id
    LEFT JOIN ( SELECT a.id AS address_id,
           cast(jr.code as varchar) AS gbloc, ga.department_indicator
           FROM addresses a
           JOIN re_oconus_rate_areas ora ON a.us_post_region_cities_id = ora.us_post_region_cities_id
           JOIN gbloc_aors ga ON ora.id = ga.oconus_rate_area_id
           JOIN jppso_regions jr ON ga.jppso_regions_id = jr.id
        		) pctg_oconus_bos ON pctg_oconus_bos.address_id = sh.pickup_address_id
          				and case when s.affiliation = 'AIR_FORCE' THEN 'AIR_AND_SPACE_FORCE'
           				         when s.affiliation = 'SPACE_FORCE' THEN 'AIR_AND_SPACE_FORCE'
           				         else s.affiliation
          				    end = pctg_oconus_bos.department_indicator
    LEFT JOIN ( SELECT a.id AS address_id,
           cast(pctg_1.code as varchar) AS gbloc, ga.department_indicator
           FROM addresses a
           JOIN re_oconus_rate_areas ora ON a.us_post_region_cities_id = ora.us_post_region_cities_id
           JOIN gbloc_aors ga ON ora.id = ga.oconus_rate_area_id
           JOIN jppso_regions pctg_1 ON ga.jppso_regions_id = pctg_1.id
         		) pctg_oconus ON pctg_oconus.address_id = sh.pickup_address_id and pctg_oconus.department_indicator is null
     WHERE sh.deleted_at IS NULL
     ORDER BY sh.move_id, sh.created_at) as m;


-- used for the destination queue
CREATE OR REPLACE VIEW move_to_dest_gbloc
AS
SELECT distinct move_id, gbloc FROM (
  SELECT sh.move_id, s.affiliation,
    COALESCE(case when s.affiliation = 'MARINES' then 'USMC' else pctg.gbloc end, coalesce(pctg_oconus_bos.gbloc, coalesce(pctg_oconus.gbloc, pctg_ppm.gbloc))) AS gbloc
  FROM mto_shipments sh
  JOIN moves m ON sh.move_id = m.id
  JOIN orders o on m.orders_id = o.id
  JOIN service_members s on o.service_member_id = s.id
    LEFT JOIN ( SELECT a.id AS address_id,
           pctg_1.gbloc, pctg_1.postal_code
           FROM addresses a
           JOIN postal_code_to_gblocs pctg_1 ON a.postal_code::text = pctg_1.postal_code::text) pctg ON pctg.address_id = sh.destination_address_id
    LEFT JOIN ( SELECT ppm.shipment_id,
           pctg_1.gbloc
           FROM ppm_shipments ppm
           JOIN addresses ppm_address ON ppm.destination_postal_address_id = ppm_address.id
           JOIN postal_code_to_gblocs pctg_1 ON ppm_address.postal_code::text = pctg_1.postal_code::text) pctg_ppm ON pctg_ppm.shipment_id = sh.id
    LEFT JOIN ( SELECT a.id AS address_id,
           cast(jr.code as varchar) AS gbloc, ga.department_indicator
           FROM addresses a
           JOIN re_oconus_rate_areas ora ON a.us_post_region_cities_id = ora.us_post_region_cities_id
           JOIN gbloc_aors ga ON ora.id = ga.oconus_rate_area_id
           JOIN jppso_regions jr ON ga.jppso_regions_id = jr.id
        		) pctg_oconus_bos ON pctg_oconus_bos.address_id = sh.destination_address_id
          				and case when s.affiliation = 'AIR_FORCE' THEN 'AIR_AND_SPACE_FORCE'
           				         when s.affiliation = 'SPACE_FORCE' THEN 'AIR_AND_SPACE_FORCE'
           				         else s.affiliation
          				    end = pctg_oconus_bos.department_indicator
    LEFT JOIN ( SELECT a.id AS address_id,
           cast(pctg_1.code as varchar) AS gbloc, ga.department_indicator
           FROM addresses a
           JOIN re_oconus_rate_areas ora ON a.us_post_region_cities_id = ora.us_post_region_cities_id
           JOIN gbloc_aors ga ON ora.id = ga.oconus_rate_area_id
           JOIN jppso_regions pctg_1 ON ga.jppso_regions_id = pctg_1.id
         		) pctg_oconus ON pctg_oconus.address_id = sh.destination_address_id and pctg_oconus.department_indicator is null
     WHERE sh.deleted_at IS NULL) as m;

-- database function that returns a list of moves that have destination requests
-- this includes shipment address update requests, destination SIT, & destination shuttle
CREATE OR REPLACE FUNCTION get_destination_queue(
    user_gbloc TEXT DEFAULT NULL,
    customer_name TEXT DEFAULT NULL,
    edipi TEXT DEFAULT NULL,
    emplid TEXT DEFAULT NULL,
    m_status TEXT[] DEFAULT NULL,
    move_code TEXT DEFAULT NULL,
    requested_move_date TIMESTAMP DEFAULT NULL,
    date_submitted TIMESTAMP DEFAULT NULL,
    branch TEXT DEFAULT NULL,
    origin_duty_location TEXT DEFAULT NULL,
    counseling_office TEXT DEFAULT NULL,
    too_assigned_user TEXT DEFAULT NULL,
    page INTEGER DEFAULT 1,
    per_page INTEGER DEFAULT 20,
    sort TEXT DEFAULT NULL,
    sort_direction TEXT DEFAULT NULL
)
RETURNS TABLE (
    id UUID,
    show BOOLEAN,
    locator TEXT,
    submitted_at TIMESTAMP WITH TIME ZONE,
    orders_id UUID,
    status TEXT,
    locked_by UUID,
    too_assigned_id UUID,
    counseling_transportation_office_id UUID,
    orders JSONB,
    mto_shipments JSONB,
    counseling_transportation_office JSONB,
    too_assigned JSONB,
    total_count BIGINT
) AS $$
DECLARE
    sql_query TEXT;
    offset_value INTEGER;
    sort_column TEXT;
    sort_order TEXT;
BEGIN
    IF page < 1 THEN
        page := 1;
    END IF;

    IF per_page < 1 THEN
        per_page := 20;
    END IF;

    -- OFFSET for pagination
    offset_value := (page - 1) * per_page;

    sql_query := '
        SELECT
            moves.id AS id,
            moves.show AS show,
            moves.locator::TEXT AS locator,
            moves.submitted_at::TIMESTAMP WITH TIME ZONE AS submitted_at,
            moves.orders_id AS orders_id,
            moves.status::TEXT AS status,
            moves.locked_by AS locked_by,
            moves.too_assigned_id AS too_assigned_id,
            moves.counseling_transportation_office_id AS counseling_transportation_office_id,
            json_build_object(
                ''id'', orders.id,
                ''origin_duty_location_gbloc'', orders.gbloc,
                ''service_member'', json_build_object(
                    ''id'', service_members.id,
                    ''first_name'', service_members.first_name,
                    ''last_name'', service_members.last_name,
                    ''edipi'', service_members.edipi,
                    ''emplid'', service_members.emplid,
                    ''affiliation'', service_members.affiliation
                ),
                ''origin_duty_location'', json_build_object(
                    ''name'', origin_duty_locations.name
                )
            )::JSONB AS orders,
            COALESCE(
                (
                    SELECT json_agg(
                        json_build_object(
                            ''id'', ms.id,
                            ''shipment_type'', ms.shipment_type,
                            ''status'', ms.status,
                            ''requested_pickup_date'', TO_CHAR(ms.requested_pickup_date, ''YYYY-MM-DD"T00:00:00Z"''),
                            ''scheduled_pickup_date'', TO_CHAR(ms.scheduled_pickup_date, ''YYYY-MM-DD"T00:00:00Z"''),
                            ''approved_date'', TO_CHAR(ms.approved_date, ''YYYY-MM-DD"T00:00:00Z"''),
                            ''prime_estimated_weight'', ms.prime_estimated_weight
                        )
                    )
                    FROM (
                        SELECT DISTINCT ON (mto_shipments.id) mto_shipments.*
                        FROM mto_shipments
                        WHERE mto_shipments.move_id = moves.id
                    ) AS ms
                ),
                ''[]''
            )::JSONB AS mto_shipments,
            json_build_object(
                ''name'', counseling_offices.name
            )::JSONB AS counseling_transportation_office,
            json_build_object(
                ''first_name'', too_user.first_name,
                ''last_name'', too_user.last_name
            )::JSONB AS too_assigned,
            COUNT(*) OVER() AS total_count
        FROM moves
        JOIN orders ON moves.orders_id = orders.id
        JOIN mto_shipments ON mto_shipments.move_id = moves.id
        LEFT JOIN ppm_shipments ON ppm_shipments.shipment_id = mto_shipments.id
        JOIN mto_service_items ON mto_shipments.id = mto_service_items.mto_shipment_id
        JOIN re_services ON mto_service_items.re_service_id = re_services.id
        JOIN service_members ON orders.service_member_id = service_members.id
        JOIN duty_locations AS new_duty_locations ON orders.new_duty_location_id = new_duty_locations.id
        JOIN duty_locations AS origin_duty_locations ON orders.origin_duty_location_id = origin_duty_locations.id
        LEFT JOIN office_users AS too_user ON moves.too_assigned_id = too_user.id
        LEFT JOIN office_users AS locked_user ON moves.locked_by = locked_user.id
        LEFT JOIN transportation_offices AS counseling_offices
            ON moves.counseling_transportation_office_id = counseling_offices.id
        LEFT JOIN shipment_address_updates ON shipment_address_updates.shipment_id = mto_shipments.id
        JOIN move_to_dest_gbloc ON move_to_dest_gbloc.move_id = moves.id
        WHERE moves.show = TRUE
    ';

    IF user_gbloc IS NOT NULL THEN
        sql_query := sql_query || ' AND move_to_dest_gbloc.gbloc = $1 ';
    END IF;

    IF customer_name IS NOT NULL THEN
        sql_query := sql_query || ' AND (
            service_members.first_name || '' '' || service_members.last_name ILIKE ''%'' || $2 || ''%''
            OR service_members.last_name || '' '' || service_members.first_name ILIKE ''%'' || $2 || ''%''
        )';
    END IF;

    IF edipi IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.edipi ILIKE ''%'' || $3 || ''%'' ';
    END IF;

    IF emplid IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.emplid ILIKE ''%'' || $4 || ''%'' ';
    END IF;

    IF m_status IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.status = ANY($5) ';
    END IF;

    IF move_code IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.locator ILIKE ''%'' || $6 || ''%'' ';
    END IF;

    IF requested_move_date IS NOT NULL THEN
        sql_query := sql_query || ' AND (
            mto_shipments.requested_pickup_date::DATE = $7::DATE
            OR ppm_shipments.expected_departure_date::DATE = $7::DATE
            OR (mto_shipments.shipment_type = ''HHG_OUTOF_NTS'' AND mto_shipments.requested_delivery_date::DATE = $7::DATE)
        )';
    END IF;

    IF date_submitted IS NOT NULL THEN
        sql_query := sql_query || ' AND moves.submitted_at::DATE = $8::DATE ';
    END IF;

    IF branch IS NOT NULL THEN
        sql_query := sql_query || ' AND service_members.affiliation ILIKE ''%'' || $9 || ''%'' ';
    END IF;

    IF origin_duty_location IS NOT NULL THEN
        sql_query := sql_query || ' AND origin_duty_locations.name ILIKE ''%'' || $10 || ''%'' ';
    END IF;

    IF counseling_office IS NOT NULL THEN
        sql_query := sql_query || ' AND counseling_offices.name ILIKE ''%'' || $11 || ''%'' ';
    END IF;

    IF too_assigned_user IS NOT NULL THEN
        sql_query := sql_query || ' AND (too_user.first_name || '' '' || too_user.last_name) ILIKE ''%'' || $12 || ''%'' ';
    END IF;

    -- add destination queue-specific filters (pending dest address requests, dest SIT & dest shuttle service items)
    sql_query := sql_query || '
        AND (
            shipment_address_updates.status = ''REQUESTED''
            OR (
                mto_service_items.status = ''SUBMITTED''
                AND re_services.code IN (''DDFSIT'', ''DDASIT'', ''DDDSIT'', ''DDSHUT'', ''DDSFSC'', ''IDFSIT'', ''IDASIT'', ''IDDSIT'', ''IDSHUT'')
            )
        )
    ';

    -- default sorting values if none are provided (move.id)
    sort_column := 'id';
    sort_order := 'ASC';

    IF sort IS NOT NULL THEN
        CASE sort
            WHEN 'locator' THEN sort_column := 'moves.locator';
            WHEN 'status' THEN sort_column := 'moves.status';
            WHEN 'customerName' THEN sort_column := 'service_members.last_name';
            WHEN 'edipi' THEN sort_column := 'service_members.edipi';
            WHEN 'emplid' THEN sort_column := 'service_members.emplid';
            WHEN 'requestedMoveDate' THEN sort_column := 'COALESCE(mto_shipments.requested_pickup_date, ppm_shipments.expected_departure_date, mto_shipments.requested_delivery_date)';
            WHEN 'appearedInTooAt' THEN sort_column := 'COALESCE(moves.submitted_at, moves.approvals_requested_at)';
            WHEN 'branch' THEN sort_column := 'service_members.affiliation';
            WHEN 'originDutyLocation' THEN sort_column := 'origin_duty_locations.name';
            WHEN 'counselingOffice' THEN sort_column := 'counseling_offices.name';
            WHEN 'assignedTo' THEN sort_column := 'too_user.last_name';
            ELSE
                sort_column := 'moves.id';
        END CASE;
    END IF;

    IF sort_direction IS NOT NULL THEN
        IF LOWER(sort_direction) = 'desc' THEN
            sort_order := 'DESC';
        ELSE
            sort_order := 'ASC';
        END IF;
    END IF;

    sql_query := sql_query || '
        GROUP BY
            moves.id,
            moves.show,
            moves.locator,
            moves.submitted_at,
            moves.orders_id,
            moves.status,
            moves.locked_by,
            moves.too_assigned_id,
            moves.counseling_transportation_office_id,
            mto_shipments.requested_pickup_date,
            mto_shipments.requested_delivery_date,
            ppm_shipments.expected_departure_date,
            orders.id,
            service_members.id,
            service_members.first_name,
            service_members.last_name,
            service_members.edipi,
            service_members.emplid,
            service_members.affiliation,
            origin_duty_locations.name,
            counseling_offices.name,
            too_user.first_name,
            too_user.last_name';
    sql_query := sql_query || format(' ORDER BY %s %s ', sort_column, sort_order);
    sql_query := sql_query || ' LIMIT $13 OFFSET $14 ';

    RETURN QUERY EXECUTE sql_query
    USING user_gbloc, customer_name, edipi, emplid, m_status, move_code, requested_move_date, date_submitted,
          branch, origin_duty_location, counseling_office, too_assigned_user, per_page, offset_value;

END;
$$ LANGUAGE plpgsql;
