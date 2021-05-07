ALTER TABLE duty_stations ADD COLUMN provides_services_counseling boolean DEFAULT false NOT NULL;

-- Column comment
COMMENT ON COLUMN duty_stations.provides_services_counseling IS 'Indicates whether a duty station provides services counseling or not';
