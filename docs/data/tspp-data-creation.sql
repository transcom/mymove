-- First, let's get those discount rates in. This file includes the linehaul and SIT discounts and may have a name like "2018 Code 2 Peak Rates.txt". The "rates" part is what you're looking for: key columns are LH_RATE and SIT_RATE.
-- If this isn't your first time at the data-loading rodeo today:
DROP TABLE IF EXISTS temp_tsp_discount_rates;

-- Loading data that includes discount rates (linehaul and SIT)
  -- First, convert discount rate Excel files to CSVs, if needed.

-- Duplicates format of discount rates CSVs
CREATE TABLE temp_tsp_discount_rates (
  rate_cycle text,
  origin text,
  destination text,
  cos text,
  scac text,
  lh_rate numeric(6,2),
  sit_rate numeric(6,2)
);

-- /copy in psql terminal is simpler because it requires less in the way of user permissions (the way a COPY command can be). Use your absolute path for where you stored those CSV files.
\copy temp_tsp_discount_rates FROM '/add/filename/for/discount/rates/file.csv' WITH CSV HEADER;

-- Now, let's get those best value scores. This file will likely have "TDL scores" in the title. Key columns are RANK and BVS.
-- We received these as txt files before - a quick change of file extension (.txt to .csv) will get you where you need to go.

-- In case you already made this table...
DROP TABLE IF EXISTS temp_tdl_scores;

-- Duplcates format of TDL scores CSVs
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

-- /copy in psql terminal is simpler because it requires less in the way of user permissions (the way a COPY command can be). Use your absolute path for where you stored those CSV files.
\copy temp_tdl_scores FROM '/add/filename/for/tdl/scores/file.csv' WITH CSV HEADER;

-- Let's combine the important parts of both data sources into one that we'll begin to shape into a full set of TSPP data.

-- If you wanted to create the empty table separately of the command below, these are the column details.
-- DROP TABLE IF EXISTS tdl_scores_and_discounts;
-- CREATE TABLE tdl_scores_and_discounts (
--  rate_cycle text,
--  market text,
--  origin text,
--  destination text,
--  cos text,
--  quartile int,
--  rank int,
--  scac text,
--  lh_rate numeric(6,2),
--  sit_rate numeric(6,2) ,
--  svy_score numeric(8,4),
--  rate_score numeric(8,4),
--  bvs numeric(8,4)
-- );

-- This will create and populate the table described above with the relevant, overlapping details from the table imported earlier with BVSes and the table with discount rates.
CREATE TABLE tdl_scores_and_discounts AS
  SELECT s.market, s.origin, s.destination, s.cos, s.scac, s.bvs, dr.lh_rate, dr.sit_rate FROM temp_tdl_scores AS s
  LEFT JOIN temp_tsp_discount_rates as dr
  ON s.origin = dr.origin
  AND s.destination = dr.destination
  AND s.cos = dr.cos
  AND s.scac = dr.scac;

-- Add TDL ID column to fill with this next update.
ALTER TABLE tdl_scores_and_discounts ADD COLUMN tdl_id uuid;

-- This has a side effect of creating rows with a TDL ID and no other data; fix before PR.
UPDATE tdl_scores_and_discounts as tsd
SET    tdl_id = tdl.id
FROM   traffic_distribution_lists tdl
WHERE  tdl.source_rate_area = tsd.origin
  AND tdl.destination_region = tsd.destination
  AND tdl.code_of_service = tsd.cos;

-- This fixes that side effect but is obviously not ideal.
DELETE FROM tdl_scores_and_discounts WHERE market is null;

-- Import TSP IDs
ALTER TABLE tdl_scores_and_discounts ADD COLUMN tsp_id uuid;

UPDATE tdl_scores_and_discounts as tsd
SET tsp_id = tsp.id
FROM transportation_service_providers tsp
WHERE tsd.scac = tsp.standard_carrier_alpha_code;

-- Update BVS, LH discount rate, and SIT discount rate to be numerics, not ints (lest we lose a LOT of important detail)
ALTER TABLE transportation_service_provider_performances ALTER COLUMN best_value_score TYPE numeric;
ALTER TABLE transportation_service_provider_performances ALTER COLUMN linehaul_rate TYPE numeric;
ALTER TABLE transportation_service_provider_performances ALTER COLUMN sit_rate TYPE numeric;

-- Put that big old product of joins into the TSPP table. Use your data's current rate cycle and performance period date in lieu of the hard-coded dates below.
INSERT INTO transportation_service_provider_performances (id, performance_period_start, performance_period_end, traffic_distribution_list_id, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate)
  SELECT uuid_generate_v4() as id, '2018-05-15' as performance_period_start, '2018-07-31' as performance_period_end, tdl_id, 0 as offer_count, bvs, tsp_id, now() as created_at, now() as updated_at, '2018-05-15' as rate_cycle_start, '2018-09-30' as rate_cycle_end, lh_rate, sit_rate
  FROM tdl_scores_and_discounts;

DROP TABLE tdl_scores_and_discounts;
DROP TABLE temp_tdl_scores;
DROP TABLE temp_tsp_discount_rates;
