-- New Column
ALTER TABLE mto_service_item_customer_contacts
ADD COLUMN date_of_contact timestamptz NOT NULL;

-- Comment On Column
COMMENT ON COLUMN mto_service_item_customer_contacts.date_of_contact IS 'The date of attempted contact with the customer by the prime corresponding to the time_military column';
COMMENT ON COLUMN mto_service_item_customer_contacts.time_military IS 'The time of attempted contact with the customer by the prime, in military format (HHMMZ), corresponding to the date_of_contact column';
