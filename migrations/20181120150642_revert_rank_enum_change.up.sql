-- Commissioned officers and warrant officers are different paygrades.
-- They were conflated previously solely because they share entitlements,
-- at least for now.
-- Separating these ranks will protect us in the future from losing information
-- in electronic orders, will ensure that we address users respectfully with
-- the correct rank, and will make it simple to give them separate entitlements
-- if necessary.
-- At last check, none of the users who have enrolled in MilMove are warrant
-- officers. Therefore, it's safe and accurate to consider all users with
-- O_n_W_n ranks to be commissioned officers.
UPDATE service_members SET rank='O_1_ACADEMY_GRADUATE' WHERE rank = 'O_1_W_1_ACADEMY_GRADUATE';
UPDATE service_members SET rank='O_2' WHERE rank = 'O_2_W_2';
UPDATE service_members SET rank='O_3' WHERE rank = 'O_3_W_3';
UPDATE service_members SET rank='O_4' WHERE rank = 'O_4_W_4';
UPDATE service_members SET rank='O_5' WHERE rank = 'O_5_W_5';
-- Academy cadets are West Pointers, while midshipmen are Naval Academy cadets.
-- They were conflated previously solely because they share entitlements, at
-- least for now. That entitlement is not inherent to being a cadet, as
-- aviation cadets have a different entitlement.
UPDATE service_members SET rank='ACADEMY_CADET' WHERE rank = 'ACADEMY_CADET_MIDSHIPMAN';
