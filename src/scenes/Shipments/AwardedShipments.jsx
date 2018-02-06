// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import ShipmentCards from 'scenes/Shipments/ShipmentCards';

import { loadAwardedShipments } from './ducks';

export class AwardedShipments extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Awarded Shipments';
    this.props.loadAwardedShipments();
  }
  render() {
    const { shipments, hasError } = this.props;
    return (
      <div className="usa-grid">
        <h1>Awarded Shipments</h1>
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

AwardedShipments.propTypes = {
  loadAwardedShipments: PropTypes.func.isRequired,
  shipments: PropTypes.array,
  hasError: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return {
    shipments: state.awardedShipments.shipments,
    hasError: state.awardedShipments.hasError,
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadAwardedShipments }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(AwardedShipments);
