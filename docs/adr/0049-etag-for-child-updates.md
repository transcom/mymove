# Do not update child records using parent's E-tag

**User Story:** [Jira Story](https://dp3.atlassian.net/browse/MB-2566)

When we have an endpoint that updates a record in the db, it's sometimes desirable to update a child record as well.

Generally, to update a record, the caller must provide an E-tag, passed in the header `If-Match` that matches that of the record in the db.

However the parent and child have two different E-tags, and only the parent's E-tag is passed in a sole parameter in the header.

Therefore, it's not possible to pass in the child and parent E-tag cleanly.

## Considered Alternatives

* Make a new endpoint for child updates so they can be updated separately with the correct E-tag.
* Pass a second E-tag in the body if the child is to be updated.
* Bubble up a child's updated_at value to the parent, so that the child and parent will have one E-tag.

## Decision Outcome

We will make a new endpoint for child updates so they can be updated separately with the correct E-tag.

Currently this is just true for address and agent updates.

## Pros and Cons of the Alternatives

### Make a new endpoint for child updates

* `+` The mechanism for optimistic locking stays the same across all endpoints, so it's understandable for the Prime.
* `+` The `updated_at` value for parent and child record will correctly state the last time that record was updated.
* `-` More endpoints to create and maintain.

### Pass a second E-tag in the body if the child is to be updated

* `-` The mechanism differs when you want to update a child, as the E-tag is passed in body instead of header. Makes the mechanism inconsistent, harder to reason about and harder to explain to Prime.
* `+` Fewer endpoints to maintain

### Bubble up a child's updated_at value to the parent

* `-` Adds complexity because the child may have multiple parents and Prime would not realize that they have unwittingly updated unrelated records.
* `-` The mechanism differs from the norm and making exceptions for certain updates, will make it hard to be consistent across the codebase.
* `+` Fewer endpoints to maintain

