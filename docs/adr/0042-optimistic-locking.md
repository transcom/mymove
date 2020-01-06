# Use Last-Modified / If-Unmodified-Since for optimistic locking

**User Story:** [MB-651](https://dp3.atlassian.net/browse/MB-651?atlOrigin=eyJpIjoiODhkY2Y3ZTRjZGY3NDcxZjlmNTdmODZmNGUxMDZlN2UiLCJwIjoiaiJ9)

The system needs to be robust in the face of conflicting attempts to update/access the same data.

## Desired Outcomes

The data in the system is always coherent and current, even or especially when two users attempt to modify the same entity, e.g. a TOO and an agent for the Prime access the same Service Item.

Any user who is attempting to update stale data, i.e. a record that has already been updated by another user, is given and meaningful message telling them what has happened and a clear way to recover from the situation, e..g the current data is reloaded and they are able to continue with their intended changes.

## Considered Alternatives

- Use Last-Modified / If-Unmodified-Since for optimistic locking
- Use E-tags for optimistic locking
- do nothing or last request wins


## Decision Outcome

### Chosen Alternative: Use Last-Modified / If-Unmodified-Since for optimistic locking

Using last modified timestamps will meet our desired outcomes and is easier to implement than the other suitable alternative.

## Pros and Cons of the Alternatives

### Use Last-Modified / If-Unmodified-Since for optimistic locking

Every response (for single resources) contains a `Last-Modified header` with a HTTP date. (If aggregate results are sent, the last modified is an attribute included with each item.) When requesting an update using a PUT or PATCH request, the client has to provide the Last-Modified value of the resource via the header If-Unmodified-Since. The server rejects the request (via HTTP/1.1 412 Precondition failed), if the last modified date of the entity is after the given date in the header.

- `+` This follows standard/best conventions for REST APIs.
- `+` This guarantees that data is coherent and current.
- `+` We already store the last modified timestamp.
- `-` This will require more work on the client to handle rejection. In particular we will need to make sure the prime understands how this works.
- `-` Updates will require a read to confirm that the client was working with the most recent copy.

### Use E-tags for optimistic locking

We implement optimistic locking using the `If-Match` header. If the ETag header does not match the value of the resource on the server, the server rejects the change with a 412 Precondition Failed error. The client is therefore notified of the error, and can try the request again after updating their local copy of the resource.

- `+` This takes advantage of ETag HTTP header and follows standard/best conventions for REST APIs.
- `+` This guarantees that data is coherent and current.
- `-` This will require more work on the client to handle rejection. In particular we will need to make sure the prime understands how this works.
- `-` Updates will require a read to confirm that the client was working with the most recent copy.
- `-` We would need to do work to calculate the hash of resources and to store that value.


### do nothing or last request wins

- `+` Locking may not be needed if most resources are immutable, or if changes made by different parties are on different attributes
- `+` Since we will have a complete change record, we will be able to detect cases where one user has overridden the changes of another user
- `-` since the contractor will be developing an independent system there is a greater, uncontrollable risk that they will work with stale data

## References

[Optimistic Locking in a REST API](https://sookocheff.com/post/api/optimistic-locking-in-a-rest-api/)
[RESTful API and Event Scheme Guidelines](https://opensource.zalando.com/restful-api-guidelines/index.html#optimistic-locking)
