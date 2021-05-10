# _Use OpenTelemetry to instrument code for distributed tracing_

**User Story:** *[Distributed Tracing ADR](https://dp3.atlassian.net/browse/MB-8053)*

_Give developers and operators an easier way to understand the behavior and structure of running systems by instrumenting code for distributed tracing._

## Considered Alternatives

* Use OpenTelemetry
* Use a vendor's instrumentation libraries
* Do not instrument

## Decision Outcome

* Chosen Alternative: *Use OpenTelemetry*
* OpenTelemetry is an emerging industry standard
  * vendors find benefit of being in the OpenTelemetry ecosystem because they
  no longer have to create or support instrumentation libraries in an every
  growing array of languages, i.e. as soon as language library exists for
  OpenTelemetry, the vendors automatically become available to support that
  given language.
* OpenTelemetry is vendor agnostic
  * tracing information can be sent to hosted services (e.g. Honeycomb.io, AWS
  X-Ray, etc) or self-hosted Open Source implementations (e.g. Zipkin, Jaeger,
  etc)
  * if left unconfigured, OpenTelemetry instrumentation calls default to
  lightweight/noop executions
* OpenTelemetry has well-maintained libraries for the languages used in the
layers of the MilMove project
  * i.e. Go (back-end); JavaScript (front-end); Python (load testing); etc
* Easily swappable back-ends
  * e.g. could choose a local Docker version of OpenZipkin for an all-local
  development environment
  * e.g. can use Honeycomb.io in the experimental commercial-cloud hosted
  environment
  * e.g. can swap in AWS X-Ray for use in GovCloud hosted environments
* Cons
  * as an abstraction layer, OpenTelemetry may prohibit usage of vendor-
  specific capabilities
  * some OpenTelemetry libraries and tools may trail their vendor-supported
  counterparts
  * instrumentation for tracing may be a vector for performance overhead

## Pros and Cons of the Alternatives

### ~~Use a vendor's instrumentation libraries~~

* `+` Enables accessing vendor-specific capabilities
* `-` Vendor lock-in in code
  * lock-in may be somewhat mitigated by translation layers available within
  the OpenTelemetry ecosystem, at the expense of increased configuration burden
  * example - choosing AWS X-Ray would work well in the deployed GovCloud
  environments, but it does not scale down to exclusively local development
  environments, i.e. X-Ray does not provide a UI for browsing distributed
  traces with their local X-Ray daemon
  * example - choosing Honeycomb.io's instrumentation libraries adds a lot of
  nice auto-instrumentation capabilities over OpenTelemetry, but since
  Honeycomb does not have FedRAMP (nor do most of their peers), the distributed
  tracing could not be enabled in the GovCloud deployed environments
  * example - using an open source tool (e.g. OpenZipkin) can scale down to
  local development, but would require more infrastructure support to self-
  host the data storage and UI tools in the GovCloud environments

### ~~Do not instrument~~

* `+` No work to be done
* `-` developers and operators continue to use current methods to build their
understanding of the system, which is likely slower and less complete than when
using distributed tracing

