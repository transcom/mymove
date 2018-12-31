# Turning TDL scores and TSP discounts into transportation service provider performances

This outlines the steps you need to do to join the two data sources we've traditionally gotten - CSVs or text files of best value scores tide to TDLs, exported one code of service at a time, and CSVs or text files of TSP discount rates, organized by the three pieces of data that make up a TDL (origin, destination, and code of service). If anything behaves in a surprising way, double check the schema detailed here against the organization of your input files. No step of this should alter zero rows, for instance.

Before you begin this process, convert discount rate Excel files or txt files to CSVs, if needed.

> We will use the `\copy` `psql` command throughout this guide.
>
> `\copy` is a simpler way of getting this into the db because it requires less in the way of user permissions (unlike the `COPY` command). Use your absolute path for where you stored those CSV files.

Note: If you wish to view existing data from production, use the command `bin/run-prod-migrations`. The local development `dev-db`
does not contain the full set of data. You should not need to run the command `bin/run-prod-migrations` to complete the steps outlined here.

## Verify Input Files

Check that the files you are about to import have roughly the correct number of lines in them:

### TDL Scores

```text
496592   TDL Scores - 1Aug2018 PP - NP - Code 2.csv
496592   TDL Scores - 1Aug2018 PP - PK - Code 2.csv
565673   TDL Scores - 1Aug2018 PP - NP - Code D.csv
565673   TDL Scores - 1Aug2018 PP - PK - Code D.csv
```

### TSP Rates

```text
496593   2018 Code 2 NonPeak Rates.txt
496593   2018 Code 2 Peak Rates.txt
565674   2018 Code D NonPeak Rates.txt
565674   2018 Code D Peak Rates.txt
```

### TDL Performance Dates vs TSP Rate Dates

Note that Rates overlap Performance Periods. You may get a new set of TDLs and will have to use existing (Non)Peak Rates.
E.g., To load the performance data for `Performance Period 1` 2019 the 2018 `* NonPeak Rates.txt` files were used.

```text
Rate Cycle:
Peak: 5/15 to 9/30

Performance Periods
1: 1/1   to  5/14
2: 5/15  to  7/31
3: 8/1   to  9/30
4: 10/1  to  12/31

+--------------------------------------------------------------------------------------------+-----------------------------+
| 2018                                                                                       | 2019                        |
+--------------------------------------------------------------------------------------------+-----------------------------+
| Rate Cycle Rate - Non Peak (2017)     | Rate Cycle Rate - Peak (2018)    | Rate Cycle Rate - Non Peak (2018)          |  |
+---------------------------------------+----------------------------------+--------------------------------------------+--+
| Perf Period 1                         | Perf Period 2    | Perf Period 3 | Perf Period 4   | Perf Period 1            |  |
+------------------------------------------+---------------+---------------+-----------------+-----------------------------+
| Jan    | Feb    | Mar    | Apr    | May  | Jun   | Jul   | Aug   | Sept  | Oct | Nov | Dec | Jan | Feb | Mar | Apr | May |
+--------+--------+--------+--------+------+-------+-------+-------+-------+-----+-----+-----+-----+-----+-----+-----+-----+
```

## Load TSP Discount Rates

> If this isn't your first time at the data-loading rodeo today:
> `DROP TABLE IF EXISTS temp_tsp_discount_rates;`

Create a table to hold the incoming data:

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

The files you need now will include the linehaul and SIT discounts and may have a name like `2018 Code 2 Peak Rates.txt`. The "rates" part is what you're looking for: key columns are LH_RATE and SIT_RATE.

You will need to import **two** files, one for each code of service in the part of the rate cycle that applies to the TDL data you just imported.

* If your TDL data is during the peak part of the rate cycle (May 15th - September 30th), import the **peak** rates.
* Otherwise, import the **nonpeak** rates.

```sql
\copy temp_tsp_discount_rates FROM '/add/filename/for/discount/rates/2018 Code D Peak Rates.csv' WITH CSV HEADER;
\copy temp_tsp_discount_rates FROM '/add/filename/for/discount/rates/2018 Code 2 Peak Rates.csv' WITH CSV HEADER;
```

## Load TSP Best Value Scores from TDL data

Now, let's get those best value scores. This file will likely have "TDL scores" in the title. Key columns are RANK and BVS.

> In case you already made this table...
> `DROP TABLE IF EXISTS temp_tdl_scores;`

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

Use the `\copy` command in psql again to import the TDL scores. Again, you'll need to import **two** files based on whether the dates you are working with are in the peak or nonpeak season. Use your absolute path for where you stored those CSV files.

```sql
\copy temp_tdl_scores FROM '/add/filename/for/tdl/scores/TDL Scores - 1Aug2018 PP - PK - Code D.csv' WITH CSV HEADER;
\copy temp_tdl_scores FROM '/add/filename/for/tdl/scores/TDL Scores - 1Aug2018 PP - PK - Code 2.csv' WITH CSV HEADER;
```

## Combining Scores and Discounts

Now let's combine the important parts of both data sources into one table, which we'll begin to shape into a full set of TSPP data.

The following command will create and populate the table described above with the relevant, overlapping details from the table imported earlier with BVSes and the table with discount rates:

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

```sql
ALTER TABLE tdl_scores_and_discounts ADD COLUMN tdl_id uuid;
```

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

```sql
ALTER TABLE tdl_scores_and_discounts ADD COLUMN tsp_id uuid;
```

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

Now we're ready to combine the datasets together into one table. First, be sure to clear out the `transportation_service_provider_performances` table in case it already contains data:

```sql
DELETE FROM transportation_service_provider_performances;
```

The following command will fill the TSPP table with data. Use your data's current rate cycle and performance period date in lieu of the hard-coded dates below.

> _Rate cycle_ in this context means the rate cycle **period**, so either the peak or non-peak part of the annual rate cycle and **not** the rate cycle itself.
>
> [This document](https://docs.google.com/document/d/1BsE_yIx5_6URs4Kp7baRVMhqJ4Ec_q3hdCmwOYR4EIc) specifies the date ranges for both the performance periods and the rate cycle periods.

```SQL
INSERT INTO
  transportation_service_provider_performances (id, performance_period_start, performance_period_end, traffic_distribution_list_id, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate)
SELECT
  uuid_generate_v4() as id, '2018-08-01' as performance_period_start, '2018-09-30' as performance_period_end, tdl_id, 0 as offer_count, bvs, tsp_id, now() as created_at, now() as updated_at, '2018-05-15' as rate_cycle_start, '2018-09-30' as rate_cycle_end, lh_rate/100, sit_rate/100
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

Run this in your terminal to dump the contents of the `transportation_service_provider_performances` table for use elsewhere. Double-check your local db name before assuming this will work.
`pg_dump -h localhost -U postgres -W dev_db --table transportation_service_provider_performances --data-only > tspp_data_dump.pgsql`

Et voilà: TSPPs!

## Data Validation

The following SQL statements can be used to verify that the above process has been completed successfully. Some numbers may be slightly off
due to natural changes in the data, but any large discrepancies are a potential signal that something has gone wrong somewhere along the way.

```text
dev_db=# SELECT COUNT(id) FROM transportation_service_provider_performances;
  count
---------
 1062265
(1 row)

dev_db=# select min(best_value_score), max(best_value_score) from transportation_service_provider_performances;
   min   |  max
---------+--------
 61.0964 | 99.223
(1 row)


dev_db=# select min(sit_rate), max(sit_rate) from transportation_service_provider_performances;
 min | max
-----+------
 0.4 | 0.65
(1 row)

dev_db=# select min(linehaul_rate), max(linehaul_rate) from transportation_service_provider_performances;
 min | max
-----+------
 0.4 | 0.69
(1 row)

dev_db=# select code_of_service, count(tspp.id) from transportation_service_provider_performances tspp left join traffic_distribution_lists tdl ON tspp.traffic_distribution_list_id=tdl.id group by code_of_service;
 code_of_service | count
-----------------+--------
 2               | 496592
 D               | 565673
(2 rows)

dev_db=# SELECT min(count), max(count) FROM (
    SELECT transportation_service_provider_id, COUNT(id) FROM transportation_service_provider_performances
    GROUP BY transportation_service_provider_id
    ) as tspp;
 min | max
-----+------
   1 | 1592
(1 row)

dev_db=# SELECT count(DISTINCT transportation_service_provider_id) FROM transportation_service_provider_performances;
 count
-------
   851
(1 row)

SELECT CONCAT(((bucket -1) * 100)::text, '-', (bucket * 100)::text) as rows, count(transportation_service_provider_id) as tsps FROM (
    SELECT transportation_service_provider_id, width_bucket(COUNT(id), 0, 2000, 20) as bucket FROM transportation_service_provider_performances
    GROUP BY transportation_service_provider_id
    ) as tspp
   GROUP BY bucket;

   rows    | tsps
-----------+-------
 0-100     |   105
 100-200   |     3
 600-700   |    34
 700-800   |    18
 800-900   |    14
 900-1000  |    20
 1000-1100 |    10
 1100-1200 |    17
 1200-1300 |    10
 1300-1400 |    16
 1400-1500 |    59
 1500-1600 |   545
(12 rows)

-- Spot check the data by picking a row from the TDL and TSP text/CSV files and verifying the data:

SELECT source_rate_area, destination_region, code_of_service, performance_period_start, performance_period_end,
       best_value_score, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate, standard_carrier_alpha_code,
       tdl.created_at
FROM traffic_distribution_lists AS tdl
LEFT JOIN transportation_service_provider_performances on tdl.id = transportation_service_provider_performances.traffic_distribution_list_id
LEFT JOIN transportation_service_providers on transportation_service_provider_performances.transportation_service_provider_id = transportation_service_providers.id
WHERE performance_period_start='2019-01-01' and performance_period_end='2019-05-14'
  AND standard_carrier_alpha_code='ABCD'
  AND destination_region='14' AND source_rate_area='US11'
  AND code_of_service='D';

```

## Create Secure Migrations

You will have to create a secure migration for this data import. Two files will need to be created,
the file that contains the real data and a local migration (dummy file for dev). Follow the instructions
at [docs/database.md#secure-migrations](https://github.com/transcom/mymove/blob/master/docs/database.md#secure-migrations)

### Some tips for creating the dummy file

You will need to scrub the data that is in the dummy file. The fields: `linehaul_rate`, `sit_rate`, and `best_value_score`
are company competition sensitive data and needs to scrubbed.

The file will also need to be reduced. Currently, we are picking 2 TSPs per TDL.

The following SQL can be used to do the above mentioned:

* Truncate the table transportation_service_provider_performances:

```sql
TRUNCATE transportation_service_provider_performances CASCADE;
```

* Load the file created from the `pg_dump`:

```sh
bin/psql < tspp_data_dump.pgsql
```

* Reduce the number of TSPs to two (2) TSPs per TDL:

```sql
DELETE FROM transportation_service_provider_performances
WHERE id not in (
      SELECT id  FROM
          (SELECT id, traffic_distribution_list_id, transportation_service_provider_id, performance_period_start, ROW_NUMBER() OVER
            (PARTITION BY (traffic_distribution_list_id, performance_period_start, performance_period_end)) rn
           FROM transportation_service_provider_performances
          ) tmp WHERE (rn = 1 OR rn = 2)
    );
```

* Scrub the data:

```sql
UPDATE transportation_service_provider_performances
SET linehaul_rate=random(),
    sit_rate=random(),
    best_value_score=random_between(60,70);
```

* Run the `pg_dump` again to capture the new local migration file:

```sql
pg_dump -h localhost -U postgres -W dev_db --table transportation_service_provider_performances --data-only > local_migration_tspp_data_dump.pgsql
```

Complete the steps from [docs/database.md#secure-migrations](https://github.com/transcom/mymove/blob/master/docs/database.md#secure-migrations) to
submit both migration files.
