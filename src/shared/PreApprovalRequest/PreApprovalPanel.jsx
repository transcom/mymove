import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import BasicPanel from 'shared/BasicPanel';
import { isOfficeSite } from 'shared/constants.js';
import Alert from 'shared/Alert';
import PreApprovalTable from 'shared/PreApprovalRequest/PreApprovalTable.jsx';
import Creator from 'shared/PreApprovalRequest/Creator';

import {
  createShipmentLineItem,
  createShipmentLineItemLabel,
  deleteShipmentLineItem,
  deleteShipmentLineItemLabel,
  approveShipmentLineItem,
  approveShipmentLineItemLabel,
  updateShipmentLineItem,
  updateShipmentLineItemLabel,
} from 'shared/Entities/modules/shipmentLineItems';
import { selectSortedPreApprovalShipmentLineItems } from 'shared/Entities/modules/shipmentLineItems';
import { selectSortedPreApprovalTariff400ngItems } from 'shared/Entities/modules/tariff400ngItems';

export class PreApprovalPanel extends Component {
  constructor() {
    super();
    this.state = {
      isRequestActionable: true,
      isCreatorActionable: true,
      error: null,
    };
  }
  closeError = () => {
    this.setState({ error: null });
  };
  onSubmit = createPayload => {
    return this.props.createShipmentLineItem(createShipmentLineItemLabel, this.props.shipmentId, createPayload);
  };
  onEdit = (shipmentLineItemId, editPayload) => {
    this.props.updateShipmentLineItem(updateShipmentLineItemLabel, shipmentLineItemId, editPayload);
  };
  onDelete = shipmentLineItemId => {
    this.props.deleteShipmentLineItem(deleteShipmentLineItemLabel, shipmentLineItemId).catch(err => {
      this.setState({ error: true });
    });
  };
  onApproval = shipmentLineItemId => {
    this.props.approveShipmentLineItem(approveShipmentLineItemLabel, shipmentLineItemId);
  };
  onFormActivation = isFormActive => {
    this.setState({ isRequestActionable: !isFormActive });
  };
  onRequestActivation = isRequestActive => {
    this.setState({ isCreatorActionable: !isRequestActive });
  };
  render() {
    return (
      <div className="pre-approval-panel">
        <BasicPanel title={'Pre-Approval Requests'}>
          {this.state.error && (
            <Alert type="error" heading="Oops, something went wrong!" onRemove={this.closeError}>
              <span className="warning--header">Please refresh the page and try again.</span>
            </Alert>
          )}
          <PreApprovalTable
            tariff400ngItems={this.props.tariff400ngItems}
            shipmentLineItems={this.props.shipmentLineItems}
            isActionable={this.state.isRequestActionable}
            onRequestActivation={this.onRequestActivation}
            onEdit={this.onEdit}
            onDelete={this.onDelete}
            onApproval={isOfficeSite ? this.onApproval : null}
          />
          {this.state.isCreatorActionable && (
            <Creator
              tariff400ngItems={this.props.tariff400ngItems}
              savePreApprovalRequest={this.onSubmit}
              onFormActivation={this.onFormActivation}
            />
          )}
        </BasicPanel>
      </div>
    );
  }
}

PreApprovalPanel.propTypes = {
  shipmentLineItems: PropTypes.array,
  tariff400ngItems: PropTypes.array,
  shipmentId: PropTypes.string,
};

function mapStateToProps(state, ownProps) {
  return {
    shipmentLineItems: selectSortedPreApprovalShipmentLineItems(state, ownProps.shipmentId),
    tariff400ngItems: selectSortedPreApprovalTariff400ngItems(state),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { createShipmentLineItem, deleteShipmentLineItem, approveShipmentLineItem, updateShipmentLineItem },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(PreApprovalPanel);
