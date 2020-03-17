# How to run the GHC Pricing Import and Verify Data

To support loading GHC Pricing data you can use the `bin/ghc-pricing-parser` to do so.

## Running the parser

You will need to build the parser first.

```sh
make bin/ghc-pricing-parser
```

Once built you can run the command. This command will take some time, the sample data XLSX used below takes 5-6 minutes to complete.

```sh
ghc-pricing-parser --filename pkg/parser/pricing/fixtures/pricing_template_2019-09-19_fake-data.xlsx --contract-code=UNIQUECODE --contract-name="Unique Name"
```

Once complete move on to the next section to verify the import

## Verifying the data

The script will output the summary of the staging tables and the rate engine tables that were used for the import. The summary will include the total number of rows inserted as well as a sample of two rows. The goal here is to spot check the data as an additional verification of the data import.

To do the verification follow the below steps for each of the `re_*` tables. It's not required to do so for the `stage_*` temporary tables but if there is a discrepancy a summary of those is also printed out to help in finding where the issue is. The examples below use the `re_shipment_type_prices` table as an example.

### Tips

If you are having trouble locating the start of the Rate Engine Table summary you can search for `Stage table import into rate engine tables complete` in the output.
You only need to look into the Stage / Temp table summary if you wish to debug why the data was inaccurately parsed into the rate engine tables. The heading for those is `XLSX to stage table parsing complete`

### 1. Make sure table total row count matches expectation

The pricing parser will output total rows imported, compare this to the number of rows with data in the spreadsheet. For example for the `re_shipment_type_prices` table the data is found in the pricing template sheet `5a) Access. and Add. Prices`. This can be determined by checking the matching import file `pkg/services/ghcimport/import_re_shipment_type_prices.go` and seeing which staging data table feeds this rate engine table.

A couple of notes, there are many service items that represent one row in the original XLSX. For example `re_domestic_accessorial_prices` comes from a table in sheet *5a* of the spreadsheet and has the service item *Shuttle Service*. The Service item for shuttle service is split into Origin and Destination records. Second note some rate engine tables require data from more than one sheet.

Once you find the main source of the information you can verify that the number of rows reported in the summary is the same as the number of rows in the matching table.

Pricing parser output example:

```sh
2020-02-27T17:13:21.044Z  INFO  ghc-pricing-parser/main.go:273  ----
2020-02-27T17:13:21.049Z  INFO  ghc-pricing-parser/main.go:312  re_shipment_type_prices (ReShipmentTypePrice)  {"row count": 7}
```

### 2. Verify two row matches

If the number of rows matches you can then move to verifying the two rows are as expected.

Pricing parser output example with first and second row (note that these are two sample rows and not
in any particular order relative to the spreadsheet):

```sh
2020-02-27T17:13:21.044Z  INFO  ghc-pricing-parser/main.go:273  ----
2020-02-27T17:13:21.049Z  INFO  ghc-pricing-parser/main.go:312  re_shipment_type_prices (ReShipmentTypePrice)  {"row count": 7}
2020-02-27T17:13:21.049Z  INFO  ghc-pricing-parser/main.go:314  first:  {"ReShipmentTypePrice": {"id":"b93c75b2-559b-4990-8a24-a4ac9b40d7c4","contract_id":"7beb7e1b-b5d7-48e4-bd62-82ebf2f6bd96","service_id":"dbd3a39a-6bb9-42da-b81a-9229df7019cf","market":"C","factor":1.2,"created_at":"2020-02-27T17:13:20.884717Z","updated_at":"2020-02-27T17:13:20.88472Z","Contract":{"id":"00000000-0000-0000-0000-000000000000","code":"","name":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"Service":{"id":"00000000-0000-0000-0000-000000000000","code":"","name":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}}
2020-02-27T17:13:21.049Z  INFO  ghc-pricing-parser/main.go:317  second:  {"ReShipmentTypePrice": {"id":"e4b94491-072f-40d5-8915-7877c0a64014","contract_id":"7beb7e1b-b5d7-48e4-bd62-82ebf2f6bd96","service_id":"0e45b6f5-f2f5-4235-94e4-7b4cb899eb5d","market":"C","factor":1.1,"created_at":"2020-02-27T17:13:20.888991Z","updated_at":"2020-02-27T17:13:20.888993Z","Contract":{"id":"00000000-0000-0000-0000-000000000000","code":"","name":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"Service":{"id":"00000000-0000-0000-0000-000000000000","code":"","name":"","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}}
```

## Useful Command Options

You can run the parser with the `--help` flag to see all possible options. Below is a selection of the most commonly needed flags:

* `--filename string` **Required**
  * Filename (including path) of the XLSX to parse for the GHC rate engine data import
* `--contract-code string` **Required**
  * Contract code to use for this import
* `--contract-name string`
  * Contract name to use for this import; if not provided, the contract-code value will be used
* `--display`
  * Display output of parsed info (default false)
* `--save-csv`
  * Save output of XLSX sheets to CSV file (default false)
* `--verify`
  * Perform sheet format verification -- but does not validate data (default true)
* `--re-import`
  * Perform the import from staging tables to GHC rate engine tables (default true)
* `--use-temp-tables`
  * Make the staging tables be temp tables that don't persist after import (default true)
* `--drop`
  * Drop any existing staging tables prior to creating them; useful when turning `--use-temp-tables` off (default false)
* `--db-env string`
  * Database environment: container, test, development (default "development")
* `--db-name string`
  * Database name (default "dev_db")
* `--db-host string`
  * Database hostname (default "localhost")
* `--db-port int`
  * Database port (default 5432)
* `--db-user string`
  * Database username (default "postgres")
* `--db-password string`
  * Database password
