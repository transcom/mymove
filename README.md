# Personal Property Prototype

[![Build status](https://img.shields.io/circleci/project/github/transcom/mymove/master.svg)](https://circleci.com/gh/transcom/mymove/tree/master)

[![GoDoc](https://godoc.org/github.com/transcom/mymove?status.svg)](https://godoc.org/github.com/transcom/mymove)

This repository contains the application source code for the Personal Property Prototype, a possible next generation version of the Defense Personal Property System (DPS). DPS is an online system managed by the U.S. [Department of Defense](https://www.defense.gov/) (DoD) [Transportation Command](http://www.ustranscom.mil/) (USTRANSCOM) and is used by service members and their families to manage household goods moves.

This prototype was built by a [Defense Digital Service](https://www.dds.mil/) team in support of USTRANSCOM's mission.

## Table of Contents

<!-- Table of Contents auto-generated with `scripts/generate-md-toc` -->

<!-- toc -->

* [Supported Browsers](#supported-browsers)
* [Application Setup](#application-setup)
  * [Setup: Developer Setup](#setup-developer-setup)
  * [Setup: Git](#setup-git)
  * [Setup: Golang](#setup-golang)
  * [Setup: Project Checkout](#setup-project-checkout)
  * [Setup: Project Layout](#setup-project-layout)
  * [Setup: Editor Config](#setup-editor-config)
  * [Setup: Makefile](#setup-makefile)
  * [Setup: Quick Initial Setup](#setup-quick-initial-setup)
  * [Setup: Prerequisites](#setup-prerequisites)
  * [Setup: Direnv](#setup-direnv)
  * [Setup: Pre-Commit](#setup-pre-commit)
  * [Setup: Hosts](#setup-hosts)
  * [Setup: Dependencies](#setup-dependencies)
  * [Setup: Build Tools](#setup-build-tools)
  * [Setup: Database](#setup-database)
  * [Setup: Server](#setup-server)
  * [Setup: MilMoveLocal Client](#setup-milmovelocal-client)
  * [Setup: OfficeLocal client](#setup-officelocal-client)
  * [Setup: TSPLocal client](#setup-tsplocal-client)
  * [Setup: AdminLocal client](#setup-adminlocal-client)
  * [Setup: DPS user](#setup-dps-user)
  * [Setup: Orders Gateway](#setup-orders-gateway)
  * [Setup: AWS Services (Optional)](#setup-aws-services-optional)
* [Development](#development)
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
  * [GoLand](#goland)
  * [Troubleshooting](#troubleshooting)
    * [Postgres Issues](#postgres-issues)
    * [Development Machine Timezone Issues](#development-machine-timezone-issues)
    * [Linters & Pre-commit Hooks](#linters--pre-commit-hooks)

Regenerate with "scripts/generate-md-toc"

<!-- tocstop -->

## Supported Browsers

As of 3/6/2018, DDS has confirmed that support for IE is limited to IE 11 and Edge or newer versions. Currently, the intention is to encourage using Chrome and Firefox instead, with specific versions TBD. Research is incomplete on mobile browsers, but we are assuming support for iOS and Android. For more information please read [ADR0016 Browser Support](./docs/adr/0016-Browser-Support.md).

## Application Setup

### Setup: Developer Setup

There are a number of things you'll need at a minimum to be able to check out, develop and run this project.

* Install [Homebrew](https://brew.sh)
  * Use the following command `/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"`
* We always use the latest version of Go unless there's a known conflict (which will be announced by the team).
  * Install it with Homebrew: `brew install go`
  * Pin it, so that you don't accidentally upgrade before we upgrade the project: `brew pin go`
  * When we upgrade the project's go version, unpin, upgrade, and then re-pin: `brew unpin go; brew upgrade go; brew pin go`
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
* `local_migrations`: Database migrations used locally in place of secure migrations
* `migrations`: Database migrations
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

1. `make prereqs`
1. `direnv allow`
1. `make ensure_pre_commit`
1. `make deps`
1. `make db_dev_run`
1. `make db_dev_migrate`
1. `make server_run`
1. `make client_build`
1. `make client_run`

### Setup: Prerequisites

Run `make prereqs` and install everything it tells you to. Most of the prerequisites need you to use `brew install <package>`.

**NOTE:** Do not configure PostgreSQL to automatically start at boot time or the DB commands will not work correctly!

### Setup: Direnv

For managing local environment variables, we're using [direnv](https://direnv.net/). You need to [configure your shell to use it](https://direnv.net/). For bash, add the command `eval "$(direnv hook bash)"` to whichever file loads upon opening bash (likely `~./bash_profile`, though instructions say `~/.bashrc`).

Run `direnv allow` to load up the `.envrc` file. It should complain that you have missing variables which you will rectify in one of the following ways.

You can add a `.envrc.local` file. One way to do this is to run `chamber env app-devlocal >> .envrc.local`. If you don't have access to chamber you can also `touch .envrc.local` and add any values that the output from direnv asks you to define. Instructions are in the error messages.

If you wish to not maintain a `.envrc.local` you can alternatively run `cp .envrc.chamber.template .envrc.chamber` to enable getting secret values from `chamber`. **Note** that this method does not work for users of the `fish` shell unless you replace `direnv allow` with `direnv export fish | source`.

### Setup: Pre-Commit

Run `pre-commit install` to install a pre-commit hook into `./git/hooks/pre-commit`.  This is different than `brew install pre-commit` and must be done so that the hook will check files you are about to commit to the repository.  Next install the pre-commit hook libraries with `pre-commit install-hooks`.

Before running `pre-commit run -a` you will need to install Javascript dependencies and generate some golang code from Swagger files. An easier way to handle this is by running `make pre_commit_tests` or `make server generate client_deps && pre-commit run -a`. But it's early to do this so you can feel free to skip running the pre-commit checks at this time.

### Setup: Hosts

You need to modify your `/etc/hosts` file. This is a tricky file to modify and you will need to use `sudo` to edit it.
Here are the steps:

  ```bash
  echo "127.0.0.1 milmovelocal" | sudo tee -a /etc/hosts
  echo "127.0.0.1 officelocal" | sudo tee -a /etc/hosts
  echo "127.0.0.1 tsplocal" | sudo tee -a /etc/hosts
  echo "127.0.0.1 orderslocal" | sudo tee -a /etc/hosts
  echo "127.0.0.1 adminlocal" | sudo tee -a /etc/hosts
  ```

Check that the file looks correct with `cat /etc/hosts`:

  ```text
  ##
  # Host Database
  #
  # localhost is used to configure the loopback interface
  # when the system is booting.  Do not change this entry.
  ##
  255.255.255.255 broadcasthost
  ::1             localhost
  127.0.0.1   localhost
  127.0.0.1   milmovelocal
  127.0.0.1   officelocal
  127.0.0.1   tsplocal
  127.0.0.1   orderslocal
  127.0.0.1   adminlocal
  ```

### Setup: Dependencies

Run `make deps`. This will check your system for any setup issues. Then it will ensure that you have installed pre-commit
and go on to install the client (javascript) and server (golang) dependencies for you.

### Setup: Build Tools

Run `make build_tools` to get all the server and tool dependencies built. These will be needed in future steps to not only generate test data but also to interact with the database and more.

### Setup: Database

You will need to setup a local database before you can begin working on the local server / client. Docker will need to be running for any of this to work.

1. `make db_dev_run`: Creates a PostgreSQL docker container if you haven't made one yet

1. `make db_dev_migrate`:  Runs all existing database migrations, which does things like creating table structures, etc. You will run this command again anytime you add new migrations to the app (see below for more)

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

If both the server and client are running, you should be able to view the Swagger UI at <http://milmovelocal:3000/api/v1/docs>.  If it does not, try running `make client_build` (this only needs to be run the first time).

Dependencies are managed by yarn. To add a new dependency, use `yarn add`

### Setup: OfficeLocal client

1. Ensure that you have a test account which can log into the office site...
    * run `generate-test-data --named-scenario e2e_basic` to load test data
    * Log into "Local Sign In" and either select a pre-made user or use the button to create a new user
2. `make office_client_run`
3. Login with the email used above to access the office

### Setup: TSPLocal client

1. Ensure that you have a test account which can log into the TSP site...
    * run `generate-test-data --named-scenario e2e_basic` to load test data
    * Log into "Local Sign In" and either select a pre-made user or use the button to create a new user
2. `make tsp_client_run`
3. Login with the email used above to access the TSP

### Setup: AdminLocal client

1. `make admin_client_run`

### Setup: DPS user

1. Ensure that you have a login.gov test account
2. run `make-dps-user -email <email>` to set up a DPS user associated with that email address

### Setup: Orders Gateway

Nothing to do.

### Setup: AWS Services (Optional)

If you want to develop against AWS services you will need an AWS user account with `engineering` privileges. Then you will need to configure the `PPP_INFRA_PATH` in your `.envrc.local`.

AWS credentials are managed via `aws-vault`. See the [the instructions in transcom-ppp](https://github.com/transcom/ppp-infra/blob/master/transcom-ppp/README.md#setup) to set things up.

## Development

### TSP Award Queue

This background job is built as a separate binary which can be built using `make tsp_run`.

### Test Data Generator

When creating new features, it is helpful to have sample data for the feature to interact with. The TSP Award Queue is an example of that--it matches shipments to TSPs, and it's hard to tell if it's working without some shipments and TSPs in the database!

* `make bin/generate-test-data` will build the fake data generator binary
* `bin/generate-test-data -named-scenario="e2e_basic"` will populate the database with a handful of users in various stages of progress along the flow. The emails are named accordingly (see [`e2ebasic.go`](https://github.com/transcom/mymove/blob/master/pkg/testdatagen/scenario/e2ebasic.go)). Alternatively, run `make db_dev_e2e_populate` to reset your db and populate it with e2e user flow cases.
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

The Admin API is defined in the file `swagger/admin.yaml` and served at `/admin/v1/swagger.yaml`.

You can view the API's documentation (powered by Swagger UI) at <http://localhost:3000/api/v1/docs> when a local server is running.

### Testing

There are a few handy targets in the Makefile to help you run tests:

* `make client_test`: Run front-end testing suites.
* `make server_test`: Run back-end testing suites. [Additional info for running go tests](https://github.com/transcom/mymove/blob/master/docs/how-to/run-go-tests.md)
* `make e2e_test`: Run e2e testing suite.
* `make test`: Run e2e, client- and server-side testing suites.

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

To add new regular and/or secure migrations, see the [database development guide](./docs/database.md)

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
with `chamber write app-devlocal <key> <value>`. For long blocks of text like certificates you can write them with
`echo "$LONG_VALUE" | chamber write app-devlocal <key> -`.

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

We use a number of linters for formatting, security and error checking. Please see this [how-to document](./docs/how-to/run-pre-commit-hooks.md) for a list of linters and troubleshooting tips.
