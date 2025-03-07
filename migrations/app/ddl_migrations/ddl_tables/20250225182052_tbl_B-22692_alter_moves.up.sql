ALTER TABLE public.moves
ADD COLUMN IF NOT EXISTS prime_acknowledged_at TIMESTAMPTZ NULL;

COMMENT ON COLUMN public.moves.prime_acknowledged_at IS 'The date and time the prime acknowledged the move';