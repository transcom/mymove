--
-- Add service items for international NTS shipments.
--
INSERT INTO re_service_items
(id, service_id, shipment_type, market_code, is_auto_approved, created_at, updated_at, sort)
VALUES
    --ISLH International Shipping & Linehaul
    ('2a560507-db09-4be1-b809-49c0f515b31b', '9f3d551a-0725-430e-897e-80ee9add3ae9' ,'HHG_INTO_NTS', 'i', true, now(), now(), 1),
    --PODFSC International POD Fuel Surcharge
    ('e702818f-defd-452c-81a3-865b902e7dd0', '388115e8-abe9-441d-96cf-a39f24baa0a3' ,'HHG_INTO_NTS', 'i', true, now(), now(), 2),
    --INPK International NTS packing
    ('366ee5a4-eb61-4ded-a68c-ddc29fe1a886', '874cb86a-bc39-4f57-a614-53ee3fcacf14' ,'HHG_INTO_NTS', 'i', true, now(), now(), 3),
    --ICRT International crating
    ('aac4e95e-27ed-4f09-9b6b-384d8542f410', '86203d72-7f7c-49ff-82f0-5b95e4958f60' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IDASIT International destination add'l day SIT
    ('010f2f91-7381-4149-8d74-8eb5f593a864', '806c6d59-57ff-4a3f-9518-ebf29ba9cb10' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IDDSIT International destination SIT delivery
    ('a41966b7-b83a-4eaf-8e68-d5e884777102', '28389ee1-56cf-400c-aa52-1501ecdd7c69' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IDFSIT International destination 1st day SIT
    ('14c77957-3c76-465a-bb07-c98d36ef1e54', 'bd6064ca-e780-4ab4-a37b-0ae98eebb244' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IDSHUT International destination shuttle service
    ('d52d2d03-100a-4ed9-b2de-16eac63a375f', '22fc07ed-be15-4f50-b941-cbd38153b378' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IOASIT International origin add'l day SIT
    ('7fd91408-7d69-4375-b7e6-5b2ff714206b', 'bd424e45-397b-4766-9712-de4ae3a2da36' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IOFSIT International origin 1st day SIT
    ('b3dc509d-d652-4300-a702-a1ddce6255b6', 'b488bf85-ea5e-49c8-ba5c-e2fa278ac806' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IOPSIT International origin SIT pickup
    ('001eadb6-3526-45b9-96e0-0648bb481e86', '6f4f6e31-0675-4051-b659-89832259f390' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IOSHUT International origin shuttle service
    ('b991c5c9-af2c-4146-b999-1d0bdf91de3f', '624a97c5-dfbf-4da9-a6e9-526b4f95af8d' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IUCRT International uncrating
    ('5a89315a-257b-4ef0-92cb-4c06aa6f1332', '4132416b-b1aa-42e7-98f2-0ac0a03e8a31' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IOFSC International Origin SIT Fuel Surcharge
    ('d4a98dea-a5f7-4b92-b5de-e6350ab07824', '81e29d0c-02a6-4a7a-be02-554deb3ee49e' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL),
    --IDSFSC International Destination SIT Fuel Surcharge
    ('eaea90c2-93d3-4db9-89cd-23ac57ec9ce1', '690a5fc1-0ea5-4554-8294-a367b5daefa9' ,'HHG_INTO_NTS', 'i', false, now(), now(), NULL);
