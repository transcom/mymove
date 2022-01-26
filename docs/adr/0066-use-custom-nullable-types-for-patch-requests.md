# Use custom nullable types for patch requests

**User Story:** [MB-10592](https://dp3.atlassian.net/browse/MB-10592)

There are several instances within the milmove app where a user will need to delete a field or reset a field to null on an object. Currently using patch requests this is not possible because there is no way to differentiate between a field that was omitted or a JSON field that is intentionally set to null. Go-swagger maps the existing YAML type pointers to Go primitive types, where nil could mean omitted or set to null.

## Considered Alternatives

* Use custom nullable types
* Use PUT requests in addition to PATCH requests

## Decision Outcome

* Chosen Alternative: Use custom nullable types

To solve this engineers will now be able to use a custom, nullable type that indicates if the field was set to explicitly set to null [delete the data] or implicitly [don't change anything]. This solution was selected because the implementation of creating the nullable types is relatively quick, and it allows MilMove to work with a single update request type: PATCH. Using the custom types will require engineers to alter how they access the value of the nullable fields, and understanding of how this type works. The maintenance of the custom code should be minimal after implementation.

## Pros and Cons of the Alternatives

### Use custom nullable types

* `+` Easy and quick to implement
* `+` Allows MilMove to use a single update request type
* `+` Allows API users to set fields to null without having to send all object data
* `-` Custom code will require milmove engineers to be aware of when and why to use nullable types. There may be a small learning curve.

### Use PUT requests in addition to PATCH requests

* `+` PUT requests are a widely understood and known method of making updates
* `+` Does not require custom code to implement
* `-` Requires the API user to return all object data when updating a single field
* `-` Will require adding additional support for another endpoint for a given resource (new handler, service, changes to YAML, FE calls/payloads, etc.)

## Additional Context

* [Slack thread discussing alternatives](https://ustcdp3.slack.com/archives/CP6PTUPQF/p1638833895016700)
* [Article describing the approach to creating the custom types](https://romanyx90.medium.com/handling-json-null-or-missing-values-with-go-swagger-4d7f37a2a7ca)
* [Github PR prototyping the custom, nullable string type](https://github.com/transcom/mymove/pull/7881)
