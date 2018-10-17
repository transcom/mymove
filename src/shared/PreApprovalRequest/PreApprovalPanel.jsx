import React, { Component } from 'react';
import BasicPanel from 'shared/BasicPanel';
import PropTypes from 'prop-types';
import { isOfficeSite } from 'shared/constants.js';

import PreApprovalTable from 'shared/PreApprovalRequest/PreApprovalTable.jsx';
import Creator from 'shared/PreApprovalRequest/Creator';

export class PreApprovalPanel extends Component {
  constructor() {
    super();
    this.state = {
      isRequestActionable: true,
      isCreatorActionable: true,
    };
  }
  onSubmit = values => {
    return new Promise(function(resolve, reject) {
      // do a thing, possibly async, then…
      setTimeout(function() {
        console.log('onSubmit async', values);
        resolve('success');
      }, 50);
    });
  };
  onEdit = () => {
    console.log('onEdit hit');
  };
  onDelete = () => {
    console.log('onDelete hit');
  };
  onApproval = () => {
    console.log('onApproval hit');
  };
  onFormActivation = isFormActive => {
    this.setState({ isRequestActionable: !isFormActive });
  };
  onRequestActivation = isRequestActive => {
    this.setState({ isCreatorActionable: !isRequestActive });
  };
  render() {
    return (
      <div>
        <BasicPanel title={'Pre-Approval Requests'}>
          <PreApprovalTable
            shipment_accessorials={this.props.shipment_accessorials}
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
  shipment_accessorials: PropTypes.array,
  tariff400ngItems: PropTypes.array,
};
export default PreApprovalPanel;
