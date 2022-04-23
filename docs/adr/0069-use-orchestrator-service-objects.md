# Use orchestrator service objects

**NOTE:** [ADR-0033 Service Object Layer](./0033-service-object-layer.md) is the one that started us on using service
objects. This ADR doesn't supersede it, but it may be helpful to read the other one if you want more background on
why we use service objects.

## Problem statement

### Summary

There are actions that require several things to happen across related models, e.g. creating an `MTOShipment` may also
need `MTOAgents` and/or `MTOServiceItems` to be created. The question is then, should we have a single service object
that knows how to do all the things, or break things down and orchestrate service objects somehow?

There are also actions, like routing a shipment, that have their own service objects, but then how should those be
handled when it comes to the existing `MTOShipment` service objects?

### Story Time

When we started working on implementing the logic to handle PPMs, we decided to start a pattern of having shipment types
have their own model that will be a child of the `MTOShipment` model (see
[ADR-0067 Add a child table to mto_shipments for PPMs](./0067-ppm-db-design.md) for more details).

This meant we needed to add logic for creating and updating `PPMShipments`, but we weren't sure where to add it because
we don't have a set standard for this. The ADR that started us using service objects (see note at top) is more focused
on the fact that we should use service objects, but doesn't really get into when a new service object should be created
vs expanding an existing one, or how to handle the times that you need multiple service objects in order to perform a
full action. If we look at the existing code base, we have implemented a variety of different patterns, which makes it
harder to maintain, and more confusing for folks adding new logic.

For `PPMShipments`, we opted to create separate service objects from the existing `MTOShipment` service objects, but due
to still being in the middle of discussing this ADR, we ended up with two different implementations for how they
interact with the `MTOShipment` service objects. The `PPMShipment` creator internally calls the `MTOShipment` creator,
while the `PPMShipment` updater expects the `MTOShipment` updater to be called separately (currently being done in the
handler). The [Considered Alternatives](#considered-alternatives) section will use these as examples.

## Measuring success

Initial success would be:

* Have our first orchestration service objects set up for managing shipment-related service objects.
  * At first this would primarily be focused on coordinating existing `MTOShipment` and `PPMShipment` service objects.
* Team members know about this ADR and have a good example to work from (the first orchestration object mentioned above).
* Have documentation around these service objects to help people when they need to work on them.

Long-term, we would ideally:

* Expand documentation to have helpful information to help people decide when to use orchestration service objects vs
  having service objects take other service objects as dependencies or some other method.
* Evaluate how the existing `MTOShipment` service objects handle `MTOAgents`, `MTOServiceItems`, and shipment routing
  and see if they could/should move to using the orchestration service objects.
* Evaluate other service objects to see if the can benefit from having a higher level orchestration service object.

### Observability

_How will this change be observed by other team members?_

This ADR will be announced in meetings where most of the MilMove team is present. If approved, as orchestration service
objects are added or updated, they could be announced at internal demos and/or BE check in meetings.

### Implementation Plan

1. Create new orchestration service objects for managing existing `MTOShipment` and `PPMShipment` service objects.
2. Update existing usages of the `MTOShipment` and `PPMShipment` service objects to use the orchestration service
  objects instead.
3. Update documentation around service objects to include info about this new kind of service object, with examples as
  appropriate.
4. Announce new service objects at demo and BE check in.

### Ownership

AppEng would own this since it's related to the way we work with the business logic of our application code.

As the need for orchestration service objects arises, it would be up to the person writing code and reviewers to ensure
orchestration service objects are used, added, or update as needed.

## Considered Alternatives

<!-- The order of the alternatives is used in the pros/cons section below, so if they change here, update there accordingly. -->

For these alternatives, we'll keep talking about `MTOShipment` and `PPMShipment` service objects to give concrete
examples, but the idea for this ADR is for the decision to be applicable in to any situation where we have multiple
service objects that are needed for an action.

1. Use separate service objects and the handler decides which to call.
    1. In this one, the `PPMShipment` service objects will internally call the `MTOShipment` service objects to do
       whatever needs to be done to the parent MTO object.
    1. This is how `PPMShipmentCreator` works now.
1. Use separate service objects and the handler orchestrates the calls.
    1. So for example, we need an `MTOShipment` to exist before we create a PPMShipment, so the handler would call one,
       then the other.
    1. This is how `PPMShipmentUpdater` works now.
1. Update existing service objects to contain the new logic.
    1. An example of this is how the `MTOShipment` service objects handle `MTOAgents` and `MTOServiceItems` internally
       rather than calling other service objects.
1. Have separate service objects, but the handler just calls the existing service objects. Then we update the existing
   service objects to call the new service objects as needed.
    1. This is how routing works currently. The `MTOShipment` service objects will call the shipment router service
       object as needed.
1. Have composable service objects and create orchestrator service objects that then call the appropriate service
   objects.
    1. Not sure if we have examples of this.
1. Leave things as they are.

## Decision Outcome

Chosen Alternative -> Option 5: Have composable service objects and create orchestrator service objects that then call
the appropriate service objects.

While it would take some work up-front, and might add work as new service objects are created, this seems to be the
cleanest path for what we need. It keeps business logic out of our handlers, while at the same time still keeping our
service objects from getting too large, since each can focus on just what _it_ needs to do. This also leaves us with a
codebase that's easier to test since each service object can be isolated as needed in tests.

## Pros and Cons of the Alternatives

### Option 1: Use separate service objects and the handler decides which to call

* `+` Keeps some service objects from getting too large since they only have to worry about their own thing,
  calling other service objects as needed.
* `+` This makes service objects work as standalone objects. E.g. the `MTOShipment` service objects can do their thing
  without worrying about other things, while the `PPMShipment` service objects would know to call the `MTOShipment`
  service objects as needed so they wouldn't need something else to coordinate the calls.
* `+` It is easier to mock out service objects if they are separate objects than if they are all lumped into a single one.
* `-` This puts business logic and database details/transactions in the handlers. Ideally handlers should only focus on
  translating data from the protocol layer to the service layer and vice-versa.

### Option 2: Use separate service objects and the handler orchestrates the calls

* `+` Keeps service objects from getting too large since each will encapsulate the logic they need for their own work.
* `+` It is easier to mock out service objects if they are separate objects than if they are all lumped into a single one.
* `-` This puts business logic and database details/transactions in the handlers. Ideally handlers should only focus on
  translating data from the protocol layer to the service layer and vice-versa.
* `-` Testing at the handler level takes more setup than testing at a service object level so testing the orchestration
  of service objects would take more work this way.

### Option 3: Update existing service objects to contain the new logic

* `+` Leaves business logic and database details/transactions out of the handlers. Ideally handlers should only focus on
  translating data from the protocol layer to the service layer and vice-versa.
* `-` This makes our service objects become incredibly large and hard to maintain. As it is, we already have some
  service objects, like the `MTOShipment` ones, that contain a large amount of logic and can be hard to parse through.
* `-` Related to the point above, but finding the logic for a specific type of change, e.g. tracking how PPMs change,
  would be harder to do if it's all in one service object.
* `-` Testing large service objects is hard because you need to account for many branches of code.

### Option 4: Have separate service objects, but the handler just calls the existing service objects. Then we update the existing service objects to call the new service objects as needed

* `+` Leaves business logic and database details/transactions out of the handlers. Ideally handlers should only focus on
  translating data from the protocol layer to the service layer and vice-versa.
* `+` This makes service objects work as standalone objects. E.g. the `PPMShipment` service objects can do their thing
  without worrying about other things, while the `MTOShipment` service objects would know to call the `PPMShipment`
  service objects as needed so they wouldn't need something else to coordinate the calls.
* `+` Might keep some service objects from getting too large since they only have to worry about their own thing,
  calling other service objects as needed.
* `+` It is easier to mock out service objects if they are separate objects than if they are all lumped into a single one.
* `-` Related to the point above, there is the potential for service objects to still get large though as more
  connections arise. E.g. if every shipment type had its own service objects, the `MTOShipment` service objects would
  need to have their own base logic, plus the logic for calling each of the related service objects correctly, which
  could expand to be too large.

### Option 5: Have composable service objects and create orchestrator service objects that then call the appropriate service objects

* `+` Leaves business logic and database details/transactions out of the handlers. Ideally handlers should only focus on
  translating data from the protocol layer to the service layer and vice-versa.
* `+` Keeps service objects from getting too large since each can focus on doing their own thing.
* `+` The orchestration service objects could serve as a nice way of viewing all the steps needed for an action at a
  high level.
* `+` It is easier to mock out service objects if they are separate objects than if they are all lumped into a single one.
* `-` Requires work to create the orchestration service objects that we need now.
* `-` Potentially adds more work when creating new service objects if you need to create both the main service objects
  you're focusing on, plus orchestration service objects if needed and they don't already exist.

### Option 6: Leave things as they are

* `+` No extra work is needed right now.
* `-` We are left with code that implements most of the options and is inconsistent even within the same service objects.
* `-` This leaves some cases of business logic and database details/transactions in the handlers. Ideally handlers
  should only focus on translating data from the protocol layer to the service layer and vice-versa.
* `-` We are left without a standard that would help guide future folks.

## Resources

* [Slack thread where we discussed options for handling related service objects](https://ustcdp3.slack.com/archives/C022R8A9FC4/p1645221666532619)
* [Back-end check-in notes where we discussed the options](https://dp3.atlassian.net/wiki/spaces/MT/pages/1661665292/2022-02-24+Meeting+notes)
* [Existing documentation on service objects](https://transcom.github.io/mymove-docs/docs/backend/guides/service-objects)
