CREATE TABLE entitlements
(
	id uuid PRIMARY KEY NOT NULL,
	dependents_authorized bool,
	total_dependents integer,
	non_temporary_storage bool,
	privately_owned_vehicle bool,
	pro_gear_weight integer,
	pro_gear_weight_spouse integer,
	storage_in_transit integer,
	created_at timestamp WITH TIME ZONE NOT NULL,
	updated_at timestamp WITH TIME ZONE NOT NULL
);
