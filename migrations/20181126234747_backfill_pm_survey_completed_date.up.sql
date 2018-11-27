-- Loop through and backfill PmSurveyCompletedDate using PmSurveyConductedDate
UPDATE shipments
SET pm_survey_completed_date = pm_survey_conducted_date
WHERE pm_survey_completed_date IS NULL
	AND pm_survey_conducted_date IS NOT NULL;
