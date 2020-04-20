# Nesting Swagger paths in the Prime API with parent IDs

The Prime API manages Move Task Orders and the child objects for these orders, such as shipments and payment requests.
When writing the paths for the endpoints to access these objects, two distinct strategies have been used:

1. Nest the path for the new endpoint by including the ID of the parent Move Task Order for the object, e.g.
`/move-task-orders/{moveTaskOrderID}/mto-shipments/{mtoShipmentID}`
2. Do not nest under the parent MTO and only include the ID of the object being accessed in the path, e.g.
`/payment-requests/{paymentRequestID}`

The first strategy is problematic because the `moveTaskOrderID` in the path is not validated, so there is no meaningful
connection between that value and the ID of the child object. Only the ID of the child object - the `mtoShipmentID` in this example -
is validated and used for updates. The UUIDs also make the paths extremely long.

The second strategy is problematic because it is inconsistent with the other the endpoints. It also obfuscates the relationships
between objects.

## Considered Alternatives

* Leave the codebase as-is
* Nest all paths with the IDs of the parent objects
* Do not nest paths with parents - start a new root

## Decision Outcome

### Chosen Alternative: *Do not nest paths with parents - start a new root*

* **Justification:** Using a new root for accessing each objects simplifies the endpoint paths dramatically. It also
eliminates the overhead of having to grab more ID values for testing and the burden of validating those values in handlers.
Furthermore, because we already tag all of the endpoints with the object type being accessed, the structure of the generated
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
            description: This may be a MTOServiceItemBasic, MTOServiceItemDOFSIT or etc.
            $ref: '#/definitions/MTOServiceItem'
      responses:
        [...]
```

Because of the different values in the `tags` attribute, they will be grouped separately in the generated docs despite the detailed nesting in the path.

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
* `+` Provides an extra value to check that the correct record is being updated
* `-` Paths with multiple UUIDs quickly become long and unwieldy
* `-` Development burden to validate that the parent ID matches the child's ID
* `-` Testing burden to always require the parent ID
