import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import './InvoicePanel.css';

export const PAYMENT_IN_CONFIRMATION = 'IN_CONFIRMATION_FLOW';
export const PAYMENT_IN_PROCESSING = 'IN_PROCESSING_FLOW';
export const PAYMENT_APPROVED = 'IN_APPROVED_FLOW';
export const PAYMENT_FAILED = 'IN_FAILURE_FLOW';

class InvoicePayment extends PureComponent {
  constructor(props) {
    super(props);

    this.state = {
      draftInvoice: false,
    };
  }

  draftInvoice = () => {
    this.setState({ draftInvoice: true });
  };

  cancelPayment = () => {
    this.setState({ draftInvoice: false });
  };

  approvePayment = () => {
    this.props.approvePayment();
    this.cancelPayment();
  };

  render() {
    const status = this.props.createInvoiceStatus;
    let paymentAlert;
    const allowPayment = this.props.allowPayment && !status.isLoading;

    let header = (
      <div className="invoice-panel-header-cont">
        <div className="usa-width-one-half">
          <h5>Unbilled line items</h5>
        </div>
        {allowPayment && (
          <div className="usa-width-one-half align-right">
            <button className="button button-secondary" onClick={this.draftInvoice}>
              Approve Payment
            </button>
          </div>
        )}
      </div>
    );

    if (this.state.draftInvoice) {
      header = (
        <Alert type="warning" heading="Approve payment?">
          <span className="warning--header">Please make sure you've double-checked everything.</span>
          <button className="button usa-button-secondary" onClick={this.cancelPayment}>
            Cancel
          </button>
          <button className="button usa-button-primary" onClick={this.approvePayment}>
            Approve
          </button>
        </Alert>
      );
    }

    if (status.error) {
      paymentAlert = (
        <Alert type="error" heading="Oops, something went wrong!">
          <span className="warning--header">Please try again.</span>
        </Alert>
      );
    } else if (status.isLoading) {
      paymentAlert = (
        <Alert type="loading" heading="Creating invoice">
          <span className="warning--header">Sending information to USBank/Syncada.</span>
        </Alert>
      );
    } else if (status.isSuccess) {
      paymentAlert = (
        <div>
          <Alert type="success" heading="Success!">
            <span className="warning--header">The invoice has been created and will be paid soon.</span>
          </Alert>
        </div>
      );
    }

    return (
      <div>
        {paymentAlert}
        {header}
      </div>
    );
  }
}

InvoicePayment.propTypes = {
  approvePayment: PropTypes.func,
  cancelPayment: PropTypes.func,
  // paymentStatus: PropTypes.string,
  createInvoiceStatus: PropTypes.object,
  allowPayment: PropTypes.bool,
};

export default InvoicePayment;
