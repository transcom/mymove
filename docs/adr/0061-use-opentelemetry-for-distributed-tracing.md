# _Use OpenTelemetry to instrument code for distributed tracing_

**User Story:** *[MB-8053](https://dp3.atlassian.net/browse/MB-8053)*

_Give developers and operators an easier way to understand the behavior and structure of running systems by instrumenting code for distributed tracing._

*[decision drivers | forces]* <!-- optional -->

## Considered Alternatives

* Use OpenTelementry
* Use a vendor's instrumentation libraries
* Do not instrument

## Decision Outcome

* Chosen Alternative: *Use OpenTelemetry*
* OpenTelemetry is an emerging industry standard
  * vendors find benefit of being in the OpenTelemetry ecosystem because they no longer have to create or support instrumentation libraries in an every growing array of languages, i.e. as soon as language library exists for OpenTelemetry, the vendors automatically become available to support that given language.
* OpenTelemetry is vendor agnostic
  * tracing information can be sent to hosted services (e.g. Honeycomb.io, AWS X-Ray, etc) or self-hosted Open Source implementations (e.g. Zipkin, Jaeger. etc)
  * if left unconfigured, OpenTelemetry instrumentation calls default to lightweight/noop executions
* OpenTelemetry has well-maintained libraries for the languages used in the layers of the MilMove project
  * i.e. Go (back-end); JavaScript (front-end); Python (load testing); etc
* Easily swappable back-ends
  * e.g. could choose a local Docker version of OpenZipkin for an all-local development environment
  * e.g. can use Honeycomb.io in the experimental commercial-cloud hosted environment
  * e.g. can swap in AWS X-Ray for use in GovCloud hosted environments
* Cons
  * as an abstraction layer, OpenTelemetry may prohibit usage of vendor-specific capabilities
  * some OpenTelemetry libraries and tools may trail their vendor-supported counterparts
  * instrumentation for tracing may be a vector for performance overhead

## Pros and Cons of the Alternatives

### ~~Use a vendor's instrumentation libraries~~

* `+` Enables accessing vendor-specific capabilities
* `-` Vendor lock-in in code 
  * this may be mitigated by translation layers available within the OpenTelemetry ecosystem

### ~~Do not instrument~~

* `+` No work to be done
* `-` developers and operators continue to use current methods to build their understanding of the system

