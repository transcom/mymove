# How To Store UI State in Redux

State that is specific to the UI should be set by dispatching an action and accessed using a selector. Here is an
example of how this might work for managing which of a list of Shipments is currently selected in the UI:

```javascript
import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component} from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { lastError } from 'shared/Swagger/ducks';
import { allShipments, selectShipment } from 'shared/Entities/modules/shipments';
import { setCurrentShipmentID, currentShipmentID } from 'shared/UI/ducks';

const requestLabel = 'ShipmentForm.loadShipments';

export class ShipmentList extends Component {
    componentDidMount() {
        const id = get(this.props, 'shipmentID');
        if (!id) return;

        this.props.listShipments(requestLabel);
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
    selectedShipment: selectShipment(state, currentShipmentID(state)),
    error: lastError(state, requestLabel),
  };
}
```

Note that the above use of defining an inline event handler for `onClick` is not considered a
best practice. This technique is used above only for its brevity.
