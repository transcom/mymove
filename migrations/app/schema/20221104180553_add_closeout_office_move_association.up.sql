ALTER TABLE moves ADD COLUMN closeout_office_id uuid REFERENCES transportation_offices(id);
