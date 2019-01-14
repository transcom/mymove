# How To Run One Off Data Validation (Cypress) Tests

## To validate rate for Pre Approval Requests follow these steps

1. Run prod migrations script against your dev database using ```DB_NAME=dev_db ./bin/run-prod-migrations```
2. Seed data with e2e: ```./bin/generate-test-data -named-scenario="e2e_basic" -env="development"```
3. Setup cypress: ```npx cypress open```
4. Remove ```skip``` from ```cypress/data-validation/<type>/<filename>``` file so tests can be run
5. Search and select ```testName``` to run the data validation test.
6. **PII data purging** please make sure you remove all data from your dev db by running this command: ```db_dev_reset```

## Validate expected data and behavior

If test(s) fails, check for bad or invalid data in the database