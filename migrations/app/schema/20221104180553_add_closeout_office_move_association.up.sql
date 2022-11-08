ALTER TABLE moves ADD COLUMN closeout_office_id uuid REFERENCES transportation_offices(id);

COMMENT on COLUMN moves.closeout_office_id IS 'The ID of the assiociated transportation office that is the closeout office for a move.';
