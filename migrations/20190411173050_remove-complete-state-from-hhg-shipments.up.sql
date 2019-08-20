-- Removing references to COMPLETED state for all moves
UPDATE public.shipments
SET status = 'DELIVERED'
WHERE status = 'COMPLETED';

UPDATE public.moves
SET status = 'APPROVED'
WHERE status = 'COMPLETED';

