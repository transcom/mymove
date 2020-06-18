# Scripts

This directory holds the scripts that are not compiled go code. For
compiled go code please look in the `bin/` directory of the project.

If you want to see if scripts are not listed in this file you can run
`find-scripts-missing-in-readme`.

## Dev Environment

These scripts are primarily used for managing the developers
environment.

| Script Name | Description |
| --- | --- |
| `check-aws-cli-version` | checks the awscli version required for the project |
| `check-aws-vault-version` | checks the aws-vault version required for the project |
| `check-bash-version` | Script helps ensure that /etc/shells has all the correct entries in it |
| `check-chamber-version` | checks the chamber version required for the project |
| `check-go-version` | checks the go version required for the project |
| `check-gopath` | checks the go path is correct for the project |
| `check-hosts-file` | Script helps ensure that /etc/hosts has all the correct entries in it |
| `check-node-version` | checks the node version required for the project |
| `check-opensc-version` | checks the opensc version required for the project |
| `kill-process-on-port` | asks to kill a process running on the specified port |
| `prereqs` | validate if all prerequisite programs have been installed |

## AWS Scripts

These scripts are used for interacting with AWS or secrets in the AWS System Manager Parameter Store

| Script Name | Description |
| --- | --- |
| `aws` | Linked to aws-vault-wrapper. Runs the aws binary |
| `aws-vault-wrapper` | A wrapper to ensure AWS credentials are in the environment |
| `chamber` | Linked to aws-vault-wrapper. Runs chamber binary |
| `check-aws-vault-version` | Checks the aws-vault version required for the project |

## Operations Scripts

These scripts are used to operate the system.

| Script Name | Description |
| --- | --- |
| `deploy-app` | Deploy the app |
| `deploy-app-client-tls` | Deploy the app client-tls |
| `deploy-app-migrations` | Deploy the app migrations |
| `deploy-app-tasks` | Deploy the app tasks |
| `download-alb-logs` | Download alb logs for the given environment and dates to a local folder |
| `dupe-secrets` | Dupes experimental secrets to target params |
| `health-tls-check` | Run health and TLS version checks. |
| `scan-alb-logs` | Scan alb logs for specific http codes. |

## Pre-commit Scripts

These scripts are used primarily to check our code before
committing.

| Script Name | Description |
| --- | --- |
| `commit-msg` | Ensure JIRA issue is tagged to commit message |
| `gen-docs-index` | generate index for documents |
| `pre-commit-go-imports` | modify imports in go files |
| `pre-commit-go-lint` | modify go files with linting rules |
| `pre-commit-go-mod` | modify `go.mod` and `go.sum` to match whats in the project |
| `pre-commit-go-vet` | analyze code with `go vet` |
| `pre-commit-swagger-validate` | Pre-commit hook to run `swagger validate` that is different based on environment |

## CircleCI Scripts

These scripts are primarily used for CircleCI workflows.

| Script Name | Description |
| --- | --- |
| `check-deployed-commit` | checks that the deployed commit and given commit match. |
| `check-generated-code` | checks that the generated code has not changed |
| `circleci-announce-broken-branch` | announce that a branch is broken |
| `compare-deployed-commit` | checks that the given commit is ahead of the currently deployed commit |
| `do-exclusively` | CircleCI's current recommendation for roughly serializing a subset of build commands for a given branch |
| `ecr-describe-image-scan-findings` | Checks an uploaded image scan results |
| `ecs-deploy-service-container` | Updates the named service with the given name, image, and environment. |
| `ecs-deploy-task-container` | Updates the named task with the given name, image, and environment. |
| `ecs-restart-services` | Restarted the ECS services associated with the given environment. |
| `ecs-run-app-migrations-container` | Creates and runs a migration task using the given container definition. |
| `rds-snapshot-app-db` | Creates a snapshot of the app database for the given environment. |
| `push-storybook-assets` | Pushes static build of Story Book to AWS S3 for hosting. |

## Development Scripts

These scripts are primarily used for developing the application and
application testing

| Script Name | Description |
| --- | --- |
| `check-docker-size` | Script to check the available disk space Docker has used |
| `ensure-application` | Ensure APPLICATION is set to `app` or `orders` and matches input value |
| `export-obfuscated-tspp-sample` | Export a subset of rows from the `transportation_service_provider_performances` table |
| `find-invoices` |  This script will use available API endpoints to find invoices in whatever environment you specify|
| `generate-devlocal-cert` | Convenience script for creating a new certificate signed by the DevLocal CA. |
| `go-find-pattern` |  searches over all our go source code files for a regex pattern |
| `merge-pr` |  A script to automate the landing of your GitHub pull requests. |
| `make-test` | A script to test common developer make targets. |
| `prime-api` | A script to connect to endpoints on the PRIME API. |
| `update-docker-compose` | Update branch name before running docker-compose |

### Building

This subset of development scripts is used primarily for building the app.

| Script Name | Description |
| --- | --- |
| `copy-swagger-ui` |  Copies the assets (other than xxx.html) into the public directory |
| `gen-assets` | generate assets from packages using go-bindata |
| `gen-server` | generate swagger code from yaml files |

### Testing

This subset of development scripts is used for testing

| Script Name | Description |
| --- | --- |
| `run-e2e-test` | Runs cypress tests with interactive GUI |
| `run-e2e-test-docker` | Runs cypress tests entirely inside docker containers like in CircleCI |
| `run-e2e-mtls-test-docker` | Runs integration tests for mtls endpoints inside docker containers like in CircleCI |
| `run-server-test` | Run golang server tests |
| `run-server-test-in-circle-container` | Executed in docker-compose.circle.yml to run the `make server_test` task in a CircleCI container |
| `run-storybook-tests` | Runs the storybook tests in docker containers for CircleCI testing. |
| `start-storybook-tests` | Intended to run inside a docker container, will wait for a 'storybook' container to have started before running the tests. |

### Secure Migrations

This subset of development scripts is used in developing secure
migrations.

| Script Name | Description |
| --- | --- |
| `download-secure-migration` |  A script to download secure migrations from all environments |
| `generate-secure-migration` |  A script to help manage the creation of secure migrations |
| `update-s3-sql-files` | A script to rename secure migration files to have .up.sql extension |
| `upload-secure-migration` | A script to upload secure migrations to all environments |

### Database Scripts

These scripts are primarily used for working with the database

| Script Name | Description |
| --- | --- |
| `db-backup` |  Backup the contents of the development database for later restore. |
| `db-cleanup` | Remove the database backup. |
| `db-restore` |  Restore the contents of the development database from an earlier backup. |
| `psql-dev` | Convenience script to drop into development postgres DB |
| `psql-deployed-migrations` | Convenience script to drop into deployed migrations postgres DB |
| `psql-schema` | Convenience script to dump the schema from the postgres DB |
| `psql-test` | Convenience script to drop into testing postgres DB |
| `psql-wrapper` | A wrapper around `psql` that sets correct values |
| `redis-dev` | Convenience script to drop into redis-cli |
| `update-migrations-manifest` | Update manifest for migrations |
| `wait-for-db` |  waits for an available database connection, or until a timeout is reached |
| `wait-for-db-docker` |  waits for an available database connection, or until a timeout is reached using docker |

### CAC Scripts

These scripts are primarily used for working with a CAC and the Orders API

| Script Name | Description |
| --- | --- |
| `cac-extract-cert` | Get a certificate from CAC |
| `cac-extract-fingerprint` | Get SHA 256 fingerprint of the public certificate from CAC |
| `cac-extract-pubkey` | Get a public key from CAC |
| `cac-extract-subject` | Get Subject from CAC |
| `cac-extract-token-label` | Get the Token Label from CAC |
| `cac-info` | Get general information from a CAC |
| `cac-prereqs` | Check the prereqs for CAC |

### Mutual TLS

These scripts are primarily for working with Mutual TLS certificates

| Script Name | Description |
| --- | --- |
| `mutual-tls-extract-fingerprint` | Get SHA 256 fingerprint of the public certificate from a cert file |
| `mutual-tls-extract-subject` | Get a sha256 hash of the certificate from a cert file|

### Amazon Console Scripts

These scripts are used for quickly opening up tools in the AWS Console

| Script Name | Description |
| --- | --- |
| `cloudwatch-logs` | Open up the CloudWatch logs group page |

### Vulnerability Scanning

These scripts are used to do vulnerability scanning on our code

| Script Name | Description |
| --- | --- |
| `anti-virus` | Scan the source code for viruses |
| `anti-virus-whitelists` | Create anti-virus whitelist database files for files and signatures to ignore |

### Prime Scripts

These scripts are primarily used for working with the Prime API

| Script Name | Description |
| --- | --- |
| `run-prime-docker` | Runs a docker container allowing access to the Prime API |


### Pricing & Rate Engine Scripts

These scripts are primarily for working with the pricing and rate engine

| Script Name | Description |
| --- | --- |
| `pricing-import` | Creates secure migration to move pricing data into production |