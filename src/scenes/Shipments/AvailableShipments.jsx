// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import ShipmentCards from 'scenes/Shipments/ShipmentCards';

import { loadAvailableShipments } from './ducks';

export class AvailableShipments extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Available Shipments';
    this.props.loadAvailableShipments();
  }
  render() {
    const { shipments, hasError } = this.props;
    return (
      <div className="usa-grid">
        <h1>Available Shipments</h1>
        {hasError && (
          <Alert type="error" heading="Server Error">
            There was a problem loading the shipments from the server.
          </Alert>
        )}
        {!hasError && <ShipmentCards shipments={shipments} />}
      </div>
    );
  }
}

AvailableShipments.propTypes = {
  loadAvailableShipments: PropTypes.func.isRequired,
  shipments: PropTypes.array,
  hasError: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return {
    shipments: state.availableShipments.shipments,
    hasError: state.availableShipments.hasError,
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadAvailableShipments }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(AvailableShipments);
