--insert missing pay_grades
INSERT INTO public.pay_grades
(id, grade, grade_description, created_at, updated_at)
VALUES('9a892c59-48d5-4eba-b5f9-193716da8827', 'O_1', 'Officer Grade O_1', now(), now());

-- Army
INSERT INTO pay_grade_ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    (uuid_generate_v4(), '6cb785d0-cabf-479a-a36d-a6aec294a4d0', 'ARMY', 'PVT', 'Private', 1, now(), now()),
    (uuid_generate_v4(), '5f871c82-f259-43cc-9245-a6e18975dde0', 'ARMY', 'PV2', 'Private Second Class', 2, now(), now()),
    (uuid_generate_v4(), '862eb395-86d1-44af-ad47-dec44fbeda30', 'ARMY', 'PFC', 'Private First Class', 3, now(), now()),
    (uuid_generate_v4(), 'bb55f37c-3165-46ba-ad3f-9a477f699990', 'ARMY', 'SPC', 'Specialist', 4, now(), now()),
    (uuid_generate_v4(), 'bb55f37c-3165-46ba-ad3f-9a477f699990', 'ARMY', 'CPL', 'Corporal', 5, now(), now()),
    (uuid_generate_v4(), '3f142461-dca5-4a77-9295-92ee93371330', 'ARMY', 'SGT', 'Sergeant', 6, now(), now()),
    (uuid_generate_v4(), '541aec36-bd9f-4ad2-abb4-d9b63e29dc80', 'ARMY', 'SSG', 'Staff Sergeant', 7, now(), now()),
    (uuid_generate_v4(), '523d57a1-529c-4dfd-8c33-9cb169fd29a0', 'ARMY', 'SFC', 'Sergeant First Class', 8, now(), now()),
    (uuid_generate_v4(), '1d909db0-602f-4724-bd43-8f90a6660460', 'ARMY', 'MSG', 'Master Sergeant', 9, now(), now()),
    (uuid_generate_v4(), '1d909db0-602f-4724-bd43-8f90a6660460', 'ARMY', '1SG', 'First Sergeant', 10, now(), now()),
    (uuid_generate_v4(), 'a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'ARMY', 'SGM', 'Sergeant Major', 11, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'ARMY', 'CSM', 'Command Sergeant Major', 12, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'ARMY', 'SMA', 'Sergeant Major of the Army', 13, now(), now()),
    (uuid_generate_v4(), '6badf8a0-b0ef-4e42-b827-7f63a3987a4b', 'ARMY', 'WO1', 'Warrant Officer 1', 14, now(), now()),
    (uuid_generate_v4(), 'a687a2e1-488c-4943-b9d9-3d645a2712f4', 'ARMY', 'CW2', 'Chief Warrant Officer 2', 15, now(), now()),
    (uuid_generate_v4(), '5a65fb1f-4245-4178-b6a7-cc504c9cbb37', 'ARMY', 'CW3', 'Chief Warrant Officer 3', 16, now(), now()),
    (uuid_generate_v4(), '74db5649-cf66-4af8-939b-d3d7f1f6b7c6', 'ARMY', 'CW4', 'Chief Warrant Officer 4', 17, now(), now()),
    (uuid_generate_v4(), 'ea8cb0e9-15ff-43b4-9e41-7168d01e7553', 'ARMY', 'CW5', 'Chief Warrant Officer 5', 18, now(), now()),
    (uuid_generate_v4(), '9a892c59-48d5-4eba-b5f9-193716da8827', 'ARMY', '2LT', 'Second Lieutenant', 19, now(), now()),
    (uuid_generate_v4(), 'd1b76a01-d8e4-4bd3-98ff-fa93ff7bc790', 'ARMY', '1LT', 'First Lieutenant', 20, now(), now()),
    (uuid_generate_v4(), '5658d67b-d510-4226-9e56-714403ba0f10', 'ARMY', 'CPT', 'Captain', 21, now(), now()),
    (uuid_generate_v4(), 'e83d8f8d-f70b-4db1-99cc-dd983d2fd250', 'ARMY', 'MAJ', 'Major', 22, now(), now()),
    (uuid_generate_v4(), '3bc4b197-7897-4105-80a1-39a0378d7730', 'ARMY', 'LTC', 'Lieutenant Colonel', 23, now(), now()),
    (uuid_generate_v4(), '455a112d-d1e0-4559-81e8-6df664638f70', 'ARMY', 'COL', 'Colonel', 24, now(), now()),
    (uuid_generate_v4(), 'cf664124-9baf-4187-8f28-0908c0f0a5e0', 'ARMY', 'BG', 'Brigadier General', 25, now(), now()),
    (uuid_generate_v4(), '6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20', 'ARMY', 'MG', 'Major General', 26, now(), now()),
    (uuid_generate_v4(), '1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80', 'ARMY', 'LTG', 'Lieutenant General', 27, now(), now()),
    (uuid_generate_v4(), '7fa938ab-1c34-4666-a878-9b989c916d1a', 'ARMY', 'GEN', 'General', 28, now(), now());

-- USAF
INSERT INTO pay_grade_ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    (uuid_generate_v4(), '6cb785d0-cabf-479a-a36d-a6aec294a4d0', 'AIR_FORCE', 'AB', 'Airman Basic', 1, now(), now()),
    (uuid_generate_v4(), '5f871c82-f259-43cc-9245-a6e18975dde0', 'AIR_FORCE', 'Amn', 'Airman', 2, now(), now()),
    (uuid_generate_v4(), '862eb395-86d1-44af-ad47-dec44fbeda30', 'AIR_FORCE', 'A1C', 'Airman First Class', 3, now(), now()),
    (uuid_generate_v4(), 'bb55f37c-3165-46ba-ad3f-9a477f699990', 'AIR_FORCE', 'SrA', 'Senior Airman', 4, now(), now()),
    (uuid_generate_v4(), '3f142461-dca5-4a77-9295-92ee93371330', 'AIR_FORCE', 'SSgt', 'Staff Sergeant', 5, now(), now()),
    (uuid_generate_v4(), '541aec36-bd9f-4ad2-abb4-d9b63e29dc80', 'AIR_FORCE', 'TSgt', 'Technical Sergeant', 6, now(), now()),
    (uuid_generate_v4(), '523d57a1-529c-4dfd-8c33-9cb169fd29a0', 'AIR_FORCE', 'MSgt', 'Master Sergeant', 7, now(), now()),
    (uuid_generate_v4(), '523d57a1-529c-4dfd-8c33-9cb169fd29a0', 'AIR_FORCE', '1st Sgt', 'First Sergeant', 8, now(), now()),
    (uuid_generate_v4(), '1d909db0-602f-4724-bd43-8f90a6660460', 'AIR_FORCE', 'SMSgt', 'Senior Master Sergeant', 9, now(), now()),
    (uuid_generate_v4(), '1d909db0-602f-4724-bd43-8f90a6660460', 'AIR_FORCE', '1st Sgt', 'First Sergeant', 10, now(), now()),
    (uuid_generate_v4(), 'a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'AIR_FORCE', 'CMSgt', 'Chief Master Sergeant', 11, now(), now()),
    (uuid_generate_v4(), 'a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'AIR_FORCE', '1st Sgt', 'First Sergeant', 12, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'AIR_FORCE', 'CCM', 'Command Chief Master Sergeant', 13, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'AIR_FORCE', 'CMSAF', 'Chief Master Sergeant of the Air Force', 14, now(), now()),
    (uuid_generate_v4(), '6badf8a0-b0ef-4e42-b827-7f63a3987a4b', 'AIR_FORCE', 'WO', 'Warrant Officer 1', 15, now(), now()),
    (uuid_generate_v4(), 'a687a2e1-488c-4943-b9d9-3d645a2712f4', 'AIR_FORCE', 'CWO2', 'Chief Warrant Officer 2', 16, now(), now()),
    (uuid_generate_v4(), '5a65fb1f-4245-4178-b6a7-cc504c9cbb37', 'AIR_FORCE', 'CWO3', 'Chief Warrant Officer 3', 17, now(), now()),
    (uuid_generate_v4(), '74db5649-cf66-4af8-939b-d3d7f1f6b7c6', 'AIR_FORCE', 'CWO4', 'Chief Warrant Officer 4', 18, now(), now()),
    (uuid_generate_v4(), 'ea8cb0e9-15ff-43b4-9e41-7168d01e7553', 'AIR_FORCE', 'CWO5', 'Chief Warrant Officer 5', 19, now(), now()),
    (uuid_generate_v4(), '9a892c59-48d5-4eba-b5f9-193716da8827', 'AIR_FORCE', '2d Lt', 'Second Lieutenant', 20, now(), now()),
    (uuid_generate_v4(), 'd1b76a01-d8e4-4bd3-98ff-fa93ff7bc790', 'AIR_FORCE', '1st Lt', 'First Lieutenant', 21, now(), now()),
    (uuid_generate_v4(), '5658d67b-d510-4226-9e56-714403ba0f10', 'AIR_FORCE', 'Capt', 'Captain', 22, now(), now()),
    (uuid_generate_v4(), 'e83d8f8d-f70b-4db1-99cc-dd983d2fd250', 'AIR_FORCE', 'Maj', 'Major', 23, now(), now()),
    (uuid_generate_v4(), '3bc4b197-7897-4105-80a1-39a0378d7730', 'AIR_FORCE', 'Lt Col', 'Lieutenant Colonel', 24, now(), now()),
    (uuid_generate_v4(), '455a112d-d1e0-4559-81e8-6df664638f70', 'AIR_FORCE', 'Col', 'Colonel', 25, now(), now()),
    (uuid_generate_v4(), 'cf664124-9baf-4187-8f28-0908c0f0a5e0', 'AIR_FORCE', 'Brig Gen', 'Brigadier General', 26, now(), now()),
    (uuid_generate_v4(), '6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20', 'AIR_FORCE', 'Maj Gen', 'Major General', 27, now(), now()),
    (uuid_generate_v4(), '1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80', 'AIR_FORCE', 'Lt Gen', 'Lieutenant General', 28, now(), now()),
    (uuid_generate_v4(), '7fa938ab-1c34-4666-a878-9b989c916d1a', 'AIR_FORCE', 'Gen', 'General', 29, now(), now());

-- Marines
INSERT INTO pay_grade_ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    (uuid_generate_v4(), '6cb785d0-cabf-479a-a36d-a6aec294a4d0', 'MARINE_CORPS', 'PVT', 'Private', 1, now(), now()),
    (uuid_generate_v4(), '5f871c82-f259-43cc-9245-a6e18975dde0', 'MARINE_CORPS', 'PFC', 'Private First Class', 2, now(), now()),
    (uuid_generate_v4(), '862eb395-86d1-44af-ad47-dec44fbeda30', 'MARINE_CORPS', 'LCpl', 'Lance Corporal', 3, now(), now()),
    (uuid_generate_v4(), 'bb55f37c-3165-46ba-ad3f-9a477f699990', 'MARINE_CORPS', 'Cpl', 'Corporal', 4, now(), now()),
    (uuid_generate_v4(), '3f142461-dca5-4a77-9295-92ee93371330', 'MARINE_CORPS', 'Sgt', 'Sergeant', 5, now(), now()),
    (uuid_generate_v4(), '541aec36-bd9f-4ad2-abb4-d9b63e29dc80', 'MARINE_CORPS', 'SSgt', 'Staff Sergeant', 6, now(), now()),
    (uuid_generate_v4(), '523d57a1-529c-4dfd-8c33-9cb169fd29a0', 'MARINE_CORPS', 'GySgt', 'Gunnery Sergeant', 7, now(), now()),
    (uuid_generate_v4(), '1d909db0-602f-4724-bd43-8f90a6660460', 'MARINE_CORPS', 'MSgt', 'Master Sergeant', 8, now(), now()),
    (uuid_generate_v4(), '1d909db0-602f-4724-bd43-8f90a6660460', 'MARINE_CORPS', '1st Sgt', 'First Sergeant', 9, now(), now()),
    (uuid_generate_v4(), 'a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'MARINE_CORPS', 'MGySgt', 'Master Gunnery Sergeant', 10, now(), now()),
    (uuid_generate_v4(), 'a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'MARINE_CORPS', 'SgtMaj', 'Sergeant Major', 11, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'MARINE_CORPS', 'SgtMajMC', 'Sergeant Major of the Marine Corps', 12, now(), now()),
    (uuid_generate_v4(), '6badf8a0-b0ef-4e42-b827-7f63a3987a4b', 'MARINE_CORPS', 'WO', 'Warrant Officer 1', 13, now(), now()),
    (uuid_generate_v4(), 'a687a2e1-488c-4943-b9d9-3d645a2712f4', 'MARINE_CORPS', 'CWO2', 'Chief Warrant Officer 2', 14, now(), now()),
    (uuid_generate_v4(), '5a65fb1f-4245-4178-b6a7-cc504c9cbb37', 'MARINE_CORPS', 'CWO3', 'Chief Warrant Officer 3', 15, now(), now()),
    (uuid_generate_v4(), '74db5649-cf66-4af8-939b-d3d7f1f6b7c6', 'MARINE_CORPS', 'CWO4', 'Chief Warrant Officer 4', 16, now(), now()),
    (uuid_generate_v4(), 'ea8cb0e9-15ff-43b4-9e41-7168d01e7553', 'MARINE_CORPS', 'CWO5', 'Chief Warrant Officer 5', 17, now(), now()),
    (uuid_generate_v4(), '9a892c59-48d5-4eba-b5f9-193716da8827', 'MARINE_CORPS', '2ndLt', 'Second Lieutenant', 18, now(), now()),
    (uuid_generate_v4(), 'd1b76a01-d8e4-4bd3-98ff-fa93ff7bc790', 'MARINE_CORPS', '1stLt', 'First Lieutenant', 19, now(), now()),
    (uuid_generate_v4(), '5658d67b-d510-4226-9e56-714403ba0f10', 'MARINE_CORPS', 'Capt', 'Captain', 20, now(), now()),
    (uuid_generate_v4(), 'e83d8f8d-f70b-4db1-99cc-dd983d2fd250', 'MARINE_CORPS', 'Maj', 'Major', 21, now(), now()),
    (uuid_generate_v4(), '3bc4b197-7897-4105-80a1-39a0378d7730', 'MARINE_CORPS', 'LtCol', 'Lieutenant Colonel', 22, now(), now()),
    (uuid_generate_v4(), '455a112d-d1e0-4559-81e8-6df664638f70', 'MARINE_CORPS', 'Col', 'Colonel', 23, now(), now()),
    (uuid_generate_v4(), 'cf664124-9baf-4187-8f28-0908c0f0a5e0', 'MARINE_CORPS', 'BGen', 'Brigadier General', 24, now(), now()),
    (uuid_generate_v4(), '6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20', 'MARINE_CORPS', 'MajGen', 'Major General', 25, now(), now()),
    (uuid_generate_v4(), '1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80', 'MARINE_CORPS', 'LtGen', 'Lieutenant General', 26, now(), now()),
    (uuid_generate_v4(), '7fa938ab-1c34-4666-a878-9b989c916d1a', 'MARINE_CORPS', 'Gen', 'General', 27, now(), now());

-- Navy
INSERT INTO pay_grade_ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    (uuid_generate_v4(), '6cb785d0-cabf-479a-a36d-a6aec294a4d0', 'NAVY', 'SR', 'Seaman Recruit', 1, now(), now()),
    (uuid_generate_v4(), '5f871c82-f259-43cc-9245-a6e18975dde0', 'NAVY', 'SA', 'Seaman Apprentice', 2, now(), now()),
    (uuid_generate_v4(), '862eb395-86d1-44af-ad47-dec44fbeda30', 'NAVY', 'SN', 'Seaman', 3, now(), now()),
    (uuid_generate_v4(), 'bb55f37c-3165-46ba-ad3f-9a477f699990', 'NAVY', 'PO3', 'Petty Officer Third Class', 4, now(), now()),
    (uuid_generate_v4(), '3f142461-dca5-4a77-9295-92ee93371330', 'NAVY', 'PO2', 'Petty Officer Second Class', 5, now(), now()),
    (uuid_generate_v4(), '541aec36-bd9f-4ad2-abb4-d9b63e29dc80', 'NAVY', 'PO1', 'Petty Officer First Class', 6, now(), now()),
    (uuid_generate_v4(), '523d57a1-529c-4dfd-8c33-9cb169fd29a0', 'NAVY', 'CPO', 'Chief Petty Officer', 7, now(), now()),
    (uuid_generate_v4(), '1d909db0-602f-4724-bd43-8f90a6660460', 'NAVY', 'SCPO', 'Senior Chief Petty Officer', 8, now(), now()),
    (uuid_generate_v4(), 'a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'NAVY', 'MCPO', 'Master Chief Petty Officer', 9, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'NAVY', 'CMDCM', 'Command Master Chief Petty Officer', 10, now(), now()),
    (uuid_generate_v4(), 'a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'NAVY', 'FLTCM', 'Fleet Master Chief Petty Officer', 11, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'NAVY', 'FORCM', 'Force Master Chief Petty Officer', 12, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'NAVY', 'MCPON', 'Master Chief Petty Officer of the Navy', 13, now(), now()),
    (uuid_generate_v4(), '6badf8a0-b0ef-4e42-b827-7f63a3987a4b', 'NAVY', 'WO', 'Warrant Officer 1', 14, now(), now()),
    (uuid_generate_v4(), '6badf8a0-b0ef-4e42-b827-7f63a3987a4b', 'NAVY', 'CWO2', 'Chief Warrant Officer 2', 15, now(), now()),
    (uuid_generate_v4(), 'a687a2e1-488c-4943-b9d9-3d645a2712f4', 'NAVY', 'CWO3', 'Chief Warrant Officer 3', 16, now(), now()),
    (uuid_generate_v4(), '5a65fb1f-4245-4178-b6a7-cc504c9cbb37', 'NAVY', 'CWO4', 'Chief Warrant Officer 4', 17, now(), now()),
    (uuid_generate_v4(), '74db5649-cf66-4af8-939b-d3d7f1f6b7c6', 'NAVY', 'CWO5', 'Chief Warrant Officer 5', 18, now(), now()),
    (uuid_generate_v4(), '9a892c59-48d5-4eba-b5f9-193716da8827', 'NAVY', 'ENS', 'Ensign', 19, now(), now()),
    (uuid_generate_v4(), 'd1b76a01-d8e4-4bd3-98ff-fa93ff7bc790', 'NAVY', 'LTJG', 'Lieutenant Junior Grade', 20, now(), now()),
    (uuid_generate_v4(), '5658d67b-d510-4226-9e56-714403ba0f10', 'NAVY', 'LT', 'Lieutenant', 21, now(), now()),
    (uuid_generate_v4(), 'e83d8f8d-f70b-4db1-99cc-dd983d2fd250', 'NAVY', 'LCDR', 'Lieutenant Commander', 22, now(), now()),
    (uuid_generate_v4(), '3bc4b197-7897-4105-80a1-39a0378d7730', 'NAVY', 'CDR', 'Commander', 23, now(), now()),
    (uuid_generate_v4(), '455a112d-d1e0-4559-81e8-6df664638f70', 'NAVY', 'CAPT', 'Captain', 24, now(), now()),
    (uuid_generate_v4(), 'cf664124-9baf-4187-8f28-0908c0f0a5e0', 'NAVY', 'RDML', 'Rear Admiral Lower Half', 25, now(), now()),
    (uuid_generate_v4(), '6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20', 'NAVY', 'RADM', 'Rear Admiral Upper Half', 26, now(), now()),
    (uuid_generate_v4(), '1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80', 'NAVY', 'VADM', 'Vice Admiral', 27, now(), now()),
    (uuid_generate_v4(), '7fa938ab-1c34-4666-a878-9b989c916d1a', 'NAVY', 'ADM', 'Admiral', 28, now(), now());

-- Coast Guard
INSERT INTO pay_grade_ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    (uuid_generate_v4(), '6cb785d0-cabf-479a-a36d-a6aec294a4d0', 'COAST_GUARD', 'SR', 'Seaman Recruit', 1, now(), now()),
    (uuid_generate_v4(), '5f871c82-f259-43cc-9245-a6e18975dde0', 'COAST_GUARD', 'SA', 'Seaman Apprentice', 2, now(), now()),
    (uuid_generate_v4(), '862eb395-86d1-44af-ad47-dec44fbeda30', 'COAST_GUARD', 'SN', 'Seaman', 3, now(), now()),
    (uuid_generate_v4(), 'bb55f37c-3165-46ba-ad3f-9a477f699990', 'COAST_GUARD', 'PO3', 'Petty Officer Third Class', 4, now(), now()),
    (uuid_generate_v4(), '3f142461-dca5-4a77-9295-92ee93371330', 'COAST_GUARD', 'PO2', 'Petty Officer Second Class', 5, now(), now()),
    (uuid_generate_v4(), '541aec36-bd9f-4ad2-abb4-d9b63e29dc80', 'COAST_GUARD', 'PO1', 'Petty Officer First Class', 6, now(), now()),
    (uuid_generate_v4(), '523d57a1-529c-4dfd-8c33-9cb169fd29a0', 'COAST_GUARD', 'CPO', 'Chief Petty Officer', 7, now(), now()),
    (uuid_generate_v4(), '1d909db0-602f-4724-bd43-8f90a6660460', 'COAST_GUARD', 'SCPO', 'Senior Chief Petty Officer', 8, now(), now()),
    (uuid_generate_v4(), 'a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'COAST_GUARD', 'MCPO', 'Master Chief Petty Officer', 9, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'COAST_GUARD', 'CMC', 'Command Master Chief Petty Officer', 10, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'COAST_GUARD', 'MCPOCG', 'Master Chief Petty Officer of the Coast Guard', 11, now(), now()),
    (uuid_generate_v4(), 'a687a2e1-488c-4943-b9d9-3d645a2712f4', 'COAST_GUARD', 'CWO2', 'Chief Warrant Officer 2', 12, now(), now()),
    (uuid_generate_v4(), '5a65fb1f-4245-4178-b6a7-cc504c9cbb37', 'COAST_GUARD', 'CWO3', 'Chief Warrant Officer 3', 13, now(), now()),
    (uuid_generate_v4(), '74db5649-cf66-4af8-939b-d3d7f1f6b7c6', 'COAST_GUARD', 'CWO4', 'Chief Warrant Officer 4', 14, now(), now()),
    (uuid_generate_v4(), '9a892c59-48d5-4eba-b5f9-193716da8827', 'COAST_GUARD', 'ENS', 'Ensign', 15, now(), now()),
    (uuid_generate_v4(), 'd1b76a01-d8e4-4bd3-98ff-fa93ff7bc790', 'COAST_GUARD', 'LTJG', 'Lieutenant Junior Grade', 16, now(), now()),
    (uuid_generate_v4(), '5658d67b-d510-4226-9e56-714403ba0f10', 'COAST_GUARD', 'LT', 'Lieutenant', 17, now(), now()),
    (uuid_generate_v4(), 'e83d8f8d-f70b-4db1-99cc-dd983d2fd250', 'COAST_GUARD', 'LCDR', 'Lieutenant Commander', 18, now(), now()),
    (uuid_generate_v4(), '3bc4b197-7897-4105-80a1-39a0378d7730', 'COAST_GUARD', 'CDR', 'Commander', 19, now(), now()),
    (uuid_generate_v4(), '455a112d-d1e0-4559-81e8-6df664638f70', 'COAST_GUARD', 'CAPT', 'Captain', 20, now(), now()),
    (uuid_generate_v4(), 'cf664124-9baf-4187-8f28-0908c0f0a5e0', 'COAST_GUARD', 'RDML', 'Rear Admiral Lower Half', 21, now(), now()),
    (uuid_generate_v4(), '6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20', 'COAST_GUARD', 'RADM', 'Rear Admiral Upper Half', 22, now(), now()),
    (uuid_generate_v4(), '1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80', 'COAST_GUARD', 'VADM', 'Vice Admiral', 23, now(), now()),
    (uuid_generate_v4(), '7fa938ab-1c34-4666-a878-9b989c916d1a', 'COAST_GUARD', 'ADM', 'Admiral', 24, now(), now());

-- Space Force
INSERT INTO pay_grade_ranks (id, pay_grade_id, affiliation, rank_abbv, rank_name, rank_order, created_at, updated_at) VALUES
    (uuid_generate_v4(), '6cb785d0-cabf-479a-a36d-a6aec294a4d0', 'SPACE_FORCE', 'Spc1', 'Specialist 1', 1, now(), now()),
    (uuid_generate_v4(), '5f871c82-f259-43cc-9245-a6e18975dde0', 'SPACE_FORCE', 'Spc2', 'Specialist 2', 2, now(), now()),
    (uuid_generate_v4(), '862eb395-86d1-44af-ad47-dec44fbeda30', 'SPACE_FORCE', 'Spc3', 'Specialist 3', 3, now(), now()),
    (uuid_generate_v4(), 'bb55f37c-3165-46ba-ad3f-9a477f699990', 'SPACE_FORCE', 'Spc4', 'Specialist 4', 4, now(), now()),
    (uuid_generate_v4(), '3f142461-dca5-4a77-9295-92ee93371330', 'SPACE_FORCE', 'Sgt', 'Sergeant', 5, now(), now()),
    (uuid_generate_v4(), '541aec36-bd9f-4ad2-abb4-d9b63e29dc80', 'SPACE_FORCE', 'TSgt', 'Technical Sergeant', 6, now(), now()),
    (uuid_generate_v4(), '523d57a1-529c-4dfd-8c33-9cb169fd29a0', 'SPACE_FORCE', 'MSgt', 'Master Sergeant', 7, now(), now()),
    (uuid_generate_v4(), '1d909db0-602f-4724-bd43-8f90a6660460', 'SPACE_FORCE', 'SMSgt', 'Senior Master Sergeant', 8, now(), now()),
    (uuid_generate_v4(), 'a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'SPACE_FORCE', 'CMSgt', 'Chief Master Sergeant', 9, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'SPACE_FORCE', 'CCM', 'Command Chief Master Sergeant', 10, now(), now()),
    (uuid_generate_v4(), '911208cc-3d13-49d6-9478-b0a3943435c0', 'SPACE_FORCE', 'CMSSF', 'Chief Master Sergeant of the Space Force', 11, now(), now()),
    (uuid_generate_v4(), '9a892c59-48d5-4eba-b5f9-193716da8827', 'SPACE_FORCE', '2d Lt', 'Second Lieutenant', 12, now(), now()),
    (uuid_generate_v4(), 'd1b76a01-d8e4-4bd3-98ff-fa93ff7bc790', 'SPACE_FORCE', '1st Lt', 'First Lieutenant', 13, now(), now()),
    (uuid_generate_v4(), '5658d67b-d510-4226-9e56-714403ba0f10', 'SPACE_FORCE', 'Capt', 'Captain', 14, now(), now()),
    (uuid_generate_v4(), 'e83d8f8d-f70b-4db1-99cc-dd983d2fd250', 'SPACE_FORCE', 'Maj', 'Major', 15, now(), now()),
    (uuid_generate_v4(), '3bc4b197-7897-4105-80a1-39a0378d7730', 'SPACE_FORCE', 'Lt Col', 'Lieutenant Colonel', 16, now(), now()),
    (uuid_generate_v4(), '455a112d-d1e0-4559-81e8-6df664638f70', 'SPACE_FORCE', 'Col', 'Colonel', 17, now(), now()),
    (uuid_generate_v4(), 'cf664124-9baf-4187-8f28-0908c0f0a5e0', 'SPACE_FORCE', 'Brig Gen', 'Brigadier General', 18, now(), now()),
    (uuid_generate_v4(), '6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20', 'SPACE_FORCE', 'Maj Gen', 'Major General', 19, now(), now()),
    (uuid_generate_v4(), '1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80', 'SPACE_FORCE', 'Lt Gen', 'Lieutenant General', 20, now(), now()),
    (uuid_generate_v4(), '7fa938ab-1c34-4666-a878-9b989c916d1a', 'SPACE_FORCE', 'Gen', 'General', 21, now(), now());

--add pay_grade_rank_id to orders table
alter table orders drop if exists pay_grade_rank_id;
alter table pay_grade_ranks drop constraint if exists pay_grade_rank_id;

ALTER TABLE orders
   ADD pay_grade_rank_id uuid
   	CONSTRAINT fk_orders_pay_grade_rank_id REFERENCES pay_grade_ranks
(id);

--update pay_grade_rank_id in orders where grade:rank is 1:1
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
		  from pay_grade_ranks
		 where pay_grade_id = i.pay_grade_id
		   and affiliation = i.affiliation;

		if v_count = 1 then	--if 1 rank for pay grade then assign pay_grade_rank_id

			update orders o
			   set pay_grade_rank_id = p.id
			  from pay_grade_ranks p
			 where o.id = i.orders_id
			   and p.pay_grade_id = i.pay_grade_id
			   and p.affiliation = i.affiliation
			   and o.pay_grade_rank_id is null;

		end if;

	end loop;

end ';

--remove unused pay grades
delete from pay_grades where grade in
('O_1_ACADEMY_GRADUATE',
'ACADEMY_CADET',
'MIDSHIPMAN',
'AVIATION_CADET');