import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import './InvoicePanel.css';

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
    let paymentAlert;
    const status = this.props.createInvoiceStatus;
    const allowPayments = this.props.allowPayments && !status.isLoading;

    let header = (
      <div className="invoice-panel-header-cont">
        <div className="usa-width-one-half">
          <h5>Unbilled line items</h5>
        </div>
        {allowPayments && (
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
  createInvoiceStatus: PropTypes.object,
  allowPayments: PropTypes.bool,
};

export default InvoicePayment;
