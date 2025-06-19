--insert missing pay_grades
INSERT INTO public.pay_grades
(id, grade, grade_description, created_at, updated_at)
VALUES('9a892c59-48d5-4eba-b5f9-193716da8827', 'O_1', 'Officer Grade O_1', now(), now());

-- Army
INSERT INTO ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    ('c45746c6-0711-4ef1-ab76-c4b91f7d708f','7fa938ab-1c34-4666-a878-9b989c916d1a'::uuid,'ARMY','GEN','General',1,now(),now()),
    ('0019b47f-7963-4221-a9d2-84914f69a97d','1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80'::uuid,'ARMY','LTG','Lieutenant General',2,now(),now()),
    ('48968751-4857-475d-8bcd-0410d659f206','6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20'::uuid,'ARMY','MG','Major General',3,now(),now()),
    ('6c7c17e8-db98-4df1-8665-8784fccff12d','cf664124-9baf-4187-8f28-0908c0f0a5e0'::uuid,'ARMY','BG','Brigadier General',4,now(),now()),
    ('5ad1173d-4978-492f-aaa1-69ea826f7fae','455a112d-d1e0-4559-81e8-6df664638f70'::uuid,'ARMY','COL','Colonel',5,now(),now()),
    ('4a52e2f0-0936-46fc-b6c0-a00e888359cf','3bc4b197-7897-4105-80a1-39a0378d7730'::uuid,'ARMY','LTC','Lieutenant Colonel',6,now(),now()),
    ('ecb4ec4d-6ee6-4e15-a4d7-f344b171eb4e','e83d8f8d-f70b-4db1-99cc-dd983d2fd250'::uuid,'ARMY','MAJ','Major',7,now(),now()),
    ('b34c33a2-0d62-4d9a-a6f8-7fce746443e9','5658d67b-d510-4226-9e56-714403ba0f10'::uuid,'ARMY','CPT','Captain',8,now(),now()),
    ('8d653759-0a10-4703-9b52-bf305c48c470','d1b76a01-d8e4-4bd3-98ff-fa93ff7bc790'::uuid,'ARMY','1LT','First Lieutenant',9,now(),now()),
    ('2c2e9f55-17cc-42a5-aa9a-d4efd7f2529b','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'ARMY','2LT','Second Lieutenant',10,now(),now()),
    ('461aa3f5-8390-4fc4-b6ae-1611a4d188ce','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'ARMY','OC','Officer Candidate',10,now(),now()),
    ('d447b93a-d0ae-4943-af1c-39830f5e7278','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'ARMY','CDT','Cadet',10,now(),now()),
    ('36f84d0b-12cd-4355-af02-1bfb05e28b2a','ea8cb0e9-15ff-43b4-9e41-7168d01e7553'::uuid,'ARMY','CW5','Chief Warrant Officer 5',11,now(),now()),
    ('1b6a17fc-3659-4b02-b7ed-85bd5ac4a9a3','74db5649-cf66-4af8-939b-d3d7f1f6b7c6'::uuid,'ARMY','CW4','Chief Warrant Officer 4',12,now(),now()),
    ('07b939ea-ee1a-489e-a6ff-b8a05a5fc258','5a65fb1f-4245-4178-b6a7-cc504c9cbb37'::uuid,'ARMY','CW3','Chief Warrant Officer 3',13,now(),now()),
    ('3085a98c-8feb-486b-9736-ba655907319b','a687a2e1-488c-4943-b9d9-3d645a2712f4'::uuid,'ARMY','CW2','Chief Warrant Officer 2',14,now(),now()),
    ('1dcbdbc0-a83d-469c-93f7-bd69792f06a8','6badf8a0-b0ef-4e42-b827-7f63a3987a4b'::uuid,'ARMY','WO1','Warrant Officer 1',15,now(),now()),
    ('4e9ae713-e1cb-4423-ba9c-0c7e89ffcb08','a5fc8fd2-6f91-492b-abe2-2157d03ec990'::uuid,'ARMY','SGM','Sergeant Major',17,now(),now()),
    ('1ee79fdc-01c4-4496-aebb-fb96445f7ba0','a5fc8fd2-6f91-492b-abe2-2157d03ec990'::uuid,'ARMY','CSM','Command Sergeant Major',18,now(),now()),
    ('72b796de-844f-44f6-bb4c-9c5efd6117f7','1d909db0-602f-4724-bd43-8f90a6660460'::uuid,'ARMY','MSG','Master Sergeant',19,now(),now()),
    ('50f29cf0-fd75-452a-969e-a64ac19f3775','1d909db0-602f-4724-bd43-8f90a6660460'::uuid,'ARMY','1ST','1st Sergeant',20,now(),now()),
    ('9a227f63-904b-4c71-8a64-566cd427b2c5','523d57a1-529c-4dfd-8c33-9cb169fd29a0'::uuid,'ARMY','SFC','Sergeant First Class',21,now(),now()),
    ('1f3bab7c-06f3-4a12-901c-1a52dc10bd57','523d57a1-529c-4dfd-8c33-9cb169fd29a0'::uuid,'ARMY','PSG','Platoon Sergeant',21,now(),now()),
    ('0aaed778-0ad1-4e21-b1d4-7a42e337a696','541aec36-bd9f-4ad2-abb4-d9b63e29dc80'::uuid,'ARMY','SSG','Staff Sergeant',22,now(),now()),
    ('e5f85779-6bbf-4f37-9f6f-45c4e2b89310','3f142461-dca5-4a77-9295-92ee93371330'::uuid,'ARMY','SGT','Sergeant',23,now(),now()),
    ('fb8e41e3-0d26-414b-b208-087cd8d86f7b','bb55f37c-3165-46ba-ad3f-9a477f699990'::uuid,'ARMY','CPL','Corporal',24,now(),now()),
    ('468a0614-e87d-44b4-a50f-c2aa0dc987f5','bb55f37c-3165-46ba-ad3f-9a477f699990'::uuid,'ARMY','SPC','Specialist',25,now(),now()),
    ('ca930098-115a-494d-a8f7-32058af75494','862eb395-86d1-44af-ad47-dec44fbeda30'::uuid,'ARMY','PFC','Private First Class',27,now(),now()),
    ('dba4187c-1894-45f2-b892-f19427af5408','5f871c82-f259-43cc-9245-a6e18975dde0'::uuid,'ARMY','PV2','Private',28,now(),now()),
    ('577cd5c1-8c53-4242-b3d3-b01ca13cc371','6cb785d0-cabf-479a-a36d-a6aec294a4d0'::uuid,'ARMY','PV1','Private',29,now(),now()),
    ('f3bbb4fd-b39b-4965-91ae-919d47cc1103','6cb785d0-cabf-479a-a36d-a6aec294a4d0'::uuid,'ARMY','PVT','Private',30,now(),now());

-- USAF
INSERT INTO ranks (id,pay_grade_id,affiliation,rank_abbv,rank_name,rank_order,created_at,updated_at) VALUES
    ('560ce615-e31d-4800-9273-8344508d18e2','7fa938ab-1c34-4666-a878-9b989c916d1a'::uuid,'AIR_FORCE','Gen','General',1,now(),now()),
    ('d84742d4-cbbc-48ce-af03-3eada00aa145','1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80'::uuid,'AIR_FORCE','Lt Gen','Lieutenant General',2,now(),now()),
    ('5db514e4-5808-4211-aea7-866ac6a74dc4','6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20'::uuid,'AIR_FORCE','Maj Gen','Major General',3,now(),now()),
    ('daae28b0-9b9d-45bf-9ba9-b33e079f6105','cf664124-9baf-4187-8f28-0908c0f0a5e0'::uuid,'AIR_FORCE','Brig Gen','Brigadier General',4,now(),now()),
    ('f6157a08-c804-4d61-867d-806b4680a4de','455a112d-d1e0-4559-81e8-6df664638f70'::uuid,'AIR_FORCE','Col','Colonel',5,now(),now()),
    ('78eb194c-e7c3-4a47-a6c6-913de1a689f8','3bc4b197-7897-4105-80a1-39a0378d7730'::uuid,'AIR_FORCE','Lt Col','Lieutenant Colonel',6,now(),now()),
    ('299eb095-0dc9-4789-9549-3aaa147a2b81','e83d8f8d-f70b-4db1-99cc-dd983d2fd250'::uuid,'AIR_FORCE','Maj','Major',7,now(),now()),
    ('8850863d-c5e1-49dd-9c47-58a3289b122c','5658d67b-d510-4226-9e56-714403ba0f10'::uuid,'AIR_FORCE','Capt','Captain',8,now(),now()),
    ('c98ba41b-3fcd-474f-98ac-51764cc1f0e5','d1b76a01-d8e4-4bd3-98ff-fa93ff7bc790'::uuid,'AIR_FORCE','1st Lt','First Lieutenant',9,now(),now()),
    ('43b69e7b-99a3-488f-8cea-3abe63d6f20a','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'AIR_FORCE','AVC','Aviation Cadet',10,now(),now()),
    ('914db640-0de5-494d-b917-b4a44e022f4b','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'AIR_FORCE','2d Lt','Second Lieutenant',10,now(),now()),
    ('2cf8e36a-20fb-41fe-9268-d3d1f0219d1a','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'AIR_FORCE','AFC','Air Force Academy Cadet',10,now(),now()),
    ('6317aedf-73b6-4926-b763-44ff1ba0c00a','ea8cb0e9-15ff-43b4-9e41-7168d01e7553'::uuid,'AIR_FORCE','CWO5','Chief Warrant Officer 5',11,now(),now()),
    ('3b55fae1-bd2e-4e70-b154-d0df35cd706a','74db5649-cf66-4af8-939b-d3d7f1f6b7c6'::uuid,'AIR_FORCE','CWO4','Chief Warrant Officer 4',12,now(),now()),
    ('13cf90df-c8a9-47ae-bd78-e0854db429aa','5a65fb1f-4245-4178-b6a7-cc504c9cbb37'::uuid,'AIR_FORCE','CWO3','Chief Warrant Officer 3',13,now(),now()),
    ('5af407d9-5586-42ec-a1e7-f9d11283641a','a687a2e1-488c-4943-b9d9-3d645a2712f4'::uuid,'AIR_FORCE','CWO2','Chief Warrant Officer 2',14,now(),now()),
    ('6a819381-85d6-45fe-ad61-88aeb1d5f91f','6badf8a0-b0ef-4e42-b827-7f63a3987a4b'::uuid,'AIR_FORCE','WO','Warrant Officer 1',15,now(),now()),
    ('85578266-a86c-42ce-b740-e62614361114','a5fc8fd2-6f91-492b-abe2-2157d03ec990'::uuid,'AIR_FORCE','CMSgt','Chief Master Sergeant',17,now(),now()),
    ('dda2b553-99be-439b-b088-b7608b48eff0','1d909db0-602f-4724-bd43-8f90a6660460'::uuid,'AIR_FORCE','SMSgt','Senior Master Sergeant',18,now(),now()),
    ('43dd0d76-1f2f-45fb-a73e-404fe2ab93f0','523d57a1-529c-4dfd-8c33-9cb169fd29a0'::uuid,'AIR_FORCE','MSgt','Master Sergeant',19,now(),now()),
    ('0472a25d-b1a0-451c-9895-110dfe44496a','541aec36-bd9f-4ad2-abb4-d9b63e29dc80'::uuid,'AIR_FORCE','TSgt','Technical Sergeant',20,now(),now()),
    ('ae9f9d91-b049-4f60-bdc9-e441a7b3cb30','3f142461-dca5-4a77-9295-92ee93371330'::uuid,'AIR_FORCE','SSgt','Staff Sergeant',21,now(),now()),
    ('753f82f9-27e1-4ee7-9b57-bfef3c83656b','bb55f37c-3165-46ba-ad3f-9a477f699990'::uuid,'AIR_FORCE','SrA','Senior Airman',22,now(),now()),
    ('3aca9ba8-3b84-42bf-8f2f-5ef02587ba89','862eb395-86d1-44af-ad47-dec44fbeda30'::uuid,'AIR_FORCE','A1C','Airman First Class',23,now(),now()),
    ('cb0ee2b8-e852-40fe-b972-2730b53860c7','5f871c82-f259-43cc-9245-a6e18975dde0'::uuid,'AIR_FORCE','Amn','Airman',24,now(),now()),
    ('f6dbd496-8f71-487b-a432-55b60967f474','6cb785d0-cabf-479a-a36d-a6aec294a4d0'::uuid,'AIR_FORCE','AB','Airman Basic',25,now(),now());
-- Marines
INSERT INTO ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    ('535c23c3-95e6-4795-9c5f-52e1dceae375','7fa938ab-1c34-4666-a878-9b989c916d1a'::uuid,'MARINES','Gen','General',1,now(),now()),
    ('e15498d9-0015-457b-8a10-fd42e6ed02b1','1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80'::uuid,'MARINES','LtGen','Lieutenant General',2,now(),now()),
    ('602a512b-fee8-41c6-a65f-042318a3924b','6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20'::uuid,'MARINES','MajGen','Major General',3,now(),now()),
    ('ac532793-77d8-40c0-b25e-a98249520bc7','cf664124-9baf-4187-8f28-0908c0f0a5e0'::uuid,'MARINES','BGen','Brigadier General',4,now(),now()),
    ('0fe1d0c4-a734-472d-8b16-7dfbe881aa5a','455a112d-d1e0-4559-81e8-6df664638f70'::uuid,'MARINES','Col','Colonel',5,now(),now()),
    ('63080039-cf5c-4faf-9354-e0eff69b2133','3bc4b197-7897-4105-80a1-39a0378d7730'::uuid,'MARINES','LtCol','Lieutenant Colonel',6,now(),now()),
    ('66c6f1e8-4d82-48ff-ae4a-622ec54a2b75','e83d8f8d-f70b-4db1-99cc-dd983d2fd250'::uuid,'MARINES','Maj','Major',7,now(),now()),
    ('53c41df9-dd4b-4821-83ce-3927e39cf9d4','5658d67b-d510-4226-9e56-714403ba0f10'::uuid,'MARINES','Capt','Captain',8,now(),now()),
    ('5a21a044-207c-4c4d-bfe2-7a73c717fed6','d1b76a01-d8e4-4bd3-98ff-fa93ff7bc790'::uuid,'MARINES','1stLt','First Lieutenant',9,now(),now()),
    ('d01f62c6-14d6-4f5d-8b42-f2319df6730d','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'MARINES','2ndLt','Second Lieutenant',10,now(),now()),
    ('deedacab-55f6-4a0a-8042-9a6b45d4f819','ea8cb0e9-15ff-43b4-9e41-7168d01e7553'::uuid,'MARINES','CWO5','Chief Warrant Officer 5',11,now(),now()),
    ('fd522412-7b6d-410b-9629-1acba8b63108','74db5649-cf66-4af8-939b-d3d7f1f6b7c6'::uuid,'MARINES','CWO4','Chief Warrant Officer 4',12,now(),now()),
    ('3821df00-1b7a-449d-b529-ac789e72f52c','5a65fb1f-4245-4178-b6a7-cc504c9cbb37'::uuid,'MARINES','CWO3','Chief Warrant Officer 3',13,now(),now()),
    ('62197a37-3b56-430d-a084-0b73170b3003','a687a2e1-488c-4943-b9d9-3d645a2712f4'::uuid,'MARINES','CWO2','Chief Warrant Officer 2',14,now(),now()),
    ('8df6c604-f085-4aa3-885b-6522b88125bd','6badf8a0-b0ef-4e42-b827-7f63a3987a4b'::uuid,'MARINES','WO','Warrant Officer 1',15,now(),now()),
    ('97c779cf-d554-4d5e-9b39-6d38fc382ef4','a5fc8fd2-6f91-492b-abe2-2157d03ec990'::uuid,'MARINES','SgtMaj','Sergeant Major',17,now(),now()),
    ('dcb640a6-c0ed-442a-b02e-e2f39105f3c6','a5fc8fd2-6f91-492b-abe2-2157d03ec990'::uuid,'MARINES','MGySgt','Master Gunnery Sergeant',18,now(),now()),
    ('d6404c12-3e86-4db7-879b-89e6856142a5','1d909db0-602f-4724-bd43-8f90a6660460'::uuid,'MARINES','MSGgt','Master Sergeant',19,now(),now()),
    ('3258e3ba-b151-42d4-9d2a-3f75fa4fd1e6','1d909db0-602f-4724-bd43-8f90a6660460'::uuid,'MARINES','1st Sgt','1st Sergeant',20,now(),now()),
    ('87570382-a2ed-447a-a3c6-cdfda1f8d753','523d57a1-529c-4dfd-8c33-9cb169fd29a0'::uuid,'MARINES','GySgt','Gunnery Sergeant',21,now(),now()),
    ('c2557065-3e79-4b4c-b0c7-5cc18243cb5a','541aec36-bd9f-4ad2-abb4-d9b63e29dc80'::uuid,'MARINES','SSgt','Staff Sergeant',22,now(),now()),
    ('eb404355-ed3d-4bf0-81cb-fda983829830','3f142461-dca5-4a77-9295-92ee93371330'::uuid,'MARINES','Sgt','Sergeant',23,now(),now()),
    ('d2f5a512-50fd-46d4-9bfa-745f6f73f856','bb55f37c-3165-46ba-ad3f-9a477f699990'::uuid,'MARINES','Cpl','Corporal',24,now(),now()),
    ('4b99e53c-d103-45bf-8c83-b15b456f8f29','862eb395-86d1-44af-ad47-dec44fbeda30'::uuid,'MARINES','LCpl','Lance Corporal',26,now(),now()),
    ('6ffb7477-6078-49ed-9103-c75bdc6d0818','5f871c82-f259-43cc-9245-a6e18975dde0'::uuid,'MARINES','PFC','Private First Class',27,now(),now()),
    ('ceff4230-04ec-4ce1-93d5-920c99af4991','6cb785d0-cabf-479a-a36d-a6aec294a4d0'::uuid,'MARINES','PVT','Private',28,now(),now());

-- Navy
INSERT INTO ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    ('1ad0a9ac-f886-4c36-9c7a-43e4f993a9df','7fa938ab-1c34-4666-a878-9b989c916d1a'::uuid,'NAVY','ADM','Admiral',1,now(),now()),
    ('4f111a4b-651c-4f96-ad8a-3ec65777e6c0','1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80'::uuid,'NAVY','VADM','Vice Admiral',2,now(),now()),
    ('1b5913ef-6675-4e11-9871-077ee19eacaa','6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20'::uuid,'NAVY','RADM','Rear Admiral (Upper Half)',3,now(),now()),
    ('581d1318-29b0-4f11-933a-e716bc3fc8fd','cf664124-9baf-4187-8f28-0908c0f0a5e0'::uuid,'NAVY','RDML','Rear Admiral (Lower Half)',4,now(),now()),
    ('a28777d4-181f-4d78-869d-f21d895b671e','455a112d-d1e0-4559-81e8-6df664638f70'::uuid,'NAVY','CAPT','Captain',5,now(),now()),
    ('a216481b-ceee-4175-bcc3-1afeec36f5fa','3bc4b197-7897-4105-80a1-39a0378d7730'::uuid,'NAVY','CDR','Commander',6,now(),now()),
    ('4616c1c9-6794-4009-a8ff-e7d5b4a5f9e4','e83d8f8d-f70b-4db1-99cc-dd983d2fd250'::uuid,'NAVY','LCDR','Lieutenant Commander',7,now(),now()),
    ('6d1c27bc-245b-452a-9cd3-6da5540cf6db','5658d67b-d510-4226-9e56-714403ba0f10'::uuid,'NAVY','LT','Lieutenant',8,now(),now()),
    ('7ff96b44-c9b2-4db8-91a0-94671f987faf','d1b76a01-d8e4-4bd3-98ff-fa93ff7bc790'::uuid,'NAVY','LTJG','Lieutenant JG',9,now(),now()),
    ('ffc413c3-59df-4c5c-848f-b530a1bee691','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'NAVY','ENS','Ensign',10,now(),now()),
    ('0dc31054-0939-44ff-80c4-114b80f40895','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'NAVY','MID','Midshipman',10,now(),now()),
    ('d5a88410-076d-42aa-9889-5e90b87f4821','ea8cb0e9-15ff-43b4-9e41-7168d01e7553'::uuid,'NAVY','CWO5','Chief Warrant Officer 5',11,now(),now()),
    ('9d40f511-123a-4e22-8fb9-f5107c3dc63e','74db5649-cf66-4af8-939b-d3d7f1f6b7c6'::uuid,'NAVY','CWO4','Chief Warrant Officer 4',12,now(),now()),
    ('0739d2e0-267c-4b11-abd9-0baec2b44cc7','5a65fb1f-4245-4178-b6a7-cc504c9cbb37'::uuid,'NAVY','CWO3','Chief Warrant Officer 3',13,now(),now()),
    ('772db820-df8e-471b-a77e-65704e701862','a687a2e1-488c-4943-b9d9-3d645a2712f4'::uuid,'NAVY','CWO2','Chief Warrant Officer 2',14,now(),now()),
    ('298a491b-24ba-48a7-9e1b-5593dbd7387c','6badf8a0-b0ef-4e42-b827-7f63a3987a4b'::uuid,'NAVY','WO','Warrant Officer 1',15,now(),now()),
    ('a56fec50-4639-426f-8abd-85eedf485723','a5fc8fd2-6f91-492b-abe2-2157d03ec990'::uuid,'NAVY','MCPO','Master Chief Petty Officer',17,now(),now()),
    ('e5d90322-b55d-4156-a9a7-a8331b2ff92f','1d909db0-602f-4724-bd43-8f90a6660460'::uuid,'NAVY','SCPO','Senior Chief Petty Officer',18,now(),now()),
    ('34ebdf38-d762-47ad-b521-e3a56e9716f2','523d57a1-529c-4dfd-8c33-9cb169fd29a0'::uuid,'NAVY','CPO','Chief Petty Officer',19,now(),now()),
    ('64b449ea-145a-4097-b25c-cff7fe61bf21','541aec36-bd9f-4ad2-abb4-d9b63e29dc80'::uuid,'NAVY','PO1','Petty Officer First Class',20,now(),now()),
    ('55e27ead-15e9-479b-bf84-fc57d2cecb7f','3f142461-dca5-4a77-9295-92ee93371330'::uuid,'NAVY','PO2','Petty Officer Second Class',21,now(),now()),
    ('24b49ba8-1ce9-4456-bbdc-64ae9c08c156','bb55f37c-3165-46ba-ad3f-9a477f699990'::uuid,'NAVY','PO3','Petty Officer Third Class',22,now(),now()),
    ('6b8f9f9d-d018-48b6-8339-9ab28060e838','862eb395-86d1-44af-ad47-dec44fbeda30'::uuid,'NAVY','SN','Seaman',24,now(),now()),
    ('4d025a52-ca65-4dfe-9ba7-58f253a85660','5f871c82-f259-43cc-9245-a6e18975dde0'::uuid,'NAVY','SA','Seaman Apprentice',25,now(),now()),
    ('2839c1f7-cf58-4603-9806-bf8c924949b8','6cb785d0-cabf-479a-a36d-a6aec294a4d0'::uuid,'NAVY','SR','Seaman Recruit',26,now(),now());

-- Coast Guard
INSERT INTO ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    ('377cde5f-6e62-400e-8749-fd7526efd60f','7fa938ab-1c34-4666-a878-9b989c916d1a'::uuid,'COAST_GUARD','ADM','Admiral',1,now(),now()),
    ('87d7a277-ba7c-48a0-99ec-6da4bd321f3a','1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80'::uuid,'COAST_GUARD','VADM','Vice Admiral',2,now(),now()),
    ('ff535837-7522-4096-b089-e58708401002','6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20'::uuid,'COAST_GUARD','RADM','Rear Admiral (Upper Half)',3,now(),now()),
    ('dfbbadcb-4075-4b4a-b62d-c2b6bbf91a7a','cf664124-9baf-4187-8f28-0908c0f0a5e0'::uuid,'COAST_GUARD','RDML','Rear Admiral (Lower Half)',4,now(),now()),
    ('791962d3-7ad9-466d-9b06-5d3045570884','455a112d-d1e0-4559-81e8-6df664638f70'::uuid,'COAST_GUARD','CAPT','Captain',5,now(),now()),
    ('786807ed-e002-4c04-8c0b-8332cfbab600','3bc4b197-7897-4105-80a1-39a0378d7730'::uuid,'COAST_GUARD','CDR','Commander',6,now(),now()),
    ('685cfe77-a633-4a9a-8a17-adb9fa71e5f0','e83d8f8d-f70b-4db1-99cc-dd983d2fd250'::uuid,'COAST_GUARD','LCDR','Lieutenant Commander',7,now(),now()),
    ('4f0dafd9-038c-4eb7-afd4-866abf3bd056','5658d67b-d510-4226-9e56-714403ba0f10'::uuid,'COAST_GUARD','LT','Lieutenant',8,now(),now()),
    ('03e42915-41e9-4246-9aec-26662148e6c6','d1b76a01-d8e4-4bd3-98ff-fa93ff7bc790'::uuid,'COAST_GUARD','LTJG','Lieutenant JG',9,now(),now()),
    ('9f86e4ba-29f3-47db-b9a2-da5de1d585a4','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'COAST_GUARD','MID','Midshipman',10,now(),now()),
    ('9d2342de-b299-45d1-b521-96e0d1435ac5','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'COAST_GUARD','ENS','Ensign',10,now(),now()),
    ('446c41b6-04bd-4847-80ca-95a6dc5430fd','74db5649-cf66-4af8-939b-d3d7f1f6b7c6'::uuid,'COAST_GUARD','CWO4','Chief Warrant Officer 4',11,now(),now()),
    ('ee2fe271-e68f-4958-bc18-f7e03b05b4b0','5a65fb1f-4245-4178-b6a7-cc504c9cbb37'::uuid,'COAST_GUARD','CWO3','Chief Warrant Officer 3',12,now(),now()),
    ('d12ac35b-07e8-485a-968c-2a1a912e14d4','a687a2e1-488c-4943-b9d9-3d645a2712f4'::uuid,'COAST_GUARD','CWO2','Chief Warrant Officer 2',13,now(),now()),
    ('c0cf6c2e-5bbb-4c16-9fcd-3e0a7bb41700','a5fc8fd2-6f91-492b-abe2-2157d03ec990'::uuid,'COAST_GUARD','CPM','Master Chief Petty Officer',16,now(),now()),
    ('d472e7ac-dc61-48df-9809-5f47cf46cebd','1d909db0-602f-4724-bd43-8f90a6660460'::uuid,'COAST_GUARD','CPS','Senior Chief Petty Officer',17,now(),now()),
    ('f11e1111-993b-46ad-8180-03648a018e4f','523d57a1-529c-4dfd-8c33-9cb169fd29a0'::uuid,'COAST_GUARD','CPO','Chief Petty Officer',18,now(),now()),
    ('7dc6066d-22b9-4c01-bf6b-0d895f3f206b','541aec36-bd9f-4ad2-abb4-d9b63e29dc80'::uuid,'COAST_GUARD','PO1','Petty Officer First Class',19,now(),now()),
    ('14c26a18-0d50-4581-9be6-a38c4de2a6a9','3f142461-dca5-4a77-9295-92ee93371330'::uuid,'COAST_GUARD','PO2','Petty Officer Second Class',20,now(),now()),
    ('b95491ab-b51a-4cca-95f6-aca9edaf9c53','bb55f37c-3165-46ba-ad3f-9a477f699990'::uuid,'COAST_GUARD','PO3','Petty Officer Third Class',21,now(),now()),
    ('6a54692b-f745-4aae-8e9e-b770eb794be9','862eb395-86d1-44af-ad47-dec44fbeda30'::uuid,'COAST_GUARD','SN','Seaman',23,now(),now()),
    ('ece77940-430d-4ba1-8057-8933b841ea41','5f871c82-f259-43cc-9245-a6e18975dde0'::uuid,'COAST_GUARD','SA','Seaman Apprentice',24,now(),now()),
    ('2587c794-f20d-42c5-84c9-3b4b3fa788f2','6cb785d0-cabf-479a-a36d-a6aec294a4d0'::uuid,'COAST_GUARD','SR','Seaman Recruit',25,now(),now());

-- Space Force
INSERT INTO ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    ('b0d56e0d-c5e3-480a-b839-b20f59c03884','7fa938ab-1c34-4666-a878-9b989c916d1a'::uuid,'SPACE_FORCE','Gen','General',1,now(),now()),
    ('60163e2c-e822-4e85-ba0c-595e5b5a52df','1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80'::uuid,'SPACE_FORCE','Lt Gen','Lieutenant General',2,now(),now()),
    ('1c4feb55-5aa5-4057-aba2-8f4cdf0c2c2c','6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20'::uuid,'SPACE_FORCE','Maj Gen','Major General',3,now(),now()),
    ('90c2dce3-fc9a-407f-8868-5d8e45757826','cf664124-9baf-4187-8f28-0908c0f0a5e0'::uuid,'SPACE_FORCE','Brig Gen','Brigadier General',4,now(),now()),
    ('9327cd6c-ebf1-44dc-98aa-ee6ed2604af8','455a112d-d1e0-4559-81e8-6df664638f70'::uuid,'SPACE_FORCE','Col','Colonel',5,now(),now()),
    ('a907574a-0841-4e00-a08d-ccaa1cb43feb','3bc4b197-7897-4105-80a1-39a0378d7730'::uuid,'SPACE_FORCE','Lt Col','Lieutenant Colonel',6,now(),now()),
    ('12a68180-8299-4ea0-b8a6-50b79ac6f74e','e83d8f8d-f70b-4db1-99cc-dd983d2fd250'::uuid,'SPACE_FORCE','Maj','Major',7,now(),now()),
    ('dc82cc34-8fd1-49ff-8bf7-31793fe93316','5658d67b-d510-4226-9e56-714403ba0f10'::uuid,'SPACE_FORCE','Capt','Captain',8,now(),now()),
    ('103d519d-1b5d-4226-8c71-d34f366a0583','d1b76a01-d8e4-4bd3-98ff-fa93ff7bc790'::uuid,'SPACE_FORCE','1st Lt','First Lieutenant',9,now(),now()),
    ('aa604956-fc8b-4f97-a66f-16d6d875c564','9a892c59-48d5-4eba-b5f9-193716da8827'::uuid,'SPACE_FORCE','2d Lt','Second Lieutenant',10,now(),now()),
    ('903c8212-19b7-40e4-9f4b-150f6f76b8aa','a5fc8fd2-6f91-492b-abe2-2157d03ec990'::uuid,'SPACE_FORCE','CMSgt','Chief Master Sergeant',12,now(),now()),
    ('0cb0fe2c-7d45-4167-bce7-2ae85b6debd8','1d909db0-602f-4724-bd43-8f90a6660460'::uuid,'SPACE_FORCE','SMSgt','Senior Master Sergeant',13,now(),now()),
    ('0d35702e-5cd1-4051-a049-45e32a9116f0','523d57a1-529c-4dfd-8c33-9cb169fd29a0'::uuid,'SPACE_FORCE','MSgt','Master Sergeant',14,now(),now()),
    ('d5a09a55-ffa3-46fb-941e-7c37b762970e','541aec36-bd9f-4ad2-abb4-d9b63e29dc80'::uuid,'SPACE_FORCE','TSgt','Technical Sergeant',15,now(),now()),
    ('f3628803-59be-456d-87c3-fa0df971d04c','3f142461-dca5-4a77-9295-92ee93371330'::uuid,'SPACE_FORCE','Sgt','Sergeant',16,now(),now()),
    ('73b61e4b-54c4-435a-9a7b-f7d6e835c06e','bb55f37c-3165-46ba-ad3f-9a477f699990'::uuid,'SPACE_FORCE','Spc4','Specialist 4',17,now(),now()),
    ('dbd26213-fd1c-4078-9a4e-a321f4c3060a','862eb395-86d1-44af-ad47-dec44fbeda30'::uuid,'SPACE_FORCE','Spc3','Specialist 3',18,now(),now()),
    ('7a66c9db-5c87-4c15-b584-647dbfa5f6eb','5f871c82-f259-43cc-9245-a6e18975dde0'::uuid,'SPACE_FORCE','Spc2','Specialist 2',19,now(),now()),
    ('61c647fa-5325-45b9-8d6f-30a2aaa06308','6cb785d0-cabf-479a-a36d-a6aec294a4d0'::uuid,'SPACE_FORCE','Spc1','Specialist 1',20,now(),now());

INSERT INTO ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    ('d3aa6931-7858-4123-be0b-f3242a49e9f7', '9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9'::uuid,'AIR_FORCE','CIV','Civilian',null,now(),now()),
    ('0be75ca3-5226-447b-ad5f-d73205946bcb', '9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9'::uuid,'ARMY','CIV','Civilian',null,now(),now()),
    ('4af7df34-a5a3-448e-93eb-7bed05704cd0', '9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9'::uuid,'COAST_GUARD','CIV','Civilian',null,now(),now()),
    ('676ae8f9-0ca6-4919-b502-ef46c1cfaa48', '9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9'::uuid,'MARINES','CIV','Civilian',null,now(),now()),
    ('4a1dfa4b-f051-4c2b-a7d1-7279171e02d2', '9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9'::uuid,'NAVY','CIV','Civilian',null,now(),now()),
    ('3a4fcdf6-9ddf-4e31-b4da-d9823cfbb9b6', '9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9'::uuid,'SPACE_FORCE','CIV','Civilian',null,now(),now()),
    ('6557dfe3-35b7-44c5-adee-43a29dfcf289', '9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9'::uuid,'OTHER','CIV','Civilian',null,now(),now());

INSERT INTO ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    ('4e75f79c-2a2a-442e-a7c9-bd3bcf57e3c4','911208cc-3d13-49d6-9478-b0a3943435c0'::uuid,'AIR_FORCE','CMSAF','Chief Master Sergeant of the Air Force',16,now(),now()),
    ('2a1b2c87-9861-4900-9d1f-6844a3dd4e9d','911208cc-3d13-49d6-9478-b0a3943435c0'::uuid,'ARMY','SMA','Sergeant Major of the Army',16,now(),now()),
    ('621104d4-e717-414c-8016-b7edf0683ad7','911208cc-3d13-49d6-9478-b0a3943435c0'::uuid,'COAST_GUARD','MCPOCG','Master Chief Petty Officer of the Coast Guard',15,now(),now()),
    ('028097e4-2b35-4a20-b6db-3b436dc56f45','911208cc-3d13-49d6-9478-b0a3943435c0'::uuid,'MARINES','SgtMajMC','Sergeant Major of the Marine Corps',16,now(),now()),
    ('41023a85-bbbf-4ee9-a7ce-ae5bab910cce','911208cc-3d13-49d6-9478-b0a3943435c0'::uuid,'NAVY','MCPON','Master Chief Petty Officer of the Navy',16,now(),now()),
    ('b333751f-4831-4068-91ca-25fbe60c4b3a','911208cc-3d13-49d6-9478-b0a3943435c0'::uuid,'SPACE_FORCE','CMASSF','Chief Master Sergeant of the Space Force',11,now(),now());

INSERT INTO ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    ('5324bc17-9c46-486e-9d7f-1ed9f1f76dcb', '9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9', 'CIVILIAN', 'CIV', 'Civilian', 1, now(), now());


--add pay_grade_rank_id to orders table
alter table orders drop if exists rank_id;
alter table ranks drop constraint if exists rank_id;

ALTER TABLE orders
   ADD rank_id uuid
   	CONSTRAINT fk_orders_rank_id REFERENCES ranks
(id);

--update rank_id in orders where grade:rank is 1:1
do '
declare
	i record;
	v_count int;
begin
	for i in (
		select pg.id pay_grade_id, pg.grade, o.id orders_id, sm.affiliation
			from pay_grades pg, orders o, service_members sm
			where pg.grade = o.grade
			  and o.service_member_id = sm.id)
	loop
		select count(*) into v_count
		  from ranks
		 where pay_grade_id = i.pay_grade_id
		   and affiliation = i.affiliation;
		if v_count = 1 then	--if 1 rank for pay grade then assign rank_id
			update orders o
			   set rank_id = p.id
			  from ranks p
			 where o.id = i.orders_id
			   and p.pay_grade_id = i.pay_grade_id
			   and p.affiliation = i.affiliation
			   and o.rank_id is null;
		end if;
	end loop;
end ';

--update grade and rank on orders
update orders
   set grade = 'O_1',
       rank_id = '2cf8e36a-20fb-41fe-9268-d3d1f0219d1a' --O_1/AFC
 where grade in ('ACADEMY_CADET','O_1_ACADEMY_GRADUATE')
   and service_member_id in (select id from service_members where affiliation = 'AIR_FORCE');

update orders
   set grade = 'O_1',
       rank_id = 'd447b93a-d0ae-4943-af1c-39830f5e7278' --O_1/CDT
 where grade in ('ACADEMY_CADET','O_1_ACADEMY_GRADUATE')
   and service_member_id in (select id from service_members where affiliation = 'ARMY');

update orders
   set grade = 'O_1',
       rank_id = '0dc31054-0939-44ff-80c4-114b80f40895' --O_1/MID
 where grade = 'MIDSHIPMAN'
   and service_member_id in (select id from service_members where affiliation = 'NAVY');

--remove unused pay grades
delete from pay_grades where grade in
('O_1_ACADEMY_GRADUATE',
'ACADEMY_CADET',
'MIDSHIPMAN',
'AVIATION_CADET');

INSERT INTO hhg_allowances (
        id,
        pay_grade_id,
        total_weight_self,
        total_weight_self_plus_dependents,
        pro_gear_weight,
        pro_gear_weight_spouse
    )
VALUES (
        '9a892c59-48d5-4eba-b5f9-193716da8827',
        (
            SELECT id
            FROM pay_grades
            WHERE grade = 'O_1'
        ),
        10000,
        12000,
        2000,
        500
    );
