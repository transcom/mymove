--B-23635   Jonathan Spight   Alter pay grades table

--Enlisted pay Grades
UPDATE pay_grades
    SET grade = 'E-1', grade_description = 'Enlisted Grade E-1', "order" = 0 WHERE id = '6cb785d0-cabf-479a-a36d-a6aec294a4d0';
UPDATE pay_grades
    SET grade = 'E-2', grade_description = 'Enlisted Grade E-2', "order" = 1 WHERE id = '5f871c82-f259-43cc-9245-a6e18975dde0';
UPDATE pay_grades
    SET grade = 'E-3', grade_description = 'Enlisted Grade E-3', "order" = 2 WHERE id = '862eb395-86d1-44af-ad47-dec44fbeda30';
UPDATE pay_grades
    SET grade = 'E-4', grade_description = 'Enlisted Grade E-4', "order" = 3 WHERE id = 'bb55f37c-3165-46ba-ad3f-9a477f699990';
UPDATE pay_grades
    SET grade = 'E-5', grade_description = 'Enlisted Grade E-5', "order" = 4 WHERE id = '3f142461-dca5-4a77-9295-92ee93371330';
UPDATE pay_grades
    SET grade = 'E-6', grade_description = 'Enlisted Grade E-6', "order" = 5 WHERE id = '541aec36-bd9f-4ad2-abb4-d9b63e29dc80';
UPDATE pay_grades
    SET grade = 'E-7', grade_description = 'Enlisted Grade E-7', "order" = 6 WHERE id = '523d57a1-529c-4dfd-8c33-9cb169fd29a0';
UPDATE pay_grades
    SET grade = 'E-8', grade_description = 'Enlisted Grade E-8', "order" = 7 WHERE id = '1d909db0-602f-4724-bd43-8f90a6660460';
UPDATE pay_grades
    SET grade = 'E-9', grade_description = 'Enlisted Grade E-9', "order" = 8 WHERE id = 'a5fc8fd2-6f91-492b-abe2-2157d03ec990';
UPDATE pay_grades
    SET grade = 'E-9-SPECIAL-ENIOR-ENLISTED', grade_description = 'Enlisted Grade E-9 Special Senior Enlisted', "order" = 9 WHERE id = '911208cc-3d13-49d6-9478-b0a3943435c0';

--Warrant pay Grades
UPDATE pay_grades
    SET grade = 'W-1', grade_description = 'Warrant Officer W-1', "order" = 10 WHERE id = '6badf8a0-b0ef-4e42-b827-7f63a3987a4b';
UPDATE pay_grades
    SET grade = 'W-2', grade_description = 'Warrant Officer W-2', "order" = 11 WHERE id = 'a687a2e1-488c-4943-b9d9-3d645a2712f4';
UPDATE pay_grades
    SET grade = 'W-3', grade_description = 'Warrant Officer W-3', "order" = 12 WHERE id = '5a65fb1f-4245-4178-b6a7-cc504c9cbb37';
UPDATE pay_grades
    SET grade = 'W-4', grade_description = 'Warrant Officer W-4', "order" = 13 WHERE id = '74db5649-cf66-4af8-939b-d3d7f1f6b7c6';
UPDATE pay_grades
    SET grade = 'W-5', grade_description = 'Warrant Officer W-5', "order" = 14 WHERE id = 'ea8cb0e9-15ff-43b4-9e41-7168d01e7553';

--Officer pay Grades
UPDATE pay_grades
    SET grade = 'O-1', grade_description = 'Officer Grade O-1', "order" = 15 WHERE id = 'b25998f4-4715-4f41-8986-4c5c8e59fc80';
UPDATE pay_grades
    SET grade = 'O-2', grade_description = 'Officer Grade O-2', "order" = 16 WHERE id = 'd1b76a01-d8e4-4bd3-98ff-fa93ff7bc790';
UPDATE pay_grades
    SET grade = 'O-3', grade_description = 'Officer Grade O-3', "order" = 17 WHERE id = '5658d67b-d510-4226-9e56-714403ba0f10';
UPDATE pay_grades
    SET grade = 'O-4', grade_description = 'Officer Grade O-4', "order" = 18 WHERE id = 'e83d8f8d-f70b-4db1-99cc-dd983d2fd250';
UPDATE pay_grades
    SET grade = '0-5', grade_description = 'Officer Grade O-5', "order" = 19 WHERE id = '3bc4b197-7897-4105-80a1-39a0378d7730';
UPDATE pay_grades
    SET grade = '0-6', grade_description = 'Officer Grade O-6', "order" = 20 WHERE id = '455a112d-d1e0-4559-81e8-6df664638f70';
UPDATE pay_grades
    SET grade = '0-7', grade_description = 'Officer Grade O-7', "order" = 21 WHERE id = 'cf664124-9baf-4187-8f28-0908c0f0a5e0';
UPDATE pay_grades
    SET grade = '0-8', grade_description = 'Officer Grade O-8', "order" = 22 WHERE id = '6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20';
UPDATE pay_grades
    SET grade = '0-9', grade_description = 'Officer Grade O-9', "order" = 23 WHERE id = '1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80';
UPDATE pay_grades
    SET grade = '0-10', grade_description = 'Officer Grade O-10', "order" = 24 WHERE id = '7fa938ab-1c34-4666-a878-9b989c916d1a';

--Other pay Grades
INSERT INTO pay_grades (id, grade, grade_description,created_at, updated_at, "order")VALUES ('63998729-ef74-486e-beea-5b519fa3812f', 'MIDSHIPMAN', 'Midshipman', NOW(), NOW(),25);
INSERT INTO pay_grades (id, grade, grade_description,created_at, updated_at, "order")VALUES ('df749d7e-5007-43cd-8715-2875d281f817', 'AVIATION_CADET', 'Aviation Cadet', NOW(), NOW(),26);
INSERT INTO pay_grades (id, grade, grade_description,created_at, updated_at, "order")VALUES ('8d8c82ea-ea8f-4d7f-9d84-8d186ab7a7c0', 'ACADEMY_CADET', 'Academy Cadet', NOW(), NOW(),27);
UPDATE pay_grades
    SET grade = 'CIVILIAN', grade_description = 'Civilian', "order" = 29 WHERE id = '9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9';