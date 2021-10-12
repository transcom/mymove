# Run tests within transactions

**NOTE:** This ADR updates and supersedes [ADR0010 Isolate Test Access to Database](./0010-isolate-test-access-to-database.md).

Ever since ADR 0010, we've been truncating the DB in between tests as a way to
start each test with a clean state, and to allow tests to run in parallel. Now
that our test suite and number of packages are much larger, the time lost to
cloning and truncating the DB outweighs the initial benefit of simplicity.

Instead, we propose using the [go-txdb](https://github.com/DATA-DOG/go-txdb)
tool to run tests within a transaction, that gets rolled back at the end of the
test. This is much faster than truncating the DB, and the work to implement this
is already done in [PR #6650](https://github.com/transcom/mymove/pull/6650).

As part of configuring go-txdb in PR #6650, we also converted a few packages to
start running tests in transaction, which serve as examples of how to convert
the remaining packages. We also have a [Docs with instructions](https://transcom.github.io/mymove-docs/docs/dev/testing/writing-tests/Running-server-tests-inside-a-transaction).

## Considered Alternatives

* **Wrap test execution in a database transaction**: Use go-txdb to run tests
inside a transaction, and roll it back at the end of the test.

* **Truncate tables between tests**: Keep truncating the DB as we've been doing
all along.

## Decision Outcome

* Chosen Alternative: **Wrap test execution in a database transaction**

* The reason for this choice is to stop losing time. So far, with only a few
packages using transactions, we are saving 33 seconds when running
`make server_test_standalone` locally. When running individual packages or test,
the cumulative savings are even greater. If we stick with the status quo, every
36 seconds we could be saving but don't by not doing anything translates to
**3 person-weeks lost per year!**

## Pros and Cons of the Alternatives

### Wrap test execution in a database transaction

* `+` Allows multiple tests to interact with the database at once without
interfering with each other.
* `+` Saves us at least 3 person-weeks per year.
* `+` Allows us to discover and fix bugs due to false-positive tests
* `+` Does not break existing tests. Using transactions is an opt-in feature per
package, and we can speed up each package incrementally.
* `-` Requires updating existing tests, although in most cases, this is minimal
work.

### Keeping the status quo: truncating the DB in between tests

* `-` We lose at least 3 person-weeks per year waiting for tests to run.
* `-` Makes it easy to abuse DB setup in tests, resulting in false-positive
tests, and tests that cannot be run in isolation.
* `+` Does not require any additional work
