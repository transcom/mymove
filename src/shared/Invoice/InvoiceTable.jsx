import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';

import { formatDateTime } from 'shared/formatters';
import {
  selectInvoiceShipmentLineItems,
  selectTotalFromInvoicedLineItems,
} from 'shared/Entities/modules/shipmentLineItems';
import LineItemTable from 'shared/Invoice/LineItemTable';

import './InvoicePanel.css';

class InvoiceTable extends PureComponent {
  render() {
    const tableTitle = (
      <div className="invoice-panel-header-cont">
        <div>
          <h5 data-cy="invoice--detail">
            Invoice {this.props.invoice.invoice_number}{' '}
            <span className="detail">
              Approved: <strong>{formatDateTime(this.props.invoice.invoiced_date)}</strong> by{' '}
              {this.props.invoice.approver_first_name} {this.props.invoice.approver_last_name}
            </span>
          </h5>
        </div>
      </div>
    );

    return (
      <div className="invoice-panel-table-cont" data-cy="invoice-table">
        <LineItemTable
          shipmentLineItems={this.props.lineItems}
          totalAmount={this.props.lineItemsTotal}
          title={tableTitle}
        />
      </div>
    );
  }
}

InvoiceTable.propTypes = {
  invoice: PropTypes.object.isRequired,
};

const mapStateToProps = (state, ownProps) => {
  return {
    lineItems: selectInvoiceShipmentLineItems(state, ownProps.invoice.id),
    lineItemsTotal: selectTotalFromInvoicedLineItems(state, ownProps.invoice.id) || 0,
  };
};

export default connect(mapStateToProps, null)(InvoiceTable);
