# Use If-Match / E-tags for optimistic locking

The system needs to be robust in the face of conflicting attempts to update/access the same data.

## Desired Outcomes

The data in the system is always coherent and current, even or especially when two users attempt to modify the same entity, e.g. a TOO and an agent for the Prime access the same Service Item.

Any user who is attempting to update stale data, i.e. a record that has already been updated by another user, is given and meaningful message telling them what has happened and a clear way to recover from the situation, e.g. the current data is reloaded and they are able to continue with their intended changes.

## Considered Alternatives

- Use Last-Modified / If-Unmodified-Since for optimistic locking
- Use E-tags for optimistic locking
- do nothing or last request wins

## Decision Outcome

### Chosen Alternative: Use E-tags for optimistic locking

Using E-tags provides us with a simple way to prevent users from updating stale data while avoiding the risk of different programming languages parsing timestamps differently. For example, Ruby's default `DateTime.parse` method doesn't includes microseconds by default. Go's `time.Time()` conversion function does. To avoid this, we're doing the following:

- Each record has an `updated_at` timestamp, which is effectively a version number.
- Before sending a record's payload to the client, we base64 encode that timestamp so that the client doesn't know they're dealing with a timestamp - they only have to worry about storing a string.
- The client sends a `PUT` or `PATCH` along with the record's base64 encoded E-tag in the `If-Match` header of the request.
- We query for the given record's `updated_at` value in the database, base64 encode it and compare it against the `If-Match` header.
- If they don't match, we return a `412 Precondition Failed`, indicating to the client that they need to re-fetch the record that they're attempting to update.
- If they do match, we update the record and return the updated payload **including the new E-tag as a part of the payload body**

## Pros and Cons of the Alternatives

### Use Last-Modified / If-Unmodified-Since for optimistic locking

Every response (for single resources) contains a `Last-Modified header` with a HTTP date. (If aggregate results are sent, the last modified is an attribute included with each item.) When requesting an update using a PUT or PATCH request, the client has to provide the Last-Modified value of the resource via the header If-Unmodified-Since. The server rejects the request (via HTTP/1.1 412 Precondition failed), if the last modified date of the entity is after the given date in the header.

- `+` This follows standard/best conventions for REST APIs.
- `+` This guarantees that data is coherent and current.
- `+` We already store the last modified timestamp.
- `-` This will require more work on the client to handle rejection. In particular we will need to make sure the prime understands how this works.
- `-` Updates will require a read to confirm that the client was working with the most recent copy.
- `-` Different programming languages have different default behavior when dealing with timestamps. This presents a risk since we don't know what language the Prime will build their API client in. This puts even more responsibility on clients of the API to tread carefully.

### Use E-tags for optimistic locking

We implement optimistic locking using the `If-Match` header. If the ETag header does not match the value of the resource on the server, the server rejects the change with a 412 Precondition Failed error. The client is therefore notified of the error, and can try the request again after updating their local copy of the resource.

- `+` This takes advantage of ETag HTTP header and follows standard/best conventions for REST APIs.
- `+` This guarantees that data is coherent and current.
- `+` Identifiers are opaque, allowing us to avoid the issue of different date parsing rules in other programming environments.
- `-` This will require more work on the client to handle rejection. In particular we will need to make sure the prime understands how this works.
- `-` Updates will require a read to confirm that the client was working with the most recent copy.


### do nothing or last request wins

- `+` Locking may not be needed if most resources are immutable, or if changes made by different parties are on different attributes
- `+` Since we will have a complete change record, we will be able to detect cases where one user has overridden the changes of another user
- `-` since the contractor will be developing an independent system there is a greater, uncontrollable risk that they will work with stale data

## References

- [Optimistic Locking in a REST API](https://sookocheff.com/post/api/optimistic-locking-in-a-rest-api/)
- [RESTful API and Event Scheme Guidelines](https://opensource.zalando.com/restful-api-guidelines/index.html#optimistic-locking)
