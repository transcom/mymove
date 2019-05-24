# Scripts

This directory holds the scripts that are not compiled go code. For
compiled go code please look in the `bin/` directory of the project.

## Dev Environment

These scripts are primarily used for managing the developers
environment.

| Script Name | Description |
| --- | --- |
| check-bash-version | Script helps ensure that /etc/shells has all the correct entries in it |
| check-go-version | checks the go version required for the project |
| check-gopath | checks the go path is correct for the project |
| check-hosts-file | Script helps ensure that /etc/hosts has all the correct entries in it |
| prereqs | validate if all prerequisite programs have been installed |

## Pre-commit Scripts

These scripts are used primarily to check our code before
committing.

| Script Name | Description |
| --- | --- |
| gen-docs-index | generate index for documents |
| generate-md-toc |  Wrapper script to generate table of contents on Markdown files. |
| pre-commit-circleci-validate | validate circleci `config.yml` file |
| pre-commit-go-imports | modify imports in go files |
| pre-commit-go-lint | modify go files with linting rules |
| pre-commit-go-mod | modify `go.mod` and `go.sum` to match whats in the project |
| pre-commit-go-vet | analyze code with `go vet` |
| pre-commit-spellcheck | run spell checker against code |

## CircleCI Scripts

These scripts are primarily used for CircleCI workflows.

| Script Name | Description |
| --- | --- |
| check-deployed-commit |  checks that the deployed commit and given commit match. |
| circleci-announce-broken-branch | announce that a branch is broken |
| circleci-push-dependency-updates | Updates dependencies on the repo |
| compare-deployed-commit | checks that the given commit is ahead of the currently deployed commit |
| do-exclusively | CircleCI's current recommendation for roughly serializing a subset of build commands for a given branch |
| ecs-deploy-service-container |  Updates the named service with the given container definition template, image, and environment |
| ecs-restart-services | Restarted the ECS services associated with the given environment. |
| ecs-run-app-migrations-container | Creates and runs a migration task using the given container definition. |
| rds-snapshot-app-db | Creates a snapshot of the app database for the given environment. |

## Bug Triaging Scripts

These scripts are used primarily for bug triaging in AWS.

| Script Name | Description |
| --- | --- |
| ecs-show-service-logs |  Show logs from the containers running for the named service. |
| ecs-show-service-stopped-logs |  Show logs from the most recently stopped app tasks. |

## Development Scripts

These scripts are primarily used for developing the application and
application testing

| Script Name | Description |
| --- | --- |
| dump-function-calls |  Show all used functions in our codebase. |
| dump-packages |  Show all used packges in our codebase. |
| export-obfuscated-tspp-sample | Export a subset of rows from the transportation_service_provider_performances table |
| find-invoices |  This script will use available API endpoints to find invoices in whatever environment you specify|
| generate-devlocal-cert | Convenience script for creating a new certificate signed by the devlocal CA. |
| go-find-pattern |  searchs over all our go source code files for a regex pattern |
| merge-pr |  A script to automate the landing of your GitHub pull requests. |
| make-test | A script to test common developer make targets. |

### Building

This subset of development scripts is used primarily for building the app.

| Script Name | Description |
| --- | --- |
| copy-swagger-ui |  Copies the assets (other than xxx.html) into the public directory |
| gen-model | generate models using soda |
| gen-server | generate swagger code from yaml files |

### Testing

This subset of development scripts is used for testing

| Script Name | Description |
| --- | --- |
| gen-e2e-migration | generate migrations for cypress |
| go-test | runs go test but with the correct DB port |
| run-e2e-test | Runs cypress tests with interactive GUI |
| run-e2e-test-docker | Runs cypress tests entirely inside docker containers like in CircleCI |

### Secure Migrations

This subset of development scripts is used in developing secure
migrations.

| Script Name | Description |
| --- | --- |
| apply-secure-migration.sh | Executes an SQL file from S3 against the environment's database. |
| generate-secure-migration |  A script to help manage the creation of secure migrations |
| download-secure-migration |  A script to download secure migrations from all environments |
| run-prod-migrations |  A script to apply all migrations, including secure migrations, to a local database. |
| upload-secure-migration | A script to upload secure migrations to all environments |

### Database Scripts

These scripts are primarily used for working with the database

| Script Name | Description |
| --- | --- |
| db-backup |  Backup the contents of the development database for later restore. |
| db-cleanup | Remove the database backup. |
| db-restore |  Restore the contents of the development database from an earlier backup. |
| psql-dev | Convenience script to drop into development postgres DB |
| psql-prod-migrations | Convenience script to drop into production migrations postgres DB |
| psql-test | Convenience script to drop into testing postgres DB |
| psql-wrapper | A wrapper around `psql` that sets correct values |
| wait-for-db |  waits for an available database connection, or until a timeout is reached |
| wait-for-db-docker |  waits for an available database connection, or until a timeout is reached using docker |
