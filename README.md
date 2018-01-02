# Personal Property Prototype

[![Build status](https://img.shields.io/circleci/project/github/transcom/ppp/master.svg)](https://circleci.com/gh/transcom/ppp/tree/master)

This repository contains the application source code for the Personal Property Prototype, a possible next generation version of the Defense Personal Property System (DPS). DPS is an online system managed by the U.S. [Department of Defense](https://www.defense.gov/) (DoD) [Transportation Command](http://www.ustranscom.mil/) (USTRANSCOM) and is used by service members and their families to manage household goods moves.

This prototype was built by a [Defense Digital Service](https://www.dds.mil/) team in support of USTRANSCOM's mission.

## Development

### Prerequisites

Run `bin/prereqs` and install everything it tells you to. Then run `make client deps` and `make server deps`.

### Setup: Database

You will need to setup a local database before you can begin working on the local server / client.

1. `make db_dev_init`: initializes a Docker container with a Postgres database.
1. `make db_migrate`: runs all existing database migrations, which do things like creating table structures, etc.

### Setup: Server

`make server_run`: installs dependencies and builds both the client and the server, then runs the server.

For faster development, use `make server_run_dev`. This builds and runs the server but skips updating dependences and re-building the client. Those tasks can be accomplished as needed with `make server_deps` and `make client_build`

You can verify the server is working as follows:

`> curl http://localhost:8080/api/v1/issues --data "{ \"issue\": \"This is a test issue\"}"`

from which the response should be

`{"id":1}`

Dependencies are managed by glide. To add a new dependency:
`GOPATH=/path/to/dp3 glide get new/dependency`

### Setup: Client

`make server_run`
`make client_run_dev`

The above will start the server running and starts the webpack dev server, proxied to our running go server.

Dependencies are managed by yarn

### Database

#### Dev Commands

There are a few handy targets in the Makefile to help you interact with the dev database:

* `make db_dev_init`: Initializes a new postgres Docker container with a test database and runs it. You must do this before any other database operations.
* `make db_dev_run`: Starts the previously initialized Docker container if it has been stopped.
* `make db_dev_reset`: Destroys your database container. Useful if you want to start from scratch.
* `make db_dev_migrate`: Applies database migrations against your running database container.

#### Migrations

If you need to change the datata base schema, you'll need to write a migration. Migrations are handled by [Flyway](https://flywaydb.org), though we've added some wrappers around it.

1. Use `bin/db-generate-migraion $NAME` to create a new migration file. You should supply a descriptive name, such as "create_users_table" or "add_description_column_to_user".
1. Edit the file that was created in the previous step, and write the appropriate SQL to make the needed changes.
1. Use `make db_dev_migrate` to apply the migration and test your changes.
