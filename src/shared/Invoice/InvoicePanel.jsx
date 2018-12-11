import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';
import { get } from 'lodash';

import './InvoicePanel.css';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { isOfficeSite, isDevelopment } from 'shared/constants.js';

import BasicPanel from 'shared/BasicPanel';
import {
  selectUnbilledShipmentLineItems,
  selectTotalFromUnbilledLineItems,
} from 'shared/Entities/modules/shipmentLineItems';
import { sendHHGInvoice, draftInvoice, resetInvoiceFlow } from 'scenes/Office/ducks';
import InvoiceTable from 'shared/Invoice/InvoiceTable';
import {
  default as InvoicePayment,
  PAYMENT_IN_PROCESSING,
  PAYMENT_IN_CONFIRMATION,
  PAYMENT_FAILED,
  PAYMENT_APPROVED,
} from './InvoicePayment';

export class InvoicePanel extends PureComponent {
  approvePayment = () => {
    return this.props.sendHHGInvoice(this.props.shipmentId);
  };

  render() {
    let invoicingContent = <span className="empty-content">No line items</span>;
    let title = 'Invoicing';

    if (this.props.unbilledShipmentLineItems.length > 0) {
      let tableTitle = <h5>Unbilled line items</h5>;
      let paymentStatus = null;

      //figure out payment status
      if (this.props.isInvoiceInDraft) {
        paymentStatus = PAYMENT_IN_CONFIRMATION;
      } else if (this.props.invoiceProcessing) {
        paymentStatus = PAYMENT_IN_PROCESSING;
      } else if (this.props.invoiceSentSuccessfully) {
        paymentStatus = PAYMENT_APPROVED;
      } else if (this.props.invoiceHasFailed) {
        paymentStatus = PAYMENT_FAILED;
      }

      if (
        isOfficeSite && //user is an office user
        isDevelopment && //only for development env
        (!paymentStatus || paymentStatus === PAYMENT_FAILED)
      ) {
        //payment status is empty or invoice payment has failed
        tableTitle = (
          <div className="invoice-panel-header-cont">
            <div className="usa-width-one-half">
              <h5>Unbilled line items</h5>
            </div>
            <div className="usa-width-one-half align-right">
              <button className="button button-secondary" onClick={this.props.draftInvoice}>
                Approve Payment
              </button>
            </div>
          </div>
        );
      }

      invoicingContent = (
        <div>
          <InvoicePayment
            cancelPayment={this.props.resetInvoiceFlow}
            approvePayment={this.approvePayment}
            isDelivered={this.props.isShipmentDelivered}
            paymentStatus={paymentStatus}
          />
          <InvoiceTable
            shipmentLineItems={this.props.unbilledShipmentLineItems}
            totalAmount={this.props.lineItemsTotal}
            title={tableTitle}
          />
        </div>
      );
      title = <span>Invoicing</span>;
    }

    return (
      <div className="invoice-panel">
        <BasicPanel title={title}>{invoicingContent}</BasicPanel>
      </div>
    );
  }
}

InvoicePanel.propTypes = {
  unbilledShipmentLineItems: PropTypes.array,
  shipmentId: PropTypes.string,
  shipmentStatus: PropTypes.string,
  lineItemsTotal: PropTypes.number,
  onApprovePayment: PropTypes.func,
  canApprove: PropTypes.bool,
  isShipmentDelivered: PropTypes.bool,
};

const mapStateToProps = (state, ownProps) => {
  const isShipmentDelivered = ownProps.shipmentStatus.toUpperCase() === 'DELIVERED';
  return {
    unbilledShipmentLineItems: isShipmentDelivered ? selectUnbilledShipmentLineItems(state, ownProps.shipmentId) : [],
    lineItemsTotal: isShipmentDelivered ? selectTotalFromUnbilledLineItems(state, ownProps.shipmentId) : 0,
    isShipmentDelivered: ownProps.shipmentStatus.toUpperCase() === 'DELIVERED',
    invoiceProcessing: get(state, 'office.hhgInvoiceIsSending'),
    invoiceSentSuccessfully: get(state, 'office.hhgInvoiceHasSendSuccess'),
    invoiceHasFailed: get(state, 'office.hhgInvoiceHasFailure'),
    isInvoiceInDraft: get(state, 'office.hhgInvoiceInDraft'),
  };
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ sendHHGInvoice, draftInvoice, resetInvoiceFlow }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(InvoicePanel);
