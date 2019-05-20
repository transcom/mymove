# How To Call Swagger Endpoints from React

## 1. Verify the Schema is Defined

For each model type returned by the backend, there needs to be an `Entity` defined and exported in `src/shared/Entities/schema.js`.

Here is the definition for `Shipment`:

```javascript
export const shipment = new schema.Entity('shipments');

// add any embedded objects that should be extracted during normalization
shipment.define({
  pickup_address: address,
  secondary_pickup_address: address,
  delivery_address: address,
  partial_sit_delivery_address: address,
});

export const shipments = new schema.Array(shipment);
```

## 2. Call the Swagger Operation

Add a function to `src/shared/Entities/modules/$MODEL.js` that calls the `operationId` defined in the
swagger YAML. Action creator functions should take a `label` argument, which will be used to allow the calling component to determine the status of any requests with that label.

`swaggerRequest` returns a promise, so it is possible to chain behavior onto its result, for example to perform a few requests in sequence.

```javascript
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

export function getShipment(label, shipmentId, moveId) {
  return swaggerRequest(
    getClient,                  // function returning a promise that will resolve to a Swagger client instance
    'shipments.getShipment',    // what operation to perform, including tag namespace
    { moveId, shipmentId },     // parameters to pass to the operation
    { label },                  // optional params for swaggerRequest, such as label
  );
}
```

By directing all Swagger Client calls through the `swaggerRequest` function, we can have a centralized place to manage how to track
the lifecycle of the request. This allows us to dispatch actions to Redux that represent these events, currently `@@swagger/${operation}/START`, `@@swagger/${operation}/SUCCESS` and `@@swagger/${operation}/FAILURE`. These actions will appear in the Redux debugger along with any other state changes.

## 3. Dispatch an Action when Component Mounts

The following pattern, using `onDidMount`, allows the data fetching to be defined outside the component.

```javascript
export class ShipmentDisplay extends Component {

    componentDidMount() {
        this.props.onDidMount && this.props.onDidMount();
    }

    render {
        const { shipment } = this.props;

        return (
            <div>
                <p>You are moving on { shipment.requested_move_date }.</p>
            </div>
        );
    }
}

ShipmentDisplay.propTypes = {
    shipmentID: PropTypes.string.isRequired,

    onDidMount: PropTypes.function,
    shipment: PropTypes.object,
};

const requestLabel = 'ShipmentDisplay.getShipment';

function mapDispatchToProps(dispatch, ownProps) {
    return {
        onDidMount: function() {
            dispatch(getShipment(requestLabel, ownProps.shipmentID));        }
    };
}

function mapStateToProps(state, ownProps) {
  return {
    shipment: selectShipment(ownProps.shipmentID),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ShipmentDisplay);
```

If you need to load data based on a value that isn't passed in as a `prop`, it's best to embed another component and pass that value into it as a `prop`. This can be thought of as an extension of the container pattern.

## 4. Use a Selector to Access the Data

All data access should be done through selectors and not by directly accessing the global Redux state.

Add a function to `src/shared/Entities/modules/$MODEL.js` that returns the value from Redux. This example uses `denormalize`:

```javascript
// Return a shipment identified by its ID
export function selectShipment(state, id) {
  if (!id) {
    return null;
  }
  return denormalize([id], shipments, state.entities)[0];
}
```

This one returns a value that doesn't need `denormalize`:

```javascript
// Return a shipment identified by its ID
export function selectShipment(state, id) {
  if (!id) {
    return null;
  }
  return get(state, `entities.shipments.${id}`);
}
```

## 5. Handle Fetch Errors

The `lastError` selector provides access to the last error for a specified request label.

```javascript
import { lastError } from 'shared/Swagger/selectors';

export class ShipmentDisplay extends Component {

    componentDidMount() {
        this.props.onDidMount && this.props.onDidMount();
    }

    render {
        const { shipment, error } = this.props;

        return (
            { error && <p>An error has occurred.</p> }

            <div>
                <p>You are moving on { shipment.requested_move_date }.</p>
            </div>
        );
    }
}

ShipmentDisplay.propTypes = {
    shipmentID: PropTypes.string.isRequired,

    onDidMount: PropTypes.function,
    shipment: PropTypes.object,
    error: PropTypes.object,
};

const requestLabel = 'ShipmentDisplay.getShipment';

function mapDispatchToProps(dispatch, ownProps) {
    return {
        onDidMount: function() {
            dispatch(getShipment(requestLabel, ownProps.shipmentID));
        }
    };
}

function mapStateToProps(state, ownProps) {
  return {
    shipment: selectShipment(ownProps.shipmentID),
    error: lastError(state, requestLabel),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ShipmentDisplay);
```
