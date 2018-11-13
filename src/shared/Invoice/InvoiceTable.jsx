import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import { formatFromBaseQuantity } from 'shared/formatters';
import './InvoicePanel.css';

class InvoiceTable extends PureComponent {
  render() {
    return (
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
                <td>{item.location[0]}</td>
                <td>{formatFromBaseQuantity(item.quantity_1)}</td>
                <td>{formatFromBaseQuantity(item.amount)}</td>
              </tr>
            );
          })}
          <tr>
            <td />
            <td>Total</td>
            <td />
            <td />
            <td>{this.props.totalAmount}</td>
          </tr>
        </tbody>
      </table>
    );
  }
}

InvoiceTable.propTypes = {
  shipmentLineItems: PropTypes.array,
  totalAmount: PropTypes.number,
};

export default InvoiceTable;
