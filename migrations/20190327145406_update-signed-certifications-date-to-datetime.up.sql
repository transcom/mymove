ALTER TABLE signed_certifications
	ALTER COLUMN date TYPE timestamp USING date::timestamp;
