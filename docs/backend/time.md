# Time in Golang

## Table of Contents

<!-- toc -->

* [Clock Dependency](#clock-dependency)
  * [Setting the Mock Clock](#setting-the-mock-clock)
* [MilMove Calendar Utils](#milmove-calendar-utils)

Regenerate with "bin/generate-md-toc.sh"

<!-- tocstop -->

## Clock Dependency

`time.Now()` can cause a lot of side effects in a codebase.
One example is
that you can't test the "current" time
that happened in a function you called in the past

For example, let's say we have the following:

```go
package mypackage

import "time"

func MyTimeFunc() time.Time {
    return time.Now()
}

func TestMyTimeFunc(t *testing.T) {
    if MyTimeFunc() != time.Now() {
        // This will error!
        // The time in the function and the test happen at different times
        t.Errorf("time was not now")
  }
}
```

How do we test the contents of the return here?
If we want to assert the time
we need a way to know what `time.Now()` was when the function was called.

Instead of directly using the `time` package,
we can pass a clock as a dependency and call `.Now()` on that.
Then in our tests, we can assert against that clock!
The clock can be anything as long as it adheres to the `clock.Clock` interface
as defined in the
[facebookgo clock package](https://godoc.org/github.com/facebookgo/clock#Clock).
We could, for example,
make the clock always return the year 0,
or the 2019 New Year,
or maybe your birthday!
In this clock package,
there are two clocks.

* The real clock where `clock.Now()` will call `time.Now()`.
* A mock clock where `clock.Now()` always returns epoch time.
  We'll show later how to change that!

Let's look at the example above with the `clock` package.

```go
package mypackage

import "fmt"
import "time"

import "github.com/facebookgo/clock"

func MyTimeFunc(clock clock.Clock) time.Time {
    return clock.Now()
}

// Then our caller
func main() {
    // clock.New() creates a clock that uses the time package
    // it will output current time when .Now() is called
    fmt.Print(MyTimeFunc(clock.New()))
}
```

Then in our tests we can use a mock clock that freezes `.Now()` at epoch time:

```go
func TestMyTimeFunc(t *testing.T) {
    testClock := clock.NewMock()
    if MyTimeFunc(testClock) != testClock.Now() {
        // both should equal epoch time, we won't hit this error
        t.Errorf("time was not now")
  }
}
```

Cool, but what if I want to use a different date?
Say my test relies on our `TestYear` constant.
The [clock.Mock clock](https://godoc.org/github.com/facebookgo/clock#Mock)
allows us to add durations to the clock and set the current time.
Note that the `clock.Clock` interface does not allow this,
it needs to happen before passing the mock clock through the interface parameter.

### Setting the Mock Clock

Here's an example using the test above and setting the time to September 30 of TestYear:

```go
func TestMyTimeFunc(t *testing.T) {
    testClock := clock.NewMock()
    dateToTest := time.Date(TestYear, time.September, 30, 0, 0, 0, 0, time.UTC)
    timeDiff := dateToTest.Sub(c.Now())
    testClock.Add(timeDiff)
    if MyTimeFunc(testClock) != testClock.Now() {
        // both will now be September 30 of TestYear
        // we'll pass the test again
        t.Errorf("time was not now")
  }
}
```

## MilMove Calendar Utils

The MilMove project has a set of date/calendar util
to help develop and test.
You can find them in the [dates package](../../pkg/dates)

For testing, we also have `TestYear`
in the [constants package](../../pkg/testdatagen/constants.go)
which should be used instead of the current year.
