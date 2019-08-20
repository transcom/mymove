-- Create new storage_in_transit_number_trackers table to support SIT number generation.
CREATE TABLE storage_in_transit_number_trackers (
	year            integer NOT NULL,
	day_of_year     integer NOT NULL CHECK (day_of_year BETWEEN 1 and 366),
	sequence_number integer NOT NULL,
	PRIMARY KEY (year, day_of_year)
);
