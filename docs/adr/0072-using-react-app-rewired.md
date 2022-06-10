# Using React-App-Rewired

**User Story:** *[MB-9033](https://dp3.atlassian.net/browse/MB-9033)* <!-- optional -->

## Background

The MilMove client application uses Create-React-App and React-Scripts for its
build toolchain. This build toolchain has many benefits to developing the client
application, but also has limitations when it comes to updating the
configuration of various build tools that are used in the development of the
client application. These tools such as webpack, ESlint, and Babel are
configured with pre-determined configurations that are in-accessible without
_ejecting_ from the Create-React-App toolchain.

Facebook's own documentation mentions that _If you aren't satisfied with the
build tool and configuration choices, you can `eject` at any time. This command
will remove the single build dependency from your project._ This is not an
option currently on MilMove as ejecting would add some maintainability overhead
that the development team may address at a later time or a later ADR.



*[context and problem statement]*
*[context and problem statement]*
*[decision drivers | forces]* <!-- optional -->

## Considered Alternatives

* *[alternative 1]*
* *[alternative 2]*
* *[alternative 3]*
* *[...]* <!-- numbers of alternatives can vary -->

## Decision Outcome

* Chosen Alternative: *[alternative 1]*
* *[justification. e.g., only alternative, which meets KO criterion decision driver | which resolves force force | ... | comes out best (see below)]*
* *[consequences. e.g., negative impact on quality attribute, follow-up decisions required, ...]* <!-- optional -->

## Pros and Cons of the Alternatives <!-- optional -->

### *[alternative 1]*

* `+` *[argument 1 pro]*
* `+` *[argument 2 pro]*
* `-` *[argument 1 con]*
* *[...]* <!-- numbers of pros and cons can vary -->

### *[alternative 2]*

* `+` *[argument 1 pro]*
* `+` *[argument 2 pro]*
* `-` *[argument 1 con]*
* *[...]* <!-- numbers of pros and cons can vary -->

### *[alternative 3]*

* `+` *[argument 1 pro]*
* `+` *[argument 2 pro]*
* `-` *[argument 1 con]*
* *[...]* <!-- numbers of pros and cons can vary -->
