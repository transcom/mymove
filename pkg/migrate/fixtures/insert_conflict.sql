INSERT INTO office_users (id, user_id, last_name, first_name, middle_initials, email, telephone, transportation_office_id, created_at, updated_at, disabled)
  VALUES('dbb40692-9e52-4c3c-86fe-70971f94cd89', NULL, 'Foo', 'Bar', '', 'foo.bar@example.com', '333-333-3333', 'c219d9e5-2659-427d-be33-bf439251b7f3', NOW(), NOW(), false)
  ON CONFLICT DO NOTHING;
INSERT INTO office_users (id, user_id, last_name, first_name, middle_initials, email, telephone, transportation_office_id, created_at, updated_at, disabled)
  VALUES('dbb40692-9e52-4c3c-86fe-70971f94cd89', NULL, 'Foo', 'Bar', '', 'foo.bar@example.com', '333-333-3333', 'c219d9e5-2659-427d-be33-bf439251b7f3', NOW(), NOW(), false)
  ON CONFLICT DO NOTHING;
