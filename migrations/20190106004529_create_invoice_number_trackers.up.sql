-- Create new invoice_number_trackers table to support invoice number generation.
CREATE TABLE invoice_number_trackers (
    standard_carrier_alpha_code text NOT NULL,
    year integer NOT NULL,
    sequence_number integer NOT NULL,
    PRIMARY KEY(standard_carrier_alpha_code, year)
);
