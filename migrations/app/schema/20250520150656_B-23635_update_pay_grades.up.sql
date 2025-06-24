--B-23635   Jonathan Spight  Alter pay grades table


WITH updates (id, grade, grade_description, updated_at, "sort_order") AS (
    VALUES
        --Enlisted pay Grades
        ('6cb785d0-cabf-479a-a36d-a6aec294a4d0', 'E-1', 'E-1', NOW(), 0),
        ('5f871c82-f259-43cc-9245-a6e18975dde0', 'E-2', 'E-2', NOW(), 1),
        ('862eb395-86d1-44af-ad47-dec44fbeda30', 'E-3', 'E-3', NOW(), 2),
        ('bb55f37c-3165-46ba-ad3f-9a477f699990', 'E-4', 'E-4', NOW(), 3),
        ('3f142461-dca5-4a77-9295-92ee93371330', 'E-5', 'E-5', NOW(), 4),
        ('541aec36-bd9f-4ad2-abb4-d9b63e29dc80', 'E-6', 'E-6', NOW(), 5),
        ('523d57a1-529c-4dfd-8c33-9cb169fd29a0', 'E-7', 'E-7', NOW(), 6),
        ('1d909db0-602f-4724-bd43-8f90a6660460', 'E-8', 'E-8', NOW(), 7),
        ('a5fc8fd2-6f91-492b-abe2-2157d03ec990', 'E-9', 'E-9', NOW(), 8),
        ('911208cc-3d13-49d6-9478-b0a3943435c0', 'E-9-SPECIAL-SENIOR-ENLISTED', 'E-9 Special Senior Enlisted', NOW(), 9),
        --Warrant pay Grades
        ('6badf8a0-b0ef-4e42-b827-7f63a3987a4b', 'W-1', 'W-1', NOW(), 10),
        ('a687a2e1-488c-4943-b9d9-3d645a2712f4', 'W-2', 'W-2', NOW(), 11),
        ('5a65fb1f-4245-4178-b6a7-cc504c9cbb37', 'W-3', 'W-3', NOW(), 12),
        ('74db5649-cf66-4af8-939b-d3d7f1f6b7c6', 'W-4', 'W-4', NOW(), 13),
        ('ea8cb0e9-15ff-43b4-9e41-7168d01e7553', 'W-5', 'W-5', NOW(), 14),
        --Officer pay Grades,
        ('b25998f4-4715-4f41-8986-4c5c8e59fc80', 'O-1', 'O-1', NOW(), 15),
        ('d1b76a01-d8e4-4bd3-98ff-fa93ff7bc790', 'O-2', 'O-2', NOW(), 16),
        ('5658d67b-d510-4226-9e56-714403ba0f10', 'O-3', 'O-3', NOW(), 17),
        ('e83d8f8d-f70b-4db1-99cc-dd983d2fd250', 'O-4', 'O-4', NOW(), 18),
        ('3bc4b197-7897-4105-80a1-39a0378d7730', 'O-5', 'O-5', NOW(), 19),
        ('455a112d-d1e0-4559-81e8-6df664638f70', 'O-6', 'O-6', NOW(), 20),
        ('cf664124-9baf-4187-8f28-0908c0f0a5e0', 'O-7', 'O-7', NOW(), 21),
        ('6e50b04a-52dc-45c9-91d9-4a7b4fa1ab20', 'O-8', 'O-8', NOW(), 22),
        ('1d6e34c3-8c6c-4d4f-8b91-f46bed3f5e80', 'O-9', 'O-9', NOW(), 23),
        ('7fa938ab-1c34-4666-a878-9b989c916d1a', 'O-10','O-10', NOW(), 24),
        --Other pay Grades
        ('63998729-ef74-486e-beea-5b519fa3812f', 'MIDSHIPMAN', 'Midshipman', NOW(), 25),
        ('df749d7e-5007-43cd-8715-2875d281f817', 'AVIATION_CADET', 'Aviation Cadet', NOW(), 26),
        ('8d8c82ea-ea8f-4d7f-9d84-8d186ab7a7c0', 'ACADEMY_CADET', 'Academy Cadet', NOW(), 27),
        ('9e2cb9a5-ace3-4235-9ee7-ebe4cc2a9bc9', 'CIVILIAN_EMPLOYEE', 'Civilian', NOW(), 28)
)
UPDATE pay_grades
SET
    grade = updates.grade,
    grade_description = updates.grade_description,
    "sort_order" = updates."sort_order"
FROM updates
WHERE pay_grades.id = updates.id::uuid;

INSERT INTO pay_grades (id, grade, grade_description,created_at, updated_at, "sort_order")VALUES ('ec620134-d40f-4ebb-bfeb-0e4e0ef06d14', 'ACADEMY_GRADUATE', 'Academy Graduate', NOW(), NOW(), 29);