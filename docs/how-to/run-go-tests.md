# How To Run Go Tests

## Run All Go Tests

```console
$ make server_test
```

## Run Specific Tests

All of the commands in this section assume that `test_db` is setup properly. This can be done using:

```console
$ make db_test_reset && make db_test_migrate
```

## Run Acceptance Tests

If you're adding a feature that requires new or modified configuration, it's a good idea to run acceptance tests against our environments before you merge into master.  You can run acceptance tests against an environment with:

```console
$ TEST_ACC_ENV=experimental make webserver_test
```

This command will first load the variables from the `config/env/*.env` file and then run `chamber exec` to pull the environments from AWS.  You can run acceptance tests for the database, DOD certificates, and Honeycomb through environment variables with `TEST_ACC_DATABASE=1`, `TEST_ACC_DOD_CERTIFICATES=1`, and `TEST_ACC_HONEYCOMB=1`, respectively.

For example to run acceptance tests against staging, including DOD certificate parsing, use:

```console
$ TEST_ACC_ENV=staging TEST_ACC_DOD_CERTIFICATES=1 make webserver_test
```

### Run All Tests in a Single Package

```console
$ go test ./pkg/handlers/internalapi/
```

### Run Tests with Names Matching a String

The following will run any Testify tests that have a name matching `Test_Name` in the `handlers/internalapi` package:

```console
$ go test ./pkg/handlers/internalapi/ -testify.m Test_Name
```

## Run Tests when a File Changes

You'll need to install [ripgrep](https://github.com/BurntSushi/ripgrep) and [entr](http://www.entrproject.org/) (`brew install ripgrep entr`):

```console
$ rg -t go --files | entr -c $YOUR_TEST_COMMAND
```

Here is an example that will run all model tests when any Go file in the project is changed:

```console
$ rg -t go --files | entr -c go test ./pkg/models
```

There is generally no need to be any more specific than `rg -t go`, as watching all `.go` files is plenty fast enough.
