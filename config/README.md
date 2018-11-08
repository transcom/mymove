# Configuration

## Container Definitions

The `*container-definition*` files define how the ECS containers are configured. They are primarily used to set non-secret environment variables.

## Database

The database is configured using environment variables. The environments "development", "test", and "container" each
have specific modifications to the database connection. Specifically, "test" environment will always use the database
name `test_db` no matter what the value of `$DB_NAME` is set to. The "container" environment will require that SSL
mode is required.

## TLS cert/key (optional)

The `devlocal-https.*` files are a self-signed TLS cert/key pair. They are a [snake oil](https://en.wikipedia.org/wiki/Snake_oil_(cryptography)) certificate used to locally run the webserver during development. They are included as a convenience so engineers don't have to generate their own.
