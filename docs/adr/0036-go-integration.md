# Use Separate Integration Package for Go Integration Tests

**User Story:**

MilMove Engineers have started to decouple parts of the Go codebase.
Namely, we now have a core setup of
`publicapi` (similarly `internalapi` and `adminapi`), `services`, and `models`.
This helps us decouple implementations in each of those layers
and better unit test specific functionality.

One of the reasons for adding `services`
was to speed up handler tests by mocking service implementations
and removing the constant setup and tear down of the database.
However, by mocking services,
we lose the integration style testing that our previous handler tests supplied.
These are the tests that ran against the handler,
executing all the code as if it were production with no mocking.
Because of the unit testing style changes with services,
old integration tests are being run in tandem with new mocked handler unit tests
(causing performance issues)
or not being written for new code
(causing code quality/safety issues).

This ADR looks to address the gaps created by mock driven unit tests
and create a strategy for organizing integration tests on critical code paths.
It does **not** answer the question of what we should be integration testing.

Note that we are also not referring to Cypress end-to-end tests
which are sometimes referred to as integration tests.

## Considered Alternatives

- **Run integration in same suite as handlers and flag with `testing.Short`**

An example of the suite setup and test:

```go
// SetupTest sets up the test suite by preparing the DB
func (suite *HandlerSuite) SetupTest() {
  if !testing.Short() {

	suite.DB().TruncateAll()
  }
}

func (suite *HandlerSuite) TestCreateShipmentHandlerEmpty() {
  if testing.Short() {
      suite.T().Skip("Skip database integration test")
  }
  // Rest of test
}
```

Tests would be run with `go test -short`
if developers need quick unit testing in development environments.
Or synonymously with a `make server_test -short` target.

- **Use separate suites for integration tests, but within handlers package**

In this example,
integration tests would use `IntegrationHandlerSuite`
and the existing `HandlerSuite` would no longer have database tests requiring setup/teardown.

- **Move integration tests to separate package and flag the suite with `testing.Short`**

In this example,
we move all tests requiring external dependencies (database, APIs) to an `integration` package.
This package would have its own suite.

The suite test would be prefaced as follows so that it can be excluded with `-short`.

```go
func TestHandlerSuite(t *testing.T) {
  if testing.Short() {
      suite.T().Skip("Skip database integration test")
  }
  // rest of suite instantiation
}
```

## Decision Outcome

- Chosen Alternative: **Move integration tests to separate package and flag the suite with `testing.Short`**

Moving all dependency wiring to a separate package allows unit tests to stay decoupled from dependencies.
Namely, `handlers` is only required to pull in packages that mimics an interface dependency based production structure
that our service layer pushes to achieve.
Integration tests then become more synonymous with wiring up our webserver with dependencies.
`Short()` tests are also easier to reason about at the package/suite,
rather than test level.

## Pros and Cons of the Alternatives

### Run integration in same suite as handlers and flag with `testing.Short`

- `+` Keeps current structure of tests (no migration).
- `+` Integration tests live close to their unit counterparts,
  increasing developer awareness.
- `-` Dependency packages are imported into handlers package for testing.
  Note that this currently happens in `api.go`,
  but individual handlers are unaware of dependency implementations.
- `-` Suite setup/teardown happen on the suite level,
  but short flagging would be required within individual tests and suite setup.
  In other words, short logic would be sparse across tests.

### Use separate suites for integration tests, but within handlers package

- `+` Integration tests live close to their unit counterparts,
  increasing developer awareness.
- `+` Short flagging happens on suite/package level.
- `-` Dependency packages are imported into handlers package for testing.
  Note that this currently happens in `api.go`,
  but individual handlers are unaware of dependency implementations.
- `-` Multiple suites in a single package can be confusing.

### Move integration tests to separate package and flag the suite with `testing.Short`

- `+` Integration tests with dependencies are clearly marked and easily tested as separate from unit tests.
  Dependencies are no longer imported for unit tests.
- `+` Short flagging happens on suite/package level.
  Tests can also be run per directory.
- `-` Requires developer awareness.
  This might be moot,
  as our current instructions for testing with services ignore integration tests.
