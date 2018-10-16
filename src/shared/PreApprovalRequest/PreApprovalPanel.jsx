import React, { Component } from 'react';
import BasicPanel from 'shared/BasicPanel';
import PropTypes from 'prop-types';
import { isOfficeSite } from 'shared/constants.js';

import PreApprovalTable from 'shared/PreApprovalRequest/PreApprovalTable.jsx';
import Creator from 'shared/PreApprovalRequest/Creator';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import {
  createShipmentAccessorial,
  createShipmentAccessorialLabel,
} from 'shared/Entities/modules/shipmentAccessorials';
import { selectShipmentAccessorials } from 'shared/Entities/modules/shipmentAccessorials';
import { selectTariff400ngItems } from 'shared/Entities/modules/tariff400ngItems';

export class PreApprovalPanel extends Component {
  constructor() {
    super();
    this.state = {
      isActionable: true,
    };
  }
  onSubmit = createPayload => {
    return new Promise(
      function(resolve, reject) {
        // do a thing, possibly async, then…
        this.props.createShipmentAccessorial(createShipmentAccessorialLabel, this.props.shipmentId, createPayload);
        setTimeout(function() {
          resolve('success');
        }, 50);
      }.bind(this),
    );
  };
  onEdit = () => {
    console.log('onEdit hit', this.props);
  };
  onDelete = () => {
    console.log('onDelete hit');
  };
  onApproval = () => {
    console.log('onApproval hit');
  };
  onFormActivation = active => {
    this.setState({ isActionable: active });
  };
  render() {
    return (
      <div>
        <BasicPanel title={'Pre-Approval Requests'}>
          <PreApprovalTable
            shipment_accessorials={this.props.shipment_accessorials}
            isActionable={this.state.isActionable}
            onEdit={this.onEdit}
            onDelete={this.onDelete}
            onApproval={isOfficeSite ? this.onApproval : null}
          />
          <Creator
            tariff400ngItems={this.props.tariff400ngItems}
            savePreApprovalRequest={this.onSubmit}
            onFormActivation={this.onFormActivation}
          />
        </BasicPanel>
      </div>
    );
  }
}

PreApprovalPanel.propTypes = {
  shipment_accessorials: PropTypes.array,
  tariff400ngItems: PropTypes.array,
  shipmentId: PropTypes.string,
};

function mapStateToProps(state) {
  return {
    shipmentAccessorials: selectShipmentAccessorials(state),
    tariff400ngItems: selectTariff400ngItems(state),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createShipmentAccessorial }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(PreApprovalPanel);
