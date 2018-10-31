import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import BasicPanel from 'shared/BasicPanel';
import { isOfficeSite } from 'shared/constants.js';
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
import { selectSortedShipmentLineItems } from 'shared/Entities/modules/shipmentLineItems';
import { selectTariff400ngItems } from 'shared/Entities/modules/tariff400ngItems';

export class PreApprovalPanel extends Component {
  constructor() {
    super();
    this.state = {
      isRequestActionable: true,
      isCreatorActionable: true,
    };
  }
  onSubmit = createPayload => {
    return this.props.createShipmentLineItem(createShipmentLineItemLabel, this.props.shipmentId, createPayload);
  };
  onEdit = (shipmentLineItemId, editPayload) => {
    this.props.updateShipmentLineItem(updateShipmentLineItemLabel, shipmentLineItemId, editPayload);
  };
  onDelete = shipmentLineItemId => {
    this.props.deleteShipmentLineItem(deleteShipmentLineItemLabel, shipmentLineItemId);
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

function mapStateToProps(state) {
  return {
    shipmentLineItems: selectSortedShipmentLineItems(state),
    tariff400ngItems: selectTariff400ngItems(state),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { createShipmentLineItem, deleteShipmentLineItem, approveShipmentLineItem, updateShipmentLineItem },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(PreApprovalPanel);
