-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.

-- README: These values are modified in the `secure/` migration file of the
-- same name in AWS S3. The information here was discussed in this DP3 Slack
-- thread:
-- https://ustcdp3.slack.com/archives/CTNTFJSBA/p1691434456374919?thread_ts=1691432974.070869&cid=CTNTFJSBA
-- This email address is elsewhere in the codebase, but the phone number is a placeholder in this file.

INSERT INTO organizations (id, name, created_at, updated_at, poc_email, poc_phone)
VALUES ('580020C5-58B9-490F-9092-BEBB47CFADB5', 'USTC', now(), now(), 'transcom.scott.tcj9.mbx.mil-move@mail.mil', '(555) 555-5555');
