# Back-end Programming Guide

## Table of Contents

<!-- toc -->

* [Go](#go)
  * [Acronyms](#acronyms)
  * [Style and Conventions](#style-and-conventions)
  * [Importing Dependencies](#importing-dependencies)
  * [Querying the Database Safely](#querying-the-database-safely)
  * [`models.Fetch*` functions](#modelsfetch-functions)
  * [Logging](#logging)
    * [Logging Levels](#logging-levels)
  * [Errors](#errors)
    * [Don't bury your errors in underscores](#dont-bury-your-errors-in-underscores)
    * [Log at the top level; create and pass along errors below](#log-at-the-top-level-create-and-pass-along-errors-below)
    * [Use `errors.Wrap()` when using external libraries](#use-errorswrap-when-using-external-libraries)
    * [Don't `fmt` errors; log instead](#dont-fmt-errors-log-instead)
    * [If some of your errors are predictable, pattern match on them to provide more error detail](#if-some-of-your-errors-are-predictable-pattern-match-on-them-to-provide-more-error-detail)
  * [Libraries](#libraries)
    * [Pop](#pop)
  * [Learning](#learning)
  * [Testing](#testing)
    * [General](#general)
    * [Coverage](#coverage)
    * [Models](#models)
  * [Time](#time)
  * [Miscellaneous Tips](#miscellaneous-tips)
* [Environment settings](#environment-settings)
  * [Adding `ulimit`](#adding-ulimit)

Regenerate with "bin/generate-md-toc.sh"

<!-- tocstop -->

## Go

### Acronyms

Domain concepts should be used without abbreviation when used alone.

Do:

* `TransportationServiceProvider`
* `TrafficDistributionList`

Avoid:

* `TSP`
* `TDL`

However, when used as a specifier or part of another name, names that have existing acronyms should use the acronym for brevity.

Do:

* `TSPPerformance`

Avoid:

* `TransportationServiceProviderPerformance`

Acronyms should always be either all caps or all lower-cased.

Do:

* `TSPPerformance`
* `tspPerformance`

Avoid:

* `tSPPerformance`
* `tspperformance`
* `TspPerformance`

### Style and Conventions

Generally speaking, we will follow the recommendations laid out in [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments). By its own admission, this page:

> _...collects common comments made during reviews of Go code, so that a single detailed explanation can be referred to by shorthands. This is a laundry list of common mistakes, not a style guide._

Despite not being an official style guide, it covers a good amount of scope in a concise format, and should be able to keep our project code fairly consistent.

Beyond what is described above, the following contain additional insights into how to write better Go code.

* [What's in a name?](https://talks.golang.org/2014/names.slide#1) (how to name things in Go)
* [Go best practices, six years in](https://peter.bourgon.org/go-best-practices-2016/)
* [A theory of modern Go](https://peter.bourgon.org/blog/2017/06/09/theory-of-modern-go.html)

### Importing Dependencies

Dependencies are managed by [dep](https://github.com/golang/dep). New dependencies are automatically detected in import statements. To add a new dependency to the project:

1. Add the package to the import statement of a Go file.
1. `make clean`
1. `make server_generate`
1. `dep check` (to verify what's missing.) If it looks reasonable then...
1. `dep ensure`

### Querying the Database Safely

* SQL statements _must_ use PostgreSQL-native parameter replacement format (e.g. `$1`, `$2`, etc.) and _never_ interpolate values into SQL fragments in any other way.
* SQL statements must only be defined in the `models` package.

Here is an example of a safe query for a single `Shipment`:

```golang
// db is a *pop.Connection

id := "0186ad95-14ed-4c9b-9f62-d5bd124f62a1"

query := db.Where("id = $1", id)

shipment := &models.Shipment{}
if err = query.First(shipment); err == nil {
  pp.Println(shipment)
}
```

### `models.Fetch*` functions

Functions that return model structs should return pointers.

Do:

```go
func FetchShipment(db *pop.Connection, id uuid.UUID) (*Shipment, error) {}
```

Avoid:

```go
func FetchShipment(db *pop.Connection, id uuid.UUID) (Shipment, error) {}
```

This is for a few reasons:

* Many Pop methods that accept models need to use pointers in order to modify the struct being passed in. Using them everywhere
  models are used makes the code more consistent and avoids needing to convert to and from pointers.
* Any methods on struct model instances need to have pointer receivers in order to mutate the struct. This is a common point of confusion
  and an easy way to introduce bugs into a codebase.

### Logging

We use the [Zap](https://github.com/uber-go/zap) logging framework from Uber to produce structured log records.
To this end, code should avoid using the [SugaredLogger](https://godoc.org/go.uber.org/zap#Logger.Sugar)s without
a very explicit reason which should be record in an inline comment where the SugaredLogger is constructed.

#### Logging Levels

Another reason to use the Zap logging package is that it provides more nuanced logging levels than the basic Go logging package.
That said, leveled logging is only meaningful if there is a common pattern in the usage of each logging level. To that end,
the following indicates when each level of Logging should be used.

* **Fatal** This should never be used outside of the server startup code in main.go. Fatal log messages call `sys.exit(1)` which unceremoniously kills the server without running any clean up code. This is almost certainly never what you want in production.

* **Error** Reserved for system failures, e.g. cannot reach a database, a DB insert which was expected to work failed,
  or an `"enum"` has an unexpected value. In production an Error logging message should alert the team that
  all is not well with the server, so avoid being the 'Boy Who Cried Wolf'. In particular, if there is an API which takes an object ID as part of the URL,
  then passing a bad value in should NOT log an Error message. It should log Info and then return 404.

* **Warn** Don't use Warn - it rarely, if ever, adds any meaningful signal to the logs.

* **Info** Use for recording 'Normal' events at a granularity that may be helpful to tracing and debugging requests,
  e.g. 404's from requests with bad IDs, authentication events (user logs in/out), authorization failures etc.

* **Debug** Debug events are of questionable value and should be used during development, but probably best removed
  before landing changes. The issue with them is, if they are left in the code, they quickly become so dense in the logs
  as to obscure other debug log entries. This leads to people an arms race of folks adding 'XXXXXXX' to comments
  in order to identify their log items. If you must use them, I suggest adding an, e.g. zap.String("owner", "nick")

### Errors

Some general guidelines for errors:

#### Don't bury your errors in underscores

If a function or other action generates an error, assign it to a variable and either return it as part of your function's output or handle it in place (`if err != nil`, etc.). There will be the very occasional exception to this - one is within tests, depending on the test's goal. If you find yourself typing that underscore, take a moment to ask yourself why you're choosing that option. On those very rare occasions when it is the correct behavior, please add a comment explaining why.

_Don't:_
`myVal, _ := functionThatShouldReturnAnInt()`

_Do:_

```golang
    myVal, err := functionThatShouldReturnAnInt()
    if err != nil {
         return myVal, errors.Wrap(err, "function didn't return an int")
   }
```

#### Log at the top level; create and pass along errors below

If you're creating a query (1) that is called by a function (2) that is in turned called by another function (3), create and return errors at levels 1 and 2 (and possibly handle them immediately after creation, if needed), and log them at level 3. Logs should be created at the top level and contain context about what created them. This is more difficult if logs are being created in every function and file that supports the operation you're working on. Here's an example of when to create errors and when to handle them:

In `pkg/models/blackout_dates.go`, an error is created and returned:

```golang
func FetchTSPBlackoutDates(tx *pop.Connection, tspID uuid.UUID, shipment Shipment) ([]BlackoutDate, error) {
  ...
  err = query.All(&blackoutDates)
  if err != nil {
    return blackoutDates, errors.Wrap(err, "Blackout dates query failed")
  }

  return blackoutDates, err
}
```

In `pkg/awardqueue/awardqueue.go`, `FetchTSPBlackoutDates` is called, and any possible error is handled. This function also returns an error.

```golang
func ShipmentWithinBlackoutDates(tspID uuid.UUID, shipment models.Shipment) (bool, error) {
  blackoutDates, err := models.FetchTSPBlackoutDates(aq.db, tspID, shipment)

  if err != nil {
    return false, errors.Wrap(err, "Error retrieving blackout dates from database")
  }

  return len(blackoutDates) != 0, nil
}
```

Finally, at the top level in `attemptShipmentOffer` in the same file, any errors bubbled up from `ShipmentWithinBlackoutDates` or `FetchTSPBlackoutDates` are handled definitively, halting the progress of the longer function if the underlying processes and queries didn't complete as expected in the functions being called:

```golang
func (aq *AwardQueue) attemptShipmentOffer(shipment models.Shipment) (*models.ShipmentOffer, error) {
  aq.logger.Info("Attempting to offer shipment", zap.Any("shipment_id", shipment.ID))
  ...
  isAdministrativeShipment, err := aq.ShipmentWithinBlackoutDates(tsp.ID, shipment)
  if err != nil {
    aq.logger.Error("Failed to determine if shipment is within TSP blackout dates", zap.Error(err))
    return err
  }
  ...
}
```

The error is created and passed along at the lowest level, logged and passed along at the middle level (along with other errors that can happen within that function), and logged again at the highest level before finally halting the progress of the process if an error is present.

#### Use `errors.Wrap()` when using external libraries

[`errors.Wrap()`](https://godoc.org/github.com/pkg/errors) provides greater error context and a stack trace, making it especially useful when dealing with the opacity that sometimes comes with external libraries. `errors.Wrap()` takes two parameters: the error and a string to provide context and explanation. Keep the string brief and clear, assuming that the fuller cause will be provided by the context `errors.Wrap()` brings. It can also add useful context for errors related to internal code if there might otherwise be unhelpful opacity. `errors.Errorf()` and `errors.Wrapf()` also capture stack traces with the additional function of string substitution/formatting for output. Instead of just returning the error, offer greater context with something like this:

```golang
if err != nil {
        return errors.Wrap(err, "Pop validate failed")
}
```

#### Don't `fmt` errors; log instead

`fmt` can provide useful error handling during initial debugging, but we strongly suggest logging instead, from when you write the initial lines of a new function. Using logging creates structured logs instead of the unstructured, human-friendly-only output that `fmt` does. If an `fmt` statement offers usefulness beyond your initial troubleshooting while working, switch it to `errors.Wrap()` or `logger.Error()`, perhaps with [Zap](https://github.com/uber-go/zap).

_Don't:_
`fmt.Println("Blackout dates fetch failed: ", err)`

_Do:_
`logger.Error("Blackout dates fetch failed: ", err)`

#### If some of your errors are predictable, pattern match on them to provide more error detail

Some errors are predictable, such as those from the database that Pop returns to us. This gives you the option to use those predictable errors to give yourself and fellow maintainers of code more detail than you might get otherwise, like so:

```golang
// FetchServiceMemberForUser returns a service member only if it is allowed for the given user to access that service member.
 func FetchServiceMemberForUser(ctx context.Context, db *pop.Connection, user User, id uuid.UUID) (ServiceMember, error) {
  var serviceMember ServiceMember
  err := db.Eager().Find(&serviceMember, id)
  if err != nil {
    if errors.Cause(err).Error() == recordNotFoundErrorString {
      return ServiceMember{}, ErrFetchNotFound
    }
    // Otherwise, it's an unexpected err so we return that.
    return ServiceMember{}, err
  }

  if serviceMember.UserID != user.ID {
    return ServiceMember{}, ErrFetchForbidden
  }

  return serviceMember, nil
 }
```

You can also use `errors.Wrap()` in this situation to provide access to even more information, beyond the breadcrumbs left here.

### Libraries

#### Pop

We use Pop as the ORM(-ish) to mediate access to the database. [The Unofficial Pop Book](https://andrew-sledge.gitbooks.io/the-unofficial-pop-book/content/) is a valuable companion to Pop’s [GitHub documentation](https://github.com/gobuffalo/pop).

### Learning

If you are new to Go, you should work your way through all of these resources (in this order, ideally):

1. [A Tour of Go](https://tour.golang.org) (in-browser interactive language tutorial)
1. [How to Write Go Code](https://golang.org/doc/code.html) (info about the Go environment, testing, etc.)
1. [Effective Go](https://golang.org/doc/effective_go.html) (how to do things “the Go way”)
1. [Daily Dep documentation](https://golang.github.io/dep/docs/daily-dep.html) (common tasks you’ll encounter with our dependency manager)
1. [Exercism](http://exercism.io/languages/go/about) offers a series of exercises with gradually increasing complexity

Additional resources:

* [GoDoc](https://godoc.org/) (where you can read the docs for nearly any Go package)
* Check out the [Go wiki](https://github.com/golang/go/wiki/Learn)
* Advanced Testing with Go [Video](https://www.youtube.com/watch?v=yszygk1cpEc) and [Article](https://about.sourcegraph.com/go/advanced-testing-in-go) (great overview of useful techniques, useful for all Go programmers)
* _Book_: [The Go Programming Language](http://www.gopl.io/)
* _Article_: [Copying data from S3 to EBS 30x faster using Golang](https://medium.com/@venks.sa/copying-data-from-s3-to-ebs-30x-faster-using-go-e2cdb1093284)

### Testing

Knowing what deserves a test and what doesn’t can be tricky, especially early on when a project’s testing conventions haven’t been established. Use the following guidelines to determine if and how some code should be tested.

#### General

* Use table-driven tests where appropriate.
* Make judicious use of helper functions so that the intent of a test is not lost in a sea of error checking and boilerplate. Use [`t.Helper()`](https://golang.org/pkg/testing/#T.Helper) in your test helper functions to keep stack traces clean.

#### Coverage

* Always test exported functions.
  Exported functions should be treated as an API layer for other packages.
  Cover the expected behavior and error scenarios as a user of that API.
* Try not to test unexported functions.
  Unexported functions are implementation details of exported ones
  and should not change the intended usage.
  If you find that an unexported function is complex and needs testing,
  it might mean it needs to be refactored as it's exported function elsewhere.

#### Models

In general, focus on testing non-trivial behavior.

* Structs do not need to be tested as they have no behavior of their own.
* Struct methods warrant a unit test if they contain important behavior, e.g. validations.
* Avoid testing functionality of libraries, e.g. model saving and loading (which is provided by Pop)
* Try to leverage the type system to ensure that components are “hooked up correctly” instead of writing integration tests.

### Time

Some helpful tips on dealing with time
in the MilMove Go codebase
can be found in [this doc](backend/time.md)

### Miscellaneous Tips

* Use `golang` instead of `go` in Google searches.
* Try to use the standard lib as much as possible, especially when learning.

## Environment settings

### Adding `ulimit`

Dep appears to open many files simultaneously, particularly as the project matures and depends on more and more third-party repositories. You may encounter a message like this as a `dep-status` hook error when trying to commit locally: `remote repository does not exist, or is inaccessible: : pipe: too many open files`.

To fix this, run `ulimit -n 5000` in your terminal. This increases the number of file handles (the details on files that a process holds open) on your system. You can run `ulimit -n` to see how many are currently allowed; you may see a number like 128 or 256. (Run `ulimit -a` to see all current limits on your system, including pipe size, stack size, and user processes.) If running this in your terminal window allows you to complete your commit, you may wish to add it to `/.bash_profile` or whichever system file you use for your terminal settings.
