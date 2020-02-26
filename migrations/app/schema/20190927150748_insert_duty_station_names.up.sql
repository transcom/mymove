-- air force
INSERT INTO duty_station_names VALUES ('5c04e8a5-4379-4f7a-9138-1c2ce1e53eb4', 'JBER', (SELECT id FROM duty_stations WHERE name = 'JB Elmendorf-Richardson'), now(), now());
INSERT INTO duty_station_names VALUES ('e7cc7468-9e4d-439e-bd91-019597c332aa', 'JBLE', (SELECT id FROM duty_stations WHERE name = 'JB Langley-Eustis'), now(), now());
INSERT INTO duty_station_names VALUES ('f23b6966-8664-4031-9d17-c2416ed2c063', 'JBLM', (SELECT id FROM duty_stations WHERE name = 'JB Lewis-McChord'), now(), now());
INSERT INTO duty_station_names VALUES ('38f1478b-6e29-4912-94bc-24b8e4c938e0', 'JBMDL', (SELECT id FROM duty_stations WHERE name = 'JB McGuire-Dix-Lakehurst'), now(), now());
INSERT INTO duty_station_names VALUES ('86160638-004a-44fa-9ef3-dcd688244567', 'LAAFB', (SELECT id FROM duty_stations WHERE name = 'Los Angeles AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('49b03d29-1510-47cd-8064-ab546156fd4a', 'Las Vegas Army Airfield', (SELECT id FROM duty_stations WHERE name = 'Nellis AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('fd6b122b-89fc-4e64-9e60-41e1f7e3be38', 'WPAFB', (SELECT id FROM duty_stations WHERE name = 'Wright-Patterson AFB'), now(), now());

-- army
INSERT INTO duty_station_names VALUES ('6c87ba18-f8e8-47d9-8270-4a8e769276d7', 'APG', (SELECT id FROM duty_stations WHERE name = 'Aberdeen Proving Ground'), now(), now());
INSERT INTO duty_station_names VALUES ('d7f7f420-9b68-4cf6-86c5-9bffb6079b84', 'BGAD', (SELECT id FROM duty_stations WHERE name = 'Blue Grass Army Depot'), now(), now());
INSERT INTO duty_station_names VALUES ('910c68c4-d743-4b41-ab08-d62630e65d83', 'DPG', (SELECT id FROM duty_stations WHERE name = 'Dugway Proving Ground'), now(), now());
INSERT INTO duty_station_names VALUES ('3baf7673-8acf-4ddb-8a00-e0a760289ff8', 'National Training Center', (SELECT id FROM duty_stations WHERE name = 'Fort Irwin'), now(), now());
INSERT INTO duty_station_names VALUES ('f3e5ffa1-669c-45e9-9e69-db302a5979f4', 'NTC', (SELECT id FROM duty_stations WHERE name = 'Fort Irwin'), now(), now());
INSERT INTO duty_station_names VALUES ('1ef11c12-7129-4832-b70a-cbb0346a3806', 'United States Army Garrison Alaska', (SELECT id FROM duty_stations WHERE name = 'Fort Wainwright'), now(), now());
INSERT INTO duty_station_names VALUES ('d2740c2f-ad44-47b4-abfc-76880de1aaa5', 'USARAK', (SELECT id FROM duty_stations WHERE name = 'Fort Wainwright'), now(), now());
INSERT INTO duty_station_names VALUES ('744e735c-acd2-431b-a88c-2720d6d5a67c', 'PBA', (SELECT id FROM duty_stations WHERE name = 'Pine Bluff Arsenal'), now(), now());
INSERT INTO duty_station_names VALUES ('5f6541ed-8550-409a-a48a-2b009453f599', 'RRAD', (SELECT id FROM duty_stations WHERE name = 'Red River Army Depot'), now(), now());
INSERT INTO duty_station_names VALUES ('293971e9-8f38-4a3e-9dbc-78c979d15882', 'RSA', (SELECT id FROM duty_stations WHERE name = 'Rock Island Arsenal'), now(), now());
INSERT INTO duty_station_names VALUES ('0c3f4676-8f1c-4057-8a61-5e1e638059fa', 'TEAD', (SELECT id FROM duty_stations WHERE name = 'Tooele Army Depot'), now(), now());
INSERT INTO duty_station_names VALUES ('70f4d237-5393-413f-b23e-e1f85102342a', 'WSMR', (SELECT id FROM duty_stations WHERE name = 'White Sands Missile Range'), now(), now());

-- coast guard
INSERT INTO duty_station_names VALUES ('ef1b5238-51d0-4643-9c93-70142ddd2a12', 'Baltimore', (SELECT id FROM duty_stations WHERE name = 'Coast Guard Yard'), now(), now());
INSERT INTO duty_station_names VALUES ('058dff3d-d42e-4448-9a49-e13410a9000e', 'Buzzards Bay', (SELECT id FROM duty_stations WHERE name = 'Base Cape Cod'), now(), now());
INSERT INTO duty_station_names VALUES ('b9448609-6893-46f4-b658-ef431059a05d', 'Staten Island', (SELECT id FROM duty_stations WHERE name = 'Station New York'), now(), now());

-- Navy
INSERT INTO duty_station_names VALUES ('6a80c7ac-59e1-45c3-8335-607b03744d32', 'USNA', (SELECT id FROM duty_stations WHERE name = 'US Naval Academy'), now(), now());
INSERT INTO duty_station_names VALUES ('c466de05-e9b3-463b-8fd3-e9c2d64cdeec', 'Annapolis', (SELECT id FROM duty_stations WHERE name = 'US Naval Academy'), now(), now());
INSERT INTO duty_station_names VALUES ('e5abd627-a113-493e-b11c-56add695ad2d', 'JBAB', (SELECT id FROM duty_stations WHERE name = 'JB Anacostia–Bolling'), now(), now());
INSERT INTO duty_station_names VALUES ('00c38cae-c20f-4bcd-b83d-58dbe2716a04', '29 Palms', (SELECT id FROM duty_stations WHERE name = 'MCAGCC Twentynine Palms'), now(), now());
INSERT INTO duty_station_names VALUES ('fadfa776-534a-480a-8ba6-7d2f23dbf8a6', 'Carswell Field', (SELECT id FROM duty_stations WHERE name = 'NAS Fort Worth JRB'), now(), now());
INSERT INTO duty_station_names VALUES ('02a4ca7b-9d72-4e4e-8ff6-12c00d56f6da', 'NSAB', (SELECT id FROM duty_stations WHERE name = 'NSA Bethesda'), now(), now());
INSERT INTO duty_station_names VALUES ('f8348194-999f-4df8-a286-fbf4d59f96ea', 'NPS', (SELECT id FROM duty_stations WHERE name = 'Naval Postgraduate School'), now(), now());
INSERT INTO duty_station_names VALUES ('82fd8439-f7c5-4c96-8d85-d12d162b2d65', 'NSAMS', (SELECT id FROM duty_stations WHERE name = 'NSA Mid-South'), now(), now());
INSERT INTO duty_station_names VALUES ('6267be24-1e28-4870-81ff-a04de3c83264', 'NSAPC', (SELECT id FROM duty_stations WHERE name = 'NSA Panama City'), now(), now());
INSERT INTO duty_station_names VALUES ('21a74468-20d7-46a3-921c-99ddbfae89ff', 'NASWI', (SELECT id FROM duty_stations WHERE name = 'NAS Whidbey Island'), now(), now());
INSERT INTO duty_station_names VALUES ('b9ff06bc-6af7-4d95-9d28-778384ccd870', 'NBVC', (SELECT id FROM duty_stations WHERE name = 'NB Ventura County'), now(), now());
INSERT INTO duty_station_names VALUES ('ccd08394-9aae-4f05-a37f-a612868b738e', 'Naval Base Ventura County', (SELECT id FROM duty_stations WHERE name = 'NB Ventura County'), now(), now());
INSERT INTO duty_station_names VALUES ('bb6cdedf-b7af-4868-b089-3e9e838f6fae', 'PNS', (SELECT id FROM duty_stations WHERE name = 'Portsmouth Naval Shipyard'), now(), now());
INSERT INTO duty_station_names VALUES ('7cccb2bd-c639-4362-a7a5-7f04e67f70d3', 'Portsmouth Navy Yard', (SELECT id FROM duty_stations WHERE name = 'Portsmouth Naval Shipyard'), now(), now());

-- AFB -> Air Force Base
INSERT INTO duty_station_names VALUES ('1529D2E4-9DC0-3A16-D49F-19069AB83C69', 'Altus Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Altus AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('AE8659D5-E3F5-D042-42D8-80A946F9C5A0', 'Barksdale Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Barksdale AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('1F07A018-2FB6-3086-1BC4-29E31BAC4C5F', 'Beale Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Beale AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('BF08B82C-BF4F-3069-3FCD-BC84157C7CE8', 'Buckley Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Buckley AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('85C307D6-6A4C-9F39-5279-A8501C8DB9E3', 'Cannon Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Cannon AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('29465B61-1FC4-4D1A-B059-727D8640AED1', 'Columbus Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Columbus AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('4C16D282-3AD4-1EB7-9A7D-C1D27CA36302', 'Creech Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Creech AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('C2534719-DBA0-2C71-1D1B-3C021E578909', 'Davis-Monthan Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Davis-Monthan AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('1B27A62F-71CF-D405-27DA-36BF72F01E02', 'Dover Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Dover AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('A9FB1924-4F01-9545-FEF5-3682B6C630B8', 'Dyess Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Dyess AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('FAF3E6D9-4A86-30CE-97F8-95317AD5AD09', 'Edwards Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Edwards AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('52A654CB-B0AC-C90A-8F6E-3D174F4A2435', 'Eglin Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Eglin AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('D389FAF6-9C87-4C2D-E1D1-E8EFC7FD3410', 'Eielson Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Eielson AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('6BBE3F0B-84A8-5301-D904-E25634BE9F25', 'Ellsworth Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Ellsworth AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('D5454B14-D0EF-EF17-EF84-496FD9279049', 'Fairchild Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Fairchild AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('59630B31-6380-C9FE-A2C0-7F440AF9D816', 'F.E. Warren Air Force Base', (SELECT id FROM duty_stations WHERE name = 'F.E. Warren AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('B1C1CF7A-C8BF-DC97-3D13-529F9A431EDA', 'Goodfellow Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Goodfellow AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('F6293464-F2A9-C1C1-17E6-AFD0761842B5', 'Grand Forks Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Grand Forks AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('5E46D734-7A7C-CB2C-B427-AC2C168C0416', 'Griffiss Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Griffiss AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('2A75052F-164D-637B-236E-3A95E1D950D3', 'Hanscom Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Hanscom AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('7D4C8074-D59E-E34A-DF3D-A01EC25E75F6', 'Hill Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Hill AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('7F746CFE-F2EC-8D5C-2195-B86A842783D0', 'Holloman Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Holloman AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('2D53A546-AC05-C210-A8B3-47360F0497E3', 'Hurlburt Field Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Hurlburt Field AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('2AC5BD5E-F83B-7A1C-97F3-375868F26EB7', 'Keesler Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Keesler AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('420CAF15-6DCF-C901-1EF5-AD549E716C50', 'Kirtland Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Kirtland AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('74645DAE-2308-3EAB-CBC9-C8C1608BF27E', 'Laughlin Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Laughlin AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('5873C0F7-E3AC-B78A-41BA-7FC737F27B9A', 'Little Rock Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Little Rock AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('B9FA36CB-C3E5-B608-CD25-8B867CD21AD4', 'Los Angeles Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Los Angeles AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('F1CF786F-291D-5079-5A32-29C9EB3B097E', 'Luke Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Luke AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('5D3EA15F-4327-F70B-E32B-A46182189AD6', 'MacDill Air Force Base', (SELECT id FROM duty_stations WHERE name = 'MacDill AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('6148FC76-9CF4-EF97-EC20-F28B159EC8CD', 'Malmstrom Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Malmstrom AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('F14BFC76-62B3-C6F0-B532-86D186D6BF5C', 'Maxwell Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Maxwell AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('D1D6F7CF-CB60-248F-1959-D97DE9014530', 'McConnell Air Force Base', (SELECT id FROM duty_stations WHERE name = 'McConnell AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('1E04196E-DC05-D84C-6042-FC205362ADC3', 'Minot Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Minot AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('AF8461C7-6E92-64DF-120C-9C16AC1B42BD', 'Moody Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Moody AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('6E854F86-961D-D758-9E3D-615C25A01034', 'Mountain Home Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Mountain Home AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('D32F7B7D-3542-F261-FDEF-F3A875CA7EDE', 'Nellis Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Nellis AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('719B9D19-6DE4-B4EF-F127-9DED26201C85', 'Offutt Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Offutt AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('E49DCE78-6E14-3DC1-E5A2-8348641C9B25', 'Patrick Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Patrick AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('9B147CD9-CFB5-2CE5-6839-F093AB8C315F', 'Peterson Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Peterson AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('E6E163FD-D5C0-5910-EC71-485E0D8F4295', 'Robins Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Robins AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('14DB3B70-6325-DE12-75C3-B279BAEA680D', 'Schriever Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Schriever AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('7E4ABC42-902E-A013-1985-FBA3C710D0B9', 'Scott Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Scott AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('E57D1F5D-38A7-5BEF-3148-67D7E82A25E8', 'Seymour Johnson Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Seymour Johnson AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('80E9A502-4935-25F9-D23B-CB6828B968FE', 'Shaw Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Shaw AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('B165C8D4-529B-D43F-5840-2C6F50502735', 'Sheppard Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Sheppard AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('CB379D87-9C8F-A36C-A8C4-CD6172DC2CF0', 'Tinker Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Tinker AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('CE9A59DE-B083-CA35-F242-90BE79A69B59', 'Travis Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Travis AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('31C0D594-1A73-F04B-7452-DAE547E01E6A', 'Tyndall Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Tyndall AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('48618CFE-9B34-E3AE-1AFC-82FABEC96E23', 'Vance Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Vance AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('1536A80C-7630-428D-2936-D0F9B2C0AF14', 'Vandenberg Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Vandenberg AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('4BF1C231-C65F-2543-A86C-6F5DAD9B6867', 'Whiteman Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Whiteman AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('C41A354B-6BCA-40A7-4315-6F6E9E1C592D', 'Wright-Patterson Air Force Base', (SELECT id FROM duty_stations WHERE name = 'Wright-Patterson AFB'), now(), now());

-- JB -> Joint Base
INSERT INTO duty_station_names VALUES ('10B92CD4-9EC8-8B6C-9EC4-1AC65940BA7E', 'Joint Base Langley-Eustis', (SELECT id FROM duty_stations WHERE name = 'JB Langley-Eustis'), now(), now());
INSERT INTO duty_station_names VALUES ('9FE90C7F-2E94-4B0F-A974-6A7ACFE6E2D2', 'Joint Base McGuire-Dix-Lakehurst', (SELECT id FROM duty_stations WHERE name = 'JB McGuire-Dix-Lakehurst'), now(), now());
INSERT INTO duty_station_names VALUES ('6A01C134-E8CF-FC71-910E-D13DA1E2A425', 'Joint Base Charleston', (SELECT id FROM duty_stations WHERE name = 'JB Charleston'), now(), now());
INSERT INTO duty_station_names VALUES ('C498DC65-2E24-FD3A-3083-BC7DA836F8D3', 'Joint Base Andrews', (SELECT id FROM duty_stations WHERE name = 'JB Andrews'), now(), now());
INSERT INTO duty_station_names VALUES ('DC6D2018-F4A0-5B9F-BDCF-87EC9EF763B9', 'Joint Base Myer-Henderson Hall', (SELECT id FROM duty_stations WHERE name = 'JB Myer-Henderson Hall'), now(), now());
INSERT INTO duty_station_names VALUES ('76EB3F9D-27A8-8A72-FD28-4C961CE1F2BF', 'Joint Base Elmendorf-Richardson', (SELECT id FROM duty_stations WHERE name = 'JB Elmendorf-Richardson'), now(), now());
INSERT INTO duty_station_names VALUES ('BF2E54AF-1BE8-E46A-F9BF-71027B045F63', 'Joint Base Lewis-McChord', (SELECT id FROM duty_stations WHERE name = 'JB Lewis-McChord'), now(), now());
INSERT INTO duty_station_names VALUES ('B9F23528-CAD0-F691-32B4-640E9679F753', 'Joint Base Anacostia–Bolling', (SELECT id FROM duty_stations WHERE name = 'JB Anacostia–Bolling'), now(), now());

-- Fort -> Ft
INSERT INTO duty_station_names VALUES ('EDCDF262-5A87-CF41-B762-2DBD96A6FB30', 'Ft Belvoir', (SELECT id FROM duty_stations WHERE name = 'Fort Belvoir'), now(), now());
INSERT INTO duty_station_names VALUES ('D30B8CA7-E8D7-1D4F-54B4-795CF3EFE40E', 'Ft Benning', (SELECT id FROM duty_stations WHERE name = 'Fort Benning'), now(), now());
INSERT INTO duty_station_names VALUES ('7E9684E2-2D56-AC41-C7C5-8E583652853A', 'Ft Bliss', (SELECT id FROM duty_stations WHERE name = 'Fort Bliss'), now(), now());
INSERT INTO duty_station_names VALUES ('C8395D2F-2FB7-64F0-D458-1B8CA4A42303', 'Ft Bragg', (SELECT id FROM duty_stations WHERE name = 'Fort Bragg'), now(), now());
INSERT INTO duty_station_names VALUES ('E749FABC-B9AB-B076-CE98-CAF15284BEF0', 'Ft Campbell', (SELECT id FROM duty_stations WHERE name = 'Fort Campbell'), now(), now());
INSERT INTO duty_station_names VALUES ('E5261D8F-905D-BC45-34C4-649EC73E12B2', 'Ft Carson', (SELECT id FROM duty_stations WHERE name = 'Fort Carson'), now(), now());
INSERT INTO duty_station_names VALUES ('39E0756A-5E65-9CBD-C291-CDC4B754B0B4', 'Ft Detrick', (SELECT id FROM duty_stations WHERE name = 'Fort Detrick'), now(), now());
INSERT INTO duty_station_names VALUES ('2D763707-52C3-20CA-2FB9-6562CD51720E', 'Ft Drum', (SELECT id FROM duty_stations WHERE name = 'Fort Drum'), now(), now());
INSERT INTO duty_station_names VALUES ('F041483A-E7F1-1078-4926-17C1915C3FCD', 'Ft George G. Meade', (SELECT id FROM duty_stations WHERE name = 'Fort George G. Meade'), now(), now());
INSERT INTO duty_station_names VALUES ('9708F6D8-9AB7-3BC6-1D25-D7B40AC8AF80', 'Ft Gordon', (SELECT id FROM duty_stations WHERE name = 'Fort Gordon'), now(), now());
INSERT INTO duty_station_names VALUES ('4D82CE33-B230-76EB-A8FD-CB34FCF0ADF6', 'Ft Greely', (SELECT id FROM duty_stations WHERE name = 'Fort Greely'), now(), now());
INSERT INTO duty_station_names VALUES ('9C653091-FB68-B3D7-9C2A-B902DF8A8534', 'Ft Hamilton', (SELECT id FROM duty_stations WHERE name = 'Fort Hamilton'), now(), now());
INSERT INTO duty_station_names VALUES ('845A4C54-EC89-C721-41E2-5FA168BFEF19', 'Ft Hood', (SELECT id FROM duty_stations WHERE name = 'Fort Hood'), now(), now());
INSERT INTO duty_station_names VALUES ('8F3EDA9E-4EB0-3BEA-C94A-B9DE4EB3F761', 'Ft Huachuca', (SELECT id FROM duty_stations WHERE name = 'Fort Huachuca'), now(), now());
INSERT INTO duty_station_names VALUES ('427A96C7-F8D9-2CB6-C16F-D5BDF8B5DE1D', 'Ft Irwin', (SELECT id FROM duty_stations WHERE name = 'Fort Irwin'), now(), now());
INSERT INTO duty_station_names VALUES ('721C060A-BA3B-C579-1568-45214217532B', 'Ft Jackson', (SELECT id FROM duty_stations WHERE name = 'Fort Jackson'), now(), now());
INSERT INTO duty_station_names VALUES ('B8CA0A24-9EC4-BA89-7CEF-59D287504C80', 'Ft Knox', (SELECT id FROM duty_stations WHERE name = 'Fort Knox'), now(), now());
INSERT INTO duty_station_names VALUES ('5B8CA984-618E-E5B9-4D4A-A8C5CDEA04BF', 'Ft Leavenworth', (SELECT id FROM duty_stations WHERE name = 'Fort Leavenworth'), now(), now());
INSERT INTO duty_station_names VALUES ('953458BC-1E47-C27A-C989-532E5ACF2503', 'Ft Lee', (SELECT id FROM duty_stations WHERE name = 'Fort Lee'), now(), now());
INSERT INTO duty_station_names VALUES ('8CE01C29-1A60-BA45-7B08-2C359E8FD2F5', 'Ft Leonard Wood', (SELECT id FROM duty_stations WHERE name = 'Fort Leonard Wood'), now(), now());
INSERT INTO duty_station_names VALUES ('3183E825-6BA4-5C1D-23C5-8A5E52E1326E', 'Ft McCoy', (SELECT id FROM duty_stations WHERE name = 'Fort McCoy'), now(), now());
INSERT INTO duty_station_names VALUES ('D89E7920-69F2-6B51-D256-A7C5B5AEA43F', 'Ft Polk', (SELECT id FROM duty_stations WHERE name = 'Fort Polk'), now(), now());
INSERT INTO duty_station_names VALUES ('7EB4D3D7-B061-1F98-CD9B-FE9D62E7E0A1', 'Ft Riley', (SELECT id FROM duty_stations WHERE name = 'Fort Riley'), now(), now());
INSERT INTO duty_station_names VALUES ('A05AD5B7-BA13-CE91-297E-F7352806D5B9', 'Ft Rucker', (SELECT id FROM duty_stations WHERE name = 'Fort Rucker'), now(), now());
INSERT INTO duty_station_names VALUES ('735D6B03-7E97-1360-B7DC-12C02C4D6C93', 'Ft Sill', (SELECT id FROM duty_stations WHERE name = 'Fort Sill'), now(), now());
INSERT INTO duty_station_names VALUES ('5E9485A5-12E3-E50A-D380-E2EF79580E5D', 'Ft Stewart-Hunter', (SELECT id FROM duty_stations WHERE name = 'Fort Stewart-Hunter'), now(), now());
INSERT INTO duty_station_names VALUES ('FABFB16B-6B63-F9EB-8207-AB01F41A4958', 'Ft Wainwright', (SELECT id FROM duty_stations WHERE name = 'Fort Wainwright'), now(), now());

-- nas -> naval air station
INSERT INTO duty_station_names VALUES ('AD125170-56F3-5B31-D611-98A025E0F6B6', 'Naval Air Station Corpus Christi', (SELECT id FROM duty_stations WHERE name = 'NAS Corpus Christi'), now(), now());
INSERT INTO duty_station_names VALUES ('E8A21F8D-32F5-D1FA-D3E1-EA7A15B249CA', 'Naval Air Station Fallon', (SELECT id FROM duty_stations WHERE name = 'NAS Fallon'), now(), now());
INSERT INTO duty_station_names VALUES ('4072375E-FBA9-109D-5A15-870C87CEC618', 'Naval Air Station Key West', (SELECT id FROM duty_stations WHERE name = 'NAS Key West'), now(), now());
INSERT INTO duty_station_names VALUES ('73EC9651-DFE7-C56D-BF65-8BA97B180B52', 'Naval Air Station Lemoore', (SELECT id FROM duty_stations WHERE name = 'NAS Lemoore'), now(), now());
INSERT INTO duty_station_names VALUES ('1F21E047-264F-F937-7D12-15F375E7A67A', 'Naval Air Station Meridian', (SELECT id FROM duty_stations WHERE name = 'NAS Meridian'), now(), now());
INSERT INTO duty_station_names VALUES ('A1B41F6F-5C89-E9C5-C58C-F42A31A434F1', 'Naval Air Station Patuxent River', (SELECT id FROM duty_stations WHERE name = 'NAS Patuxent River'), now(), now());
INSERT INTO duty_station_names VALUES ('E51A89CB-DE3A-36E4-12A8-88EA24132325', 'Naval Air Station Pensacola', (SELECT id FROM duty_stations WHERE name = 'NAS Pensacola'), now(), now());
INSERT INTO duty_station_names VALUES ('A84D2EA3-81CA-1E6F-E65E-3F0F5480D87E', 'Naval Air Station Whidbey Island', (SELECT id FROM duty_stations WHERE name = 'NAS Whidbey Island'), now(), now());

-- nas -> naval air station && jrb -> joint reserve base
INSERT INTO duty_station_names VALUES ('4F56350E-A60C-D734-F983-3E08A3F0909A', 'Naval Air Station Fort Worth Joint Reserve Base', (SELECT id FROM duty_stations WHERE name = 'NAS Fort Worth JRB'), now(), now());
INSERT INTO duty_station_names VALUES ('A1F6B892-ABCA-31D5-61D7-9BD01B4CD49E', 'Naval Air Station Joint Reserve Base New Orleans', (SELECT id FROM duty_stations WHERE name = 'NAS JRB New Orleans'), now(), now());

-- JBSA -> Joint Base San Antonio
INSERT INTO duty_station_names VALUES ('1B353E47-96FA-FA7F-1ECE-3D84278B7416', 'Joint Base San Antonio Randolph', (SELECT id FROM duty_stations WHERE name = 'JBSA Randolph'), now(), now());
INSERT INTO duty_station_names VALUES ('D8D89E38-7294-16D1-CA8D-DB20F1D3FC07', 'Joint Base San Antonio Lackland', (SELECT id FROM duty_stations WHERE name = 'JBSA Lackland'), now(), now());
INSERT INTO duty_station_names VALUES ('45A42E26-18AC-1D9F-AFA7-ADA374B84DF9', 'Joint Base San Antonio Fort Sam Houston', (SELECT id FROM duty_stations WHERE name = 'JBSA Fort Sam Houston'), now(), now());

-- NAVSTA -> Naval Station
INSERT INTO duty_station_names VALUES ('23D253A3-E7E3-2EAB-8B09-5E381C6A3B4C', 'Naval Station Everett', (SELECT id FROM duty_stations WHERE name = 'NAVSTA Everett'), now(), now());
INSERT INTO duty_station_names VALUES ('C1E27469-F12C-4B03-6E35-DA87405E63FA', 'Naval Station Newport', (SELECT id FROM duty_stations WHERE name = 'NAVSTA Newport'), now(), now());

-- NSA -> Naval Support Activity
INSERT INTO duty_station_names VALUES ('A464DAD7-108A-74AC-38F4-8715D637EFD4', 'Naval Support Activity Pensacola', (SELECT id FROM duty_stations WHERE name = 'NAS Pensacola'), now(), now());
INSERT INTO duty_station_names VALUES ('C1C0D5ED-B2B6-CD38-825A-A2E8E23CD482', 'Naval Support Activity Bethesda', (SELECT id FROM duty_stations WHERE name = 'NSA Bethesda'), now(), now());
INSERT INTO duty_station_names VALUES ('563F06BD-B0B5-4FE2-CBD9-74C4D25F1712', 'Naval Support Activity Mid-South', (SELECT id FROM duty_stations WHERE name = 'NSA Mid-South'), now(), now());
INSERT INTO duty_station_names VALUES ('A1B8C310-A13A-7B40-F026-36E601AF7363', 'Naval Support Activity Panama City', (SELECT id FROM duty_stations WHERE name = 'NSA Panama City'), now(), now());
INSERT INTO duty_station_names VALUES ('B01DEBAD-1B62-4A0A-E0A0-3BED71B5E65A', 'Naval Support Activity Saratoga Springs', (SELECT id FROM duty_stations WHERE name = 'NSA Saratoga Springs'), now(), now());

-- ANGBB -> Air National Guard Base
INSERT INTO duty_station_names VALUES ('7D151E41-3C3C-3051-FEB6-7F63897C0A98', 'Ellington Field Air National Guard Base', (SELECT id FROM duty_stations WHERE name = 'Ellington Field ANGB'), now(), now());

-- AAF -> Army Airfield
INSERT INTO duty_station_names VALUES ('98196F72-9797-4784-CA1B-C2ECE9E1BE4D', 'Hunter Army Airfield', (SELECT id FROM duty_stations WHERE name = 'Ellington Field ANGB'), now(), now());

-- AAP -> Army Ammunition Plant
INSERT INTO duty_station_names VALUES ('6FBE04ED-8D5F-DF5E-1AFE-A8A05F78D268', 'McAlester Army Ammunition Plant', (SELECT id FROM duty_stations WHERE name = 'McAlester AAP'), now(), now());

-- MCAGCC -> Marine Corps Air Ground Combat Center
INSERT INTO duty_station_names VALUES ('1E36395E-E410-DAB5-B6B1-D7C4EB6DB6B6', 'Marine Corps Air Ground Combat Center Twentynine Palms', (SELECT id FROM duty_stations WHERE name = 'MCAGCC Twentynine Palms'), now(), now());

-- NAWS -> Naval Air Weapons Station
INSERT INTO duty_station_names VALUES ('97956DC5-94FA-68C2-DACB-79F7AB8EF77F', 'Naval Air Weapons Station China Lake', (SELECT id FROM duty_stations WHERE name = 'NAWS China Lake'), now(), now());

-- FLC -> Fleet Logistics Center && NAVSUP -> Naval Supply
INSERT INTO duty_station_names VALUES ('64FC9C32-CD8A-39E3-86FC-D6357B0BC7AC', 'Naval Supply Fleet Logistics Center Puget Sound', (SELECT id FROM duty_stations WHERE name = 'NAVSUP FLC Puget Sound'), now(), now());

-- NSB -> Naval Subarine Base
INSERT INTO duty_station_names VALUES ('BA450413-2377-5308-96D6-8D7A4C1DA6B6', 'Naval Submarine Base New London', (SELECT id FROM duty_stations WHERE name = 'NSB New London'), now(), now());

-- NCBC -> Naval Construction Battalion Center
INSERT INTO duty_station_names VALUES ('74DF9109-4383-A7EC-8FA7-F4594EB726E1', 'Naval Construction Battalion Center Gulfport', (SELECT id FROM duty_stations WHERE name = 'NCBC Gulfport'), now(), now());

-- NS -> Naval Statiaon
INSERT INTO duty_station_names VALUES ('CD092AC2-9529-E7AD-2D54-B340424D8BF4', 'Naval Station Norfolk', (SELECT id FROM duty_stations WHERE name = 'NS Norfolk'), now(), now());
