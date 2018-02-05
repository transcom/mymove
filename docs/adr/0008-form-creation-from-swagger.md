# Generate forms from swagger definitions of payload

**User Story:** #154407746

For the first phase development we wanted to create very basic forms to enter data for a model that matches the printed forms 1299, 1797, and 2278. The first of these had many fields and so we decided that it would be nice if we could generate the view based on the swagger model definition and potentially another source of data for ui specific information (such as grouping). (Manually generating this form would be error-prone and tedious.) There are a number of libraries out there that support doing this with JsonSchema (of which Swagger is a subset), but it might be better for us to roll our own so that we can support `redux-forms` and `uswds`.

## Considered Alternatives

* In house form creation
* [Uniforms](https://github.com/vazco/uniforms)
* [react-jsonschema-form](https://github.com/mozilla-services/react-jsonschema-form)

## Decision Outcome

* In House
* The other libraries were not easy to integrate with `redux-forms` and expected other styling libraries (such as bootstrap).
* This will require more maintenance on our part and we may want to reconsider should our needs become more advanced.

## Pros and Cons of the Alternatives <!-- optional -->

### _[alternative 1]_

* `+` complete control over look and feel
* `+` current implementation is very simple
* `-` we may need more advanced features for future work

### Uniforms

* `+` support for many schema formats
* `+` under active development
* `+` has rich language
* `-` built for meteor and may have features we don't need

### react-jsonschema-form

* `+` nice addition of separate schema for ui specific concerns
* `+` under active development
* `-` use of bootstrap did not play well with uswds in exploration
