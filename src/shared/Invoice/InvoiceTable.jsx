import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import { formatFromBaseQuantity, formatCents } from 'shared/formatters';
import Alert from 'shared/Alert';

import { isOfficeSite, isDevelopment } from 'shared/constants.js';
import './InvoicePanel.css';

const PAYMENT_PENDING = 'IN_PENDING_FLOW';
const PAYMENT_IN_CONFIRMATION = 'IN_CONFIRMATION_FLOW';
const PAYMENT_IN_PROCESSING = 'IN_PROCESSING_FLOW';
const PAYMENT_APPROVED = 'IN_APPROVED_FLOW';
const PAYMENT_FAILED = 'IN_FAILURE_FLOW';

class InvoiceTable extends PureComponent {
  state = {
    paymentStatus: this.props.paymentStatus,
  };

  static defaultProps = {
    paymentStatus: PAYMENT_PENDING,
  };

  approvePayment = () => {
    this.setState({ paymentStatus: PAYMENT_IN_CONFIRMATION });
  };

  cancelPayment = () => {
    this.setState({ paymentStatus: PAYMENT_PENDING });
  };

  submitPayment = () => {
    //dispatch action to submit invoice to GEX
    this.props.approvePayment();
    this.setState({ paymentStatus: PAYMENT_IN_PROCESSING });
    setTimeout(this.invoiceSuccess, 5000);
  };

  invoiceSuccess = () => {
    this.setState({ paymentStatus: PAYMENT_APPROVED });
  };

  invoiceFail = () => {
    this.setState({ paymentStatus: PAYMENT_FAILED });
  };

  render() {
    let paymentContainer = null;
    let isDelivered = this.props.shipmentStatus.toUpperCase() === 'DELIVERED';

    //calculate what payment status view to display
    switch (this.state.paymentStatus) {
      default:
        paymentContainer = null;
        break;
      case PAYMENT_PENDING:
        if (isOfficeSite && isDelivered) {
          paymentContainer = (
            <div className="invoice-panel-header-cont">
              <div className="usa-width-one-half">
                <h5>Unbilled line items</h5>
              </div>
              <div className="usa-width-one-half align-right">
                <button
                  className="button button-secondary"
                  disabled={!this.props.canApprove || !isDevelopment}
                  onClick={this.approvePayment}
                >
                  Approve Payment
                </button>
              </div>
            </div>
          );
        }
        break;
      case PAYMENT_IN_CONFIRMATION:
        if (isOfficeSite && isDelivered) {
          paymentContainer = (
            <div>
              <Alert type="warning" heading="Approve payment?">
                <span className="warning--header">Please make sure you've double-checked everything.</span>
                <button className="button usa-button-secondary" onClick={this.cancelPayment}>
                  Cancel
                </button>
                <button className="button usa-button-primary" onClick={this.submitPayment}>
                  {' '}
                  Approve
                </button>
              </Alert>
              <div className="usa-grid-full invoice-panel-header-cont">
                <div className="usa-width-one-half">
                  <h5>Unbilled line items</h5>
                </div>
              </div>
            </div>
          );
        }
        break;
      case PAYMENT_IN_PROCESSING:
        if (isOfficeSite && isDelivered) {
          paymentContainer = (
            <div>
              <Alert type="loading" heading="Creating invoice">
                <span className="warning--header">Sending information to USBank/Syncada.</span>
              </Alert>
              <div className="usa-grid-full invoice-panel-header-cont">
                <div className="usa-width-one-half">
                  <h5>Unbilled line items</h5>
                </div>
              </div>
            </div>
          );
        }
        break;
      case PAYMENT_APPROVED:
        if (isOfficeSite && isDelivered) {
          paymentContainer = (
            <div>
              <Alert type="success" heading="Success!">
                <span className="warning--header">The invoice has been created and will be paid soon.</span>
              </Alert>
            </div>
          );
        }
        break;
      case PAYMENT_FAILED:
        if (isOfficeSite && isDelivered) {
          paymentContainer = (
            <div>
              <Alert type="error" heading="Oops, something went wrong!">
                <span className="warning--header">Please try again.</span>
              </Alert>
              <div className="invoice-panel-header-cont">
                <div className="usa-width-one-half">
                  <h5>Unbilled line items</h5>
                </div>
                <div className="usa-width-one-half align-right">
                  <button
                    className="button button-secondary"
                    disabled={!this.props.canApprove || !isDevelopment}
                    onClick={this.approvePayment}
                  >
                    Approve Payment
                  </button>
                </div>
              </div>
            </div>
          );
        }
        break;
    }
    return (
      <div className="invoice-panel-table-cont">
        {paymentContainer}
        <table cellSpacing={0}>
          <tbody>
            <tr>
              <th>Code</th>
              <th>Item</th>
              <th>Loc</th>
              <th>Base quantity</th>
              <th>Inv amt</th>
            </tr>
            {this.props.shipmentLineItems.map(item => {
              return (
                <tr key={item.id}>
                  <td>{item.tariff400ng_item.code}</td>
                  <td>{item.tariff400ng_item.item}</td>
                  <td>{item.location[0]}</td>
                  <td>{formatFromBaseQuantity(item.quantity_1)}</td>
                  <td>${formatCents(item.amount_cents)}</td>
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
  shipmentStatus: PropTypes.string,
  totalAmount: PropTypes.number,
  approvePayment: PropTypes.func,
  canApprove: PropTypes.bool,
};

export default InvoiceTable;
