# Improved Swagger/Redux collaboration

**User Story:** _[ticket/issue-number]_ <!-- optional -->

## Accessing models by ID

All data access should be done through selectors and not by directly accessing the global Redux state. Here's an example of a fictitious `ShipmentDisplay` component fetching a specific shipment using its ID (which may have been extracted from a URL).

Define an action creator function in `shared/Entities/modules/shipment.js` to connect

```js
export function getShipment(label, shipmentId, moveId) {
  return swaggerRequest(
    'shipments.getShipment',
    { moveId, shipmentId },
    { label },
  );
}
```

This is also where action creators that need to delegate to multiple other actions live. Here is an example that either updates or creates a shipment based on if `id` is provided:

```js
export function createOrUpdateShipment(label, moveId, shipment, id) {
  if (id) {
    return updateShipment(label, moveId, id, shipment);
  } else {
    return createShipment(label, moveId, shipment);
  }
}
```

Action creator functions should take a `label` argument, which will be used to allow the calling component to determine the status of any requests with that label.

`swaggerRequest` returns a promise, so it is possible to chain behavior onto its result, for example to perform a few requests in sequence.

To access the value after it has been fetched and stored in Redux, we'll need to also create a selector:

```js
// Return a shipment identified by its ID
export function selectShipment(state, id) {
  if (!id) {
    return null;
  }
  return denormalize([id], shipments, state.entities)[0];
}
```

Here is an example of a component using `getShipment` and `selectShipment` as defined above:

```jsx
// import { get } from 'lodash';
// import PropTypes from 'prop-types';
// import React, { Component} from 'react';
// import { connect } from 'react-redux';
// import { bindActionCreators } from 'redux';

import { request } from 'shared/api';
import { lastError } from 'shared/Swagger/ducks';
import { getShipment } from 'shared/Entities/modules/shipments';

// This value is used to identify requests made by this component, since the
// same Swagger operation may be called from multiple components.
const requestLabel = 'ShipmentForm.loadShipment';

export class ShipmentDisplay extends Component {

    componentDidMount() {
        const id = get(this.props, 'shipmentID');
        if (!id) return;

        this.props.getShipment(requestLabel, id);
    }

    render {
        const { shipment, error } = this.props;

        return (
            <div>
                { error && <p>An error has occurred.</p> }

                <p>You are moving on { shipment.requested_move_date }.</p>
            </div>
        );
    }
}

ShipmentDisplay.propTypes = {
    shipment: PropTypes.object,
    error: PropTypes.object,
    shipmentID: string.required,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ request, getShipment }, dispatch);
}

function mapStateToProps(state, props) {
  return {
    shipment: selectShipment(props.shipmentID),
    error: lastError(state, requestLabel),
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(ShipmentForm);
```

Note the use of `REQUEST_LABEL` to allow for notification of any errors to requests.

## Storing UI state

State that is specific to the UI should be set by dispatching an action and accessed using a selector. Here is an
example of how this might work for managing which of a list of Shipments is currently selected in the UI:

```javascript
// import { get } from 'lodash';
// import PropTypes from 'prop-types';
// import React, { Component} from 'react';
// import { connect } from 'react-redux';
// import { bindActionCreators } from 'redux';

import { request } from 'shared/api';
import { lastError } from 'shared/Swagger/ducks';
import { allShipments, selectShipment } from 'shared/Entities/modules/shipments';
import { setCurrentShipmentID, currentShipmentID } from 'shared/UI/ducks';

const requestLabel = 'ShipmentForm.loadShipments';

export class ShipmentList extends Component {
    componentDidMount() {
        const id = get(this.props, 'shipmentID');
        if (!id) return;

        this.props.request(requestLabel, 'shipments.listShipment');
    }

    shipmentClicked = (id) => {
       this.props.setCurrentShipmentID(id);
    }

    render {
        const { shipments, selectedShipment, error } = this.props;

        return (
            <div>
                { error && <p>An error has occurred.</p> }

                <ul>
                    { shipments.map(shipment => (<li>
                        <button onClick={this.shipmentClicked.bind(shipment.id)}> { shipment.id } </button>
                      </li>))}
                </ul>

                <p>The selected shipment is { selectedShipment.id }.</p>
            </div>
        );
    }
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ request, setCurrentShipmentID }, dispatch);
}

function mapStateToProps(state) {
  return {
    shipments: allShipments(state),
    selectedShipment: selectShipment(currentShipmentID(state)),
    error: lastError(state, requestLabel),
  };
}
```

## Improved debugging experience

By directing all Swagger Client calls through the `request` function, we can have a centralized place to manage how to track
the lifecycle of the request. This allows us to dispatch actions to Redux that represent these events, currently `@@swagger/${operation}/START`, `@@swagger/${operation}/SUCCESS` and `@@swagger/${operation}/FAILURE`. These actions will appear in the Redux debugger along with any other state changes.

## Redux store data layout

The patterns described above will utilize the following internal Redux store layout. This should
generally be considered an implementation detail and we should strive to avoid coupling any Components
to this structure directly.

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
