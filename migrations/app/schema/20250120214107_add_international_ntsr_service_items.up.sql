--
-- Add service items for international NTS-R shipments.
--
INSERT INTO re_service_items
(id, service_id, shipment_type, market_code, is_auto_approved, created_at, updated_at, sort)
VALUES
    --ISLH International Shipping & Linehaul
    ('bf76fb0f-408a-4391-8aa7-92908f3c027a', '9f3d551a-0725-430e-897e-80ee9add3ae9' ,'HHG_OUTOF_NTS', 'i', true, now(), now(), 1),
    --PODFSC International POD Fuel Surcharge
    ('db2106c8-887c-4304-aad2-c7413de13cc4', '388115e8-abe9-441d-96cf-a39f24baa0a3' ,'HHG_OUTOF_NTS', 'i', true, now(), now(), 2),
	--POEFSC International POD Fuel Surcharge
    ('509a491d-cddd-476c-9ba0-65077cf93e58', 'f75758d8-2fcd-40ba-9432-3ff3032a71d1' ,'HHG_OUTOF_NTS', 'i', true, now(), now(), 2),
    --INPK International NTS packing
    ('4d348ec0-a278-4038-b061-6a4e17ea6721', '874cb86a-bc39-4f57-a614-53ee3fcacf14' ,'HHG_OUTOF_NTS', 'i', true, now(), now(), 3),
    --ICRT International crating
    ('b8f4e434-0912-44c5-b824-c60e5c5dffee', '86203d72-7f7c-49ff-82f0-5b95e4958f60' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IDASIT International destination add'l day SIT
    ('7135540f-602c-4d02-ba66-403b7252738e', '806c6d59-57ff-4a3f-9518-ebf29ba9cb10' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IDDSIT International destination SIT delivery
    ('5d3261a5-a7af-4133-b2e0-2a06f694c551', '28389ee1-56cf-400c-aa52-1501ecdd7c69' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IDFSIT International destination 1st day SIT
    ('90e022da-9944-4563-98a3-6eb7cabb017e', 'bd6064ca-e780-4ab4-a37b-0ae98eebb244' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IDSHUT International destination shuttle service
    ('ba6c218b-dd99-4ef4-87ed-421581218bbf', '22fc07ed-be15-4f50-b941-cbd38153b378' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IOASIT International origin add'l day SIT
    ('ded116f5-ca9d-465e-acd8-3eee899e9713', 'bd424e45-397b-4766-9712-de4ae3a2da36' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IOFSIT International origin 1st day SIT
    ('03114a6f-22ef-4664-9269-b76636466285', 'b488bf85-ea5e-49c8-ba5c-e2fa278ac806' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IOPSIT International origin SIT pickup
    ('03e7b2fd-431d-4ce6-a640-78de846095c9', '6f4f6e31-0675-4051-b659-89832259f390' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IOSHUT International origin shuttle service
    ('94b1786f-86ec-4736-8dbc-6a0e29b64272', '624a97c5-dfbf-4da9-a6e9-526b4f95af8d' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IUCRT International uncrating
    ('cefe3094-5670-41c3-b6cd-72730bbb8fc7', '4132416b-b1aa-42e7-98f2-0ac0a03e8a31' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IOFSC International Origin SIT Fuel Surcharge
    ('6506987f-925d-473d-9872-b94e0279c1af', '81e29d0c-02a6-4a7a-be02-554deb3ee49e' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL),
    --IDSFSC International Destination SIT Fuel Surcharge
    ('ca34445f-3e42-4e7e-b631-b4ae81c813c9', '690a5fc1-0ea5-4554-8294-a367b5daefa9' ,'HHG_OUTOF_NTS', 'i', false, now(), now(), NULL);
