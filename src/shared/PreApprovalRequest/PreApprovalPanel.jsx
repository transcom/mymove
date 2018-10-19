import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import BasicPanel from 'shared/BasicPanel';
import { isOfficeSite } from 'shared/constants.js';
import PreApprovalTable from 'shared/PreApprovalRequest/PreApprovalTable.jsx';
import Creator from 'shared/PreApprovalRequest/Creator';

import {
  createShipmentAccessorial,
  createShipmentAccessorialLabel,
  deleteShipmentAccessorial,
  deleteShipmentAccessorialLabel,
  approveShipmentAccessorial,
  approveShipmentAccessorialLabel,
  updateShipmentAccessorial,
  updateShipmentAccessorialLabel,
} from 'shared/Entities/modules/shipmentAccessorials';
import { selectShipmentAccessorials } from 'shared/Entities/modules/shipmentAccessorials';
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
    return this.props.createShipmentAccessorial(createShipmentAccessorialLabel, this.props.shipmentId, createPayload);
  };
  onEdit = (shipmentAccessorialId, editPayload) => {
    this.props.updateShipmentAccessorial(updateShipmentAccessorialLabel, shipmentAccessorialId, editPayload);
  };
  onDelete = shipmentAccessorialId => {
    this.props.deleteShipmentAccessorial(deleteShipmentAccessorialLabel, shipmentAccessorialId);
  };
  onApproval = shipmentAccessorialId => {
    this.props.approveShipmentAccessorial(approveShipmentAccessorialLabel, shipmentAccessorialId);
  };
  onFormActivation = isFormActive => {
    this.setState({ isRequestActionable: !isFormActive });
  };
  onRequestActivation = isRequestActive => {
    this.setState({ isCreatorActionable: !isRequestActive });
  };
  render() {
    return (
      <div className="accessorial-panel">
        <BasicPanel title={'Pre-Approval Requests'}>
          <PreApprovalTable
            tariff400ngItems={this.props.tariff400ngItems}
            shipmentAccessorials={this.props.shipmentAccessorials}
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
  shipmentAccessorials: PropTypes.array,
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
    { createShipmentAccessorial, deleteShipmentAccessorial, approveShipmentAccessorial, updateShipmentAccessorial },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(PreApprovalPanel);
