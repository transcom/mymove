-- We have a few office_users with e-mail accounts that are mixed case.
UPDATE office_users SET email = lower(email);
