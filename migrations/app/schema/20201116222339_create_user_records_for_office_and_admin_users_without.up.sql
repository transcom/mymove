CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Remove local office and admin users created in secure migrations
DELETE FROM office_users
WHERE email IN ('test1@example.com','test2@example.com','test3@example.com','test4@example.com','test5@example.com','test6@example.com','test7@example.com','test9@example.com','test10@example.com','test11@example.com','test12@example.com','test13@example.com','test14@example.com','test15@example.com','test16@example.com','test17@example.com','test18@example.com','test19@example.com','test20@example.com','test21@example.com','test22@example.com','test23@example.com','test24@example.com','test25@example.com','test26@example.com','test27@example.com','test28@example.com','test29@example.com','test30@example.com','test31@example.com','testy.civ@example.com','test-93780977-321c-463d-9d55-3dc1bdcaa6de@example.com','test-d4078459-54b7-4206-a496-1be0561337fc@example.com','test-7aacde3a-3789-49d8-8acf-ba54a74aa31f@example.com','bilbo@example.com','test-7f94b9c1-5b5f-416e-ac88-98a9db620499@example.com','test-b551a309-832b-4cf1-8091-50c4d0ee15d5@example.com','test-0fe08591-2931-42b7-bfb2-bd23ab48038a@example.com','test-8a9f544e-45f6-449f-a069-97cfde8d5237@example.com','test-6f43cb4e-1dc7-4a40-b292-48b247664810@example.com','test-6c783441-50f1-453a-b25d-ab6a2eee4ac6@example.com','test-1a49f523-8763-450f-bb31-1d505852eb84@example.com','test-27cfa761-c7c5-4852-8e68-958b4f17ab49@example.com','test-a9f9d1f7-7056-4397-a597-37b277f5de29@example.com','test-cb3a6858-2b6b-4312-a395-d1888b7cef25@example.com','test-9a5c760b-b07d-4754-902f-64a0a1de4a8f@example.com','test-89ddd73a-776f-49fb-aecb-968e4fc745f7@example.com','test-7a841b0f-661c-4b71-8bfb-3141d2b2670a@example.com','test-86aef93c-be95-4428-aab8-5a2f3d9e1006@example.com','test-d0488f9a-053c-47af-b850-c02dcc4304e0@example.com','test-7936285c-bf16-4f1e-9037-419450153e45@example.com','test-116f9f21-8652-4502-8329-61066a3ff11c@example.com','test-86da6c6e-db97-4a19-94f8-69277eecd3be@example.com','test-574171d5-db19-422f-833f-aa26814af546@example.com','test-e90be61b-726f-4327-9a5e-b6e25e7fa588@example.com','test-7dfd83b6-4ee9-456e-a83c-61cd98d2d84b@example.com','test-8327d2ee-be49-4be2-8db6-7bc5d6faf09c@example.com','bones@example.com','test-7a3a8cb3-fa7e-4069-b8ae-b5ec10c87501@example.com','test-8ee447f7-fd28-43fc-a5a2-a86376dcf52a@example.com','test-16dd34b1-bb58-44f7-bea8-bbfa8b94f34e@example.com','test-10fbbd60-c720-4208-bb06-836d2240c8df@example.com','test-9bbc0283-25cb-41d8-93be-cd5d80c55294@example.com','test-2cd7478e-d2d3-4b84-b058-ca64548efba4@example.com','test-89bf33d8-8bd8-4024-8b84-6d63fd9719c4@example.com','test-e27037ca-0fbd-4195-82b0-894fe96132ed@example.com','test-d8599360-ffab-46ce-89e6-11d193eba525@example.com','test-6588e234-f606-4094-8339-9034de23a7b0@example.com','frodo@example.com','test32@test.com','ripley@nostromo.space','ash@hyperdyne.biz','alyssa@example.com','thadiun@example.com','marie@example.com','obrien@example.com','bashir@example.com','randj@example.com','test33@example.com','test34@example.com','test35@example.com','test36@example.com','test37@example.com','test38@example.com','test39@example.com','test40@example.com','test41@example.com','test42@example.com','test44@example.com','gibbons@example.com','lumbergh@example.com','waddams@example.com','smykowski@example.com','daniel@example.com','davenbusters7@test.com','dannybegood@example.com','sherriberry@example.com','rosia@mail.com');

DELETE FROM admin_users
WHERE email IN ('example1@truss.works','example2@truss.works','example3@truss.works');

-- Finds office users who have never signed in so we must create a user record
-- for them that can be updated on first sign in
WITH office_never_signed_in AS (
    SELECT *
    FROM office_users
             LEFT JOIN users ON office_users.email = users.login_gov_email
    WHERE users.login_gov_email IS NULL
)
INSERT INTO users (id, login_gov_email, active, created_at, updated_at)
SELECT uuid_generate_v4(), email, TRUE, now(), now()
FROM office_never_signed_in;

-- Now that we've created the user record, update the office_user record with
-- the new user.id
WITH office_associate_user AS (
    SELECT users.*
    FROM office_users
    JOIN users ON office_users.email = users.login_gov_email
    WHERE office_users.user_id IS NULL
)
UPDATE office_users
SET user_id = office_associate_user.id
FROM office_associate_user
WHERE office_users.email = office_associate_user.login_gov_email;

-- Finds admin users who have never signed in so we must create a user record
-- for them that can be updated on first sign in
WITH admin_never_signed_in AS (
    SELECT *
    FROM admin_users
             LEFT JOIN users ON admin_users.email = users.login_gov_email
    WHERE users.login_gov_email IS NULL
)
INSERT INTO users (id, login_gov_email, active, created_at, updated_at)
SELECT uuid_generate_v4(), email, TRUE, now(), now()
FROM admin_never_signed_in;

-- Now that we've created the user record, update the admin_user record with
-- the new user.id
WITH admin_associate_user AS (
    SELECT users.*
    FROM admin_users
    JOIN users ON admin_users.email = users.login_gov_email
    WHERE admin_users.user_id IS NULL
)
UPDATE admin_users
SET user_id = admin_associate_user.id
FROM admin_associate_user
WHERE admin_users.email = admin_associate_user.login_gov_email;