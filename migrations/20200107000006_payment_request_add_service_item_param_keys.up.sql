INSERT INTO service_item_param_keys
(id,key,description,type,origin,created_at,updated_at)
VALUES
('164050e3-e35b-480d-bf6e-ed2fab86f370','RequestedPickupDate', 'Customer requested pick up date', 'DATE', 'PRIME', now(), now()),
('8c0c57e2-0acc-4319-9757-59608936a154','NumberDaysSIT', 'Customer requested pick up date', 'STRING', 'PRIME', now(), now()),
('02e6ab7a-818e-4a90-92a3-709019ba36a1','MarketOrigin', 'Origin Market', 'STRING', 'PRIME', now(), now()),
('6a6cb0e2-d9e4-4155-b57c-0af0ea87dc06','MarketDest', 'Dest Market', 'STRING', 'PRIME', now(), now()),
('1df2468b-8fb4-4371-b8ce-3c05d7da2050','CanStandAlone', 'Can Stand Alone', 'STRING', 'PRIME', now(), now()),
('b9739817-6408-4829-8719-1e26f8a9ceb3','BilledActualWeight', 'Billed Actual Weight', 'INTEGER', 'SYSTEM', now(), now()),
('1fe986ae-dbff-4fe1-b528-714560f7d2f5','BilledCubicFeet', 'Billed Cubic Feet', 'INTEGER', 'SYSTEM', now(), now()),
('7a99efc3-df2b-401f-ae56-f293517afbde','DistanceZip3', 'Distance using Zip3', 'INTEGER', 'SYSTEM', now(), now()),
('60b0d960-eb2e-4597-846b-d97720493799','DistanceZip5', 'Distance using Zip5', 'INTEGER', 'SYSTEM', now(), now()),
('f9753611-4b3e-4bf5-8e00-6d9ce9900f50','DistanceZip5SITOrigin', 'Distance from pickup address to SIT facility ZIP', 'INTEGER', 'SYSTEM', now(), now()),
('45ede48b-364d-473a-8c61-0f520a6a4e04','DistanceZip5SITDest', 'Distance from  SIT facility ZIP to destination address', 'INTEGER', 'SYSTEM', now(), now()),
('adeb57e5-6b1c-4c0f-b5c9-9e57e600303f','EIAFuelPrice', 'EIA diesel fuel price', 'DECIMAL', 'SYSTEM', now(), now()),
('599bbc21-8d1d-4039-9a89-ff52e3582144','ServiceAreaOrigin', 'Origin Service Area', 'STRING', 'SYSTEM', now(), now()),
('af92f0ca-f669-4483-95d2-d66e9c0c69e4','ServiceAreaDest', 'Destination Service Area', 'STRING', 'SYSTEM', now(), now()),
('edeb108a-3aa8-4e7c-9571-de81951cbb51','SITScheduleOrigin', 'Origin SIT Schedule', 'STRING', 'SYSTEM', now(), now()),
('5a3cf5bb-2bbf-4af6-8966-143d52d8f94b','SITScheduleDest', 'Dest SIT Schedule', 'STRING', 'SYSTEM', now(), now()),
('cd37b2a6-ac7d-4c93-a148-ca67f7f67cff','ServicesScheduleOrigin', 'Origin Services Schedule', 'STRING', 'SYSTEM', now(), now()),
('33751a6a-0a20-4ee2-9e46-ece2c1e997ba','ServicesScheduleDest', 'Dest Services Schedule', 'STRING', 'SYSTEM', now(), now()),
('86f9b2b6-6315-409f-a563-e91998a5043f','PriceAreaOrigin', 'Origin Price Area', 'STRING', 'SYSTEM', now(), now()),
('07e3ff12-0717-4661-a4e1-dca325da0fcc','PriceAreaDest', 'Dest Price Area', 'STRING', 'SYSTEM', now(), now()),
('6d44624c-b91b-4226-8fcd-98046e2f433d','PriceAreaIntlOrigin', 'Intl Origin Price Area', 'STRING', 'SYSTEM', now(), now()),
('4736f489-dfda-4df1-a303-8c434a120d5d','PriceAreaIntlDest', 'Intl Dest Price Area', 'STRING', 'SYSTEM', now(), now()),
('0f564bdd-153c-44af-97ae-117b62339c02','RateAreaNonStdOrigin', 'Non Std Rate Area', 'STRING', 'SYSTEM', now(), now()),
('73060cd4-15fd-4d6d-9d38-9f362c7da7af','RateAreaNonStdDest', 'Non Std Rate Area', 'STRING', 'SYSTEM', now(), now()),
('b051e180-4435-4dc1-b3b1-dbde8f1ba8ae','LinehaulDom', 'Domestic Linehaul', 'DECIMAL', 'SYSTEM', now(), now()),
('30eb55c4-f8e2-4b83-8e40-3e734c292841','LinehaulShort', 'Short Linehaul', 'DECIMAL', 'SYSTEM', now(), now()),
('4f029437-2ba3-4d5d-b25f-c43c6aeb7077','PriceDomOrigin', 'Domestic Origin Price', 'DECIMAL', 'SYSTEM', now(), now()),
('7ae62314-55f1-4fef-8b3d-c83edda9a4fd','PriceDomDest', 'Domestic Destination Price', 'DECIMAL', 'SYSTEM', now(), now()),
('2e913e12-ac88-4425-bc89-a6bccb12734b','ShippingLinehaulIntlCO', 'Intl CO Shipping & Linehaul', 'DECIMAL', 'SYSTEM', now(), now()),
('d5900ed8-f5e6-49a1-af4c-6958c1190500','ShippingLinehaulIntlOC', 'Intl OC Shipping & Linehaul', 'DECIMAL', 'SYSTEM', now(), now()),
('421a4330-accc-44f4-9418-86fd5fc6195a','ShippingLinehaulIntlOO', 'Intl OO Shipping & Linehaul', 'DECIMAL', 'SYSTEM', now(), now()),
('ce6cd6db-d2af-4f0a-a17b-63bfb12f47a0','PackingDom', 'Domestic Packing', 'DECIMAL', 'SYSTEM', now(), now()),
('5ec93806-b0a2-48e2-b7c7-fd1fce4eead0','PackingHHGIntl', 'Intl HHG Packing', 'DECIMAL', 'SYSTEM', now(), now());


INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
-- Shipment Mgmt Services
('10b87be8-d4ed-408c-925c-8e6377f914c7',(SELECT id FROM re_services WHERE code='MS'), (SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),

-- Counseling Services
('4df01039-854d-4c7e-8655-39fc3dec3237',(SELECT id FROM re_services WHERE code='CS'), (SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),

-- Dom Linehaul
('01f34b3c-7491-49da-862d-c8e2de8d64c9',(SELECT id FROM re_services WHERE code='DLH'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('48f59a0b-3afe-499c-a0c0-917d7691161f',(SELECT id FROM re_services WHERE code='DLH'),(SELECT id FROM service_item_param_keys where key='DistanceZip3'), now(), now()),
('3dab2ac4-3ef1-423e-9500-6c778230de2a',(SELECT id FROM re_services WHERE code='DLH'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('512825bd-98f6-4b1c-aa1a-d1f5b2405934',(SELECT id FROM re_services WHERE code='DLH'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),

-- Dom Shorthaul
('452f9b29-e17c-499e-aa32-83eea5757ed4',(SELECT id FROM re_services WHERE code='DSH'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('f5067fa1-5f7b-4a8f-927b-386774cbc7dd',(SELECT id FROM re_services WHERE code='DSH'),(SELECT id FROM service_item_param_keys where key='DistanceZip5'), now(), now()),
('d40c57d1-9777-4c53-ae50-5af782ec9e0e',(SELECT id FROM re_services WHERE code='DSH'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('75a6a5a7-33bb-48c1-bf46-5e485a18ef69',(SELECT id FROM re_services WHERE code='DSH'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),

-- Dom Origin Price
('105d9f5a-003e-4f33-a3e8-323d60e3aea8',(SELECT id FROM re_services WHERE code='DOP'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('e779e7aa-791a-4775-8088-2e11e518d801',(SELECT id FROM re_services WHERE code='DOP'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('35891772-d117-4340-aac3-26048516cfa3',(SELECT id FROM re_services WHERE code='DOP'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),

-- Dom Destination Price
('b6dfdb0d-83aa-4f24-bbb4-fb83af8148a3',(SELECT id FROM re_services WHERE code='DDP'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('1d67b4e2-bbbb-4aa0-9e24-9c9fe7287ccc',(SELECT id FROM re_services WHERE code='DDP'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('7b5bb8a7-62f6-47f4-9e1a-ef3d1dcbff90',(SELECT id FROM re_services WHERE code='DDP'),(SELECT id FROM service_item_param_keys where key='ServiceAreaDest'), now(), now()),

-- Dom Origin 1st Day SIT
('4207774c-166e-490f-ae8a-69eb29f92ca8',(SELECT id FROM re_services WHERE code='DOFSIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('4f6c0dac-6865-4c0a-99e3-4d2da49860e2',(SELECT id FROM re_services WHERE code='DOFSIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('421e6c66-8961-4372-abdf-4137f75488cf',(SELECT id FROM re_services WHERE code='DOFSIT'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),

-- Dom. Destination 1st Day SIT
('0439e51f-f99f-4914-824a-c1c19513436a',(SELECT id FROM re_services WHERE code='DDFSIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('7c73ec45-a00f-4d4c-b9d5-adb49995252c',(SELECT id FROM re_services WHERE code='DDFSIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('092e23a1-6842-4967-806b-6f1b899349e2',(SELECT id FROM re_services WHERE code='DDFSIT'),(SELECT id FROM service_item_param_keys where key='ServiceAreaDest'), now(), now()),

-- Dom. Origin Add'l SIT
('4e8e84f5-d254-4237-8a47-fa0cda5f7a5c',(SELECT id FROM re_services WHERE code='DOASIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('5f51ba93-d0e8-4ae0-8195-b937ee4bde6a',(SELECT id FROM re_services WHERE code='DOASIT'),(SELECT id FROM service_item_param_keys where key='NumberDaysSIT'), now(), now()),
('678b8488-47c3-4f05-a9a0-6155fd062c8a',(SELECT id FROM re_services WHERE code='DOASIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('4df8aa63-d87b-4579-9603-20b522d27371',(SELECT id FROM re_services WHERE code='DOASIT'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),

-- Destination Add'l SIT
('65cfc862-051c-4fd3-bea5-092da05a86c7',(SELECT id FROM re_services WHERE code='DDASIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('82e8672c-23e4-4aea-88f4-7ba40d7f9efc',(SELECT id FROM re_services WHERE code='DDASIT'),(SELECT id FROM service_item_param_keys where key='NumberDaysSIT'), now(), now()),
('9de41f0c-19d1-4e3a-b626-e90a9646eaf3',(SELECT id FROM re_services WHERE code='DDASIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('a934d40b-4fff-4c53-91f7-ace45ffee655',(SELECT id FROM re_services WHERE code='DDASIT'),(SELECT id FROM service_item_param_keys where key='ServiceAreaDest'), now(), now()),

-- Dom. Origin SIT Pickup
('dc9a0e5a-a5b8-4b45-88ba-44284af4d064',(SELECT id FROM re_services WHERE code='DOPSIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('694988ab-93ad-46f4-b654-343eda690a9e',(SELECT id FROM re_services WHERE code='DOPSIT'),(SELECT id FROM service_item_param_keys where key='DistanceZip5SITOrigin'), now(), now()),
('ec2dc91f-5b8e-47d4-bac6-b219cb8fb9bc',(SELECT id FROM re_services WHERE code='DOPSIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('782d4060-7528-4b46-8618-775ce184ee2f',(SELECT id FROM re_services WHERE code='DOPSIT'),(SELECT id FROM service_item_param_keys where key='SITScheduleOrigin'), now(), now()),

-- Dom. Destination SIT Delivery
('70a2d389-9c72-4d44-8348-bfe3ddd182b1',(SELECT id FROM re_services WHERE code='DDDSIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('36876790-1efb-4c51-bd74-ee63e0eb72ee',(SELECT id FROM re_services WHERE code='DDDSIT'),(SELECT id FROM service_item_param_keys where key='DistanceZip5SITDest'), now(), now()),
('cbb1dcac-71dc-41b1-a285-29310003a2d4',(SELECT id FROM re_services WHERE code='DDDSIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('f8e9e01d-cd1d-4293-85fc-c00dff309671',(SELECT id FROM re_services WHERE code='DDDSIT'),(SELECT id FROM service_item_param_keys where key='SITScheduleDest'), now(), now()),

-- Service Item	Dom. Packing
('47b81436-f1df-4c61-810c-939467d02826',(SELECT id FROM re_services WHERE code='DPK'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('0a21ebb3-4de1-4527-9dd7-594100477dd4',(SELECT id FROM re_services WHERE code='DPK'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('7dcc822c-7a5c-4622-a95d-5d662bbd439f',(SELECT id FROM re_services WHERE code='DPK'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleOrigin'), now(), now()),

-- Service Item	Dom. Unpacking
('6a40fd50-f0fd-4067-a9a7-d5258983db31',(SELECT id FROM re_services WHERE code='DUPK'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('2ba6276a-2c72-47f9-8550-aee7cb7e1c18',(SELECT id FROM re_services WHERE code='DUPK'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('fea8e93a-7f4a-4b51-88ca-f57de797333f',(SELECT id FROM re_services WHERE code='DUPK'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleDest'), now(), now()),

-- Service Item	Dom. Crating
('769f15bb-8b22-403a-a444-975101855013',(SELECT id FROM re_services WHERE code='DCRT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('246d5628-16ad-4b37-89e2-db20963f7db1',(SELECT id FROM re_services WHERE code='DCRT'),(SELECT id FROM service_item_param_keys where key='CanStandAlone'), now(), now()),
('336f790e-ec2c-4587-99b4-6b75d278ecb9',(SELECT id FROM re_services WHERE code='DCRT'),(SELECT id FROM service_item_param_keys where key='BilledCubicFeet'), now(), now()),
('a4145ea0-be9d-4b4c-8cd9-bacde58dfc4c',(SELECT id FROM re_services WHERE code='DCRT'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleOrigin'), now(), now()),

-- Service Item	Dom. Uncrating
('a5a7b555-f135-470e-b94b-251f966b5a57',(SELECT id FROM re_services WHERE code='DUCRT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('15ddfdba-9907-4f79-ba5f-87af493f5e3a',(SELECT id FROM re_services WHERE code='DUCRT'),(SELECT id FROM service_item_param_keys where key='BilledCubicFeet'), now(), now()),
('5190e621-fe26-4e10-9ff6-faf8488b0157',(SELECT id FROM re_services WHERE code='DUCRT'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleDest'), now(), now()),

-- Dom. Origin Shuttle Service
('d1dc3d9d-e0ce-4409-a783-bc87702abbd0',(SELECT id FROM re_services WHERE code='DOSHUT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('26de29be-736a-4746-b9c4-0831a4a29850',(SELECT id FROM re_services WHERE code='DOSHUT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('61476f61-1ab5-4336-b38d-010e6502c4c4',(SELECT id FROM re_services WHERE code='DOSHUT'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleOrigin'), now(), now()),

-- Dom. Destination Shuttle Service
('1a340b6e-1d72-4967-ad3f-3e1097f22ac5',(SELECT id FROM re_services WHERE code='DDSHUT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('76cf0d69-ded1-4988-98d9-06efb4d731e8',(SELECT id FROM re_services WHERE code='DDSHUT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('603934a7-1fcd-455a-ba04-12766b8c607d',(SELECT id FROM re_services WHERE code='DDSHUT'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleDest'), now(), now()),

-- Int'l. O->O Shipping & LH
('c298d508-b922-49c3-be29-a344a6cb133a',(SELECT id FROM re_services WHERE code='IOOLH'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('06af387d-7eae-46f8-b038-9865e6989558',(SELECT id FROM re_services WHERE code='IOOLH'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('a1c4ab39-e918-45c7-a0fe-b31e65a745ac',(SELECT id FROM re_services WHERE code='IOOLH'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlOrigin'), now(), now()),
('ed73e791-f997-4267-97a7-769a8fc56e40',(SELECT id FROM re_services WHERE code='IOOLH'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlDest'), now(), now()),

-- Service Item	Int'l. O->O UB
('44954698-2a70-4e05-8ad9-1308bafb1b6c',(SELECT id FROM re_services WHERE code='IOOUB'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('4931b3ae-e4a6-4b06-a292-3c5761d36f0f',(SELECT id FROM re_services WHERE code='IOOUB'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('0110646d-de0c-46ba-b810-d5b617340be0',(SELECT id FROM re_services WHERE code='IOOUB'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlOrigin'), now(), now()),
('51403057-2c10-4704-ac98-36b47241574c',(SELECT id FROM re_services WHERE code='IOOUB'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlDest'), now(), now()),

-- Int'l. C->O Shipping & LH
('b0c1a2da-40dc-4ffc-8daf-639512a86832',(SELECT id FROM re_services WHERE code='ICOLH'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('943b00bd-05fe-480e-ab7e-b44075b3584e',(SELECT id FROM re_services WHERE code='ICOLH'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('985f5bf2-48d2-4bb1-afcb-0ecb8d379dbe',(SELECT id FROM re_services WHERE code='ICOLH'),(SELECT id FROM service_item_param_keys where key='PriceAreaOrigin'), now(), now()),
('2f32e977-eb18-48d3-82ab-f80b79920b63',(SELECT id FROM re_services WHERE code='ICOLH'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlDest'), now(), now()),

-- Int'l. C->O UB
('6c923165-d8f3-495f-a962-fc9aeb5e0b15',(SELECT id FROM re_services WHERE code='ICOUB'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('a7f25d48-0b4a-4acd-bbd0-cba600aa86af',(SELECT id FROM re_services WHERE code='ICOUB'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('3f997e31-6949-4e49-aebc-1307b1f911ae',(SELECT id FROM re_services WHERE code='ICOUB'),(SELECT id FROM service_item_param_keys where key='PriceAreaOrigin'), now(), now()),
('c1d576bf-6846-4cab-9ed2-292fe6b31517',(SELECT id FROM re_services WHERE code='ICOUB'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlDest'), now(), now()),

-- Int'l. O->C Shipping & LH
('0375a7e4-0b45-409d-9502-4cc35a630d58',(SELECT id FROM re_services WHERE code='IOCLH'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('ca0b1132-a862-435a-be49-ae183c744033',(SELECT id FROM re_services WHERE code='IOCLH'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('ad5a77f6-8b54-4579-a408-197158a76e6f',(SELECT id FROM re_services WHERE code='IOCLH'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlOrigin'), now(), now()),
('27df3253-0986-488f-8636-b98a8d8b326c',(SELECT id FROM re_services WHERE code='IOCLH'),(SELECT id FROM service_item_param_keys where key='PriceAreaDest'), now(), now()),

-- Int'l. O->C UB
('6b1d22b3-3511-4a06-9462-6e1a42a16526',(SELECT id FROM re_services WHERE code='IOCUB'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('cdbc5b22-69bc-4474-8270-ba12255b2df2',(SELECT id FROM re_services WHERE code='IOCUB'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('ce7c56c9-da99-4488-bf7e-830d39d6c375',(SELECT id FROM re_services WHERE code='IOCUB'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlOrigin'), now(), now()),
('4ad7b8b4-424f-4198-91d1-bbc40e71a4e3',(SELECT id FROM re_services WHERE code='IOCUB'),(SELECT id FROM service_item_param_keys where key='PriceAreaDest'), now(), now()),

-- Int'l. HHG Pack
('fdb72da3-6297-4748-b22c-c4ceb546525e',(SELECT id FROM re_services WHERE code='IHPK'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('79b0b82c-9e78-4514-810b-b91bcc339b8a',(SELECT id FROM re_services WHERE code='IHPK'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('ebfdcfb1-e612-483e-87d4-36f12d336ea9',(SELECT id FROM re_services WHERE code='IHPK'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlOrigin'), now(), now()),

-- Item	Int'l. HHG Unpack
('3fa3a40e-6bd3-4bb8-8d83-06c93bd23b63',(SELECT id FROM re_services WHERE code='IHUPK'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('6b6e9766-1292-44de-b434-5cc1b80d80de',(SELECT id FROM re_services WHERE code='IHUPK'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('d1df9449-e6a9-4a4a-8117-8daf70c4475b',(SELECT id FROM re_services WHERE code='IHUPK'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlDest'), now(), now()),

-- Item	Int'l. UB Pack
('45f304ff-2597-416b-8aa6-120265f56b90',(SELECT id FROM re_services WHERE code='IUBPK'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('7308af31-df0f-4584-834b-bf7480aedb9f',(SELECT id FROM re_services WHERE code='IUBPK'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('a8d577ba-f277-4f7b-9170-b36d67c7b738',(SELECT id FROM re_services WHERE code='IUBPK'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlOrigin'), now(), now()),

-- Int'l. UB Unpack
('a3df75f7-216f-4083-b638-52a5e9400a0a',(SELECT id FROM re_services WHERE code='IUBUPK'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('4c9bb859-9f31-4b70-8d30-21ddf8b964a7',(SELECT id FROM re_services WHERE code='IUBUPK'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('88323d4b-46b5-46dd-840c-6039cb48f93b',(SELECT id FROM re_services WHERE code='IUBUPK'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlDest'), now(), now()),

-- Int'l. Origin 1st Day SIT
('72c52a8c-1b9d-4097-b061-7cb0ef185eb5',(SELECT id FROM re_services WHERE code='IOFSIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('01da041c-82b1-4893-aba1-12eaa0e1f25e',(SELECT id FROM re_services WHERE code='IOFSIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('7dcd21ac-e327-4183-8b15-a4d770807a9c',(SELECT id FROM re_services WHERE code='IOFSIT'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlOrigin'), now(), now()),

-- Int'l. Destination 1st Day SIT
('b40286a9-afc0-45a1-b3bf-59a2f53284d2',(SELECT id FROM re_services WHERE code='IDFSIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('62412254-7bc5-4d99-8026-430fcc03743d',(SELECT id FROM re_services WHERE code='IDFSIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('23483138-1b7f-4ace-9b21-d07d1fecb734',(SELECT id FROM re_services WHERE code='IDFSIT'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlDest'), now(), now()),

-- Int'l. Origin Add'l Day SIT
('8678b18e-80de-4e17-b426-261aad3575c3',(SELECT id FROM re_services WHERE code='IOASIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('30dfd247-3f02-49c3-ab3e-87e381256655',(SELECT id FROM re_services WHERE code='IOASIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('55d73836-c567-4abe-9d38-77d4eb7f68d7',(SELECT id FROM re_services WHERE code='IOASIT'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlOrigin'), now(), now()),

-- Int'l. Destination Add'l Day SIT
('8adcf10e-8be7-4fd8-88e8-b59284c29821',(SELECT id FROM re_services WHERE code='IDASIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('579d7e93-088d-4bf3-8c5e-f3763dc2e29c',(SELECT id FROM re_services WHERE code='IDASIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('081a809d-3aaf-4818-a546-d6e1dbecc499',(SELECT id FROM re_services WHERE code='IDASIT'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlDest'), now(), now()),

-- Int'l. Origin SIT Pickup
('3a2b5fe8-4fe1-436d-b03f-f8a5e773dc74',(SELECT id FROM re_services WHERE code='IOPSIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('6fbff63f-f63e-4de1-8429-28dea833171c',(SELECT id FROM re_services WHERE code='IOPSIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('459c72da-abb7-435f-84d8-e46262caba49',(SELECT id FROM re_services WHERE code='IOPSIT'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlOrigin'), now(), now()),

-- Int'l. Destination SIT Delivery
('f2800b0f-6bf3-425c-9b91-a7a59b39abc2',(SELECT id FROM re_services WHERE code='IDDSIT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('e23e5b16-6759-4957-aa2c-4138c928ddcd',(SELECT id FROM re_services WHERE code='IDDSIT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('af5a41d2-2711-4265-8752-f46e3e17a936',(SELECT id FROM re_services WHERE code='IDDSIT'),(SELECT id FROM service_item_param_keys where key='PriceAreaIntlDest'), now(), now()),

-- Int'l. Crating
('62d41787-716a-4302-b7a6-4afd344a0079',(SELECT id FROM re_services WHERE code='ICRT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('61e7fd20-70fa-4464-83fb-7ab995178f26',(SELECT id FROM re_services WHERE code='ICRT'),(SELECT id FROM service_item_param_keys where key='BilledCubicFeet'), now(), now()),
('d755c944-7380-4ca3-a01e-b273edf2a7cc',(SELECT id FROM re_services WHERE code='ICRT'),(SELECT id FROM service_item_param_keys where key='MarketOrigin'), now(), now()),

-- Int'l. Uncrating
('18f7be16-1fc3-449e-a240-586c2c17f707',(SELECT id FROM re_services WHERE code='IUCRT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('9a98ceef-5629-435e-8214-8c30511467e1',(SELECT id FROM re_services WHERE code='IUCRT'),(SELECT id FROM service_item_param_keys where key='BilledCubicFeet'), now(), now()),
('b34b9c68-21b0-4b01-a553-78d4cbf0a516',(SELECT id FROM re_services WHERE code='IUCRT'),(SELECT id FROM service_item_param_keys where key='MarketDest'), now(), now()),

-- Int'l. Origin Shuttle Service
('100b3224-de5e-4242-b160-fc71d59a19de',(SELECT id FROM re_services WHERE code='IOSHUT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('9c3ce9ca-df95-4a0a-bf40-8ad5ee54f692',(SELECT id FROM re_services WHERE code='IOSHUT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('3876a014-0e63-45cf-82dc-bc2c8e92d1b0',(SELECT id FROM re_services WHERE code='IOSHUT'),(SELECT id FROM service_item_param_keys where key='MarketOrigin'), now(), now()),

-- Int'l. Destination Shuttle Service
('eb253cd4-b27b-4393-8ecc-e29a0d3f4f8f',(SELECT id FROM re_services WHERE code='IDSHUT'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('d45f0292-d252-4eff-91fa-f77680ac3d6f',(SELECT id FROM re_services WHERE code='IDSHUT'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('54de2e15-f0c7-490e-9770-22d2660eabc1',(SELECT id FROM re_services WHERE code='IDSHUT'),(SELECT id FROM service_item_param_keys where key='MarketDest'), now(), now()),

-- NonStd. HHG
('83caa7e3-625d-4722-b3d8-81bf90d5b53a',(SELECT id FROM re_services WHERE code='NSTH'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('3d6fa54f-8f74-4fcf-89a7-77e495aff7be',(SELECT id FROM re_services WHERE code='NSTH'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('0e7ae010-ebb1-41ad-8138-8d2dd2e780d0',(SELECT id FROM re_services WHERE code='NSTH'),(SELECT id FROM service_item_param_keys where key='RateAreaNonStdOrigin'), now(), now()),
('1d24be3a-0540-49dd-bd49-23d0ede9ea23',(SELECT id FROM re_services WHERE code='NSTH'),(SELECT id FROM service_item_param_keys where key='RateAreaNonStdDest'), now(), now()),

-- NonStd. UB
('f07868ad-c7f9-491d-9709-efc5a5f9acce',(SELECT id FROM re_services WHERE code='NSTUB'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('7ab04a2a-7a25-4ce2-b575-aad16ea3145a',(SELECT id FROM re_services WHERE code='NSTUB'),(SELECT id FROM service_item_param_keys where key='BilledActualWeight'), now(), now()),
('5fa25515-ff12-4043-a837-254926cf5b6d',(SELECT id FROM re_services WHERE code='NSTUB'),(SELECT id FROM service_item_param_keys where key='RateAreaNonStdOrigin'), now(), now()),
('99cbf769-24d1-469b-ac29-b3c4c4e417c9',(SELECT id FROM re_services WHERE code='NSTUB'),(SELECT id FROM service_item_param_keys where key='RateAreaNonStdDest'), now(), now()),

-- Fuel Surcharge
('f84159a2-62b9-468b-a1cd-4014f8fb0075',(SELECT id FROM re_services WHERE code='FSC'),(SELECT id FROM service_item_param_keys where key='LinehaulDom'), now(), now()),
('1d861c97-05d4-446f-8d2c-6b0502a74f53',(SELECT id FROM re_services WHERE code='FSC'),(SELECT id FROM service_item_param_keys where key='LinehaulShort'), now(), now()),
('64d4e3ae-c69b-4c43-9a64-2bff0f8b29ea',(SELECT id FROM re_services WHERE code='FSC'),(SELECT id FROM service_item_param_keys where key='EIAFuelPrice'), now(), now()),

-- Dom. Mobile Home Factor
('cc52706b-5f40-4a6a-9650-723035a722d9',(SELECT id FROM re_services WHERE code='DMHF'),(SELECT id FROM service_item_param_keys where key='LinehaulDom'), now(), now()),
('b654fa42-5fd3-4124-961a-49e817651533',(SELECT id FROM re_services WHERE code='DMHF'),(SELECT id FROM service_item_param_keys where key='LinehaulShort'), now(), now()),
('78837c64-1953-49c2-9446-dfc92a30c028',(SELECT id FROM re_services WHERE code='DMHF'),(SELECT id FROM service_item_param_keys where key='PriceDomOrigin'), now(), now()),
('f2ad93c6-11e4-4b31-b026-15b2878a617c',(SELECT id FROM re_services WHERE code='DMHF'),(SELECT id FROM service_item_param_keys where key='PriceDomDest'), now(), now()),

-- Dom. Tow Away Boat Factor
('5d750e0c-9907-44e5-84a5-d0e9e670fe53',(SELECT id FROM re_services WHERE code='DBTF'),(SELECT id FROM service_item_param_keys where key='LinehaulDom'), now(), now()),
('ee123a61-7ac8-40a6-9300-92ca13271667',(SELECT id FROM re_services WHERE code='DBTF'),(SELECT id FROM service_item_param_keys where key='LinehaulShort'), now(), now()),
('ab07b8b4-22d5-455f-9c43-ad264890ea73',(SELECT id FROM re_services WHERE code='DBTF'),(SELECT id FROM service_item_param_keys where key='PriceDomOrigin'), now(), now()),
('8d4df03a-88a5-4f44-b776-26056a4e5265',(SELECT id FROM re_services WHERE code='DBTF'),(SELECT id FROM service_item_param_keys where key='PriceDomDest'), now(), now()),

-- Dom. Haul Away Boat Factor
('62e7091d-3c0d-4e5c-ab38-cceeeac10393',(SELECT id FROM re_services WHERE code='DBHF'),(SELECT id FROM service_item_param_keys where key='LinehaulDom'), now(), now()),
('8dc69dda-e979-487f-ac5f-52feadc2ba00',(SELECT id FROM re_services WHERE code='DBHF'),(SELECT id FROM service_item_param_keys where key='LinehaulShort'), now(), now()),
('b7b8fb51-ce90-4e3f-ba5b-810e8448a911',(SELECT id FROM re_services WHERE code='DBHF'),(SELECT id FROM service_item_param_keys where key='PriceDomOrigin'), now(), now()),
('a4fc9c74-bb72-4351-8920-ec228aa1e3a9',(SELECT id FROM re_services WHERE code='DBHF'),(SELECT id FROM service_item_param_keys where key='PriceDomDest'), now(), now()),

-- Int’l. Tow Away Boat Factor
('b19959da-375c-4b71-b088-37968bc7e4e7',(SELECT id FROM re_services WHERE code='IBTF'),(SELECT id FROM service_item_param_keys where key='ShippingLinehaulIntlOO'), now(), now()),
('1600c27b-e7c5-4548-a926-f87ba863e9ff',(SELECT id FROM re_services WHERE code='IBTF'),(SELECT id FROM service_item_param_keys where key='ShippingLinehaulIntlCO'), now(), now()),
('746ca650-f544-41e3-a599-fb02eb5cb8d7',(SELECT id FROM re_services WHERE code='IBTF'),(SELECT id FROM service_item_param_keys where key='ShippingLinehaulIntlOC'), now(), now()),

-- Int’l. Haul Away Boat Factor
('3f9364df-eebe-47ba-bf53-b5a7973e159f',(SELECT id FROM re_services WHERE code='IBHF'),(SELECT id FROM service_item_param_keys where key='ShippingLinehaulIntlOO'), now(), now()),
('1759cc7f-1d00-4b27-ac3e-52b123842e1c',(SELECT id FROM re_services WHERE code='IBHF'),(SELECT id FROM service_item_param_keys where key='ShippingLinehaulIntlCO'), now(), now()),
('5b920850-16b2-4af6-8553-1fa897bdd013',(SELECT id FROM re_services WHERE code='IBHF'),(SELECT id FROM service_item_param_keys where key='ShippingLinehaulIntlOC'), now(), now()),

-- Dom. NTS Packing Factor
('b82e1e45-8286-4247-b1c0-7f18d15a5896',(SELECT id FROM re_services WHERE code='DNPKF'),(SELECT id FROM service_item_param_keys where key='PackingDom'), now(), now()),

-- Int’l. NTS Packing Factor
('1105dad9-e167-4fea-a84d-fc5db52802bd',(SELECT id FROM re_services WHERE code='INPKF'),(SELECT id FROM service_item_param_keys where key='PackingHHGIntl'), now(), now());
