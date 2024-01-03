ALTER TABLE sit_extensions
ADD COLUMN customer_expense BOOLEAN DEFAULT FALSE;
COMMENT on COLUMN sit_extensions.customer_expense IS 'Denotes that the TOO rejected this extension request AND converted it to member''s expense (could be used in MTO view/history to show exactly when a shipment was converted)';

ALTER TABLE mto_service_items
ADD COLUMN customer_expense BOOLEAN DEFAULT FALSE;
COMMENT on COLUMN mto_service_items.customer_expense IS 'Whether or not the service member is responsible for expenses of SIT (i.e. if SIT extension request was denied). Only applicable to DOFSIT items.';

-- Ensures that customer_expense is not NULL
CREATE OR REPLACE FUNCTION check_customer_expense() RETURNS TRIGGER AS $$
BEGIN
  IF NEW.customer_expense IS NULL THEN
    NEW.customer_expense := FALSE;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER check_customer_expense_on_update
  BEFORE UPDATE ON mto_service_items
  FOR EACH ROW EXECUTE FUNCTION check_customer_expense();
CREATE TRIGGER check_customer_expense_on_insert
  BEFORE INSERT ON mto_service_items
  FOR EACH ROW EXECUTE FUNCTION check_customer_expense();
