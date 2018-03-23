# Test Suites

**User Story:** [155076695](https://www.pivotaltracker.com/story/show/155076695)

Go's built in `testing` package does not support setup, teardown, etc. functionality that is provided by xUnit testing frameworks. In order to ensure that our test are running in a clean database, we want to call [`TruncateAll()`](https://godoc.org/github.com/gobuffalo/pop#Connection.TruncateAll) before each test executes (for more background, see ADR #0010).

Running code before or after *all* tests is supported by `go test` by defining a `TestMain` function, but this only allows for code to be run once at the very beginning and end of test execution and not before or after each test. Our requirement is for per-test setup and teardown, which is not supported.

## Considered Alternatives

* Use the [suite](https://godoc.org/github.com/stretchr/testify/suite) from [Testify](https://github.com/stretchr/testify)
* Use subtests in the style suggested [on the Go Blog](https://blog.golang.org/subtests#TOC_6.)
* Write our own suite-like functionality
* Require that tests perform all setup and teardown themselves

We did not consider adopting [another Go testing framework](https://awesome-go.com/#testing), as we wanted as thin of a layer on top of `testing` as possible.

## Decision Outcome

* Chosen Alternative: **Use the suite functionality from Testify**

* The API provided by `testify.suite`, which involves using a `struct` to represent the suite and declaring tests as methods on it fits nicely with our use case. It also provides a place for suite- or test-level state should we need it.

* Although a bit of extra code is needed so that it is picked up by `go test`, once a suite is setup the developer experience of *put a test in the right place and it will be detected and run* is preserved.

* It is worth noting that Testify is maintained but not as actively as we might like. This is a risk, but it would be possible to replace its functionality with our own code without too much difficulty. We will address this in the future once we have some experience with using this approach.

* Providing setup/teardown per-test functionality is familiar to developers who have worked with xUnit-style testing libraries in the past.

* Testify includes a lot of functionality that we don't wish to use at this time, such as its optional assertion packages `assert` and `require`.

## Pros and Cons of the Alternatives <!-- optional -->

### Use subtests in the style suggested [on the Go Blog](https://blog.golang.org/subtests#TOC_6.)

* `+` Does not rely on external dependencies.
* `+` Appeals to our desire to do things "the Go way".
* `-` Would require us to write some custom code to wrap each test as subtests still only execute the setup or teardown once.
* `-` Adds friction to writing a test as new tests may not be automatically detected and will need to be manually registered in order to run.
* `-` Test-driven tests (another possible way to structure tests) don't seem like a good fit for our integration-style tests that tend to be very long.

### Write our own suite-like functionality

* `+` Avoids adding an external dependency to the project.
* `+` Provides an opportunity to customize the suite API to suit the project's needs.
* `-` Requires writing non-trivial code for a problem we are just starting to understand.

### Require that tests perform all setup and teardown themselves

* `+` Makes all behavior associated with a test visible within the test body.
* `-` Obscures the intent of a test by requiring boilerplate to be places in each test.
* `-` Changing setup and teardown that affect multiple tests can require making changes to the body of tests, leading to maintenance issues in the long term.
