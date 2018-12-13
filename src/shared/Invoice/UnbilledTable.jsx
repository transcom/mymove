import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';
import { withContext } from 'shared/AppContext';

import { isOfficeSite } from 'shared/constants.js';
import LineItemTable from 'shared/Invoice/LineItemTable';
import InvoicePayment from './InvoicePayment';

import './InvoicePanel.css';

export class UnbilledTable extends PureComponent {
  render() {
    const allowPayments =
      this.props.allowPayments &&
      isOfficeSite && //user is an office user
      this.props.context.flags.allowHhgInvoicePayment;

    let itemsComponent = <span className="empty-content">No line items</span>;
    if (this.props.lineItems.length) {
      itemsComponent = (
        <LineItemTable shipmentLineItems={this.props.lineItems} totalAmount={this.props.lineItemsTotal} />
      );
    }

    return (
      <div className="invoice-panel-table-cont">
        <InvoicePayment
          allowPayments={Boolean(allowPayments)}
          cancelPayment={this.props.cancelPayment}
          approvePayment={this.props.approvePayment}
          createInvoiceStatus={this.props.createInvoiceStatus}
        />
        {itemsComponent}
      </div>
    );
  }
}

UnbilledTable.propTypes = {
  lineItems: PropTypes.array,
  lineItemsTotal: PropTypes.number,
  approvePayment: PropTypes.func,
  cancelPayment: PropTypes.func,
  allowPayments: PropTypes.bool,
  createInvoiceStatus: PropTypes.object,
};

export default withContext(UnbilledTable);
