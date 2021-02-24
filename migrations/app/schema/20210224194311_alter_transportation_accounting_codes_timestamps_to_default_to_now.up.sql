-- Defaulting these to NOW() to reduce size of input file since there are so many rows.

ALTER TABLE public.transportation_accounting_codes ALTER COLUMN created_at SET DEFAULT NOW();
ALTER TABLE public.transportation_accounting_codes ALTER COLUMN updated_at SET DEFAULT NOW();
