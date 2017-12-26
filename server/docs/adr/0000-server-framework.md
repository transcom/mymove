# Use Truss' [golang](https://golang.org/) web server skeleton to build API for dp3

The Personal Property Prototype project needs an API and services in place to support

1. The React [client](../../client/README.md) application
1. Integration from Transport Service Providers (TSPs) who want to service the moves

Because of 2. above the API will need to be fully documented, secured and readily accessible from a variety of client applications

## Considered Alternatives

* Using [gRpc](https://grpc.io/) as a way of publishing the APIs
* Using [OpenAPI](https://www.openapis.org/) within the framework of a [Buffalo](https://gobuffalo.io/) golang application
* Using [OpenAPI](https://www.openapis.org/) within the framework of a custom golang application
* Using [Aws Lambda](https://aws.amazon.com/lambda/) to host a 'serverless' API
* Using another language and web framework, e.g. Python/Django, Ruby/Rails or Phoenix/Elixir

## Decision Outcome

* Chosen Alternative: Using OpenAPI within the context of a custom golang application
* Golang is a fairly straightforward choice for the implementation language bringing together Type safety, Active support and development, and enough community experience to be relatively low risk
* gRpc is not (yet) well suited to Web applications (see below) so was ruled out by our need to support the React client
* Buffalo brought too many 'opinions'/baggage in terms of web pipeline, lack of support for React out of the box
* Lambda is very tempting for simplifying deployment and management, but has too many unknowns for the team in terms of cost and performance.

* Consequence: need to rapidly evaluate OpenAPI code generation tools and test/confirm the belief that they can be intergrated into our application framework

## Pros and Cons of the Alternatives

### [gRpc](https://grpc.io/)

* `+` High performance RPC mechanism with active support from google
* `-` gRpc doesn't currently have good [support for web clients](https://improbable.io/games/blog/grpc-web-moving-past-restjson-towards-type-safe-web-apis) and while Improbable is driving this and has a solution, there is no official implementation yet.

### [OpenAPI](https://www.openapis.org/) within the framework of a [Buffalo](https://gobuffalo.io)

* `+` Buffalo is the most complete golang web framework we have found and pulls together solutions to most of the commong web framework concerns (authentication, authorization, middleware injection, database management and mapping)
* `+` It has an active community supporting it
* `-` The built in web pipeline does not (out of the box) have support for a React pipeline, and in particular the Create React App framework.
* `-` Adapting the framework to accomodate out client work would undermine many of the benefits of using an off the shelf web framework without mitigating any of the dependecy risks
* `-` It is not clear how OpenAPI code generation tools might co-exist with the buffalo framework

### [OpenAPI](https://www.openapis.org/) within the framework of a custom golang application

* `+` Truss has already done most of the preparatory work for this in our Prototype Web Application
* `+` We have control of a known codebase
* `-` It is not clear how best to intergrate OpenAPI code generation tools into our current workflow.
* `-` We need to manage the evaluation, selection and integration of all the parts of the system, e.g. [authentication](https://github.com/markbates/goth), [database management](https://github.com/markbates/pop) etc

### [Aws Lambda](https://aws.amazon.com/lambda/)

* `+` Is GovCloud approved and so requires no extra work if we need to operate in that environment
* `+` Removes the need for server infrastructure.
* `-` Pricing can be a suprise as it is based on traffic (according to anecdata from people Truss has asked)
* `-` We have no experience building this way so it carries risks of longer dev ramp up and potential performance/reliability suprises.

### Other languages/web frameworks

* `+` Some team members are more familiar with, e.g. Django or Rails
* `+` Some frameworks have more longjevity in the field
* `-` Django and Rails are less compelling in a single-page responsive client app environment
* `-` Type safety has repeatedly shown
* `-` React/go is what was proposed, with good reason, during the bid process for this contract. Adopting another approach without a compelling reason seems contrary.
