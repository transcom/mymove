import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';

import { isOfficeSite, isDevelopment } from 'shared/constants.js';
import './InvoicePanel.css';

const PAYMENT_PENDING = 'IN_PENDING_FLOW';
const PAYMENT_IN_CONFIRMATION = 'IN_CONFIRMATION_FLOW';
const PAYMENT_IN_PROCESSING = 'IN_PROCESSING_FLOW';
const PAYMENT_APPROVED = 'IN_APPROVED_FLOW';
const PAYMENT_FAILED = 'IN_FAILURE_FLOW';

class InvoicePayment extends PureComponent {
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
    this.setState({ paymentStatus: PAYMENT_IN_PROCESSING });

    //dispatch action to submit invoice to GEX
    this.props.approvePayment().then(status => {
      //this is a temp workaround until invoice table gets refactored
      //and invoice status starts coming from redux store instead of being in-state
      if (status.type === 'SEND_HHG_INVOICE_SUCCESS') {
        this.invoiceSuccess();
      } else {
        this.invoiceFail();
      }
    });
  };

  invoiceSuccess = () => {
    this.setState({ paymentStatus: PAYMENT_APPROVED });
  };

  invoiceFail = () => {
    this.setState({ paymentStatus: PAYMENT_FAILED });
  };

  render() {
    let paymentContainer = null;

    //calculate what payment status view to display
    switch (this.state.paymentStatus) {
      default:
        paymentContainer = null;
        break;
      case PAYMENT_PENDING:
        if (isOfficeSite && this.props.isDelivered) {
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
        if (isOfficeSite && this.props.isDelivered) {
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
              <div className="invoice-panel-header-cont">
                <div className="usa-width-one-half">
                  <h5>Unbilled line items</h5>
                </div>
              </div>
            </div>
          );
        }
        break;
      case PAYMENT_IN_PROCESSING:
        if (isOfficeSite && this.props.isDelivered) {
          paymentContainer = (
            <div>
              <Alert type="loading" heading="Creating invoice">
                <span className="warning--header">Sending information to USBank/Syncada.</span>
              </Alert>
              <div className="invoice-panel-header-cont">
                <div className="usa-width-one-half">
                  <h5>Unbilled line items</h5>
                </div>
              </div>
            </div>
          );
        }
        break;
      case PAYMENT_APPROVED:
        if (isOfficeSite && this.props.isDelivered) {
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
        if (isOfficeSite && this.props.isDelivered) {
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
    return { paymentContainer };
  }
}

InvoicePayment.propTypes = {
  approvePayment: PropTypes.func,
  canApprove: PropTypes.bool,
  isDelivered: PropTypes.bool,
};

export default InvoicePayment;
