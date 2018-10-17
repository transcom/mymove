-- There were a few duty stations whose addresses were invalid. They had leading zeros truncated.
-- Given that the address values were generated randomly, we can't update by address ID. This
-- query only will change postal_codes that are both: less that 5 characters long, and are in the
-- same city as the affected duty stations.

UPDATE addresses
SET postal_code = concat('0', postal_code)
WHERE
	char_length(postal_code) < 5
	AND
	city IN ('Buzzards Bay', 'Groton', 'Kittery', 'Naval Station Newport', 'Hanscom AFB');
