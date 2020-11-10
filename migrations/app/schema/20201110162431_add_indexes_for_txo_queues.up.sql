CREATE INDEX on service_members (affiliation);
CREATE INDEX on service_members (last_name text_pattern_ops);
CREATE INDEX on service_members (edipi);
CREATE INDEX on payment_requests (created_at);
CREATE INDEX on payment_requests (status);
