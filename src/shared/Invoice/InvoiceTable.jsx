import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import { formatFromBaseQuantity, formatCents } from 'shared/formatters';
import Alert from 'shared/Alert';

import { isOfficeSite, isDevelopment } from 'shared/constants.js';
import './InvoicePanel.css';

class InvoiceTable extends PureComponent {
  state = {
    isConfirmationFlowVisible: this.props.isConfirmationFlowVisible,
  };

  static defaultProps = {
    isConfirmationFlowVisible: false,
  };

  approvePayment = () => {
    this.setState({ isConfirmationFlowVisible: true });
  };

  cancelPayment = () => {
    this.setState({ isConfirmationFlowVisible: false });
  };

  render() {
    return (
      <div className="invoice-panel-table-cont">
        {isOfficeSite && this.state.isConfirmationFlowVisible ? (
          <div>
            <Alert type="warning" heading="Approve payment?">
              <span className="warning--header">Please make sure you've double-checked everything.</span>
              <button className="button usa-button-secondary" onClick={this.cancelPayment}>
                Cancel
              </button>
              <button className="button usa-button-primary"> Approve</button>
            </Alert>
          </div>
        ) : null}
        <div className="usa-grid-full invoice-panel-header-cont">
          <div className="usa-width-one-half">
            <h5>Unbilled line items</h5>
          </div>
          <div className="usa-width-one-half align-right">
            {isOfficeSite &&
              !this.state.isConfirmationFlowVisible &&
              this.props.shipmentStatus.toUpperCase() === 'DELIVERED' && (
                <button
                  className="button button-secondary"
                  disabled={!this.props.canApprove || !isDevelopment}
                  onClick={this.approvePayment}
                >
                  Approve Payment
                </button>
              )}
          </div>
        </div>
        <table cellSpacing={0}>
          <tbody>
            <tr>
              <th>Code</th>
              <th>Item</th>
              <th>Loc.</th>
              <th>Base Quantity</th>
              <th>Inv. amt.</th>
            </tr>
            {this.props.shipmentLineItems.map(item => {
              return (
                <tr key={item.id}>
                  <td>{item.tariff400ng_item.code}</td>
                  <td>{item.tariff400ng_item.item}</td>
                  <td>{item.location[0] + item.location.substring(1).toLowerCase()}</td>
                  <td>{formatFromBaseQuantity(item.quantity_1)}</td>
                  <td>${formatCents(item.amount_cents)}</td>
                </tr>
              );
            })}
            <tr>
              <td />
              <td>Total</td>
              <td />
              <td />
              <td>${formatCents(this.props.totalAmount)}</td>
            </tr>
          </tbody>
        </table>
      </div>
    );
  }
}

InvoiceTable.propTypes = {
  shipmentLineItems: PropTypes.array,
  shipmentStatus: PropTypes.string,
  totalAmount: PropTypes.number,
  approvePayment: PropTypes.func,
  canApprove: PropTypes.bool,
};

export default InvoiceTable;
