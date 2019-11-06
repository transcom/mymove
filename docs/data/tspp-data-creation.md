# Turning TDL scores and TSP discounts into transportation service provider performances

This outlines the steps you need to do to join the two data sources we've traditionally gotten - CSVs or text files of best value scores tied to TDLs, exported one code of service at a time, and CSVs or text files of TSP discount rates, organized by the three pieces of data that make up a TDL (origin, destination, and code of service). If anything behaves in a surprising way, double check the schema detailed here against the organization of your input files. No step of this should alter zero rows, for instance.

Before you begin this process, convert discount rate Excel files or txt files to CSVs, if needed. **Verify that values for SVY_SCORE, RATE_SCORE, and BVS are decimal values (should be formatted like `77.3456`).**

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

### Add column to hold TDL IDs

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

Check for null TDL IDs:

```sql
SELECT count(DISTINCT scac) FROM tdl_scores_and_discounts WHERE tdl_id IS NULL;
```

#### If TDL ID is still null

If count returns anything but 0, you'll need to add new TDL entries.
Check for new entries on the
[Domestic Channel Control List](https://move.mil/sites/default/files/2019-01/2019%20Domestic%20Channel%20Control%20List.pdf).
They'll be highlighted in red.
Create a new temp table for TDLs and add the new entries as follows:

```sql
CREATE TABLE temp_tdls AS SELECT * FROM traffic_distribution_lists;

ALTER TABLE temp_tdls ADD COLUMN import boolean;

INSERT INTO temp_tdls (id, source_rate_area, destination_region, code_of_service, created_at, updated_at, import)
VALUES
  (uuid_generate_v4(), 'US4965500', '1', '2', now(), now(), true),
  (uuid_generate_v4(), 'US4965500', '1', 'D', now(), now(), true),
  /* ... */
  (uuid_generate_v4(), 'US4965500', '10', '2', now(), now(), true);
```

This will add the new entries to the temporary TDL table,
forcing them to adhere to any table constraints
and generating new UUIDs to be consistent across environments.
For info on why having consistent UUIDs is important [see this document](docs/how-to/create-or-deactivate-users.md#a-note-about-uuid_generate_v4)

We'll now [create a new migration](../how-to/migrate-the-database.md#how-to-migrate-the-database) with that data (replace your migration filename):

```bash
make bin/milmove
milmove gen migration -n add_new_tdls
echo -e "INSERT INTO traffic_distribution_lists (id, source_rate_area, destination_region, code_of_service, created_at, updated_at) \nVALUES\n$(
./scripts/psql-dev "\copy (SELECT id, source_rate_area, destination_region, code_of_service FROM temp_tdls WHERE import = true) TO stdout WITH (FORMAT CSV, FORCE_QUOTE *, QUOTE '''');" \
  | awk '{print "  ("$0", now(), now()),"}' \
  | sed '$ s/.$//');" \
  > migrations/20190410152949_add_new_tdls.up.sql
```

This will copy all rows from the table that were included in the new TDL import
and create an insert statement for the data.
You can also use `pg_dump` to generate this migration,
however replacing the timestamps with `now()` allows the environments
to have true `created_at` and `updated_at` timestamps.
Not your locally inserted time.

Once this migration is written, run it and rejoin the TDLs as above.

----

### Add column to hold TSP IDs

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

Similar to TDLs,
there may be missing TSPs.
Currently, we're not using any of this TSP data for production moves,
but we have to satisfy the foreign key constraints for the TSPP data.

Check for missing TSP IDs:

```sql
SELECT count(DISTINCT scac) FROM tdl_scores_and_discounts WHERE tsp_id IS NULL;
```

#### If TSP ID is still null

Note we use GENERATED_UUID4_VAL here to represent a generated UUID, read [this doc](docs/how-to/create-or-deactivate-users.md#a-note-about-uuid_generate_v4) for details.
If this is not 0, add the TSPs:

```sql
CREATE TABLE temp_tsps AS SELECT * FROM transportation_service_providers;

ALTER TABLE temp_tsps ADD COLUMN import boolean;

INSERT INTO temp_tsps (standard_carrier_alpha_code, id, import)
  SELECT DISTINCT ON (scac) scac AS standard_carrier_alpha_code, GENERATED_UUID4_VAL AS id, true AS import
  FROM tdl_scores_and_discounts
  WHERE tsp_id IS NULL;
```

[Generate the migration](../how-to/migrate-the-database.md#how-to-migrate-the-database) (replacing your migration filename):

```bash
make bin/milmove
milmove gen migration -n add_new_scacs
echo -e "INSERT INTO transportation_service_providers (id, standard_carrier_alpha_code, created_at, updated_at) \nVALUES\n$(
./scripts/psql-dev "\copy (SELECT id, standard_carrier_alpha_code FROM temp_tsps WHERE import = true) TO stdout WITH (FORMAT CSV, FORCE_QUOTE *, QUOTE '''');" \
  | awk '{print "  ("$0", now(), now()),"}' \
  | sed '$ s/.$//');" \
  > migrations/20190409010258_add_new_scacs.up.sql
```

Run this migration and rejoin the TSP IDs as above.

----

### Generate data for production import

Now we're ready to combine the datasets together into one table. First, be sure to clear out the `transportation_service_provider_performances` table in case it already contains data:

```sql
DELETE FROM transportation_service_provider_performances;
```

The following command will fill the TSPP table with data. Use your data's current rate cycle and performance period date in lieu of the hard-coded dates below.

> _Rate cycle_ in this context means the rate cycle **period**, so either the peak or non-peak part of the annual rate cycle and **not** the rate cycle itself.
>
> [This document](https://docs.google.com/document/d/12AN1igDt9Acxm9cu1cJA0fiWQIMI3XGs5u_jhCHLo6I) specifies the date ranges for both the performance periods and the rate cycle periods.

Note we use GENERATED_UUID4_VAL here to represent a generated UUID, read [this doc](docs/how-to/create-or-deactivate-users.md#a-note-about-uuid_generate_v4) for details.

```SQL
INSERT INTO
  transportation_service_provider_performances (id, performance_period_start, performance_period_end, traffic_distribution_list_id, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate)
SELECT
  GENERATED_UUID4_VAL as id, '2018-08-01' as performance_period_start, '2018-09-30' as performance_period_end, tdl_id, 0 as offer_count, bvs, tsp_id, now() as created_at, now() as updated_at, '2018-05-15' as rate_cycle_start, '2018-09-30' as rate_cycle_end, lh_rate/100, sit_rate/100
FROM
  tdl_scores_and_discounts;
```

The `/100` of the `sit_rate` and `linehaul_rate` columns accounts for the differences in representing percentages/decimals across sources. This changes integers into decimal representations that fit into our calculations of rates and reimbursements.

### Export TSPP Data

Run this in your terminal to dump the contents of the `transportation_service_provider_performances` table for use elsewhere. Double-check your local db name before assuming this will work.

```pg_dump -h localhost -U postgres -W dev_db --table transportation_service_provider_performances --data-only > tspp_data_dump.pgsql```

Et voilà: TSPPs!

Note that the above `pg_dump` command will generate a file that uses a single `COPY ... FROM stdin` to load data
as opposed to a series of `INSERT` statements.  Using `COPY` can be dramatically faster than `INSERT` -- around 100 times
faster in some cases.  We generally [prefer `INSERT`](https://github.com/transcom/mymove/pull/2670/files#r322981353)
but the amount of data being loaded may make using it simply too expensive for a migration.

**WARNING:** If the generated file is larger than 250 MB then you will not be able to upload the file. This size limit is in place so that the anti-virus software can scan the file. In the case where a file has been generated that is larger than this size you'll need to split the migration file into multiple migration files each with a size smaller than 250 MB.

## Data Validation

The following SQL statements can be used to verify that the above process has been completed successfully. Some numbers may be slightly off
due to natural changes in the data, but any large discrepancies are a potential signal that something has gone wrong somewhere along the way.

NOTE: As of the updates for 2019-10-01 we only import the BVSes for the top performer, instead of everyone. This could lead to more variation in the numbers than in past updates. The numbers below have been updated to reflect the reduced imports.

```text
dev_db=# SELECT COUNT(id) FROM transportation_service_provider_performances;
  count
---------
 847
(1 row)

dev_db=# select min(best_value_score), max(best_value_score) from transportation_service_provider_performances;
   min   |  max
---------+--------
  91.0   | 100.0
(1 row)


dev_db=# select min(sit_rate), max(sit_rate) from transportation_service_provider_performances;
  min | max
------+------
 0.45 | 0.63
(1 row)

dev_db=# select min(linehaul_rate), max(linehaul_rate) from transportation_service_provider_performances;
 min | max
-----+------
 0.4 | 0.68
(1 row)

dev_db=# SELECT min(count), max(count) FROM (
    SELECT transportation_service_provider_id, COUNT(id) FROM transportation_service_provider_performances
    GROUP BY transportation_service_provider_id
    ) as tspp;
 min | max
-----+------
   1 | 511
(1 row)

dev_db=# SELECT count(DISTINCT transportation_service_provider_id) FROM transportation_service_provider_performances;
 count
-------
   35
(1 row)

SELECT CONCAT(((bucket -1) * 100)::text, '-', (bucket * 100)::text) as rows, count(transportation_service_provider_id) as tsps FROM (
    SELECT transportation_service_provider_id, width_bucket(COUNT(id), 0, 2000, 20) as bucket FROM transportation_service_provider_performances
    GROUP BY transportation_service_provider_id
    ) as tspp
   GROUP BY bucket;

   rows    | tsps
-----------+-------
 0-100     |   33
 200-300   |    1
 500-600   |    1
(3 rows)

-- Spot check the data by picking a row from the TDL and TSP text/CSV files and verifying the data:

SELECT source_rate_area, destination_region, code_of_service, performance_period_start, performance_period_end,
       best_value_score, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate, standard_carrier_alpha_code,
       tdl.created_at
FROM traffic_distribution_lists AS tdl
LEFT JOIN transportation_service_provider_performances on tdl.id = transportation_service_provider_performances.traffic_distribution_list_id
LEFT JOIN transportation_service_providers on transportation_service_provider_performances.transportation_service_provider_id = transportation_service_providers.id
WHERE performance_period_start='2019-10-01' and performance_period_end='2019-12-31'
  AND destination_region='1' AND source_rate_area='US11'
  AND code_of_service='D';

```

## Temp Data Clean Up

Vacuum up now that the party's over. Only required if you haven't reset the local database already.

```SQL
DROP TABLE tdl_scores_and_discounts;
DROP TABLE temp_tdl_scores;
DROP TABLE temp_tsp_discount_rates;
```

## Create Secure Migrations

You will have to create a secure migration for this data import. Two files will need to be created,
the file that contains the real data and a local migration (dummy file for dev). Follow the
[secure migration steps](../how-to/migrate-the-database.md#secure-migrations).

### How to create the dummy file

You will need to scrub the data that is in the dummy file. The fields: `linehaul_rate`, `sit_rate`, and `best_value_score`
are company competition sensitive data and needs to scrubbed.

The file will also need to be reduced. Currently, we are picking 2 TSPs per TDL.

We have a [script](../../scripts/export-obfuscated-tspp-sample) to help with this process. The script will backup the TSPP table, make the appropriate reduction of
data and scrubbing of key columns, output the results, then restore the original TSPP table.  You can run it like so:

```sh
./scripts/export-obfuscated-tspp-sample <filename>
```

Complete the [secure migration steps](../how-to/migrate-the-database.md#secure-migrations) to
submit both migration files.
