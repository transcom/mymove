# *Service Object Layer*

Currently the web service is built as two layers, Web Handlers ([pkg/handlers](https://github.com/transcom/mymove/tree/master/pkg/handlers)) which implement interfaces based on the swagger definitions of the services provided by the server and Model Objects ([pkg/models]((https://github.com/transcom/mymove/tree/master/pkg/models))) which marshal object representations of data in and out of the database.

We are currently coming across a number of issues which suggest that we have reached the limits of what such a naive, two-layer design can easily support, viz:

It is not clear where Authorization code should live, i.e. code which enforces that logged in users only see and can access the data pertinent to them. Currently this is in the models (see [ADR 0024](https://github.com/transcom/mymove/blob/master/docs/adr/0024-model-authorization-and-handler-design.md)) but that means that models cannot be used for tools applications with different authorization controls, e.g. bulk loaders or admin interfaces.
Furthermore, there is no place for code which touches multiple models, but is used by more than one handler, e.g. enforcing coherent state for multiple object relating to the same move (aka 'state machines') or making sure invoices line items are consistent between the GBL and the invoice.
There is little or no encapsulation in the layers, so details of pop (database ORM) usage are in the handlers and equally swagger details appear in the model code. These examples and others show how painful our experiences with testing and refactoring could be.

Ultimately, this lead to a discussion around Business Logic and Service Layers. Jim drew the teams attention to the Service Object pattern from rails. Looking for a similar pattern for go, it was suggested that we simply implement the Service Object pattern. For further context on the pattern that inspires our approach, please see this [article](https://medium.com/selleo/essential-rubyonrails-patterns-part-1-service-objects-1af9f9573ca1).

This, in turn, lead to a search for a Dependency Injection framework for golang which could be used in place of the global state used in Rails. Nick T. then went on to complete a spike that investigated Dig, a dependency injection framework, resulting in a very high value, but risky [PR](https://github.com/transcom/mymove/pull/1118) with over 132 different file changes. The [service object layer and dependency injection design document](https://docs.google.com/document/d/1xlqgVSTf9JUhZfWR18rvzaGPg2iHcF7uKRahGvrO45E/edit#) ultimately provided a plan for integration - first provide and integrate examples of service objects, provide training on using service objects, then finally adding the dependency injection framework later. This ADR is primarily concerned with the decisions behind adding a service object layer.

## Decision Drivers

* Ease of adoption
* Minimal impact
* Provide encapsulation of logic
* Better ability to test and refactor encapsulations of logic
* Code re-usability

## Considered Options

* Service Objects
* Do nothing

## Decision Outcome

Adopt service object layer, an architectural pattern for writing code that allows for encapsulation of logic, code reusability, ultimately keeping our handler code much less complex and more lightweight.

Resources:

* [Essential Ruby On Rails Patterns - Part 1: Service Objects](https://medium.com/selleo/essential-rubyonrails-patterns-part-1-service-objects-1af9f9573ca1)
* [Using the Service Object Pattern in Go](https://www.calhoun.io/using-the-service-object-pattern-in-go/)
* [Service Object Layer & Dependency Injection Design Document](https://docs.google.com/document/d/1xlqgVSTf9JUhZfWR18rvzaGPg2iHcF7uKRahGvrO45E/edit#)

## Pros and Cons of the Alternatives

### *Service Objects*

* `+` Allows better organization of business logic
* `+` Keeps API handler endpoints less complex by writing less code
* `+` Improves maintainability as ease of refactoring is increased
* `+` Allows better unit testing
* `+` Allows encapsulation of logic
* `+` Provides code re-usability
* `+` Provides a pattern for writing better code as more conventions are introduced defining codification of services
* `+` Provides easier scalability
* `-` New learning as team must adopt new pattern
* `-` Dependency management can become difficult as services become more complex

### *Do Nothing*

* `-` Maintains everything as-is and we reap none of the above benefits.
