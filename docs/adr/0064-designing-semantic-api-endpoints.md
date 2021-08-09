# Designing semantic API endpoints for logical semantic actions

**User Story:** *[MB-8929][jira]* <!-- optional -->

[jira]: https://dp3.atlassian.net/browse/MB-8929 "MB-8929 Jira Ticket"

## Background

This ADR is intended to work through the design of the GHC Prime API (Prime)
when it comes to actions represented in our endpoints. Currently, the Prime
contains endpoints that require multiple API calls to perform a single action
such as _Divert Shipment_. A _Divert Shipment_ action mentioned in this ADR is
considered a logical semantic action. This means that it is an action that
requires multiple steps in order to be performed. These types of actions are
easily understood by a user of the Prime API, but require low-level coordination
such as adding or removing records


To divert a shipment today, it involves a number of
steps using different API endpoints.

1. Update an existing shipment with `{ "diversion": true }` to
   `updateMTOShipment`.
1. Updating the shipment's _destination address_ to be the diversion point
   address using `updateMTOShipmentAddress`.
1. Creating a new shipment with `createMTOShipment` that has all the same
   information as the previous shipment but with `{ "diversion": true }` and
   uses the diversion point address as the _pickup address_.

This multi-pronged approach to updating a shipment record leaves a lot of
responsibility to the user of the Prime. This also allows for the handler
methods of the Prime to do more work when it comes to modifying data.

### Decision drivers and forces

The decision for this change is to have the Prime endpoints focus on the changes
based on specific actions that the Prime will need to handle. To borrow the
example from above, this ADR pushes for a new endpoint named `divertMTOShipment`
which would do all the actions above but by only retrieving the shipment once
from the database and creating and/or updating the appropriate shipments tied to
the divergence.

1. Update a shipment with a diversion by sending an original shipment ID and
   diversion point to a `divertMTOShipment`. The Prime method handler would then
   perform the necessary steps outlined above. The users of the Prime would not
   have to keep track of associated shipments.

This approach would allow the users of the Prime to focus on semantic actions on
resources in a one-to-one manner rather than having the user of the Prime to
have to manage updating resources and associating themselves.

To remove a diversion of a shipment, an endpoint named `cancelMTOShipment` would
allow the Prime user to remove a shipment. The Prime method handler would check
if `{ "diversion": true }` is part of the shipment and then check for associated
shipments by matching the _pickup_ or the _destination_ address. Depending on
which address is matched, the existing shipments would be updated with the
appropriate address for the shipment. The Prime API would be required to ensure
that shipments being cancelled no longer have a `{ "diversion": true }` property
and ensure that the appropriate _destination_ or _pickup_ address is being
updated.

### Additional actions on resources in the future

This ADR is not an exhaustive list of possible actions that fall under this
category. Actions may be added in the future that are not outlined explicitly.
In short, this ADR is proposing obfuscating low-level associations that are not
easily conveyed with a single verb on a given resource.

There are some additional actions that may be performed using the Prime API
around the state of a shipment or move. These actions are `cancelMTOShipment`,
`updateShuttleWeights`, and possibly others. Any additional actions needed in
the future would follow the design concepts of actions outlined in this ADR.

### What this isn't

This ADR is not intended to update existing endpoints but rather add new ones
that perform logical semantic actions that would normally require multiple calls
to the Prime API by the user.

## Considered Alternatives

- *Do nothing to the Prime and document these edge cases of logical semantic
  actions on a resource*.
- ✅ *Create new semantic methods for single actions on resources*

## Decision Outcome

### Chosen Alternative: *Create new semantic methods for single actions on resources*

This solution will add some complexity to the Prime handler methods since the
Prime is currently written to handle actions on a resource rather than inferring
the resources necessary to perform a given action. In other words, we have
handler methods in the Prime which expect to read from the database for each
request. For these new semantic actions, the Prime method handlers will need to
resolve all of the associated records itself and have documentation that
explains to the user the potential and real side effects for perform these
semantic actions.

The semantic method handlers for the Prime will also require a failsafe
mechanism in case actions cannot be performed within a given action/transaction.
This would happen automatically for the Prime user and would not require an
additional API requests to rollback/undo multi-pronged actions.

#### Pros and cons

- ➕ Less documentation to outline necessary steps to perform logical semantic
    actions
- ➕ The ability for Prime engineers to reuse methods for Prime handlers that
    may work on disparate resources.
- ➕ Better logical relationship with a given action and resource from the
    perspective of the Prime user (e.g. a given action performed on a given
    resource without exposing details to the Prime user).
- ➕ Ability for the Prime engineers to perform extra actions without having to
    modify actions if the underlying details of an action change over time.
- ➖ Method handlers for the Prime would need to handle specific edge cases for
    the user in case of a failure during a logical semantic action (e.g.
    ensuring a new shipment ID is not created if the method fails to update a
    shipment ID with a `diversion` property)
- ➖ Logical semantic actions create methods which may have side effects that
    are not apparent to the user due to obfuscation of the underlying resources
    that are being modified.
- ➖ Logical semantic actions create scenarios that give the Prime method
    handlers a lot of different edge cases that would otherwise be better
    handled by specific steps when modifying shipments.

## Pros and cons of the alternatives

### _Do nothing and document_

- ➕ Less engineering work is needed from Roci engineers.
- ➕ The Prime API remains flexible and does not assume what the Prime will
    want.
- ➕ RESTful APIs are easier to maintain than RPC-style APIs.
- ➕ Leaving the API as is allows the API to follow standard MilMove patterns.
    This reduces the amount of onboarding and confusion for new team members.
- ➖ More work is necessary from Roci content designers.
- ➖ More documentation is necessary to outline the necessary steps to perform
    logical semantic steps from both Roci engineers and designers.
- ➖ More work for users of the Prime API to perform logical semantic
    actions that may leverage disparate resources with multiple API calls.
