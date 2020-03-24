# How to run the GHC Transit Time Import and Verify Data

To support loading GHC Transit Time data you can use the `bin/ghc-transit-time-parser` to do so.

## Running the parser

You will need to build the parser first.

```sh
make bin/ghc-transit-time-parser
```

Once built you can run the command.

```sh
bin/ghc-transit-time-parser --display --filename [path to your file]/Appendix_C\(i\)_-_Transit_Time_Tables.xlsx
```

Once complete move on to the next section to verify the import

## Verifying the data

The script will output the transit time table model.

To do the verification, open the newly created csv file located in the directory where the parser was ran.

### 1. Make sure csv data matches expectation

The pricing parser will output the csv file, compare this to the data in the spreadsheet. For example the csv data is found in the pricing template sheet `domestic`.

Once you find the main source of the information you can verify that the number of rows reported in the summary is the same as the number of rows in the matching table.

Pricing parser output example:

```sh
2020-03-18T02:59:31.535Z        INFO    transittime/parse_transit_times.go:69           {"DomesticTransitTime": {"ID":"493fae77-e55c-4d9d-aa3a-1641558e8a2b","MaxDaysTransitTime":32,"WeightLbsLower":8000,"WeightLbsUpper":0,"DistanceMilesLower":6751,"DistanceMilesUpper":7000}}
2020/03/18 02:59:31 File created:
2020/03/18 02:59:31 1_hhg_domestic_transit_times_domestic_20200318025931.csv
2020/03/18 02:59:31 Completed processing sheet index 1 with Description HHG Domestic Transit Times
```

## Useful Command Options

You can run the parser with the `--help` flag to see all possible options. Below is a selection of the most commonly needed flags:

* `--filename string` **Required**
  * Filename (including path) of the XLSX to parse for the GHC transit time data import
* `--save-csv`
  * Save output of XLSX sheets to CSV file (default true)
* `--display`
  * Display output of parsed info (default false)
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
