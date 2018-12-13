import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import BasicPanel from 'shared/BasicPanel';
import {
  selectUnbilledShipmentLineItems,
  selectTotalFromUnbilledLineItems,
  getAllShipmentLineItems,
  getShipmentLineItemsLabel,
} from 'shared/Entities/modules/shipmentLineItems';
import { selectSortedInvoices, createInvoice, createInvoiceLabel } from 'shared/Entities/modules/invoices';
import { getRequestStatus } from 'shared/Swagger/selectors';
import UnbilledTable from 'shared/Invoice/UnbilledTable';
import InvoiceTable from 'shared/Invoice/InvoiceTable';

import './InvoicePanel.css';

export class InvoicePanel extends PureComponent {
  approvePayment = () => {
    return this.props.createInvoice(createInvoiceLabel, this.props.shipmentId).then(() => {
      return this.props.getAllShipmentLineItems(getShipmentLineItemsLabel, this.props.shipmentId);
    });
  };

  render() {
    // For now we're only allowing one invoice to be generated
    const allowPayments = !this.props.invoices || !this.props.invoices.length;
    return (
      <div className="invoice-panel">
        <BasicPanel title="Invoicing">
          <UnbilledTable
            lineItems={this.props.unbilledShipmentLineItems}
            lineItemsTotal={this.props.unbilledLineItemsTotal}
            cancelPayment={this.props.resetInvoiceFlow}
            approvePayment={this.approvePayment.bind(this)}
            allowPayments={allowPayments}
            createInvoiceStatus={this.props.createInvoiceStatus}
          />

          {this.props.invoices &&
            this.props.invoices.map(invoice => {
              return <InvoiceTable invoice={invoice} key={invoice.id} />;
            })}
        </BasicPanel>
      </div>
    );
  }
}

InvoicePanel.propTypes = {
  unbilledShipmentLineItems: PropTypes.array,
  unbilledLineItemsTotal: PropTypes.number,
  shipmentId: PropTypes.string,
  shipmentStatus: PropTypes.string,
  isShipmentDelivered: PropTypes.bool,
};

const mapStateToProps = (state, ownProps) => {
  const isShipmentDelivered = ownProps.shipmentStatus.toUpperCase() === 'DELIVERED';
  return {
    invoices: isShipmentDelivered ? selectSortedInvoices(state, ownProps.shipmentId) : [],
    unbilledShipmentLineItems: isShipmentDelivered ? selectUnbilledShipmentLineItems(state, ownProps.shipmentId) : [],
    unbilledLineItemsTotal: isShipmentDelivered ? selectTotalFromUnbilledLineItems(state, ownProps.shipmentId) : 0,
    isShipmentDelivered: isShipmentDelivered,
    createInvoiceStatus: getRequestStatus(state, createInvoiceLabel),
  };
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createInvoice, getAllShipmentLineItems }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(InvoicePanel);
