# Turning TDL scores and TSP discounts into transportation service provider performances

This outlines the steps you need to do to join the two data sources we've traditionally gotten - CSVs or text files of best value scores tide to TDLs, exported one code of service at a time, and CSVs or text files of TSP discount rates, organized by the three pieces of data that make up a TDL (origin, destination, and code of service). If anything behaves in a surprising way, double check the schema detailed here against the organization of your input files. No step of this should alter zero rows, for instance.

If this isn't your first time at the data-loading rodeo today:
`DROP TABLE IF EXISTS temp_tsp_discount_rates;`

## Loading data that includes discount rates (linehaul and SIT)

Before you do this, convert discount rate Excel files or txt files to CSVs, if needed.

Duplicate the format of discount rates CSVs:

```SQL
CREATE TABLE temp_tsp_discount_rates (
  rate_cycle text,
  origin text,
  destination text,
  cos text,
  scac text,
  lh_rate numeric(6,2),
  sit_rate numeric(6,2)
);
```

Now, let's get those discount rates in. The file you need now will include the linehaul and SIT discounts and may have a name like `2018 Code 2 Peak Rates.txt`. The "rates" part is what you're looking for: key columns are LH_RATE and SIT_RATE.

`\copy` in psql terminal is a simpler way of getting this into the db because it requires less in the way of user permissions (unlike the `COPY` command). Use your absolute path for where you stored those CSV files.

`\copy temp_tsp_discount_rates FROM '/add/filename/for/discount/rates/file.csv' WITH CSV HEADER;`
(You'll need to do this for each file you're importing and processing. You _could_ copypasta them together, but these rows tend to have 500,000 rows or more. Let Postgres do the work.)

Now, let's get those best value scores. This file will likely have "TDL scores" in the title. Key columns are RANK and BVS.

We received these as txt files before - a quick change of file extension (.txt to .csv) will get you where you need to go.

In case you already made this table...
`DROP TABLE IF EXISTS temp_tdl_scores;`

Duplicate the format of TDL scores CSVs:

```SQL
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
```

Use the `\copy` command in psql again to import the TDL scores. Use your absolute path for where you stored those CSV files.
`\copy temp_tdl_scores FROM '/add/filename/for/tdl/scores/file.csv' WITH CSV HEADER;`

Now let's combine the important parts of both data sources into one table, which we'll begin to shape into a full set of TSPP data.

If you wanted to create the empty table separately of the command below, these are the column details. You don't need to, though.

```SQL
DROP TABLE IF EXISTS tdl_scores_and_discounts;
CREATE TABLE tdl_scores_and_discounts (
  rate_cycle text,
  market text,
  origin text,
  destination text,
  cos text,
  quartile int,
  rank int,
  scac text,
  lh_rate numeric(6,2),
  sit_rate numeric(6,2) ,
  svy_score numeric(8,4),
  rate_score numeric(8,4),
  bvs numeric(8,4)
);
```

Do this instead. This will create and populate the table described above with the relevant, overlapping details from the table imported earlier with BVSes and the table with discount rates:

```SQL
CREATE TABLE tdl_scores_and_discounts AS
  SELECT
    s.market, s.origin, s.destination, s.cos, s.scac, s.bvs, dr.lh_rate, dr.sit_rate FROM temp_tdl_scores AS s
  LEFT JOIN
    temp_tsp_discount_rates as dr
  ON
    s.origin = dr.origin
  AND
    s.destination = dr.destination
  AND
    s.cos = dr.cos
  AND
    s.scac = dr.scac;
  ```

Add a TDL ID column to fill with this next update:
`ALTER TABLE tdl_scores_and_discounts ADD COLUMN tdl_id uuid;`

Sometimes the data provided to us represents fields (destination, most recently) in different ways. Here's how to alter the destination column to match other data sources (most notably the TDL table) - to change 'REGION 1' to just '1' to make the next step work:

```SQL
UPDATE
  tdl_scores_and_discounts
SET
  destination = RIGHT(destination, char_length(destination) - 7);
  ```

Add TDL IDs to the rows in our interim table:

```SQL
UPDATE
  tdl_scores_and_discounts as tsd
SET
  tdl_id = tdl.id
FROM
  traffic_distribution_lists as tdl
WHERE
  tdl.source_rate_area = tsd.origin
AND
  tdl.destination_region = tsd.destination
AND
  tdl.code_of_service = tsd.cos;
  ```

Make room for TSP IDs:
`ALTER TABLE tdl_scores_and_discounts ADD COLUMN tsp_id uuid;`

Import the TSP IDs:

```SQL
UPDATE
  tdl_scores_and_discounts as tsd
SET
  tsp_id = tsp.id
FROM
  transportation_service_providers tsp
WHERE
  tsd.scac = tsp.standard_carrier_alpha_code;
```

Check the types of the BVS, LH discount rate, and SIT discount rate fields in the existing TSPP table. They need to be numerics, not ints, lest we lose a LOT of important detail:

```SQL
ALTER TABLE
  transportation_service_provider_performances ALTER COLUMN best_value_score TYPE numeric;
```

Let's put it all into the TSPP table. Use your data's current rate cycle and performance period date in lieu of the hard-coded dates below:

```SQL
INSERT INTO
  transportation_service_provider_performances (id, performance_period_start, performance_period_end, traffic_distribution_list_id, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate)
SELECT
  uuid_generate_v4() as id, '2018-05-15' as performance_period_start, '2018-07-31' as performance_period_end, tdl_id, 0 as offer_count, bvs, tsp_id, now() as created_at, now() as updated_at, '2018-05-15' as rate_cycle_start, '2018-09-30' as rate_cycle_end, lh_rate/100, sit_rate/100
FROM
  tdl_scores_and_discounts;
```

The `/100` of the `sit_rate` and `linehaul_rate` columns accounts for the differences in representing percentages/decimals across sources. This changes integers into decimal representations that fit into our calculations of rates and reimbursements.

Vacuum up now that the party's over.

```SQL
DROP TABLE tdl_scores_and_discounts;
DROP TABLE temp_tdl_scores;
DROP TABLE temp_tsp_discount_rates;
```

Run this in your terminal to dump your pretty new table for use elsewhere. Double-check your local db name before assuming this will work.
`pg_dump -h localhost -U postgres -W dev_db --table transportation_service_provider_performances > tspp_data_dump.pgsql`

Et voilà: TSPPs!
