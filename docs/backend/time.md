# Time in Golang

## Table of Contents

<!-- toc -->

* [Clock Dependency](#clock-dependency)
* [MilMove Calendar Utils](#milmove-calendar-utils)

Regenerate with "bin/generate-md-toc.sh"

<!-- tocstop -->

## Clock Dependency

For example, let's say we have the following:

```go
package mypackage

import "time"

type MyObject struct {
    myTime time.Time
}

func MyTimeFunc() MyObject {
    return MyObject{myTime: time.Now()}
}
```

How do we test the contents of `MyObject`?
If we want to assert `myTime`
we need a way to know what `time.Now()` was when the function was called.

Instead of directly using the `time` package,
we can pass a clock as a dependency and call `.Now()` on that.
Then in our tests, we can assert against that clock!
Let's look at the example above with the `clock` package.

```go
package mypackage

import "fmt"
import "time"

import "github.com/facebookgo/clock"

type MyObject struct {
    myTime time.Time
}

func MyTimeFunc(clock clock.Clock) MyObject {
    return MyObject{myTime: clock.Now()}
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
    myObject := MyTimeFunc(testClock)
    if myObject.myTime != testClock.Now() {
        // both should equal epoch time!
        t.Errorf("time was not now, expected: %v, is: ", testClock.Now(), myObject.myTime)
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
