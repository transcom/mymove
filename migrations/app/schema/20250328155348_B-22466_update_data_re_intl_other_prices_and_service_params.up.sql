------------------------------------------------------------------------------------------------
-- update re_intl_other_prices with data for new re_intl_other_prices.is_less_50_miles field
------------------------------------------------------------------------------------------------

--set existing IOPSIT prices as <= 50 miles
update re_intl_other_prices set is_less_50_miles = true where service_id = '6f4f6e31-0675-4051-b659-89832259f390';

--set existing IDDSIT prices as not <= 50 miles
update re_intl_other_prices set is_less_50_miles = false where service_id = '28389ee1-56cf-400c-aa52-1501ecdd7c69';

--insert IDDSIT recs for <= 50 miles
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('1db01885-d7a1-42d9-a3a7-d2d6d7f1fe8f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'a761a482-2929-4345-8027-3c6258f0c8dd'::uuid,1817,true,now(),now()),
	 ('ace0f878-58ef-4758-92dc-cfa37c546a5d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'a761a482-2929-4345-8027-3c6258f0c8dd'::uuid,1867,true,now(),now()),
	 ('299677d6-04bc-48d0-83f3-e0dc2add5ead','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'91eb2878-0368-4347-97e3-e6caa362d878'::uuid,1833,true,now(),now()),
	 ('892a4dff-04f1-4b5f-9aef-9fc83f3881b9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'91eb2878-0368-4347-97e3-e6caa362d878'::uuid,1883,true,now(),now()),
	 ('65558666-c6af-4172-aba1-1701384fa2a0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'635e4b79-342c-4cfc-8069-39c408a2decd'::uuid,3302,true,now(),now()),
	 ('62f8c644-7e1b-4f6e-a4bc-7036023b5111','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'635e4b79-342c-4cfc-8069-39c408a2decd'::uuid,3302,true,now(),now()),
	 ('52b8cf03-bf84-44eb-925a-b0fd7f4f8ae0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'9bb87311-1b29-4f29-8561-8a4c795654d4'::uuid,3486,true,now(),now()),
	 ('8ea5645d-6ba1-42e4-a1a9-4d1e6e260e96','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'9bb87311-1b29-4f29-8561-8a4c795654d4'::uuid,3486,true,now(),now()),
	 ('c054f6ee-b6d8-4234-8d37-382d97893e63','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'b80a00d4-f829-4051-961a-b8945c62c37d'::uuid,2830,true,now(),now()),
	 ('e1d3b5f4-8863-4439-91ef-ab4bb12ccb9e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'b80a00d4-f829-4051-961a-b8945c62c37d'::uuid,2830,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('ab39f0b3-f322-4cf9-9c48-4c7eda4bf259','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'5a27e806-21d4-4672-aa5e-29518f10c0aa'::uuid,2830,true,now(),now()),
	 ('b59248f7-1328-4dc6-8cb6-4521be0df368','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'5a27e806-21d4-4672-aa5e-29518f10c0aa'::uuid,2830,true,now(),now()),
	 ('bd279a0e-d6c7-4f9a-a12b-be510c6c05a9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'71755cc7-0844-4523-a0ac-da9a1e743ad1'::uuid,731,true,now(),now()),
	 ('d6978875-8589-4174-8100-1471138714ee','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'71755cc7-0844-4523-a0ac-da9a1e743ad1'::uuid,731,true,now(),now()),
	 ('09f20b44-42be-4d36-8903-c761874a7de6','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'829d8b45-19c1-49a3-920c-cc0ae14e8698'::uuid,626,true,now(),now()),
	 ('623ec13f-4579-46e0-8eda-cbf81f0887ea','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'829d8b45-19c1-49a3-920c-cc0ae14e8698'::uuid,626,true,now(),now()),
	 ('86b78bf2-144b-4aae-ba34-25a8b00bbbc5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'03dd5854-8bc3-4b56-986e-eac513cc1ec0'::uuid,626,true,now(),now()),
	 ('3676c15f-3d25-4de1-a9ec-e2f7f1d02fcf','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'03dd5854-8bc3-4b56-986e-eac513cc1ec0'::uuid,626,true,now(),now()),
	 ('d45f3aa1-453f-40ab-bf9e-2be724aeb2ca','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'ee0ffe93-32b3-4817-982e-6d081da85d28'::uuid,626,true,now(),now()),
	 ('a2ee779b-076b-48a5-90e6-648f69395342','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'ee0ffe93-32b3-4817-982e-6d081da85d28'::uuid,626,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('3f3ea5e8-4d82-40d4-b3fb-d0b58be5e8be','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'0506bf0f-bc1c-43c7-a75f-639a1b4c0449'::uuid,626,true,now(),now()),
	 ('d63dc23d-ac7d-496d-a583-c8b75b2ad5c9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'0506bf0f-bc1c-43c7-a75f-639a1b4c0449'::uuid,626,true,now(),now()),
	 ('a4a8250d-8fa5-4bcb-8e12-a93f091831a8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'6a0f9a02-b6ba-4585-9d7a-6959f7b0248f'::uuid,626,true,now(),now()),
	 ('00c775e7-5504-403f-87ff-d36d04069450','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'6a0f9a02-b6ba-4585-9d7a-6959f7b0248f'::uuid,626,true,now(),now()),
	 ('d6c23404-2617-412d-9718-492d1d2687bd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'c9036eb8-84bb-4909-be20-0662387219a7'::uuid,626,true,now(),now()),
	 ('99b71b47-b366-4809-bfba-0c57ab8054c5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'c9036eb8-84bb-4909-be20-0662387219a7'::uuid,626,true,now(),now()),
	 ('38c56cd3-3e1d-4e89-8d8a-5293a2e601e4','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'fe76b78f-67bc-4125-8f81-8e68697c136d'::uuid,626,true,now(),now()),
	 ('a9eba65a-ad89-4038-b2e4-94373b98de23','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'fe76b78f-67bc-4125-8f81-8e68697c136d'::uuid,626,true,now(),now()),
	 ('8541beee-fcfc-4a8c-a582-61b74992d851','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'2b1d1842-15f8-491a-bdce-e5f9fea947e7'::uuid,626,true,now(),now()),
	 ('b077a60a-f3ce-4e13-af65-efaec688b770','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'2b1d1842-15f8-491a-bdce-e5f9fea947e7'::uuid,626,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('0c9a68fa-7df8-4de6-b1b0-0af6f2fb39e5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'612c2ce9-39cc-45e6-a3f1-c6672267d392'::uuid,626,true,now(),now()),
	 ('0a116af3-df79-4a2e-873a-9502fcda5baf','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'612c2ce9-39cc-45e6-a3f1-c6672267d392'::uuid,626,true,now(),now()),
	 ('f322afb8-6b4e-4809-9f3e-b2fb55806711','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'7d0fc5a1-719b-4070-a740-fe387075f0c3'::uuid,626,true,now(),now()),
	 ('5ccb7f6b-976c-4985-826c-7e0476bd8731','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'7d0fc5a1-719b-4070-a740-fe387075f0c3'::uuid,626,true,now(),now()),
	 ('a231eab2-33bd-4c33-b6d0-58aab2032cf9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'4f16c772-1df4-4922-a9e1-761ca829bb85'::uuid,626,true,now(),now()),
	 ('47d90486-44c4-4e3c-bc98-bc74f18462d2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'4f16c772-1df4-4922-a9e1-761ca829bb85'::uuid,626,true,now(),now()),
	 ('b8564a0d-be24-4105-89a2-3a4dffc535fa','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'7675199b-55b9-4184-bce8-a6c0c2c9e9ab'::uuid,626,true,now(),now()),
	 ('ce4b3b42-471f-4290-a9a6-926b61de020d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'7675199b-55b9-4184-bce8-a6c0c2c9e9ab'::uuid,626,true,now(),now()),
	 ('659f1ec1-5944-4a5f-bd1a-b2e67aec48b1','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'4fb560d1-6bf5-46b7-a047-d381a76c4fef'::uuid,2087,true,now(),now()),
	 ('645327e8-d29f-4299-99fb-0d5bc397c0fc','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'4fb560d1-6bf5-46b7-a047-d381a76c4fef'::uuid,2087,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('a8fedb7b-beea-4209-b304-76f25ea160dd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'f42c9e51-5b7e-4ab3-847d-fd86b4e90dc1'::uuid,1796,true,now(),now()),
	 ('236db6d4-20d0-4e4d-8636-b40e3556c55c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'f42c9e51-5b7e-4ab3-847d-fd86b4e90dc1'::uuid,1891,true,now(),now()),
	 ('593ea3ff-756b-4334-a0d3-580f15a9b728','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'47cbf0b7-e249-4b7e-8306-e5a2d2b3f394'::uuid,1821,true,now(),now()),
	 ('b221b9fa-c8db-4d43-b12b-267d0b8a3e88','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'47cbf0b7-e249-4b7e-8306-e5a2d2b3f394'::uuid,1881,true,now(),now()),
	 ('38989ad8-73bc-48eb-a79d-18682fca8f4d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'cfca47bf-4639-4b7c-aed9-5ff87c9cddde'::uuid,1833,true,now(),now()),
	 ('376a9fbf-a295-4166-b6fb-27a272c8b808','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'cfca47bf-4639-4b7c-aed9-5ff87c9cddde'::uuid,1918,true,now(),now()),
	 ('7c450d52-cdc6-431b-b280-b021de6187bc','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'10644589-71f6-4baf-ba1c-dfb19d924b25'::uuid,1766,true,now(),now()),
	 ('1e6de2ce-593c-4483-a88a-2dbf4efe7f02','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'10644589-71f6-4baf-ba1c-dfb19d924b25'::uuid,1869,true,now(),now()),
	 ('fbcada9b-b2f1-48db-bf50-ff1979e75ff0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'e337daba-5509-4507-be21-ca13ecaced9b'::uuid,1808,true,now(),now()),
	 ('73009f66-93ae-4cc7-a618-f292adf29d0f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'e337daba-5509-4507-be21-ca13ecaced9b'::uuid,1878,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('41044cdd-d548-4178-9d78-8317c76a0c77','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'cfe9ab8a-a353-433e-8204-c065deeae3d9'::uuid,1820,true,now(),now()),
	 ('db68a641-c657-4dd6-97a1-9e22e15453ae','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'cfe9ab8a-a353-433e-8204-c065deeae3d9'::uuid,1898,true,now(),now()),
	 ('0a786119-ee24-4b8b-b7c3-59f4cb95dc32','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'9a4aa0e1-6b5f-4624-a21c-3acfa858d7f3'::uuid,1820,true,now(),now()),
	 ('99a152f0-0000-43d8-a7a0-4787ea3460f3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'9a4aa0e1-6b5f-4624-a21c-3acfa858d7f3'::uuid,1867,true,now(),now()),
	 ('9fc6fd2f-f20b-4456-9e66-a8e8b78bcbab','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'7582d86d-d4e7-4a88-997d-05593ccefb37'::uuid,1826,true,now(),now()),
	 ('b0203174-2b2f-4cba-ae5e-00eb203e34d3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'7582d86d-d4e7-4a88-997d-05593ccefb37'::uuid,1889,true,now(),now()),
	 ('dd43f23d-3248-43b9-ae4e-8d6c6f027a82','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'5802e021-5283-4b43-ba85-31340065d5ec'::uuid,1777,true,now(),now()),
	 ('3df9e8c6-e8c2-4b8d-81e6-20610a1c945a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'5802e021-5283-4b43-ba85-31340065d5ec'::uuid,1886,true,now(),now()),
	 ('44200b99-d538-4857-8e20-e960abe69c27','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'027f06cd-8c82-4c4a-a583-b20ccad9cc35'::uuid,1827,true,now(),now()),
	 ('2f5bd5d3-0cb4-41dc-9afb-dd06880ec184','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'027f06cd-8c82-4c4a-a583-b20ccad9cc35'::uuid,1898,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('c86672da-4cbf-4b5a-8db9-d254c2c2322b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'4a239fdb-9ad7-4bbb-8685-528f3f861992'::uuid,1819,true,now(),now()),
	 ('fc9cefb8-3b2c-4662-b590-42241a67b974','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'4a239fdb-9ad7-4bbb-8685-528f3f861992'::uuid,1910,true,now(),now()),
	 ('c6820a70-5c7d-45cb-a561-826c7ad5e43f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'ea0fa1cc-7d80-4bd9-989e-f119c33fb881'::uuid,1830,true,now(),now()),
	 ('8c5d322e-4716-4d0e-9a5f-cc753c89c566','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'ea0fa1cc-7d80-4bd9-989e-f119c33fb881'::uuid,1885,true,now(),now()),
	 ('191b7dfc-23f7-4513-8ee3-699ceb8e9c66','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'fd89694b-06ef-4472-ac9f-614c2de3317b'::uuid,1862,true,now(),now()),
	 ('100de098-fdd9-40bf-a94e-4fb1b596f6bb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'fd89694b-06ef-4472-ac9f-614c2de3317b'::uuid,1909,true,now(),now()),
	 ('f44d7663-68ab-47ff-8ce2-e992b830ba59','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'01d0be5d-aaec-483d-a841-6ab1301aa9bd'::uuid,1843,true,now(),now()),
	 ('7c715f31-5ac9-41b1-812b-bd2184c99d81','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'01d0be5d-aaec-483d-a841-6ab1301aa9bd'::uuid,1910,true,now(),now()),
	 ('7821db48-29b8-43bc-9bbc-2ad6faf6efc4','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'ddd74fb8-c0f1-41a9-9d4f-234bd295ae1a'::uuid,1873,true,now(),now()),
	 ('d6bef20b-1800-4940-b1f2-e41c0c0ca809','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'ddd74fb8-c0f1-41a9-9d4f-234bd295ae1a'::uuid,1900,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('14e4b653-1a3d-457e-b2dd-55f8eccc265c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'c68492e9-c7d9-4394-8695-15f018ce6b90'::uuid,1821,true,now(),now()),
	 ('1b02fbb2-03b0-4a6a-8c8b-3ba19c77a1a7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'c68492e9-c7d9-4394-8695-15f018ce6b90'::uuid,1891,true,now(),now()),
	 ('1fb67b8b-a577-4e71-a5d9-9b8f8874a349','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'19ddeb7f-91c1-4bd0-83ef-264eb78a3f75'::uuid,1876,true,now(),now()),
	 ('f0e6f4aa-5f6a-4a99-8b3d-d8d6d3f01d4e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'19ddeb7f-91c1-4bd0-83ef-264eb78a3f75'::uuid,1915,true,now(),now()),
	 ('e43cf22a-5719-4681-b5a3-4cce245d8d30','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'a7f17fd7-3810-4866-9b51-8179157b4a2b'::uuid,1784,true,now(),now()),
	 ('d8c91af2-e473-4df1-9bfe-b7e54564b9ab','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'a7f17fd7-3810-4866-9b51-8179157b4a2b'::uuid,1889,true,now(),now()),
	 ('2a690ad9-ef3e-4c14-a597-5022b88c1249','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'40da86e6-76e5-443b-b4ca-27ad31a2baf6'::uuid,1771,true,now(),now()),
	 ('48bcc15e-a5bf-4f34-b253-46de62aca32e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'40da86e6-76e5-443b-b4ca-27ad31a2baf6'::uuid,1861,true,now(),now()),
	 ('3543ec4f-0dc0-4adb-be1f-c764f18a03a3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'f79dd433-2808-4f20-91ef-6b5efca07350'::uuid,1759,true,now(),now()),
	 ('61181f81-3ebd-47a9-8cbd-2462a82e8384','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'f79dd433-2808-4f20-91ef-6b5efca07350'::uuid,1876,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('1ca89947-3a4c-4b57-9ce4-622f3b847a87','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'47e88f74-4e28-4027-b05e-bf9adf63e572'::uuid,1780,true,now(),now()),
	 ('ecd9fe89-9caa-4cb2-800c-23c9c35ad2b1','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'47e88f74-4e28-4027-b05e-bf9adf63e572'::uuid,1876,true,now(),now()),
	 ('e239c048-8383-497f-a7ef-dc07748e6700','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'afb334ca-9466-44ec-9be1-4c881db6d060'::uuid,1810,true,now(),now()),
	 ('961bcd26-b814-4acf-845f-30f11705be56','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'afb334ca-9466-44ec-9be1-4c881db6d060'::uuid,1910,true,now(),now()),
	 ('c87f1c04-c588-42c1-925b-d41f018871bd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'46c16bc1-df71-4c6f-835b-400c8caaf984'::uuid,1818,true,now(),now()),
	 ('837a722a-d005-4b04-acc8-15e78b04d38c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'46c16bc1-df71-4c6f-835b-400c8caaf984'::uuid,1877,true,now(),now()),
	 ('9f068ab6-dca0-49a6-ab77-b52d4749b7ac','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'816f84d1-ea01-47a0-a799-4b68508e35cc'::uuid,1807,true,now(),now()),
	 ('4b40344b-6a03-483f-8480-2110a62aa29e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'816f84d1-ea01-47a0-a799-4b68508e35cc'::uuid,1919,true,now(),now()),
	 ('36763835-05ca-4263-b0b0-e8d5d196dc8b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'b194b7a9-a759-4c12-9482-b99e43a52294'::uuid,1808,true,now(),now()),
	 ('2e15c91f-060a-4ef0-9848-13d848352762','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'b194b7a9-a759-4c12-9482-b99e43a52294'::uuid,1904,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('ad7cfa21-365a-4f61-9d1a-255eae8f41f8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'8abaed50-eac1-4f40-83db-c07d2c3a123a'::uuid,1822,true,now(),now()),
	 ('237ec393-67f1-4400-b6f3-db60f71aab2a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'8abaed50-eac1-4f40-83db-c07d2c3a123a'::uuid,1888,true,now(),now()),
	 ('3400945a-8a5f-4a89-9193-276753c79296','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'6e43ffbc-1102-45dc-8fb2-139f6b616083'::uuid,1818,true,now(),now()),
	 ('50bfc2a9-19d4-444b-ba0d-fae4ea89ac65','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'6e43ffbc-1102-45dc-8fb2-139f6b616083'::uuid,1897,true,now(),now()),
	 ('5aa3d402-8985-4d64-8871-c080eead5e57','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'1beb0053-329a-4b47-879b-1a3046d3ff87'::uuid,1810,true,now(),now()),
	 ('b7b01dd4-1752-4e99-b3ff-c705fe8e0837','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'1beb0053-329a-4b47-879b-1a3046d3ff87'::uuid,1888,true,now(),now()),
	 ('925bbef1-5c83-4d4d-9dec-434ef70ecb16','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'c18e25f9-ec34-41ca-8c1b-05558c8d6364'::uuid,1844,true,now(),now()),
	 ('fcc8e4fb-1067-4dbf-8719-1060a3640004','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'c18e25f9-ec34-41ca-8c1b-05558c8d6364'::uuid,1900,true,now(),now()),
	 ('bb1b9230-ae25-4f62-a715-262d4adffaee','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'5bf18f68-55b8-4024-adb1-c2e6592a2582'::uuid,1887,true,now(),now()),
	 ('6bf7c3c4-3452-488c-8d70-c1bf247e6443','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'5bf18f68-55b8-4024-adb1-c2e6592a2582'::uuid,1869,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('8c9aef8e-919d-4fe0-a974-05dcd80db66c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'508d9830-6a60-44d3-992f-3c48c507f9f6'::uuid,1813,true,now(),now()),
	 ('37ae8a68-dad6-4205-8f60-cc0a9d95b547','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'508d9830-6a60-44d3-992f-3c48c507f9f6'::uuid,1889,true,now(),now()),
	 ('f44403e2-a398-4b24-9e70-bd69a90cd8f8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'b80251b4-02a2-4122-add9-ab108cd011d7'::uuid,1878,true,now(),now()),
	 ('f4d701d6-345a-475c-8825-fb3cf6ad6ea3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'d45cf336-8c4b-4651-b505-bbd34831d12d'::uuid,1804,true,now(),now()),
	 ('4b1f21bb-b877-456e-a7a3-35857c348ec8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'d45cf336-8c4b-4651-b505-bbd34831d12d'::uuid,1852,true,now(),now()),
	 ('8f6a8028-56c8-487c-9301-4ae6d1524683','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'3ece4e86-d328-4206-9f81-ec62bdf55335'::uuid,1823,true,now(),now()),
	 ('a07b2398-79c2-4a7a-ba4f-ae2dbc8bc15f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'3ece4e86-d328-4206-9f81-ec62bdf55335'::uuid,1895,true,now(),now()),
	 ('7c939efc-e66d-4fba-9679-ab2e2dc485f0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'243e6e83-ff11-4a30-af30-8751e8e63bd4'::uuid,1825,true,now(),now()),
	 ('421f867f-2100-43dd-abff-fe6bd867bf92','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'243e6e83-ff11-4a30-af30-8751e8e63bd4'::uuid,1926,true,now(),now()),
	 ('900e62e4-b152-41c0-9740-80a348f751ff','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'531e3a04-e84c-45d9-86bf-c6da0820b605'::uuid,1779,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('674ebad0-a736-4004-ad16-4ee07b8fae14','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'531e3a04-e84c-45d9-86bf-c6da0820b605'::uuid,1905,true,now(),now()),
	 ('6952c051-b8e1-446e-8910-5f83c929db70','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'c7442d31-012a-40f6-ab04-600a70db8723'::uuid,1781,true,now(),now()),
	 ('a4356c38-61da-47ea-a2b1-ee9378f46fce','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'c7442d31-012a-40f6-ab04-600a70db8723'::uuid,1868,true,now(),now()),
	 ('547f668c-a489-49ed-995a-d904989f917a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'6e802149-7e46-4d7a-ab57-6c4df832085d'::uuid,626,true,now(),now()),
	 ('e9d44a4c-9495-4a73-84a6-281b21207bf5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'6e802149-7e46-4d7a-ab57-6c4df832085d'::uuid,626,true,now(),now()),
	 ('1dea7c3e-2ad7-4353-acb1-8a290a2fa9ed','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'311e5909-df08-4086-aa09-4c21a48b5e6e'::uuid,626,true,now(),now()),
	 ('4ca0ea10-4c51-4ba3-9f36-d536ea19c35a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'311e5909-df08-4086-aa09-4c21a48b5e6e'::uuid,626,true,now(),now()),
	 ('7575b362-612e-4c57-9aa5-7515f9585064','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'709dad47-121a-4edd-ad95-b3dd6fd88f08'::uuid,626,true,now(),now()),
	 ('77af6c81-fd2c-414c-b51b-90942365cbec','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'709dad47-121a-4edd-ad95-b3dd6fd88f08'::uuid,626,true,now(),now()),
	 ('242c1539-ccf7-4c8b-95e8-ad7d75e17155','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'2124fcbf-be89-4975-9cc7-263ac14ad759'::uuid,626,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('751d2ffc-7497-4a70-9824-0c4f51beb5d0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'2124fcbf-be89-4975-9cc7-263ac14ad759'::uuid,626,true,now(),now()),
	 ('27a3e21b-0443-4010-9abc-f62ba2570aa8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'c5aab403-d0e2-4e6e-b3f1-57fc52e6c2bd'::uuid,626,true,now(),now()),
	 ('f02ac392-edc5-4d4b-bbbb-7532b23d7e9c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'c5aab403-d0e2-4e6e-b3f1-57fc52e6c2bd'::uuid,626,true,now(),now()),
	 ('676e19e4-9917-4e69-95ae-9a10c5918cd2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'649f665a-7624-4824-9cd5-b992462eb97b'::uuid,626,true,now(),now()),
	 ('101d2e31-16ef-4971-89d1-d88fb17a4ca3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'649f665a-7624-4824-9cd5-b992462eb97b'::uuid,626,true,now(),now()),
	 ('4c21775c-5613-444e-86a8-eb57a2504f8e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'c68e26d0-dc81-4320-bdd7-fa286f4cc891'::uuid,626,true,now(),now()),
	 ('b27c7194-2b2c-4feb-87af-9cc1805be5fa','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'c68e26d0-dc81-4320-bdd7-fa286f4cc891'::uuid,626,true,now(),now()),
	 ('bccdbdcd-1412-40f9-a0f2-dbf41b2fc2f5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'b7329731-65df-4427-bdee-18a0ab51efb4'::uuid,626,true,now(),now()),
	 ('c3c072fb-e77d-4baa-8384-068b9a65969e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'b7329731-65df-4427-bdee-18a0ab51efb4'::uuid,626,true,now(),now()),
	 ('ba1a8792-a3d1-46d0-b3ee-1b61a111076e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'2c144ea1-9b49-4842-ad56-e5120912fd18'::uuid,1809,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('a568056a-1cf1-4c7b-875a-7542a6482c35','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'2c144ea1-9b49-4842-ad56-e5120912fd18'::uuid,1902,true,now(),now()),
	 ('c42be09a-81ab-48df-ac2d-c9dc4ddc699f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'3733db73-602a-4402-8f94-36eec2fdab15'::uuid,1964,true,now(),now()),
	 ('0bc90c93-10d6-4a4a-ba97-4cc13ad287d7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'c3c46c6b-115a-4236-b88a-76126e7f9516'::uuid,1814,true,now(),now()),
	 ('ea2b7184-b08d-4d49-b8b4-ec6ab50d2bb9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'c3c46c6b-115a-4236-b88a-76126e7f9516'::uuid,1891,true,now(),now()),
	 ('883e5cce-678e-4e14-982a-d2465f809c5d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'fd57df67-e734-4eb2-80cf-2feafe91f238'::uuid,1781,true,now(),now()),
	 ('a9e012bc-c93b-420e-a96c-f518a519a4fa','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'fd57df67-e734-4eb2-80cf-2feafe91f238'::uuid,1911,true,now(),now()),
	 ('467fd21d-13b9-4104-9480-17dc7fc29e19','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'146c58e5-c87d-4f54-a766-8da85c6b6b2c'::uuid,1835,true,now(),now()),
	 ('75a95b4e-bf6d-459f-8b15-d2b462c40779','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'146c58e5-c87d-4f54-a766-8da85c6b6b2c'::uuid,1879,true,now(),now()),
	 ('4cccd035-d967-47ab-b547-5b82f3dbdf42','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'ba215fd2-cdfc-4b98-bd78-cfa667b1b371'::uuid,1847,true,now(),now()),
	 ('b962f326-64f0-4f52-8dc2-bcf962a22015','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'ba215fd2-cdfc-4b98-bd78-cfa667b1b371'::uuid,1899,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('ea6192cb-dfa1-4172-9790-6a5fd9acfd98','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'e5d41d36-b355-4407-9ede-cd435da69873'::uuid,1849,true,now(),now()),
	 ('08cf90cd-b71c-4da0-9201-1b0577e7fae9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'e5d41d36-b355-4407-9ede-cd435da69873'::uuid,1946,true,now(),now()),
	 ('9f479d89-29a0-4936-954d-3c16578255cb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'9893a927-6084-482c-8f1c-e85959eb3547'::uuid,1821,true,now(),now()),
	 ('937538eb-a735-4ca8-a53c-db27a37f694a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'9893a927-6084-482c-8f1c-e85959eb3547'::uuid,1895,true,now(),now()),
	 ('8ad9d809-573f-405c-baa0-96c2644557be','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'ca72968c-5921-4167-b7b6-837c88ca87f2'::uuid,1826,true,now(),now()),
	 ('19bc4aae-0418-4c33-a61c-03fc27d60b50','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'ca72968c-5921-4167-b7b6-837c88ca87f2'::uuid,1893,true,now(),now()),
	 ('2ccaa63a-567b-4ef0-a636-c7f8e1c0447b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'64265049-1b4a-4a96-9cba-e01f59cafcc7'::uuid,1799,true,now(),now()),
	 ('ebc8075d-b91a-4682-8485-7751dc37948b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'64265049-1b4a-4a96-9cba-e01f59cafcc7'::uuid,1893,true,now(),now()),
	 ('401ac592-b984-4e0b-a402-0b0ebb5eed54','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'c4c73fcb-be11-4b1a-986a-a73451d402a7'::uuid,1815,true,now(),now()),
	 ('03a1a05d-05ac-401e-9e03-095fada989da','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'c4c73fcb-be11-4b1a-986a-a73451d402a7'::uuid,1874,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('af0a2c20-4dc1-4ba3-8ca0-4de565d94943','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'098488af-82c9-49c6-9daa-879eff3d3bee'::uuid,1768,true,now(),now()),
	 ('4a9132ce-f789-4e2d-a30f-6e75b3986d39','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'098488af-82c9-49c6-9daa-879eff3d3bee'::uuid,1949,true,now(),now()),
	 ('7e98c5b4-e300-435e-b758-4f4699070820','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'b3911f28-d334-4cca-8924-7da60ea5a213'::uuid,626,true,now(),now()),
	 ('062783c1-578a-4288-a0b7-7950c096f87b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'b3911f28-d334-4cca-8924-7da60ea5a213'::uuid,626,true,now(),now()),
	 ('0533bceb-fbe9-415d-9653-70ac8cf1cc99','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'f18133b7-ef83-4b2b-beff-9c3b5f99e55a'::uuid,626,true,now(),now()),
	 ('aebf537c-8e64-47dc-82d7-9e16fe494e66','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'30040c3f-667d-4dee-ba4c-24aad0891c9c'::uuid,626,true,now(),now()),
	 ('8ba18af9-34cd-4342-856c-5cb007b400fc','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'a4fa6b22-3d7f-4d56-96f1-941f9e7570aa'::uuid,626,true,now(),now()),
	 ('7c3ce74d-16e1-4937-b917-8887d204e9bb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'a4fa6b22-3d7f-4d56-96f1-941f9e7570aa'::uuid,626,true,now(),now()),
	 ('e709cc83-b4df-494e-9437-835d845b5b02','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'9b6832a8-eb82-4afa-b12f-b52a3b2cda75'::uuid,626,true,now(),now()),
	 ('98733252-2530-45b4-be00-fff1932a93b6','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'9b6832a8-eb82-4afa-b12f-b52a3b2cda75'::uuid,626,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('b7743775-4711-4605-ac1a-9c33172a71bf','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'cae0eb53-a023-434c-ac8c-d0641067d8d8'::uuid,2087,true,now(),now()),
	 ('181ef4ff-061d-4780-8b3e-cfea4365bdd1','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'cae0eb53-a023-434c-ac8c-d0641067d8d8'::uuid,2087,true,now(),now()),
	 ('840686aa-fc11-4c37-b24b-4bf11a5e2442','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'2a1b3667-e604-41a0-b741-ba19f1f56892'::uuid,626,true,now(),now()),
	 ('45d61ed9-524d-4ca2-81e4-b23568def085','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'2a1b3667-e604-41a0-b741-ba19f1f56892'::uuid,626,true,now(),now()),
	 ('54d0de38-68cc-4584-b604-19d0ece348c5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'d53d6be6-b36c-403f-b72d-d6160e9e52c1'::uuid,626,true,now(),now()),
	 ('4eed4ebf-e2f5-47b7-ae2b-3a0d175f37b3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'d53d6be6-b36c-403f-b72d-d6160e9e52c1'::uuid,626,true,now(),now()),
	 ('893039aa-8055-4b08-ae4b-7dfc2a968dce','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'1a170f85-e7f1-467c-a4dc-7d0b7898287e'::uuid,626,true,now(),now()),
	 ('699be45b-9ab8-412a-ad34-187d1062c63a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'1a170f85-e7f1-467c-a4dc-7d0b7898287e'::uuid,626,true,now(),now()),
	 ('0db0d0ef-c601-4283-97db-59fafce04b46','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'5e8d8851-bf33-4d48-9860-acc24aceea3d'::uuid,626,true,now(),now()),
	 ('35e5c53a-72c3-49fc-a93e-0b43d4221e75','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'5e8d8851-bf33-4d48-9860-acc24aceea3d'::uuid,626,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('6f4a112a-8ddc-424e-9b33-3b0605efb322','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'7ac1c0ec-0903-477c-89e0-88efe9249c98'::uuid,626,true,now(),now()),
	 ('b368379b-a90b-41af-b960-eb38457405ac','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'7ac1c0ec-0903-477c-89e0-88efe9249c98'::uuid,626,true,now(),now()),
	 ('ccff33b9-6f34-48b7-8eda-842dd218abf6','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'0cb31c3c-dfd2-4b2a-b475-d2023008eea4'::uuid,626,true,now(),now()),
	 ('80110aac-8a24-4a43-9ce5-adcd201c95d7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'0cb31c3c-dfd2-4b2a-b475-d2023008eea4'::uuid,626,true,now(),now()),
	 ('ea9dffb9-b614-4f2e-a81e-70ee11c5ed44','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'8eb44185-f9bf-465e-8469-7bc422534319'::uuid,2087,true,now(),now()),
	 ('09d6c427-61f5-410f-8194-20dbf2f7768b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'8eb44185-f9bf-465e-8469-7bc422534319'::uuid,2087,true,now(),now()),
	 ('cff61a75-f833-4250-928e-5ad3f7ee848e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'6530aaba-4906-4d63-a6d3-deea01c99bea'::uuid,626,true,now(),now()),
	 ('4bca64d4-6e7a-431c-a463-3070e61c59d8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'6530aaba-4906-4d63-a6d3-deea01c99bea'::uuid,626,true,now(),now()),
	 ('7c8785c1-4a40-4831-8142-b6a82b3635e1','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'182eb005-c185-418d-be8b-f47212c38af3'::uuid,626,true,now(),now()),
	 ('72265281-94c5-462b-9488-d0960d754b64','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'182eb005-c185-418d-be8b-f47212c38af3'::uuid,626,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('05d592fa-c5b0-41b6-9e9f-5b0c9d1e7f9c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'4a366bb4-5104-45ea-ac9e-1da8e14387c3'::uuid,2012,true,now(),now()),
	 ('46997f8d-12c7-4e42-bf95-ccfdaf80f5ba','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'4a366bb4-5104-45ea-ac9e-1da8e14387c3'::uuid,2012,true,now(),now()),
	 ('1b5994bf-af3e-4f29-89de-25ba880e8c68','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'899d79f7-8623-4442-a398-002178cf5d94'::uuid,2012,true,now(),now()),
	 ('18290dcb-b92e-48cd-9931-14851246e6fd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'899d79f7-8623-4442-a398-002178cf5d94'::uuid,2012,true,now(),now()),
	 ('fdabde94-e9b8-4fd0-b4d5-486ff9deaef4','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'dd6c2ace-2593-445b-9569-55328090de99'::uuid,2012,true,now(),now()),
	 ('c21e567f-e573-4ab5-a467-c1d128000b75','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'dd6c2ace-2593-445b-9569-55328090de99'::uuid,2012,true,now(),now()),
	 ('8d8afda5-ac0b-48ab-918e-b3b2443ac9b6','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'02cc7df6-83d0-4ff1-a5ea-8240f5434e73'::uuid,2012,true,now(),now()),
	 ('a7f147fb-88ad-429e-b364-781318a16365','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'02cc7df6-83d0-4ff1-a5ea-8240f5434e73'::uuid,2012,true,now(),now()),
	 ('36d49eea-a2b3-4c39-8a51-ce72e20ab1b6','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'93052804-f158-485d-b3a5-f04fd0d41e55'::uuid,704,true,now(),now()),
	 ('4fa250a4-dd78-443b-9a12-01227d46dbf6','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'93052804-f158-485d-b3a5-f04fd0d41e55'::uuid,704,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('20b5ae39-b6bf-4bc8-af3e-233b9b581f01','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'6455326e-cc11-4cfe-903b-ccce70e6f04e'::uuid,704,true,now(),now()),
	 ('3fa63ad2-ae44-44f4-8bac-8ef70606e28c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'6455326e-cc11-4cfe-903b-ccce70e6f04e'::uuid,704,true,now(),now()),
	 ('52e9595c-02c6-488f-9717-e257a2ad62f5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'0ba534f5-0d24-4d7c-9216-d07f57cd8edd'::uuid,704,true,now(),now()),
	 ('ca69b727-4174-49e6-b52d-62f40eeee8ad','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'0ba534f5-0d24-4d7c-9216-d07f57cd8edd'::uuid,704,true,now(),now()),
	 ('2ca69b23-5995-4af2-b18d-b8a46bd81a71','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'433334c3-59dd-404d-a193-10dd4172fc8f'::uuid,2012,true,now(),now()),
	 ('a6ee9f60-f39f-4957-a6c2-25699c209669','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'433334c3-59dd-404d-a193-10dd4172fc8f'::uuid,2012,true,now(),now()),
	 ('2a48615e-b324-4870-baad-5b7219029d07','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'e3071ca8-bedf-4eff-bda0-e9ff27f0e34c'::uuid,2012,true,now(),now()),
	 ('429c725d-bc8b-4e51-a629-42c50283d876','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'e3071ca8-bedf-4eff-bda0-e9ff27f0e34c'::uuid,2012,true,now(),now()),
	 ('ed24df00-f624-4f9a-99b2-85d9a7e8d0bc','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'3320e408-93d8-4933-abb8-538a5d697b41'::uuid,704,true,now(),now()),
	 ('f1a0cf6d-ce3f-4ab2-a574-9649e2f610f8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'3320e408-93d8-4933-abb8-538a5d697b41'::uuid,704,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('ba1d21f1-ce72-44f2-8ad3-496893b30aaa','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'9a9da923-06ef-47ea-bc20-23cc85b51ad0'::uuid,704,true,now(),now()),
	 ('1358f2af-44c1-4a9d-be4c-6fa9e72fbf90','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'1e23a20c-2558-47bf-b720-d7758b717ce3'::uuid,704,true,now(),now()),
	 ('05b23ae8-566a-42c4-ab85-a03c2d7ad61a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'1e23a20c-2558-47bf-b720-d7758b717ce3'::uuid,704,true,now(),now()),
	 ('ac5a529b-04f4-47f6-be7e-5d416f7c247d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'58dcc836-51e1-4633-9a89-73ac44eb2152'::uuid,704,true,now(),now()),
	 ('90f23adc-86e6-4e05-a134-cbfeea1c5240','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'7ee486f1-4de8-4700-922b-863168f612a0'::uuid,704,true,now(),now()),
	 ('7445e033-e801-4d6b-8ed4-c03da0f66d7a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'7ee486f1-4de8-4700-922b-863168f612a0'::uuid,704,true,now(),now()),
	 ('6c976876-cb58-4ad2-97f9-6deb41b1139e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'3ec11db4-f821-409f-84ad-07fc8e64d60d'::uuid,704,true,now(),now()),
	 ('58657686-74b8-4823-9d4f-b03fd7c66f85','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'3ec11db4-f821-409f-84ad-07fc8e64d60d'::uuid,704,true,now(),now()),
	 ('613efb6f-b6df-45a5-8550-302b01ef5e84','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'dcc3cae7-e05e-4ade-9b5b-c2eaade9f101'::uuid,626,true,now(),now()),
	 ('315f0903-99e3-4a0f-b9af-b8b9e23bb499','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'dcc3cae7-e05e-4ade-9b5b-c2eaade9f101'::uuid,626,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('b2bc4cad-0b6c-49f5-b2ad-64725c5d91b2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'6f0e02be-08ad-48b1-8e23-eecaab34b4fe'::uuid,626,true,now(),now()),
	 ('bc52d19c-4ef1-4c52-bb89-e789dd939046','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'6f0e02be-08ad-48b1-8e23-eecaab34b4fe'::uuid,626,true,now(),now()),
	 ('f16ac3e5-f803-490c-af96-4eed0f129290','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'4f2e3e38-6bf4-4e74-bd7b-fe6edb87ee42'::uuid,626,true,now(),now()),
	 ('2beba562-89d6-4e18-a8fa-57e97692f776','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'4f2e3e38-6bf4-4e74-bd7b-fe6edb87ee42'::uuid,626,true,now(),now()),
	 ('c75ccb1b-b4cc-4b96-82fb-7b3825dd9a06','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'760f146d-d5e7-4e08-9464-45371ea3267d'::uuid,626,true,now(),now()),
	 ('c780e148-e057-477f-a335-e617a8e6d554','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'760f146d-d5e7-4e08-9464-45371ea3267d'::uuid,626,true,now(),now()),
	 ('bce8a2d6-c72d-4440-9fa2-a1d9b2b32ccc','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'a2fad63c-b6cb-4b0d-9ced-1a81a6bc9985'::uuid,626,true,now(),now()),
	 ('ab393354-95bf-4bf7-8db5-4e3ba26e5d6b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'0026678a-51b7-46de-af3d-b49428e0916c'::uuid,626,true,now(),now()),
	 ('99b3e68e-7ad4-4ebc-b3fd-073757116cb5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'0026678a-51b7-46de-af3d-b49428e0916c'::uuid,626,true,now(),now()),
	 ('f0b3517f-3534-450d-871a-f19265a2ff4c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'40ab17b2-9e79-429c-a75d-b6fcbbe27901'::uuid,626,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('a8bf3539-68ad-4ff0-b406-5f4b02968f2f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'43a09249-d81b-4897-b5c7-dd88331cf2bd'::uuid,626,true,now(),now()),
	 ('3ef0aa10-9d25-4bb6-be53-06bb1c541095','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'43a09249-d81b-4897-b5c7-dd88331cf2bd'::uuid,626,true,now(),now()),
	 ('a6af47f0-8e5e-49ac-907d-7d3ec9a21d66','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'a2fad63c-b6cb-4b0d-9ced-1a81a6bc9985'::uuid,626,true,now(),now()),
	 ('6b637056-fea4-4732-9e2f-c0deda00c8a7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'422021c7-08e1-4355-838d-8f2821f00f42'::uuid,626,true,now(),now()),
	 ('b339996a-3c47-4d5c-a15f-09e222c9932e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'422021c7-08e1-4355-838d-8f2821f00f42'::uuid,626,true,now(),now()),
	 ('1ab59d56-9a5a-4fb0-bd9d-c05214468455','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'e4e467f2-449d-46e3-a59b-0f8714e4824a'::uuid,626,true,now(),now()),
	 ('58a35e85-33ce-4a11-ac42-55228e391de1','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'e4e467f2-449d-46e3-a59b-0f8714e4824a'::uuid,626,true,now(),now()),
	 ('547b3980-d179-41b8-941f-2f113c839867','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'535e6789-c126-405f-8b3a-7bd886b94796'::uuid,626,true,now(),now()),
	 ('e963bb44-48f0-4183-959f-e7e24ae1a3a7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'40ab17b2-9e79-429c-a75d-b6fcbbe27901'::uuid,626,true,now(),now()),
	 ('1befc86e-de46-4b15-a4fc-e12be04aec55','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'30040c3f-667d-4dee-ba4c-24aad0891c9c'::uuid,626,true,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('ecce7b25-8274-4ab4-8d0d-0d2f0cc207ee','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'def8c7af-d4fc-474e-974d-6fd00c251da8'::uuid,1716,true,now(),now()),
	 ('319e5861-0b8d-4429-920e-69cbfd9683d4','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'def8c7af-d4fc-474e-974d-6fd00c251da8'::uuid,1875,true,now(),now()),
	 ('bea2cc95-be2b-4fe6-aa58-db118056673c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'3733db73-602a-4402-8f94-36eec2fdab15'::uuid,1733,true,now(),now()),
	 ('cf8cecc2-bc71-4bbb-af30-8ac6c9af036d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'b80251b4-02a2-4122-add9-ab108cd011d7'::uuid,1795,true,now(),now()),
	 ('f83fc4e4-ec7b-4786-9773-afdd85d54f8a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'811a32c0-90d6-4744-9a57-ab4130091754'::uuid,1820,true,now(),now()),
	 ('76749413-80bd-4e08-b5df-ab4e2290e11d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'811a32c0-90d6-4744-9a57-ab4130091754'::uuid,1916,true,now(),now()),
	 ('c15d6b3f-4269-410a-b84f-a2bc5d7ca2ff','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'9a9da923-06ef-47ea-bc20-23cc85b51ad0'::uuid,704,true,now(),now()),
	 ('796f250f-a466-4dce-9dfa-e19101bd3e94','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',true,'535e6789-c126-405f-8b3a-7bd886b94796'::uuid,626,true,now(),now()),
	 ('209c8784-15a4-433e-81e9-f96021007a7a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'f18133b7-ef83-4b2b-beff-9c3b5f99e55a'::uuid,626,true,now(),now()),
	 ('b215b0ec-f126-43eb-8994-bab53bf19998','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'28389ee1-56cf-400c-aa52-1501ecdd7c69',false,'58dcc836-51e1-4633-9a89-73ac44eb2152'::uuid,704,true,now(),now());

--insert IOPSIT recs for > 50 miles
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('52efa9a7-4ada-4037-aa18-0ce7cf1945cf','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'a761a482-2929-4345-8027-3c6258f0c8dd'::uuid,54,false,now(),now()),
	 ('6eb777bd-23a2-4466-88f7-5eda2536fdd1','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'a761a482-2929-4345-8027-3c6258f0c8dd'::uuid,54,false,now(),now()),
	 ('308cca12-5095-47a0-95d8-bfeb0853c041','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'91eb2878-0368-4347-97e3-e6caa362d878'::uuid,54,false,now(),now()),
	 ('a86bbaa0-caf6-48e0-a613-76b3cd3f7cc7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'91eb2878-0368-4347-97e3-e6caa362d878'::uuid,54,false,now(),now()),
	 ('71bb51d2-fdb9-4c1d-853f-3af1619d1941','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'635e4b79-342c-4cfc-8069-39c408a2decd'::uuid,16,false,now(),now()),
	 ('39f906c5-334b-423f-9731-1b90b277bafe','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'635e4b79-342c-4cfc-8069-39c408a2decd'::uuid,16,false,now(),now()),
	 ('b7990076-2f7d-44c7-b969-6c0203679917','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'9bb87311-1b29-4f29-8561-8a4c795654d4'::uuid,16,false,now(),now()),
	 ('870eb56a-031f-4316-8766-4b64e5f478ed','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'9bb87311-1b29-4f29-8561-8a4c795654d4'::uuid,16,false,now(),now()),
	 ('e3a021cc-ca18-48de-a8a7-897f052b1fdb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'b80a00d4-f829-4051-961a-b8945c62c37d'::uuid,16,false,now(),now()),
	 ('cdfae55e-b566-4733-bb8b-d7c19787c51e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'b80a00d4-f829-4051-961a-b8945c62c37d'::uuid,16,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('ec3f6ebc-26a0-4126-94f4-0e311393a29d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'5a27e806-21d4-4672-aa5e-29518f10c0aa'::uuid,16,false,now(),now()),
	 ('b293c90e-ae6d-48d0-94cb-1ea6c7f4b94e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'5a27e806-21d4-4672-aa5e-29518f10c0aa'::uuid,16,false,now(),now()),
	 ('1c19f52c-1bf2-4004-8037-967e2e1b8149','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'71755cc7-0844-4523-a0ac-da9a1e743ad1'::uuid,42,false,now(),now()),
	 ('a06ea893-7e1b-41d3-8bb5-66ee2b51ff7b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'71755cc7-0844-4523-a0ac-da9a1e743ad1'::uuid,52,false,now(),now()),
	 ('2b78fe96-2f5e-48bd-b428-c7ccceb5e2e5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'829d8b45-19c1-49a3-920c-cc0ae14e8698'::uuid,21,false,now(),now()),
	 ('9e3bf610-e5c0-46bc-9274-374c7166e2dc','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'829d8b45-19c1-49a3-920c-cc0ae14e8698'::uuid,21,false,now(),now()),
	 ('f4fcf5c0-73d4-418d-b321-d4381939b96f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'03dd5854-8bc3-4b56-986e-eac513cc1ec0'::uuid,21,false,now(),now()),
	 ('0486d516-72b7-42e7-9101-a5fdd94fe277','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'03dd5854-8bc3-4b56-986e-eac513cc1ec0'::uuid,21,false,now(),now()),
	 ('4f9de586-9c50-4a2d-8175-8169ff02bfbe','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'ee0ffe93-32b3-4817-982e-6d081da85d28'::uuid,21,false,now(),now()),
	 ('f878bf82-ae75-4a6a-ae1e-a9932e3fc73b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'ee0ffe93-32b3-4817-982e-6d081da85d28'::uuid,52,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('c0d871b0-5744-4315-a037-ea677742551f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'0506bf0f-bc1c-43c7-a75f-639a1b4c0449'::uuid,21,false,now(),now()),
	 ('c70a8c5f-bcd9-4d09-ac9d-87ef585a5f68','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'0506bf0f-bc1c-43c7-a75f-639a1b4c0449'::uuid,21,false,now(),now()),
	 ('8b1539c2-bacc-44c4-ac57-f971aa880f30','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'6a0f9a02-b6ba-4585-9d7a-6959f7b0248f'::uuid,21,false,now(),now()),
	 ('beb2b518-56c4-43de-82d6-29d5d3deccbb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'6a0f9a02-b6ba-4585-9d7a-6959f7b0248f'::uuid,21,false,now(),now()),
	 ('645b6393-46a8-4d97-bbfb-8b4663031444','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'c9036eb8-84bb-4909-be20-0662387219a7'::uuid,21,false,now(),now()),
	 ('c4911bf6-12db-4163-a0e7-faab45aec016','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'c9036eb8-84bb-4909-be20-0662387219a7'::uuid,21,false,now(),now()),
	 ('b165c788-6adb-4818-baaa-6f9ff6615e10','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'fe76b78f-67bc-4125-8f81-8e68697c136d'::uuid,21,false,now(),now()),
	 ('85da2f00-cd81-4c8d-874a-430eb77f76ec','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'fe76b78f-67bc-4125-8f81-8e68697c136d'::uuid,21,false,now(),now()),
	 ('76324579-4776-4a05-98d8-db8dfa6b2013','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'2b1d1842-15f8-491a-bdce-e5f9fea947e7'::uuid,21,false,now(),now()),
	 ('a110dc84-1256-48d9-89d7-e6c524f1678d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'2b1d1842-15f8-491a-bdce-e5f9fea947e7'::uuid,21,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('c19b80ae-0dd5-4050-90be-581df84e932b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'612c2ce9-39cc-45e6-a3f1-c6672267d392'::uuid,21,false,now(),now()),
	 ('01dfe0e0-9728-489d-b5d6-0cb7a3c944dd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'612c2ce9-39cc-45e6-a3f1-c6672267d392'::uuid,21,false,now(),now()),
	 ('9aeb25c9-a150-438d-b19d-48c3ae2d5a20','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'7d0fc5a1-719b-4070-a740-fe387075f0c3'::uuid,21,false,now(),now()),
	 ('290ca426-9ae8-4941-aba2-7011a9364c80','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'7d0fc5a1-719b-4070-a740-fe387075f0c3'::uuid,21,false,now(),now()),
	 ('a5ddc793-0908-43e2-9d1a-0a038841db43','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'4f16c772-1df4-4922-a9e1-761ca829bb85'::uuid,21,false,now(),now()),
	 ('092fcdc7-4526-4e71-816e-47ade6a6542c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'4f16c772-1df4-4922-a9e1-761ca829bb85'::uuid,21,false,now(),now()),
	 ('a9abb681-4bf4-419f-bd91-216c6df4726a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'7675199b-55b9-4184-bce8-a6c0c2c9e9ab'::uuid,21,false,now(),now()),
	 ('74ba1dba-8937-402f-9939-310e4fd2f869','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'7675199b-55b9-4184-bce8-a6c0c2c9e9ab'::uuid,21,false,now(),now()),
	 ('8a8fc5df-f670-43fa-988f-6b3fe410cfdd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'4fb560d1-6bf5-46b7-a047-d381a76c4fef'::uuid,42,false,now(),now()),
	 ('35860792-688d-4e0a-bb93-f167581a909f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'4fb560d1-6bf5-46b7-a047-d381a76c4fef'::uuid,42,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('d91821a2-264e-4247-b0ff-be24d221a696','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'f42c9e51-5b7e-4ab3-847d-fd86b4e90dc1'::uuid,54,false,now(),now()),
	 ('8f9c639c-deb6-4015-879b-b5d0fb4b717e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'f42c9e51-5b7e-4ab3-847d-fd86b4e90dc1'::uuid,54,false,now(),now()),
	 ('70a22d67-30d6-4e1f-a21a-75a32e3b694f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'47cbf0b7-e249-4b7e-8306-e5a2d2b3f394'::uuid,54,false,now(),now()),
	 ('579d66ab-5493-4bc1-973a-83b04f60e76a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'47cbf0b7-e249-4b7e-8306-e5a2d2b3f394'::uuid,54,false,now(),now()),
	 ('f74de640-a7d1-429b-9af0-eaec814e566b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'cfca47bf-4639-4b7c-aed9-5ff87c9cddde'::uuid,54,false,now(),now()),
	 ('873a3664-d688-4172-a1e5-5f64bd9dade0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'cfca47bf-4639-4b7c-aed9-5ff87c9cddde'::uuid,54,false,now(),now()),
	 ('4ea26f5c-4ca6-4f0e-96d1-aaf856bb2bd2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'10644589-71f6-4baf-ba1c-dfb19d924b25'::uuid,54,false,now(),now()),
	 ('ae6cfad5-4e92-4e45-a758-40289c6fc3f9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'10644589-71f6-4baf-ba1c-dfb19d924b25'::uuid,54,false,now(),now()),
	 ('a1154f74-05d8-4ca0-a6f1-945407fc16b9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'e337daba-5509-4507-be21-ca13ecaced9b'::uuid,54,false,now(),now()),
	 ('d118e381-0e73-4b60-bcb2-e067746674cd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'e337daba-5509-4507-be21-ca13ecaced9b'::uuid,54,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('b2d32030-ec39-4e59-830e-297923d6f226','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'cfe9ab8a-a353-433e-8204-c065deeae3d9'::uuid,54,false,now(),now()),
	 ('39b5b2be-5a62-4876-9bda-e32d0b2a1c5e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'cfe9ab8a-a353-433e-8204-c065deeae3d9'::uuid,54,false,now(),now()),
	 ('e3bd9abd-2542-469a-8040-93bcdb461629','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'9a4aa0e1-6b5f-4624-a21c-3acfa858d7f3'::uuid,54,false,now(),now()),
	 ('aa5b1622-ecf7-4f73-8a45-2fac8a1935f2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'9a4aa0e1-6b5f-4624-a21c-3acfa858d7f3'::uuid,54,false,now(),now()),
	 ('813d8c75-f9e9-4a0a-b296-112bceee15d7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'7582d86d-d4e7-4a88-997d-05593ccefb37'::uuid,54,false,now(),now()),
	 ('b18ddeec-c129-4959-8756-f3dfb6093420','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'7582d86d-d4e7-4a88-997d-05593ccefb37'::uuid,54,false,now(),now()),
	 ('f1d0353d-1042-49fb-bcf5-50733229c2dd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'5802e021-5283-4b43-ba85-31340065d5ec'::uuid,54,false,now(),now()),
	 ('f082ba1c-1bf7-4367-8f2b-fd7aea6f5dc2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'5802e021-5283-4b43-ba85-31340065d5ec'::uuid,54,false,now(),now()),
	 ('7261d6d7-1553-4ad8-9ac1-9dfcac1a448a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'027f06cd-8c82-4c4a-a583-b20ccad9cc35'::uuid,54,false,now(),now()),
	 ('39a19734-96d4-4f33-af0c-9ae107a977bb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'027f06cd-8c82-4c4a-a583-b20ccad9cc35'::uuid,54,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('6dfbddac-2bb4-437a-ab76-eec0212b81a4','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'4a239fdb-9ad7-4bbb-8685-528f3f861992'::uuid,54,false,now(),now()),
	 ('323d2d23-7974-4b77-8e22-59cffebbbc0b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'4a239fdb-9ad7-4bbb-8685-528f3f861992'::uuid,54,false,now(),now()),
	 ('fbcf1982-28ee-45e0-93a2-f2c61c67f466','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'ea0fa1cc-7d80-4bd9-989e-f119c33fb881'::uuid,54,false,now(),now()),
	 ('e73a08c9-b3d4-4a78-bc26-36b38014752c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'ea0fa1cc-7d80-4bd9-989e-f119c33fb881'::uuid,54,false,now(),now()),
	 ('6e021cb4-9a7c-41dd-9c79-727590beb6f3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'fd89694b-06ef-4472-ac9f-614c2de3317b'::uuid,54,false,now(),now()),
	 ('0f6dfe1c-c58b-4e57-af2e-7881559d8640','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'fd89694b-06ef-4472-ac9f-614c2de3317b'::uuid,54,false,now(),now()),
	 ('eb79a687-b876-47a6-850d-01e999b9a50a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'01d0be5d-aaec-483d-a841-6ab1301aa9bd'::uuid,54,false,now(),now()),
	 ('99d0c80f-e9b7-43bd-9c45-800441ea43d7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'01d0be5d-aaec-483d-a841-6ab1301aa9bd'::uuid,54,false,now(),now()),
	 ('dcd74fb3-1ef8-420d-aa0f-c1a9f9c9fab3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'ddd74fb8-c0f1-41a9-9d4f-234bd295ae1a'::uuid,54,false,now(),now()),
	 ('6dc2e847-f72e-48e9-b2d1-efec573a53cd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'ddd74fb8-c0f1-41a9-9d4f-234bd295ae1a'::uuid,54,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('eb436276-eb00-4b82-93d5-bc4b0db8b243','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'c68492e9-c7d9-4394-8695-15f018ce6b90'::uuid,54,false,now(),now()),
	 ('776af8d0-48c2-4640-ae5f-fe2e479972d8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'c68492e9-c7d9-4394-8695-15f018ce6b90'::uuid,54,false,now(),now()),
	 ('486e52fc-102a-441f-9d1e-eab3025e19b7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'19ddeb7f-91c1-4bd0-83ef-264eb78a3f75'::uuid,54,false,now(),now()),
	 ('8ed20096-f27a-4c74-adf5-cf94717e817e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'19ddeb7f-91c1-4bd0-83ef-264eb78a3f75'::uuid,54,false,now(),now()),
	 ('8737a9d7-94f4-4fe1-a477-7772921139b2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'a7f17fd7-3810-4866-9b51-8179157b4a2b'::uuid,54,false,now(),now()),
	 ('e225efa5-df98-46a2-8f15-bc0c259958ce','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'a7f17fd7-3810-4866-9b51-8179157b4a2b'::uuid,54,false,now(),now()),
	 ('c5a02139-3ace-447d-91dc-89b2a6daf572','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'40da86e6-76e5-443b-b4ca-27ad31a2baf6'::uuid,54,false,now(),now()),
	 ('e87992f8-4bb8-4645-96d3-f9efcb518339','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'40da86e6-76e5-443b-b4ca-27ad31a2baf6'::uuid,54,false,now(),now()),
	 ('d5ea9953-efc3-4aac-8075-966a4f1888dc','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'f79dd433-2808-4f20-91ef-6b5efca07350'::uuid,54,false,now(),now()),
	 ('73d55cfa-e7b2-411c-9e36-6ee4f51dfbcb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'f79dd433-2808-4f20-91ef-6b5efca07350'::uuid,54,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('4be167aa-ba0b-476b-a586-f9ae030619e0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'47e88f74-4e28-4027-b05e-bf9adf63e572'::uuid,54,false,now(),now()),
	 ('a48a695b-30f2-4e72-9ff6-3d199072c21b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'47e88f74-4e28-4027-b05e-bf9adf63e572'::uuid,54,false,now(),now()),
	 ('94c720eb-27b4-436d-8168-68320b17a539','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'afb334ca-9466-44ec-9be1-4c881db6d060'::uuid,54,false,now(),now()),
	 ('1d3ffb6c-6642-4c37-9294-f6f18472b4fd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'afb334ca-9466-44ec-9be1-4c881db6d060'::uuid,54,false,now(),now()),
	 ('3294e933-357e-45b8-809c-72cbdf4b7cb6','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'46c16bc1-df71-4c6f-835b-400c8caaf984'::uuid,54,false,now(),now()),
	 ('37fd6596-06d9-46a5-8cd6-e8505f7200f0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'46c16bc1-df71-4c6f-835b-400c8caaf984'::uuid,54,false,now(),now()),
	 ('5d987f31-8132-4bf7-a21e-095415dcdd86','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'816f84d1-ea01-47a0-a799-4b68508e35cc'::uuid,54,false,now(),now()),
	 ('f201bba7-cb68-4e17-94bf-94f5d38d7c42','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'816f84d1-ea01-47a0-a799-4b68508e35cc'::uuid,54,false,now(),now()),
	 ('180188d1-1df8-47c0-98c1-f29433dcdb6f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'b194b7a9-a759-4c12-9482-b99e43a52294'::uuid,54,false,now(),now()),
	 ('1a3a2e6d-5b37-406a-a48b-2632ef6e208d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'b194b7a9-a759-4c12-9482-b99e43a52294'::uuid,54,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('d21f6f26-833f-472e-be0e-1ce9315608dc','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'8abaed50-eac1-4f40-83db-c07d2c3a123a'::uuid,54,false,now(),now()),
	 ('d56942de-2384-4f3f-93f7-85cc709c2a87','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'8abaed50-eac1-4f40-83db-c07d2c3a123a'::uuid,54,false,now(),now()),
	 ('86409760-be84-4a64-8554-2c9b50963988','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'6e43ffbc-1102-45dc-8fb2-139f6b616083'::uuid,54,false,now(),now()),
	 ('fc8ab61c-212c-4530-be5c-8fb38cc6a973','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'6e43ffbc-1102-45dc-8fb2-139f6b616083'::uuid,54,false,now(),now()),
	 ('ab6e3c45-9482-4b2a-8bf9-a5cb59d8e515','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'1beb0053-329a-4b47-879b-1a3046d3ff87'::uuid,54,false,now(),now()),
	 ('64c64eb7-d911-44a3-85ea-d94fff4c46ef','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'1beb0053-329a-4b47-879b-1a3046d3ff87'::uuid,54,false,now(),now()),
	 ('1399418b-0657-4351-8f44-e60cfa6f67e7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'c18e25f9-ec34-41ca-8c1b-05558c8d6364'::uuid,54,false,now(),now()),
	 ('c5a9bcd7-1f44-4a25-8163-39e5a6c5ac9f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'c18e25f9-ec34-41ca-8c1b-05558c8d6364'::uuid,54,false,now(),now()),
	 ('11effbae-5ff9-4432-a25f-d9f1a8649cd1','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'5bf18f68-55b8-4024-adb1-c2e6592a2582'::uuid,54,false,now(),now()),
	 ('d1e07618-ff01-4c65-9885-2e79a0059053','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'5bf18f68-55b8-4024-adb1-c2e6592a2582'::uuid,54,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('e41736dc-d4e0-402e-86ac-4c026cec514b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'508d9830-6a60-44d3-992f-3c48c507f9f6'::uuid,54,false,now(),now()),
	 ('63e128f9-4828-482a-b717-f45063d09824','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'508d9830-6a60-44d3-992f-3c48c507f9f6'::uuid,54,false,now(),now()),
	 ('ad9b65c8-cef2-4ff6-8b7f-5e273ef8bb1f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'b80251b4-02a2-4122-add9-ab108cd011d7'::uuid,54,false,now(),now()),
	 ('277699dc-26bc-4f1b-a815-5aaa6ea12d58','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'d45cf336-8c4b-4651-b505-bbd34831d12d'::uuid,54,false,now(),now()),
	 ('f5afc3af-5c71-487f-8e06-bff20cb89052','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'d45cf336-8c4b-4651-b505-bbd34831d12d'::uuid,54,false,now(),now()),
	 ('14684d1d-685c-43ee-b962-69a0567bc4a3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'3ece4e86-d328-4206-9f81-ec62bdf55335'::uuid,54,false,now(),now()),
	 ('7858993b-2558-4d3e-b799-0bca5b45de32','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'3ece4e86-d328-4206-9f81-ec62bdf55335'::uuid,54,false,now(),now()),
	 ('29d96b4b-2bb0-420f-adf2-b462d7d282ef','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'243e6e83-ff11-4a30-af30-8751e8e63bd4'::uuid,54,false,now(),now()),
	 ('5acf7efb-e0c4-4819-87ee-4539fd73d611','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'243e6e83-ff11-4a30-af30-8751e8e63bd4'::uuid,54,false,now(),now()),
	 ('9220fb1c-6aa2-4399-b753-29951e3d3b35','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'531e3a04-e84c-45d9-86bf-c6da0820b605'::uuid,54,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('7f1576d6-9dcb-448b-ae66-793332254654','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'531e3a04-e84c-45d9-86bf-c6da0820b605'::uuid,54,false,now(),now()),
	 ('f0aa8bdb-d1eb-4b05-8296-2f07a47ae291','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'c7442d31-012a-40f6-ab04-600a70db8723'::uuid,54,false,now(),now()),
	 ('3735a631-0a35-41ef-8c35-e89da0406a72','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'c7442d31-012a-40f6-ab04-600a70db8723'::uuid,54,false,now(),now()),
	 ('281ddf44-89e6-496f-907b-2a3c880b6655','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'6e802149-7e46-4d7a-ab57-6c4df832085d'::uuid,21,false,now(),now()),
	 ('847d90fe-7755-4e43-8fe2-b6ccabe9c7cf','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'6e802149-7e46-4d7a-ab57-6c4df832085d'::uuid,21,false,now(),now()),
	 ('de8a1cc8-861e-4dd0-9599-1fc9eae275ff','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'311e5909-df08-4086-aa09-4c21a48b5e6e'::uuid,21,false,now(),now()),
	 ('efe42269-dd78-4565-900e-637a35f74b6c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'311e5909-df08-4086-aa09-4c21a48b5e6e'::uuid,21,false,now(),now()),
	 ('102be076-949c-4f8d-99e7-ffc52ef25553','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'709dad47-121a-4edd-ad95-b3dd6fd88f08'::uuid,21,false,now(),now()),
	 ('12b4f819-bb8d-4248-9437-96d412ceba45','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'709dad47-121a-4edd-ad95-b3dd6fd88f08'::uuid,21,false,now(),now()),
	 ('c1e9cf3d-ba50-4ee6-b10c-2d6b89d02785','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'2124fcbf-be89-4975-9cc7-263ac14ad759'::uuid,21,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('08706ef6-f126-4f51-85a4-ce0356c2f4bb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'2124fcbf-be89-4975-9cc7-263ac14ad759'::uuid,21,false,now(),now()),
	 ('20944767-8c16-4f37-b238-ae83de0e2dd2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'c5aab403-d0e2-4e6e-b3f1-57fc52e6c2bd'::uuid,21,false,now(),now()),
	 ('53f7ff35-2859-494c-ab31-ae34cca4b9db','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'c5aab403-d0e2-4e6e-b3f1-57fc52e6c2bd'::uuid,21,false,now(),now()),
	 ('540cb843-f255-4cb6-981f-7b8a0ecd9265','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'649f665a-7624-4824-9cd5-b992462eb97b'::uuid,21,false,now(),now()),
	 ('33020469-ad7b-4b40-a5d4-8794cac952fa','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'649f665a-7624-4824-9cd5-b992462eb97b'::uuid,21,false,now(),now()),
	 ('fdb3962e-f90a-4810-9f6a-0d343cbb6989','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'c68e26d0-dc81-4320-bdd7-fa286f4cc891'::uuid,21,false,now(),now()),
	 ('9648a495-76df-4eff-aef9-d256c237c7cd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'c68e26d0-dc81-4320-bdd7-fa286f4cc891'::uuid,21,false,now(),now()),
	 ('ab28ed7c-9809-4ed7-8ccd-f051dc9b2349','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'b7329731-65df-4427-bdee-18a0ab51efb4'::uuid,21,false,now(),now()),
	 ('4d37a380-c2be-4c39-8c79-478f905fb105','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'b7329731-65df-4427-bdee-18a0ab51efb4'::uuid,21,false,now(),now()),
	 ('5b70448b-0d4c-4efa-b5ec-5cfb2dbde5fd','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'2c144ea1-9b49-4842-ad56-e5120912fd18'::uuid,54,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('ab6e2d63-8b1d-4565-a974-c8781fea90a6','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'2c144ea1-9b49-4842-ad56-e5120912fd18'::uuid,54,false,now(),now()),
	 ('818d731f-a304-4ddf-85fe-18150b46121c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'3733db73-602a-4402-8f94-36eec2fdab15'::uuid,54,false,now(),now()),
	 ('dcaa1cca-87d7-42d2-965c-e1fa6d8a7396','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'c3c46c6b-115a-4236-b88a-76126e7f9516'::uuid,54,false,now(),now()),
	 ('582049ca-0556-4cde-9658-3703bb052dc4','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'c3c46c6b-115a-4236-b88a-76126e7f9516'::uuid,54,false,now(),now()),
	 ('fe610a8d-c366-482d-931e-fde66b77ba10','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'fd57df67-e734-4eb2-80cf-2feafe91f238'::uuid,54,false,now(),now()),
	 ('fa74ef0a-874b-4f69-8af8-d21745b770a9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'fd57df67-e734-4eb2-80cf-2feafe91f238'::uuid,54,false,now(),now()),
	 ('fa2e8029-781d-4adb-adad-090cb02b8f25','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'146c58e5-c87d-4f54-a766-8da85c6b6b2c'::uuid,54,false,now(),now()),
	 ('517d392e-4299-4657-9683-548d4864d801','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'146c58e5-c87d-4f54-a766-8da85c6b6b2c'::uuid,54,false,now(),now()),
	 ('9706a4b5-9b9b-43f4-a193-38c311451973','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'ba215fd2-cdfc-4b98-bd78-cfa667b1b371'::uuid,54,false,now(),now()),
	 ('1a2bdd92-d2ac-4914-bb8a-8bca78d2a3a9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'ba215fd2-cdfc-4b98-bd78-cfa667b1b371'::uuid,54,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('1895384b-9fbf-4bdf-8eee-51df10b33c64','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'e5d41d36-b355-4407-9ede-cd435da69873'::uuid,54,false,now(),now()),
	 ('1d1860ef-4c87-4482-9cf8-66c647595685','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'e5d41d36-b355-4407-9ede-cd435da69873'::uuid,54,false,now(),now()),
	 ('467790e5-671d-4b25-b6cc-5ae18850e152','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'9893a927-6084-482c-8f1c-e85959eb3547'::uuid,54,false,now(),now()),
	 ('bd8ea688-f79b-458c-920d-8de68821ab97','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'9893a927-6084-482c-8f1c-e85959eb3547'::uuid,54,false,now(),now()),
	 ('5db42b4a-ea03-4b9a-b616-33a477bafa67','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'ca72968c-5921-4167-b7b6-837c88ca87f2'::uuid,54,false,now(),now()),
	 ('41ae2e0d-8ebf-4008-b784-82fdd00ef6e7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'ca72968c-5921-4167-b7b6-837c88ca87f2'::uuid,54,false,now(),now()),
	 ('1800df09-5374-47f8-847d-a21c65973dcb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'64265049-1b4a-4a96-9cba-e01f59cafcc7'::uuid,54,false,now(),now()),
	 ('b0939b03-1ae6-432c-8622-d236e5b49161','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'64265049-1b4a-4a96-9cba-e01f59cafcc7'::uuid,54,false,now(),now()),
	 ('833a70a9-802a-434e-93b3-32aaece3ef2a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'c4c73fcb-be11-4b1a-986a-a73451d402a7'::uuid,54,false,now(),now()),
	 ('cc0a4776-7701-4767-b4d2-4f5b82541f40','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'c4c73fcb-be11-4b1a-986a-a73451d402a7'::uuid,54,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('ef546ec6-fda1-4f57-af8c-ce682dbda223','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'098488af-82c9-49c6-9daa-879eff3d3bee'::uuid,54,false,now(),now()),
	 ('81e703cf-011c-4e2a-b712-a8ceb70d2b18','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'098488af-82c9-49c6-9daa-879eff3d3bee'::uuid,54,false,now(),now()),
	 ('61fecaa7-53a5-400e-b6df-3e17c6db52fb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'b3911f28-d334-4cca-8924-7da60ea5a213'::uuid,21,false,now(),now()),
	 ('ee75e97a-246e-4111-b92c-e292cd1bf370','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'b3911f28-d334-4cca-8924-7da60ea5a213'::uuid,21,false,now(),now()),
	 ('2ee4e217-d2a9-4e9a-9046-add54a35e60b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'f18133b7-ef83-4b2b-beff-9c3b5f99e55a'::uuid,21,false,now(),now()),
	 ('20c85962-08c3-4670-9173-f465916b2e3e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'30040c3f-667d-4dee-ba4c-24aad0891c9c'::uuid,21,false,now(),now()),
	 ('f217ed40-d07f-460b-ac81-1a3fe91682d1','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'a4fa6b22-3d7f-4d56-96f1-941f9e7570aa'::uuid,21,false,now(),now()),
	 ('e7155cf4-aa64-4242-8bcc-a4028f88d490','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'a4fa6b22-3d7f-4d56-96f1-941f9e7570aa'::uuid,21,false,now(),now()),
	 ('5f615022-9836-4679-ae22-2ede6db3d39c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'9b6832a8-eb82-4afa-b12f-b52a3b2cda75'::uuid,21,false,now(),now()),
	 ('7210643e-6e81-4c3e-bbe1-1eca271dffb2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'9b6832a8-eb82-4afa-b12f-b52a3b2cda75'::uuid,21,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('0cfd4a13-8bc6-4116-86f8-f2066f68d00a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'cae0eb53-a023-434c-ac8c-d0641067d8d8'::uuid,42,false,now(),now()),
	 ('ea67540c-4a54-4a6f-b96e-a38db8bb9824','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'cae0eb53-a023-434c-ac8c-d0641067d8d8'::uuid,42,false,now(),now()),
	 ('e6f7904d-5ad6-4008-a338-3fd38a9f5320','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'2a1b3667-e604-41a0-b741-ba19f1f56892'::uuid,21,false,now(),now()),
	 ('1c0d45f0-c22c-4c3f-8982-9a5c19202950','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'2a1b3667-e604-41a0-b741-ba19f1f56892'::uuid,21,false,now(),now()),
	 ('9218f896-a474-4321-bf09-16eb7527b3ff','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'d53d6be6-b36c-403f-b72d-d6160e9e52c1'::uuid,21,false,now(),now()),
	 ('1d69ca96-b1bc-4e33-91e8-9c09b0a9d8b5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'d53d6be6-b36c-403f-b72d-d6160e9e52c1'::uuid,21,false,now(),now()),
	 ('71d5314e-8202-490c-8ac9-d14b6a969b4f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'1a170f85-e7f1-467c-a4dc-7d0b7898287e'::uuid,21,false,now(),now()),
	 ('b3ed4ea4-12a6-420c-b77e-2abf10e22e75','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'1a170f85-e7f1-467c-a4dc-7d0b7898287e'::uuid,21,false,now(),now()),
	 ('feaf7434-3c57-4f55-9de2-eaabfecb6051','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'5e8d8851-bf33-4d48-9860-acc24aceea3d'::uuid,21,false,now(),now()),
	 ('88acb22e-1a4e-46be-9825-09bbe5de628f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'5e8d8851-bf33-4d48-9860-acc24aceea3d'::uuid,21,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('9da18e0e-509d-4216-a77e-dda86173b56c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'7ac1c0ec-0903-477c-89e0-88efe9249c98'::uuid,21,false,now(),now()),
	 ('48ab6006-072f-4976-8fa2-c9f89f33caa2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'7ac1c0ec-0903-477c-89e0-88efe9249c98'::uuid,21,false,now(),now()),
	 ('48d74193-bd1f-4068-b855-f33ca4168ba1','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'0cb31c3c-dfd2-4b2a-b475-d2023008eea4'::uuid,21,false,now(),now()),
	 ('f40f3b06-7202-4e98-b87d-6e956aa3201e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'0cb31c3c-dfd2-4b2a-b475-d2023008eea4'::uuid,21,false,now(),now()),
	 ('dd0e39df-261d-4ca1-8e7e-1bcc92a5fb2d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'8eb44185-f9bf-465e-8469-7bc422534319'::uuid,42,false,now(),now()),
	 ('8fd7dd20-1ec9-4ab6-ac3a-947e81ddc6a0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'8eb44185-f9bf-465e-8469-7bc422534319'::uuid,42,false,now(),now()),
	 ('479efc57-7cf4-4c62-8fba-61c323c9eb2f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'6530aaba-4906-4d63-a6d3-deea01c99bea'::uuid,21,false,now(),now()),
	 ('3c5a1ed5-332f-4009-93a2-3655823a1aa0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'6530aaba-4906-4d63-a6d3-deea01c99bea'::uuid,21,false,now(),now()),
	 ('099dfb7f-8efc-4b1f-8a57-11bdeb353d8c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'182eb005-c185-418d-be8b-f47212c38af3'::uuid,21,false,now(),now()),
	 ('238d64db-0537-4c6f-990a-d7c0aafbe8b7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'182eb005-c185-418d-be8b-f47212c38af3'::uuid,21,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('d232d016-7866-433b-8685-067ba722118a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'4a366bb4-5104-45ea-ac9e-1da8e14387c3'::uuid,40,false,now(),now()),
	 ('9c3f5b86-c9de-4975-bb17-fb937ffd7190','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'4a366bb4-5104-45ea-ac9e-1da8e14387c3'::uuid,40,false,now(),now()),
	 ('89606308-6be5-4738-9dc3-21a339914fc8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'899d79f7-8623-4442-a398-002178cf5d94'::uuid,40,false,now(),now()),
	 ('0bef9d40-3082-4bbd-ab46-20c77a541375','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'899d79f7-8623-4442-a398-002178cf5d94'::uuid,40,false,now(),now()),
	 ('ccc86bef-2736-467e-8906-eb6b2c2e7844','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'dd6c2ace-2593-445b-9569-55328090de99'::uuid,40,false,now(),now()),
	 ('d98cfb35-a8ac-4710-8f4a-d95d9ee39047','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'dd6c2ace-2593-445b-9569-55328090de99'::uuid,40,false,now(),now()),
	 ('d2299b47-ec93-4144-91af-ab3ab77d7bf7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'02cc7df6-83d0-4ff1-a5ea-8240f5434e73'::uuid,20,false,now(),now()),
	 ('c0e31a12-e018-47c2-b438-c4892d352042','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'02cc7df6-83d0-4ff1-a5ea-8240f5434e73'::uuid,20,false,now(),now()),
	 ('1fde495e-4bf5-4429-bac7-d0663ea68d1d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'93052804-f158-485d-b3a5-f04fd0d41e55'::uuid,20,false,now(),now()),
	 ('81cc2a5e-d119-4b3c-b20e-35291844bdeb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'93052804-f158-485d-b3a5-f04fd0d41e55'::uuid,20,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('8c03c054-1787-44b5-93c1-980dd89a2da1','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'6455326e-cc11-4cfe-903b-ccce70e6f04e'::uuid,20,false,now(),now()),
	 ('84f1eddd-c52a-4cfc-82bc-11b23767b16b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'6455326e-cc11-4cfe-903b-ccce70e6f04e'::uuid,20,false,now(),now()),
	 ('5cc506f7-6030-4a0c-831b-93df65371da7','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'0ba534f5-0d24-4d7c-9216-d07f57cd8edd'::uuid,20,false,now(),now()),
	 ('54af4304-a88b-4b6a-bcbf-2e129215945e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'0ba534f5-0d24-4d7c-9216-d07f57cd8edd'::uuid,20,false,now(),now()),
	 ('16d99935-b381-4120-8fa5-af4a4d5d4859','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'433334c3-59dd-404d-a193-10dd4172fc8f'::uuid,20,false,now(),now()),
	 ('bdcf228c-d71f-4c70-a90c-d48d0e023556','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'433334c3-59dd-404d-a193-10dd4172fc8f'::uuid,20,false,now(),now()),
	 ('b496c4fd-a7fe-4d7a-baa7-9fe8e7d21ab3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'e3071ca8-bedf-4eff-bda0-e9ff27f0e34c'::uuid,20,false,now(),now()),
	 ('0f37f000-1999-43ef-89c1-19281250123c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'e3071ca8-bedf-4eff-bda0-e9ff27f0e34c'::uuid,20,false,now(),now()),
	 ('81b9ba17-005f-48e3-9ca8-882f3a19e3b0','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'3320e408-93d8-4933-abb8-538a5d697b41'::uuid,20,false,now(),now()),
	 ('58317f34-e277-4ee4-9855-2c205358207a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'3320e408-93d8-4933-abb8-538a5d697b41'::uuid,20,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('676ceac8-37c1-42f9-a17c-127c0cd7b163','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'9a9da923-06ef-47ea-bc20-23cc85b51ad0'::uuid,20,false,now(),now()),
	 ('cb61d35a-037a-49e1-82ed-c097768aaa99','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'1e23a20c-2558-47bf-b720-d7758b717ce3'::uuid,20,false,now(),now()),
	 ('5099e3c6-9a6b-485c-b70c-8675ce5c2886','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'1e23a20c-2558-47bf-b720-d7758b717ce3'::uuid,20,false,now(),now()),
	 ('5deff542-6a01-466e-89fb-63f3c6851d8b','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'58dcc836-51e1-4633-9a89-73ac44eb2152'::uuid,20,false,now(),now()),
	 ('c59010f9-e636-4c5f-8d1a-f556fc18a7d3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'7ee486f1-4de8-4700-922b-863168f612a0'::uuid,40,false,now(),now()),
	 ('678e0858-96b7-4a41-9137-01efe6f65833','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'7ee486f1-4de8-4700-922b-863168f612a0'::uuid,40,false,now(),now()),
	 ('6c5f38e3-1036-45de-9195-041797d01fb8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'3ec11db4-f821-409f-84ad-07fc8e64d60d'::uuid,40,false,now(),now()),
	 ('cb08dfc4-4d43-49e6-aa49-42a5507c4053','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'3ec11db4-f821-409f-84ad-07fc8e64d60d'::uuid,40,false,now(),now()),
	 ('ab5cc8a2-f054-4856-9528-d61380e46f38','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'dcc3cae7-e05e-4ade-9b5b-c2eaade9f101'::uuid,21,false,now(),now()),
	 ('59013d73-b3d3-4b12-b396-82db23f15f9d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'dcc3cae7-e05e-4ade-9b5b-c2eaade9f101'::uuid,21,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('3b20b61f-c0c4-495f-8b41-b203a240bb03','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'6f0e02be-08ad-48b1-8e23-eecaab34b4fe'::uuid,21,false,now(),now()),
	 ('2358302c-0201-4e48-90e0-92b7ad320092','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'6f0e02be-08ad-48b1-8e23-eecaab34b4fe'::uuid,21,false,now(),now()),
	 ('2cd487ac-dd6d-4f07-9606-338966518b98','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'4f2e3e38-6bf4-4e74-bd7b-fe6edb87ee42'::uuid,21,false,now(),now()),
	 ('c34aa8b6-344c-4304-a09e-910ad675d995','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'4f2e3e38-6bf4-4e74-bd7b-fe6edb87ee42'::uuid,21,false,now(),now()),
	 ('171b5613-9dcd-4ed5-ba6c-66b412e261ed','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'760f146d-d5e7-4e08-9464-45371ea3267d'::uuid,21,false,now(),now()),
	 ('7728b498-029b-47e0-be02-3f081449094f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'760f146d-d5e7-4e08-9464-45371ea3267d'::uuid,21,false,now(),now()),
	 ('8f6dc5bd-8e9c-40b9-b5e1-9f176425ebfa','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'a2fad63c-b6cb-4b0d-9ced-1a81a6bc9985'::uuid,21,false,now(),now()),
	 ('cbb46ae7-14e4-4c97-a993-45caedb2bad4','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'0026678a-51b7-46de-af3d-b49428e0916c'::uuid,21,false,now(),now()),
	 ('a799e032-0925-421c-b552-ec8f53995e4d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'0026678a-51b7-46de-af3d-b49428e0916c'::uuid,21,false,now(),now()),
	 ('2aa9adb4-6980-4a27-ae03-5c906bfea3d3','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'40ab17b2-9e79-429c-a75d-b6fcbbe27901'::uuid,21,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('7e864210-2b97-4a14-9c14-c15becfa3fa8','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'43a09249-d81b-4897-b5c7-dd88331cf2bd'::uuid,21,false,now(),now()),
	 ('eb11aab1-abd4-41d9-8bb6-c901a0c8e8b9','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'43a09249-d81b-4897-b5c7-dd88331cf2bd'::uuid,21,false,now(),now()),
	 ('78410663-770a-4a95-bf39-06d52bcb5231','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'a2fad63c-b6cb-4b0d-9ced-1a81a6bc9985'::uuid,21,false,now(),now()),
	 ('f4117bcb-0391-466b-9c29-1687bee1d4c2','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'422021c7-08e1-4355-838d-8f2821f00f42'::uuid,21,false,now(),now()),
	 ('a0139ddc-d7d6-4882-abd8-81a459574e7c','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'422021c7-08e1-4355-838d-8f2821f00f42'::uuid,21,false,now(),now()),
	 ('87e30935-dd1e-4df7-baeb-20645ab2c1a5','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'e4e467f2-449d-46e3-a59b-0f8714e4824a'::uuid,21,false,now(),now()),
	 ('49075c0c-d769-4c0e-995a-4bb661c11959','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'e4e467f2-449d-46e3-a59b-0f8714e4824a'::uuid,21,false,now(),now()),
	 ('67b8fa57-202a-4af4-b546-01547cc8070a','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'535e6789-c126-405f-8b3a-7bd886b94796'::uuid,21,false,now(),now()),
	 ('25362cbe-db22-4607-91bd-4301d73bc577','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'40ab17b2-9e79-429c-a75d-b6fcbbe27901'::uuid,21,false,now(),now()),
	 ('81fbd38e-40f6-42de-bd50-de8b818cd2e6','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'30040c3f-667d-4dee-ba4c-24aad0891c9c'::uuid,21,false,now(),now());
INSERT INTO re_intl_other_prices (id,contract_id,service_id,is_peak_period,rate_area_id,per_unit_cents,is_less_50_miles,created_at,updated_at) VALUES
	 ('47a63271-a273-41c8-9385-b25e5562c115','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'def8c7af-d4fc-474e-974d-6fd00c251da8'::uuid,54,false,now(),now()),
	 ('390ee2b8-af89-4a0d-b869-31e6a0dcb04e','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'def8c7af-d4fc-474e-974d-6fd00c251da8'::uuid,54,false,now(),now()),
	 ('ae5550c5-6404-45d9-9218-9a3f8b08c83d','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'3733db73-602a-4402-8f94-36eec2fdab15'::uuid,54,false,now(),now()),
	 ('af8f6f53-0abb-43c7-b6e4-c8ff6ea93fc4','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'b80251b4-02a2-4122-add9-ab108cd011d7'::uuid,54,false,now(),now()),
	 ('900b1b87-a1cb-4ffd-9f23-b2fb0daf54b4','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'811a32c0-90d6-4744-9a57-ab4130091754'::uuid,54,false,now(),now()),
	 ('91b1f184-ad47-4f91-8bc2-5448f0c1cdde','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'811a32c0-90d6-4744-9a57-ab4130091754'::uuid,54,false,now(),now()),
	 ('3ec1eceb-4ddc-498c-8830-8be883a67101','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'9a9da923-06ef-47ea-bc20-23cc85b51ad0'::uuid,20,false,now(),now()),
	 ('ed0ad9a3-80c4-4e8d-9c6e-8a1817f6a297','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',true,'535e6789-c126-405f-8b3a-7bd886b94796'::uuid,21,false,now(),now()),
	 ('8028ae10-24e2-46eb-a984-f47f296fea8f','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'f18133b7-ef83-4b2b-beff-9c3b5f99e55a'::uuid,21,false,now(),now()),
	 ('dd428ffc-f0fc-447b-918d-1e72d88144bb','070f7c82-fad0-4ae8-9a83-5de87a56472e'::uuid,'6f4f6e31-0675-4051-b659-89832259f390',false,'58dcc836-51e1-4633-9a89-73ac44eb2152'::uuid,20,false,now(),now());


------------------------------------------------------------------------
-- Map service params to payment request pricing lookup handler
------------------------------------------------------------------------
-- IOPSIT: ContractCode
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at,updated_at, is_optional) VALUES
('b3e3ab12-4588-488b-afd9-3594f4203a29'::uuid, (select id from re_services where code = 'IOPSIT'), (select id from service_item_param_keys where key = 'ContractCode'), now(), now(), false);

-- IOPSIT: PerUnitCents
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at,updated_at, is_optional) VALUES
('6b98203c-9115-4cb8-beef-31d715b8c06e'::uuid, (select id from re_services where code = 'IOPSIT'), (select id from service_item_param_keys where key = 'PerUnitCents'), now(), now(), false);

-- IDDSIT: ContractCode
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at,updated_at, is_optional) VALUES
('68251b89-9b17-42f9-88f4-6cfa8d2a00c8'::uuid, (select id from re_services where code = 'IDDSIT'), (select id from service_item_param_keys where key = 'ContractCode'), now(), now(), false);

-- IDDSIT: PerUnitCents
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at,updated_at, is_optional) VALUES
('dba6e023-dd5d-4a80-a43b-74f5751234af'::uuid, (select id from re_services where code = 'IDDSIT'), (select id from service_item_param_keys where key = 'PerUnitCents'), now(), now(), false);
