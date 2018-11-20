import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import { formatFromBaseQuantity, formatCents } from 'shared/formatters';
import { isOfficeSite, isDevelopment } from 'shared/constants.js';
import './InvoicePanel.css';

class InvoiceTable extends PureComponent {
  render() {
    return (
      <div>
        <div className="usa-grid">
          <div className="usa-width-one-half">
            <h5>Unbilled line items</h5>
          </div>
          <div className="usa-width-one-half align-right">
            {isOfficeSite &&
              this.props.shipmentStatus.toUpperCase() === 'DELIVERED' && (
                <button
                  className="button button-secondary"
                  disabled={!this.props.canApprove || !isDevelopment}
                  onClick={this.props.approvePayment}
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
                  <td>${formatCents(item.amount)}</td>
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
