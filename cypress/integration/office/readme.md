# Use this file to run one off data validation scripts

To validate rate for Pre Approval Requests follow these steps:

1. Run prod migrations script against your dev database using ```DB_NAME=dev_db ./bin/run-prod-migrations```
2. Seed data with e2e: ```./bin/generate-test-data -named-scenario="e2e_basic" -env="development"```
3. Setup cypress: ```npx cypress open```
4. Remove ```skip``` from ```PreApprovalFull.js``` file so tests can be run
5. Search and select ```PreApprovalFull.js``` to run the data validation test. You should be able to invoice with all pre approval requests

Expect all PAR's to be accepted and invoiced. If one or more PAR's fail to be accepted or invoiced you might be missing data for that PAR.