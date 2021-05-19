# Archiving Old Pricing Data

**[JIRA Story](https://dp3.atlassian.net/browse/MB-4014)**

We load pricing data via database migrations (some public, some private/secure). Because this includes historical
pricing data going back to 2017, running migrations on a new database can take a while, particularly when testing
deployed migrations. In anticipation of possibly squashing our migrations, we would like to evaluate this historical
pricing data and determine what, if any, of it can be deleted or archived in order to speed up the migration process.

Below are tables with pricing data and notes on each:

- `transportation_service_provider_performances` table
  - private data; contains BVS (best value score) and linehaul/SIT discount rates for each TSP (transportation service provider) in each TDL (traffic distribution list)
  - the BVS was used by the old award queue; discount rates used in pricing PPM moves
  - we have 4 periods of data each year; data currently loaded starts on 05/15/2018 (inclusive) and ends on 12/31/2020 (inclusive)
  - since we are not currently processing PPMs in production, we have hard-coded discount rates in `FetchDiscountRates` to avoid the need to load data for new periods but still be able to test in staging
  - no foreign keys to this table
  - about 20,000 records in dev; about 6 million in each deployed environment
- `tariff400ng_*` tables
  - public data; contains rate information used in pricing pre-GHC moves
  - the following tables include effective date ranges for their records:
    - `tariff400ng_full_pack_rates`
    - `tariff400ng_full_unpack_rates`
    - `tariff400ng_item_rates`
    - `tariff400ng_linehaul_rates`
    - `tariff400ng_service_areas`
    - `tariff400ng_shorthaul_rates`
  - we have one set of data each year; data currently loaded starts on 05/15/2017 (inclusive) and ends on 05/15/2021 (exclusive)
  - staging (for now) will try to read these tables if you choose the "move yourself" path; some endpoints will fail, however,
    if you pick a date of 05/15/2021 or later.
  - no foreign keys to these tables
  - about 50,000 records in dev and each deployed environment
- `re_*` tables (except `re_services`)
  - private data; currently used for GHC pricing
  - we only load test (fake) data in non-production environments; covers 06/01/2019 through 05/31/2027
  - no foreign keys to these tables
  - about 68,000 records in non-production environments
  - not eligible for archiving since in active use

## Considered Alternatives

- Delete all `transportation_service_provider_performances` and `tariff400ng_*` records as noted above
- Archive migrations and/or databases, then delete all `transportation_service_provider_performances` and `tariff400ng_*` records
- Delete migrations for `transportation_service_provider_performances` and `tariff400ng_*` records, but don't delete any data from deployed environments
- Do nothing

## Decision Outcome

- Chosen Alternative: Archive migrations and/or databases, then delete all `transportation_service_provider_performances` and `tariff400ng_*` records
- `+` By archiving, we have a way to get back to historical pricing data if needed
- `+` Squashed migrations would be easier to deal with and run faster without the additional 6 million+ rows
- `+` Queries that reference these tables -- particular those that count records -- would run faster
- `-` An archive is not as convenient as being able to query/migrate the data in the live system

## Pros and Cons of the Alternatives

### Delete all `transportation_service_provider_performances` and `tariff400ng_*` records as noted above

- `+` Squashed migrations would be easier to deal with and run faster without the additional 6 million+ rows
- `+` Queries that reference these tables -- particular those that count records -- would run faster
- `-` No straightforward way to get back to the historical pricing data if needed

### Delete migrations for `transportation_service_provider_performances` and `tariff400ng_*` records, but don't delete any data from deployed environments

- `+` Historical pricing data still available in production environments if needed
- `+` Squashed migrations would be easier to deal with and run faster without the additional 6 million+ rows
- `-` The migration process would no longer exactly match the migrations already applied in the deployed environments
- `-` May only be useful when squashing migrations since existing migrations may mix schema changes with data loads from various tables (i.e., you may not be able to just remove an entire migration)
- `-` Queries would still need to deal with the large amounts of data in the deployed environments
- `-` Some records (like `tariff400ng_*` ones) are not currently exposed by the admin UI or any endpoint, minimizing the benefit of keeping the data in the system

### Do nothing

- `+` No additional work to delete historical pricing data
- `+` Historical pricing data still available in production environments if needed
- `-` More data to have to squash, so migrations still take a long time to process
- `-` Queries would still need to deal with the large amounts of data in the deployed environments
- `-` Some records (like `tariff400ng_*` ones) are not currently exposed by the admin UI or any endpoint, minimizing the benefit of keeping it around
