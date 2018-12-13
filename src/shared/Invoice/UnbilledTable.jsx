import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import { isOfficeSite, isDevelopment } from 'shared/constants.js';
import LineItemTable from 'shared/Invoice/LineItemTable';
import InvoicePayment from './InvoicePayment';

import './InvoicePanel.css';

export class UnbilledTable extends PureComponent {
  render() {
    let itemsComponent = <span className="empty-content">No line items</span>;
    const allowPayment =
      this.props.allowPayments &&
      isOfficeSite && //user is an office user
      isDevelopment; //only for development env

    if (this.props.lineItems.length) {
      itemsComponent = (
        <LineItemTable shipmentLineItems={this.props.lineItems} totalAmount={this.props.lineItemsTotal} />
      );
    }

    return (
      <div className="invoice-panel-table-cont">
        <InvoicePayment
          allowPayment={Boolean(allowPayment)}
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
  isInvoiceInDraft: PropTypes.bool,
  lineItems: PropTypes.array,
  lineItemsTotal: PropTypes.number,
  draftInvoice: PropTypes.func,
  approvePayment: PropTypes.func,
  cancelPayment: PropTypes.func,
  allowPayments: PropTypes.bool,
  createInvoiceStatus: PropTypes.object,
};

export default UnbilledTable;
