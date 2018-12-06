-- Loop through and backfill PmSurveyCompletedAt using PmSurveyConductedDate
UPDATE shipments
SET pm_survey_completed_at = pm_survey_conducted_date
WHERE pm_survey_completed_at IS NULL
	AND pm_survey_conducted_date IS NOT NULL;
