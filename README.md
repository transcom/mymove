# Personal Property Prototype

[![Build status](https://img.shields.io/circleci/project/github/transcom/mymove/master.svg)](https://circleci.com/gh/transcom/mymove/tree/master)

[![GoDoc](https://godoc.org/github.com/transcom/mymove?status.svg)](https://godoc.org/github.com/transcom/mymove)

This repository contains the application source code for the Personal Property Prototype, a possible next generation version of the Defense Personal Property System (DPS). DPS is an online system managed by the U.S. [Department of Defense](https://www.defense.gov/) (DoD) [Transportation Command](http://www.ustranscom.mil/) (USTRANSCOM) and is used by service members and their families to manage household goods moves.

This prototype was built by a [Defense Digital Service](https://www.dds.mil/) team in support of USTRANSCOM's mission.

## Table of Contents

<!-- Table of Contents auto-generated with `bin/generate-md-toc.sh` -->

<!-- toc -->

* [Supported clients](#supported-clients)
* [Client Network Dependencies](#client-network-dependencies)
* [Development](#development)
  * [Git](#git)
  * [Project location](#project-location)
  * [Project Layout](#project-layout)
  * [Setup: Initial Setup](#setup-initial-setup)
  * [Setup: Prerequisites](#setup-prerequisites)
  * [Setup: Database](#setup-database)
  * [Setup: Server](#setup-server)
  * [Setup: MilMoveLocal Client](#setup-milmovelocal-client)
  * [Setup: OfficeLocal client](#setup-officelocal-client)
  * [Setup: TSPLocal client](#setup-tsplocal-client)
  * [Setup: DPS user](#setup-dps-user)
  * [Setup: Orders Gateway](#setup-orders-gateway)
  * [Setup: S3](#setup-s3)
  * [TSP Award Queue](#tsp-award-queue)
  * [Test Data Generator](#test-data-generator)
  * [API / Swagger](#api--swagger)
  * [Testing](#testing)
  * [Logging](#logging)
  * [Database](#database)
    * [Dev DB Commands](#dev-db-commands)
    * [Test DB Commands](#test-db-commands)
    * [Migrations](#migrations)
  * [Environment Variables](#environment-variables)
  * [Documentation](#documentation)
  * [Spellcheck](#spellcheck)
    * [Tips for staying sane](#tips-for-staying-sane)
  * [Troubleshooting](#troubleshooting)
    * [Postgres Issues](#postgres-issues)
    * [Development Machine Timezone Issues](#development-machine-timezone-issues)

Regenerate with "bin/generate-md-toc.sh"

<!-- tocstop -->

## Supported clients

As of 3/6/2018, DDS has confirmed that support for IE is limited to IE 11 and Edge or newer versions. Currently, the intention is to encourage using Chrome and Firefox instead, with specific versions TBD. Research is incomplete on mobile browsers, but we are assuming support for iOS and Android.

## Client Network Dependencies

The client application (i.e. website) makes outbound requests to the following domains in its normal operation. If you have a firewall in place, it will need to be configured to allow outbound access to them for the application to operate.

* S3 for document downloads; exact domains TBD.
* Honeycomb for server-side debugging and observability. Currently being tested in staging and experimental environments.

## Development

### Git

Use your work email when making commits to our repositories. The simplest path to correctness is setting global config:

    git config --global user.email "trussel@truss.works"
    git config --global user.name "Trusty Trussel"

If you drop the `--global` flag these settings will only apply to the current repo. If you ever re-clone that repo or clone another repo, you will need to remember to set the local config again. You won't. Use the global config. :-)

For web-based Git operations, GitHub will use your primary email unless you choose "Keep my email address private". If you don't want to set your work address as primary, please [turn on the privacy setting](https://github.com/settings/emails).

Note that with 2-factor-authentication enabled, in order to push local code to GitHub through HTTPS, you need to [create a personal access token](https://gist.github.com/ateucher/4634038875263d10fb4817e5ad3d332f) and use that as your password.

### Project location

All of Go's tooling expects Go code to be checked out in a specific location. Please read about [Go workspaces](https://golang.org/doc/code.html#Workspaces) for a full explanation. If you just want to get started, then decide where you want all your go code to live and configure the GOPATH environment variable accordingly. For example, if you want your go code to live at `~/code/go`, you should add the following like to your `.bash_profile`:

```bash
export GOPATH=~/code/go
```

A few of our custom tools expect the `GOPATH` environment variable to be defined.  If you'd like to use the default location, then add the following to your `.bash_profile` or hardcode the default value.  This line will set the GOPATH environment variable to the value of `go env GOPATH` if it is not already set.

```bash
export GOPATH=${GOPATH:-$(go env GOPATH)}
```

_Regardless of where your go code is located_, you need to add `$GOPATH/bin` to your `PATH` so that executables installed with the go tooling can be found. Add the following to your `.bash_profile`:

```bash
export PATH=$(go env GOPATH)/bin:$PATH
```

Once that's done, you have go installed, and you've re-sourced your profile, you can checkout this repository by running `go get github.com/transcom/mymove/cmd/webserver` (This will emit an error "can't load package:" or multiple errors with "Cannot find package" but will have cloned the source correctly). You will then find the code at `$GOPATH/src/github.com/transcom/mymove`

If you have already checked out the code somewhere else, you can just move it to be in the above location and everything will work correctly.

### Project Layout

All of our code is intermingled in the top level directory of mymove. Here is an explanation of what some of these directories contain:

`bin`: A location for tools helpful for developing this project \
`build`: The build output directory for the client. This is what the development server serves \
`cmd`: The location of main packages for any go binaries we build \
`config`: Config files can be dropped here \
`docs`: A location for docs for the project. This is where ADRs are \
`migrations`: Database migrations live here \
`node_modules`: Cached dependencies for the client \
`pkg`: The location of all of our go libraries, most of our go code lives here \
`public`: The client's static resources \
`src`: The react source code for the client \
`vendor`: Cached dependencies for the server

### Setup: Initial Setup

The following commands will get mymove running on your machine for the first time. Please read below for explanations of each of the commands.

1. `./bin/prereqs`
1. `make db_dev_run`
1. `make db_dev_migrate`
1. `make server_run`
1. `make client_build`
1. `make client_run`

### Setup: Prerequisites

* Install Go version 1.11.4 with Homebrew. Make sure you do not have other installations.
  * `brew install go@1.11.4`
  * If a later Go version exists, Brew will warn you that "go@1.11.4 is keg-only, which means it was not symlinked". If that happens, add the following to your bash config: `export PATH="/usr/local/opt/go@1.11.4/bin:$PATH"`. This line needs to appear in the file before your Go paths are declared.
* Run `bin/prereqs` and install everything it tells you to. _Do not configure PostgreSQL to automatically start at boot time or the DB commands will not work correctly!_
* For managing local environment variables, we're using [direnv](https://direnv.net/). You need to [configure your shell to use it](https://direnv.net/). For bash, add the command `eval "$(direnv hook bash)"` to whichever file loads upon opening bash (likely `~./bash_profile`, though instructions say `~/.bashrc`).
* Run `direnv allow` to load up the `.envrc` file. Add a `.envrc.local` file with any values it asks you to define.
* Run `make deps`.
* [EditorConfig](http://editorconfig.org/) allows us to manage editor configuration (like indent sizes,) with a [file](https://github.com/transcom/ppp/blob/master/.editorconfig) in the repo. Install the appropriate plugin in your editor to take advantage of that.
* Run `pre-commit install` to install a pre-commit hook into `./git/hooks/pre-commit`.  This is different than `brew install pre-commit` and must be done so that the hook will check files you are about to commit to the repository.  Also, using this hook is much faster than attempting to create your own with `pre-commit run -a`.

### Setup: Database

You will need to setup a local database before you can begin working on the local server / client. Docker will need to be running for any of this to work.

1. `make db_dev_run`: Creates a PostgreSQL docker container if you haven't made one yet

1. `make db_dev_migrate`:  Runs all existing database migrations, which does things like creating table structures, etc. You will run this command again anytime you add new migrations to the app (see below for more)

You can validate that your dev database is running by running `bin/psql-dev`. This puts you in a PostgreSQL shell. Type `\dt` to show all tables, and `\q` to quit.
You can validate that your test database is running by running `bin/psql-test`. This puts you in a PostgreSQL shell. Type `\dt` to show all tables, and `\q` to quit.

If you are stuck on this step you may need to see the section on Troubleshooting.

### Setup: Server

1. `make server_run`: installs dependencies, then builds and runs the server using `gin`, which is a hot reloading go server. It will listen on 8080 and will rebuild the actual server any time a go file changes. Pair this with `make client_run` to have hot reloading of the entire application.

In rare cases, you may want to run the server standalone, in which case you can run `make server_run_standalone`. This will build both the client and the server and this invocation can be relied upon to be serving the client JS on its own rather than relying on webpack doing so as when you run `make client_run`. You can run this without running `make client_run` and the whole app should work.

Dependencies are managed by [dep](https://github.com/golang/dep). New dependencies are automatically detected in import statements. To add a new dependency to the project, import it in a source file and then run `dep ensure`

### Setup: MilMoveLocal Client

1. add the following line to /etc/hosts
    `127.0.0.1 milmovelocal`
1. `make client_build` (if setting up for first time)
1. `make client_run`

The above will start the webpack dev server, serving the front-end on port 3000. If paired with `make server_run` then the whole app will work, the webpack dev server proxies all API calls through to the server.

If both the server and client are running, you should be able to view the Swagger UI at <http://milmovelocal:3000/api/v1/docs>.  If it does not, try running `make client_build` (this only needs to be run the first time).

Dependencies are managed by yarn. To add a new dependency, use `yarn add`

### Setup: OfficeLocal client

1. add the following line to /etc/hosts
    `127.0.0.1 officelocal`
2. Ensure that you have a test account which can log into the office site...
    * `make build_tools` to build the tools
    * run `bin/make-office-user -email <email>` to set up an office user associated with that email address
3. `make office_client_run`
4. Login with the email used above to access the office

### Setup: TSPLocal client

1. add the following line to /etc/hosts
    `127.0.0.1 tsplocal`
2. Ensure that you have a test account which can log into the TSP site...
    * `make build_tools` to build the tools
    * run `./bin/generate-test-data -scenario=7` to load test data
    * run `bin/make-tsp-user -email <email>` to set up a TSP user associated with that email address
3. `make tsp_client_run`
4. Login with the email used above to access the TSP

### Setup: DPS user

1. Ensure that you have a login.gov test account
    * `make build_tools` to build the tools
    * run `bin/make-dps-user -email <email>` to set up a DPS user associated with that email address

### Setup: Orders Gateway

1. add the following line to /etc/hosts
    `127.0.0.1 orderslocal`

### Setup: S3

If you want to develop against the live S3 service, you will need to configure the following values in your `.envrc`:

```text
AWS_S3_BUCKET_NAME
AWS_S3_KEY_NAMESPACE
AWS_REGION
AWS_PROFILE
PPP_INFRA_PATH
```

AWS credentials are managed via `aws-vault`. See the [the instructions in transcom-ppp](https://github.com/transcom/ppp-infra/blob/master/transcom-ppp/README.md#setup) to set things up.

### TSP Award Queue

This background job is built as a separate binary which can be built using
`make build_tools` and run using `make tsp_run`.

### Test Data Generator

When creating new features, it is helpful to have sample data for the feature to interact with. The TSP Award Queue is an example of that--it matches shipments to TSPs, and it's hard to tell if it's working without some shipments and TSPs in the database!

* `make build_tools` will build the fake data generator binary
* `bin/generate-test-data -named-scenario="e2e_basic"` will populate the database with a handful of users in various stages of progress along the flow. The emails are named accordingly (see [`e2ebasic.go`](https://github.com/transcom/mymove/blob/master/pkg/testdatagen/scenario/e2ebasic.go)). Alternatively, run `make db_populate_e2e` to reset your db and populate it with e2e user flow cases.
* `bin/generate-test-data` will run binary and create a preconfigured set of test data. To determine the data scenario you'd like to use, check out scenarios in the `testdatagen` package. Each scenario contains a description of what data will be created when the scenario is run. Pass the scenario in as a flag to the generate-test-data function. A sample command: `./bin/generate-test-data -scenario=2`.

There is also a package (`/pkg/testdatagen`) that can be imported to create arbitrary test data. This could be used in tests, so as not to duplicate functionality.

Currently, scenarios have the following numbers:

* `-scenario=1` for Award Queue Scenario 1
* `-scenario=2` for Award Queue Scenario 2
* `-scenario=3` for Duty Station Scenario
* `-scenario=4` for PPM or PPM SIT Estimate Scenario (can also use Rate Engine Scenarios for Estimates)
* `-scenario=5` for Rate Engine Scenario 1
* `-scenario=6` for Rate Engine Scenario 2
* `-scenario=7` for TSP test data

### API / Swagger

The public API is defined in a single file: `swagger/api.yaml` and served at `/api/v1/swagger.yaml`. This file is the single source of truth for the public API.

In addition, internal services, i.e. endpoints only intended for use by the React client are defined in `swagger/internal.yaml` and served at `/internal/swagger.yaml`. These are, as the name suggests, internal endpoints and not intended for use by external clients.

The Orders Gateway's API is defined in the file `swagger/orders.yaml` and served at `/orders/v0/orders.yaml`.

You can view the API's documentation (powered by Swagger UI) at <http://localhost:3000/api/v1/docs> when a local server is running.

### Testing

There are a few handy targets in the Makefile to help you run tests:

* `make client_test`: Run front-end testing suites.
* `make server_test`: Run back-end testing suites.
* `make e2e_test`: Run e2e testing suite. To run locally, add an environment variable called SAUCE_ACCESS_KEY, which you can find in team DP3 Engineering Vault of 1Password under Sauce Labs or by logging in to Sauce itself. In 1Password, the access key is labeled SAUCE_ACCESS_KEY. This will run against our staging environment. If you want to point to another instance, add an environment variable called E2E_BASE with the base url for the instance. Note that to test a development instance, you must run `make server_run_standalone` and set up a tunnel (via ngrok or localtunnel).
* `make test`: Run e2e, client- and server-side testing suites.

To run an individual test: `go test ./pkg/rateengine/ -testify.m Test_Scenario1`

### Logging

We are using [zap](https://github.com/uber-go/zap) as a logger in this project. We currently rely on its built-in `NewDevelopment()` and `NewProduction()` default configs, which are enabled in any of the executable packages that live in `cmd`.

This means that logging *is not* set up from within models or other packages unless the files in `cmd` are also being loaded. *If you attempt to call `zap.L()` or `zap.S()` without a configured logger, nothing will appear on the screen.*

If you need to see some output during the development process (say, for debugging purposes), it is best to use the standard lib `fmt` package to print to the screen. You will also need to pass `-v` to `go test` so that it prints all output, even from passing tests. The simplest way to do this is to run `go test` yourself by passing it which files to run, e.g. `go test pkg/models/* -v`.

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
* `make db_dev_e2e_populate`: Populate data with data used to run e2e tests

#### Test DB Commands

The Dev Commands are used to talk to the dev DB.  If you were working with the test DB you would use these commands:

* `make db_test_run`
* `make db_test_create`
* `make db_test_reset`
* `make db_test_migrate`
* `make db_test_e2e_populate`

The test DB commands all talk to the DB over localhost.  But in a docker-only environment (like CircleCI) you may not be able to use those commands, which is why `*_docker` versions exist for all of them:

* `make db_test_run_docker`
* `make db_test_create_docker`
* `make db_test_reset_docker`
* `make db_test_migrate_docker`

#### Migrations

To add new regular and/or secure migrations, see the [database development guide](./docs/database.md)

Running migrations in local development:

Use `make db_dev_migrate` to run migrations against your local dev environment.

Running migrations on Staging / Production:

Migrations are run automatically by CircleCI as part of the standard deploy process.

1. CircleCI builds and registers a container that includes the `soda` binary, along with migrations files.
1. CircleCI deploys this container to ECS and runs it as a one-off 'task'.
1. Migrations run inside the container against the environment's database.
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
    require NEW_ENV_VAR "Look for info on this value in Google Drive"
    ```

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

### Troubleshooting

* Random problems may arise if you have old Docker containers running. Run `docker ps` and if you see containers unrelated to our app, consider stopping them.
* If you happen to have installed pre-commit in a virtual environment not with brew, running bin/prereqs will not alert you. You may run into issues when running `make deps`. To install pre-commit: `brew install pre-commit`.
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
