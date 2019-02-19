import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';
import PreApprovalRequest from 'shared/PreApprovalRequest/PreApprovalRequest.jsx';

import './PreApprovalRequest.css';

export class PreApprovalTable extends PureComponent {
  state = { actionRequestId: null };
  isRequestActive = id => {
    return isActive => {
      this.props.onRequestActivation(isActive);
      if (isActive) {
        this.setState({ actionRequestId: id });
      } else {
        this.setState({ actionRequestId: null });
      }
    };
  };
  render() {
    const { shipmentLineItems, isActionable, onEdit, onApproval, onDelete } = this.props;
    // If there are no shipment line items, don't show the table at all.
    return (
      <div className="pre-approval-panel-table-cont">
        {shipmentLineItems.length > 0 && (
          <table cellSpacing={0}>
            <tbody>
              <tr>
                <th style={{ width: '5%' }}>Code</th>
                <th style={{ width: '26%' }}>Item</th>
                <th style={{ width: '5%' }}>Loc</th>
                <th style={{ width: '30%' }}>Details</th>
                <th style={{ width: '12%' }}>Submitted</th>
                <th style={{ width: '12%' }}>Status</th>
                <th style={{ width: '12%' }}>&nbsp;</th>
              </tr>
              {shipmentLineItems.map(row => {
                let requestIsActionable =
                  isActionable && (this.state.actionRequestId === null || this.state.actionRequestId === row.id);
                return (
                  <PreApprovalRequest
                    key={row.id}
                    shipmentLineItem={row}
                    onEdit={onEdit}
                    onApproval={onApproval}
                    onDelete={onDelete}
                    isActive={this.isRequestActive(row.id)}
                    isActionable={requestIsActionable}
                    tariff400ngItems={this.props.tariff400ngItems}
                  />
                );
              })}
            </tbody>
          </table>
        )}
      </div>
    );
  }
}

PreApprovalTable.propTypes = {
  shipmentLineItems: PropTypes.array,
  tariff400ngItems: PropTypes.array,
  isActionable: PropTypes.bool,
  onEdit: PropTypes.func,
  onRequestActivation: PropTypes.func,
  onDelete: PropTypes.func,
  onApproval: PropTypes.func,
};

export default PreApprovalTable;
