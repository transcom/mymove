DROP INDEX moves.moves_show_idx;
CREATE INDEX moves_available_to_prime_and_show_idx ON moves(show, available_to_prime_at);
