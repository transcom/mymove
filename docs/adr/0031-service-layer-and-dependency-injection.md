# Service Layer and Dependency injection

## Context and Problem Statement

Currently the web service is built as two layers, Web Handlers (`pkg/handlers`) which implement interfaces based on the
swagger definitions of the services provided by the server and Model Objects (`pkg/models`) which marshal object representations
of data in and out of the database.

We are currently coming across a number of issues which suggest that we have reached the limits of what such a naive,
two-layer design can easily support, viz:

* it is not clear where Authorization code should live, i.e. code which enforces that logged in users only see and can access the data pertinent to them. Currently this is in the models (see ADR 0024)  but that means that models cannot be used for tools applications with different authorization controls, e.g. bulk loaders or admin interfaces.
* there is no place for code which touched multiple models and but is used by more than one handler, e.g. enforcing coherent state for multiple object relating to the same move (aka 'state machines') or making sure invoices line items are consistent between the GBL and the invoice.
* there is little or no encapsulation in the layers, so details of pop (database ORM) usage are in the handlers and equally swagger details appear in the model code. This makes testing and refactoring painful.

These variously lead to discussion around [Business Logic](https://en.wikipedia.org/wiki/Business_logic) and
[Service Layers](https://en.wikipedia.org/wiki/Service_layer_pattern). [Jim](https://github.com/jim) drew the teams attention
 to the [Service Object](https://medium.com/selleo/essential-rubyonrails-patterns-part-1-service-objects-1af9f9573ca1)
 pattern from rails. Looking for a similar pattern for go, it was suggested that we simply implement the Service Object
 pattern describe in the medium article in go.

This, in turn, lead to a search for a [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection) framework
for golang which could be used in place of the global state used in Rails.

This ADR explains the choice of DI framework [DIG](https://github.com/uber-go/dig) and details conventions for naming and
using objects in the new 3-layer design.

## Decision Drivers

* Maintained (new commits less than 6 months ago)
* Minimally intrusive on the code
* Little or no runtime impact on main code path for requests

## Considered Options

* Manually hooking up objects
* [Inject](https://github.com/facebookgo/inject) from Facebook
* [DIG](https://github.com/uber-go/dig) from Uber
* [Wire](https://github.com/google/go-cloud/tree/master/wire) from Google

## Decision Outcome

[DIG](https://github.com/uber-go/dig) is the chosen solution. It is actively maintained with a simple model which needs
no specific code/build processes changes for the objects being managed\

## Pros and Cons of the Options

### Manually hooking up objects

Currently the code does some dependency injection in the form of the HandlerContext objects which are used for the request handlers.

* Good - conceptually very simple to understand
* Good - has allowed us to move quickly
* Bad - doesn't scale well as number of objects managed increases. Currently context is turning into a dumping bag where each handler gets the context for every handler.
* Bad - top level context only - no idea of intermediate dependencies for, e.g. service layer

### [Inject](https://github.com/facebookgo/inject)

Originally from Parse Inject is a very simple injection framework which relies of struct tags & reflect to connect objects

* Good - small library which is recommended by folks who used it at Parse
* Bad - Orphaned. Last commit was 3 years ago and eng team who developed it are no longer at Facebook
* Bad - Adding tags to code to identify connections feels intrusive and fragile (what's the point of using a typed language if you don't use the compiler to check connections)

### [DIG](https://github.com/uber-go/dig)

From Uber (who also produce the logging framework Zap that we use)

* Good - Actively maintained. At time of writing last commit to master was 6 days ago
* Good - Minimally intrusive. Uses type system to connect object providers with object consumers
* Good - Simple. Complete documentation is a single readable [page](https://godoc.org/go.uber.org/dig).

### [Wire](https://github.com/google/go-cloud/tree/master/wire)

Build time tool to generate code to connect up dependencies

* Good - Build time tool minimizes the runtime overhead of connecting up object graph
* Bad - Build time tool would require build tool changes and maintenance
* Bad - Build time tool requires upfront investment in tooling before being able to evaluate and test
* Good - Seems consistent with Dig in terms of dependency resolution, so could probably be used to replace DIG at a later time if performance becomes an issue.
