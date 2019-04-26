# Importing tariff400ng data for the year

## Tables that need to be updated with the new data

1. `tariff400ng_full_pack_rates`
2. `tariff400ng_full_unpack_rates`
3. `tariff400ng_linehaul_rates`
4. `tariff400ng_service_areas`
5. `tariff400ng_shorthaul_rates`
6. `tariff400ng_item_rates`

## Importing `full_pack_rates`, `full_unpack_rates`, `linehaul_rates`, `service_areas`, and `shorthaul_rates`

1. Clone the Truss fork of the [move.mil repository](https://github.com/trussworks/move.mil)
2. Run `bin/setup` on the command line and make sure there were no errors in populating the seed data.
3. Add the new `xlsx` file to the `lib/data` directory in the following format: `{YEAR} 400NG Baseline Rate.xlsx`.
4. Open `db/seeds.rb`
5. Near the bottom of the file, you'll see some commented code that imports baseline rates for previous years. Add the following and change the date range as needed:
    ```ruby
    puts '-- Seeding 2019 400NG baseline rates...'
    Seeds::BaselineRates.new(
      date_range: Range.new(Date.parse('2019-05-15'), Date.parse('2020-05-14')),
      file_path: Rails.root.join('lib', 'data', '2019 400NG Baseline Rates.xlsx')
    ).seed!
    ```
6. Run `rails db:reset` to drop the database, re-run migrations, and re-run the seeds import.
7. Dump the tables: `pg_dump --inserts -t full_packs -t full_unpacks -t linehauls -t service_areas -t shorthauls move_mil_development`

## Importing `item_rates`

We're going to make use of the work that Patrick Stanger delivered in [this PR](https://github.com/transcom/mymove/pull/1286).

1. Open [this google sheet](https://docs.google.com/spreadsheets/d/1z1O6hvditeVE4AX1UI-XGu0puwIidXA08tVT6VkG254/edit#gid=138983343) alongside the 400ng data you have received for the upcoming year.
2. Visit the `Accessorials` tab in both spreadsheets.
3. In the new data sheet, within the main section and the Alaska waterhaul section, copy all the values to the left of where it says "weight".
4. Paste those values into the corresponding `Accessorials` tab in the other sheet.
5. Repeat this same process for the `Additonal Rates` tab.
6. Head over to the `migration work` tab. Here, you'll find that queries have been generated for you to insert records into the milmove database.
7. Copy all of the values in the `query` column for both the `Additional Rates` table at the top of the sheet and the `Accessorials` table below it.
