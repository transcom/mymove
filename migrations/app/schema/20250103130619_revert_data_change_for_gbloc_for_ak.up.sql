drop view move_to_gbloc;
CREATE OR REPLACE VIEW move_to_gbloc AS
SELECT DISTINCT ON (sh.move_id) sh.move_id AS move_id, COALESCE(pctg.gbloc, pctg_ppm.gbloc) AS gbloc
FROM mto_shipments sh
     -- try the pickup_address path
     LEFT JOIN
     (
        SELECT a.id address_id, pctg.gbloc
        FROM addresses a
        JOIN postal_code_to_gblocs pctg ON a.postal_code = pctg.postal_code
     ) pctg ON pctg.address_id = sh.pickup_address_id
     -- try the ppm_shipments path
     LEFT JOIN
     (
        SELECT ppm.shipment_id, pctg.gbloc
        FROM ppm_shipments ppm
        JOIN addresses ppm_address ON ppm.pickup_postal_address_id = ppm_address.id
        JOIN postal_code_to_gblocs pctg ON ppm_address.postal_code = pctg.postal_code
     ) pctg_ppm ON pctg_ppm.shipment_id = sh.id
WHERE sh.deleted_at IS NULL
ORDER BY sh.move_id, sh.created_at;

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99501', 'MBFL', now(), now(), 'd8697416-e345-46a8-9767-47b7abbb3c06');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99502', 'MBFL', now(), now(), '0e060122-7cea-4bcd-b636-31c453088a5d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99503', 'MBFL', now(), now(), '1bc710db-3f4f-4177-9656-a99956a8b06e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99504', 'MBFL', now(), now(), 'e35c4fd3-b8d7-46f4-a559-ddfc8e0ad4e9');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99505', 'MBFL', now(), now(), 'f2e40ed3-bc7b-428c-9693-3bc1cc1a9c57');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99506', 'MBFL', now(), now(), 'cf3890e6-16df-46a7-aabb-57010b999ee7');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99507', 'MBFL', now(), now(), 'c869c7aa-e0fd-4933-a4e9-09bd6191b25a');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99508', 'MBFL', now(), now(), 'ca8554c3-d21c-4e26-a77f-6e965c7d31c5');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99509', 'MBFL', now(), now(), '54ea9592-d93d-4102-84ce-7a9c57b6aff8');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99510', 'MBFL', now(), now(), 'f8d0e922-ab5d-4871-bd0c-1f7484b13981');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99511', 'MBFL', now(), now(), '7098c10e-edbf-4cf6-82dd-c09dd7a8f226');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99513', 'MBFL', now(), now(), '8f2d3b79-718e-4d96-b5d2-6b0c40487332');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99514', 'MBFL', now(), now(), '8ff6d888-935c-488c-9b62-33b8f4053b4a');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99515', 'MBFL', now(), now(), 'e3927bca-6677-49d7-a622-ee20f386c28c');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99516', 'MBFL', now(), now(), 'fc5ee3c2-d0b6-4e56-9bcd-e802e14dca7e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99517', 'MBFL', now(), now(), '84c51c78-972a-4714-8218-8169ee541159');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99518', 'MBFL', now(), now(), '7b3e4eed-3584-4036-9da5-9bc6a7163e4e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99519', 'MBFL', now(), now(), '62b97b86-e8e7-410a-a267-d101ce5f9a3c');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99520', 'MBFL', now(), now(), 'bce10502-90fb-4f28-8489-f5e6d2ac00ab');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99521', 'MBFL', now(), now(), 'b8c51570-d048-41f8-a1fc-58e6b41d367a');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99522', 'MBFL', now(), now(), '6311d142-f02b-4f49-b719-608d78c91489');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99523', 'MBFL', now(), now(), '967028d9-3a7b-4949-b1a4-e94a1aa25a73');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99524', 'MBFL', now(), now(), 'afd3cc8b-2c02-4e63-8fd1-eb81a1a19b5d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99529', 'MBFL', now(), now(), 'f08e218c-e55e-43fc-bbdf-bc57aaa65726');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99530', 'MBFL', now(), now(), '65eeaf36-e786-43dd-a18c-fb470a4468f9');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99540', 'MBFL', now(), now(), '2380496a-9087-48bc-a0e2-9cd5ac19f470');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99545', 'MBFL', now(), now(), '6df5025e-e730-427a-8c10-ff6d131d0567');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99546', 'MBFL', now(), now(), '02163620-1da6-4509-b8bb-d8b3336e5c2a');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99547', 'MBFL', now(), now(), 'f869b729-88cd-4770-a435-2dfd4b50c330');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99548', 'MBFL', now(), now(), 'db80b48b-25bd-4ef1-8d5a-33d740af6b0e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99549', 'MBFL', now(), now(), '09983e53-ae83-41b1-b459-e0ab411ea87e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99550', 'MAPS', now(), now(), '643db073-5e92-43c0-9b42-a8d9cc800623');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99551', 'MBFL', now(), now(), '9bc29931-880b-42df-a758-017d5abee32d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99552', 'MBFL', now(), now(), 'f6fa0bb7-fd5c-494e-b62c-b7bc509b632f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99553', 'MBFL', now(), now(), 'abcb2cbf-389d-4830-ae16-037c55b7bc2c');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99554', 'MBFL', now(), now(), '2d4d9111-f59b-4506-935d-b6601275311a');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99555', 'MBFL', now(), now(), '00c2d5ed-2364-490d-86b9-40cbad4679fb');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99556', 'MBFL', now(), now(), 'b5c75c77-59cf-4089-815b-c0d15d9440f0');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99557', 'MBFL', now(), now(), '999e2c9e-04ed-47da-aca1-2b320ed52666');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99558', 'MBFL', now(), now(), 'cdc2383c-9ead-4501-a6cc-b736ea72ecef');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99559', 'MBFL', now(), now(), '2ca0657c-c121-41e4-a741-b563dee6838e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99561', 'MBFL', now(), now(), '0c8df634-b6e6-457f-93d7-7d64909bc7cb');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99563', 'MBFL', now(), now(), 'e9a0ebee-3ae3-4d66-b18b-cc3ea417887d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99564', 'MBFL', now(), now(), '77dd08d9-0997-4d16-babd-1ffccc4222d4');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99565', 'MBFL', now(), now(), '59b6a797-b2e8-4bd7-affc-a5639308a9b3');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99566', 'MBFL', now(), now(), 'cf5b630e-964a-4389-9fa2-e1ed0e6a3041');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99567', 'MBFL', now(), now(), 'f71a916b-e55f-4230-869e-038c319037ce');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99568', 'MBFL', now(), now(), '9719d0c5-1bd0-42fe-89ab-f4dbddcfa588');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99569', 'MBFL', now(), now(), 'fba4bdcf-7707-4826-bf2a-3e7aaf81f309');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99571', 'MBFL', now(), now(), '75e598ff-6d6c-4a80-80bc-de36bc37a3be');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99572', 'MBFL', now(), now(), 'f7b9b7d5-6730-4e8f-9364-e5b1eb1d2c2b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99573', 'MBFL', now(), now(), '74dcb198-d7c8-4dbc-81ab-33a16559e100');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99574', 'MAPK', now(), now(), '9b9d03a0-e069-48cd-ae65-a66dbd8a3214');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99575', 'MBFL', now(), now(), 'f6ffb9e8-8976-4288-88cd-d10420d1894e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99576', 'MBFL', now(), now(), '0883cf7f-c2f9-4865-91c3-34df59486f01');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99577', 'MBFL', now(), now(), 'a5e58cce-5375-498b-bfdb-d2e4b50ee2b7');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99578', 'MBFL', now(), now(), '63113a71-6f32-4531-991b-015d51be0ef7');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99579', 'MBFL', now(), now(), 'be2f7dfb-a020-45b5-afa6-0e3b6fd669bf');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99580', 'MBFL', now(), now(), '29b1f268-8ad5-400c-8aa6-0b4b851f8588');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99581', 'MBFL', now(), now(), '5e1dee94-93f9-4b7d-8445-7de967031bfa');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99583', 'MBFL', now(), now(), '2ab06e55-ae2b-4c19-b679-81ee89d098a5');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99585', 'MBFL', now(), now(), 'd177b210-8362-4073-9862-43333203abfc');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99586', 'MBFL', now(), now(), 'c8bdc299-d65d-4b4b-a0ad-0ef6c8ca905d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99587', 'MBFL', now(), now(), '9543c328-16d2-484b-9e62-1965a63f1d2e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99588', 'MBFL', now(), now(), '08d82a7f-8417-407d-b352-9e0f3ef9e9cd');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99589', 'MBFL', now(), now(), 'e09e4805-d0a9-4347-8f63-7e221c4b6d7e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99590', 'MBFL', now(), now(), '3022675c-9e50-403e-987c-7e79762274c2');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99591', 'MBFL', now(), now(), '57eebbd9-4b15-4b6f-923a-1d1db65d72bc');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99599', 'MBFL', now(), now(), '78a46448-b22a-4da0-a42a-07af993848bc');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99602', 'MBFL', now(), now(), '101c94fc-11cf-4317-a63b-c3c903d1f9b3');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99603', 'MBFL', now(), now(), '67ebaf88-8131-4215-b6a4-3b3ecb3862f0');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99604', 'MBFL', now(), now(), '45e5dfdd-9dbc-439a-96ac-292b820aa292');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99605', 'MBFL', now(), now(), 'd403b760-8ad8-409b-b8d4-e0ad9dfd3947');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99606', 'MBFL', now(), now(), 'ac7a083e-3300-4881-874b-6bbd206a4e92');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99607', 'MBFL', now(), now(), 'f8d072a2-099d-4767-9b1d-b5dd48a79a19');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99608', 'MAPS', now(), now(), '889e4f60-6352-443c-8705-734ed91dc7a4');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99609', 'MBFL', now(), now(), 'da00dfc3-a372-4e9f-8ddb-f0d9c792ab36');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99610', 'MBFL', now(), now(), '53b8b7d8-337a-43fe-9a97-e2f5083210aa');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99611', 'MBFL', now(), now(), 'f7211d4f-8ffb-482f-b268-07b9d79c467e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99612', 'MBFL', now(), now(), 'e7340a1b-3278-42aa-a067-021863d676dd');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99613', 'MBFL', now(), now(), '7de673ca-c91d-4bc1-b450-fdaa0c199723');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99614', 'MBFL', now(), now(), 'fd72416c-c88b-4c0f-a3e0-28fd924fd7d7');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99615', 'MAPS', now(), now(), 'b0597a99-16d7-4e1e-9820-042971b2551f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99619', 'MAPS', now(), now(), '8d05d190-fda5-4930-838d-76f37892417c');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99620', 'MBFL', now(), now(), 'fc2456f7-d2ca-4cb3-b16b-f50388c09a12');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99621', 'MBFL', now(), now(), 'da74789a-b710-4332-a973-cc8ff104f80f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99622', 'MBFL', now(), now(), '52b0c551-7d40-496b-9e5a-3aeb7e68f415');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99623', 'MBFL', now(), now(), '4613c99b-da96-4649-8de5-8e445af78f99');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99624', 'MAPS', now(), now(), '13e0d2ba-852e-4970-b60c-464740975141');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99625', 'MBFL', now(), now(), '772070b7-2029-40cf-82b7-e8b839f2740f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99626', 'MBFL', now(), now(), 'a6a4b405-e1af-455f-89d5-088de0afebfa');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99627', 'MBFL', now(), now(), 'e0778102-5fb8-4366-bbd2-413a96946a6d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99628', 'MBFL', now(), now(), '4dad6942-608b-415b-86b1-4575cbe92b14');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99629', 'MBFL', now(), now(), '534a58f3-e96b-41cd-86f5-d3c3f7553063');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99630', 'MBFL', now(), now(), 'aca18cb8-d348-4dfa-bfce-3508ed050ce2');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99631', 'MBFL', now(), now(), '90fd8694-e43b-49f7-a22b-27f95170069c');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99632', 'MBFL', now(), now(), 'b6736f3e-b561-4658-a22a-822b4d141db6');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99633', 'MBFL', now(), now(), 'e622bb3a-8c83-4230-b953-a256a7ed2b50');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99634', 'MBFL', now(), now(), 'a042efbf-24da-489f-ab6d-15179847cec1');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99635', 'MBFL', now(), now(), '153a77b9-ae9f-4e3c-8b7e-35631300491d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99636', 'MBFL', now(), now(), '988b8373-ca1d-459f-a0cc-598c20ad1e07');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99637', 'MBFL', now(), now(), 'eb47a507-5f66-4344-bed9-d165d8275fa9');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99638', 'MBFL', now(), now(), '0858bd59-a16b-47ca-9cb8-670411737d4d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99639', 'MBFL', now(), now(), '52b0ac8b-3931-47d1-aa66-2283a9a1650a');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99640', 'MBFL', now(), now(), '9040a7e0-de08-494b-8d69-27cd4a53bc90');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99641', 'MBFL', now(), now(), '6e1d8641-0448-40cc-88e3-67d4f8466a8e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99643', 'MAPS', now(), now(), '15c1aea9-9b9b-4f6e-93cf-a114ac9175f8');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99644', 'MAPS', now(), now(), '6f294618-1e61-42fc-8d6b-cd870b0b08fc');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99645', 'MBFL', now(), now(), '82d2557f-263b-4108-a75b-ff745aea08f0');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99647', 'MBFL', now(), now(), 'c432e858-d177-4faa-8d0b-97c5305b61e7');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99648', 'MBFL', now(), now(), '273ecc7f-0b37-439e-97e6-86a2f7ea7114');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99649', 'MBFL', now(), now(), '213d21e6-0b04-4f5c-90f8-4ddc9938381d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99650', 'MBFL', now(), now(), 'f2a5f071-205f-4cc1-9448-44058d9af3f5');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99651', 'MBFL', now(), now(), 'c16c38b5-eb50-4dc7-8a40-aa9e7e80bf2f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99652', 'MBFL', now(), now(), 'ef17b53b-f840-4b2d-9f66-3f80c413b3c1');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99653', 'MBFL', now(), now(), '7307956d-84f3-4c22-ad2f-40b62dc0a0d7');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99654', 'MBFL', now(), now(), '47ccb8d8-8df2-4293-9f86-f5703aea0dba');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99655', 'MBFL', now(), now(), '69d0d1e7-8c8e-4cde-82dc-b8374b8b4861');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99656', 'MBFL', now(), now(), 'e0be218f-2d2f-4bac-b8b5-e09ef3f28742');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99657', 'MBFL', now(), now(), 'bd121ed0-6592-427a-9ba9-c06ab3a703dd');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99658', 'MBFL', now(), now(), '6acfa9fe-137e-4a92-b21c-39e627cc378a');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99659', 'MBFL', now(), now(), 'f131070d-b5d8-453e-bac7-ad436ea124dc');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99660', 'MBFL', now(), now(), 'cab88f64-c8c1-42f1-be62-7b4cce31aeb3');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99661', 'MBFL', now(), now(), 'e5a62394-0278-42a3-84e3-57ff2240dc60');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99662', 'MBFL', now(), now(), '8c0d1e5d-6afe-4552-b7cc-3b5a6fbdc3db');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99663', 'MBFL', now(), now(), '44beec24-6856-43f7-ab55-b48286d2829b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99664', 'MBFL', now(), now(), 'c35ed3e5-c3a6-43e2-9794-56709f93620f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99665', 'MBFL', now(), now(), '06b6f8ed-7f4e-4b7a-94b6-def7b619a03d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99666', 'MBFL', now(), now(), '959c959f-b829-471c-a63e-b3cffa15f6c7');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99667', 'MBFL', now(), now(), '75ef0a3a-cc89-4a3b-809f-6cdb6753bc25');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99668', 'MBFL', now(), now(), 'c2f65e88-5c6a-4e46-95c6-89e79385c163');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99669', 'MBFL', now(), now(), '6c1d51b0-9213-404a-b7d8-c4fc384bd304');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99670', 'MBFL', now(), now(), '523b6896-d660-42c7-a8f9-e4603e75ccc9');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99671', 'MBFL', now(), now(), '2c1d481b-6ef3-49ab-8ae9-a9e6429575fd');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99672', 'MBFL', now(), now(), '702c69d2-5231-4267-8d8e-85d5817077d7');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99674', 'MBFL', now(), now(), '66b7f4fb-9098-40a8-a018-a5217ed9878d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99675', 'MBFL', now(), now(), '2ab46ff9-9120-48af-ac05-73663dfd1571');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99676', 'MBFL', now(), now(), '488978c8-3c52-4a1a-9966-c04bba23ee37');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99677', 'MBFL', now(), now(), '30ceccbf-25b6-4862-b358-78f12e22ba06');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99678', 'MBFL', now(), now(), 'a4ed1190-85d9-4b20-8390-9a1b9a3bc0c4');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99679', 'MBFL', now(), now(), '1c87ddc2-ad88-4607-b826-f37da5fc8762');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99680', 'MBFL', now(), now(), '3b5dca39-d435-413c-9edf-3bb3a2b5005d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99681', 'MBFL', now(), now(), '93acc862-3107-4f22-b270-e0c7d1f5f80b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99682', 'MBFL', now(), now(), 'd5224770-83e3-4dc8-97af-24c147669dcc');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99683', 'MBFL', now(), now(), 'b9dcbc66-53d9-4f14-b51d-02171a3fdf5f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99684', 'MBFL', now(), now(), '65a07b11-25cf-4fae-baf4-59d89677cfc2');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99685', 'MBFL', now(), now(), '0cc24afe-3579-41bf-8181-08725b45e3b0');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99686', 'MBFL', now(), now(), 'fd40734b-f683-43b2-9e6e-93eeda219aa9');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99687', 'MBFL', now(), now(), 'fe3c676b-a3d7-4d58-88d4-5a6ab3dc3aae');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99688', 'MBFL', now(), now(), '83cc5417-31e0-468c-a5b3-0495734e9504');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99689', 'MAPK', now(), now(), '68962696-2ba8-40ae-946d-8ab29dd3cd4b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99690', 'MBFL', now(), now(), '92d71dc8-6601-4cc0-bd76-3584e402ceed');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99691', 'MBFL', now(), now(), '839bfff4-f486-462f-bf2c-560aa95f5b6e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99692', 'MBFL', now(), now(), 'b3e45fb6-3e5f-452e-b355-08fa49e1f52b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99693', 'MBFL', now(), now(), 'd3843234-0168-402a-821d-9a4666a54273');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99694', 'MBFL', now(), now(), '3ca39ea7-1f6d-4112-886e-cd4c6fc07476');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99695', 'MBFL', now(), now(), 'b865dfa9-c1c4-42bd-9253-8fa643ebf582');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99697', 'MAPS', now(), now(), '8e71753d-c45a-4a84-912e-cb6bb8a0ccef');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99701', 'MBFL', now(), now(), '951282d4-0523-450b-a636-c2bdadf2a38a');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99702', 'MBFL', now(), now(), 'fcab1e10-4645-4d3e-a4b7-00a9f54a2157');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99703', 'JEAT', now(), now(), '5f47676a-520d-4106-9889-009b0ab29726');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99704', 'MBFL', now(), now(), '328fd045-2e70-4bc2-a466-1319e2425b1f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99705', 'MBFL', now(), now(), '37727f3a-6c6e-4d9e-890c-ec0279c7de0b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99706', 'MBFL', now(), now(), 'c84810a5-660f-4200-a86f-46a17f038adf');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99707', 'MBFL', now(), now(), '7d090646-098b-4096-9012-b68182cf2e35');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99708', 'MBFL', now(), now(), '7a16c7b2-e624-4933-9dc4-b3be69849fe5');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99709', 'MBFL', now(), now(), 'bd590cc5-a624-4fb2-b37b-4b37a7ac863a');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99710', 'MBFL', now(), now(), '5fab3c20-38ee-44b6-9ae1-384b93462637');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99711', 'MBFL', now(), now(), '1d3a88e9-a0e0-4cce-b20a-9f0b6996b7e1');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99712', 'MBFL', now(), now(), 'acec5c52-6709-48cb-a6cb-0e58b05860ec');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99714', 'MBFL', now(), now(), 'ce621370-6656-480d-a071-89de4607d716');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99716', 'MBFL', now(), now(), 'b65f161a-0f33-4360-8905-c716b17943d4');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99720', 'MBFL', now(), now(), '19aa130f-0dee-46db-a054-729cddf9a257');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99721', 'MBFL', now(), now(), '715fde06-a1a7-4240-8e59-4b76dffdb144');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99722', 'MBFL', now(), now(), '4be97e0b-4237-4906-bb0d-8e1704e0b65b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99723', 'MBFL', now(), now(), 'aad84c1d-7ba2-40de-92e5-74b2dd558fa9');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99724', 'MBFL', now(), now(), '8b8157b2-0c64-477f-a173-e26818caeffd');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99725', 'MBFL', now(), now(), '8f258372-966a-4565-9e22-61cdb5f88ef9');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99726', 'MBFL', now(), now(), '02f07ef0-ebad-4fb8-a50c-bdd140d87c49');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99727', 'MBFL', now(), now(), '63de42c5-eb6e-4471-9170-8d4d3fd88d87');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99729', 'MBFL', now(), now(), 'c12cf5c5-b3e0-4dff-a88f-39b8a22a950b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99730', 'MBFL', now(), now(), '820b8fe5-2c4c-459a-850a-8326fd50ea57');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99731', 'JEAT', now(), now(), 'dd366219-6c3b-45c1-a59a-e04d0efe1a6d');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99732', 'MBFL', now(), now(), '32e2f9ec-805a-419b-a850-e4868c8eb13f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99733', 'MBFL', now(), now(), '9762e6a0-a2dd-4da7-9cc5-4dfdfd21fad2');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99734', 'MBFL', now(), now(), 'f35587d4-41be-4e7a-8cc3-07e1f2a9e4b9');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99736', 'MBFL', now(), now(), 'a71189c1-784c-4d9c-925b-dcc6432c6a28');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99737', 'MBFL', now(), now(), '80a884f0-058a-46c0-8fdb-90d1d027f002');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99738', 'MBFL', now(), now(), '34cf2f93-66b3-43ab-a428-2cc2e4c6b143');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99739', 'MBFL', now(), now(), 'eadc1ca9-6de2-46ab-a015-acb813599670');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99740', 'MBFL', now(), now(), 'cc3f07af-485b-4256-8703-bf41a5967497');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99741', 'MBFL', now(), now(), '61409776-e0fa-47f4-825e-9c6f91ce2b10');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99742', 'MBFL', now(), now(), 'b46d81ef-9faa-49ec-85a3-4c3734d07126');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99743', 'MBFL', now(), now(), 'acb28741-502a-40c2-aa45-4a9b42d39405');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99744', 'MBFL', now(), now(), 'baf46567-31ea-481d-9ab9-e974ea261c2f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99745', 'MBFL', now(), now(), 'babaac73-1c6b-4209-baf3-1e0c5d41dc7f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99746', 'MBFL', now(), now(), 'ddc32535-73e0-4f13-ad2e-6294700c17cc');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99747', 'MBFL', now(), now(), '40a3c816-79a0-4f26-aa6a-2847835a45d1');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99748', 'MBFL', now(), now(), '5a3a2cce-4bc0-4bd2-98fe-7f71bc09c99e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99749', 'MBFL', now(), now(), '212653a0-4a5a-4e48-b301-78b3f294fe16');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99750', 'MBFL', now(), now(), '2013baa1-1f2f-47f7-beb0-98f3171c6bf1');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99751', 'MBFL', now(), now(), 'a19e22ef-194e-4b48-858c-59ea04d48fcf');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99752', 'MBFL', now(), now(), 'e1ce6564-37d3-4fe4-8c19-f39e31c9d7cd');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99753', 'MBFL', now(), now(), '2dd6633b-4a67-4bff-879f-32bf88ec64a5');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99754', 'MBFL', now(), now(), '5fe4e773-573e-4052-b4b7-d5686366d03b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99755', 'MBFL', now(), now(), '0383ed06-ff31-4fa7-851a-962dc7adf7d3');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99756', 'MBFL', now(), now(), '61726154-40a0-4f7d-82d0-b09f28a2f209');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99757', 'MBFL', now(), now(), '3bd35ddb-4d89-4e01-8dd1-e938bc5c0388');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99758', 'MBFL', now(), now(), '31c3a401-cc17-4e5a-bf22-d66e00a05528');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99759', 'MBFL', now(), now(), 'cc534b49-e72f-4a19-bbad-4dc5101c1373');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99760', 'MBFL', now(), now(), '56e2c9c0-f914-40f3-988c-7b88593f5513');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99761', 'MBFL', now(), now(), 'ccd849d7-e2ed-40ce-a353-9c93d57b4f92');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99762', 'MBFL', now(), now(), '1981bcff-5f4b-4358-92ea-93177783215f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99763', 'MBFL', now(), now(), '298d0d8a-a19a-4498-937f-1ab233a85730');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99764', 'MBFL', now(), now(), '37352211-832f-4af7-8655-414ac02b8544');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99765', 'MBFL', now(), now(), 'c2dd65ee-9875-4e7d-9d34-46f529f8959b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99766', 'MBFL', now(), now(), 'ad3283f4-74cb-403f-84a2-c14f510f83cd');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99767', 'MBFL', now(), now(), 'eb25b5f7-ed3b-4860-b337-5ee68b6c4534');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99768', 'MBFL', now(), now(), '92116116-8a7a-420c-83cd-5f9977401941');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99769', 'MBFL', now(), now(), 'c5a38d76-3bef-4e9d-8244-1741aa20ed10');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99770', 'MBFL', now(), now(), 'e7f0f266-019d-4bce-8921-6804e5e8cb0c');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99771', 'MBFL', now(), now(), '74f9319a-0970-4721-8fda-bfe4f777896b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99772', 'MBFL', now(), now(), 'db7dc663-1443-4209-9a01-82f694a4dca1');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99773', 'MBFL', now(), now(), '474d4d25-1314-4bbf-ae04-25ffae18a7d4');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99774', 'MBFL', now(), now(), 'ea503792-8747-465c-ad89-c513004e1345');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99775', 'MBFL', now(), now(), 'e33ba93c-8ba2-4a3b-b8bc-f331a989e388');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99776', 'MBFL', now(), now(), 'ffc00def-ec8b-413d-9601-bf5587b510f4');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99777', 'MBFL', now(), now(), 'ce0ff3db-7d59-43af-923d-efdf4feb72cf');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99778', 'MBFL', now(), now(), '4235599e-4fd3-46e9-b733-de29dca74700');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99780', 'MBFL', now(), now(), '0aa94d7f-81bb-42f4-be8a-e14b2bc630c8');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99781', 'MBFL', now(), now(), '1b98f0bd-711c-48f0-87ed-9e5bff30ad60');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99782', 'MBFL', now(), now(), 'c7f7dd5b-7674-4b78-b1a5-0013fd979231');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99783', 'MBFL', now(), now(), '20152935-a8c8-45ad-a4cd-af2a945a8a17');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99784', 'MBFL', now(), now(), '58ec9aaa-e2d0-4732-9ffd-3b7563bb1c05');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99785', 'MBFL', now(), now(), '13f3eb0b-1534-4773-9fae-a86d39a829ae');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99786', 'MBFL', now(), now(), 'eeb04da3-f067-469b-8703-c3e9de365ad7');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99788', 'MBFL', now(), now(), '224088e2-afba-4b23-8361-9126c00445a8');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99789', 'MBFL', now(), now(), '8c123f36-f68d-461b-bd60-01cbd501ae17');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99790', 'MBFL', now(), now(), '1492fc25-8a4a-4fbc-a9be-433ad03d878e');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99791', 'MBFL', now(), now(), 'fda0d6c4-7a68-425b-9567-da8d368dbfa0');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99801', 'MAPK', now(), now(), '095789bc-136f-44e6-bf3f-b5a5a4ffad5c');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99802', 'MAPK', now(), now(), '1e73fbe5-0924-42dc-aedd-9c4aab511b90');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99803', 'MAPK', now(), now(), '281f58c0-fe6e-4e77-9035-6f1db1d4a8fd');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99811', 'MAPK', now(), now(), '6dcd4248-9c66-4bf3-8f4b-30dfd4add731');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99812', 'MAPK', now(), now(), '82579769-f363-4d68-b6f4-2586ac3f974c');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99820', 'MAPK', now(), now(), 'aecd46d3-dea5-4d0e-8409-b5b300ddded8');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99821', 'MAPK', now(), now(), '5a3ff675-8c07-4058-b5c6-8c19400f2f07');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99824', 'MAPK', now(), now(), 'e0cd0a98-4c4e-4d21-9adc-416fb4ce2449');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99825', 'MAPK', now(), now(), 'aa515936-aa46-4016-a487-c925b4dc16b2');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99826', 'MAPK', now(), now(), '553f410a-ed32-4517-a766-c665369d9369');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99827', 'MAPK', now(), now(), 'a275e877-48ba-4473-8064-43b3ab7b042b');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99829', 'MAPK', now(), now(), '7cb1babd-df8b-400f-9e5c-6c595bbcaa2f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99830', 'MAPK', now(), now(), 'dcc9dd70-0984-4555-a371-e8f8bf583ecf');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99832', 'MAPK', now(), now(), 'a3740468-c12a-465d-b627-fd80bc921eb5');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99833', 'MAPK', now(), now(), '38eecaf7-4c2d-44c5-a8bd-ccde44ec6b89');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99835', 'MAPK', now(), now(), 'c914333a-4c88-404b-a17a-f51b4df9dac7');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99836', 'MAPK', now(), now(), 'fe065064-61f7-4c01-80f4-945fb1d6d26f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99840', 'MAPK', now(), now(), '7c72b250-4d4b-41a2-b548-97a5b1e1c526');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99841', 'MAPK', now(), now(), '1a1c28b2-d5c7-4180-93ea-d66887ce7e26');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99850', 'MAPK', now(), now(), '21f9854e-2a61-4810-a4e3-ea09f1f2c59f');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99901', 'MAPK', now(), now(), '36ef16e6-d9a5-4706-a115-e11691ff1c37');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99903', 'MAPK', now(), now(), '76d5df09-e4f1-4b96-8df4-68bbeee29b00');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99918', 'MAPK', now(), now(), '4b6b34d9-40ec-459a-a337-5fbdde3aaa84');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99919', 'MAPK', now(), now(), 'af79bdf9-f4cf-494e-ad67-ddbb836e98d0');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99921', 'MAPK', now(), now(), '742e78df-6615-406b-b0f2-55deaec1fc09');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99922', 'MAPK', now(), now(), '99701216-788c-41d7-bdbb-98ff0f8e2224');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99923', 'MAPK', now(), now(), '3978cb58-70b6-4c23-a97b-d6e00cf7abcf');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99925', 'MAPK', now(), now(), 'bc7cb635-6287-4ecb-8f8c-9e6f66858441');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99926', 'MAPK', now(), now(), '9abec1b8-3f19-4616-aae2-8c8a54b44e48');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99927', 'MAPK', now(), now(), 'ba7adb39-9897-4e4d-b042-c4360cc25691');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99928', 'MAPK', now(), now(), '4f721fb2-7bcc-4a2a-a24b-cd3d00532311');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99929', 'MAPK', now(), now(), '0ed1ee93-f720-4e4c-b28a-4572e847ad15');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99950', 'MAPK', now(), now(), '23c0f299-2b5c-44ff-b604-5efd70019f8d');