INSERT INTO office_users
  VALUES('c219d9e5-2659-427d-be33-bf439251b7f3', NULL, 'Foo', 'Bar', '', 'foo.bar@example.com', '', 'c219d9e5-2659-427d-be33-bf439251b7f3', NOW(), NOW())
  ON CONFLICT DO NOTHING;
INSERT INTO office_users
  VALUES('c219d9e5-2659-427d-be33-bf439251b7f3', NULL, 'Foo', 'Bar', '', 'foo.bar@example.com', '', 'c219d9e5-2659-427d-be33-bf439251b7f3', NOW(), NOW())
  ON CONFLICT DO NOTHING;
