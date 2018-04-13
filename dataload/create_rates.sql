DROP TABLE IF EXISTS temp_tsp_discount_rates;

CREATE TABLE temp_tsp_discount_rates (
	rate_cycle text,
	origin text,
	destination text,
	cos text,
	scac text,
	lh_rate numeric(6,2),
	sit_rate numeric(6,2)
);

\copy temp_tsp_discount_rates FROM '/Users/breanneboland/Desktop/datadump/2018 Code 2 Peak Rates.txt' WITH CSV HEADER;
\copy temp_tsp_discount_rates FROM '/Users/breanneboland/Desktop/datadump/2018 Code 2 NonPeak Rates.txt' WITH CSV HEADER;
\copy temp_tsp_discount_rates FROM '/Users/breanneboland/Desktop/datadump/2018 Code D Peak Rates.txt' WITH CSV HEADER;
\copy temp_tsp_discount_rates FROM '/Users/breanneboland/Desktop/datadump/2018 Code D NonPeak Rates.txt' WITH CSV HEADER;

DROP TABLE IF EXISTS temp_tdl_scores;

CREATE TABLE temp_tdl_scores (
	market text,
	origin text,
	destination text,
	cos text,
	quartile int,
	rank int,
	scac text,
	svy_score numeric(8,4),
	rate_score numeric(8,4),
	bvs numeric(8,4)
);

\copy temp_tdl_scores FROM '/Users/breanneboland/Desktop/datadump/(Pre-Decisional FOUO) TDL Scores 15May18-30September18 - Code 2.csv' WITH CSV HEADER;
\copy temp_tdl_scores FROM '/Users/breanneboland/Desktop/datadump/(Pre-Decisional FOUO) TDL Scores 15May18-30September18 - Code D.csv' WITH CSV HEADER;
\copy temp_tdl_scores FROM '/Users/breanneboland/Desktop/datadump/(Pre-Decisional FOUO) TDL Scores - 1Jan-31Jul PP - 2018 NP - Code 2.csv' WITH CSV HEADER;
\copy temp_tdl_scores FROM '/Users/breanneboland/Desktop/datadump/(Pre-Decisional FOUO) TDL Scores - 1Jan-31Jul PP - 2018 NP - Code D.csv' WITH CSV HEADER;

DROP TABLE IF EXISTS tdl_scores_and_discounts;

-- CREATE TABLE tdl_scores_and_discounts (
-- 	rate_cycle text,
-- 	market text,
-- 	origin text,
-- 	destination text,
-- 	cos text,
-- 	quartile int,
-- 	rank int,
-- 	scac text,
-- 	lh_rate numeric(6,2),
-- 	sit_rate numeric(6,2)	,
-- 	svy_score numeric(8,4),
-- 	rate_score numeric(8,4),
-- 	bvs numeric(8,4),
-- 	rate_cycle_start date DEFAULT 2018-05-15,
-- 	rate_cycle_end date DEFAULT 2018-09-30
-- );


CREATE TABLE tdl_scores_and_discounts AS
	SELECT dr.rate_cycle, s.market, s.origin, s.destination, s.cos, s.quartile, s.rank, s.scac, s.bvs, dr.lh_rate, dr.sit_rate FROM temp_tdl_scores AS s
	LEFT JOIN temp_tsp_discount_rates as dr
 	ON s.origin = dr.origin
 	AND s.destination = dr.destination
	AND s.cos = dr.cos
	AND s.scac = dr.scac;

-- Adds rate cycle and performance period dates for TSPP matching
ALTER TABLE tdl_scores_and_discounts ADD COLUMN rate_cycle_start date DEFAULT '2018-05-15';
ALTER TABLE tdl_scores_and_discounts ADD COLUMN rate_cycle_end date DEFAULT '2018-09-30';
ALTER TABLE tdl_scores_and_discounts ADD COLUMN performance_period_start date DEFAULT '2018-5-15';
ALTER TABLE tdl_scores_and_discounts ADD COLUMN performance_period_end date DEFAULT '2018-7-31';

    -- created_at timestamp without time zone NOT NULL,
    -- updated_at timestamp without time zone NOT NULL

-- To assign UUIDs to records as they're created on import. Example here: https://github.com/transcom/mymove/pull/338/files
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
+
+-- Pack rates
+SELECT
+    uuid_generate_v4() as id,
