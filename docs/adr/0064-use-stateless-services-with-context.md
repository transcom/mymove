# Use stateless services with context

## Problem statement

We want our services to be composable, so that one service can call
another. We also want to be able to have a per request trace id
associated with a logger so that we can correlate log messages in a
single request.

Right now, most services are initialized with a database connection pool
and a logger.

When using the `database/sql` package, each database
request uses a different connection to the database. Some of our
services start transactions and then want to call other services. The
"sub-service" uses its own connection pool and thus does not have
visibility to the changes made within the transaction.

Some services have a way to set the connection used the service, but
since we have a single service object, it seems almost certain that
with goroutines handling each request that we'll have errors
when multiple requests are running simultaneously.

The same problem exists for logging because services are logging using
the "global" logger and not one initialized per request. That means
service logs don't include the per request trace id.

## Considered Alternatives

- Modify all service methods to accept *both* a `Context` and a custom interface
- Modify all service methods to accept a `Context`
- Modify all service methods to accept a custom interface
- Do nothing

## Decision Outcome

- Chosen Alternative: Modify all service methods to accept *both* a `Context` and a custom interface

## Pros and Cons of the Alternatives

### Modify all service methods to accept *both* a `Context` and a custom interface

- `+` Services can be composed, calling each other and passing
  transaction connections around if necessary
- `+` Modifying all service methods means implementers and consumers
  don't have to think about what the method signature should be. It
  also ensures the necessary info is available in case the
  implementation changes and they are now needed
- `+` Services that log now can include the per request trace id
- `+` For instrumentation, the [opentelemetry
  api](https://opentelemetry.io/docs/go/getting-started/) needs a
  context
- `+` A `Context` does not provide type safety, so passing both
  provides a *lot* more compile time checking.
- `-` Passing two arguments instead of one

The need for the context for instrumentation means we really have to
pass the context around. Passing two arguments to get more type safety
is a good trade off.

### Modify all service methods to accept a `Context`

- `+` Services can be composed, calling each other and passing
  transaction connections around if necessary
- `+` A custom interface provides ease of use and type safety, easy
  extensibility in the future, and mocking if necessary
- `+` Modifying all service methods means implementers and consumers don't have to
  think about what the method signature should be. It also ensures all are available in case the implementation changes and
  they are now needed
- `+` Services that log now can include the per request trace id
- `+` For instrumentation, the [opentelemetry
  api](https://opentelemetry.io/docs/go/getting-started/) needs a
  context
- `+` A single argument is passed
- `-` A `Context` does not provide type safety, so there's no way to
  know at compile time if a `Context` has the required parameters
  (connection, logger). Because this is a massive change (~8000 lines
  changed), not having compile time safety significantly increases the
  risk.

The loss of type safety is not worth the "cost" of having to provide
 two arguments instead of one.

### Modify all service methods to accept a custom interface

- `+` Services can be composed, calling each other and passing
  transaction connections around if necessary
- `+` A custom interface provides ease of use and type safety, easy
  extensibility in the future, and mocking if necessary
- `+` Modifying all service methods means implementers and consumers don't have to
  think about what the method signature should be. It also ensures all are available in case the implementation changes and
  they are now needed
- `+` Services that log now can include the per request trace id
- `+` A single argument is passed
- `-` For instrumentation, the [opentelemetry
  api](https://opentelemetry.io/docs/go/getting-started/) needs a
  context. Go best practice is pretty clear that a `Context` should
  not be stored in a struct, so including it in the custom interface
  doesn't seem to be an option.

The need for the context for instrumentation is a deal breaker for
this option.

### Do nothing

- `+` Makes things no worse
- `-` Our current approach will almost certainly result in errors when
  multiple requests are handled simultaneously
- `-` Our current approach doesn't allow our services to call other
  services from inside transactions
- `-` Our current approach doesn't allow for instrumentation or
  logging per request trace ids.

This option really isn't realistic.
