# Model Authorization and Handler Design

Users should only have access to records that they are authorized to see. Its important that our system for implementing that constraint is robust and easy to use. In fact, it's important that the system be difficult to misuse and to encourage handlers that correctly protect our data.

## Considered Alternatives

We conducted a white-boarding session to demonstrate three different options:

1. All validation happens in the models. Handlers call model functions to fetch and create models that are scoped to the current user.

2. Validation happens in per-base-model middleware (e.g. everything that is hung off a move would share one piece of middleware) passing the validated model into the handler and rejecting unauthorized requests.

3. A generalized middleware that uses reflection to reject unauthorized requests before the handler gets run.

## Decision Outcome

* Handlers get models through fetchers that accept an authenticated user, fetchers handle authorization
* Fetchers return generic errors, handlers have helper function to convert generic errors into responses
* Creator functions are hung from parent models, if available, e.g. `serviceMember.CreateBackupContact()`

## Pros and Cons

### Validation happens in the models

* `+` No additional middleware is needed, we are on rails with go-swagger
* `+` These patterns can be used no matter where in the stack models are needed. (e.g. the rate engine, or any other entry point besides handlers)
* `-` A handler could ignore those patterns and access data incorrectly

### Validation happens in per-base-model middleware

* `+` Handlers are only run when access is allowed, it's much more difficult for a handler to access data it is not supposed to.
* `+` Handlers are slimmer and their requirements are encoded in their arguments
* `-` It's difficult to know how to implement this with go-swagger
* `-` data accessors other than handlers have to re-implement this.

### Validation happens in a generalized middleware

* `+` We could write this once and never have to update it
* `+` Putting the fetched models into the context is a common go pattern
* `-` It might be complex to implement, and reflection is always worrisome
* `-` data accessors other than handlers have to re-implement this.

something something
