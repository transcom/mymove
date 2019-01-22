# How To Store Data in Redux

The specific layout of data within the Redux store should generally be considered an implementation detail and we should strive to avoid coupling any Components to this structure directly. Selectors provide the best way to decouple component data access from store layout.

The current layout of data in Redux, however, is the following:

```javascript
{
    entities: {
        shipments: {
            '123e4567-e89b-12d3-a456-426655440000': { /* shipment properties */ },
        },
        addresses: {
            '123e4567-e89b-12d3-a456-426655440000': { /* address properties */ },
        }
    },
    requests: {
       byID: {
           'req_0': { /* request properties */},
           'req_1': { /* request properties */},
       },
       errored: {
           'req_1': { /* request properties */},
       },
       lastErrorByLabel: {
           'ShipmentForm.loadShipments': { /* error properties */ },
       }
    },
    ui: {
        'currentShipmentID': '123e4567-e89b-12d3-a456-426655440000',
    },
}
```
