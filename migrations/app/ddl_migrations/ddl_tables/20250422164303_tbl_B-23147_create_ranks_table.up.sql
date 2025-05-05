--create table ranks to associate pay grade and rank by affiliation
CREATE TABLE IF NOT EXISTS ranks (
	id 	uuid NOT NULL,
	pay_grade_id  uuid NOT NULL
		CONSTRAINT fk_ranks_pay_grade_id REFERENCES pay_grades (id),
	affiliation TEXT NOT NULL,
	rank_abbv  TEXT NOT NULL,
	rank_name  TEXT,
	rank_order INT,
	created_at  timestamp,
	updated_at  timestamp,
	CONSTRAINT ranks_pkey PRIMARY KEY (id),
	CONSTRAINT unique_ranks UNIQUE (pay_grade_id, affiliation, rank_abbv));

COMMENT ON TABLE ranks IS 'Stores ranks and associated pay grades by branch of service';
COMMENT ON COLUMN ranks.pay_grade_id IS 'ID for pay_grade record';
COMMENT ON COLUMN ranks.affiliation IS 'Branch of service';
COMMENT ON COLUMN ranks.rank_abbv IS 'Rank abbreviation';
COMMENT ON COLUMN ranks.rank_name IS 'Rank name';
COMMENT ON COLUMN ranks.rank_order IS 'Rank order';