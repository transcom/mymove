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
  * [Creating alternative users with the same email address](#creating-alternative-users-with-the-same-email-address)
* [Project Layout](#project-layout)
* [Application Setup](#application-setup)
  * [Setup: Base Setup](#setup-base-setup)
    * [Homebrew](#homebrew)
    * [Setup: Git](#setup-git)
    * [Setup: Project Checkout](#setup-project-checkout)
    * [Setup: Editor Config](#setup-editor-config)
  * [Setup: Nix](#setup-nix)
    * [Nix: Initial Setup](#nix-initial-setup)
    * [Nix: Clean Up Local Env](#nix-clean-up-local-env)
    * [Nix: Installing Dependencies](#nix-installing-dependencies)
  * [Setup: Manual](#setup-manual)
    * [Manual: Prerequisites](#manual-prerequisites)
  * [Setup: Shared](#setup-shared)
    * [Setup: AWS Services](#setup-aws-services)
    * [Setup: Direnv](#setup-direnv)
      * [Helpful variables for `.envrc.local`](#helpful-variables-for-envrclocal)
      * [Troubleshooting direnv & chamber](#troubleshooting-direnv--chamber)
    * [Setup: Run the app](#setup-run-the-app)
    * [Setup: Dependencies](#setup-dependencies)
      * [Setup: Pre-Commit](#setup-pre-commit)
        * [Pre-Commit Troubleshooting (Manual): Process hanging on install hooks](#pre-commit-troubleshooting-manual-process-hanging-on-install-hooks)
        * [Pre-Commit Troubleshooting (Nix): SSL: CERTIFICATE VERIFY FAILED](#pre-commit-troubleshooting-nix-ssl-certificate-verify-failed)
      * [Setup: Database](#setup-database)
      * [Setup: Server](#setup-server)
        * [Server Dependencies](#server-dependencies)
      * [Setup: MilMove Local Client](#setup-milmove-local-client)
  * [Other Possible Setups](#other-possible-setups)
    * [Setup: Office Local client](#setup-office-local-client)
    * [Setup: Admin Local client](#setup-admin-local-client)
    * [Setup: DPS user](#setup-dps-user)
    * [Setup: Orders Gateway](#setup-orders-gateway)
    * [Setup: Prime API](#setup-prime-api)
* [Development](#development)
  * [Makefile](#makefile)
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
  * [GoLand](#goland)
    * [Goland: Nix](#goland-nix)
  * [Storybook](#storybook)
  * [Troubleshooting](#troubleshooting)
    * [Postgres Issues](#postgres-issues)
    * [Development Machine Timezone Issues](#development-machine-timezone-issues)
    * [Linters & Pre-commit Hooks](#linters--pre-commit-hooks)
  * [Manual Redeploys and Other Helpful Information in an Emergency](#manual-redeploys-and-other-helpful-information-in-an-emergency)
  * [PII Best Practices](#pii-best-practices)
    * [More about content dispositions](#more-about-content-dispositions)
    * [More about browser settings](#more-about-browser-settings)

<!-- Regenerate with "pre-commit run -a markdown-toc" -->

<!-- tocstop -->

## Overview

Please check the [MilMove Development Documentation](https://transcom.github.io/mymove-docs/docs) for details on the project itself.

## Supported Browsers

As of 3/6/2018, DDS has confirmed that support for IE is limited to IE 11 and Edge or newer versions. Currently, the intention is to encourage using Chrome and Firefox instead, with specific versions TBD. Research is incomplete on mobile browsers, but we are assuming support for iOS and Android. For more information please read [ADR0016 Browser Support](./docs/adr/0016-Browser-Support.md).

## Login.gov

You'll need accounts for login.gov and the login.gov sandbox.  These will
require two-factor authentication, so have your second factor (one of: phone,
authentication app, security key, CAC) on hand.  To create an account at
login.gov, use your regular `truss.works` email and follow [the official
instructions](https://login.gov/help/creating-an-account/how-to-create-an-account/).
To create an account in the sandbox, follow the same instructions, but [in the
sandbox server](https://idp.int.identitysandbox.gov/sign_up/enter_email).  Do
_not_ use your regular email address in the sandbox.

### Creating alternative users with the same email address

You can use the plus sign `+` to create a new Truss email address.
`name+some_string@truss.works` will be treated as a new address, but will be
routed to your `name@truss.works` email automatically. Don't use this for the
office-side of account creation. It's helpful to use these types of accounts for
the customer-side accounts.

## Project Layout

All of our code is intermingled in the top level directory of `mymove`. Here is an explanation of what some of these
directories contain:

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

## Application Setup

Note: These instructions are a living document and often fall out-of-date.
If you run into anything that needs correcting or updating, please create
a PR with those changes to help those coming after you.

There are two main ways we have for setting up local development:

* Using `nix` with a bit of `homebrew`
* Using primarily only `homebrew`

Both need a bit of base setup before, but then you can follow whichever path you prefer after that. There are also a few parts that may be shared between both setups.

### Setup: Base Setup

There are a number of things you'll need at a minimum to be able to work with this project.

#### Homebrew

We use [Homebrew](https://brew.sh) to manage a few of the packages we need for this project.

Whether or not you already have Homebrew installed, you'll need to make sure it's
up to date and ready to brew:

```shell
SKIP_LOCAL=true /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/trussworks/fresh-brew/main/fresh-press)"
```

If you're using the Fish shell, run this command:

```shell
SKIP_LOCAL=true bash (curl -fsSL https://raw.githubusercontent.com/trussworks/fresh-brew/main/fresh-press | psub)
```

See the [trussworks/fresh-brew repo](https://github.com/trussworks/fresh-brew)
for more information.

#### Setup: Git

Use your work email when making commits to our repositories. The simplest path to correctness is setting global config:

```shell
git config --global user.email "trussel@truss.works"
```

```shell
git config --global user.name "Trusty Trussel"
```

If you drop the `--global` flag these settings will only apply to the current repo. If you ever re-clone that repo or
clone another repo, you will need to remember to set the local config again. You won't. Use the global config. :-)

For web-based Git operations, GitHub will use your primary email unless you choose "Keep my email address private".
If you don't want to set your work address as primary, please
[turn on the privacy setting](https://github.com/settings/emails).

Note if you want use HTTPS instead of SSH for working with git, since we want 2-factor-authentication enabled, you need
to [create a personal access token](https://gist.github.com/ateucher/4634038875263d10fb4817e5ad3d332f) and use that as
your password.

#### Setup: Project Checkout

You can checkout this repository by running

```shell
git clone git@github.com:transcom/mymove.git
```

Please check out the code in a directory like `~/Projects/mymove`. You can check the code out anywhere EXCEPT inside your `$GOPATH`. As an example:

```shell
mkdir -p ~/Projects && cd ~/Projects
```

```shell
git clone git@github.com:transcom/mymove.git
```

```shell
cd mymove
```

You will then find the code at `~/Projects/mymove`.

#### Setup: Editor Config

[EditorConfig](http://editorconfig.org/) allows us to manage editor configuration (like indent sizes) with a
[file](https://github.com/transcom/ppp/blob/master/.editorconfig) in the repo. Install the appropriate plugin in your
editor to take advantage of that if you wish.

### Setup: Nix

If you need help with this setup, you can ask for help in the
[Truss slack #code-nix channel](https://trussworks.slack.com/archives/C01KTH6HP7D).

1. [Initial Setup](#nix-initial-setup)
1. [Clean Up Local Env](#nix-clean-up-local-env)
1. [Install Dependencies](#nix-installing-dependencies)
1. [Run the app](#setup-run-the-app)

#### Nix: Initial Setup

1. First read the overview in the
   [Truss Engineering Playbook](https://github.com/trussworks/Engineering-Playbook/tree/main/developing/nix).
1. Follow the installation instructions in the playbook.

#### Nix: Clean Up Local Env

This section is only if you had previously set up any of these tools/packages. It is also optional, with the following
the caveat of this note:

:warning: NOTE: If you need any of the packages/tools for other things that you won't use `nix` for, you can set things
up so that they both work side by side, but you'll just have to set up your `PATH` properly. And even then, there may be
other steps necessary which aren't documented here.

1. Disable or uninstall `nodenv`, `asdf` or any other version switchers for `mymove`.

   1. `nodenv`:
      1. TLDR (disable only): remove `eval "$(nodenv init -)"` from `.zshrc` (or your shell's config file)
      1. Full instructions: [Uninstalling nodenv](https://github.com/nodenv/nodenv#uninstalling-nodenv)
   1. `asdf`:

      1. See [Remove asdf](https://asdf-vm.com/#/core-manage-asdf?id=remove)
      1. Remove setting of `GOPATH` and putting `GOPATH` in `PATH` in `.zshrc` (or your shell's config file). Looks
         something like this:

         ```shell
         export GOPATH=~/dev/go
         export PATH=$(go env GOPATH)/bin:$PATH
         ```

#### Nix: Installing Dependencies

1. Install a few MilMove dependencies:

   ```shell
   nix-env -i aws-vault chamber direnv bash
   ```

1. [Set up AWS services](#setup-aws-services)

1. Configure direnv:

   1. [Set up direnv](#setup-direnv)
   1. In `.zshrc` (or the relevant one for you), the `nix` setup line (inserted by the `nix` installation) needs to run
      before the `direnv` hook setup.

1. Run `./nix/update.sh`

   1. NOTE: If the nix dependencies change, you should see a warning from direnv:

   ```text
   direnv: WARNING: nix packages out of date. Run nix/update.sh
   ```

1. Run

   ```shell
   make deps_nix
   ```

   1. This will install some things like `pre-commit` hooks, `node_modules`, etc. You can see
      [Setup: Dependencies](#setup-dependencies) for more info on some of the parts.

### Setup: Manual

1. [Set up AWS services](#setup-aws-services)
1. [Prerequisites](#manual-prerequisites)
1. [Set up direnv](#setup-direnv)
1. [Run the app](#setup-run-the-app)

#### Manual: Prerequisites

We have scripts that will install all the dependencies for you, as well as configure your shell file with all the required commands:

```shell
SKIP_CHECKS=true make prereqs
```

This will install everything listed in `Brewfile.local`, as well as Docker.

**Note**: The script might ask you for your macOS password at certain points, like when installing opensc, or when it needs to write to your `/etc/hosts` file.

Once this script is finished, quit and restart your terminal, then complete the
installation:

```shell
make deps
```

This will install `pre-commit` hooks and frontend client dependencies. See [Setup: Dependencies](#setup-dependencies) for more info.

**Note that installing and configuring pre-commit the first time takes about 3 minutes.**

Going forward, feel free to run `make prereqs` or `make deps` as often as you'd like to keep your system up to date. Whenever we update the app to a newer version of Go or Node, all you have to run is `make prereqs` and it will update everything for you.

### Setup: Shared

#### Setup: AWS Services

This project uses AWS services which means you'll need an account to work with parts of it. AWS credentials are managed
via `aws-vault`. Once you have received AWS credentials (which are provided by the infrastructure team), you can follow
these instructions to
[finish setting up AWS](https://dp3.atlassian.net/wiki/spaces/MT/pages/1250066433/0029+AWS+Organization+Authentication).

#### Setup: Direnv

For managing local environment variables, we're using [direnv](https://direnv.net/).

1. Run

    ```shell
    direnv allow
    ```

    1. This will load up the `.envrc` file. It should complain that you have missing variables. We'll fix that next.

To fix the missing variables issue, you can do one of the following things:

* Let `direnv` get secret values with `chamber`. To enable this, run:

  ```shell
  cp .envrc.chamber.template .envrc.chamber
  ```

  * **Note** that this method does not work for users of the `fish` shell unless you replace `direnv allow` with

    ```shell
    direnv export fish | source
    ```

  * **Note also** if you have a very poor internet connection, this method may be
  problematic to you.

* An alternative is to add a `.envrc.local` file. Then run:

  ```shell
  DISABLE_AWS_VAULT_WRAPPER=1 AWS_REGION=us-gov-west-1 aws-vault exec transcom-gov-dev -- chamber env app-devlocal >> .envrc.local
  ```

* If you don't have access to `chamber`, you can also run

  ```shell
  touch .envrc.local
  ```

  then add any values that the output from `direnv` asks you to define.

##### Helpful variables for `.envrc.local`

* Increase concurrency of `golangci-lint`; defaults to 6 on dev machines and to 1 in CircleCI.

  ```shell
  export GOLANGCI_LINT_CONCURRENCY=8
  ```

* Enable go code debugging in goland

  ```shell
  export GOLAND=1
  ```

* Silence SQL logs locally; we default this to be true in `.envrc`

  ```shell
  export DB_DEBUG=0
  ```

##### Troubleshooting direnv & chamber

Make sure you have the latest version of Chamber that supports the `env` command
option. You may run into the following error if the version of Chamber you have
installed does not support `env`. The error presents itself because of the
`chamber` commands that `direnv` runs as part of the `.envrc.*` files shown
above.

```shell
>_ cd mymove
direnv: loading .envrc.chamber
Error: unknown command "env" for "chamber"
Run 'chamber --help' for usage.
```

#### Setup: Run the app

**If this is your very first time setting up this project, you'll need to launch Docker first, follow the prompts to allow macOS to open it, and agree to Docker's terms of service.**

You might also need to launch Docker if you restarted your computer and you configured Docker to not automatically launch after a restart.

Once Docker is up and running, the following commands will get `mymove` running on your machine.

1. Run the backend server

   ```shell
   make server_run
   ```

   This command also ensures the database is up and running and that the
   latest migrations are applied. See [Setup: Database](#setup-database) and
   [Setup: Server](#setup-server) for more details.

1. Run the frontend client **in a separate terminal tab**

   ```shell
   make client_run
   ```

   This will ensure the frontend dependencies are installed and will
   automatically launch the browser and open the app at milmovelocal:3000.
   See [Setup: MilMove Local Client](#setup-milmove-local-client) for more details.

#### Setup: Dependencies

This step will check your system for any setup issues. Then it will ensure that you have installed `pre-commit`
and go on to install the client (javascript) and server (golang) dependencies for you. If you are interested in
more details, you can look at the sections under this one, but it's not required.

##### Setup: Pre-Commit

Part of the `pre-commit` setup run by the `make deps` or `make deps_nix` commands.
They in turn run

```shell
pre-commit install
```

to install a pre-commit hook into `./git/hooks/pre-commit`. This must be done so
that the hook will check files you are about to commit to the repository.

Next it installs the `pre-commit` hook libraries with

```shell
pre-commit install-hooks
```

If you ever want to run the `pre-commit` hooks for all files, you can run

```shell
pre-commit run -a
```

though before you can do that, you'll need to have installed the `javascript` dependencies and generated some `golang`
code from Swagger files. Once you've finished setting up your project locally, you should be good to go. If you want
to skip ahead and be able to run `pre-commit` checks since now, you can run

```shell
make pre_commit_tests
```

or

```shell
make server_generate client_deps && pre-commit run -a
```

###### Pre-Commit Troubleshooting (Manual): Process hanging on install hooks

If any pre-commit commands (or `make deps`) result in hanging or incomplete
installation, remove the pre-commit cache and the `.client_deps.stamp` and try again:

```shell
rm -rf ~/.cache/pre-commit
rm .client_deps.stamp
```

###### Pre-Commit Troubleshooting (Nix): SSL: CERTIFICATE VERIFY FAILED

This can happen because of the way certs need to be handled in this project and `nix`. To get around this issue, you
can try running:

```shell
NIX_SSL_CERT_FILE=$HOME/.nix-profile/etc/ssl/certs/ca-bundle.crt <pre-commit related command>
```

E.g.

```shell
NIX_SSL_CERT_FILE=$HOME/.nix-profile/etc/ssl/certs/ca-bundle.crt pre-commit install-hooks
```

##### Setup: Database

You will need to setup a local database before you can begin working on the local server / client. Docker will need to
be running for any of this to work.

1. Creates a PostgreSQL docker container for dev, if it doesn't exist already, and starts/runs it.

    ```shell
    make db_dev_run
    ```

1. Runs all existing database migrations for dev database, which does things like creating table structures, etc.
   You will run this command again anytime you add new migrations to the app (see below for more).

    ```shell
    make db_dev_migrate
    ```

You can validate that your dev database is running by running

```shell
psql-dev
```

This puts you in a PostgreSQL shell. To show all the tables, type

```shell
\dt
```

If you want to exit out of the PostgreSQL shell, type

```shell
\q
```

If you are stuck on this step you may need to see the section on Troubleshooting.

##### Setup: Server

This step installs dependencies, then builds and runs the server using `gin`, which is a hot reloading go server.
It will listen on port `8080` and will rebuild the actual server any time a go file changes.

```shell
make server_run
```

To have hot reloading of the entire application (at least for the customer side), pair the above with

```shell
make client_run
```

In rare cases, you may want to run the server standalone, in which case you can run

```shell
make server_run_standalone
```

This will build both the client and the server and this invocation can be relied upon to be serving the client JS on
its own rather than relying on webpack doing so. You can run this without running `make client_run` and the whole app
should work.

###### Server Dependencies

Dependencies are managed by [go modules](https://github.com/golang/go/wiki/Modules). New dependencies are automatically
detected in import statements and added to `go.mod` when you run

```shell
go build
```

or

```shell
go run
```

You can also manually edit `go.mod` as needed.

If you need to add a Go-based tool dependency that is otherwise not imported by our code, import it in
`pkg/tools/tools.go`.

After importing _any_ go dependency it's a good practice to run

```shell
go mod tidy
```

which prunes unused dependencies and calculates dependency requirements for all possible system architectures.

##### Setup: MilMove Local Client

Commands in this section:

```shell
make client_build
```

and

```shell
make client_run
```

These will start the webpack dev server, serving the frontend on port 3000. If paired with

```shell
make server_run
```

then the whole app will work, the webpack dev server proxies all API calls through to the server.

If both the server and client are running, you should be able to view the Swagger UI at
<http://milmovelocal:3000/swagger-ui/internal.html>. If it does not, try running

```shell
make client_build
```

(this only needs to be run the first time).

Dependencies are managed by `yarn`. To add a new dependency, use

```shell
yarn add
```

### Other Possible Setups

The instructions so far have been for getting the project up and running, but focused on the client/customer side.
There are more things you can set up in the following sections.

#### Setup: Office Local client

1. Ensure that you have a test account which can log into the office site. To load test data, run:

    ```shell
    make db_dev_e2e_populate
    ```

1. Run

    ```shell
    make office_client_run
    ```

1. Log into "Local Sign In" and either select a pre-made user or use the button to create a new user

#### Setup: Admin Local client

Run

```shell
make admin_client_run
````

#### Setup: DPS user

1. Ensure that you have a login.gov test account
2. Log into [MilMove Devlocal Auth](http://milmovelocal:3000/devlocal-auth/login) and create a new DPS user from the
   interface.

#### Setup: Orders Gateway

Nothing to do.

#### Setup: Prime API

The API that the Prime will use is authenticated via mutual TSL so there are a few things you need to do to interact
with it in a local environment.

1. Make sure that the `primelocal` alias is setup for localhost
    1. Check your `/etc/hosts` file for an entry for `primelocal`.
2. Run

    ```shell
    make server_run
    ```

3. Access the Prime API using the devlocal-mtls certs. There is a script that shows you how to do this with curl
   at `./scripts/prime-api`. For instance to call the `move-task-orders` endpoint, run

    ```shell
    ./scripts/prime-api move-task-orders
    ```

## Development

### Makefile

The primary way to interact with the project is via the `Makefile`. The `Makefile` contains a number of handy
targets (you can think of these as commands) that make interacting with the project easier. Each target manages
its own dependencies so that you don't have to. This is how you'll do common tasks like build the project, run
the server and client, and manage the database.

The fastest way to get familiar with the `Makefile` is to use the command `make help`. You can also type `make`
and it will default to calling `make help` target on your behalf.  The `Makefile` is important to this project
so take the time to understand what it does.

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

Internal services (i.e. endpoints only intended for use by the React client) are
defined in `swagger-def/internal.yaml` and served from the value of the
`basePath:` stanza at the root of the generated `./swagger/internal.yaml` file.
**Internal endpoints are not intended for use by external clients**.

The Orders Gateway's API is defined in the file `swagger-def/orders.yaml` and
served from the value of the `basePath:` stanza at the root of the generated
`./swagger/orders.yaml` file.

The Admin API is defined in the file `swagger-def/admin.yaml` and served from
the value of the `basePath:` stanza at the root of the generated
`./swagger/admin.yaml` file.

The Prime API is defined in the file `./swagger-def/prime.yaml` and served from
the value of the `basePath:` stanza at the root of the generated
`./swagger/prime.yaml` file.

You can view the documentation for the following APIs (powered by Swagger UI) at
the following URLS with a local client and server running:

* internal API: [http://milmovelocal:3000/swagger-ui/internal.html](http://milmovelocal:3000/swagger-ui/internal.html)
* admin API: [http://milmovelocal:3000/swagger-ui/admin.html](http://milmovelocal:3000/swagger-ui/admin.html)
* GHC API: [http://milmovelocal:3000/swagger-ui/ghc.html](http://milmovelocal:3000/swagger-ui/ghc.html)
* Prime API: [http://primelocal:3000/swagger-ui/prime.html](http://primelocal:3000/swagger-ui/prime.html)

For more information on _API / Swagger_ definitions, please review the README
documentation found in the `./swagger/README.md` and `./swagger-def/README.md`
files.

### Testing

There are a few handy targets in the Makefile to help you run tests:

* `make client_test`: Run front-end testing suites.
* `make server_test`: Run back-end testing suites. [Additional info for running go tests](https://transcom.github.io/mymove-docs/docs/dev/testing/running-tests/run-go-tests)
* `make e2e_test`: Run end-to-end testing suite. [Additional info for running E2E tests](https://transcom.github.io/mymove-docs/docs/dev/testing/running-tests/run-e2e-tests)
  * Note: this will not necessarily reflect the same results as in the CI
  environment, run with caution. One is your `.envrc` is going to
  populate your dev environment with a bunch of values that `make e2e_test_docker`
  won't have.
* `make e2e_test_docker`: Run e2e testing suite in the same docker container as is run in CircleCI.
  * Note: this runs with a full clean/rebuild, so it is not great for fast iteration.
  Use `make e2e_test` to pick individual tests from the Cypress UI.
* `make test`: Run e2e, client- and server-side testing suites.
* `yarn test:e2e`: Useful for debugging. This opens the cypress test runner
against your already running dev servers and inspect/run individual e2e tests.
  * Note: You must already have the servers running for this to work!

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

* Read [Querying the Database Safely](https://transcom.github.io/mymove-docs/docs/dev/contributing/backend/Backend-Programming-Guide/#querying-the-database-safely) to prevent SQL injections! *

A few commands exist for starting and stopping the DB docker container:

* `make db_run`: Starts the DB docker container if one doesn't already exist
* `make db_destroy`: Stops and removes the DB docker container

#### Dev DB Commands

There are a few handy targets in the Makefile to help you interact with the dev
database. During your day-to-day, the only one you will typically need regularly
is `make db_dev_e2e_populate`. The others are for reference, or if something
goes wrong.

* `make db_dev_e2e_populate`: Populates the dev DB with data to facilitate
verification of your work when using the app locally. It seeds the DB with various
service members at different stages of the onboarding process, various office
users, moves, payment requests, etc. The data is defined in the `devseed.go` file.
* `make db_dev_run`: Initializes a new database if it does not exist and runs it,
or starts the previously initialized Docker container if it has been stopped.
You typically only need this after a computer restart, or if you manually quit
Docker or otherwise stopped the DB.
* `make db_dev_create`: Waits to connect to the DB and will create a DB if one
doesn't already exist (this is automatically run as part of `db_dev_run`).
* `make db_dev_fresh`: Destroys your database container, runs the DB, and
applies the migrations. Useful if you want to start from scratch when the DB is
not working properly. This runs `db_dev_reset` and `db_dev_migrate`.
* `make db_dev_migrate_standalone`: Applies database migrations against your
running database container but will not check for server dependencies first.


#### Test DB Commands

These commands are available for the Test DB. You will rarely need to use these
individually since the commands to run tests already set up the test DB properly.
One exception is `make db_test_run`, which you'll need to run after restarting
your computer.

* `make db_test_run`
* `make db_test_create`
* `make db_test_reset`
* `make db_test_migrate`
* `make db_test_migrate_standalone`
* `make db_test_e2e_backup`
* `make db_test_e2e_restore`
* `make db_test_e2e_cleanup`

The test DB commands all talk to the DB over localhost.  But in a docker-only environment (like CircleCI) you may not be able to use those commands, which is why `*_docker` versions exist for all of them:

* `make db_test_run_docker`
* `make db_test_create_docker`
* `make db_test_reset_docker`
* `make db_test_migrate_docker`

#### Migrations

To add new regular and/or secure migrations, see the [database development guide](https://transcom.github.io/mymove-docs/docs/dev/contributing/database/Database-Migrations)

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
}
```

### Documentation

You can view the project's GoDoc on [godoc.org](https://godoc.org/github.com/transcom/mymove).

Alternatively, run the documentation locally using:

```shell
# within the project's root dir
$ godoc -http=:6060
```

Then visit <http://localhost:6060/pkg/github.com/transcom/mymove/>

### GoLand

GoLand supports
[attaching the debugger to a running process](https://blog.jetbrains.com/go/2019/02/06/debugging-with-goland-getting-started/#debugging-a-running-application-on-the-local-machine),
however this requires that the server has been built with specific flags. If you wish to use this feature in
development add the following line `export GOLAND=1` to your `.envrc.local`. Once the server starts follow the steps
outlined in the article above and you should now be able to set breakpoints using the GoLand debugger.

#### Goland: Nix

To get Goland to play nicely with `nix`, there's a few things you can set up:

* Update `GOROOT` to `/nix/var/nix/profiles/mymove/bin/go`
  * Note that once you add it, Goland will resolve it to the actual path (the one above is a link), so it’ll look
    something like `/nix/store/rv16prybnsmav8w1sqdgr80jcwsja98q-go-1.17.7/bin/go`
* Update `GOPATH` to point to the `.gopath` dir in the `mymove` repo
  * You may need to create the `.gopath` dir yourself.
* Update Node and NPM:
  * Node interpreter: `/nix/var/nix/profiles/mymove/bin/node`
  * Package manager:
    * This might be fixed automatically, but if not, you can point it `/nix/var/nix/profiles/mymove/bin/yarn`
    * Similar to `GOROOT`, it will resolve to something that looks like
      `/nix/store/cnmxp5isc3ck1bm11zryy8dnsbnm87wk-yarn-1.22.10/libexec/yarn`

### Storybook

We use [Storybook](https://storybook.js.org) for reviewing our
component library. The current components are deployed to
[https://storybook.dp3.us](https://storybook.dp3.us) after each build
of the master branch.

Each PR saves storybook as an artifact in CircleCI. Find the
`build_storybook` task and then go to the "ARTIFACTS" tab. Find the
link to `storybook/index.html` and click on it.

### Troubleshooting

* Random problems may arise if you have old Docker containers running. Run `docker ps` and if you see containers unrelated to our app, consider stopping them.
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

We use a number of linters for formatting, security and error checking. Please see the [pre-commit documentation](https://transcom.github.io/mymove-docs/docs/dev/contributing/code-analysis/run-pre-commit-hooks) for a list of linters and troubleshooting tips.

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
