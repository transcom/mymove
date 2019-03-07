import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import {
  getTspForShipment,
  selectTransportationServiceProviderForShipment,
} from 'shared/Entities/modules/transportationServiceProviders';
import { getPublicShipment } from 'shared/Entities/modules/shipments';

export class TransportationServiceProviderContactInfo extends Component {
  componentDidMount() {
    const shipmentId = this.props.shipmentId;
    this.props.getTspForShipment(shipmentId);
  }

  render() {
    const { transportationServiceProvider, showFileAClaimInfo } = this.props;
    if (showFileAClaimInfo) {
      return (
        <div className="step">
          <div className="title">File a Claim</div>
          <div>
            If you have household goods damaged or lost during the move, contact {transportationServiceProvider.name} to
            file a claim: {transportationServiceProvider.poc_general_phone}. If, after attempting to work with them, you
            do not feel that you are receiving adequate compensation, contact the Military Claims Office for help.
          </div>
        </div>
      );
    } else {
      return (
        <div className="titled_block transportation-service-provider-contact-info">
          <div>
            <strong>{transportationServiceProvider.name}</strong>
          </div>
          <div>{transportationServiceProvider.poc_general_phone}</div>
        </div>
      );
    }
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
