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

\copy temp_tsp_discount_rates FROM '/Users/jimb/go/src/github.com/transcom/mymove/tsp_rates/2018 Code 2 Peak Rates.txt' WITH CSV HEADER;

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

\copy temp_tdl_scores FROM '/Users/jimb/go/src/github.com/transcom/mymove/tsp_rates/(Pre-Decisional FOUO) TDL Scores 15May18-30September18 - Code 2.csv' WITH CSV HEADER;


/* SELECT * FROM temp_tdl_scores AS s */
/* 	LEFT JOIN temp_tsp_discount_rates AS dr */
/* 	ON s.origin = dr.origin */
/* 		AND s.destination = dr.destination */
/* 		AND s.cos = dr.cos */
/* 		AND s.scac = dr.scac; */

/* effective_date_lower date, */
/* effective_date_upper date */
/* UPDATE temp_tdl_scores (effective_date_lower, effective_date_upper) VALUES ('2018-05-15', '2018-09-30'); */
