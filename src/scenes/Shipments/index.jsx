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
    this.props.loadShipments();
  }
  render() {
    const { shipments, hasError } = this.props;
    const shipmentsStatus = this.props.match.params.shipmentsStatus;

    // Title with capitalized shipment status
    const capShipmentsStatus = shipmentsStatus.charAt(0).toUpperCase() + shipmentsStatus.slice(1);

    // Handle cases of users entering invalid shipment types
    if (shipmentsStatus !== 'awarded' && shipmentsStatus !== 'available' && shipmentsStatus !== 'all') {
      return (
        <Alert type="error" heading="Invalid Shipment Type Error">
          You've attempted to access an inaccessible route. Invalid Shipment Status: {shipmentsStatus}.
        </Alert>
      );
    }

    // TODO: Move to reducer and memoize this, possibly including tdl grouping.
    // Inquire with Erin how should we allow users to sort.
    const filteredShipments = shipments.filter(shipment => {
      return (
        shipmentsStatus === 'all' ||
        (shipment.transportation_service_provider_id && shipmentsStatus === 'awarded') ||
        (!shipment.transportation_service_provider_id && shipmentsStatus === 'available')
      );
    });

    const groupedShipments = filteredShipments.reduce((groups, shipment) => {
      groups[shipment.traffic_distribution_list_id] = groups[shipment.traffic_distribution_list_id] || [];
      groups[shipment.traffic_distribution_list_id].push(shipment);

      return groups;
    }, {});

    const cards = [];
    for (let tdl in groupedShipments) {
      const tdlID = tdl.substr(0, 6);
      const shipments = groupedShipments[tdl]; // eslint-disable-line security/detect-object-injection
      cards.push(
        <div className="tdl-box" key={tdlID}>
          <h3>TDL #{tdlID}</h3>
          <ShipmentCards shipments={shipments} />
        </div>,
      );
    }

    return (
      <div className="usa-grid">
        <h1>{capShipmentsStatus} Shipments</h1>
        {hasError && (
          <Alert type="error" heading="Server Error">
            There was a problem loading the shipments from the server.
          </Alert>
        )}
        {!hasError && cards}
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
    shipments: state.shipments.shipments,
    hasError: state.shipments.hasError,
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadShipments }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Shipments);
