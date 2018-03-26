# Isolate Test Access to Database

**User Story:** [155076695](https://www.pivotaltracker.com/story/show/155076695)

Our tests currently execute against a single database and do not clean up after themselves. As each tests runs, any mutations made to the database persist. This lack of isolation between tests can lead to test reliability issues.

Our priorities are, in order:

* Reliable tests (tests that only sometimes pass don't provide value).
* Developer experience (we don't want to add a lot of boilerplate to each test that obscures its true intent.)
* Minimizing test run duration (we don't want to optimize for speed this early in the project.)

## Considered Alternatives

* **Wrap test execution in a database transaction** Execute tests that leverage the database within SQL transactions, calling ROLLBACK() at the end of each test to prevent any mutations from persisting after the test completes.

* **Truncate tables between tests** Delete all rows from the database using `TRUNCATE` after each test executes and ensure that only a single test is communicating with the test database at one time.

## Decision Outcome

* Chosen Alternative: **Truncate tables between tests**

* The main driver behind this decision was simplicity. We anticipate using transactions within our application (although we aren't already), and wrapping tests in transactions would complicate matters as transactions can't be simply nested. It is possible to emulate nested transactions using [`SAVEPOINT`](https://www.postgresql.org/docs/8.1/static/sql-savepoint.html), but this is not supported by Pop or sqlx and pursuing this route would require us to write and maintain additional non-trivial code in the project.

* `go test` executes tests from a single package serially unless there are tests that are explicitly marked as able to run in parallel using [`t.Parallel()`](https://golang.org/pkg/testing/#T.Parallel). We do, however, have multiple packages with tests that use the database (currently `models` and `handlers`), so we will additionally need to pass `-test.parallel 1` to `go test` so that packages as well as individual tests are executed serially. Otherwise, multiple tests using the database will encounter collisions.

  Executing each package's tests serially will increase the time required for a project-wide test run, but at this time this is a reasonable trade-off as the suite duration time is currently acceptable. If speed of tests execution become a concern in the future, there are a few approaches we could pursue without too much effort, including having a test database per package with tests hitting the database.

* This strategy requires the ability to run code before or after each test, which will be addressed in a future ADR.

* [Pop uses `TRUNCATE`](https://github.com/gobuffalo/pop/blob/9f77e19c929eda4c13f525296fe751a90de86619/postgresql.go#L232-L248) to implement [`TruncateAll()`](https://godoc.org/github.com/gobuffalo/pop#Connection.TruncateAll).

## Pros and Cons of the Alternatives

### Wrap test execution in a database transaction

* `+` Allows multiple tests to interact with the database at once without interfering with each other.
* `-` Complicates testing code that itself uses transactions (see discussion re: transaction nesting above).
* `-` Requires writing additional testing helpers to manage transactions.
* `-` Requires that all tests use the correct `*pop.Connection` that has an open transaction, which is difficult to verify without manual code inspection.
* `+` Is considered to be faster than truncation (we did not perform benchmarking as speed was not a top priority).
