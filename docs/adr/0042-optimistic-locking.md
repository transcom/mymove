# _Use E-tags for optimistic locking_

**User Story:** _[MB-651](https://dp3.atlassian.net/browse/MB-651?atlOrigin=eyJpIjoiODhkY2Y3ZTRjZGY3NDcxZjlmNTdmODZmNGUxMDZlN2UiLCJwIjoiaiJ9)_ <!-- optional -->

The system needs to be robust in the face of conflicting attempts to update/access the same data.

## Desired Outcomes

The data in the system is always coherent and current, even or especially when two users attempt to modify the same entity, e.g. a TOO and an agent for the Prime access the same Service Item.

Any user who is attempting to update stale data, i.e. a record that has already been updated by another user, is given and meaningful message telling them what has happened and a clear way to recover from the situation, e..g the current data is reloaded and they are able to continue with their intended changes.

## Considered Alternatives

- _[Use E-tags for optimistic locking]_
- _[Opt-in concurrency control]_
- _[do nothing or last request wins]_


## Decision Outcome

- Chosen Alternative: _[Use E-tags for optimistic locking]_
- _[justification. e.g., only alternative, which meets KO criterion decision driver | which resolves force force | ... | comes out best (see below)]_
- _[consequences. e.g., negative impact on quality attribute, follow-up decisions required, ...]_ <!-- optional -->

## Pros and Cons of the Alternatives <!-- optional -->

### _[Use E-tags for optimistic locking]_

We implement optimistic locking using the `If-Match` header. If the ETag header does not match the value of the resource on the server, the server rejects the change with a 412 Precondition Failed error. The client is therefore notified of the error, and can try the request again after updating their local copy of the resource.

- `+` _[This takes advantage of ETag HTTP header and follows standard/best conventions for REST APIs ]_
- `+` _[This guarantees that data is coherent and current]_
- `-` _[This will require more work on the client to handle rejection. In particular we will need to make sure the prime understands how this works]_
- `-` _[Updates will require a read to confirm that the client was working with the most recent copy]_


### _[opt-in concurrency control]_

The user can pass in an optional tag on update requests. If tag is provided, it is checked as above. If not, then the update happens.

- `+` _[it's a common solution]_
- `-` _[this is not a sound strategy if multiple parties can manipulate the same resource.]_

### _[do nothing or last request wins]_

- `+` _[Locking may not be needed if most resources are immutable, or if changes made by different parties are on different attributes]_
- `+` _[Since we will have a complete change record, we will be able to detect cases where one user has overridden the changes of another user]_
- `-` _[since the contractor will be developing an independent system there is a greater, uncontrollable risk that they will work with stale data]_

## References

[Optimistic Locking in a REST API](https://sookocheff.com/post/api/optimistic-locking-in-a-rest-api/)
