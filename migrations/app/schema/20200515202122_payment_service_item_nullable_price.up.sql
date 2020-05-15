ALTER TABLE payment_service_items
    ALTER COLUMN price_cents DROP NOT NULL;

-- Since we don't have pricing yet, any existing records with a price_cents of 0
-- should really be null instead
UPDATE payment_service_items
SET price_cents = NULL
WHERE price_cents = 0;
