# How To Run Go Tests

## Run All Go Tests

```console
$ make server_test
```

## Run Specific Tests

All of the commands in this section assume that `test_db` is setup properly. This can be done using:

```console
$ make db_test_reset
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