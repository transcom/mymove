import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import './InvoicePanel.css';

export const PAYMENT_IN_CONFIRMATION = 'IN_CONFIRMATION_FLOW';
export const PAYMENT_IN_PROCESSING = 'IN_PROCESSING_FLOW';
export const PAYMENT_APPROVED = 'IN_APPROVED_FLOW';
export const PAYMENT_FAILED = 'IN_FAILURE_FLOW';

class InvoicePayment extends PureComponent {
  render() {
    let paymentContainer = null;

    //calculate what payment status view to display
    switch (this.props.paymentStatus) {
      case PAYMENT_IN_CONFIRMATION:
        paymentContainer = (
          <Alert type="warning" heading="Approve payment?">
            <span className="warning--header">Please make sure you've double-checked everything.</span>
            <button className="button usa-button-secondary" onClick={this.props.cancelPayment}>
              Cancel
            </button>
            <button className="button usa-button-primary" onClick={this.props.approvePayment}>
              Approve
            </button>
          </Alert>
        );
        break;
      case PAYMENT_IN_PROCESSING:
        paymentContainer = (
          <Alert type="loading" heading="Creating invoice">
            <span className="warning--header">Sending information to USBank/Syncada.</span>
          </Alert>
        );
        break;
      case PAYMENT_APPROVED:
        paymentContainer = (
          <div>
            <Alert type="success" heading="Success!">
              <span className="warning--header">The invoice has been created and will be paid soon.</span>
            </Alert>
          </div>
        );
        break;
      case PAYMENT_FAILED:
        paymentContainer = (
          <Alert type="error" heading="Oops, something went wrong!">
            <span className="warning--header">Please try again.</span>
          </Alert>
        );
        break;
      default:
        // unknown status
        paymentContainer = null;
        break;
    }
    return paymentContainer;
  }
}

InvoicePayment.propTypes = {
  approvePayment: PropTypes.func,
  cancelPayment: PropTypes.func,
  paymentStatus: PropTypes.string,
};

export default InvoicePayment;
