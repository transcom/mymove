# How To Run End to End (Cypress) Tests

## Run All End to End Tests

The following command to open the Cypress UI, from which you can choose "Run all Specs":

```console
$ make e2e_test
```

If you instead would like to run the tests in the console, use the following command:

```console
$ make e2e_test_docker
```

## Run Specific Tests

### Using the Cypress UI

To run a single test, use `make e2e_test` and choose the tests you wish to run from the Cypress UI.

If you have already run tests in the current database, you will need to reset the database to known good state:

```console
$ make db_e2e_reset
```

### Run End to End Tests in a File

```console
$ yarn cypress run --spec cypress/integration/path/to/file.jsx
```

If you have already run tests in the current database, you will need to reset the database to known good state:

```console
$ make db_e2e_reset
```

### Run End to End Tests with Docker

To run just the office tests:

```console
$ SPEC=cypress/integration/office/**/* make e2e_test_docker
```
