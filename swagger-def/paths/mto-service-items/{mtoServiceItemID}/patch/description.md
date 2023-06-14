Updates `MTOServiceItems` after creation. Not all service items or fields may be
updated, please see details below.

This endpoint supports different body definitions. In the `modelType` field below,
select the `modelType` corresponding to the service item you wish to update and
the documentation will update with the new definition.

- Addresses: You can add a new SIT Destination final address using this endpoint
  (and must use this endpoint to do so), but you cannot update an existing one.
  Please use the
  [createSITAddressUpdateRequest](#operation/createSITAddressUpdateRequest)
  endpoint instead.

To create a service item, please use
[createMTOServiceItem](#operation/createMTOServiceItem) endpoint.
