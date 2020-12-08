# Nesting Swagger paths in the Prime API with multiple IDs

The Prime API manages Move Task Orders and the child objects for these orders, such as shipments and payment requests.
The relationship between these objects is a straight-forward Many-to-One (many shipments or payment requests
to one MTO), and the IDs for the child objects can be used alone to uniquely identify a record. The endpoints for this
API follow a standard pattern of basing updates and reads off of one record at a time (sometimes returning multiple
records when children of the base object are included). When writing the paths for these endpoints, two distinct
strategies have been used:

1. Nest the path for the new endpoint by including the ID of the parent Move Task Order for the object, e.g.
`/move-task-orders/{moveTaskOrderID}/mto-shipments/{mtoShipmentID}`
2. Do not nest under the parent MTO and only include the ID of the object being accessed in the path, e.g.
`/payment-requests/{paymentRequestID}`

The first strategy is problematic because the `moveTaskOrderID` in the path is not validated, so there is no meaningful
connection between that value and the ID of the child object. Only the ID of the child object - the `mtoShipmentID` in
this example - is needed to fetch, validate, and update the record. The UUIDs also make the paths extremely long.

The second strategy is problematic because it is inconsistent with the other the endpoints. It also obfuscates the
relationships between objects.

Ultimately, the question is whether or not we should be including IDs that are functionally unnecessary in the path for
the sake of a clear hierarchical structure in the API. This ADR sets this question within the context of the Prime API,
but aims to propose a generic solution that is acceptable for most other APIs as well.

## Considered Alternatives

* Leave the codebase as-is
* Nest all paths with the IDs of the parent objects
* Do not include multiple IDs unless functionally necessary. Start a new root for an object that is uniquely
identifiable by one ID.

## Decision Outcome

### Chosen Alternative: *Do not include multiple IDs unless functionally necessary*

* **Justification:** Philosophically, we should not require input that is unnecessary for our processing. In the Prime
API, we can grab all the information we need for security and general validation using the UUID of the child object.
Requiring extra IDs in the paths for the endpoints therefore becomes an aesthetic decision, and provides little benefit
to the usability of the API. The clarity it provides to object relationships could (and perhaps should) be subsumed into
a well-documented ERD.

Functionally, using a new root for accessing each object simplifies the endpoint paths dramatically. It also eliminates
the overhead of having to grab more ID values for testing and the burden of validating those values in handlers.
Furthermore, because we already tag the endpoints with the object type being accessed, the structure of the generated
Swagger docs should not change.

For example, given the two endpoints:

```yaml
  '/move-task-orders/{moveTaskOrderID}/mto-shipments/{mtoShipmentID}':
    put:
      consumes:
        - application/json
      produces:
        - application/json
      summary: Updates mto shipment
      operationId: updateMTOShipment
      tags:
        - mtoShipment
      parameters:
        - in: path
          name: moveTaskOrderID
          required: true
          format: uuid
          type: string
        - in: path
          name: mtoShipmentID
          required: true
          format: uuid
          type: string
        - in: body
          name: body
          required: true
          schema:
            $ref: '#/definitions/MTOShipment'
        - in: header
          name: If-Match
          type: string
          required: true
      responses:
        [...]
  '/move-task-orders/{moveTaskOrderID}/mto-shipments/{mtoShipmentID}/mto-service-items':
    post:
      consumes:
        - application/json
      produces:
        - application/json
      summary: Creates mto service items
      operationId: createMTOServiceItem
      tags:
        - mtoServiceItem
      parameters:
        - in: path
          name: moveTaskOrderID
          required: true
          format: uuid
          type: string
        - in: path
          name: mtoShipmentID
          required: true
          format: uuid
          type: string
        - in: body
          name: body
          schema:
            description: This may be a MTOServiceItemBasic, MTOServiceItemOriginSIT or etc.
            $ref: '#/definitions/MTOServiceItem'
      responses:
        [...]
```

Because of the different values in the `tags` attribute, they will be grouped separately in the generated docs despite
the detailed nesting in the path.

Therefore, in this context, writing an endpoint path like

`/child-object/{:childID}`

is functionally and structurally congruent to

`/parent-object/{:parentID}/child-object/{:childID}`,

and the former benefits from more simplicity and ease of use.

It is important to note that this endpoint structure is only valid because of the Many-to-One relationship of these
objects. For an API being designed with different data models, such as a Many-to-Many relationship that might need to be
represented with a query like the following:

```sql
SELECT mto_shipments.id, mto_shipments.move_task_order_id, move_task_orders.available_to_prime
FROM mto_shipments
JOIN move_task_orders ON move_task_orders.id = mto_shipments.move_task_order_id
WHERE mto_shipments.id = id_from_path
AND mto_shipments.move_task_order_id = mto_id_from_path
```

it may indeed be necessary to include both IDs in the path. As such, this ADR is making the distinction that endpoint
paths should not be nested with IDs *that are not functionally necessary for security, validation, or retrieval.*

* **Consequences:** All existing nested endpoints will need the following updates:
  * The path and the parameter attributes in the .yaml file will need changes
  * `/gen/` code will need to be regenerated
  * Integration tests will need changes
  * `/cmd/` code for the API CLI will need changes
  * Handlers may not need changes, but all modified endpoints will need to be retested

## Pros and Cons of the Alternatives

### *Leave the codebase as-is*

* `+` Less work now
* `-` Endpoints are inconsistent
* `-` MTO IDs for shipments are still not being validated
* `-` Future developers will not know which convention to use

### *Nest all paths with the IDs of the parent objects*

* `+` Object relationships are explicit and clear to anyone using the API
* `+` Provides an extra value to double-check that the correct record is being updated
* `-` Parent ID value is functionally irrelevant for identifying the correct record and presents an extra failure point
for user input
* `-` Paths with multiple UUIDs quickly become long and unwieldy
* `-` Development burden to validate that the parent ID matches the child's ID
* `-` Testing burden to always require the parent ID

### *Do not include multiple IDs unless functionally necessary. Start a new root uniquely identifiable objects*

* `+` Endpoint paths are simpler and more readable
* `+` Less input required from the user
* `+` Security can be handled discretely in the handlers
* `+` It's clear which object is being accessed and updated, or what the base object is (for lists)
* `-` Relationships between objects is unclear in the yaml file
* `-` Requires a specific data model to be effective
