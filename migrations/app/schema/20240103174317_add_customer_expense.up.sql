-- Adds new columns for convert to customer expense
ALTER TABLE mto_service_items
ADD COLUMN customer_expense_reason TEXT DEFAULT NULL;
COMMENT on COLUMN mto_service_items.customer_expense_reason IS 'Reason for converting a SIT to customer expense';

-- Ensures that customer_expense is not NULL
CREATE OR REPLACE FUNCTION check_customer_expense() RETURNS TRIGGER AS $$
BEGIN
  IF NEW.customer_expense IS NULL THEN
    NEW.customer_expense := FALSE;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS check_customer_expense_on_update ON mto_service_items;
DROP TRIGGER IF EXISTS check_customer_expense_on_insert ON mto_service_items;

CREATE TRIGGER check_customer_expense_on_update
  BEFORE UPDATE ON mto_service_items
  FOR EACH ROW EXECUTE FUNCTION check_customer_expense();

CREATE TRIGGER check_customer_expense_on_insert
  BEFORE INSERT ON mto_service_items
  FOR EACH ROW EXECUTE FUNCTION check_customer_expense();
