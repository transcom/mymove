# Use swagger-client to make calls to API from client

**User Story:** _[#153793371]_ <!-- optional -->

_[context and problem statement]_
_[decision drivers | forces]_ <!-- optional -->

## Considered Alternatives

* [swagger-client](https://www.npmjs.com/package/swagger-client)
* [Swagger Codegen](https://swagger.io/swagger-codegen/)
* roll our own (fetch)

## Decision Outcome

* Chosen Alternative: _swagger client_
* Allows use to dog-food our Swagger config with minimal setup and rapidly develop client code that uses API

## Pros and Cons of the Alternatives <!-- optional -->

### _swagger-client_

* `+` dynamically generated from `swagger.yaml`
* `+` no install or make tasks required
* `-` documentation is spotty
* `-` error handling benefit is unclear until we have more API to work against

### _Swagger Codegen_

* `+` generates a module ready for publishing on NPM (which might be useful for TSPs in future)
* `+` have static code you can see
* `-` usage seemed cumbersome and verbose
* `-` have to install java or use docker container to set up

### _roll your own_

* `+` complete control over
* `-` API calls can deviate from swagger configuration
* `-` manual maintenance
