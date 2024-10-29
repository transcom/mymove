CREATE TABLE IF NOT EXISTS gbloc_aors
(id 			uuid		NOT NULL,
jppso_regions_id uuid       NOT NULL
    CONSTRAINT fk_gbloc_aors_jppso_regions_id REFERENCES jppso_regions (id),
oconus_rate_area_id	uuid		NOT NULL
	CONSTRAINT fk_gbloc_aors_oconus_rate_area_id REFERENCES re_oconus_rate_areas (id),
department_indicator varchar(255),
shipment_type   mto_shipment_type,
is_active       bool default TRUE,
created_at		timestamp	NOT NULL DEFAULT NOW(),
updated_at		timestamp	NOT NULL DEFAULT NOW(),
CONSTRAINT gbloc_aors_pkey PRIMARY KEY (id),
CONSTRAINT unique_gbloc_aors UNIQUE (jppso_regions_id, oconus_rate_area_id, department_indicator, shipment_type)
);

COMMENT ON TABLE gbloc_aors IS 'Associates a rate area to one or more gblocs. Some associations are based on branch of service or shipment type.';
COMMENT ON COLUMN gbloc_aors.jppso_regions_id IS 'The associated id for the JPPSO and GBLOC (Government Bill of Lading Office Code)';
COMMENT ON COLUMN gbloc_aors.oconus_rate_area_id IS 'The id for associated rate area and location.';
COMMENT ON COLUMN gbloc_aors.department_indicator is 'Name of the service branch from orders table.';
COMMENT ON COLUMN gbloc_aors.shipment_type IS 'The type of shipment from mto_shipments table.';
COMMENT ON COLUMN gbloc_aors.is_active IS 'Set to true if the record is active.';

INSERT INTO jppso_regions
(id, code, name, created_at, updated_at)
VALUES('98281c1b-6161-4e74-a119-177b2a2fb176', 'MBFL', 'JPPSO Elmendorf-Richardson', now(), now());

INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('27eac48d-584a-4369-be93-3a9e6e3f8550','98281c1b-6161-4e74-a119-177b2a2fb176','e418be11-2b6f-4714-b026-e293528c50bd',NULL,NULL,true,now(),now()),
	 ('f2a583af-4407-42a8-9205-b72ee26e5f0b','98281c1b-6161-4e74-a119-177b2a2fb176','b41c5636-96dd-4f0d-a18e-eebb17f97ea5',NULL,NULL,true,now(),now()),
	 ('88c4a2e6-43d4-4c86-b6b4-24321dbc6c36','98281c1b-6161-4e74-a119-177b2a2fb176','02f32a6a-0338-4545-8437-059862892d2c',NULL,NULL,true,now(),now()),
	 ('9d51fd5a-6197-4e57-8b61-1f86da5e220b','98281c1b-6161-4e74-a119-177b2a2fb176','0f128476-d52a-418c-8ba0-c8bfd1c32629',NULL,NULL,true,now(),now()),
	 ('ca2bf299-392a-4bab-9b85-d9da613f0114','98281c1b-6161-4e74-a119-177b2a2fb176','081f84c3-17ec-4ff6-97ce-d5c44a8e4a28',NULL,NULL,true,now(),now()),
	 ('7673c201-51c5-4c81-9d08-5998b4ba3b9b','98281c1b-6161-4e74-a119-177b2a2fb176','ce7cdd91-e323-43a6-a604-daaa6bf8be06',NULL,NULL,true,now(),now()),
	 ('23b4f713-9ae0-433d-a390-e09d6abc169d','98281c1b-6161-4e74-a119-177b2a2fb176','655050d4-711a-43b7-b06d-828ec990c35e',NULL,NULL,true,now(),now()),
	 ('0f61e7c2-3698-4486-b192-ac57edcf8fde','98281c1b-6161-4e74-a119-177b2a2fb176','455a34af-30a8-4a98-a62b-6f40fd7f047b',NULL,NULL,true,now(),now()),
	 ('dade08ca-5622-4f0e-8631-853ce440bd85','98281c1b-6161-4e74-a119-177b2a2fb176','1bf1daee-51bb-4c28-aac8-a126ef596486',NULL,NULL,true,now(),now()),
	 ('65890bf6-d4c1-4aa5-9e62-1e210fb3a322','98281c1b-6161-4e74-a119-177b2a2fb176','15b1d852-0dde-4e1b-b3c6-e08bbc714db3',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('1ebdc3bf-c37a-4739-92ca-be32a63e2968','98281c1b-6161-4e74-a119-177b2a2fb176','d63de5bb-379b-4f3b-b4c1-554b9746d311',NULL,NULL,true,now(),now()),
	 ('48d3af3e-d86e-4fd8-9466-15759a71d567','98281c1b-6161-4e74-a119-177b2a2fb176','95d142f8-b50a-4108-b5e2-fbd3b7022d3b',NULL,NULL,true,now(),now()),
	 ('89338aca-cf8d-4dc4-a578-c93df3afdbaf','98281c1b-6161-4e74-a119-177b2a2fb176','922bb7da-d0e9-431c-a493-95a861dca929',NULL,NULL,true,now(),now()),
	 ('208840b4-a367-470d-804d-193a08d60b9c','98281c1b-6161-4e74-a119-177b2a2fb176','79aef0ca-9065-4c0e-a134-2889b250cc38',NULL,NULL,true,now(),now()),
	 ('aeaf7f69-204c-40ca-8f77-6748d1ab6471','98281c1b-6161-4e74-a119-177b2a2fb176','1f181025-e410-41ac-935b-4b5147353f84',NULL,NULL,true,now(),now()),
	 ('0631520d-8195-43a8-9ecf-a6b335356ca6','98281c1b-6161-4e74-a119-177b2a2fb176','d9ede620-c3f1-4f8d-b60b-eb93664382f7',NULL,NULL,true,now(),now()),
	 ('97eab888-0af3-4223-9eb3-1cda1d34977c','98281c1b-6161-4e74-a119-177b2a2fb176','bd27700b-3094-46cc-807b-f18472cbfaf0',NULL,NULL,true,now(),now()),
	 ('07680ff6-2750-4a57-a341-69498aa5315a','98281c1b-6161-4e74-a119-177b2a2fb176','2f4a1689-ee65-45fa-9d0c-debd9344f8b9',NULL,NULL,true,now(),now()),
	 ('5e130565-f439-46da-ba5a-29bb84b925ab','98281c1b-6161-4e74-a119-177b2a2fb176','20c52934-2588-4d8f-b3ed-c046d771f4e9',NULL,NULL,true,now(),now()),
	 ('4819d2eb-afe2-4fca-a6a2-538c73fe119c','98281c1b-6161-4e74-a119-177b2a2fb176','15fe364b-add5-4ade-bdd9-38fe442616fb',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('90294b38-8e47-41ab-9c50-920b87bfadb2','98281c1b-6161-4e74-a119-177b2a2fb176','e5a1248e-3870-4a4c-9c9d-2056234eb64a',NULL,NULL,true,now(),now()),
	 ('9b7a69b9-ba5a-45e2-bdbd-0dc06cea9ec5','98281c1b-6161-4e74-a119-177b2a2fb176','efc92db1-4b06-49ab-a295-a98f3f6c9c04',NULL,NULL,true,now(),now()),
	 ('ddefba33-bf44-4854-9a06-4892ea6c125d','98281c1b-6161-4e74-a119-177b2a2fb176','13474ce5-8839-4af7-b975-c3f70ccdab7b',NULL,NULL,true,now(),now()),
	 ('dd323732-800c-4242-b537-46bb9b7fb05f','98281c1b-6161-4e74-a119-177b2a2fb176','d3b05c5e-6faa-4fa9-b725-0904f8a4f3d7',NULL,NULL,true,now(),now()),
	 ('af5f7ff8-1563-46cd-a730-343e4477f061','98281c1b-6161-4e74-a119-177b2a2fb176','e5f849fe-5672-4d0d-881c-078e56eea33d',NULL,NULL,true,now(),now()),
	 ('487fa2c4-31e2-4259-9911-bfc2f970ebae','98281c1b-6161-4e74-a119-177b2a2fb176','fca7c60e-fbd9-4885-afdb-ae41c521b560',NULL,NULL,true,now(),now()),
	 ('6cbd8a95-e6db-4183-a526-5ca255a9daac','98281c1b-6161-4e74-a119-177b2a2fb176','0d52a0f0-f39c-4d34-9387-2df45a5810d4',NULL,NULL,true,now(),now()),
	 ('7b0771d0-4215-4581-9abc-717ff32e9008','98281c1b-6161-4e74-a119-177b2a2fb176','af008e75-81d5-4211-8204-4964c78e70d9',NULL,NULL,true,now(),now()),
	 ('462133c7-10e2-462d-a58a-5c5556327cf9','98281c1b-6161-4e74-a119-177b2a2fb176','dd689e55-af29-4c76-b7e0-c2429ea4833c',NULL,NULL,true,now(),now()),
	 ('c2b5798f-4edd-4bbe-9580-bd67a2678078','98281c1b-6161-4e74-a119-177b2a2fb176','60e08330-6869-4586-8198-35f7a4ae9ea7',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('07804fb2-9c3d-42b8-8327-bbbd776195bb','98281c1b-6161-4e74-a119-177b2a2fb176','f7478e79-dbfe-46c8-b337-1e7c46df79dc',NULL,NULL,true,now(),now()),
	 ('b598506c-c21d-41da-b0dd-d48d9acb55a0','98281c1b-6161-4e74-a119-177b2a2fb176','486d7dd4-f51f-4b13-88fc-830b5b15f0a8',NULL,NULL,true,now(),now()),
	 ('37e93f48-fc70-4c55-8601-87eff9331e02','98281c1b-6161-4e74-a119-177b2a2fb176','f66edf62-ba5e-47f8-8264-b6dc9e7dd9ba',NULL,NULL,true,now(),now()),
	 ('43a2b490-0d10-4717-ae5a-c7f395dd6471','98281c1b-6161-4e74-a119-177b2a2fb176','0052c55a-d9d7-46b0-9328-58d3911f61b4',NULL,NULL,true,now(),now()),
	 ('68346448-ddf0-45b4-a365-34aee83ce5d3','98281c1b-6161-4e74-a119-177b2a2fb176','16a51fd1-04ed-432a-a8d7-9f17c1de095d',NULL,NULL,true,now(),now()),
	 ('22a2e822-1022-41a2-95d9-ddd040eb55c6','98281c1b-6161-4e74-a119-177b2a2fb176','64b4756f-437d-4aa5-a95e-c396f0cafcbb',NULL,NULL,true,now(),now()),
	 ('0d391be1-4b5b-4cb3-b1da-84864688fb23','98281c1b-6161-4e74-a119-177b2a2fb176','f2de1da5-a737-4493-b3f7-7700944a5b62',NULL,NULL,true,now(),now()),
	 ('490b36dd-b7cb-447b-9c74-e1b2bf81f42b','98281c1b-6161-4e74-a119-177b2a2fb176','7a9c2adb-0562-42c1-a5f1-81710dd19590',NULL,NULL,true,now(),now()),
	 ('89072d15-3e94-4ed8-be13-13fbd89ca1f8','98281c1b-6161-4e74-a119-177b2a2fb176','e713eed6-35a6-456d-bfb2-0e0646078ab8',NULL,NULL,true,now(),now()),
	 ('a51a63a1-6be7-4268-a007-df16c66e3770','98281c1b-6161-4e74-a119-177b2a2fb176','a4594957-0ec6-4edc-85c2-68871f4f6359',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('ddc5c9c8-cc48-41e1-9c70-4253a1bcc9b8','98281c1b-6161-4e74-a119-177b2a2fb176','acd4d2f9-9a49-4a73-be14-3515f19ba0d1',NULL,NULL,true,now(),now()),
	 ('b865802d-dc3a-45cd-bc0b-bbebdd60a786','98281c1b-6161-4e74-a119-177b2a2fb176','c8300ab6-519c-4bef-9b81-25e7334773ca',NULL,NULL,true,now(),now()),
	 ('1b207ece-29a0-4e84-9c02-38f4a8062010','98281c1b-6161-4e74-a119-177b2a2fb176','67eca2ce-9d88-44ca-bb17-615f3199415c',NULL,NULL,true,now(),now()),
	 ('674d99d3-ee83-44b6-9819-d125cf075a20','98281c1b-6161-4e74-a119-177b2a2fb176','149c3a94-abb1-4af0-aabf-3af019d5e243',NULL,NULL,true,now(),now()),
	 ('71051877-1c8f-44e6-abc2-6e2b2267d8f8','98281c1b-6161-4e74-a119-177b2a2fb176','87c09b47-058e-47ea-9684-37a8ac1f7120',NULL,NULL,true,now(),now()),
	 ('526c1056-b373-4b87-83c3-0d0d602257fc','98281c1b-6161-4e74-a119-177b2a2fb176','0ba81e85-f175-435a-a7c2-a22d0c44cc7b',NULL,NULL,true,now(),now()),
	 ('fac9d7f3-8780-43bd-b71d-918500310ef7','98281c1b-6161-4e74-a119-177b2a2fb176','946256dc-1572-4201-b5ff-464da876f5ff',NULL,NULL,true,now(),now()),
	 ('70b5d93e-a85e-4a3d-a470-ed739405930d','98281c1b-6161-4e74-a119-177b2a2fb176','5d381518-6beb-422e-8c57-4682b87ff1fe',NULL,NULL,true,now(),now()),
	 ('34ab419c-7b3f-4c26-9dd9-1a180b9c123f','98281c1b-6161-4e74-a119-177b2a2fb176','b07fb701-2c9b-4847-a459-76b1f47aa872',NULL,NULL,true,now(),now()),
	 ('cd44fa2b-ef6c-47d5-8c59-8b3fe7d5cc49','98281c1b-6161-4e74-a119-177b2a2fb176','b6c51584-f3e8-4c9e-9bdf-f7cac0433319',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('c9c06237-d3ac-4784-a959-cdf66d6a77bf','98281c1b-6161-4e74-a119-177b2a2fb176','6097d01d-78ef-40d5-8699-7f8a8e48f4e7',NULL,NULL,true,now(),now()),
	 ('fd37701d-136d-48e2-a56d-d2439429940c','98281c1b-6161-4e74-a119-177b2a2fb176','e48a3602-7fdd-4527-a8a7-f244fb331228',NULL,NULL,true,now(),now()),
	 ('1a20dc73-63b0-410c-87bc-f0c57bb4bde4','98281c1b-6161-4e74-a119-177b2a2fb176','df47ef29-e902-4bd3-a32d-88d6187399b3',NULL,NULL,true,now(),now()),
	 ('bf7d6ad1-1668-4771-8928-87dabc6bd5e7','98281c1b-6161-4e74-a119-177b2a2fb176','0f4eddf9-b727-4725-848e-3d9329553823',NULL,NULL,true,now(),now()),
	 ('47839d4f-b9d8-43fc-9fd3-b52ccd45bb27','98281c1b-6161-4e74-a119-177b2a2fb176','d16ebe80-2195-48cf-ba61-b5efb8212754',NULL,NULL,true,now(),now()),
	 ('63290ae8-1edb-4a90-848f-c60adde8f124','98281c1b-6161-4e74-a119-177b2a2fb176','093268d0-597b-40bd-882a-8f385480bc68',NULL,NULL,true,now(),now()),
	 ('2359610a-dab3-48db-9e21-c71331157fa6','98281c1b-6161-4e74-a119-177b2a2fb176','edb8b022-9534-44e7-87ea-e93f5451f057',NULL,NULL,true,now(),now()),
	 ('61b2b91d-6ab0-483b-a3f6-3fe1f4cde43a','98281c1b-6161-4e74-a119-177b2a2fb176','4c62bb35-4d2a-4995-9c4f-8c665f2f6d3e',NULL,NULL,true,now(),now()),
	 ('1ed10063-859f-4aa2-809c-ae50133eed37','98281c1b-6161-4e74-a119-177b2a2fb176','58ccfcc5-ded6-4f91-8cb7-8687bc56c4c6',NULL,NULL,true,now(),now()),
	 ('d4cbd597-8519-47fb-858f-081d2c4a1a8c','98281c1b-6161-4e74-a119-177b2a2fb176','dbc839c6-7b56-45b0-ab36-a0e77b0d538c',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('3aed072e-d27f-4b13-ac8c-16b6b2442e27','98281c1b-6161-4e74-a119-177b2a2fb176','c77e434a-1bf9-44c6-92aa-d377d72d1d44',NULL,NULL,true,now(),now()),
	 ('bbe5eed0-3e02-41b1-87d4-d4358e189ccc','98281c1b-6161-4e74-a119-177b2a2fb176','89533190-e613-4c36-99df-7ee3871bb071',NULL,NULL,true,now(),now()),
	 ('92e205c8-8d77-4f6f-be13-bf14ce47173e','98281c1b-6161-4e74-a119-177b2a2fb176','c86cb50a-99f6-41ed-8e9d-f096bd5cadca',NULL,NULL,true,now(),now()),
	 ('a576cd3b-1046-4001-93f3-5f93e5ec1815','98281c1b-6161-4e74-a119-177b2a2fb176','b65bd7c4-6251-4682-aa07-77f761c363af',NULL,NULL,true,now(),now()),
	 ('fbb214d1-4bcd-4846-9427-8f175567ebe9','98281c1b-6161-4e74-a119-177b2a2fb176','88df8618-798a-4a12-9b14-a7267a6cfd7f',NULL,NULL,true,now(),now()),
	 ('d7c19a96-da14-4bc8-bf9f-3a319c3d61df','98281c1b-6161-4e74-a119-177b2a2fb176','15894c4f-eb1a-4b6b-8b3f-e3abc4eeee8d',NULL,NULL,true,now(),now()),
	 ('ead9ae06-ad19-4e86-8ada-1cdc2362cdd0','98281c1b-6161-4e74-a119-177b2a2fb176','83fd3eb7-6269-46e2-84e1-f6180f40b9e8',NULL,NULL,true,now(),now()),
	 ('a99751f7-80e9-4ca2-86db-d21e79045513','98281c1b-6161-4e74-a119-177b2a2fb176','b0f95e07-e3d1-4403-a402-50be9a542cf9',NULL,NULL,true,now(),now()),
	 ('4c767bd9-5ede-4197-b957-1b91ea9236ea','98281c1b-6161-4e74-a119-177b2a2fb176','cef18ff2-8f48-47c0-8062-a40f9dec641c',NULL,NULL,true,now(),now()),
	 ('ffaf4885-a55c-4a3b-b8d8-75556df99969','98281c1b-6161-4e74-a119-177b2a2fb176','15eefe71-55ae-40b0-8bb5-f1952fcf45c8',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('2646897d-2e62-46aa-83e0-cf9c91b7b47d','98281c1b-6161-4e74-a119-177b2a2fb176','d3c5d0a7-0477-44f6-8279-f6b9fa7b3436',NULL,NULL,true,now(),now()),
	 ('1efc97f2-eaf2-443c-a9cf-ccda5131d1cf','98281c1b-6161-4e74-a119-177b2a2fb176','8b89e7d7-cb36-40bd-ba14-699f1e4b1806',NULL,NULL,true,now(),now()),
	 ('5e588e1b-4475-4e59-8d2b-720349e222ad','98281c1b-6161-4e74-a119-177b2a2fb176','5b833661-8fda-4f2a-9f2b-e1f4c21749a4',NULL,NULL,true,now(),now()),
	 ('d315a7e1-49ca-4ff7-b8a2-c3a178ecf8ec','98281c1b-6161-4e74-a119-177b2a2fb176','38bed277-7990-460d-ad27-a69559638f42',NULL,NULL,true,now(),now()),
	 ('c008d9b4-2e4b-411f-b967-01a30ad80962','98281c1b-6161-4e74-a119-177b2a2fb176','67081e66-7e84-4bac-9c31-aa3e45bbba23',NULL,NULL,true,now(),now()),
	 ('2ad45d8f-89e5-4dbd-9265-7bc3130b183a','98281c1b-6161-4e74-a119-177b2a2fb176','99e3f894-595f-403a-83f5-f7e035dd1f20',NULL,NULL,true,now(),now()),
	 ('01844c50-cb98-46e3-91e2-cd603795afd6','615673a2-d393-4d46-95f9-705fbeb6bc79','1c0a5988-492f-4bc2-8409-9a143a494248',NULL,NULL,true,now(),now()),
	 ('410070fb-6af7-49a4-b8a7-1f8a2febd79f','615673a2-d393-4d46-95f9-705fbeb6bc79','63e10b54-5348-4362-932e-6613a8db4d42',NULL,NULL,true,now(),now()),
	 ('a7f4bd72-11d4-4609-a86c-c05153dc02f4','615673a2-d393-4d46-95f9-705fbeb6bc79','437e9930-0625-4448-938d-ba93a0a98ab5',NULL,NULL,true,now(),now()),
	 ('4ef0b319-c082-459f-ae17-ba97bf38fb8d','615673a2-d393-4d46-95f9-705fbeb6bc79','87d3cb47-1026-4f42-bff6-899ff0fa7660',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('b44df240-6ee9-44d8-85d0-f99cfdfb04eb','615673a2-d393-4d46-95f9-705fbeb6bc79','0115586e-be8c-4808-b65a-d417fad19238',NULL,NULL,true,now(),now()),
	 ('5d20c776-2cf5-498c-b140-97bc7e76ac9a','615673a2-d393-4d46-95f9-705fbeb6bc79','6638a77d-4b91-48f6-8241-c87fbdddacd1',NULL,NULL,true,now(),now()),
	 ('24fb918e-e9a8-49ae-a6e0-bb6474a627ff','615673a2-d393-4d46-95f9-705fbeb6bc79','97c16bc3-e174-410b-9a09-a0db31420dbc',NULL,NULL,true,now(),now()),
	 ('8a888e31-29f8-49cd-a5b1-4930c2133f8b','615673a2-d393-4d46-95f9-705fbeb6bc79','3c6c5e35-1281-45fc-83ee-e6d656e155b6',NULL,NULL,true,now(),now()),
	 ('63c95649-3984-4559-936f-8a56703ff2e1','615673a2-d393-4d46-95f9-705fbeb6bc79','6af248be-e5a8-49e9-a9e2-516748279ab5',NULL,NULL,true,now(),now()),
	 ('14d7bcd2-8687-4792-97c9-642b81f2e93b','615673a2-d393-4d46-95f9-705fbeb6bc79','3580abe3-84da-4b46-af7b-d4379e6cff46',NULL,NULL,true,now(),now()),
	 ('b7a2da46-17c8-4609-972e-8c98e5c2142a','615673a2-d393-4d46-95f9-705fbeb6bc79','f3bb397e-04e6-4b37-9bf3-b3ebab79a9b6',NULL,NULL,true,now(),now()),
	 ('6c9b6668-b251-4238-b630-d824574fe6e5','615673a2-d393-4d46-95f9-705fbeb6bc79','f9bfe297-4ee0-4f76-a4bd-64b3a514af5d',NULL,NULL,true,now(),now()),
	 ('c86c153c-a2ce-4f9e-8fd5-ae6adc2ea727','615673a2-d393-4d46-95f9-705fbeb6bc79','70f11a71-667b-422d-ae37-8f25539f7782',NULL,NULL,true,now(),now()),
	 ('ae22641d-3501-4c12-b93f-d291854c6e03','615673a2-d393-4d46-95f9-705fbeb6bc79','baf66a9a-3c8a-49e7-83ea-841b9960e184',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('0744ce05-1b16-4afa-a11b-b7f4553f5733','615673a2-d393-4d46-95f9-705fbeb6bc79','e96d46b1-3ab7-4b29-b798-ec3b728dd6a1',NULL,NULL,true,now(),now()),
	 ('7fe40bae-483d-43df-a0d6-7430f7d0f6f6','615673a2-d393-4d46-95f9-705fbeb6bc79','e4171b5b-d26c-43b3-b41c-68b43bcbb079',NULL,NULL,true,now(),now()),
	 ('d66ce151-95c0-4c67-859e-a5f7a705907a','615673a2-d393-4d46-95f9-705fbeb6bc79','a0651bec-1258-4e36-9c76-155247e42c0a',NULL,NULL,true,now(),now()),
	 ('62c6117b-30a2-4c5a-b7fc-073e251bdf78','615673a2-d393-4d46-95f9-705fbeb6bc79','737a8e63-af19-4902-a4b5-8f80e2268e4b',NULL,NULL,true,now(),now()),
	 ('ac05e2cd-9297-42a5-9747-b01cff5cff37','615673a2-d393-4d46-95f9-705fbeb6bc79','870df26e-2c50-4512-aa2e-61094cfbc3e1',NULL,NULL,true,now(),now()),
	 ('2055d575-f002-4272-824a-9d4cea3c0819','615673a2-d393-4d46-95f9-705fbeb6bc79','09dc6547-d346-40c3-93fb-fb7be1fa1b3e',NULL,NULL,true,now(),now()),
	 ('e9582f65-b859-4876-9770-f340b42742f4','615673a2-d393-4d46-95f9-705fbeb6bc79','43d7d081-bb32-4544-84f1-419fe0cb76e1',NULL,NULL,true,now(),now()),
	 ('06533a4d-4237-4e32-b753-740351f3d778','615673a2-d393-4d46-95f9-705fbeb6bc79','d7e57942-9c83-4138-baa0-70e8b5f08598',NULL,NULL,true,now(),now()),
	 ('028b5785-1df2-4443-9376-22d239b9fd59','615673a2-d393-4d46-95f9-705fbeb6bc79','c2f06691-5989-41a3-a848-19c9f0fec5df',NULL,NULL,true,now(),now()),
	 ('2619e3c1-736c-47f2-bdce-2cdcac1045de','615673a2-d393-4d46-95f9-705fbeb6bc79','5263a9ed-ff4d-42cc-91d5-dbdefeef54d1',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('1e89bd43-301a-4680-9092-84b3ca50207c','615673a2-d393-4d46-95f9-705fbeb6bc79','54ab3a49-6d78-4922-bac1-94a722b9859a',NULL,NULL,true,now(),now()),
	 ('c1640891-bd72-4fc8-b302-5ff6ec128385','615673a2-d393-4d46-95f9-705fbeb6bc79','3b6fe2e9-9116-4a70-8eef-96205205b0e3',NULL,NULL,true,now(),now()),
	 ('7394845c-b673-40ba-8302-79c0c539e201','615673a2-d393-4d46-95f9-705fbeb6bc79','ba0017c8-5d48-4efe-8802-361d2f2bc16d',NULL,NULL,true,now(),now()),
	 ('db8758df-067c-4044-b139-a6eb314f018e','615673a2-d393-4d46-95f9-705fbeb6bc79','af65234a-8577-4f6d-a346-7d486e963287',NULL,NULL,true,now(),now()),
	 ('7a9a1d2e-fa9c-48a8-8d0a-0bfe1a6e6027','615673a2-d393-4d46-95f9-705fbeb6bc79','25f62695-ba9b-40c3-a7e4-0c078731123d',NULL,NULL,true,now(),now()),
	 ('6d03985c-e5e1-4bd3-96f4-4d1adaf45f78','615673a2-d393-4d46-95f9-705fbeb6bc79','b0422b13-4afe-443e-901f-fe652cde23a4',NULL,NULL,true,now(),now()),
	 ('a014de5c-04bf-4982-b884-fc27002ef640','615673a2-d393-4d46-95f9-705fbeb6bc79','ecbc1d89-9cf6-4f52-b453-3ba473a0ff4e',NULL,NULL,true,now(),now()),
	 ('e883b15d-698d-45a0-9cc5-9c4cd8b9d53d','615673a2-d393-4d46-95f9-705fbeb6bc79','4e85cb86-e9dd-4c7c-9677-e0a327ac895c',NULL,NULL,true,now(),now()),
	 ('fff0f9ce-247e-44c5-98db-d5f82938823d','615673a2-d393-4d46-95f9-705fbeb6bc79','7f9dd10d-100e-4252-9786-706349f456ca',NULL,NULL,true,now(),now()),
	 ('7dcdafdf-5471-4086-8ae2-200d1b2b7d2f','615673a2-d393-4d46-95f9-705fbeb6bc79','5451f8d7-60c5-4e22-bbf6-d9af8e6ace54',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('ff198351-2f14-4ee4-8f6c-509b21e7f866','615673a2-d393-4d46-95f9-705fbeb6bc79','9e9cba85-7c39-4836-809d-70b54baf392e',NULL,NULL,true,now(),now()),
	 ('e9a496a3-160e-42eb-8d9c-5738973fea95','615673a2-d393-4d46-95f9-705fbeb6bc79','f7470914-cf48-43be-b431-a3ca2fe5b290',NULL,NULL,true,now(),now()),
	 ('c9cd1c1a-5b76-430f-8a55-14a134fbc6e5','615673a2-d393-4d46-95f9-705fbeb6bc79','5a0d3cc1-b866-4bde-b67f-78d565facf3e',NULL,NULL,true,now(),now()),
	 ('3a724501-afea-4b77-b7bd-2491628ef194','615673a2-d393-4d46-95f9-705fbeb6bc79','e629b95a-ec5b-4fc4-897f-e0d1050e1ec6',NULL,NULL,true,now(),now()),
	 ('4931072b-a3d7-4d9d-8744-de7ac20124ef','615673a2-d393-4d46-95f9-705fbeb6bc79','d68a626f-935a-4eb1-ba9b-6829feeff91c',NULL,NULL,true,now(),now()),
	 ('4868c585-5cd8-4fb3-b7f7-affd3ffd56c0','98281c1b-6161-4e74-a119-177b2a2fb176','1c0a5988-492f-4bc2-8409-9a143a494248','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('a575d822-db66-4f2f-b7e6-0893516f1ed8','98281c1b-6161-4e74-a119-177b2a2fb176','63e10b54-5348-4362-932e-6613a8db4d42','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('9eb28d55-bae4-4481-9623-5811403c0173','98281c1b-6161-4e74-a119-177b2a2fb176','437e9930-0625-4448-938d-ba93a0a98ab5','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('7d82a2cd-a668-4b8a-8758-ed75742ff27e','98281c1b-6161-4e74-a119-177b2a2fb176','87d3cb47-1026-4f42-bff6-899ff0fa7660','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('42306eb8-1b76-4286-bb7e-8c573adc77f7','98281c1b-6161-4e74-a119-177b2a2fb176','0115586e-be8c-4808-b65a-d417fad19238','AIR_AND_SPACE_FORCE',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('c70c9d51-8788-459a-9b3d-8affa08198c4','98281c1b-6161-4e74-a119-177b2a2fb176','6638a77d-4b91-48f6-8241-c87fbdddacd1','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('0b39b389-f18f-4f78-a554-9f8f58a6b833','98281c1b-6161-4e74-a119-177b2a2fb176','97c16bc3-e174-410b-9a09-a0db31420dbc','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('f3ff59bf-d975-44b9-bf4f-0f67d8917ec1','98281c1b-6161-4e74-a119-177b2a2fb176','3c6c5e35-1281-45fc-83ee-e6d656e155b6','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('ecdcbd7a-307b-4b48-97a1-8c4e009917ef','98281c1b-6161-4e74-a119-177b2a2fb176','6af248be-e5a8-49e9-a9e2-516748279ab5','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('417cad20-5d71-46eb-9e12-164dc92cd068','98281c1b-6161-4e74-a119-177b2a2fb176','3580abe3-84da-4b46-af7b-d4379e6cff46','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('a4092978-8956-48d0-bcc6-8fc98be478c9','98281c1b-6161-4e74-a119-177b2a2fb176','f3bb397e-04e6-4b37-9bf3-b3ebab79a9b6','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('0781002d-477e-4f6d-87e7-7392b3722199','98281c1b-6161-4e74-a119-177b2a2fb176','f9bfe297-4ee0-4f76-a4bd-64b3a514af5d','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('79190000-1945-4d38-9827-bde5d363a72d','98281c1b-6161-4e74-a119-177b2a2fb176','70f11a71-667b-422d-ae37-8f25539f7782','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('7756f16d-ff8f-4494-b842-69b075982e1b','98281c1b-6161-4e74-a119-177b2a2fb176','baf66a9a-3c8a-49e7-83ea-841b9960e184','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('2477354b-022c-402a-8659-8bf917f6c022','98281c1b-6161-4e74-a119-177b2a2fb176','e96d46b1-3ab7-4b29-b798-ec3b728dd6a1','AIR_AND_SPACE_FORCE',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('1c501fe3-7e15-4161-99b9-5c49d7fa8ba2','98281c1b-6161-4e74-a119-177b2a2fb176','e4171b5b-d26c-43b3-b41c-68b43bcbb079','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('550a26d5-0a50-41e2-93db-5eb3198afa54','98281c1b-6161-4e74-a119-177b2a2fb176','a0651bec-1258-4e36-9c76-155247e42c0a','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('6b9c6176-a59a-41dc-9765-f5beca1b5b9c','98281c1b-6161-4e74-a119-177b2a2fb176','737a8e63-af19-4902-a4b5-8f80e2268e4b','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('b867ff82-9728-4a08-9207-9a76e39528b4','98281c1b-6161-4e74-a119-177b2a2fb176','870df26e-2c50-4512-aa2e-61094cfbc3e1','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('dd980455-ade7-40e5-9d52-b02fce51a6d0','98281c1b-6161-4e74-a119-177b2a2fb176','09dc6547-d346-40c3-93fb-fb7be1fa1b3e','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('d40ed977-db2e-4bb7-8b01-627afd7982c8','98281c1b-6161-4e74-a119-177b2a2fb176','43d7d081-bb32-4544-84f1-419fe0cb76e1','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('bc6ee1d1-c24c-4949-883b-6143811d75a0','98281c1b-6161-4e74-a119-177b2a2fb176','d7e57942-9c83-4138-baa0-70e8b5f08598','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('eb352383-c4ad-4fba-844f-bc5ea26df9b0','98281c1b-6161-4e74-a119-177b2a2fb176','c2f06691-5989-41a3-a848-19c9f0fec5df','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('f7fe54b1-1883-4254-86ec-de74b3fd3d83','98281c1b-6161-4e74-a119-177b2a2fb176','5263a9ed-ff4d-42cc-91d5-dbdefeef54d1','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('be062c87-f98e-4a7e-9c5e-08d9b0c5b5e1','98281c1b-6161-4e74-a119-177b2a2fb176','54ab3a49-6d78-4922-bac1-94a722b9859a','AIR_AND_SPACE_FORCE',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('d1f5c22f-8b02-4b54-86d4-dcda2a8f9bd4','98281c1b-6161-4e74-a119-177b2a2fb176','3b6fe2e9-9116-4a70-8eef-96205205b0e3','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('a280e65c-445d-44e2-aa2b-503d24c3a7ab','98281c1b-6161-4e74-a119-177b2a2fb176','ba0017c8-5d48-4efe-8802-361d2f2bc16d','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('79897483-cbca-4362-abed-bf675bbf8131','98281c1b-6161-4e74-a119-177b2a2fb176','af65234a-8577-4f6d-a346-7d486e963287','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('323b43dc-af34-4cd8-9f9c-ba1b78855f51','98281c1b-6161-4e74-a119-177b2a2fb176','25f62695-ba9b-40c3-a7e4-0c078731123d','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('25c85fa9-ce16-4ba6-af9c-c8a32ed6d0d8','98281c1b-6161-4e74-a119-177b2a2fb176','b0422b13-4afe-443e-901f-fe652cde23a4','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('c909eb0a-9f3a-4358-b0c0-d2dacf73f43b','98281c1b-6161-4e74-a119-177b2a2fb176','ecbc1d89-9cf6-4f52-b453-3ba473a0ff4e','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('46321b67-aa7e-44bf-a016-f8c4c4531158','98281c1b-6161-4e74-a119-177b2a2fb176','4e85cb86-e9dd-4c7c-9677-e0a327ac895c','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('ef56af9f-a331-4ed8-8fee-4e720a75c5f6','98281c1b-6161-4e74-a119-177b2a2fb176','7f9dd10d-100e-4252-9786-706349f456ca','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('7202227a-dfda-49c3-ae5c-774068548dc0','98281c1b-6161-4e74-a119-177b2a2fb176','5451f8d7-60c5-4e22-bbf6-d9af8e6ace54','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('bd7a021c-49d6-4008-addf-0fa70381ecf1','98281c1b-6161-4e74-a119-177b2a2fb176','9e9cba85-7c39-4836-809d-70b54baf392e','AIR_AND_SPACE_FORCE',NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('f60dc26f-af8c-41d6-aa00-3ee5501a0028','98281c1b-6161-4e74-a119-177b2a2fb176','f7470914-cf48-43be-b431-a3ca2fe5b290','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('13013a3c-914e-47e6-8b79-a0b889b5cce6','98281c1b-6161-4e74-a119-177b2a2fb176','5a0d3cc1-b866-4bde-b67f-78d565facf3e','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('27bc120a-f39f-46a0-83a1-d4712aceb83f','98281c1b-6161-4e74-a119-177b2a2fb176','e629b95a-ec5b-4fc4-897f-e0d1050e1ec6','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('3326f502-ca09-4454-adc7-6aed1bcd329e','98281c1b-6161-4e74-a119-177b2a2fb176','d68a626f-935a-4eb1-ba9b-6829feeff91c','AIR_AND_SPACE_FORCE',NULL,true,now(),now()),
	 ('3d6fc709-f52d-4c5f-b36e-562013a6cce8','560eb530-39d7-4363-8436-5ffb4bbe8e12','ed496f6e-fd34-48d1-8586-63a1d305c49c',NULL,NULL,true,now(),now()),
	 ('ea812cd2-315c-45f6-9392-8fca73a27da1','560eb530-39d7-4363-8436-5ffb4bbe8e12','153b62b2-b1b8-4b9d-afa5-53df4150aba4',NULL,NULL,true,now(),now()),
	 ('7d2df4af-cec5-4d7a-a71f-9c52d7eba6a6','560eb530-39d7-4363-8436-5ffb4bbe8e12','0020441b-bc0c-436e-be05-b997ca6a853c',NULL,NULL,true,now(),now()),
	 ('4b885315-1139-40bd-b7ed-923927cd6c4c','560eb530-39d7-4363-8436-5ffb4bbe8e12','a8ae9bb9-e9ac-49b4-81dc-336b8a0dcb54',NULL,NULL,true,now(),now()),
	 ('72ee960d-4fa6-46ad-a047-b3f4d8e53ea1','58510246-905b-4989-b16e-04806904134b','1aa43046-8d6b-4249-88dc-b259d86c0cb8',NULL,NULL,true,now(),now()),
	 ('901ca376-ad0a-4f03-98f2-4fc07e59d7f3','58510246-905b-4989-b16e-04806904134b','84b50723-95fc-41d1-8115-a734c7e53f66',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('ec01d299-0d2b-4ef3-991c-ba26440e32b5','58510246-905b-4989-b16e-04806904134b','63dc3f78-235e-4b1c-b1db-459d7f5ae25f',NULL,NULL,true,now(),now()),
	 ('beb504cb-9a80-4020-81c1-031b73b60dce','58510246-905b-4989-b16e-04806904134b','4bd4c579-c163-4b1b-925a-d852d5a12642',NULL,NULL,true,now(),now()),
	 ('4fdbee93-810f-4a9c-a206-b089053c6469','58510246-905b-4989-b16e-04806904134b','ce7f1fc9-5b94-43cb-b398-b31cb6350d6a',NULL,NULL,true,now(),now()),
	 ('01c4f44c-ace7-42d5-81b0-a9ed81f56131','58510246-905b-4989-b16e-04806904134b','4a6a3260-e1ec-4d78-8c07-d89f1405ca16',NULL,NULL,true,now(),now()),
	 ('f765c3ec-e147-436b-8c72-31a6a5e6497d','58510246-905b-4989-b16e-04806904134b','658274b5-720c-4a70-ba9d-249614d85ebc',NULL,NULL,true,now(),now()),
	 ('160dd837-da5f-446d-968c-553cbedbf9c9','58510246-905b-4989-b16e-04806904134b','41f0736c-6d26-4e93-9668-860e9f0e222a',NULL,NULL,true,now(),now()),
	 ('41bd6cda-4d34-4adb-8ceb-f61fbb54ee8f','58510246-905b-4989-b16e-04806904134b','c2a8e8c3-dddc-4c0f-a23a-0b4e2c20af0d',NULL,NULL,true,now(),now()),
	 ('3546beae-dfda-4c69-af31-61f3f0a8fcc9','58510246-905b-4989-b16e-04806904134b','0f1e1c87-0497-4ee2-970d-21ac2d2155db',NULL,NULL,true,now(),now()),
	 ('25da4022-0e63-45ea-ab5f-fd7f94ec6e31','58510246-905b-4989-b16e-04806904134b','e40b411d-55f0-4470-83d0-0bbe11fa77dd',NULL,NULL,true,now(),now()),
	 ('4eae036f-7cb6-4592-bea0-fc19aadccb94','58510246-905b-4989-b16e-04806904134b','2632b4e5-c6cb-4e64-8924-0b7e4b1115ec',NULL,NULL,true,now(),now());
INSERT INTO gbloc_aors (id,jppso_regions_id,oconus_rate_area_id,department_indicator,shipment_type,is_active,created_at,updated_at) VALUES
	 ('47523bda-a2d5-4cff-a24c-63a15ad60fb1','58510246-905b-4989-b16e-04806904134b','1336deb9-5c87-409d-8051-4ab9f211eb29',NULL,NULL,true,now(),now()),
	 ('252400bc-1940-4521-9907-6876396118be','58510246-905b-4989-b16e-04806904134b','7a6d3b5b-81a6-4db5-b2ab-ecbfd6bd7941',NULL,NULL,true,now(),now()),
	 ('78d31ca4-969d-4aa8-afbd-2b8ac0f9b753','58510246-905b-4989-b16e-04806904134b','91b69254-5976-4839-a31d-972e9958d9cf',NULL,NULL,true,now(),now()),
	 ('e38f8607-d6d8-44b1-8e8d-f4fce2f1c747','560eb530-39d7-4363-8436-5ffb4bbe8e12','166f3629-79b9-451a-90a3-c43680929a2f',NULL,NULL,true,now(),now()),
	 ('113aa232-3276-4e94-a88e-2e0b8cf56450','560eb530-39d7-4363-8436-5ffb4bbe8e12','4fe05eb4-1b1c-4d4a-a185-0b039ac64835',NULL,NULL,true,now(),now());