import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';
import { get } from 'lodash';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import BasicPanel from 'shared/BasicPanel';
import {
  selectUnbilledShipmentLineItems,
  selectTotalFromUnbilledLineItems,
  getAllShipmentLineItems,
  getShipmentLineItemsLabel,
} from 'shared/Entities/modules/shipmentLineItems';
import {
  selectSortedInvoices,
  createInvoice,
  createInvoiceLabel,
  getShipmentInvoicesLabel,
  getAllInvoices,
} from 'shared/Entities/modules/invoices';
import UnbilledTable from 'shared/Invoice/UnbilledTable';
import InvoiceTable from 'shared/Invoice/InvoiceTable';
import InvoicePaymentAlert from './InvoicePaymentAlert';
import { isError, isLoading, isSuccess } from 'shared/constants';
import { getLastError } from 'shared/Swagger/selectors';

import './InvoicePanel.css';

export class InvoicePanel extends PureComponent {
  constructor(props) {
    super(props);

    this.state = {
      createInvoiceRequestStatus: null,
    };
  }

  approvePayment = () => {
    this.setState({ createInvoiceRequestStatus: isLoading });
    return this.props
      .createInvoice(createInvoiceLabel, this.props.shipmentId)
      .then(() => {
        this.setState({ createInvoiceRequestStatus: isSuccess });
        return this.props.getAllShipmentLineItems(getShipmentLineItemsLabel, this.props.shipmentId);
      })
      .catch(err => {
        this.setState({ createInvoiceRequestStatus: isError });
        let httpResCode = get(err, 'response.status');
        if (httpResCode === 409) {
          this.props.getAllInvoices(getShipmentInvoicesLabel, this.props.shipmentId);
          return this.props.getAllShipmentLineItems(getShipmentLineItemsLabel, this.props.shipmentId);
        }
      });
  };

  render() {
    // For now we're only allowing one invoice to be generated
    const allowPayments = this.props.allowPayments && (!this.props.invoices || !this.props.invoices.length);
    const hasUnbilled = Boolean(get(this.props, 'unbilledShipmentLineItems.length'));
    const hasInvoices = Boolean(get(this.props, 'invoices.length'));
    return (
      <div className="invoice-panel">
        <BasicPanel title="Invoicing">
          <InvoicePaymentAlert
            createInvoiceStatus={this.state.createInvoiceRequestStatus}
            lastInvoiceError={this.props.lastInvoiceError}
          />

          {hasUnbilled && (
            <UnbilledTable
              lineItems={this.props.unbilledShipmentLineItems}
              lineItemsTotal={this.props.unbilledLineItemsTotal}
              approvePayment={this.approvePayment.bind(this)}
              allowPayments={allowPayments}
              createInvoiceStatus={this.state.createInvoiceRequestStatus}
            />
          )}

          {hasInvoices &&
            this.props.invoices.map(invoice => {
              return <InvoiceTable invoice={invoice} key={invoice.id} />;
            })}

          {!hasUnbilled && !hasInvoices && <span className="empty-content">No line items</span>}
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
  allowPayments: PropTypes.bool,
};

const mapStateToProps = (state, ownProps) => {
  const isShipmentDelivered = ownProps.shipmentStatus.toUpperCase() === 'DELIVERED';
  return {
    invoices: isShipmentDelivered ? selectSortedInvoices(state, ownProps.shipmentId) : [],
    unbilledShipmentLineItems: isShipmentDelivered ? selectUnbilledShipmentLineItems(state, ownProps.shipmentId) : [],
    unbilledLineItemsTotal: isShipmentDelivered ? selectTotalFromUnbilledLineItems(state, ownProps.shipmentId) : 0,
    isShipmentDelivered: isShipmentDelivered,
    lastInvoiceError: getLastError(state, createInvoiceLabel),
  };
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createInvoice, getAllShipmentLineItems, getAllInvoices }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(InvoicePanel);
