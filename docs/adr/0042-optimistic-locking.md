# _Use E-tags for optimistic locking_

**User Story:** _[MB-651](https://dp3.atlassian.net/browse/MB-651?atlOrigin=eyJpIjoiODhkY2Y3ZTRjZGY3NDcxZjlmNTdmODZmNGUxMDZlN2UiLCJwIjoiaiJ9)_ <!-- optional -->

The system needs to be robust in the face of conflicting attempts to update/access the same data.

## Desired Outcomes

The data in the system is always coherent and current, even or especially when two users attempt to modify the same entity, e.g. a TOO and an agent for the Prime access the same Service Item.

Any user who is attempting to update stale data, i.e. a record that has already been updated by another user, is given and meaningful message telling them what has happened and a clear way to recover from the situation, e..g the current data is reloaded and they are able to continue with their intended changes.

## Considered Alternatives

- _[alternative 1]_
- _[alternative 2]_
- _[alternative 3]_
- _[...]_ <!-- numbers of alternatives can vary -->

## Decision Outcome

- Chosen Alternative: _[alternative 1]_
- _[justification. e.g., only alternative, which meets KO criterion decision driver | which resolves force force | ... | comes out best (see below)]_
- _[consequences. e.g., negative impact on quality attribute, follow-up decisions required, ...]_ <!-- optional -->

## Pros and Cons of the Alternatives <!-- optional -->

### _[alternative 1]_

- `+` _[argument 1 pro]_
- `+` _[argument 2 pro]_
- `-` _[argument 1 con]_
- _[...]_ <!-- numbers of pros and cons can vary -->

### _[alternative 2]_

- `+` _[argument 1 pro]_
- `+` _[argument 2 pro]_
- `-` _[argument 1 con]_
- _[...]_ <!-- numbers of pros and cons can vary -->

### _[alternative 3]_

- `+` _[argument 1 pro]_
- `+` _[argument 2 pro]_
- `-` _[argument 1 con]_
- _[...]_ <!-- numbers of pros and cons can vary -->
