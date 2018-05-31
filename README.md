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
  * [Setup: Client](#setup-client)
  * [Setup: Office/admin client](#setup-officeadmin-client)
  * [Setup: S3](#setup-s3)
  * [TSP Award Queue](#tsp-award-queue)
  * [Test Data Generator](#test-data-generator)
  * [API / Swagger](#api--swagger)
  * [Testing](#testing)
  * [Logging](#logging)
  * [Database](#database)
    * [Dev Commands](#dev-commands)
    * [Migrations](#migrations)
    * [Secure Migrations](#secure-migrations)
  * [Environment Variables](#environment-variables)
  * [Documentation](#documentation)
  * [Spellcheck](#spellcheck)
    * [Tips for staying sane](#tips-for-staying-sane)
  * [Troubleshooting](#troubleshooting)

Regenerate with "bin/generate-md-toc.sh"

<!-- tocstop -->

## Supported clients

As of 3/6/2018, DDS has confirmed that support for IE is limited to IE 11 and Edge or newer versions. Currently, the intention is to encourage using Chrome and Firefox instead, with specific versions TBD. Research is incomplete on mobile browsers, but we are assuming support for iOS and Android.

## Client Network Dependencies

The client application (i.e. website) makes outbound requests to the following domains in its normal operation. If you have a firewall in place, it will need to be configured to allow outbound access to them for the application to operate.

* S3 for document downloads; exact domains TBD.

## Development

### Git

Use your work email when making commits to our repositories. The simplest path to correctness is setting global config:

    git config --global user.email "trussel@truss.works"
    git config --global user.name "Trusty Trussel"

If you drop the `--global` flag these settings will only apply to the current repo. If you ever re-clone that repo or clone another repo, you will need to remember to set the local config again. You won't. Use the global config. :-)

For web-based Git operations, GitHub will use your primary email unless you choose "Keep my email address private". If you don't want to set your work address as primary, please [turn on the privacy setting](https://github.com/settings/emails).

### Project location

All of Go's tooling expects Go code to be checked out in a specific location. Please read about [Go workspaces](https://golang.org/doc/code.html#Workspaces) for a full explanation. If you just want to get started, then decide where you want all your go code to live and configure the GOPATH environment variable accordingly. For example, if you want your go code to live at `~/code/go`, you should add the following like to your `.bash_profile`:

```bash
export GOPATH=~/code/go
```

If you are OK with using the default location for go code (`~/go`), then there is nothing to do. Since this is the default location, using it means you do not need to set `$GOPATH` yourself.

_Regardless of where your go code is located_, you need to add `$GOPATH/bin` to your `PATH` so that executables installed with the go tooling can be found. Add the following to your `.bash_profile`:

```bash
export PATH=$(go env GOPATH)/bin:$PATH
```

Once that's done, you have go installed, and you've re-sourced your profile, you can checkout this repository by running `go get github.com/transcom/mymove/cmd/webserver` (This will emit an error "can't load package:" but will have cloned the source correctly). You will then find the code at `$GOPATH/src/github.com/transcom/mymove`

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
1. `make db_dev_migrate`
1. `make server_run`
1. `make client_run`

### Setup: Prerequisites

* Install Go with Homebrew. Make sure you do not have other installations.
* Run `bin/prereqs` and install everything it tells you to. _Do not configure PostgreSQL to automatically start at boot time!_
* For managing local environment variables, we're using [direnv](https://direnv.net/). You need to [configure your shell to use it](https://direnv.net/).
* Run `direnv allow` to load up the `.envrc` file. Add a `.envrc.local` file with any values it asks you to define.
* Run `make deps`.
* [EditorConfig](http://editorconfig.org/) allows us to manage editor configuration (like indent sizes,) with a [file](https://github.com/transcom/ppp/blob/master/.editorconfig) in the repo. Install the appropriate plugin in your editor to take advantage of that.

### Setup: Database

You will need to setup a local database before you can begin working on the local server / client. Docker will need to be running for any of this to work.

1. `make db_dev_migrate`: Creates a PostgreSQL docker container if you haven't made one yet and runs all existing database migrations, which do things like creating table structures, etc. You will run this command again anytime you add new migrations to the app (see below for more)

You can validate that your dev database is running by running `bin/psql-dev`. This puts you in a PostgreSQL shell. Type `\dt` to show all tables, and `\q` to quit.
You can validate that your test database is running by running `bin/psql-test`. This puts you in a PostgreSQL shell. Type `\dt` to show all tables, and `\q` to quit.

### Setup: Server

1. `make server_run`: installs dependencies, then builds and runs the server using `gin`, which is a hot reloading go server. It will listen on 8080 and will rebuild the actual server any time a go file changes. Pair this with `make client_run` to have hot reloading of the entire application.

In rare cases, you may want to run the server standalone, in which case you can run `make server_run_standalone`. This will build both the client and the server and this invocation can be relied upon to be serving the client JS on its own rather than relying on webpack doing so as when you run `make client_run`. You can run this without running `make client_run` and the whole app should work.

You can verify the server is working as follows:

`> curl http://localhost:8080/api/v1/issues --data "{ \"description\": \"This is a test issue\"}" -H "Content-Type: application/json"`

from which the response should be like

`{"id":"d5735bc0-7553-4d80-a42d-ea1e50bbcfc4", "description": "This is a test issue", "created_at": "2018-01-04 14:47:28.894988", "updated_at": "2018-01-04 14:47:28.894988"}`

Dependencies are managed by [dep](https://github.com/golang/dep). New dependencies are automatically detected in import statements. To add a new dependency to the project, import it in a source file and then run `dep ensure`

### Setup: Client

1. `make client_run`

The above will start the webpack dev server, serving the front-end on port 3000. If paired with `make server_run` then the whole app will work, the webpack dev server proxies all API calls through to the server.

Dependencies are managed by yarn. To add a new dependency, use `yarn add`

### Setup: Office/admin client

1. add the following line to /etc/hosts
    `127.0.0.1 officelocal`
2. Ensure that you have a test account which can log into the office site...
    * `make tools_build` to build the tools
    * run `bin/make-office-user -email <email>` to set up an office user associated with that email address
3. `make office_client_run`
4. Login with the email used above to access the office

### Setup: S3

If you want to develop against the live S3 service, you will need to configure the following values in your `.envrc`:

```text
AWS_S3_BUCKET_NAME
AWS_S3_KEY_NAMESPACE
AWS_REGION
AWS_PROFILE
```

AWS credentials should *not* be added to `.envrc` and should instead be setup using [the instructions in transcom-ppp](https://github.com/transcom/ppp-infra/blob/master/transcom-ppp/README.md#setup).

### TSP Award Queue

This background job is built as a separate binary which can be built using
`make tools_build` and run using `make tsp_run`.

### Test Data Generator

When creating new features, it is helpful to have sample data for the feature to interact with. The TSP Award Queue is an example of that--it matches shipments to TSPs, and it's hard to tell if it's working without some shipments and TSPs in the database!

* `make tools_build` will build the fake data generator binary
* `bin/generate_test_data` will run binary and create a preconfigured set of test data. To determine the data scenario you'd like to use, check out scenarios in the `testdatagen` package. Each scenario contains a description of what data will be created when the scenario is run. Pass the scenario in as a flag to the generate-test-data function. A sample command: `./bin/generate-test-data -scenario=2`. If you'd like to further specify how the data should look, you can specify the number of awards for each TSP performance by use the flag `rounds` with one of three arguments: 'none', 'half', and 'full'. This will create the TSP performance records with either no rounds of awards completed, half a round, or a full round. It will default to none if not specified. To specify how many TSPs should be created, use the flag `numTSP`. It will default to 15 if not specified. A sample command: `./bin/generate-test-data -rounds=half -numTSP=6`. You can use the `numTSP` and `rounds` in conjunction, but you cannot use them with the pre-packaged scenarios.

There is also a package (`/pkg/testdatagen`) that can be imported to create arbitrary test data. This could be used in tests, so as not to duplicate functionality.

Currently, scenarios have the following numbers:

* `-scenario=1` for Award Queue Scenario 1
* `-scenario=2` for Award Queue Scenario 2
* `-scenario=3` for Duty Station Scenario
* `-scenario=4` for PPM or PPM SIT Estimate Scenario (can also use Rate Engine Scenarios for Estimates)
* `-scenario=5` for Rate Engine Scenario 1
* `-scenario=6` for Rate Engine Scenario 2

### API / Swagger

The public API is defined in a single file: `swagger/api.yaml` and served at `/api/v1/swagger.yaml`. This file is the single source of truth for the public API. In addition, internal services, i.e. endpoints only intended for use by the React client are defined in `swagger/internal.yaml` and served at `/internal/swagger.yaml`. These are, as the name suggests, internal endpoints and not intended for use by external clients.

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

#### Dev Commands

There are a few handy targets in the Makefile to help you interact with the dev database:

* `make db_dev_run`: Initializes a new database if it does not exist and runs it, or starts the previously initialized Docker container if it has been stopped.
* `make db_dev_reset`: Destroys your database container. Useful if you want to start from scratch.
* `make db_dev_migrate`: Applies database migrations against your running database container.
* `make db_dev_migrate_down`: reverts the most recently applied migration by running the down migration

#### Migrations

If you need to change the database schema, you'll need to write a migration.

Creating a migration:

Use soda (a part of [pop](https://github.com/gobuffalo/pop/)) to generate migrations. In order to make using soda easy, a wrapper is in `./bin/soda` that sets the go environment and working directory correctly.

If you are generating a new model, use `./bin/gen_model model-name column-name:type column-name:type ...`. id, created_at and updated_at are all created automatically.

If you are modifying an existing model, use `./bin/soda generate migration migration-name` and add the [Fizz instructions](https://github.com/gobuffalo/pop/blob/master/fizz/README.md) yourself to the created files.

Running migrations in local development:

1. Use `make db_dev_migrate` to run migrations against your local dev environment.
1. Use `make db_dev_migrate_down` to revert the most recently applied migration. This is useful while you are developing a migration but should not be necessary otherwise.

Running migrations on Staging / Production:

Migrations are run automatically by CircleCI as part of the standard deploy process.

1. CircleCI builds and registers a container that includes the `soda` binary, along with migrations files.
1. CircleCI deploys this container to ECS and runs it as a one-off 'task'.
1. Migrations run inside the container against the environment's database.
1. If migrations fail, CircleCI fails the deploy.
1. If migrations pass, CircleCI continues with the deploy.

#### Secure Migrations

We are piggy-backing on the migration system for importing static datasets. This approach causes problems if the data isn't public, as all of the migrations are in this open source repository. To address this, we have what are called "secure migrations."

To create a secure migration:

* You create a regular Fizz migration, as described in the previous section.
* The body of the Fizz migration should be of the format: `exec('./apply-secure-migration.sh ${FILENAME}')`.
* `${FILENAME}` should be the same as the name of the migration that Fizz created for you, except with an `.sql` extension.
* Make a test migration
  * Create an SQL file with a mock version of the data (nothing sensitive,) and save it in `local_migrations/${FILENAME}`.
  * **This migration will be run every time migrations are run against a dev or testing database.**
  * Test the migration with: `make db_dev_reset && make db_dev_migrate`.
  * You should see your test migration run.
* Create the real migration
  * Make the actual migration SQL file with the secrets in it, and give it the same filename.
  * Upload it to the `staging` and `prod` S3 buckets: `transcom-ppp-app-{{environment}}-us-west-2/secure-migrations/`.
  * You should see other migrations in the S3 bucket. If you don't, you may be looking in the wrong place.
* Voila!
  * When you land your PR, CI will run your migration against staging. If that is successful, it will then run against production.

Gory Details:

When this migration is run, `soda` will shell out to our script, `apply-secure-migration.sh`. This script will:

* Look at `$SECURE_MIGRATION_SOURCE` to determine if the migrations should be found locally (`local`, for dev & testing,) or on S3 (`s3`).
* If the file is to be found on S3, it is downloaded from `${AWS_S3_BUCKET_NAME}/secure-migrations/${FILENAME}`.
* If it is to be found locally, the script looks for it in `$SECURE_MIGRATION_DIR`.
* Regardless of where the migration comes from, it is then applied to the database by essentially doing: `psql < ${FILENAME}`.

There is an example of a secure migration [in the repo](https://github.com/transcom/mymove/blob/master/migrations/20180424010930_test_secure_migrations.up.fizz).

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
* If you have problems connecting to PostgreSQL, or running related scripts, make sure you aren't already running a PostgreSQL daemon. You can check this by typing `ps aux | grep postgres` and looking for existing processes.
* If you happen to have installed pre-commit in a virtual environment not with brew, running bin/prereqs will not alert you. You may run into issues when running `make deps`. To install pre-commit: `brew install pre-commit`.
* If you're having trouble accessing the API docs or the server is otherwise misbehaving, try stopping the server, running `make client_build`, and then running `make client_run` and `make server_run`.
