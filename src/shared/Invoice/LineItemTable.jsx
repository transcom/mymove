import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import { formatFromBaseQuantity, formatCents } from 'shared/formatters';

import './InvoicePanel.css';

class LineItemTable extends PureComponent {
  render() {
    return (
      <div>
        {this.props.title}
        <table cellSpacing={0}>
          <tbody>
            <tr data-cy="table--header">
              <th>Code</th>
              <th>Item</th>
              <th>Loc</th>
              <th>Base quantity</th>
              <th>Inv amt</th>
            </tr>
            {this.props.shipmentLineItems.map(item => {
              return (
                <tr key={item.id} data-cy="table--item">
                  <td>{item.tariff400ng_item.code}</td>
                  <td>{item.tariff400ng_item.item}</td>
                  <td>{item.location[0]}</td>
                  <td>{formatFromBaseQuantity(item.quantity_1)}</td>
                  <td>${formatCents(item.amount_cents)}</td>
                </tr>
              );
            })}
            <tr data-cy="table--total">
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

LineItemTable.propTypes = {
  title: PropTypes.element,
  shipmentLineItems: PropTypes.array,
  totalAmount: PropTypes.number,
};

export default LineItemTable;
