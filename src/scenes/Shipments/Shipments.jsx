// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import ShipmentCards from 'scenes/Shipments/ShipmentCards';

import { loadShipments } from './ducks';

export class Shipments extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Shipments';
    console.log(this.props);
    debugger;
    this.props.loadShipments();
  }
  render() {
    const { shipments, hasError } = this.props;
    return (
      <div className="usa-grid">
        <h1>Shipments</h1>
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

Shipments.propTypes = {
  loadShipments: PropTypes.func.isRequired,
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
  return bindActionCreators({ loadShipments }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Shipments);
