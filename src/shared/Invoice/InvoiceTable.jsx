import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import { formatFromBaseQuantity, formatCents } from 'shared/formatters';
import './InvoicePanel.css';

class InvoiceTable extends PureComponent {
  render() {
    return (
      <div>
        <h5>Unbilled line items</h5>
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
  totalAmount: PropTypes.number,
};

export default InvoiceTable;
