--B-23635   Jonathan Spight   Add column to pay_grades

ALTER TABLE pay_grades
    ADD COLUMN IF NOT EXISTS "order" integer NULL;

COMMENT ON COLUMN pay_grades."order" IS 'Order of the pay grade option.';