# Scripts

This directory holds the scripts that are not compiled go code. For
compiled go code please look in the `bin/` directory of the project.

If you want to see if scripts are not listed in this file you can run
`find-scripts-missing-in-readme`.

## Dev Environment

These scripts are primarily used for managing the developers
environment.

| Script Name            | Description                                                                                                                                              |
| ---------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `check-changes`        | checks for changes since the last `git pull` using `git diff` for any file changes to a given path                                                       |
| `check-go-version`     | checks the go version required for the project                                                                                                           |
| `check-gopath`         | checks the go path is correct for the project                                                                                                            |
| `check-hosts-file`     | adds necessary entries to /etc/hosts                                                                                                                     |
| `check-node-version`   | checks the node version required for the project                                                                                                         |
| `kill-process-on-port` | asks to kill a process running on the specified port                                                                                                     |
| `prereqs`              | checks if all prerequisite programs have been installed                                                                                                  |
| `server-dev`           | Runs the MilMove app server in dev using `air`. Use the `--help` flag for more information. This is similar to `make server_run` but does not use `gin`. |
| `setup`                | installs all prerequisites and sets up the shell file                                                                                                    |

## AWS Scripts

These scripts are used for interacting with AWS or secrets in the AWS System Manager Parameter Store

| Script Name         | Description                                                |
| ------------------- | ---------------------------------------------------------- |
| `aws`               | Linked to aws-vault-wrapper. Runs the aws binary           |
| `aws-vault-wrapper` | A wrapper to ensure AWS credentials are in the environment |
| `chamber`           | Linked to aws-vault-wrapper. Runs chamber binary           |

## Operations Scripts

These scripts are used to operate the system.

You will need to specify which account you're using. Do so by pre-pending
`DISABLE_AWS_VAULT_WRAPPER=1 aws-vault exec AWS_ACCOUNT --` to the script command.
For example, to run the `health-tls-check` script, you'd run:

```bash
DISABLE_AWS_VAULT_WRAPPER=1 aws-vault exec transcom-gov-milmove-exp -- scripts/health-tls-check
```

| Script Name         | Description                                                             |
| ------------------- | ----------------------------------------------------------------------- |
| `download-alb-logs` | Download alb logs for the given environment and dates to a local folder |
| `dupe-secrets`      | Dupes experimental secrets to target params                             |
| `health-tls-check`  | Run health and TLS version checks.                                      |
| `scan-alb-logs`     | Scan alb logs for specific http codes.                                  |

## Deployment Scripts

This series of scripts allows you to quickly deploy the app manually (for example, should CircleCI be down).
Run the scripts in this order with the environment you're deploying to (`loadtest`, `demo`, `exp`, `stg`, `prd`) passed in as a variable.
You will need to specify which account you're using. Do so by pre-pending
`DISABLE_AWS_VAULT_WRAPPER=1 aws-vault exec AWS_ACCOUNT --` to the script command (i.e.
`DISABLE_AWS_VAULT_WRAPPER=1 aws-vault exec transcom-gov-milmove-exp -- scripts/deploy-app exp`).
For example, to run in the `exp` environment, you'd run:

```bash
DISABLE_AWS_VAULT_WRAPPER=1 aws-vault exec transcom-gov-milmove-exp -- scripts/deploy-app-migrations exp
DISABLE_AWS_VAULT_WRAPPER=1 aws-vault exec transcom-gov-milmove-exp -- scripts/deploy-app exp
DISABLE_AWS_VAULT_WRAPPER=1 aws-vault exec transcom-gov-milmove-exp -- scripts/deploy-app-client-tls exp
DISABLE_AWS_VAULT_WRAPPER=1 aws-vault exec transcom-gov-milmove-exp -- scripts/deploy-app-tasks exp
```

| Script Name             | Description               |
| ----------------------- | ------------------------- |
| `deploy-app-migrations` | Deploy the app migrations |
| `deploy-app`            | Deploy the app            |
| `deploy-app-client-tls` | Deploy the app client-tls |
| `deploy-app-tasks`      | Deploy the app tasks      |

## Pre-commit Scripts

These scripts are used primarily to check our code before
committing.

| Script Name                   | Description                                                       |
| ----------------------------- | ----------------------------------------------------------------- |
| `commit-msg`                  | Ensure JIRA issue is tagged to commit message                     |
| `gen-docs-index`              | generate index for documents                                      |
| `pre-commit-go-custom-linter` | run a custom linter against files (passed go files by pre-commit) |
| `pre-commit-go-imports`       | modify imports in go files                                        |
| `pre-commit-go-lint`          | modify go files with linting rules                                |
| `pre-commit-go-mod`           | modify `go.mod` and `go.sum` to match whats in the project        |
| `pre-commit-go-vet`           | analyze code with `go vet`                                        |
| `pre-commit-swagger-validate` | run `swagger validate`                                            |
| `lint-yaml-with-spectral`     | run `spectral` linter on external APIs                            |

## CircleCI Scripts

These scripts are primarily used for CircleCI workflows.

| Script Name                        | Description                                                                                             |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------- |
| `check-deployed-commit`            | checks that the deployed commit and given commit match.                                                 |
| `check-failed-test`            | used for gitlab to check failed xml test results and exit 1.
| `check-generated-code`             | checks that the generated code has not changed                                                          |
| `check-tls-pair`                   | checks that the TLS CERT and KEY match using openssl                                                    |
| `circleci-announce-broken-branch`  | announce that a branch is broken                                                                        |
| `compare-deployed-commit`          | checks that the given commit is ahead of the currently deployed commit                                  |
| `do-exclusively`                   | CircleCI's current recommendation for roughly serializing a subset of build commands for a given branch |
| `ecr-describe-image-scan-findings` | Checks an uploaded image scan results                                                                   |
| `ecs-deploy-service-container`     | Updates the named service with the given name, image, and environment.                                  |
| `ecs-deploy-task-container`        | Updates the named task with the given name, image, and environment.                                     |
| `ecs-restart-services`             | Restarted the ECS services associated with the given environment.                                       |
| `ecs-run-app-migrations-container` | Creates and runs a migration task using the given container definition.                                 |
| `rds-snapshot-app-db`              | Creates a snapshot of the app database for the given environment.                                       |
| `push-storybook-assets`            | Pushes static build of Story Book to AWS S3 for hosting.                                                |

## Development Scripts

These scripts are primarily used for developing the application and
application testing

| Script Name              | Description                                                                                       |
| ------------------------ | ------------------------------------------------------------------------------------------------- |
| `check-docker-size`      | Script to check the available disk space Docker has used                                          |
| `ds-for-gbloc`           | Helper script to find duty stations with the same GBLOC as the given transportation office        |
| `ensure-application`     | Ensure APPLICATION is set to `app` or `orders` and matches input value                            |
| `find-invoices`          | This script will use available API endpoints to find invoices in whatever environment you specify |
| `generate-devlocal-cert` | Convenience script for creating a new certificate signed by the DevLocal CA.                      |
| `go-find-pattern`        | searches over all our go source code files for a regex pattern                                    |
| `merge-pr`               | A script to automate the landing of your GitHub pull requests.                                    |
| `make-test`              | A script to test common developer make targets.                                                   |
| `handle-pr-comment`      | A script to post, update, or delete a comment on a pull request when test coverage checks change. |
| `prime-api`              | A script to connect to endpoints on the Prime API.                                                |
| `pricing-acceptance`     | A script to handle the acceptance process for pricing work.                                       |
| `to-for-gbloc`           | Helper script to find transportation offices with the same GBLOC as the given duty station        |
| `update-docker-compose`  | Update branch name before running docker-compose                                                  |
| `zips-for-gbloc`         | Helper script to find ZIP codes in the given GBLOC.                                               |

### Building

This subset of development scripts is used primarily for building the app.

| Script Name       | Description                                                       |
| ----------------- | ----------------------------------------------------------------- |
| `copy-swagger-ui` | Copies the assets (other than xxx.html) into the public directory |
| `copy-react-file-viewer` | Copies react-file-viewer assets into the public directory to support dynamic importing of the pdfjs-dist library |
| `rebuild-dependencies-without-binaries` | Creates binaries for installed dependencies that don't come with one |
| `fetch-react-file-viewer-from-yarn` | Fetches react file viewer version from yarn.lock, extracts the /dist/ folder, and stores the output into public/static/react-file-viewer in order to have the ESM chunk be served to the client during runtime |
| `gen-server`      | generate swagger code from yaml files                             |
| `openapi`         | invokes the openapi redoc swagger tool                            |

### Testing

This subset of development scripts is used for testing

| Script Name               | Description                                                            |
| ------------------------- | ---------------------------------------------------------------------- |
| `ensure-go-test-coverage` | Parse the go test coverage to ensure coverage increases                |
| `ensure-js-test-coverage` | Parse the js test coverage to ensure coverage increases                |
| `ensure-spectral-lint`    | Parse the spectral lint output to ensure error/warning counts decrease |
| `run-e2e-mtls-test`       | Runs playwright integration tests                                      |
| `run-e2e-test`            | Runs integration tests for mtls endpoints                              |
| `run-server-test`         | Run golang server tests                                                |

### Secure Migrations

This subset of development scripts is used in developing secure
migrations.

| Script Name                 | Description                                                                                  |
| --------------------------- | -------------------------------------------------------------------------------------------- |
| `download-secure-migration` | A script to download secure migrations from all environments                                 |
| `generate-secure-migration` | A script to help manage the creation of secure migrations                                    |
| `upload-secure-migration`   | A script to upload secure migrations to all environments in both commercial and GovCloud AWS |
| `generate-ddl-migration`   | A script to help manage the creation of DDL migrations                                      |

### Database Scripts

These scripts are primarily used for working with the database

| Script Name                  | Description                                                                            |
| ---------------------------- | -------------------------------------------------------------------------------------- |
| `db-backup`                  | Backup the contents of the development database for later restore.                     |
| `db-cleanup`                 | Remove the database backup.                                                            |
| `db-truncate`                | Truncates the configured database. Used in testing.                                    |
| `db-restore`                 | Restore the contents of the development database from an earlier backup.               |
| `psql-dev`                   | Convenience script to drop into development postgres DB                                |
| `psql-deployed-migrations`   | Convenience script to drop into deployed migrations postgres DB                        |
| `psql-schema`                | Convenience script to dump the schema from the postgres DB                             |
| `psql-test`                  | Convenience script to drop into testing postgres DB                                    |
| `psql-wrapper`               | A wrapper around `psql` that sets correct values                                       |
| `redis-dev`                  | Convenience script to drop into redis-cli                                              |
| `update-migrations-manifest` | Update manifest for migrations                                                         |
| `wait-for-db`                | waits for an available database connection, or until a timeout is reached              |
| `wait-for-db-docker`         | waits for an available database connection, or until a timeout is reached using docker |

### CAC Scripts

These scripts are primarily used for working with a CAC and the Orders API

| Script Name               | Description                                                |
| ------------------------- | ---------------------------------------------------------- |
| `cac-extract-cert`        | Get a certificate from CAC                                 |
| `cac-extract-fingerprint` | Get SHA 256 fingerprint of the public certificate from CAC |
| `cac-extract-pubkey`      | Get a public key from CAC                                  |
| `cac-extract-subject`     | Get Subject from CAC                                       |
| `cac-extract-token-label` | Get the Token Label from CAC                               |
| `cac-info`                | Get general information from a CAC                         |
| `cac-prereqs`             | Check the prereqs for CAC                                  |

### Mutual TLS

These scripts are primarily for working with Mutual TLS certificates

| Script Name                      | Description                                                           |
| -------------------------------- | --------------------------------------------------------------------- |
| `generate-p7b-file`              | Creates a concatenated p7b file from several certificate file formats |
| `mutual-tls-extract-fingerprint` | Get SHA 256 fingerprint of the public certificate from a cert file    |
| `mutual-tls-extract-subject`     | Get a sha256 hash of the certificate from a cert file                 |

### Amazon Console Scripts

These scripts are used for quickly opening up tools in the AWS Console

| Script Name       | Description                            |
| ----------------- | -------------------------------------- |
| `cloudwatch-logs` | Open up the CloudWatch logs group page |

### Vulnerability Scanning

These scripts are used to do vulnerability scanning on our code

| Script Name             | Description                                                                   |
| ----------------------- | ----------------------------------------------------------------------------- |
| `anti-virus`            | Scan the source code for viruses                                              |
| `anti-virus-whitelists` | Create anti-virus whitelist database files for files and signatures to ignore |

### Prime Scripts

These scripts are primarily used for working with the Prime API

| Script Name        | Description                                              |
| ------------------ | -------------------------------------------------------- |
| `run-prime-docker` | Runs a docker container allowing access to the Prime API |

### Pricing & Rate Engine Scripts

These scripts are primarily for working with the pricing and rate engine

| Script Name      | Description                                                   |
| ---------------- | ------------------------------------------------------------- |
| `pricing-import` | Creates secure migration to move pricing data into production |
