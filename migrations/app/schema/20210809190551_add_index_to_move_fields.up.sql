CREATE INDEX moves_available_to_prime_and_show_idx ON moves(show, available_to_prime_at);
DROP INDEX moves_show_idx;
