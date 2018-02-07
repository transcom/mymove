# Personal Property Prototype

[![Build status](https://img.shields.io/circleci/project/github/transcom/mymove/master.svg)](https://circleci.com/gh/transcom/mymove/tree/master)

This repository contains the application source code for the Personal Property Prototype, a possible next generation version of the Defense Personal Property System (DPS). DPS is an online system managed by the U.S. [Department of Defense](https://www.defense.gov/) (DoD) [Transportation Command](http://www.ustranscom.mil/) (USTRANSCOM) and is used by service members and their families to manage household goods moves.

This prototype was built by a [Defense Digital Service](https://www.dds.mil/) team in support of USTRANSCOM's mission.

## Development

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
`cmd`: The location of main packages for any go binaries we build (right now, just webserver) \
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
* Run `bin/prereqs` and install everything it tells you to. _Do not configure postgres to automatically start at boot time!_
* Run `make deps`.
* [EditorConfig](http://editorconfig.org/) allows us to manage editor configuration (like indent sizes,) with a [file](https://github.com/transcom/ppp/blob/master/.editorconfig) in the repo. Install the appropriate plugin in your editor to take advantage of that.

### Setup: Database

You will need to setup a local database before you can begin working on the local server / client. Docker will need to be running for any of this to work.

1. `make db_dev_migrate`: Creates a postgres docker container if you haven't made one yet and runs all existing database migrations, which do things like creating table structures, etc. You will run this command again anytime you add new migrations to the app (see below for more)

You can validate that your dev database is running by running `bin/psql-dev`. This puts you in a postgres shell. Type `\dt` to show all tables, and `\q` to quit.

### Setup: Server

1. `make server_run`: installs dependencies, then builds and runs the server using `gin`, which is a hot reloading go server. It will listen on 8080 and will rebuild the actual server any time a go file changes. Pair this with `make client_run` to have hot reloading of the entire application.

In rare cases, you may want to run the server standalone, in which case you can run `make server_run_standalone`. This will build both the client and the server and this invocation can be relied upon to be serving the client js on its own rather than relying on webpack doing so as when you run `make client_run`. You can run this without running `make client_run` and the whole app should work.

You can verify the server is working as follows:

`> curl http://localhost:8080/api/v1/issues --data "{ \"description\": \"This is a test issue\"}" -H "Content-Type: application/json"`

from which the response should be like

`{"id":"d5735bc0-7553-4d80-a42d-ea1e50bbcfc4", "description": "This is a test issue", "created_at": "2018-01-04 14:47:28.894988", "updated_at": "2018-01-04 14:47:28.894988"}`

Dependencies are managed by [dep](https://github.com/golang/dep). New dependencies are automatically detected in import statements. To add a new dependency to the project, import it in a source file and then run `dep ensure`

### Setup: Client

1. `make client_run`

The above will start the webpack dev server, serving the frontend on port 3000. If paired with `make server_run` then the whole app will work, the webpack dev server proxies all api calls through to the server.

Dependencies are managed by yarn. To add a new dependency, use `yarn add`

### TSP Award Queue

This background job is built as a separate binary which can be built using
`make tsp_build` and run using `make tsp_run`.

### API / Swagger

The api is defined in a single file: ./swagger.yaml and served at /api/v1/swagger.yaml. This file is the single source of truth for the contract between the client and server.

You can view the API's documentation (powered by Swagger UI) at [http://localhost:8081/api/v1/docs](http://localhost:8081/api/v1/docs) when a local server is running.

### Testing

There are a few handy targets in the Makefile to help you run tests:

* `make client_test`: Run frontend testing suites.
* `make server_test`: Run backend testing suites.
* `make test`: Run both client- and server-side testing suites.

### Database

#### Dev Commands

There are a few handy targets in the Makefile to help you interact with the dev database:

* `make db_dev_run`: Initializes a new database if it does not exist and runs it, or starts the previously initialized Docker container if it has been stopped.
* `make db_dev_reset`: Destroys your database container. Useful if you want to start from scratch.
* `make db_dev_migrate`: Applies database migrations against your running database container.
* `make db_dev_migrate_down`: reverts the most recently applied migration by running the down migration

#### Migrations

If you need to change the database schema, you'll need to write a migration.

Creating a migration:

Use soda (a part of [pop](https://github.com/markbates/pop/)) to generate migrations. In order to make using soda easy, a wrapper is in `./bin/soda` that sets the go environment and working directory correctly.

If you are generating a new model, use `./bin/gen_model model-name column-name:type column-name:type ...`. id, created_at and updated_at are all created automatically.

If you are modifying an existing model, use `./bin/soda generate migration migration-name` and add the [Fizz instructions](https://github.com/markbates/pop/blob/master/fizz/README.md) yourself to the created files.

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

### Troubleshooting

* Random problems may arise if you have old Docker containers running. Run `docker ps` and if you see containers unrelated to our app, consider stopping them.
* If you have problems connecting to postgres, or running related scripts, make sure you aren't already running a postgres daemon. You can check this by typing `ps aux | grep postgres` and looking for existing processes.
* If you happen to have installed pre-commit in a virtual environment not with brew, running bin/prereqs will not alert you. You may run into issues when running `make deps`. To install pre-commit: `brew install pre-commit`.
