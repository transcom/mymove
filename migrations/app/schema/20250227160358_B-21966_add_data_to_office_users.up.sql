-- Adds rejectedOn date to pre-existing rejected office users.
UPDATE office_users
SET rejected_on = updated_at
WHERE
status ='REJECTED'::public."office_user_status" AND rejected_on is null;
