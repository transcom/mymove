-- Jira ticket: MB-17020
-- Updates office_phone_lines.transportation_office_id and office_emails.transportation_office_id foreign key constraint to ON DELETE CASCADE

-------------------------------------------
-- OFFICE_PHONE_LINES
-------------------------------------------
-- Drop the existing foreign key constraint
ALTER TABLE office_phone_lines
DROP CONSTRAINT office_phone_lines_transportation_office_id_fkey;

-- Add the ON DELETE CASCADE constraint
ALTER TABLE office_phone_lines
ADD CONSTRAINT office_phone_lines_transportation_office_id_fkey
FOREIGN KEY (transportation_office_id)
REFERENCES transportation_offices(id)
ON DELETE CASCADE;

-------------------------------------------
-- OFFICE_EMAILS
-------------------------------------------
-- Drop the existing foreign key constraint
ALTER TABLE office_emails
DROP CONSTRAINT office_emails_transportation_office_id_fkey;

-- Add the ON DELETE CASCADE constraint
ALTER TABLE office_emails
ADD CONSTRAINT office_emails_transportation_office_id_fkey
FOREIGN KEY (transportation_office_id)
REFERENCES transportation_offices(id)
ON DELETE CASCADE;
