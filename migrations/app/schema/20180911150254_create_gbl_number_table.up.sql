ALTER TABLE shipments ADD COLUMN gbl_number VARCHAR(255);

CREATE TABLE gbl_number_trackers (
    sequence_number INTEGER,
    gbloc VARCHAR(255) UNIQUE
);
