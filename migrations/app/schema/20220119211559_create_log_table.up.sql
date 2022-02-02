CREATE TABLE log_event_types (
	id uuid
		CONSTRAINT log_event_type_pkey PRIMARY KEY,
	event_type varchar(255),
	event_name varchar(255),
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);


CREATE TABLE activity_logs (
	id uuid
		CONSTRAINT activity_log_pkey PRIMARY KEY,
	activity_user varchar(255),
	source varchar(255),
	entity varchar(255),
	entity_id varchar(255),
	log_event_type varchar(255),
	log_data json,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);

CREATE OR REPLACE FUNCTION address_log_insert ()
	RETURNS TRIGGER
	LANGUAGE PLPGSQL
	AS $$
BEGIN
	IF NEW.postal_code <> OLD.postal_code THEN
		WITH new_record (
			street_address_1,
			street_address_2,
			city,
			state,
			country,
			postal_code
) AS (
			VALUES(NEW.street_address_1,
					NEW.street_address_2,
					NEW.city,
					NEW.state,
					NEW.country,
					NEW.postal_code)
),
old_record (
	street_address_1,
	street_address_2,
	city,
	state,
	country,
	postal_code
) AS (
	VALUES(OLD.street_address_1,
			OLD.street_address_2,
			OLD.city,
			OLD.state,
			OLD.country,
			OLD.postal_code)
),
new_address AS (
	SELECT
		unnest(ARRAY ['street_address_1', 'street_address_2', 'city', 'state', 'country','postal_code']) AS columns,
		unnest(ARRAY [street_address_1,street_address_2, city, state, country, postal_code]) AS
	VALUES
		FROM new_record
),
old_address AS (
	SELECT
		unnest(ARRAY ['street_address_1', 'street_address_2', 'city', 'state', 'country','postal_code']) AS columns,
		unnest(ARRAY [street_address_1,street_address_2, city, state, country, '90201']) AS
	VALUES
		FROM old_record
),
combined_address AS (
	SELECT
		new_address.columns,
		new_address.values,
		old_address.values AS old_values
	FROM
		new_address
		JOIN old_address ON new_address.columns = old_address.columns
),
diff AS (
	SELECT
		combined_address.columns AS column_name,
		combined_address.values AS value
	FROM
		combined_address
	WHERE
	VALUES
		<> old_values
),
json_data AS (
	SELECT
		array_to_json(array_agg(json_build_object(column_name,
					value)))
	FROM
		diff
) INSERT INTO activity_logs (id, activity_user, source, entity, entity_id, log_event_type, log_data, created_at, updated_at)
		VALUES(uuid_generate_v4 (), SESSION_USER, 'SQL Test', 'address', NEW.id, 'update', (
				SELECT
					* FROM json_data), NOW(), NOW());
	RETURN NEW;
END IF;
	RETURN NEW;
END;
$$

CREATE TRIGGER address_update_trigger
BEFORE UPDATE
ON addresses
FOR EACH ROW
EXECUTE PROCEDURE address_log_insert();

