# How to run the GHC Pricing Import and Verify Data

To support loading GHC Pricing data you can use the `bin/ghc-pricing-parser` to do so.

## Running the parser

You will need to build the parser first.

```sh
make bin/ghc-pricing-parser
```

Once built you can run the command. This command will take some time, the sample data xlsx used below takes 5-6 minutes to complete.

```sh
ghc-pricing-parser --filename pkg/parser/pricing/fixtures/pricing_template_2019-09-19_fake-data.xlsx --contract-code=UNIQUECODE --contract-name="Unique Name"
```

Once complete move on to the next section to verify the import

## Verifying the data

The script will output the summary of the staging tables and the rate engine tables that were used for the import. The summary will include the total number of rows inserted as well as the first two rows. The goal here is to spot check the data as an additional verification of the data import.

To do the verification follow the below steps for each of the `re_*` tables. It's not required to do so for the `stage_*` temporary tables but if there is a discrepancy a summary of those is also printed out to help in finding where the issue is. The examples below use the `re_shipment_type_prices` table as an example.

### Tips

If you are having trouble locating the start of the Rate Engine Table summary you can search for `Stage Table import into Rate Engine Tables Complete` in the output.
You only need to look into the Stage / Temp table summary if you wish to debug why the data was inaccurately parsed into the rate engine tables. The heading for those is `XLSX to Stage Table Parsing Complete`

### 1. Make sure table total row count matches expectation

The pricing parser will output total rows imported, compare this to the number of rows with data in the spreadsheet. For example for the `re_shipment_type_prices` table the data is found in the pricing template sheet `5a) Access. and Add. Prices`. This can be determined by checking the matching import file `pkg/services/ghcimport/import_re_shipment_type_prices.go` and seeing which staging data table feeds this rate engine table.

A couple of notes, there are many service items that represent one row in the original XLSX. For example `re_domestic_accessorial_prices` comes from a table in sheet *5a* of the spreadsheet and has the service item *Shuttle Service*. The Service item for shuttle service is split into Origin and Destination records. Second note some rate engine tables require data from more than one sheet.

Once you find the main source of the information you can verify that the number of rows reported in the summary is the same as the number of rows in the matching table.

Pricing parser output example:

```sh
2020/02/07 23:05:42    ---
2020/02/07 23:05:42    re_shipment_type_prices (ReShipmentTypePrice): 7
```

### 2. Verify first and last row matches

If the number of rows matches you can then move to verifying the first and last row are as expected.

Pricing parser output example with first and last row:

```sh
2020/02/07 23:05:42    ---
2020/02/07 23:05:42    re_shipment_type_prices (ReShipmentTypePrice): 7
2020/02/07 23:05:42      first: {ID:9af2b8c0-153f-4069-9f75-aa3983ebbecd ContractID:111058a8-a5de-424f-921a-932fa35a6a2a ServiceID:dbd3a39a-6bb9-42da-b81a-9229df7019cf Market:C Factor:1.2 CreatedAt:2020-02-07 23:05:42.034574 +0000 +0000 UpdatedAt:2020-02-07 23:05:42.034576 +0000 +0000 Contract:{ID:00000000-0000-0000-0000-000000000000 Code: Name: CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC} Service:{ID:00000000-0000-0000-0000-000000000000 Code: Name: CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC}}
2020/02/07 23:05:42       last: {ID:1900c460-1e51-478b-82d2-64a072210be8 ContractID:111058a8-a5de-424f-921a-932fa35a6a2a ServiceID:874cb86a-bc39-4f57-a614-53ee3fcacf14 Market:O Factor:1.45 CreatedAt:2020-02-07 23:05:42.065301 +0000 +0000 UpdatedAt:2020-02-07 23:05:42.065303 +0000 +0000 Contract:{ID:00000000-0000-0000-0000-000000000000 Code: Name: CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC} Service:{ID:00000000-0000-0000-0000-000000000000 Code: Name: CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC}}
```

## Useful Command Options

You can run the parser with the `--help` flag to see all possible options. Below is a selection of the most commonly needed flags

* `--filename string` **Required**
  * Filename including path of the XLSX to parse for Rate Engine GHC import
* `--contract-code string` **Required**
  * Contract code to use for this import
* `--contract-name string`
  * Contract name to use for this import
* `--display`
  * Display output of parsed info
* `--save-csv`
  * Save output of xlsx sheets to CSV file
* `--verify`
  * Default is true, if false skip sheet format verification (default true) this will verify that the xlsx looks as we expect it too, this does not validate data.
* `--re-import`
  * Run GHC Rate Engine Import (default true)
* `--use-temp-tables`
  * Default is true, if false stage tables are NOT temp tables (default true)
* `--drop`
  * Default is false, if true stage tables will be dropped if they exist this is useful in conjunction with turning `--use-temp-tables` off
* `--db-env string`
  * database environment: container, test, development (default "development")
* `--db-name string`
  * Database Name (default "dev_db")
* `--db-host string`
  * Database Hostname (default "localhost")
* `--db-port int`
  * Database Port (default 5432)
* `--db-user string`
  * Database Username (default "postgres")
* `--db-password string`
  * Database Password
