--B-23635   Jonathan Spight  Alter pay grades table


WITH updates (id, grade, grade_description, "order") AS (
    VALUES
        --Enlisted pay Grades
        ('6cb785d0-cabf-479a-a36d-a6aec294a4d0', 'E-1', 'Enlisted Grade E-1', 0),
        ('5f871c82-f259-43cc-9245-a6e18975dde0', 'E-2', 'Enlisted Grade E-2', 1),
        ('862eb395-86d1-44af-ad47-dec44fbeda30', 'E-3', 'Enlisted Grade E-3', 2),
        ('bb55f37c-3165-46ba-ad3f-9a477f699990', 'E-4', 'Enlisted Grade E-4', 3),
        ('3f142461-dca5-4a77-9295-92ee93371330', 'E-5', 'Enlisted Grade E-5', 4),
        ('541aec36-bd9f-4ad2-abb4-d9b63e29dc80', 'E-6', 'Enlisted Grade E-6', 5),
        ('523d57a1-529c-4dfd-8c33-9cb169fd29a0', 'E-7', 'Enlisted Grade E-7', 6),
        ('1d909db0-602f-4724-bd43-8f90a6660460', 'E-8', 'Enlisted Grade E-8', 7),
        ('a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'E-9', 'Enlisted Grade E-9', 8),
        ('911208cc-3d13-49d6-9478-b0a3943435c0', 'E-9-SPECIAL-SENIOR-ENLISTED', 'Enlisted Grade E-9 Special Senior Enlisted', 9),
        --Warrant pay Grades
        ('6badf8a0-b0ef-4e42-b827-7f63a3987a4b', 'W-1', 'Warrant Officer W-1', 10),
        ('a687a2e1-488c-4943-b9d9-3d645a2712f4', 'W-2', 'Warrant Officer W-2', 11),
        ('5a65fb1f-4245-4178-b6a7-cc504c9cbb37', 'W-3', 'Warrant Officer W-3', 12),
        ('74db5649-cf66-4af8-939b-d3d7f1f6b7c6', 'W-4', 'Warrant Officer W-4', 13),
        ('ea8cb0e9-15ff-43b4-9e41-7168d01e7553', 'W-5', 'Warrant Officer W-5', 14),
        --Officer pay Grades
        ('b25998f4-4715-4f41-8986-4c5c8e59fc80', 'O-1', 'Officer Grade O-1', 15),
        ('d1b76a01-d8e4-4bd3-98ff-fa93ff7bc790', 'O-2', 'Officer Grade O-2', 16),
        ('5658d67b-d510-4226-9e56-714403ba0f10', 'O-3', 'Officer Grade O-3', 17),
        ('e83d8f8d-f70b-4db1-99cc-dd983d2fd250', 'O-4', 'Officer Grade O-4', 18),
        ('3bc4b197-7897-4105-80a1-39a0378d7730', '0-5', 'Officer Grade O-5', 19),
        ('455a112d-d1e0-4559-81e8-6df664638f70', '0-6', 'Officer Grade O-6', 20),
        ('cf664124-9baf-4187-8f28-0908c0f0a5e0', '0-7', 'Officer Grade O-7', 21),
        ('6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20', '0-8', 'Officer Grade O-8', 22),
        ('1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80', '0-9', 'Officer Grade O-9', 23),
        ('7fa938ab-1c34-4666-a878-9b989c916d1a', 'O-10', 'Officer Grade O-10', 24),
        --Other pay Grades
        ('63998729-ef74-486e-beea-5b519fa3812f', 'MIDSHIPMAN', 'Midshipman', 25),
        ('df749d7e-5007-43cd-8715-2875d281f817', 'AVIATION_CADET', 'Aviation Cadet', 26),
        ('8d8c82ea-ea8f-4d7f-9d84-8d186ab7a7c0', 'ACADEMY_CADET', 'Academy Cadet', 27),
        ('9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9', 'CIVILIAN_EMPLOYEE', 'Civilian', 29)
)
UPDATE pay_grades
SET
    grade = updates.grade,
    grade_description = updates.grade_description,
    "order" = updates."order"
FROM updates
WHERE pay_grades.id = updates.id::uuid;

INSERT INTO pay_grades (id, grade, grade_description,created_at, updated_at, "order")VALUES ('ec620134-d40f-4ebb-bfeb-0e4e0ef06d14', 'ACADEMY_GRADUATE', 'Academy Graduate', NOW(), NOW(),28);