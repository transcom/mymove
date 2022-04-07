# Use Soft Delete Instead of Hard Delete

Due to our contractual obligations with the federal government, we must be able to access deleted data even several years after it’s been used in the system.

## Considered Alternatives

* Leave everything as is. Allow the system to continue hard deleting records.
* Introduce soft deletion into the system.

## Decision Outcome

* Chosen Alternative: **Introduce soft deletion into the system.**

This option allows the system to continue treating records as 'deleted' while maintaining the records. This will allow us to serve 'deleted' records when audited or other obligations demand so. Soft delete functionality has not yet been implemented throughout the entire codebase but it is expected to be the sole deletion method moving forward.

Please note that soft delete is to be treated like a hard delete in the regard that the process should never be reversed or that data can be 'un-deleted'.

## Pros and Cons of the Alternatives

### Leave everything as is. Allow the system to continue hard deleting records

* `+` No changes needed to be done.
* `-` Risk legal exposure.
* `-` Fail to comply with contractual obligations with government.

### Introduce soft deletion into the system

* `+` Complies with contractual obligations to the government.
* `+` In possession of records when asked for.
* `-` Implementing soft delete will be a long, involved process.
* `-` Database will have to deal with holding both active and 'deleted' records.

## Reference

* [Documentation](https://transcom.github.io/mymove-docs/docs/backend/guides/how-to/soft-delete/) on how to implement soft deletion