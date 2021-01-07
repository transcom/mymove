# Personal Property Prototype

[![CircleCI](https://circleci.com/gh/transcom/mymove/tree/master.svg?style=shield&circle-token=8782cc55afd824ba48e89fc4e49b466c5e2ce7b1)](https://circleci.com/gh/transcom/mymove/tree/master)

[![GoDoc](https://godoc.org/github.com/transcom/mymove?status.svg)](https://godoc.org/github.com/transcom/mymove)

This repository contains the application source code for the Personal Property Prototype, a possible next generation version of the Defense Personal Property System (DPS). DPS is an online system managed by the U.S. [Department of Defense](https://www.defense.gov/) (DoD) [Transportation Command](http://www.ustranscom.mil/) (USTRANSCOM) and is used by service members and their families to manage household goods moves.

This prototype was built by a [Defense Digital Service](https://www.dds.mil/) team in support of USTRANSCOM's mission.

## License Information

Works created by U.S. Federal employees as part of their jobs typically are not eligible for copyright in the United
States. In places where the contributions of U.S. Federal employees are not eligible for copyright, this work is in
the public domain. In places where it is eligible for copyright, such as some foreign jurisdictions, the remainder of
this work is licensed under [the MIT License](https://opensource.org/licenses/MIT), the full text of which is included
in the [LICENSE.txt](./LICENSE.txt) file in this repository.

## Table of Contents

<!-- Table of Contents auto-generated with `scripts/generate-md-toc` -->

<!-- toc -->

* [Overview](#overview)
* [Supported Browsers](#supported-browsers)
* [Login.gov](#logingov)
* [Application Setup](#application-setup)
  * [Setup: Developer Setup](#setup-developer-setup)
  * [Setup: Git](#setup-git)
  * [Setup: Golang](#setup-golang)
  * [Setup: Project Checkout](#setup-project-checkout)
  * [Setup: Project Layout](#setup-project-layout)
  * [Setup: Editor Config](#setup-editor-config)
  * [Setup: Makefile](#setup-makefile)
  * [Setup: Quick Initial Setup](#setup-quick-initial-setup)
  * [Setup: Direnv](#setup-direnv)
    * [Helpful variables for `.envrc.local`](#helpful-variables-for-envrclocal)
  * [Setup: Prerequisites](#setup-prerequisites)
  * [Setup: Pre-Commit](#setup-pre-commit)
    * [Troubleshooting install issues (process hanging on install hooks)](#troubleshooting-install-issues-process-hanging-on-install-hooks)
  * [Setup: Dependencies](#setup-dependencies)
  * [Setup: Build Tools](#setup-build-tools)
  * [Setup: Database](#setup-database)
  * [Setup: Server](#setup-server)
  * [Setup: MilMoveLocal Client](#setup-milmovelocal-client)
  * [Setup: OfficeLocal client](#setup-officelocal-client)
  * [Setup: AdminLocal client](#setup-adminlocal-client)
  * [Setup: DPS user](#setup-dps-user)
  * [Setup: Orders Gateway](#setup-orders-gateway)
  * [Setup: Prime API](#setup-prime-api)
  * [Setup: AWS Services (Optional)](#setup-aws-services-optional)
* [Development](#development)
  * [TSP Award Queue](#tsp-award-queue)
  * [Test Data Generator](#test-data-generator)
  * [API / Swagger](#api--swagger)
  * [Testing](#testing)
    * [Troubleshooting tips -- integration / e2e tests](#troubleshooting-tips----integration--e2e-tests)
  * [Logging](#logging)
    * [Log files](#log-files)
  * [Database](#database)
    * [Dev DB Commands](#dev-db-commands)
    * [Test DB Commands](#test-db-commands)
    * [Migrations](#migrations)
  * [Environment Variables](#environment-variables)
  * [Documentation](#documentation)
  * [Spellcheck](#spellcheck)
    * [Tips for staying sane](#tips-for-staying-sane)
  * [GoLand](#goland)
  * [Troubleshooting](#troubleshooting)
    * [Postgres Issues](#postgres-issues)
    * [Development Machine Timezone Issues](#development-machine-timezone-issues)
    * [Linters & Pre-commit Hooks](#linters--pre-commit-hooks)
    * [Yarn install markdown-spell (aka mdspell)](#yarn-install-markdown-spell-aka-mdspell)
  * [Manual Redeploys and Other Helpful Information in an Emergency](#manual-redeploys-and-other-helpful-information-in-an-emergency)
  * [PII Best Practices](#pii-best-practices)
    * [More about content dispositions](#more-about-content-dispositions)
    * [More about browser settings](#more-about-browser-settings)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

## Overview

Please check the [Milmove Project Wiki](https://github.com/transcom/mymove/wiki) for details on the project itself.

## Supported Browsers

As of 3/6/2018, DDS has confirmed that support for IE is limited to IE 11 and Edge or newer versions. Currently, the intention is to encourage using Chrome and Firefox instead, with specific versions TBD. Research is incomplete on mobile browsers, but we are assuming support for iOS and Android. For more information please read [ADR0016 Browser Support](./docs/adr/0016-Browser-Support.md).

## Login.gov

You'll need accounts for login.gov and the login.gov sandbox.  These will require two-factor authentication, so have your second factor (one of: phone, authentication app, security key, CAC) on hand.
To create an account at login.gov, use your regular `truss.works` email and follow [the official instructions](https://login.gov/help/creating-an-account/how-to-create-an-account/).
To create an account in the sandbox, follow the same instructions, but [in the sandbox server](https://idp.int.identitysandbox.gov/sign_up/enter_email).  Do _not_ use your regular email address in the sandbox.
**Tip**: You can use the plus sign `+` to create a new truss email address.  `name+some_string@truss.works` will be treated as a new address, but will be routed to `name@truss.works`.

## Application Setup

### Setup: Developer Setup

Note: These instructions are a living document and often fall out-of-date. If you run into anything that needs correcting or updating, please create a PR with those changes to help those coming after you.

There are a number of things you'll need at a minimum to be able to check out, develop and run this project.

* Install [Homebrew](https://brew.sh)
  * Use the following command `/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
* We normally use the latest version of Go unless there's a known conflict (which will be announced by the team) or if we're in the time period just after a new version has been released.
  * Install it with Homebrew: `brew install go`
  * **Note**: If you have previously modified your PATH to point to a specific version of go, make sure to remove that. This would be either in your `.bash_profile` or `.bashrc`, and might look something like `PATH=$PATH:/usr/local/opt/go@1.12/bin`.
* Ensure you are using the latest version of bash for this project:
  * Install it with Homebrew: `brew install bash`
  * Update list of shells that users can choose from: `[[ $(cat /etc/shells | grep /usr/local/bin/bash) ]] || echo "/usr/local/bin/bash" | sudo tee -a /etc/shells`
  * If you are using bash as your shell (and not zsh, fish, etc) and want to use the latest shell as well then change it (optional): `chsh -s /usr/local/bin/bash`
  * Ensure that `/usr/local/bin` comes before `/bin` on your `$PATH` by running `echo $PATH`. Modify your path by editing `~/.bashrc` or `~/.bash_profile` and changing the `PATH`.  Then source your profile with `source ~/.bashrc` or `~/.bash_profile` to ensure that your terminal has it.

### Setup: Git

Use your work email when making commits to our repositories. The simplest path to correctness is setting global config:

  ```bash
  git config --global user.email "trussel@truss.works"
  git config --global user.name "Trusty Trussel"
  ```

If you drop the `--global` flag these settings will only apply to the current repo. If you ever re-clone that repo or clone another repo, you will need to remember to set the local config again. You won't. Use the global config. :-)

For web-based Git operations, GitHub will use your primary email unless you choose "Keep my email address private". If you don't want to set your work address as primary, please [turn on the privacy setting](https://github.com/settings/emails).

Note that with 2-factor-authentication enabled, in order to push local code to GitHub through HTTPS, you need to [create a personal access token](https://gist.github.com/ateucher/4634038875263d10fb4817e5ad3d332f) and use that as your password.

### Setup: Golang

All of Go's tooling expects Go code to be checked out in a specific location. Please read about [Go workspaces](https://golang.org/doc/code.html#Workspaces) for a full explanation. If you just want to get started, then decide where you want all your go code to live and configure the GOPATH environment variable accordingly. For example, if you want your go code to live at `~/code/go`, you should add the following like to your `.bash_profile`:

  ```bash
  export GOPATH=~/code/go
  ```

Golang expect the `GOPATH` environment variable to be defined.  If you'd like to use the default location, then add the following to your `.bash_profile` or hardcode the default value.  This line will set the GOPATH environment variable to the value of `go env GOPATH` if it is not already set.

  ```bash
  export GOPATH=${GOPATH:-$(go env GOPATH)}
  ```

**Regardless of where your go code is located**, you need to add `$GOPATH/bin` to your `PATH` so that executables installed with the go tooling can be found. Add the following to your `.bash_profile`:

  ```bash
  export PATH=$(go env GOPATH)/bin:$PATH
  ```

Finally to have these changes applied to your shell you must `source` your profile:

  ```bash
  source ~/.bash_profile
  ```

You can confirm that the values exist with:

  ```bash
  env | grep GOPATH
  # Verify the GOPATH is correct
  env | grep PATH
  # Verify the PATH includes your GOPATH bin directory
  ```

### Setup: Project Checkout

You can checkout this repository by running `git clone git@github.com:transcom/mymove.git`. Please check out the code in a directory like `~/Projects/mymove` and NOT in your `$GOPATH`. As an example:

  ```bash
  mkdir -p ~/Projects
  git clone git@github.com:transcom/mymove.git
  cd mymove
  ```

You will then find the code at `~/Projects/mymove`. You can check the code out anywhere EXCEPT inside your `$GOPATH`. So this is customization that is up to you.

### Setup: Project Layout

All of our code is intermingled in the top level directory of mymove. Here is an explanation of what some of these directories contain:

* `.circleci`: Directory for CircleCI CI/CD configuration
* `bin`: A location for tools compiled from the `cmd` directory
* `build`: The build output directory for the client. This is what the development server serves
* `cmd`: The location of main packages for any go binaries we build
* `config`: Config files for the database and AWS ECS. Also certificates.
* `cypress`: The integration test files for the [Cypress tool](https://www.cypress.io/)
* `docs`: A location for docs for the project. This is where ADRs are
* `internal`: Generated code for duty station loader
* `migrations`: Database migrations, see [./migrations/README.md]
* `node_modules`: Cached javascript dependencies for the client
* `pkg`: The location of all of our go code for the server and various tools
* `public`: The client's static resources
* `scripts`: A location for tools helpful for developing this project
* `src`: The react source code for the client
* `swagger`: The swagger definition files for each of our APIs

### Setup: Editor Config

[EditorConfig](http://editorconfig.org/) allows us to manage editor configuration (like indent sizes,) with a [file](https://github.com/transcom/ppp/blob/master/.editorconfig) in the repo. Install the appropriate plugin in your editor to take advantage of that if you wish.

### Setup: Makefile

The primary way to interact with the project is via the `Makefile`. The `Makefile` contains a number of handy
targets (you can think of these as commands) that make interacting with the project easier. Each target manages
its own dependencies so that you don't have to. This is how you'll do common tasks like build the project, run
the server and client, and manage the database.

The fastest way to get familiar with the `Makefile` is to use the command `make help`. You can also type `make`
and it will default to calling `make help` target on your behalf.  The `Makefile` is important to this project
so take the time to understand what it does.

### Setup: Quick Initial Setup

The following commands will get mymove running on your machine for the first time. This is an abbreviated list that should get you started. Please read below for explanations of each of the commands.

1. `direnv allow`
1. `make prereqs`
1. `make ensure_pre_commit`
1. `make deps`
1. `make build_tools`
1. `make db_dev_run`
1. `make db_dev_migrate`
1. `make server_run`
1. `make client_build`
1. `make client_run`

### Setup: Direnv

For managing local environment variables, we're using [direnv](https://direnv.net/). You need to [configure your shell to use it](https://direnv.net/). For bash, add the command `eval "$(direnv hook bash)"` to whichever file loads upon opening bash (likely `~/.bash_profile`, though instructions say `~/.bashrc`). For zsh, add `eval "$(direnv hook zsh)"` to `~/.zshrc`.

Run `direnv allow` to load up the `.envrc` file. It should complain that you have missing variables which you will rectify in one of the following ways.

The first and recommended way is to use the [chamber tool](https://github.com/segmentio/chamber) to read secrets from AWS vault. We suggest installing chamber with `brew install chamber`. To use the AWS vault you will need to follow the [instructions to set up AWS first](#setup-aws-services-optional).

Once you've installed, run `cp .envrc.chamber.template .envrc.chamber` to enable getting secret values from `chamber`. **Note** that this method does not work for users of the `fish` shell unless you replace `direnv allow` with `direnv export fish | source`. **Note also** if you have a very poor internet connection, this method may be problematic to you.

The alternative is to add a `.envrc.local` file. Then run `DISABLE_AWS_VAULT_WRAPPER=1 AWS_REGION=us-gov-west-1 aws-vault exec transcom-gov-dev -- chamber env app-devlocal >> .envrc.local`. If you don't have access to chamber, you can also `touch .envrc.local` and add any values that the output from direnv asks you to define. Instructions are in the error messages.

#### Helpful variables for `.envrc.local`

* `export GOLANGCI_LINT_CONCURRENCY=8` - variable to increase concurrency of golangci-lint; defaults to 6 on dev machines and to 1 in CircleCI.
* `export GOLAND=1` - variable to enable go code debugging in goland

### Setup: Prerequisites

Run `make prereqs` and install everything it tells you to. Most of the prerequisites need you to use `brew install <package>`.

**NOTE:** Do not configure PostgreSQL to automatically start at boot time or the DB commands will not work correctly!

### Setup: Pre-Commit

Run `pre-commit install` to install a pre-commit hook into `./git/hooks/pre-commit`.  This is different than `brew install pre-commit` and must be done so that the hook will check files you are about to commit to the repository.  Next install the pre-commit hook libraries with `pre-commit install-hooks`.

You can feel free to skip running the pre-commit checks at this time. Before you do run `pre-commit run -a`, you will need to install Javascript dependencies and generate some golang code from Swagger files. An easier way to handle this is by running `make pre_commit_tests` or `make server_generate client_deps && pre-commit run -a`.

#### Troubleshooting install issues (process hanging on install hooks)

Since pre-commit uses node to hook things up in both your local repo and its cache folder (located at `~/.cache/pre-commit`),it requires a global node install. If you are using nodenv to manage multiple installed nodes, you'll need to set a global version to proceed (eg `nodenv global 12.16.3`). You can find the current supported node version [here (in `.node-version`)](./.node-version). Make sure you run `nodenv install` to install the current supported version.

### Setup: Dependencies

Run `make deps`. This will check your system for any setup issues. Then it will ensure that you have installed pre-commit
and go on to install the client (javascript) and server (golang) dependencies for you.

### Setup: Build Tools

Run `make build_tools` to get all the server and tool dependencies built. These will be needed in future steps to not only generate test data but also to interact with the database and more.

### Setup: Database

You will need to setup a local database before you can begin working on the local server / client. Docker will need to be running for any of this to work.

1. `make db_dev_run` and `make db_test_run`: Creates a PostgreSQL docker container for dev and test, if they don't already exist.

1. `make db_dev_migrate` and `make db_test_migrate`:  Runs all existing database migrations for dev and test databases, which does things like creating table structures, etc. You will run this command again anytime you add new migrations to the app (see below for more)

You can validate that your dev database is running by running `psql-dev`. This puts you in a PostgreSQL shell. Type `\dt` to show all tables, and `\q` to quit.
You can validate that your test database is running by running `psql-test`. This puts you in a PostgreSQL shell. Type `\dt` to show all tables, and `\q` to quit.

If you are stuck on this step you may need to see the section on Troubleshooting.

### Setup: Server

1. `make server_run`: installs dependencies, then builds and runs the server using `gin`, which is a hot reloading go server. It will listen on 8080 and will rebuild the actual server any time a go file changes. Pair this with `make client_run` to have hot reloading of the entire application.

In rare cases, you may want to run the server standalone, in which case you can run `make server_run_standalone`. This will build both the client and the server and this invocation can be relied upon to be serving the client JS on its own rather than relying on webpack doing so as when you run `make client_run`. You can run this without running `make client_run` and the whole app should work.

Dependencies are managed by [go modules](https://github.com/golang/go/wiki/Modules). New dependencies are automatically detected in import statements and added to `go.mod` when you run `go build` or `go run`. You can also manually edit `go.mod` as needed.

If you need to add a Go-based tool dependency that is otherwise not imported by our code, import it in `pkg/tools/tools.go`.

After importing _any_ go dependency it's a good practice to run `go mod tidy`, which prunes unused dependencies and calculates dependency requirements for all possible system architectures.

### Setup: MilMoveLocal Client

1. `make client_build` (if setting up for first time)
2. `make client_run`

The above will start the webpack dev server, serving the front-end on port 3000. If paired with `make server_run` then the whole app will work, the webpack dev server proxies all API calls through to the server.

If both the server and client are running, you should be able to view the Swagger UI at <http://milmovelocal:3000/swagger-ui/internal.html>.  If it does not, try running `make client_build` (this only needs to be run the first time).

Dependencies are managed by yarn. To add a new dependency, use `yarn add`

### Setup: OfficeLocal client

1. Ensure that you have a test account which can log into the office site:
    * run `make db_dev_e2e_populate` to load test data
2. Run `make office_client_run`
    * Log into "Local Sign In" and either select a pre-made user or use the button to create a new user

### Setup: AdminLocal client

1. `make admin_client_run`

### Setup: DPS user

1. Ensure that you have a login.gov test account
2. Log into [MilMove Devlocal Auth](http://milmovelocal:3000/devlocal-auth/login) and create a new DPS user from the interface.

### Setup: Orders Gateway

Nothing to do.

### Setup: Prime API

The API that the Prime will use is authenticated via mutual TSL so there are a few things you need to do to interact with it in a local environment.

1. Make sure that the `primelocal` alias is setup for localhost - this should have been completed in the [Setup:Prerequisites](#Setup-Prerequisites) (check your `/etc/hosts` file for an entry for `primelocal`).
2. run `make server_run`
3. Access the Prime API using the devlocal-mtls certs. There is a script that shows you how to do this with curl at `./scripts/prime-api`. For instance to call the `move-task-orders` endpoint, call `./scripts/prime-api move-task-orders`



### Setup: AWS Services (Optional)

If you want to develop against AWS services you will need an AWS user account with `engineering` privileges. You will also need to follow these steps when using Chamber and AWS vault to [set up direnv](#setup-direnv).

AWS credentials are managed via `aws-vault`. Once you have received AWS credentials (which are provided by the infrastructure team), you can follow these instructions to [finish setting up AWS](https://github.com/transcom/transcom-infrasec-gov/blob/master/docs/runbook/0001-aws-organization-authentication.md).

## Development

### TSP Award Queue

This background job is built as a separate binary which can be built using `make tsp_run`.

### Test Data Generator

When creating new features, it is helpful to have sample data for the feature to interact with. The TSP Award Queue is an example of that--it matches shipments to TSPs, and it's hard to tell if it's working without some shipments and TSPs in the database!

* `make bin/generate-test-data` will build the fake data generator binary
* `bin/generate-test-data --named-scenario="dev_seed"` will populate the development database with a handful of users in various stages of progress along the flow. The emails are named accordingly (see [`devseed.go`](https://github.com/transcom/mymove/blob/master/pkg/testdatagen/scenario/devseed.go)). Alternatively, run `make db_dev_e2e_populate` to reset your db and populate it.
* `bin/generate-test-data` will run binary and create a preconfigured set of test data. To determine the data scenario you'd like to use, check out scenarios in the `testdatagen` package. Each scenario contains a description of what data will be created when the scenario is run. Pass the scenario in as a flag to the generate-test-data function. A sample command: `./bin/generate-test-data --scenario=2`.

There is also a package (`/pkg/testdatagen`) that can be imported to create arbitrary test data. This could be used in tests, so as not to duplicate functionality.

Currently, scenarios have the following numbers:

* `--scenario=1` for Award Queue Scenario 1
* `--scenario=2` for Award Queue Scenario 2
* `--scenario=3` for Duty Station Scenario
* `--scenario=4` for PPM or PPM SIT Estimate Scenario (can also use Rate Engine Scenarios for Estimates)
* `--scenario=5` for Rate Engine Scenario 1
* `--scenario=6` for Rate Engine Scenario 2
* `--scenario=7` for TSP test data

### API / Swagger

Internal services (i.e. endpoints only intended for use by the React client) are defined in `swagger/internal.yaml` and served at `/internal/swagger.yaml`. These are, as the name suggests, internal endpoints and not intended for use by external clients.

The Orders Gateway's API is defined in the file `swagger/orders.yaml` and served at `/orders/v0/orders.yaml`.

The Admin API is defined in the file `swagger/admin.yaml` and served at `/admin/v1/swagger.yaml`.

You can view the documentation for the following APIs (powered by Swagger UI) at the following URLS with a local client and server running:

* internal API: <http://milmovelocal:3000/swagger-ui/internal.html>

* admin API: <http://milmovelocal:3000/swagger-ui/admin.html>

* GHC API: <http://milmovelocal:3000/swagger-ui/ghc.html>

### Testing

There are a few handy targets in the Makefile to help you run tests:

* `make client_test`: Run front-end testing suites.
* `make server_test`: Run back-end testing suites. [Additional info for running go tests](https://github.com/transcom/mymove/wiki/run-go-tests)
* `make e2e_test`: Run e2e testing suite.
  * Note: this will not necessarily reflect the same results as in the CI environment, run with caution. One of the reasons for this is it's pulling actual cypress latest, which as of this writing is `5.0.0`. Another reason is your `.envrc` is going to populate your dev environment with a bunch of values that `make e2e_test_docker` won't have.
  * Note also: this runs with a full clean/rebuild, so it is not great for fast iteration. Use `yarn test:e2e` when working with individual tests
* `yarn test:e2e`: Open the cypress test runner against your already running servers and inspect/run individual e2e tests. (Should better reflect CI environment than above, but not as well as below.)
  * Note: You must already have the servers running for this to work! This may not reflect the same results as CI for the same reason as the above re: `.envrc` values. However, it is __significantly__ faster because you can run individual tests and not have to deal with the clean/rebuild.
* `yarn test:e2e-clean`: Resets your dev DB to a clean state before opening the Cypress test runner.
* `make e2e_test_docker`: Run e2e testing suite in the same docker container as is run in CircleCI.
  * Note: this also runs with a full clean/rebuild, so it is not great for fast iteration. Use `yarn test:e2e` when working with individual tests.
* `make test`: Run e2e, client- and server-side testing suites.

#### Troubleshooting tips -- integration / e2e tests

When running locally, you may find that retries or successive runs have unexpected failures. Some of the integration tests are written with the assumption that they will only be run against a clean DB. If you're working with one of these and don't have time to fix them to properly set up and clean up their state, you can use this command to reset your local dev db before opening the test runner. Note that if you choose not to fix the offending test(s), you'll have to repeatedly close the test runner to re-clean the the DB. You won't be able to take advantage of Cypress's hot reloading!

If you suspect memory issues, you can further inspect this with the commands:

* `yarn test:e2e-debug`, which runs `yarn test:e2e` with DEBUG stats
* `yarn test:e2e-debug-clean`, which runs `yarn test:e2e-clean` with DEBUG stats

### Logging

We are using [zap](https://github.com/uber-go/zap) as a logger in this project. We currently rely on its built-in `NewDevelopment()` and `NewProduction()` default configs, which are enabled in any of the executable packages that live in `cmd`.

This means that logging *is not* set up from within models or other packages unless the files in `cmd` are also being loaded. *If you attempt to call `zap.L()` or `zap.S()` without a configured logger, nothing will appear on the screen.*

If you need to see some output during the development process (say, for debugging purposes), it is best to use the standard lib `fmt` package to print to the screen. You will also need to pass `-v` to `go test` so that it prints all output, even from passing tests. The simplest way to do this is to run `go test` yourself by passing it which files to run, e.g. `go test pkg/models/* -v`.

#### Log files

In development mode, logs from the `milmove` process are written to `logs/dev.log`.

### Database

* Read [Querying the Database Safely](https://github.com/transcom/mymove/blob/master/docs/backend.md#querying-the-database-safely) to prevent SQL injections! *

A few commands exist for starting and stopping the DB docker container:

* `make db_run`: Starts the DB docker container if one doesn't already exist
* `make db_destroy`: Stops and removes the DB docker container

#### Dev DB Commands

There are a few handy targets in the Makefile to help you interact with the dev database:

* `make db_dev_run`: Initializes a new database if it does not exist and runs it, or starts the previously initialized Docker container if it has been stopped.
* `make db_dev_create`: Waits to connect to the DB and will create a DB if one doesn't already exist (run usually as part of `db_dev_run`).
* `make db_dev_reset`: Destroys your database container. Useful if you want to start from scratch.
* `make db_dev_migrate`: Applies database migrations against your running database container.
* `make db_dev_migrate_standalone`: Applies database migrations against your running database container but will not check for server dependencies first.
* `make db_dev_e2e_populate`: Populate data with data used to run e2e tests

#### Test DB Commands

The Dev Commands are used to talk to the dev DB.  If you were working with the test DB you would use these commands:

* `make db_test_run`
* `make db_test_create`
* `make db_test_reset`
* `make db_test_migrate`
* `make db_test_migrate_standalone`
* `make db_test_e2e_populate`
* `make db_test_e2e_backup`
* `make db_test_e2e_restore`
* `make db_test_e2e_cleanup`

The test DB commands all talk to the DB over localhost.  But in a docker-only environment (like CircleCI) you may not be able to use those commands, which is why `*_docker` versions exist for all of them:

* `make db_test_run_docker`
* `make db_test_create_docker`
* `make db_test_reset_docker`
* `make db_test_migrate_docker`

#### Migrations

To add new regular and/or secure migrations, see the [database development guide](https://github.com/transcom/mymove/wiki/migrate-the-database)

Running migrations in local development:

Use `make db_dev_migrate` to run migrations against your local dev environment.

Running migrations on Staging / Production:

Migrations are run automatically by CircleCI as part of the standard deploy process.

1. CircleCI builds and registers a container.
1. CircleCI deploys this container to ECS and runs it as a one-off 'task'.
1. The container downloads and execute migrations against the environment's database.
1. If migrations fail, CircleCI fails the deploy.
1. If migrations pass, CircleCI continues with the deploy.

### Environment Variables

In development, we use [direnv](https://direnv.net/) to setup environment variables required by the application.

* If you want to add a new environment variable to affect only your development machine, export it in `.envrc.local`. Variables exported in this file take precedence over those in `.envrc`.
* If you want to add a new environment variable that is required by new development, it can be added to `.envrc` using one of the following:

    ```bash
    # Add a default value for all devs that can be overridden in their .envrc.local
    export NEW_ENV_VAR="default value"

    # or

    # Specify that an environment variable must be defined in .envrc.local
    require NEW_ENV_VAR "Look for info on this value in chamber and Google Drive"
    ```

Required variables should be placed in google docs and linked in `.envrc`. The value should also be placed in `chamber`
with `DISABLE_AWS_VAULT_WRAPPER=1 AWS_REGION=us-gov-west-1 aws-vault exec transcom-gov-dev -- chamber write app-devlocal <key> <value>`. For long blocks of text like certificates you can write them with
`echo "$LONG_VALUE" | DISABLE_AWS_VAULT_WRAPPER=1 AWS_REGION=us-gov-west-1 aws-vault exec transcom-gov-dev -- chamber write app-devlocal <key> -`.

For per-tier environment variables (that are not secret), simply add the variables to the relevant `config/env/[experimental|staging|prod].env` file with the format `NAME=VALUE` on each line.  Then add the relevant section to `config/app.container-definition.json`.  The deploy process uses Go's [template package](https://golang.org/pkg/text/template/) for rendering the container definition.  For example,

```bash
MY_SPECIAL_TOKEN=abcxyz
```

```json
{
  "name": "MY_SPECIAL_TOKEN",
  "value": "{{ .MY_SPECIAL_TOKEN }}"
},
```

### Documentation

You can view the project's GoDoc on [godoc.org](https://godoc.org/github.com/transcom/mymove).

Alternatively, run the documentation locally using:

```shell
# within the project's root dir
$ godoc -http=:6060
```

Then visit <http://localhost:6060/pkg/github.com/transcom/mymove/>

### Spellcheck

We use [markdown-spellcheck](https://github.com/lukeapage/node-markdown-spellcheck) as a pre-commit hook to catch spelling errors in Markdown files. To make fixing caught errors easier, there's a handy make target that runs the spellchecker in interactive mode:

* `make spellcheck`

This will let you walk through the caught spelling errors one-by-one and choose whether to fix it, add it to the dictionary, or have it be permanently ignored for that file.

#### Tips for staying sane

* If you want to use a bare hyperlink, wrap it in angle braces: `<http://example.com>`

### GoLand

* GoLand supports [attaching the debugger to a running process](https://blog.jetbrains.com/go/2019/02/06/debugging-with-goland-getting-started/#debugging-a-running-application-on-the-local-machine), however this requires that the server has been built with specific flags. If you wish to use this feature in development add the following line `export GOLAND=1` to your `.envrc.local`. Once the server starts follow the steps outlined in the article above and you should now be able to set breakpoints using the GoLand debugger.

### Troubleshooting

* Random problems may arise if you have old Docker containers running. Run `docker ps` and if you see containers unrelated to our app, consider stopping them.
* If you happen to have installed pre-commit in a virtual environment not with brew, running `make prereqs` will not alert you. You may run into issues when running `make deps`. To install pre-commit: `brew install pre-commit`.
* If you're having trouble accessing the API docs or the server is otherwise misbehaving, try stopping the server, running `make client_build`, and then running `make client_run` and `make server_run`.

#### Postgres Issues

If you have problems connecting to PostgreSQL, or running related scripts, make sure you aren't already running a PostgreSQL daemon. You may see errors like:

```text
Migrator: problem creating schema migrations: couldn't start a new transaction: could not create new transaction: pq: role "postgres" does not exist
```

or

```text
Migrator: problem creating schema migrations: couldn't start a new transaction: could not create new transaction: pq: database "dev_db" does not exist
```

You can check this by typing `ps aux | grep postgres` or `brew services list` and looking for existing processes. In the case of homebrew you can run `brew services stop postgresql` to stop the service and prevent it from running at startup.

#### Development Machine Timezone Issues

If you are experiencing problems like redux forms that are 'dirty' when they shouldn't be on your local environment, it may be due to a mismatch
of your local dev machine's timezone and the assumption of UTC made by the local database. A detailed example of this sort of issue can be found in
[this story](https://www.pivotaltracker.com/n/projects/2136865/stories/160975609). A workaround for this is to set the TZ environment variable
in Mac OS for the context of your running app. This can be done by adding the following to `.envrc.local`:

```bash
export TZ="UTC"
```

Doing so will set the timezone environment variable to UTC utilizing the same localized context as your other `.envrc.local` settings.

#### Linters & Pre-commit Hooks

We use a number of linters for formatting, security and error checking. Please see the [pre-commit documentation](https://github.com/transcom/mymove/wiki/run-pre-commit-hooks) for a list of linters and troubleshooting tips.

#### Yarn install markdown-spell (aka mdspell)

We use `mdspell` for spell checking markdown files during pre-commit hooks. You may run into an issue such as below during the installation command `yarn global add markdown-spellcheck` suggested by the Makefile.

Example error:

```sh
>$ yarn global add markdown-spellcheck

yarn global v1.19.0
[1/4] :mag:  Resolving packages...
[2/4] :truck:  Fetching packages...
error An unexpected error occurred: "https://registry.yarnpkg.com/har-validator/-/har-validator-5.1.2.tgz: Request failed \"404 Not Found\"".
info If you think this is a bug, please open a bug report with the information provided in "/Users/john/.config/yarn/global/yarn-error.log".
info Visit https://yarnpkg.com/en/docs/cli/global for documentation about this command.
```

If you do, following these steps may resolve it.

```sh
rm ~/.config/yarn/global/yarn.lock
cd ~/.config/yarn/global
yarn cache clean
yarn global add markdown-spellcheck
```

### Manual Redeploys and Other Helpful Information in an Emergency

Like many modern software developers, we rely on a number of external services to be responsible for certain repeatable processes. One such example is CircleCI, which we use for deployment. It's a great tool in many ways, and reduces the surface area of what we ourselves have to manage, ideally transferring associated risk. However, it opens us up to a different risk: namely, what happens if CircleCI goes down and we need to deploy our app? For this circumstance, we have a series of scripts that can be run manually. They live in the script directory, and you can find information about how to run them [in this README under the Deployment Scripts heading](scripts/README.md#deployment-scripts).

Please add any other fear-inducing scenarios and our mitigation attempts here.

### PII Best Practices

Server side: any downloadable content passed to the client should by default have an inline content disposition (like PDFs).

Client side: before accessing any prod environment, fix your browser settings to display PDFs, not use an external app.

#### More about content dispositions

An inline disposition tells the browser to display a file inline if it can. If you instead set an attachment disposition, the browser will download the file even when it is capable of displaying the data itself. If an engineer downloads PII instead of letting their browser display it, this will cause a security incident. So make sure content like PDFs are passed to the client with inline dispositions. You can read more in [the official Mozilla docs](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition "MDN: Content-Disposition Headers").

#### More about browser settings

Even if an inline disposition is set, most browsers still allow you to override that behavior and automatically download certain file types.

Before working with the prod environment, ensure you have changed your settings to display PDFs in the browser. In most browsers, you can find the relevant setting by searching "PDF" in the settings menu.
