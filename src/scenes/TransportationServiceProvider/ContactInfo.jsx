import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import {
  getTspForShipmentLabel,
  getTspForShipment,
  selectTransportationServiceProviderForShipment,
} from 'shared/Entities/modules/transportationServiceProviders';
import { getPublicShipment } from 'shared/Entities/modules/shipments';

export class TransportationServiceProviderContactInfo extends Component {
  componentDidMount() {
    const shipmentId = this.props.shipmentId;
    this.props.getTspForShipment(getTspForShipmentLabel, shipmentId);
    this.props.getPublicShipment('Shipments.getPublicShipment', shipmentId);
  }

  render() {
    const { transportationServiceProvider } = this.props;
    return (
      <div className="titled_block transportation-service-provider-contact-info">
        <div>
          <b>{transportationServiceProvider.name}</b>
        </div>
        <div>{transportationServiceProvider.poc_general_phone}</div>
      </div>
    );
  }
}

function mapStateToProps(state, props) {
  return {
    transportationServiceProvider: selectTransportationServiceProviderForShipment(state, props.shipmentId),
  };
}

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      getTspForShipment,
      getPublicShipment,
    },
    dispatch,
  );

export default connect(mapStateToProps, mapDispatchToProps)(TransportationServiceProviderContactInfo);
