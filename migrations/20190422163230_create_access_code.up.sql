CREATE TABLE access_codes (
	id uuid PRIMARY KEY,
	service_member_id uuid REFERENCES service_members(id),
	code text UNIQUE NOT NULL,
	move_type text NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT NOW(),
	claimed_at timestamp with time zone
);
