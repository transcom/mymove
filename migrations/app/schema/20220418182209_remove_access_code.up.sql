CREATE TABLE archived_access_codes (
  LIKE access_codes
  INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES
);

ALTER TABLE archived_access_codes ADD CONSTRAINT archived_access_codes_service_member_id FOREIGN KEY (service_member_id) REFERENCES service_members(id);

INSERT INTO archived_access_codes SELECT * FROM access_codes;

DROP TABLE access_codes CASCADE;

ALTER TABLE service_members DROP IF EXISTS requires_access_code CASCADE;