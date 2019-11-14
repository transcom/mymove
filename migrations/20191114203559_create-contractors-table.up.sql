CREATE TABLE contractor
(
	id UUID PRIMARY KEY,
    created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    contract_number text UNIQUE NOT NULL,
	type text NOT NULL
);
