CREATE TABLE payment_request_edi_files (
    id UUID PRIMARY KEY,
    payment_request_number TEXT NOT NULL,
    edi_string TEXT NOT NULL,
    file_name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
