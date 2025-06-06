--B-23635   Jonathan Spight   Add column to pay_grades

ALTER TABLE pay_grades
    ADD COLUMN IF NOT EXISTS "sort_order" integer NULL;

COMMENT ON COLUMN pay_grades."sort_order" IS 'Sorted order of the pay grade option.';