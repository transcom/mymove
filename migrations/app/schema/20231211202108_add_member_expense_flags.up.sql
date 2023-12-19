ALTER TABLE sit_extensions
ADD COLUMN customer_expense BOOLEAN DEFAULT FALSE;
COMMENT on COLUMN sit_extensions.customer_expense IS 'Denotes that the TOO rejected this extension request AND converted it to member''s expense (could be used in MTO view/history to show exactly when a shipment was converted)';

ALTER TABLE mto_service_items
ADD COLUMN customer_expense BOOLEAN DEFAULT FALSE;
COMMENT on COLUMN mto_service_items.customer_expense IS 'Whether or not the service member is responsible for expenses of SIT (i.e. if SIT extension request was denied). Only applicable to DOFSIT items.';

-- Ensures that only items with the re_service_code "DOFSIT" or "DDFSIT" can be given the "customer_expense" flag.
CREATE function check_customer_expense()
RETURNS TRIGGER AS $body$
DECLARE re_service_code VARCHAR(20);
BEGIN
  re_service_code := (SELECT code FROM re_services WHERE re_services.id =  NEW.re_service_id); -- Get the service code for the service item.
  IF re_service_code != 'DOFSIT' OR re_service_code != 'DDFSIT' THEN -- If not a domestic origin SIT 1st day, then customer_expense isn't a valid option and must be false.
    NEW.customer_expense := FALSE;
  END IF;
  RETURN NEW;
END;

$body$
language plpgsql;

CREATE TRIGGER check_customer_expense_on_update
  BEFORE UPDATE ON mto_service_items
  FOR EACH ROW EXECUTE FUNCTION check_customer_expense();

CREATE TRIGGER check_customer_expense_on_insert
  BEFORE INSERT ON mto_service_items
  FOR EACH ROW EXECUTE FUNCTION check_customer_expense();

