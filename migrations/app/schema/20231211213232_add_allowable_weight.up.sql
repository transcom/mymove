-- Adding new field to weight_tickets to allow capture of the maximum reimbursable benefit amount for PPM shipments.
ALTER TABLE public.weight_tickets ADD allowable_weight int4;

-- Column comments
COMMENT ON COLUMN public.weight_tickets.allowable_weight IS 'Stores the maximum reimbursable benefit weight';
