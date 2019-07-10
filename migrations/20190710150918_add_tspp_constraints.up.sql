-- Change new columns to be NOT NULL now that they have data.
-- Note: Using table constraint instead of column not null constraint to avoid locking entire table.
ALTER TABLE transportation_service_provider_performances
	ADD CONSTRAINT tspp_quartile_not_null CHECK (quartile IS NOT NULL) NOT VALID,
	ADD CONSTRAINT tspp_rank_not_null CHECK (rank IS NOT NULL) NOT VALID,
	ADD CONSTRAINT tspp_survey_score_not_null CHECK (survey_score IS NOT NULL) NOT VALID,
	ADD CONSTRAINT tspp_rate_score_not_null CHECK (rate_score IS NOT NULL) NOT VALID;

-- Now validate the constraint (which shouldn't block).
ALTER TABLE transportation_service_provider_performances
	VALIDATE CONSTRAINT tspp_quartile_not_null,
	VALIDATE CONSTRAINT tspp_rank_not_null,
	VALIDATE CONSTRAINT tspp_survey_score_not_null,
	VALIDATE CONSTRAINT tspp_rate_score_not_null;
