--create table pay_grade_ranks to associate pay grade and rank by bos
CREATE TABLE IF NOT EXISTS pay_grade_ranks (
	id 	uuid NOT NULL,
	pay_grade_id  uuid NOT NULL
		CONSTRAINT fk_pay_grade_ranks_pay_grade_id REFERENCES pay_grades (id),
	affiliation TEXT NOT NULL,
	rank_abbv  TEXT NOT NULL,
	rank_name  TEXT,
	rank_order INT,
	created_at  timestamp,
	updated_at  timestamp,
	CONSTRAINT pay_grade_ranks_pkey PRIMARY KEY (id),
	CONSTRAINT unique_pay_grade_ranks UNIQUE (pay_grade_id, affiliation, rank_abbv));

COMMENT ON TABLE pay_grade_ranks IS 'Stores ranks and associated pay grades by branch of service';
COMMENT ON COLUMN pay_grade_ranks.pay_grade_id IS 'ID for pay_grade record';
COMMENT ON COLUMN pay_grade_ranks.affiliation IS 'Branch of service';
COMMENT ON COLUMN pay_grade_ranks.rank_abbv IS 'Rank abbreviation';
COMMENT ON COLUMN pay_grade_ranks.rank_name IS 'Rank name';
COMMENT ON COLUMN pay_grade_ranks.rank_order IS 'Rank order';