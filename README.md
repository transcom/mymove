# Personal Property Prototype

[![Build status](https://img.shields.io/circleci/project/github/transcom/ppp/master.svg)](https://circleci.com/gh/transcom/ppp/tree/master)

This repository contains the application source code for the Personal Property Prototype, a possible next generation version of the Defense Personal Property System (DPS). DPS is an online system managed by the U.S. [Department of Defense](https://www.defense.gov/) (DoD) [Transportation Command](http://www.ustranscom.mil/) (USTRANSCOM) and is used by service members and their families to manage household goods moves.

This prototype was built by a [Defense Digital Service](https://www.dds.mil/) team in support of USTRANSCOM's mission.

## Development

### Project location

All of Go's tooling expects Go code to be checked out in a specific location. Please read [How to Write Go code](https://golang.org/doc/code.html) for a full explaination, but if you just want to get started, then decide where you want all your go code to live and configure the GOPATH envionment variable accordingly. For example, if you want your go code to live at ~/code/go, then in your .bash_profile you should add the following two lines:

```bash
export GOPATH=~/code/go
PATH=~/code/go/bin:$PATH
```

Once that's done, and you've re-sourced your profile, you can checkout this repository by running `go get github.com/transcom/mymove` You will then find the code at `$GOPATH/src/github.com/transcom/mymove`

If you have already checked out the code somewhere else, you can just move it to be in the above locaiton and everything will work correctly.

### Prerequisites

* Install Go with Homebrew. Make sure you do not have other installations.
* Run `bin/prereqs` and install everything it tells you to. Then run `make deps`.
* [EditorConfig](http://editorconfig.org/) allows us to manage editor configuration (like indent sizes,) with a [file](https://github.com/transcom/ppp/blob/master/.editorconfig) in the repo. Install the appropriate plugin in your editor to take advantage of that.

### Setup: Database

You will need to setup a local database before you can begin working on the local server / client. Docker will need to be running for any of this to work.

1. `make db_dev_init`: initializes a Docker container with a Postgres database.
1. `make db_dev_migrate`: runs all existing database migrations, which do things like creating table structures, etc.
1. You can validate that your dev database is running by running `bin/psql-dev`. This puts you in a postgres shell. Type `\dt` to show all tables, and `\q` to quit.

### Setup: Server

`make server_run`: installs dependencies and builds both the client and the server, then runs the server.

For faster development, use `make server_run_dev`. This builds and runs the server but skips updating dependences and re-building the client. Those tasks can be accomplished as needed with `make server_deps` and `make client_build`

You can verify the server is working as follows:

`> curl http://localhost:8080/api/v1/issues --data "{ \"body\": \"This is a test issue\"}"`

from which the response should be like

`{"id":"d5735bc0-7553-4d80-a42d-ea1e50bbcfc4", "body": "This is a test issue", "created_at": "2018-01-04 14:47:28.894988", "updated_at": "2018-01-04 14:47:28.894988"}`

Dependencies are managed by [dep](https://github.com/golang/dep). New dependencies are automatically detected in import statements. To add a new dependency to the project, import it in a source file and then run `dep ensure`

### Setup: Client

`make server_run`
`make client_run_dev`

The above will start the server running and starts the webpack dev server, proxied to our running go server.

Dependencies are managed by yarn

### API

The api is defined in a single file: ./swagger.yaml and served at /api/v1/swagger.yaml. it is the single source of truth for what the API contract between client and server should be.

### Testing

There are a few handy targets in the Makefile to help you run tests:

* `make client_test`: Run frontend testing suites.
* `make server_test`: Run backend testing suites.
* `make test`: Run both client- and server-side testing suites.

### Database

#### Dev Commands

There are a few handy targets in the Makefile to help you interact with the dev database:

* `make db_dev_init`: Initializes a new postgres Docker container with a test database and runs it. You must do this before any other database operations.
* `make db_dev_run`: Starts the previously initialized Docker container if it has been stopped.
* `make db_dev_reset`: Destroys your database container. Useful if you want to start from scratch.
* `make db_dev_migrate`: Applies database migrations against your running database container.

#### Migrations

If you need to change the database schema, you'll need to write a migration.

Creating a migration:

Use soda (a part of [pop](https://github.com/markbates/pop/)) to generate migrations. In order to make using soda easy, a wrapper is in `./bin/soda` that sets the go environment and working directory correctly.

If you are generating a new model, use `./bin/soda generate model model-name column-name:type column-name:type ...`. id, created_at and updated_at are all created automatically.

If you are modifying an existing model, use `./bin/soda generate migration migration-name` and add the [Fizz instructions](https://github.com/markbates/pop/blob/master/fizz/README.md) yourself to the created files.

Running migrations:

1. Use `make db_dev_migrate` to run migrations against your local dev environment. Production migrations TBD.

### Troubleshooting

* Random problems may arise if you have old Docker containers running. Run `docker ps` and if you see containers unrelated to our app, consider stopping them.
* If you have problems connecting to postgres, or running related scripts, make sure you aren't already running a postgres daemon. You can check this by typing `ps aux | grep postgres` and looking for existing processes.
* If you happen to have installed pre-commit in a virtual environment not with brew, running bin/prereqs will not alert you. You may run into issues when running `make deps`. To install pre-commit: `brew install pre-commit`.
