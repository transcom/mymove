# How to Add Application Logging

Application logging in MyMove is built on top of the [Zap](https://godoc.org/go.uber.org/zap) logger.  However, each package defines its own logging interface to self-document and enable us to extend Zap method definitions with custom implementations.  To learn how to extend Zap, see the `Embedded Structs` section below.

## New Logging

To add basic logging to a package, you create a `logger.go` file in the package that includes an interface that defines the logging methods required.

```go
package mypackage
...
import (
  "go.uber.org/zap"
)
...
// Logger is an interface that describes the logging requirements of this package.
type Logger interface {
  Error(msg string, fields ...zap.Field)
}
...
```

You can then pass the logger object through an interface to underlying function calls.

```go

func (h *MyService) Call(ctx context.Context, r *http.Request, logger Logger) {
  ...
  logger.Error("hit an error!")
  ...
}
```

## Expand Logging

If your logging requirements expand for a given package, you can simply add the functions to the package's interface.  If the require methods are not implemented yet, you have the flexibility to implement the method at the project level or create an embedded struct at the package level.

```go
type Logger interface {
  Info(msg string, fields ...zap.Field)
  Error(msg string, fields ...zap.Field)
  WithOptions(opts ...zap.Option) *zap.Logger
}
```

## Embedded Structs

Since we use an [interface](https://gobyexample.com/interfaces) instead of a single [struct type](https://tour.golang.org/moretypes/2), we can enhance loggers at the project or package level while maintaining interface compatibility.  The logger `CustomLogger` as defined below logs the number running go routines on `Fatal` messages.

```go

import (
  "go.uber.org/zap"
)
...
type CustomLogger struct {
  *zap.Logger
}

func (e *CustomLogger) Fatal(msg string, fields ...zap.Field) {
  fields = append(fields, zap.Int("goroutines", runtime.NumGoroutine()))
  e.Logger.Fatal(msg, fields...)
}
...
```
