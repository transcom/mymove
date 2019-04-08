import React, { PureComponent } from 'react';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import { withContext } from 'shared/AppContext';

import { formatCents } from 'shared/formatters';
import { displayBaseQuantityUnits } from 'shared/lineItems';

import './InvoicePanel.css';

class LineItemTable extends PureComponent {
  render() {
    const showItem35Missing = item => isRobust35A(item, this.props);
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
                  <td>
                    {item.tariff400ng_item.item}
                    {showItem35Missing(item) && (
                      <span>
                        <br />
                        <span className="shipment-line-item-warning">Missing actual amount</span>
                      </span>
                    )}
                  </td>
                  <td>{item.location[0]}</td>
                  <td>{displayBaseQuantityUnits(item)}</td>
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

function isRobust35A(item, props) {
  const robustAccessorialFlag = get(props, 'context.flags.robustAccessorial', false);
  return (
    robustAccessorialFlag &&
    item.tariff400ng_item.code === '35A' &&
    item.estimate_amount_cents &&
    !item.actual_amount_cents
  );
}

export default withContext(LineItemTable);
