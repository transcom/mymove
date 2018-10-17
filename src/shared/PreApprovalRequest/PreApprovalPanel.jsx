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
  deleteShipmentAccessorial,
  deleteShipmentAccessorialLabel,
  approveShipmentAccessorial,
  approveShipmentAccessorialLabel,
} from 'shared/Entities/modules/shipmentAccessorials';
import { selectShipmentAccessorials } from 'shared/Entities/modules/shipmentAccessorials';
import { selectTariff400ngItems } from 'shared/Entities/modules/tariff400ngItems';

export class PreApprovalPanel extends Component {
  constructor() {
    super();
    this.state = { isActionable: true };
  }
  onSubmit = createPayload => {
    return this.props.createShipmentAccessorial(createShipmentAccessorialLabel, this.props.shipmentId, createPayload);
  };
  onEdit = () => {
    console.log('onEdit hit');
  };
  onDelete = shipmentAccessorialId => {
    if (window.confirm('Are you sure you want to delete this pre approval request?')) {
      this.props.deleteShipmentAccessorial(deleteShipmentAccessorialLabel, shipmentAccessorialId);
    }
  };
  onApproval = shipmentAccessorialId => {
    let response = this.props.approveShipmentAccessorial(approveShipmentAccessorialLabel, shipmentAccessorialId);
    let resolved = result => {
      //do something here if successful
      console.log('Got response success: ');
      console.log(result);
    };
    let rejected = result => {
      //do something here if unsucessful
      console.error('Got response error: ');
      console.error(result);
    };
    response.then(resolved, rejected);
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
  return bindActionCreators(
    { createShipmentAccessorial, deleteShipmentAccessorial, approveShipmentAccessorial },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(PreApprovalPanel);
