CREATE TABLE contractor
(
	id UUID PRIMARY KEY,
    created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name varchar(80) NOT NULL,
    contract_number varchar(80) UNIQUE NOT NULL,
	type varchar(80) NOT NULL
);
