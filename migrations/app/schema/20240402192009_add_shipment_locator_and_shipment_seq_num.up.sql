ALTER TABLE moves ADD COLUMN shipment_seq_num INT DEFAULT 0;
ALTER TABLE mto_shipments ADD COLUMN shipment_locator TEXT;


CREATE OR REPLACE FUNCTION generate_shipment_locator()
RETURNS TRIGGER AS $$
DECLARE
    locator_prefix TEXT;
	shipment_seq_num_prefix INT;
BEGIN
    -- Fetch the locator from the moves table
    SELECT locator INTO locator_prefix FROM moves WHERE id = NEW.move_id;

    -- Increment the shipment_seq_num for the move and fetch the updated value
    UPDATE moves SET shipment_seq_num = COALESCE(shipment_seq_num, 0) + 1 WHERE id = NEW.move_id RETURNING shipment_seq_num INTO shipment_seq_num_prefix;

    -- Set the shipment_locator(move locator + '-' + shipment seq num) in the new shipment row
    NEW.shipment_locator := locator_prefix || '-' || LPAD(shipment_seq_num_prefix::text, 2, '0');

    -- Return the modified NEW row to be inserted
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_generate_shipment_locator
BEFORE INSERT ON mto_shipments
FOR EACH ROW
EXECUTE FUNCTION generate_shipment_locator();
