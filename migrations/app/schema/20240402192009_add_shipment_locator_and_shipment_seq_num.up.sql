SET statement_timeout = 0;

ALTER TABLE moves ADD COLUMN IF NOT EXISTS shipment_seq_num INT;
ALTER TABLE mto_shipments ADD COLUMN IF NOT EXISTS shipment_locator TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS shipment_locator_unique_idx ON mto_shipments (shipment_locator);

DO $$
BEGIN
    -- Update existing 'mto_shipments' rows with a 'shipment_locator'
    WITH calculated_values AS (
        SELECT
            s.id AS shipment_id,
            m.id AS move_id,
            m.locator || '-' || LPAD((COALESCE(m.shipment_seq_num, 0) + ROW_NUMBER() OVER (PARTITION BY m.id ORDER BY s.created_at))::text, 2, '0') AS new_locator
        FROM mto_shipments s
        JOIN moves m ON m.id = s.move_id
        WHERE s.shipment_locator IS NULL OR s.shipment_locator = ''
    )
    UPDATE mto_shipments s
    SET shipment_locator = cv.new_locator
    FROM calculated_values cv
    WHERE s.id = cv.shipment_id;

    -- Update the 'shipment_seq_num' in the 'moves' table to reflect the highest sequence number used
    WITH max_seq_nums AS (
        SELECT
            move_id,
            MAX(SUBSTRING(shipment_locator FROM '.*-(\d+)$')::INT) AS max_seq_num
        FROM mto_shipments
        GROUP BY move_id
    )
    UPDATE moves m
    SET shipment_seq_num = msn.max_seq_num
    FROM max_seq_nums msn
    WHERE m.id = msn.move_id
    AND (m.shipment_seq_num IS NULL OR m.shipment_seq_num = 0);

END $$;


CREATE OR REPLACE FUNCTION generate_shipment_locator()
RETURNS TRIGGER AS $$
DECLARE
    locator_prefix TEXT;
	shipment_seq_num_prefix INT;
BEGIN
    -- Fetch the locator from the moves table
    SELECT locator INTO locator_prefix
	FROM moves
	WHERE id = NEW.move_id
	FOR UPDATE;

    -- Increment the shipment_seq_num for the move and fetch the updated value
    UPDATE moves
	SET shipment_seq_num = COALESCE(shipment_seq_num, 0) + 1
	WHERE id = NEW.move_id
	RETURNING shipment_seq_num INTO shipment_seq_num_prefix;

    -- Set the shipment_locator(move locator + '-' + shipment seq num) in the new shipment row
    NEW.shipment_locator := locator_prefix || '-' || LPAD(shipment_seq_num_prefix::text, 2, '0');

    -- Return the modified NEW row to be inserted
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_generate_shipment_locator ON mto_shipments;
CREATE TRIGGER trigger_generate_shipment_locator
BEFORE INSERT ON mto_shipments
FOR EACH ROW
EXECUTE FUNCTION generate_shipment_locator();

COMMENT ON COLUMN mto_shipments.shipment_locator IS 'Stores the new locator for the shipment just like locator in moves table.';
COMMENT ON COLUMN moves.shipment_seq_num IS 'Keeps track of number of shipments created for a move and use it when creating a shipment_locator.';
