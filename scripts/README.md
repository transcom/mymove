# Scripts

This directory holds the scripts that are not compiled go code. For
compiled go code please look in the `bin/` directory of the project.

## Dev Environment

These scripts are primarily used for managing the developers
environment.

| Script Name | Description |
| --- | --- |
| `check-bash-version` | Script helps ensure that /etc/shells has all the correct entries in it |
| `check-go-version` | checks the go version required for the project |
| `check-gopath` | checks the go path is correct for the project |
| `check-hosts-file` | Script helps ensure that /etc/hosts has all the correct entries in it |
| `prereqs` | validate if all prerequisite programs have been installed |

## Operations Scripts

These scripts are used to operate the system.

| Script Name | Description |
| --- | --- |
| `download-alb-logs` | Download alb logs for the given environment and dates to a local folder |
| `scan-alb-logs` | Scan alb logs for specific http codes. |

## Pre-commit Scripts

These scripts are used primarily to check our code before
committing.

| Script Name | Description |
| --- | --- |
| `gen-docs-index` | generate index for documents |
| `generate-md-toc` |  Wrapper script to generate table of contents on Markdown files. |
| `pre-commit-circleci-validate` | validate CircleCI `config.yml` file |
| `pre-commit-go-imports` | modify imports in go files |
| `pre-commit-go-lint` | modify go files with linting rules |
| `pre-commit-go-mod` | modify `go.mod` and `go.sum` to match whats in the project |
| `pre-commit-go-vet` | analyze code with `go vet` |
| `pre-commit-spellcheck` | run spell checker against code |

## CircleCI Scripts

These scripts are primarily used for CircleCI workflows.

| Script Name | Description |
| --- | --- |
| `check-deployed-commit` |  checks that the deployed commit and given commit match. |
| `circleci-announce-broken-branch` | announce that a branch is broken |
| `circleci-push-dependency-updates` | Updates dependencies on the repo |
| `compare-deployed-commit` | checks that the given commit is ahead of the currently deployed commit |
| `do-exclusively` | CircleCI's current recommendation for roughly serializing a subset of build commands for a given branch |
| `ecs-deploy-service-container` |  Updates the named service with the given container definition template, image, and environment |
| `ecs-restart-services` | Restarted the ECS services associated with the given environment. |
| `ecs-run-app-migrations-container` | Creates and runs a migration task using the given container definition. |
| `rds-snapshot-app-db` | Creates a snapshot of the app database for the given environment. |
| `push-storybook-assets` | Pushes static build of Story Book to AWS S3 for hosting. |

## Development Scripts

These scripts are primarily used for developing the application and
application testing

| Script Name | Description |
| --- | --- |
| `export-obfuscated-tspp-sample` | Export a subset of rows from the `transportation_service_provider_performances` table |
| `find-invoices` |  This script will use available API endpoints to find invoices in whatever environment you specify|
| `generate-devlocal-cert` | Convenience script for creating a new certificate signed by the DevLocal CA. |
| `go-find-pattern` |  searches over all our go source code files for a regex pattern |
| `merge-pr` |  A script to automate the landing of your GitHub pull requests. |
| `make-test` | A script to test common developer make targets. |

### Building

This subset of development scripts is used primarily for building the app.

| Script Name | Description |
| --- | --- |
| `copy-swagger-ui` |  Copies the assets (other than xxx.html) into the public directory |
| `gen-server` | generate swagger code from yaml files |

### Testing

This subset of development scripts is used for testing

| Script Name | Description |
| --- | --- |
| `run-e2e-test` | Runs cypress tests with interactive GUI |
| `run-e2e-test-docker` | Runs cypress tests entirely inside docker containers like in CircleCI |
| `run-server-test-in-circle-container` | Executed in docker-compose.circle.yml to run the `make server_test` task in a CircleCI container |

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
| `psql-test` | Convenience script to drop into testing postgres DB |
| `psql-wrapper` | A wrapper around `psql` that sets correct values |
| `wait-for-db` |  waits for an available database connection, or until a timeout is reached |
| `wait-for-db-docker` |  waits for an available database connection, or until a timeout is reached using docker |

### CAC Scripts

These scripts are primarily used for working with a CAC and the Orders API

| Script Name | Description |
| --- | --- |
| `cac-extract-cert` | Get a certificate from CAC |
| `cac-extract-fingerprint` | Get SHA 256 fingerprint from CAC |
| `cac-extract-pubkey` | Get a public key from CAC |
| `cac-extract-subject` | Get Subject from CAC |
| `cac-extract-token-label` | Get the Token Label from CAC |
| `cac-info` | Get general information from a CAC |
| `cac-prereqs` | Check the prereqs for CAC |

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
